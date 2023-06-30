package bootstraplogging

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"embed"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/coreos/ignition/v2/config/v3_2/types"

	"github.com/Azure/ARO-RP/pkg/api"
	"github.com/Azure/ARO-RP/pkg/env"
	"github.com/Azure/ARO-RP/pkg/util/ignition"
	utiltls "github.com/Azure/ARO-RP/pkg/util/tls"
	"github.com/Azure/ARO-RP/pkg/util/version"
)

//go:embed staticresources
var loggingStaticFiles embed.FS

var enabledUnits = map[string]bool{
	"fluentbit.service": true,
	"mdsd.service":      true,
}

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

func Files(config *Config) ([]types.File, []types.Unit, error) {
	return ignition.GetFiles(loggingStaticFiles, map[string]*Config{"LoggingConfig": config}, enabledUnits)
}
