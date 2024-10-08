package env

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	mgmtcompute "github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/jongio/azidext/go/azidext"
	icazure "github.com/openshift/installer/pkg/asset/installconfig/azure"
	"github.com/sirupsen/logrus"

	"github.com/openshift/ARO-Installer/pkg/proxy"
	"github.com/openshift/ARO-Installer/pkg/util/azureclient/mgmt/compute"
	"github.com/openshift/ARO-Installer/pkg/util/clientauthorizer"
	"github.com/openshift/ARO-Installer/pkg/util/computeskus"
	"github.com/openshift/ARO-Installer/pkg/util/keyvault"
	"github.com/openshift/ARO-Installer/pkg/util/version"
)

type prod struct {
	Core
	proxy.Dialer

	armClientAuthorizer   clientauthorizer.ClientAuthorizer
	adminClientAuthorizer clientauthorizer.ClientAuthorizer

	acrDomain string
	vmskus    map[string]*mgmtcompute.ResourceSku

	fpCertificateRefresher CertificateRefresher
	fpClientID             string

	clusterKeyvault keyvault.Manager
	serviceKeyvault keyvault.Manager

	clusterGenevaLoggingCertificate   *x509.Certificate
	clusterGenevaLoggingPrivateKey    *rsa.PrivateKey
	clusterGenevaLoggingAccount       string
	clusterGenevaLoggingConfigVersion string
	clusterGenevaLoggingEnvironment   string
	clusterGenevaLoggingNamespace     string

	gatewayDomains []string

	log *logrus.Entry

	features map[Feature]bool
}

func newProd(ctx context.Context, log *logrus.Entry) (*prod, error) {
	for _, key := range []string{
		"ARO_AZURE_FP_CLIENT_ID",
		"ARO_DOMAIN_NAME",
	} {
		if _, found := os.LookupEnv(key); !found {
			return nil, fmt.Errorf("environment variable %q unset", key)
		}
	}

	if !IsLocalDevelopmentMode() {
		for _, key := range []string{
			"ARO_CLUSTER_MDSD_CONFIG_VERSION",
			"ARO_CLUSTER_MDSD_ACCOUNT",
			"ARO_GATEWAY_DOMAINS",
			"ARO_GATEWAY_RESOURCEGROUP",
			"ARO_MDSD_ENVIRONMENT",
			"ARO_CLUSTER_MDSD_NAMESPACE",
			"ARO_ACR_RESOURCE_ID",
		} {
			if _, found := os.LookupEnv(key); !found {
				return nil, fmt.Errorf("environment variable %q unset", key)
			}
		}
	}

	core, err := NewCore(ctx, log)
	if err != nil {
		return nil, err
	}

	dialer, err := proxy.NewDialer(core.IsLocalDevelopmentMode())
	if err != nil {
		return nil, err
	}

	p := &prod{
		Core:   core,
		Dialer: dialer,

		fpClientID: os.Getenv("ARO_AZURE_FP_CLIENT_ID"),

		clusterGenevaLoggingAccount:       os.Getenv("ARO_CLUSTER_MDSD_ACCOUNT"),
		clusterGenevaLoggingConfigVersion: os.Getenv("ARO_CLUSTER_MDSD_CONFIG_VERSION"),
		clusterGenevaLoggingEnvironment:   os.Getenv("ARO_MDSD_ENVIRONMENT"),
		clusterGenevaLoggingNamespace:     os.Getenv("ARO_CLUSTER_MDSD_NAMESPACE"),

		log: log,

		features: map[Feature]bool{},
	}

	features := os.Getenv("ARO_RP_FEATURES")
	if features != "" {
		for _, feature := range strings.Split(features, ",") {
			f, err := FeatureString("Feature" + feature)
			if err != nil {
				return nil, err
			}

			p.features[f] = true
		}
	}

	msiKVAuthorizer, err := p.NewMSIAuthorizer(MSIContextRP, p.Environment().KeyVaultScope)
	if err != nil {
		return nil, err
	}

	serviceKeyvaultURI, err := keyvault.URI(p, ServiceKeyvaultSuffix)
	if err != nil {
		return nil, err
	}

	p.serviceKeyvault = keyvault.NewManager(msiKVAuthorizer, serviceKeyvaultURI)

	p.fpCertificateRefresher = newCertificateRefresher(log, 1*time.Hour, p.serviceKeyvault, RPFirstPartySecretName)
	err = p.fpCertificateRefresher.Start(ctx)
	if err != nil {
		return nil, err
	}

	localFPKVAuthorizer, err := p.FPAuthorizer(p.TenantID(), p.Environment().KeyVaultScope)
	if err != nil {
		return nil, err
	}

	clusterKeyvaultURI, err := keyvault.URI(p, ClusterKeyvaultSuffix)
	if err != nil {
		return nil, err
	}

	p.clusterKeyvault = keyvault.NewManager(localFPKVAuthorizer, clusterKeyvaultURI)

	clusterGenevaLoggingPrivateKey, clusterGenevaLoggingCertificates, err := p.serviceKeyvault.GetCertificateSecret(ctx, ClusterLoggingSecretName)
	if err != nil {
		return nil, err
	}

	p.clusterGenevaLoggingPrivateKey = clusterGenevaLoggingPrivateKey
	p.clusterGenevaLoggingCertificate = clusterGenevaLoggingCertificates[0]

	localFPAuthorizer, err := p.FPAuthorizer(p.TenantID(), p.Environment().ResourceManagerScope)
	if err != nil {
		return nil, err
	}

	resourceSkusClient := compute.NewResourceSkusClient(p.Environment(), p.SubscriptionID(), localFPAuthorizer)
	err = p.populateVMSkus(ctx, resourceSkusClient)
	if err != nil {
		return nil, err
	}

	var acrDataDomain string
	if p.ACRResourceID() != "" { // TODO: ugh!
		acrResource, err := azure.ParseResourceID(p.ACRResourceID())
		if err != nil {
			return nil, err
		}
		p.acrDomain = acrResource.ResourceName + "." + p.Environment().ContainerRegistryDNSSuffix
		acrDataDomain = acrResource.ResourceName + "." + p.Location() + ".data." + p.Environment().ContainerRegistryDNSSuffix
	} else {
		p.acrDomain = "arointsvc." + azure.PublicCloud.ContainerRegistryDNSSuffix                             // TODO: make cloud aware once this is set up for US Gov Cloud
		acrDataDomain = "arointsvc." + p.Location() + ".data." + azure.PublicCloud.ContainerRegistryDNSSuffix // TODO: make cloud aware once this is set up for US Gov Cloud
	}

	if !p.IsLocalDevelopmentMode() {
		gatewayDomains := os.Getenv("ARO_GATEWAY_DOMAINS")
		if gatewayDomains != "" {
			p.gatewayDomains = strings.Split(gatewayDomains, ",")
		}

		for _, rawurl := range []string{
			p.Environment().ActiveDirectoryEndpoint,
			p.Environment().ResourceManagerEndpoint,
		} {
			u, err := url.Parse(rawurl)
			if err != nil {
				return nil, err
			}

			p.gatewayDomains = append(p.gatewayDomains, u.Hostname())
		}

		p.gatewayDomains = append(p.gatewayDomains, p.acrDomain, acrDataDomain)
	}

	return p, nil
}

