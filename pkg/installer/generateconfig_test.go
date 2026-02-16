package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	mgmtcompute "github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/Azure/go-autorest/autorest/to"

	"github.com/openshift/installer/pkg/types/azure"
)

func TestVMNetworkingType(t *testing.T) {
	capabilityName := azure.AcceleratedNetworkingEnabled
	for _, tt := range []struct {
		name     string
		sku      *mgmtcompute.ResourceSku
		wantType string
	}{
		{
			name: "sku with support for accelerated networking",
			sku: &mgmtcompute.ResourceSku{
				Capabilities: &([]mgmtcompute.ResourceSkuCapabilities{
					{Name: &capabilityName, Value: to.StringPtr("True")},
				}),
			},
			wantType: "Accelerated",
		}, {
			name: "sku without support for accelerated networking",
			sku: &mgmtcompute.ResourceSku{
				Capabilities: &([]mgmtcompute.ResourceSkuCapabilities{
					{Name: &capabilityName, Value: to.StringPtr("False")},
				}),
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

func TestDetermineV2SkuSupport(t *testing.T) {
	for _, tt := range []struct {
		name       string
		sku        *mgmtcompute.ResourceSku
		wantResult bool
		wantErr    string
	}{
		{
			name: "sku supports both V1 and V2",
			sku: &mgmtcompute.ResourceSku{
				Name: to.StringPtr("Standard_D8s_v3"),
				Capabilities: &[]mgmtcompute.ResourceSkuCapabilities{
					{Name: to.StringPtr("HyperVGenerations"), Value: to.StringPtr("V1,V2")},
				},
			},
			wantResult: true,
		},
		{
			name: "sku supports only V2",
			sku: &mgmtcompute.ResourceSku{
				Name: to.StringPtr("Standard_D8s_v6"),
				Capabilities: &[]mgmtcompute.ResourceSkuCapabilities{
					{Name: to.StringPtr("HyperVGenerations"), Value: to.StringPtr("V2")},
				},
			},
			wantResult: true,
		},
		{
			name: "sku supports only V1",
			sku: &mgmtcompute.ResourceSku{
				Name: to.StringPtr("Standard_D2_v2"),
				Capabilities: &[]mgmtcompute.ResourceSkuCapabilities{
					{Name: to.StringPtr("HyperVGenerations"), Value: to.StringPtr("V1")},
				},
			},
			wantResult: false,
		},
		{
			name: "sku with empty capabilities returns error",
			sku: &mgmtcompute.ResourceSku{
				Name:         to.StringPtr("Standard_Empty"),
				Capabilities: &[]mgmtcompute.ResourceSkuCapabilities{},
			},
			wantErr: "no capabilities found for SKU Standard_Empty",
		},
		{
			name: "sku missing HyperVGenerations capability returns error",
			sku: &mgmtcompute.ResourceSku{
				Name: to.StringPtr("Standard_NoHyperV"),
				Capabilities: &[]mgmtcompute.ResourceSkuCapabilities{
					{Name: to.StringPtr("AcceleratedNetworkingEnabled"), Value: to.StringPtr("True")},
					{Name: to.StringPtr("vCPUs"), Value: to.StringPtr("8")},
				},
			},
			wantErr: "could not fetch HyperV generations for SKU Standard_NoHyperV: unable to determine HyperVGeneration version",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			result, err := determineV2SkuSupport(tt.sku)

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

func TestDetermineZones(t *testing.T) {
	for _, tt := range []struct {
		name                  string
		controlPlaneSkuZones  []string
		workerSkuZones        []string
		wantControlPlaneZones []string
		wantWorkerZones       []string
		allowExpandedAZs      bool
		wantErr               string
	}{

		{
			name:                  "non-zonal control plane, zonal workers",
			controlPlaneSkuZones:  nil,
			workerSkuZones:        []string{"1", "2", "3"},
			wantControlPlaneZones: []string{""},
			wantWorkerZones:       []string{"1", "2", "3"},
		},
		{
			name:                  "zonal control plane, non-zonal workers",
			controlPlaneSkuZones:  []string{"1", "2", "3"},
			workerSkuZones:        nil,
			wantControlPlaneZones: []string{"1", "2", "3"},
			wantWorkerZones:       []string{""},
		},
		{
			name:                  "zonal control plane, zonal workers",
			controlPlaneSkuZones:  []string{"1", "2", "3"},
			workerSkuZones:        []string{"1", "2", "3"},
			wantControlPlaneZones: []string{"1", "2", "3"},
			wantWorkerZones:       []string{"1", "2", "3"},
		},
		{
			name:                  "region with 4 availability zones, expanded AZs, control plane uses first 3, workers use all",
			allowExpandedAZs:      true,
			controlPlaneSkuZones:  []string{"1", "2", "3", "4"},
			workerSkuZones:        []string{"1", "2", "3", "4"},
			wantControlPlaneZones: []string{"1", "2", "3"},
			wantWorkerZones:       []string{"1", "2", "3", "4"},
		},
		{
			name:                  "region with 4 availability zones, basic AZs only, control plane and workers use 3",
			allowExpandedAZs:      false,
			controlPlaneSkuZones:  []string{"1", "2", "3", "4"},
			workerSkuZones:        []string{"1", "2", "3", "4"},
			wantControlPlaneZones: []string{"1", "2", "3"},
			wantWorkerZones:       []string{"1", "2", "3"},
		},
		{
			name:                 "not enough control plane zones",
			controlPlaneSkuZones: []string{"1", "2"},
			workerSkuZones:       []string{"1", "2", "3"},
			wantErr:              "cluster creation with 2 zones and 3 control plane replicas is unsupported",
		},
		{
			name:                 "not enough control plane zones, basic AZs only",
			controlPlaneSkuZones: []string{"1", "2", "4"},
			workerSkuZones:       []string{"1", "2", "3"},
			wantErr:              "cluster creation with 2 zones and 3 control plane replicas is unsupported",
		},
		{
			name:                 "not enough worker zones",
			controlPlaneSkuZones: []string{"1", "2", "3"},
			workerSkuZones:       []string{"1", "2"},
			wantErr:              "cluster creation with a worker SKU available on less than 3 zones is unsupported (available: 2)",
		},
		{
			name:                 "not enough worker zones, basic AZs only",
			controlPlaneSkuZones: []string{"1", "2", "3"},
			workerSkuZones:       []string{"1", "2", "4"},
			wantErr:              "cluster creation with a worker SKU available on less than 3 zones is unsupported (available: 2)",
		},
		{
			name:                  "region with 4 availability zones, expanded AZs, control plane only available in non-consecutive 3, workers use all",
			allowExpandedAZs:      true,
			controlPlaneSkuZones:  []string{"1", "2", "4"},
			workerSkuZones:        []string{"1", "2", "3", "4"},
			wantControlPlaneZones: []string{"1", "2", "4"},
			wantWorkerZones:       []string{"1", "2", "3", "4"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			controlPlaneSku := &mgmtcompute.ResourceSku{
				LocationInfo: &[]mgmtcompute.ResourceSkuLocationInfo{
					{Zones: &tt.controlPlaneSkuZones},
				},
			}
			workerSku := &mgmtcompute.ResourceSku{
				LocationInfo: &[]mgmtcompute.ResourceSkuLocationInfo{
					{Zones: &tt.workerSkuZones},
				},
			}

			if tt.allowExpandedAZs {
				t.Setenv(ALLOW_EXPANDED_AZ_ENV, "1")
			}

			controlPlaneZones, workerZones, err := determineAvailabilityZones(controlPlaneSku, workerSku)
			if err != nil && err.Error() != tt.wantErr {
				t.Error(cmp.Diff(tt.wantErr, err))
			}

			if !reflect.DeepEqual(controlPlaneZones, tt.wantControlPlaneZones) {
				t.Error(cmp.Diff(tt.wantControlPlaneZones, controlPlaneZones))
			}

			if !reflect.DeepEqual(workerZones, tt.wantWorkerZones) {
				t.Error(cmp.Diff(tt.wantWorkerZones, workerZones))
			}
		})
	}
}
