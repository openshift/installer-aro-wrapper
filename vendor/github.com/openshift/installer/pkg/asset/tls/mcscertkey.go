package tls

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"net"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/templates/content/bootkube"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
	baremetaltypes "github.com/openshift/installer/pkg/types/baremetal"
	nutanixtypes "github.com/openshift/installer/pkg/types/nutanix"
	openstacktypes "github.com/openshift/installer/pkg/types/openstack"
	ovirttypes "github.com/openshift/installer/pkg/types/ovirt"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

// MCSCertKey is the asset that generates the MCS key/cert pair.
type MCSCertKey struct {
	SignedCertKey
}

var _ asset.Asset = (*MCSCertKey)(nil)

// Dependencies returns the dependency of the the cert/key pair, which includes
// the parent CA, and install config if it depends on the install config for
// DNS names, etc.
func (a *MCSCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{
		&RootCA{},
		&installconfig.InstallConfig{},
		&bootkube.ARODNSConfig{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *MCSCertKey) Generate(dependencies asset.Parents) error {
	ca := &RootCA{}
	installConfig := &installconfig.InstallConfig{}
	aroDNSConfig := &bootkube.ARODNSConfig{}
	dependencies.Get(ca, installConfig, aroDNSConfig)

	hostname := internalAPIAddress(installConfig.Config)

	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: "system:machine-config-server"},
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		Validity:     ValidityTenYears,
	}

	var vips []string
	switch installConfig.Config.Platform.Name() {
	case azuretypes.Name:
		if installConfig.Config.Azure.IsARO() {
			vips = []string{aroDNSConfig.APIIntIP}
		}
	case baremetaltypes.Name:
		vips = installConfig.Config.BareMetal.APIVIPs
	case nutanixtypes.Name:
		vips = installConfig.Config.Nutanix.APIVIPs
	case openstacktypes.Name:
		vips = installConfig.Config.OpenStack.APIVIPs
	case ovirttypes.Name:
		vips = installConfig.Config.Ovirt.APIVIPs
	case vspheretypes.Name:
		vips = installConfig.Config.VSphere.APIVIPs
	}

	cfg.IPAddresses = []net.IP{}
	cfg.DNSNames = []string{hostname}
	for _, vip := range vips {
		cfg.IPAddresses = append(cfg.IPAddresses, net.ParseIP(vip))
		cfg.DNSNames = append(cfg.DNSNames, vip)
	}

	return a.SignedCertKey.Generate(cfg, ca, "machine-config-server", DoNotAppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *MCSCertKey) Name() string {
	return "Certificate (mcs)"
}
