package computeskus

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	mgmtcompute "github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/Azure/go-autorest/autorest/to"
)

func TestZones(t *testing.T) {
	for _, tt := range []struct {
		name      string
		sku       *mgmtcompute.ResourceSku
		wantZones []string
	}{
		{
			name: "sku with location info present",
			sku: &mgmtcompute.ResourceSku{
				LocationInfo: &([]mgmtcompute.ResourceSkuLocationInfo{
					{Zones: &([]string{"1", "2", "3"})},
				}),
			},
			wantZones: []string{"1", "2", "3"},
		},
		{
			name: "sku with location info present, but zones field is nil",
			sku: &mgmtcompute.ResourceSku{
				LocationInfo: &([]mgmtcompute.ResourceSkuLocationInfo{
					{Zones: nil},
				}),
			},
		},
		{
			name: "sku with location info present, but empty",
			sku: &mgmtcompute.ResourceSku{
				LocationInfo: &([]mgmtcompute.ResourceSkuLocationInfo{}),
			},
		},
		{
			name: "sku with location info missing",
			sku:  &mgmtcompute.ResourceSku{},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			zones := Zones(tt.sku)

			if !reflect.DeepEqual(tt.wantZones, zones) {
				t.Error(cmp.Diff(tt.wantZones, zones))
			}
		})
	}
}

func TestHasCapability(t *testing.T) {
	fakeCapabilityName := "fakeCapability"

	for _, tt := range []struct {
		name       string
		sku        *mgmtcompute.ResourceSku
		wantResult bool
	}{
		{
			name: "sku explicitly supports capability",
			sku: &mgmtcompute.ResourceSku{
				Capabilities: &([]mgmtcompute.ResourceSkuCapabilities{
					{Name: &fakeCapabilityName, Value: to.StringPtr("True")},
				}),
			},
			wantResult: true,
		},
		{
			name: "sku explicitly does not support capability",
			sku: &mgmtcompute.ResourceSku{
				Capabilities: &([]mgmtcompute.ResourceSkuCapabilities{
					{Name: &fakeCapabilityName, Value: to.StringPtr("False")},
				}),
			},
		},
		{
			name: "sku implicitly does not support capability because it is missing from the list",
			sku: &mgmtcompute.ResourceSku{
				Capabilities: &([]mgmtcompute.ResourceSkuCapabilities{}),
			},
		},
		{
			name: "sku implicitly does not support capability, because capabilities info missing",
			sku:  &mgmtcompute.ResourceSku{},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			result := HasCapability(tt.sku, fakeCapabilityName)

			if result != tt.wantResult {
				t.Error(result)
			}
		})
	}
}

func TestFilterVmSizes(t *testing.T) {
	for _, tt := range []struct {
		name             string
		providedLocation string
		resourceType     string
		skuLocation      []string
		skuRestrictions  mgmtcompute.ResourceSkuRestrictions
		skuLocationInfo  []mgmtcompute.ResourceSkuLocationInfo
		skuCapabilities  string
		wantResult       map[string]*mgmtcompute.ResourceSku
	}{
		{
			name:             "resource type is a virtual machine",
			providedLocation: "eastus",
			resourceType:     "virtualMachines",
			skuRestrictions:  mgmtcompute.ResourceSkuRestrictions{ReasonCode: mgmtcompute.NotAvailableForSubscription},
			skuLocation:      []string{"eastus"},
			skuLocationInfo:  []mgmtcompute.ResourceSkuLocationInfo{{Zones: &[]string{"eastus-2"}}},
			skuCapabilities:  "some-capability",

			wantResult: map[string]*mgmtcompute.ResourceSku{
				"Fake_Sku": {
					Name: to.StringPtr("Fake_Sku"),
					Restrictions: &[]mgmtcompute.ResourceSkuRestrictions{{
						ReasonCode: mgmtcompute.NotAvailableForSubscription}},
					LocationInfo: &[]mgmtcompute.ResourceSkuLocationInfo{{
						Zones: &[]string{"eastus-2"}},
					},
					Capabilities: &[]mgmtcompute.ResourceSkuCapabilities{{
						Name: to.StringPtr("some-capability"),
					}},
				},
			},
		},
		{
			name:             "resource type not a virtual machine",
			providedLocation: "eastus",
			resourceType:     "disk",
			skuLocation:      []string{"eastus"},
			skuLocationInfo:  []mgmtcompute.ResourceSkuLocationInfo{{Zones: &[]string{"eastus-2"}}},
			wantResult:       map[string]*mgmtcompute.ResourceSku{},
		},
		{
			name:             "sku Location doesn't match provided location",
			providedLocation: "mars",
			resourceType:     "virtualMachines",
			skuLocation:      []string{"eastus"},
			skuLocationInfo:  []mgmtcompute.ResourceSkuLocationInfo{{Zones: &[]string{"eastus-2"}}},
			wantResult:       map[string]*mgmtcompute.ResourceSku{},
		},
		{
			name:             "sku Location has length of 0",
			providedLocation: "eastus",
			resourceType:     "virtualMachines",
			skuLocation:      []string{},
			skuLocationInfo:  []mgmtcompute.ResourceSkuLocationInfo{{Zones: &[]string{"eastus-2"}}},
			wantResult:       map[string]*mgmtcompute.ResourceSku{},
		},
		{
			name:             "sku LocationInfo has length of 0",
			providedLocation: "eastus",
			resourceType:     "virtualMachines",
			skuLocation:      []string{"eastus"},
			skuLocationInfo:  []mgmtcompute.ResourceSkuLocationInfo{},
			wantResult:       map[string]*mgmtcompute.ResourceSku{},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			sku := []mgmtcompute.ResourceSku{
				{
					Name: to.StringPtr("Fake_Sku"),
					Capabilities: &[]mgmtcompute.ResourceSkuCapabilities{
						{
							Name: to.StringPtr(tt.skuCapabilities),
						},
					},
					Locations:    &tt.skuLocation,
					Restrictions: &[]mgmtcompute.ResourceSkuRestrictions{tt.skuRestrictions},
					LocationInfo: &tt.skuLocationInfo,
					ResourceType: to.StringPtr(tt.resourceType),
				},
			}

			result := FilterVMSizes(sku, tt.providedLocation)

			if !reflect.DeepEqual(result, tt.wantResult) {
				t.Error(cmp.Diff(result, tt.wantResult))
			}
		})
	}
}

