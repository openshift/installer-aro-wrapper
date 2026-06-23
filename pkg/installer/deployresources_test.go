package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"reflect"
	"testing"

	"github.com/Azure/go-autorest/autorest/to"

	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
)

func TestZones(t *testing.T) {
	for _, tt := range []struct {
		name       string
		zones      []string
		region     string
		wantMaster *[]string
	}{
		{
			name:       "non-zonal, empty slice",
			zones:      []string{},
			wantMaster: nil,
		},
		{
			name:       "non-zonal",
			zones:      []string{""},
			wantMaster: nil,
		},
		{
			name:       "zonal",
			zones:      []string{"1", "2", "3"},
			wantMaster: &[]string{"[parameters('controlPlaneZones')[copyIndex(0)]]"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			z := zones(&installconfig.InstallConfig{
				AssetBase: installconfig.AssetBase{
					Config: &types.InstallConfig{
						ControlPlane: &types.MachinePool{
							Platform: types.MachinePoolPlatform{
								Azure: &azuretypes.MachinePool{
									Zones: tt.zones,
								},
							},
							Replicas: to.Int64Ptr(3),
						},
						Platform: types.Platform{
							Azure: &azuretypes.Platform{
								Region: tt.region,
								DefaultMachinePlatform: &azuretypes.MachinePool{
									Zones: tt.zones,
								},
							},
						},
					},
				},
			})
			if !reflect.DeepEqual(tt.wantMaster, z) {
				t.Errorf("Expected master %v, got master %v", tt.wantMaster, z)
			}
		})
	}
}

func Test_convertControlPlaneZonesToArmParameter(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{
			name: "empty slice is converted to slice with empty string",
			in:   []string{},
			want: []string{""},
		},
		{
			name: "nil slice is returned as-is",
			in:   nil,
			want: []string{""},
		},
		{
			name: "non-empty zone slice is returned unchanged",
			in:   []string{"1", "2", "3"},
			want: []string{"1", "2", "3"},
		},
		{
			name: "slice with empty string is returned unchanged",
			in:   []string{""},
			want: []string{""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertControlPlaneZonesToArmParameter(tt.in)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertControlPlaneZonesToArmParameter() = %v, want %v", got, tt.want)
			}
		})
	}
}
