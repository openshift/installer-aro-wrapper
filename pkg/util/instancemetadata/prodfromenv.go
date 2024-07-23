package instancemetadata

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"fmt"
	"os"

	"github.com/openshift/installer-aro-wrapper/pkg/util/azureclient"
)

type prodFromEnv struct {
	instanceMetadata

	Getenv    func(key string) string
	LookupEnv func(key string) (string, bool)
}

func newProdFromEnv(ctx context.Context) (InstanceMetadata, error) {
	p := &prodFromEnv{
		Getenv:    os.Getenv,
		LookupEnv: os.LookupEnv,
	}

	err := p.populateInstanceMetadata()
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *prodFromEnv) populateInstanceMetadata() error {
	for _, key := range []string{
		"ARO_AZURE_ENVIRONMENT",
		"ARO_AZURE_SUBSCRIPTION_ID",
		"ARO_AZURE_TENANT_ID",
		"ARO_LOCATION",
		"ARO_RESOURCEGROUP",
	} {
		if _, found := p.LookupEnv(key); !found {
			return fmt.Errorf("environment variable %q unset", key)
		}
	}

	// optional env variables
	// * HOSTNAME_OVERRIDE: defaults to os.Hostname()

	environment, err := azureclient.EnvironmentFromName(p.Getenv("ARO_AZURE_ENVIRONMENT"))
	if err != nil {
		return err
	}
	p.environment = &environment
	p.subscriptionID = p.Getenv("ARO_AZURE_SUBSCRIPTION_ID")
	p.tenantID = p.Getenv("ARO_AZURE_TENANT_ID")
	p.location = p.Getenv("ARO_LOCATION")
	p.resourceGroup = p.Getenv("ARO_RESOURCEGROUP")
	p.hostname = p.Getenv("ARO_HOSTNAME_OVERRIDE") // empty string returned if not set

	if p.hostname == "" {
		hostname, err := os.Hostname()
		if err == nil {
			p.hostname = hostname
		}
	}

	return nil
}