func (p *prod) InitializeAuthorizers() error {
	if !p.FeatureIsSet(FeatureEnableDevelopmentAuthorizer) {
		p.armClientAuthorizer = clientauthorizer.NewARM(p.log, p.Core)
	} else {
		armClientAuthorizer, err := clientauthorizer.NewSubjectNameAndIssuer(
			p.log,
			"/etc/aro-rp/arm-ca-bundle.pem",
			os.Getenv("ARO_ARM_API_CLIENT_CERT_COMMON_NAME"),
		)
		if err != nil {
			return err
		}

		p.armClientAuthorizer = armClientAuthorizer
	}

	adminClientAuthorizer, err := clientauthorizer.NewSubjectNameAndIssuer(
		p.log,
		"/etc/aro-rp/admin-ca-bundle.pem",
		os.Getenv("ARO_ADMIN_API_CLIENT_CERT_COMMON_NAME"),
	)
	if err != nil {
		return err
	}

	p.adminClientAuthorizer = adminClientAuthorizer
	return nil
}

func (p *prod) ArmClientAuthorizer() clientauthorizer.ClientAuthorizer {
	return p.armClientAuthorizer
}

func (p *prod) AdminClientAuthorizer() clientauthorizer.ClientAuthorizer {
	return p.adminClientAuthorizer
}

func (p *prod) ACRResourceID() string {
	return os.Getenv("ARO_ACR_RESOURCE_ID")
}

func (p *prod) ACRDomain() string {
	return p.acrDomain
}

func (p *prod) AROOperatorImage() string {
	return fmt.Sprintf("%s/aro:%s", p.acrDomain, version.GitCommit)
}

