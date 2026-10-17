package instancemetadata

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/openshift/installer-aro-wrapper/pkg/util/azureclient"
)

func setWorkloadEnvironment(t *testing.T) {
	t.Helper()
	for key, value := range map[string]string{
		"AZURE_CLIENT_ID":            "installer-hosting-client",
		"AZURE_TENANT_ID":            "hosting-tenant",
		"AZURE_FEDERATED_TOKEN_FILE": "/var/run/secrets/azure/tokens/azure-identity-token",
		"AZURE_ENVIRONMENT":          "AzurePublicCloud",
		"ARO_AZURE_SUBSCRIPTION_ID":  "rp-subscription",
		"ARO_LOCATION":               "eastus",
		"ARO_RESOURCEGROUP":          "rp-resource-group",
	} {
		t.Setenv(key, value)
	}
}

type rejectMetadataTransport struct {
	t *testing.T
}

func (r rejectMetadataTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	r.t.Errorf("unexpected network request: %s", req.URL.Host)
	return nil, fmt.Errorf("network access forbidden")
}

func TestWorkloadIdentityMetadata(t *testing.T) {
	for _, environment := range []azureclient.AROEnvironment{azureclient.PublicCloud, azureclient.USGovernmentCloud} {
		t.Run(environment.Name, func(t *testing.T) {
			setWorkloadEnvironment(t)
			t.Setenv("AZURE_ENVIRONMENT", environment.Name)
			// The hosting deployment must not be inferred from the customer cluster.
			t.Setenv("ARO_AZURE_TENANT_ID", "different-rp-tenant")
			t.Setenv(envAksShardId, "not-a-hive-shard")

			oldClient := http.DefaultClient
			http.DefaultClient = &http.Client{Transport: rejectMetadataTransport{t}}
			t.Cleanup(func() { http.DefaultClient = oldClient })

			im, err := New(context.Background(), logrus.NewEntry(logrus.New()), false)
			if err != nil {
				t.Fatal(err)
			}
			for _, check := range []struct {
				name string
				got  string
				want string
			}{
				{"tenant", im.TenantID(), "hosting-tenant"},
				{"subscription", im.SubscriptionID(), "rp-subscription"},
				{"location", im.Location(), "eastus"},
				{"resourceGroup", im.ResourceGroup(), "rp-resource-group"},
			} {
				if check.got != check.want {
					t.Errorf("%s = %q, want %q", check.name, check.got, check.want)
				}
			}
			if im.Environment().Name != environment.Name {
				t.Errorf("cloud = %q, want %q", im.Environment().Name, environment.Name)
			}
			if im.Hostname() == "" || im.AksMsiClientID() != "" {
				t.Fatal("expected pod hostname and no Hive agent-pool identity")
			}
		})
	}
}

func TestWorkloadIdentityMetadataRejectsIncompleteConfiguration(t *testing.T) {
	for _, key := range []string{
		"AZURE_CLIENT_ID", "AZURE_TENANT_ID", "AZURE_FEDERATED_TOKEN_FILE",
		"AZURE_ENVIRONMENT", "ARO_AZURE_SUBSCRIPTION_ID", "ARO_LOCATION", "ARO_RESOURCEGROUP",
	} {
		for _, unset := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/unset=%t", key, unset), func(t *testing.T) {
				setWorkloadEnvironment(t)
				t.Setenv(key, "")
				if unset {
					if err := os.Unsetenv(key); err != nil {
						t.Fatal(err)
					}
				}
				_, err := New(context.Background(), logrus.NewEntry(logrus.New()), false)
				if err == nil || !strings.Contains(err.Error(), key) {
					t.Fatalf("expected configuration error for %s, got %v", key, err)
				}
			})
		}
	}
	t.Run("unsupported cloud", func(t *testing.T) {
		setWorkloadEnvironment(t)
		t.Setenv("AZURE_ENVIRONMENT", "unsupported-cloud")
		_, err := New(context.Background(), logrus.NewEntry(logrus.New()), false)
		if err == nil || !strings.Contains(err.Error(), "unsupported by ARO") {
			t.Fatalf("expected unsupported cloud error, got %v", err)
		}
	})
}

func TestWorkloadIdentityEnabled(t *testing.T) {
	keys := []string{"ARO_INSTALLER_IDENTITY_CLIENT_ID", "AZURE_CLIENT_ID", "AZURE_TENANT_ID", "AZURE_FEDERATED_TOKEN_FILE"}
	for _, key := range keys {
		t.Setenv(key, "")
		if err := os.Unsetenv(key); err != nil {
			t.Fatal(err)
		}
	}
	if WorkloadIdentityEnabled() {
		t.Fatal("legacy Hive configuration should not enable workload identity")
	}
	for _, key := range keys {
		t.Run(key, func(t *testing.T) {
			t.Setenv(key, "")
			if !WorkloadIdentityEnabled() {
				t.Fatal("partial or empty webhook configuration must not fall back to IMDS")
			}
		})
	}
}

func TestWorkloadIdentityRequiresExpectedHostingIdentity(t *testing.T) {
	for _, clientID := range []string{"installer-hosting-client", "wrong-client", ""} {
		t.Run(clientID, func(t *testing.T) {
			setWorkloadEnvironment(t)
			t.Setenv("ARO_INSTALLER_IDENTITY_CLIENT_ID", clientID)
			_, err := New(context.Background(), logrus.NewEntry(logrus.New()), false)
			if clientID == "installer-hosting-client" {
				if err != nil {
					t.Fatal(err)
				}
			} else if err == nil || !strings.Contains(err.Error(), "match AZURE_CLIENT_ID") {
				t.Fatalf("expected hosting identity mismatch error, got %v", err)
			}
		})
	}
	t.Run("webhook did not run", func(t *testing.T) {
		setWorkloadEnvironment(t)
		t.Setenv("ARO_INSTALLER_IDENTITY_CLIENT_ID", "installer-hosting-client")
		for _, key := range []string{"AZURE_CLIENT_ID", "AZURE_TENANT_ID", "AZURE_FEDERATED_TOKEN_FILE"} {
			if err := os.Unsetenv(key); err != nil {
				t.Fatal(err)
			}
		}
		_, err := New(context.Background(), logrus.NewEntry(logrus.New()), false)
		if err == nil || !strings.Contains(err.Error(), "AZURE_CLIENT_ID") {
			t.Fatalf("expected missing webhook configuration error, got %v", err)
		}
	})
}

func TestDevelopmentMetadataTakesPrecedence(t *testing.T) {
	setWorkloadEnvironment(t)
	t.Setenv("ARO_AZURE_TENANT_ID", "development-tenant")
	t.Setenv("AZURE_FEDERATED_TOKEN_FILE", "")
	im, err := New(context.Background(), logrus.NewEntry(logrus.New()), true)
	if err != nil {
		t.Fatal(err)
	}
	if im.TenantID() != "development-tenant" {
		t.Fatalf("expected development metadata, got tenant %q", im.TenantID())
	}
}
