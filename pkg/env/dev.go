package env

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/go-autorest/autorest"
	"github.com/jongio/azidext/go/azidext"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer-aro-wrapper/pkg/util/version"
)

type dev struct {
	*prod
}

func newDev(ctx context.Context, log *logrus.Entry) (Interface, error) {
	for _, key := range []string{
		"ARO_PROXY_HOSTNAME",
	} {
		if _, found := os.LookupEnv(key); !found {
			return nil, fmt.Errorf("environment variable %q unset", key)
		}
	}

	d := &dev{}

	var err error
	d.prod, err = newProd(ctx, log)
	if err != nil {
		return nil, err
	}

	for _, feature := range []Feature{
		FeatureDisableDenyAssignments,
		FeatureDisableSignedCertificates,
		FeatureRequireD2sV3Workers,
		FeatureDisableReadinessDelay,
	} {
		d.features[feature] = true
	}

	d.prod.clusterGenevaLoggingAccount = version.DevClusterGenevaLoggingAccount
	d.prod.clusterGenevaLoggingConfigVersion = version.DevClusterGenevaLoggingConfigVersion
	d.prod.clusterGenevaLoggingEnvironment = version.DevGenevaLoggingEnvironment
	d.prod.clusterGenevaLoggingNamespace = version.DevClusterGenevaLoggingNamespace

	return d, nil
}

func (d *dev) FPAuthorizer(tenantID string, scopes ...string) (autorest.Authorizer, error) {
	tokenCredential, err := d.FPNewClientCertificateCredential(tenantID)
	if err != nil {
		return nil, err
	}

	return azidext.NewTokenCredentialAdapter(tokenCredential, scopes), nil
}

func (d *dev) FPNewClientCertificateCredential(tenantID string) (*azidentity.ClientCertificateCredential, error) {
	fpPrivateKey, fpCertificates := d.fpCertificateRefresher.GetCertificates()

	options := d.Environment().ClientCertificateCredentialOptions()
	credential, err := azidentity.NewClientCertificateCredential(tenantID, d.fpClientID, fpCertificates, fpPrivateKey, options)
	if err != nil {
		return nil, err
	}
	return credential, nil
}
