package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"testing"

	sdkcompute "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v7"
	"github.com/Azure/go-autorest/autorest/to"

	"github.com/openshift/installer/pkg/types/azure"
)

func TestVMNetworkingType(t *testing.T) {
	capabilityName := azure.AcceleratedNetworkingEnabled
	for _, tt := range []struct {
		name     string
		sku      *sdkcompute.ResourceSKU
		wantType string
	}{
		{
			name: "sku with support for accelerated networking",
			sku: &sdkcompute.ResourceSKU{
				Capabilities: []*sdkcompute.ResourceSKUCapabilities{
					{Name: &capabilityName, Value: to.StringPtr("True")},
				},
			},
			wantType: "Accelerated",
		}, {
			name: "sku without support for accelerated networking",
			sku: &sdkcompute.ResourceSKU{
				Capabilities: []*sdkcompute.ResourceSKUCapabilities{
					{Name: &capabilityName, Value: to.StringPtr("False")},
				},
			},
			wantType: "Basic",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			result := determineVMNetworkingType(tt.sku)

			if result != tt.wantType {
				t.Error(result)
			}
		})
	}
}

func TestDetermineSkuSupportsV2Only(t *testing.T) {
	for _, tt := range []struct {
		name       string
		sku        *sdkcompute.ResourceSKU
		wantResult bool
		wantErr    string
	}{
		{
			name: "sku supports both V1 and V2, does not require V2",
			sku: &sdkcompute.ResourceSKU{
				Name: to.StringPtr("Standard_D8s_v3"),
				Capabilities: []*sdkcompute.ResourceSKUCapabilities{
					{Name: to.StringPtr("HyperVGenerations"), Value: to.StringPtr("V1,V2")},
				},
			},
			wantResult: false,
		},
		{
			name: "sku supports only V2, requires V2",
			sku: &sdkcompute.ResourceSKU{
				Name: to.StringPtr("Standard_D8s_v6"),
				Capabilities: []*sdkcompute.ResourceSKUCapabilities{
					{Name: to.StringPtr("HyperVGenerations"), Value: to.StringPtr("V2")},
				},
			},
			wantResult: true,
		},
		{
			name: "sku supports only V1, does not require V2",
			sku: &sdkcompute.ResourceSKU{
				Name: to.StringPtr("Standard_D2_v2"),
				Capabilities: []*sdkcompute.ResourceSKUCapabilities{
					{Name: to.StringPtr("HyperVGenerations"), Value: to.StringPtr("V1")},
				},
			},
			wantResult: false,
		},
		{
			name: "sku with empty capabilities returns error",
			sku: &sdkcompute.ResourceSKU{
				Name:         to.StringPtr("Standard_Empty"),
				Capabilities: []*sdkcompute.ResourceSKUCapabilities{},
			},
			wantErr: "no capabilities found for SKU Standard_Empty",
		},
		{
			name: "sku missing HyperVGenerations capability returns error",
			sku: &sdkcompute.ResourceSKU{
				Name: to.StringPtr("Standard_NoHyperV"),
				Capabilities: []*sdkcompute.ResourceSKUCapabilities{
					{Name: to.StringPtr("AcceleratedNetworkingEnabled"), Value: to.StringPtr("True")},
					{Name: to.StringPtr("vCPUs"), Value: to.StringPtr("8")},
				},
			},
			wantErr: "could not fetch HyperV generations for SKU Standard_NoHyperV: unable to determine HyperVGeneration version",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			result, err := determineSkuSupportsV2Only(tt.sku)

			if tt.wantErr != "" {
				if err == nil {
					t.Errorf("expected error %q, got nil", tt.wantErr)
				} else if err.Error() != tt.wantErr {
					t.Errorf("expected error %q, got %q", tt.wantErr, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result != tt.wantResult {
				t.Errorf("expected %v, got %v", tt.wantResult, result)
			}
		})
	}
}
