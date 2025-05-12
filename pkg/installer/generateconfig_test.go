package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"reflect"
	"testing"

	mgmtcompute "github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/google/go-cmp/cmp"
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