func TestIsRestricted(t *testing.T) {
	for _, tt := range []struct {
		name       string
		location   string
		vmsize     string
		sku        map[string]*mgmtcompute.ResourceSku
		wantResult bool
	}{
		{
			name:     "sku is restricted in one location",
			location: "eastus",
			vmsize:   "Standard_Sku_1",
			sku: map[string]*mgmtcompute.ResourceSku{
				"Standard_Sku_1": {Restrictions: &[]mgmtcompute.ResourceSkuRestrictions{
					{
						RestrictionInfo: &mgmtcompute.ResourceSkuRestrictionInfo{Locations: &[]string{"eastus"}},
					},
				}},
			},
			wantResult: true,
		},
		{
			name:     "sku is restricted in multiple locations",
			location: "eastus",
			vmsize:   "Standard_Sku_1",
			sku: map[string]*mgmtcompute.ResourceSku{
				"Standard_Sku_1": {Restrictions: &[]mgmtcompute.ResourceSkuRestrictions{
					{
						RestrictionInfo: &mgmtcompute.ResourceSkuRestrictionInfo{Locations: &[]string{
							"eastus",
							"eastus2",
						}},
					},
				}},
			},
			wantResult: true,
		},
		{
			name:     "sku is not restricted",
			location: "eastus",
			vmsize:   "Standard_Sku_2",
			sku: map[string]*mgmtcompute.ResourceSku{
				"Standard_Sku_2": {Restrictions: &[]mgmtcompute.ResourceSkuRestrictions{
					{
						RestrictionInfo: &mgmtcompute.ResourceSkuRestrictionInfo{Locations: &[]string{""}},
					},
				}},
			},
			wantResult: false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRestricted(tt.sku, tt.location, tt.vmsize)

			if result != tt.wantResult {
				t.Error(result)
			}
		})
	}
}

func TestGetCapabilityValue(t *testing.T) {
	capabilityName := "TestCapability"

	for _, tt := range []struct {
		name      string
		sku       *mgmtcompute.ResourceSku
		wantValue string
		wantFound bool
	}{
		{
			name: "capability exists",
			sku: &mgmtcompute.ResourceSku{
				Capabilities: &[]mgmtcompute.ResourceSkuCapabilities{
					{Name: &capabilityName, Value: to.StringPtr("TestValue")},
				},
			},
			wantValue: "TestValue",
			wantFound: true,
		},
		{
			name: "capability does not exist",
			sku: &mgmtcompute.ResourceSku{
				Capabilities: &[]mgmtcompute.ResourceSkuCapabilities{
					{Name: to.StringPtr("OtherCapability"), Value: to.StringPtr("OtherValue")},
				},
			},
			wantValue: "",
			wantFound: false,
		},
		{
			name:      "capabilities is nil",
			sku:       &mgmtcompute.ResourceSku{},
			wantValue: "",
			wantFound: false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			value, found := GetCapabilityValue(tt.sku, capabilityName)

			if value != tt.wantValue {
				t.Errorf("expected value %q, got %q", tt.wantValue, value)
			}
			if found != tt.wantFound {
				t.Errorf("expected found %v, got %v", tt.wantFound, found)
			}
		})
	}
}

func TestRequiresHyperVGenerationV2Only(t *testing.T) {
	hyperVGenCapability := "HyperVGenerations"

	for _, tt := range []struct {
		name       string
		sku        *mgmtcompute.ResourceSku
		wantResult bool
	}{
		{
			name: "supports only V2 - requires Gen2",
			sku: &mgmtcompute.ResourceSku{
				Capabilities: &[]mgmtcompute.ResourceSkuCapabilities{
					{Name: &hyperVGenCapability, Value: to.StringPtr("V2")},
				},
			},
			wantResult: true,
		},
		{
			name: "supports V1 and V2 - does not require Gen2",
			sku: &mgmtcompute.ResourceSku{
				Capabilities: &[]mgmtcompute.ResourceSkuCapabilities{
					{Name: &hyperVGenCapability, Value: to.StringPtr("V1,V2")},
				},
			},
			wantResult: false,
		},
		{
			name: "supports only V1 - does not require Gen2",
			sku: &mgmtcompute.ResourceSku{
				Capabilities: &[]mgmtcompute.ResourceSkuCapabilities{
					{Name: &hyperVGenCapability, Value: to.StringPtr("V1")},
				},
			},
			wantResult: false,
		},
		{
			name:       "no HyperVGenerations capability",
			sku:        &mgmtcompute.ResourceSku{},
			wantResult: false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			result := RequiresHyperVGenerationV2Only(tt.sku)

			if result != tt.wantResult {
				t.Errorf("expected %v, got %v", tt.wantResult, result)
			}
		})
	}
}