func (p *prod) populateVMSkus(ctx context.Context, resourceSkusClient compute.ResourceSkusClient) error {
	// Filtering is poorly documented, but currently (API version 2019-04-01)
	// it seems that the API returns all SKUs without a filter and with invalid
	// value in the filter.
	// Filtering gives significant optimisation: at the moment of writing,
	// we get ~1.2M response in eastus vs ~37M unfiltered (467 items vs 16618).
	filter := fmt.Sprintf("location eq '%s'", p.Location())
	skus, err := resourceSkusClient.List(ctx, filter)
	if err != nil {
		return err
	}

	p.vmskus = computeskus.FilterVMSizes(skus, p.Location())

	return nil
}

func (p *prod) ClusterGenevaLoggingAccount() string {
	return p.clusterGenevaLoggingAccount
}

func (p *prod) ClusterGenevaLoggingConfigVersion() string {
	return p.clusterGenevaLoggingConfigVersion
}

func (p *prod) ClusterGenevaLoggingEnvironment() string {
	return p.clusterGenevaLoggingEnvironment
}

func (p *prod) ClusterGenevaLoggingNamespace() string {
	return p.clusterGenevaLoggingNamespace
}

func (p *prod) ClusterGenevaLoggingSecret() (*rsa.PrivateKey, *x509.Certificate) {
	return p.clusterGenevaLoggingPrivateKey, p.clusterGenevaLoggingCertificate
}

func (p *prod) ClusterKeyvault() keyvault.Manager {
	return p.clusterKeyvault
}

func (p *prod) Domain() string {
	return os.Getenv("ARO_DOMAIN_NAME")
}

func (p *prod) FeatureIsSet(f Feature) bool {
	return p.features[f]
}

func (p *prod) FPAuthorizer(tenantID string, scopes ...string) (autorest.Authorizer, error) {
	tokenCredential, err := p.FPNewClientCertificateCredential(tenantID)
	if err != nil {
		return nil, err
	}

	return azidext.NewTokenCredentialAdapter(tokenCredential, scopes), nil
}

func (p *prod) FPNewClientCertificateCredential(tenantID string) (*azidentity.ClientCertificateCredential, error) {
	fpPrivateKey, fpCertificates := p.fpCertificateRefresher.GetCertificates()
	options := p.Environment().ClientCertificateCredentialOptions()
	credential, err := azidentity.NewClientCertificateCredential(tenantID, p.fpClientID, fpCertificates, fpPrivateKey, options)
	if err != nil {
		return nil, err
	}

	return credential, nil
}

func (p *prod) FPClientID() string {
	return p.fpClientID
}

func (p *prod) Listen() (net.Listener, error) {
	return net.Listen("tcp", ":8443")
}

func (p *prod) GatewayDomains() []string {
	gatewayDomains := make([]string, len(p.gatewayDomains))

	copy(gatewayDomains, p.gatewayDomains)

	return gatewayDomains
}

func (p *prod) GatewayResourceGroup() string {
	return os.Getenv("ARO_GATEWAY_RESOURCEGROUP")
}

func (p *prod) ServiceKeyvault() keyvault.Manager {
	return p.serviceKeyvault
}

func (p *prod) VMSku(vmSize string) (*mgmtcompute.ResourceSku, error) {
	vmsku, found := p.vmskus[vmSize]
	if !found {
		return nil, fmt.Errorf("sku information not found for vm size %q", vmSize)
	}
	return vmsku, nil
}

func (p *prod) FPNewInstallConfigClientCertificateCredential(tenantID, subscriptionID string) (*icazure.Credentials, error) {
	fpPrivateKey, fpCertificates := p.fpCertificateRefresher.GetCertificates()

	clientCertificateFile, err := os.CreateTemp("/tmp", "fpClientCertificate-*.pem")
	if err != nil {
		return nil, err
	}

	defer clientCertificateFile.Close()

	for _, cer := range fpCertificates {
		pemType := "UNKNOWN"
		var rawBytes []byte

		switch t := any(cer).(type) {
		case *x509.Certificate:
			pemType = "CERTIFICATE"
			rawBytes = t.Raw
		case *x509.CertificateRequest:
			pemType = "CERTIFICATE REQUEST"
			rawBytes = t.Raw
		case *rsa.PrivateKey:
			pemType = "RSA PRIVATE KEY"
			rawBytes = x509.MarshalPKCS1PrivateKey(t)
		case *rsa.PublicKey:
			pemType = "RSA PUBLIC KEY"
			rawBytes, err = x509.MarshalPKIXPublicKey(t)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unable to identify %v", t)
		}

		pem.Encode(clientCertificateFile, &pem.Block{Type: pemType, Bytes: rawBytes})
	}

	pem.Encode(clientCertificateFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(fpPrivateKey)})

	return &icazure.Credentials{
		TenantID:              tenantID,
		SubscriptionID:        subscriptionID,
		ClientID:              p.FPClientID(),
		ClientCertificatePath: clientCertificateFile.Name(),
	}, nil
}
