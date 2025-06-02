package env

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	mgmtcompute "github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/Azure/go-autorest/autorest"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer-aro-wrapper/pkg/proxy"
	"github.com/openshift/installer-aro-wrapper/pkg/util/instancemetadata"
	"github.com/openshift/installer-aro-wrapper/pkg/util/keyvault"
)

type Feature int

// At least to start with, features are intended to be used so that the
// production default is not set (in production RP_FEATURES is unset).
const (
	FeatureDisableDenyAssignments Feature = iota
	FeatureDisableSignedCertificates
	FeatureEnableDevelopmentAuthorizer
	FeatureRequireD2sV3Workers
	FeatureDisableReadinessDelay
)

const (
	RPDevARMSecretName               = "dev-arm"
	RPFirstPartySecretName           = "rp-firstparty"
	RPServerSecretName               = "rp-server"
	ClusterLoggingSecretName         = "cluster-mdsd"
	EncryptionSecretName             = "encryption-key"
	EncryptionSecretV2Name           = "encryption-key-v2"
	FrontendEncryptionSecretName     = "fe-encryption-key"
	FrontendEncryptionSecretV2Name   = "fe-encryption-key-v2"
	DBTokenServerSecretName          = "dbtoken-server"
	PortalServerSecretName           = "portal-server"
	PortalServerClientSecretName     = "portal-client"
	PortalServerSessionKeySecretName = "portal-session-key"
	PortalServerSSHKeySecretName     = "portal-sshkey"
	ClusterKeyvaultSuffix            = "-cls"
	DBTokenKeyvaultSuffix            = "-dbt"
	GatewayKeyvaultSuffix            = "-gwy"
	PortalKeyvaultSuffix             = "-por"
	ServiceKeyvaultSuffix            = "-svc"
	RPPrivateEndpointPrefix          = "rp-pe-"
)

type Interface interface {
	IsLocalDevelopmentMode() bool
	NewMSIAuthorizer(MSIContext, ...string) (autorest.Authorizer, error)
	instancemetadata.InstanceMetadata
	proxy.Dialer

	ClusterGenevaLoggingAccount() string
	ClusterGenevaLoggingConfigVersion() string
	ClusterGenevaLoggingEnvironment() string
	ClusterGenevaLoggingNamespace() string
	ClusterGenevaLoggingSecret() (*rsa.PrivateKey, *x509.Certificate)
	Domain() string
	FeatureIsSet(Feature) bool
	FPAuthorizer(string, ...string) (autorest.Authorizer, error)
	FPCertificates() (*rsa.PrivateKey, []*x509.Certificate)
	FPNewClientCertificateCredential(string) (*azidentity.ClientCertificateCredential, error)
	FPClientID() string
	GatewayDomains() []string
	ServiceKeyvault() keyvault.Manager
	ACRDomain() string

	// VMSku returns SKU for a given vm size. Note that this
	// returns a pointer to partly populated object.
	VMSku(vmSize string) (*mgmtcompute.ResourceSku, error)
}

func NewEnv(ctx context.Context, log *logrus.Entry) (Interface, error) {
	if IsLocalDevelopmentMode() {
		return newDev(ctx, log)
	}

	return newProd(ctx, log)
}

func IsLocalDevelopmentMode() bool {
	return strings.EqualFold(os.Getenv("ARO_RP_MODE"), "development")
}
