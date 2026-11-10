package env

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"

	"github.com/openshift/installer-aro-wrapper/pkg/util/azureclient"
	mock_instancemetadata "github.com/openshift/installer-aro-wrapper/pkg/util/mocks/instancemetadata"
)

func TestHostingCredentialSelection(t *testing.T) {
	for _, tt := range []struct {
		name        string
		development bool
		workload    bool
		wantType    string
	}{
		{name: "legacy Hive", wantType: "*azidentity.ManagedIdentityCredential"},
		{name: "AKS workload identity", workload: true, wantType: "*azidentity.WorkloadIdentityCredential"},
		{name: "local development", development: true, wantType: "*azidentity.ClientSecretCredential"},
		{name: "local development with webhook variables", development: true, workload: true, wantType: "*azidentity.ClientSecretCredential"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, key := range []string{"ARO_INSTALLER_IDENTITY_CLIENT_ID", "AZURE_CLIENT_ID", "AZURE_TENANT_ID", "AZURE_FEDERATED_TOKEN_FILE"} {
				t.Setenv(key, "")
				if err := os.Unsetenv(key); err != nil {
					t.Fatal(err)
				}
			}
			if tt.workload {
				t.Setenv("AZURE_CLIENT_ID", "installer-client")
				t.Setenv("AZURE_TENANT_ID", "hosting-tenant")
				t.Setenv("AZURE_FEDERATED_TOKEN_FILE", "/projected/token")
			}
			t.Setenv("ARO_AZURE_TENANT_ID", "development-tenant")
			t.Setenv("ARO_AZURE_RP_CLIENT_ID", "development-client")
			t.Setenv("ARO_AZURE_RP_CLIENT_SECRET", "development-secret")

			im := mock_instancemetadata.NewMockInstanceMetadata(gomock.NewController(t))
			im.EXPECT().Environment().Return(&azureclient.PublicCloud)
			if !tt.development && !tt.workload {
				im.EXPECT().AksMsiClientID().Return("hive-agentpool-client")
			}
			p := &prod{InstanceMetadata: im, isLocalDevelopmentMode: tt.development}
			credential, err := p.newMSITokenCredential(MSIContextRP)
			if err != nil {
				t.Fatal(err)
			}
			if got := fmt.Sprintf("%T", credential); got != tt.wantType {
				t.Fatalf("credential = %s, want %s", got, tt.wantType)
			}
		})
	}
}

func TestHostingCredentialRejectsIncompleteWorkloadIdentity(t *testing.T) {
	for _, key := range []string{"AZURE_CLIENT_ID", "AZURE_TENANT_ID", "AZURE_FEDERATED_TOKEN_FILE"} {
		t.Run(key, func(t *testing.T) {
			t.Setenv("AZURE_CLIENT_ID", "installer-client")
			t.Setenv("AZURE_TENANT_ID", "hosting-tenant")
			t.Setenv("AZURE_FEDERATED_TOKEN_FILE", "/projected/token")
			t.Setenv(key, "")
			p := &prod{}
			_, err := p.NewMSIAuthorizer(MSIContextRP, azureclient.PublicCloud.KeyVaultScope)
			if err == nil || !strings.Contains(err.Error(), key) {
				t.Fatalf("expected configuration error for %s, got %v", key, err)
			}
		})
	}
}

func TestWorkloadIdentityAuthorizer(t *testing.T) {
	t.Setenv("AZURE_CLIENT_ID", "installer-client")
	t.Setenv("AZURE_TENANT_ID", "hosting-tenant")
	t.Setenv("AZURE_FEDERATED_TOKEN_FILE", "/projected/token")
	im := mock_instancemetadata.NewMockInstanceMetadata(gomock.NewController(t))
	im.EXPECT().Environment().Return(&azureclient.USGovernmentCloud)
	p := &prod{InstanceMetadata: im}
	credential, err := p.newMSITokenCredential(MSIContextRP)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := credential.(*azidentity.WorkloadIdentityCredential); !ok {
		t.Fatalf("expected workload identity credential, got %T", credential)
	}
	im.EXPECT().Environment().Return(&azureclient.USGovernmentCloud)
	if _, err := p.NewMSIAuthorizer(MSIContextRP, azureclient.USGovernmentCloud.KeyVaultScope); err != nil {
		t.Fatal(err)
	}
}
