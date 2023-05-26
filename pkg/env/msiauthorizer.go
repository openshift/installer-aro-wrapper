package env

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/go-autorest/autorest"
	"github.com/jongio/azidext/go/azidext"
)

type MSIContext string

const (
	MSIContextRP      MSIContext = "RP"
	MSIContextGateway MSIContext = "GATEWAY"
)

func (c *core) NewMSIAuthorizer(msiContext MSIContext, scopes ...string) (autorest.Authorizer, error) {
	var tokenCredential azcore.TokenCredential
	var err error

	if !c.IsLocalDevelopmentMode() {
		options := c.Environment().ManagedIdentityCredentialOptions()
		tokenCredential, err = azidentity.NewManagedIdentityCredential(options)
	} else {
		for _, key := range []string{
			"ARO_AZURE_" + string(msiContext) + "_CLIENT_ID",
			"ARO_AZURE_" + string(msiContext) + "_CLIENT_SECRET",
			"ARO_AZURE_TENANT_ID",
		} {
			if _, found := os.LookupEnv(key); !found {
				return nil, fmt.Errorf("environment variable %q unset (development mode)", key)
			}
		}

		options := c.Environment().ClientSecretCredentialOptions()
		tokenCredential, err = azidentity.NewClientSecretCredential(
			os.Getenv("ARO_AZURE_TENANT_ID"),
			os.Getenv("ARO_AZURE_"+string(msiContext)+"_CLIENT_ID"),
			os.Getenv("ARO_AZURE_"+string(msiContext)+"_CLIENT_SECRET"),
			options)
	}
	if err != nil {
		return nil, err
	}

	return azidext.NewTokenCredentialAdapter(tokenCredential, scopes), nil
}
