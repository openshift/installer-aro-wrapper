package bootstraplogging

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/openshift/installer/pkg/asset"

	"github.com/openshift/installer-aro-wrapper/pkg/api"
	"github.com/openshift/installer-aro-wrapper/pkg/env"
	utiltls "github.com/openshift/installer-aro-wrapper/pkg/util/tls"
	"github.com/openshift/installer-aro-wrapper/pkg/util/version"
)

// Config generates the bootstraplogging.yaml file.
type Config struct {
	Certificate       string
	Key               string
	Namespace         string
	Environment       string
	Account           string
	ConfigVersion     string
	ResourceID        string
	SubscriptionID    string
	Region            string
	ResourceName      string
	ResourceGroupName string
	File              *asset.File
	FluentbitImage    string
	MdsdImage         string
}

// GetConfig prepares a bootstraplogging.Config object based on
// the environment
func GetConfig(env env.Interface, oc *api.OpenShiftCluster) (*Config, error) {
	r, err := azure.ParseResourceID(oc.ID)
	if err != nil {
		return nil, err
	}

	key, cert := env.ClusterGenevaLoggingSecret()

	gcsKeyBytes, err := utiltls.PrivateKeyAsBytes(key)
	if err != nil {
		return nil, err
	}

	gcsCertBytes, err := utiltls.CertAsBytes(cert)
	if err != nil {
		return nil, err
	}

	return &Config{
		Certificate:       string(gcsCertBytes),
		Key:               string(gcsKeyBytes),
		Namespace:         env.ClusterGenevaLoggingNamespace(),
		Account:           env.ClusterGenevaLoggingAccount(),
		Environment:       env.ClusterGenevaLoggingEnvironment(),
		ConfigVersion:     env.ClusterGenevaLoggingConfigVersion(),
		Region:            env.Location(),
		ResourceID:        oc.ID,
		SubscriptionID:    r.SubscriptionID,
		ResourceName:      r.ResourceName,
		ResourceGroupName: r.ResourceGroup,
		MdsdImage:         version.MdsdImage(env.ACRDomain()),
		FluentbitImage:    version.FluentbitImage(env.ACRDomain()),
	}, nil
}
