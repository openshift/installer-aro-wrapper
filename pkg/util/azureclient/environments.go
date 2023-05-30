package azureclient

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/go-autorest/autorest/azure"
)

// AROEnvironment contains additional, cloud-specific information needed by ARO.
type AROEnvironment struct {
	azure.Environment
	ActualCloudName          string
	GenevaMonitoringEndpoint string
	AppSuffix                string

	Cloud cloud.Configuration
	// Microsoft identity platform scopes used by ARO
	// See https://learn.microsoft.com/EN-US/azure/active-directory/develop/scopes-oidc#the-default-scope
	ResourceManagerScope      string
	KeyVaultScope             string
	ActiveDirectoryGraphScope string
}

var (
	// PublicCloud contains additional ARO information for the public Azure cloud environment.
	PublicCloud = AROEnvironment{
		Environment:              azure.PublicCloud,
		ActualCloudName:          "AzureCloud",
		GenevaMonitoringEndpoint: "https://gcs.prod.monitoring.core.windows.net/",
		AppSuffix:                "aro.azure.com",

		ResourceManagerScope:      azure.PublicCloud.ResourceManagerEndpoint + "/.default",
		KeyVaultScope:             azure.PublicCloud.ResourceIdentifiers.KeyVault + "/.default",
		ActiveDirectoryGraphScope: azure.PublicCloud.GraphEndpoint + "/.default",
	}

	// USGovernmentCloud contains additional ARO information for the US Gov cloud environment.
	USGovernmentCloud = AROEnvironment{
		Environment:              azure.USGovernmentCloud,
		ActualCloudName:          "AzureUSGovernment",
		GenevaMonitoringEndpoint: "https://gcs.monitoring.core.usgovcloudapi.net/",
		AppSuffix:                "aro.azure.us",

		ResourceManagerScope:      azure.PublicCloud.ResourceManagerEndpoint + "/.default",
		KeyVaultScope:             azure.PublicCloud.ResourceIdentifiers.KeyVault + "/.default",
		ActiveDirectoryGraphScope: azure.PublicCloud.GraphEndpoint + "/.default",
	}
)

// EnvironmentFromName returns the AROEnvironment corresponding to the common name specified.
func EnvironmentFromName(name string) (AROEnvironment, error) {
	switch strings.ToUpper(name) {
	case "AZUREPUBLICCLOUD":
		return PublicCloud, nil
	case "AZUREUSGOVERNMENTCLOUD":
		return USGovernmentCloud, nil
	}
	return AROEnvironment{}, fmt.Errorf("cloud environment %q is unsupported by ARO", name)
}

func (e *AROEnvironment) ClientCertificateCredentialOptions() *azidentity.ClientCertificateCredentialOptions {
	return &azidentity.ClientCertificateCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: e.Cloud,
		},
	}
}

func (e *AROEnvironment) ClientSecretCredentialOptions() *azidentity.ClientSecretCredentialOptions {
	return &azidentity.ClientSecretCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: e.Cloud,
		},
	}
}

func (e *AROEnvironment) EnvironmentCredentialOptions() *azidentity.EnvironmentCredentialOptions {
	return &azidentity.EnvironmentCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: e.Cloud,
		},
	}
}

func (e *AROEnvironment) ManagedIdentityCredentialOptions() *azidentity.ManagedIdentityCredentialOptions {
	return &azidentity.ManagedIdentityCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: e.Cloud,
		},
	}
}
