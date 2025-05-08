package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"reflect"
	"time"

	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"

	"github.com/openshift/ARO-Installer/pkg/util/arm"
	"github.com/openshift/ARO-Installer/pkg/util/stringutils"
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

	t := &arm.Template{
		Schema:         "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
		ContentVersion: "1.0.0.0",
		Parameters: map[string]*arm.TemplateParameter{
			"sas": {
				Type: paramType,
			},
		},
		Resources: []*arm.Resource{
			m.networkBootstrapNIC(installConfig),
			m.networkMasterNICs(installConfig),
			m.computeBootstrapVM(installConfig),
			m.computeMasterVMs(installConfig, zones(installConfig), machineMaster),
		},
	}

	return arm.DeployTemplate(ctx, m.log, m.deployments, resourceGroup, "resources", t, params)
}

// Handle the case where nonzonal resources actually need to have an empty zone
// param instead of {""}
func zones(installConfig *installconfig.InstallConfig) *[]string {
	if reflect.DeepEqual(installConfig.Config.ControlPlane.Platform.Azure.Zones, []string{""}) {
		// Non-zonal
		return nil
	} else {
		// Use the zones we have been specified
		return &installConfig.Config.ControlPlane.Platform.Azure.Zones
	}
}
