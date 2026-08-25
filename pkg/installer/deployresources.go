package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"reflect"
	"time"

	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"

	"github.com/openshift/installer-aro-wrapper/pkg/util/arm"
	"github.com/openshift/installer-aro-wrapper/pkg/util/stringutils"
)

func (m *manager) deployResourceTemplate(ctx context.Context) error {
	resourceGroup := stringutils.LastTokenByte(m.oc.Properties.ClusterProfile.ResourceGroupID, '/')
	account := "cluster" + m.oc.Properties.StorageSuffix

	pg, err := m.graph.LoadPersisted(ctx, resourceGroup, account)
	if err != nil {
		return err
	}

	var installConfig *installconfig.InstallConfig
	var machineMaster *machine.Master
	err = pg.Get(&installConfig, &machineMaster)
	if err != nil {
		return err
	}

	params := map[string]interface{}{}
	var paramType string

	if m.oc.UsesWorkloadIdentity() {
		paramType = "secureString"
		sasURL, err := m.graph.GetUserDelegatedSASIgnitionBlobURL(ctx, resourceGroup, account, `https://cluster`+m.oc.Properties.StorageSuffix+`.blob.`+m.env.Environment().StorageEndpointSuffix+`/ignition/bootstrap.ign`, m.oc.UsesWorkloadIdentity())
		if err != nil {
			return err
		}
		params["sas"] = map[string]string{
			"value": sasURL,
		}
	} else {
		paramType = "object"
		params["sas"] = map[string]interface{}{
			"value": map[string]interface{}{
				"signedStart":         m.oc.Properties.Install.Now.Format(time.RFC3339),
				"signedExpiry":        m.oc.Properties.Install.Now.Add(24 * time.Hour).Format(time.RFC3339),
				"signedPermission":    "rl",
				"signedResourceTypes": "o",
				"signedServices":      "b",
				"signedProtocol":      "https",
			},
		}
	}

	controlPlaneZones, hasAZs := convertControlPlaneZonesToArmParameter(installConfig.Config.ControlPlane.Platform.Azure.Zones)

	params["controlPlaneZones"] = map[string]interface{}{
		"value": controlPlaneZones,
	}

	t := &arm.Template{
		Schema:         "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
		ContentVersion: "1.0.0.0",
		Parameters: map[string]*arm.TemplateParameter{
			"sas": {
				Type: paramType,
			},
			"controlPlaneZones": {
				Type: "array",
			},
		},
		Resources: []*arm.Resource{
			m.networkBootstrapNIC(installConfig),
			m.networkMasterNICs(installConfig),
			m.computeBootstrapVM(installConfig),
			m.computeMasterVMs(installConfig, zones(installConfig), machineMaster, !hasAZs),
		},
	}

	if !hasAZs {
		t.Resources = append(t.Resources, m.controlPlaneAvailabilitySet(installConfig))
	}

	return arm.DeployTemplate(ctx, m.log, m.deployments, resourceGroup, "resources", t, params)
}

// Handle the case where nonzonal resources actually need to have an empty zone
// param instead of {""}
func zones(installConfig *installconfig.InstallConfig) *[]string {
	if reflect.DeepEqual(installConfig.Config.ControlPlane.Platform.Azure.Zones, []string{""}) ||
		reflect.DeepEqual(installConfig.Config.ControlPlane.Platform.Azure.Zones, []string{}) ||
		installConfig.Config.ControlPlane.Platform.Azure.Zones == nil {
		// Non-zonal
		return nil
	} else {
		// Use the zones we have been specified
		return &[]string{"[parameters('controlPlaneZones')[copyIndex(0)]]"}
	}
}

// convertControlPlaneZonesToArmParameter makes sure that an empty zone slice
// gets changed to []string{""}, and a single zone slice is triplicated, so it
// can be safely passed in via the arm parameter. It also returns if there are
// AZs in use.
func convertControlPlaneZonesToArmParameter(in []string) ([]string, bool) {
	if in == nil {
		return []string{""}, false
	}

	if reflect.DeepEqual(in, []string{}) {
		return []string{""}, false
	}

	if len(in) == 1 && in[0] == "" {
		return []string{""}, false
	}

	if len(in) == 1 {
		return []string{in[0], in[0], in[0]}, true
	}

	return in, true
}
