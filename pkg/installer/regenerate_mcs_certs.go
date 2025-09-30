package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"fmt"
	"net"
	"strings"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	corev1 "k8s.io/api/core/v1"

	"sigs.k8s.io/yaml"

	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"

	"github.com/openshift/installer-aro-wrapper/pkg/cluster/graph"
	"github.com/openshift/installer-aro-wrapper/pkg/installer/dnsmasq"
)

const (
	mcsCertKeyFilepath = "/opt/openshift/manifests/machine-config-server-tls-secret.yaml"
	mcsKeyFile         = "/opt/openshift/tls/machine-config-server.key"
	mcsCertFile        = "/opt/openshift/tls/machine-config-server.crt"
	// header is the string that precedes the encoded data in the ignition data.
	// The data must be replaced before decoding the string, and the string must be
	// prepended to the encoded data.
	header = "data:text/plain;charset=utf-8;base64,"
)

// RegenerateSignedCertKey regenerates a cert/key pair signed by the specified parent CA.
// It does not write the cert/key pair to an asset file.
func regenerateSignedCertKey(
	cfg *tls.CertCfg,
	parentCA tls.CertKeyInterface,
	appendParent tls.AppendParentChoice,
) ([]byte, []byte, error) {
	var key *rsa.PrivateKey
	var crt *x509.Certificate
	var err error

	caKey, err := tls.PemToPrivateKey(parentCA.Key())
	if err != nil {
		logrus.Debugf("Failed to parse RSA private key: %s", err)
		return nil, nil, errors.Wrap(err, "failed to parse rsa private key")
	}

	caCert, err := tls.PemToCertificate(parentCA.Cert())
	if err != nil {
		logrus.Debugf("Failed to parse x509 certificate: %s", err)
		return nil, nil, errors.Wrap(err, "failed to parse x509 certificate")
	}

	key, crt, err = tls.GenerateSignedCertificate(caKey, caCert, cfg)
	if err != nil {
		logrus.Debugf("Failed to generate signed cert/key pair: %s", err)
		return nil, nil, errors.Wrap(err, "failed to generate signed cert/key pair")
	}

	keyRaw := tls.PrivateKeyToPem(key)
	certRaw := tls.CertToPem(crt)

	if appendParent {
		certRaw = bytes.Join([][]byte{certRaw, tls.CertToPem(caCert)}, []byte("\n"))
	}

	return keyRaw, certRaw, nil
}

// RegenerateMCSCertKey generates the cert/key pair based on input values.
func regenerateMCSCertKey(ic *installconfig.InstallConfig, ca *tls.RootCA, privateLBIP string) ([]byte, []byte, error) {
	hostname := fmt.Sprintf("api-int.%s", ic.Config.ClusterDomain())
	cfg := &tls.CertCfg{
		Subject:      pkix.Name{CommonName: "system:machine-config-server"},
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		Validity:     tls.ValidityTenYears(),
	}
	cfg.IPAddresses = []net.IP{}
	cfg.DNSNames = []string{hostname}
	cfg.IPAddresses = append(cfg.IPAddresses, net.ParseIP(privateLBIP))
	cfg.DNSNames = append(cfg.DNSNames, privateLBIP)
	return regenerateSignedCertKey(cfg, ca, tls.DoNotAppendParent)
}

func updateMCSCertKey(g graph.Graph, ic *installconfig.InstallConfig, localdnsConfig *dnsmasq.DNSConfig) error {
	if len(localdnsConfig.APIIntIP) > 0 {
		rootCA := g.Get(&tls.RootCA{}).(*tls.RootCA)
		keyRaw, certRaw, err := regenerateMCSCertKey(ic, rootCA, localdnsConfig.APIIntIP)
		if err != nil {
			return fmt.Errorf("failed to regenerate MCS Cert and Key: %w", err)
		}

		// Manipulating the bootstrap ignition
		config := &igntypes.Config{}
		bootstrap := g.Get(&bootstrap.Bootstrap{}).(*bootstrap.Bootstrap)
		config = bootstrap.Config

		for i, fileData := range config.Storage.Files {
			switch fileData.Path {
			case mcsCertKeyFilepath:
				contents := strings.Split(*config.Storage.Files[i].Contents.Source, ",")

				rawDecodedText, err := base64.StdEncoding.DecodeString(contents[1])
				if err != nil {
					return fmt.Errorf("failed to decode contents of ignition file %s: %w", mcsCertKeyFilepath, err)
				}
				mcsSecret := &corev1.Secret{}
				if err := yaml.Unmarshal(rawDecodedText, mcsSecret); err != nil {
					return fmt.Errorf("failed to unmarshal MCSCertKey within ignition: %w", err)
				}
				mcsSecret.Data[corev1.TLSCertKey] = certRaw
				mcsSecret.Data[corev1.TLSPrivateKeyKey] = keyRaw
				// convert the mcsSecret back to an encoded string
				mcsSecretContents, err := yaml.Marshal(mcsSecret)
				if err != nil {
					return fmt.Errorf("failed to marshal MCS Secret: %w", err)
				}
				encoded := fmt.Sprintf("%s%s", header, base64.StdEncoding.EncodeToString(mcsSecretContents))
				// replace the contents with the edited information
				config.Storage.Files[i].Contents.Source = &encoded

				logrus.Debugf("Updated MCSCertKey file %s with new cert and key", mcsCertKeyFilepath)
			case mcsKeyFile:
				encoded := fmt.Sprintf("%s%s", header, base64.StdEncoding.EncodeToString(keyRaw))
				// replace the contents with the edited information
				config.Storage.Files[i].Contents.Source = &encoded
				logrus.Debugf("Updated MCSKey file %s with new key", mcsKeyFile)
			case mcsCertFile:
				encoded := fmt.Sprintf("%s%s", header, base64.StdEncoding.EncodeToString(certRaw))
				// replace the contents with the edited information
				config.Storage.Files[i].Contents.Source = &encoded
				logrus.Debugf("Updated MCSCert file %s with new cert", mcsCertFile)
			}
		}
	}
	return nil
}
