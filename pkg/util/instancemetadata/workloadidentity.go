package instancemetadata

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"fmt"
	"os"

	"github.com/openshift/installer-aro-wrapper/pkg/util/azureclient"
)

// WorkloadIdentityEnabled detects even partial webhook configuration so that a
// misconfigured workload cannot fall back to the Hive node identity.
func WorkloadIdentityEnabled() bool {
	for _, key := range []string{"ARO_INSTALLER_IDENTITY_CLIENT_ID", "AZURE_CLIENT_ID", "AZURE_TENANT_ID", "AZURE_FEDERATED_TOKEN_FILE"} {
		if _, found := os.LookupEnv(key); found {
			return true
		}
	}
	return false
}

func newWorkloadIdentity() (InstanceMetadata, error) {
	for _, key := range []string{
		"AZURE_CLIENT_ID",
		"AZURE_TENANT_ID",
		"AZURE_FEDERATED_TOKEN_FILE",
		"AZURE_ENVIRONMENT",
		"ARO_AZURE_SUBSCRIPTION_ID",
		"ARO_LOCATION",
		"ARO_RESOURCEGROUP",
	} {
		if os.Getenv(key) == "" {
			return nil, fmt.Errorf("environment variable %q must be set and non-empty (workload identity mode)", key)
		}
	}

	if clientID, found := os.LookupEnv("ARO_INSTALLER_IDENTITY_CLIENT_ID"); found {
		if clientID == "" || clientID != os.Getenv("AZURE_CLIENT_ID") {
			return nil, fmt.Errorf("ARO_INSTALLER_IDENTITY_CLIENT_ID must be non-empty and match AZURE_CLIENT_ID")
		}
	}

	environment, err := azureclient.EnvironmentFromName(os.Getenv("AZURE_ENVIRONMENT"))
	if err != nil {
		return nil, err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return &instanceMetadata{
		hostname:       hostname,
		tenantID:       os.Getenv("AZURE_TENANT_ID"),
		subscriptionID: os.Getenv("ARO_AZURE_SUBSCRIPTION_ID"),
		location:       os.Getenv("ARO_LOCATION"),
		resourceGroup:  os.Getenv("ARO_RESOURCEGROUP"),
		environment:    &environment,
	}, nil
}
