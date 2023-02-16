package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"testing"

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
