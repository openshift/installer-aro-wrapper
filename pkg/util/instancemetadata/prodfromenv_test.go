package instancemetadata

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"os"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/openshift/installer-aro-wrapper/pkg/util/azureclient"
)

func TestProdEnvPopulateInstanceMetadata(t *testing.T) {
	hostname, err := os.Hostname()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name                 string
		environment          map[string]string
		wantInstanceMetadata instanceMetadata
		wantErr              string
	}{
		{
			name:    "missing environment variables",
			wantErr: "environment variable \"ARO_AZURE_ENVIRONMENT\" unset",
		},
		{
			name: "valid environment variables",
			environment: map[string]string{
				"ARO_AZURE_ENVIRONMENT":     azureclient.PublicCloud.Name,
				"ARO_AZURE_SUBSCRIPTION_ID": "some-sub-guid",
				"ARO_AZURE_TENANT_ID":       "some-tenant-guid",
				"ARO_LOCATION":              "some-region",
				"ARO_RESOURCEGROUP":         "my-resourceGroup",
			},
			wantInstanceMetadata: instanceMetadata{
				environment:    &azureclient.PublicCloud,
				subscriptionID: "some-sub-guid",
				tenantID:       "some-tenant-guid",
				location:       "some-region",
				resourceGroup:  "my-resourceGroup",
				hostname:       hostname,
			},
		},
		{
			name: "valid environment variables, but invalid Azure environment name",
			environment: map[string]string{
				"ARO_AZURE_ENVIRONMENT":     "ThisEnvDoesNotExist",
				"ARO_AZURE_SUBSCRIPTION_ID": "some-sub-guid",
				"ARO_AZURE_TENANT_ID":       "some-tenant-guid",
				"ARO_LOCATION":              "some-region",
				"ARO_RESOURCEGROUP":         "my-resourceGroup",
			},
			wantErr: "cloud environment \"ThisEnvDoesNotExist\" is unsupported by ARO",
		},
		{
			name: "valid environment variables with hostname override",
			environment: map[string]string{
				"ARO_AZURE_ENVIRONMENT":     azureclient.PublicCloud.Name,
				"ARO_AZURE_SUBSCRIPTION_ID": "some-sub-guid",
				"ARO_AZURE_TENANT_ID":       "some-tenant-guid",
				"ARO_LOCATION":              "some-region",
				"ARO_RESOURCEGROUP":         "my-resourceGroup",
				"ARO_HOSTNAME_OVERRIDE":     "my.over.ride",
			},
			wantInstanceMetadata: instanceMetadata{
				environment:    &azureclient.PublicCloud,
				subscriptionID: "some-sub-guid",
				tenantID:       "some-tenant-guid",
				location:       "some-region",
				resourceGroup:  "my-resourceGroup",
				hostname:       "my.over.ride",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := &prodFromEnv{
				Getenv: func(key string) string {
					return test.environment[key]
				},
				LookupEnv: func(key string) (string, bool) {
					value, ok := test.environment[key]
					return value, ok
				},
			}

			err := p.populateInstanceMetadata()
			if err != nil && err.Error() != test.wantErr ||
				err == nil && test.wantErr != "" {
				t.Error(err)
			}
			if !reflect.DeepEqual(p.instanceMetadata, test.wantInstanceMetadata) {
				opts := cmp.AllowUnexported(instanceMetadata{})
				t.Error(cmp.Diff(p.instanceMetadata, test.wantInstanceMetadata, opts))
			}
		})
	}
}
