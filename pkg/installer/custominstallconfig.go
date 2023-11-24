package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/releaseimage"
	"github.com/openshift/installer/pkg/asset/targets"
	"github.com/openshift/installer/pkg/asset/templates/content/bootkube"
	"github.com/pkg/errors"

	"github.com/openshift/installer-aro-wrapper/pkg/api"
	"github.com/openshift/installer-aro-wrapper/pkg/bootstraplogging"
	"github.com/openshift/installer-aro-wrapper/pkg/cluster/graph"
	"github.com/openshift/installer-aro-wrapper/pkg/patches"
	"github.com/openshift/installer-aro-wrapper/pkg/patches/mtu3900"
	utilignition "github.com/openshift/installer-aro-wrapper/pkg/util/ignition"
)

// applyInstallConfigCustomisations modifies the InstallConfig and creates
// parent assets, then regenerates the InstallConfig for use for Ignition
// generation, etc.
func (m *manager) applyInstallConfigCustomisations(installConfig *installconfig.InstallConfig, image *releaseimage.Image) (graph.Graph, error) {
	clusterID := &installconfig.ClusterID{
		UUID:    m.clusterUUID,
		InfraID: m.oc.Properties.InfraID,
	}

	bootstrapLoggingConfig, err := bootstraplogging.GetConfig(m.env, m.oc)
	if err != nil {
		return nil, err
	}

	httpSecret := make([]byte, 64)
	_, err = rand.Read(httpSecret)
	if err != nil {
		return nil, err
	}

	imageRegistryConfig := &bootkube.AROImageRegistryConfig{
		AccountName:   m.oc.Properties.ImageRegistryStorageAccountName,
		ContainerName: "image-registry",
		HTTPSecret:    hex.EncodeToString(httpSecret),
	}

	dnsConfig := &bootkube.ARODNSConfig{
		APIIntIP:  m.oc.Properties.APIServerProfile.IntIP,
		IngressIP: m.oc.Properties.IngressProfiles[0].IP,
	}

	if m.oc.Properties.NetworkProfile.GatewayPrivateEndpointIP != "" {
		dnsConfig.GatewayPrivateEndpointIP = m.oc.Properties.NetworkProfile.GatewayPrivateEndpointIP
		dnsConfig.GatewayDomains = append(m.env.GatewayDomains(), m.oc.Properties.ImageRegistryStorageAccountName+".blob."+m.env.Environment().StorageEndpointSuffix)
	}

	g := graph.Graph{}
	g.Set(installConfig, image, clusterID, bootstrapLoggingConfig, dnsConfig, imageRegistryConfig)

	m.log.Print("resolving graph")
	for _, a := range targets.Cluster {
		err = g.Resolve(a)
		if err != nil {
			return nil, err
		}
	}

	// apply any patches to the generated graph
	err = m.applyPatches(g)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (m *manager) applyPatches(g graph.Graph) error {
	// Get the bootstrap configuration
	bootstrap := g.Get(&bootstrap.Bootstrap{}).(*bootstrap.Bootstrap)

	ignitionConfigs := make([]patches.IgnitionPatch, 0)

	// Handle MTU3900 feature flag
	if m.oc.Properties.NetworkProfile.MTUSize == api.MTU3900 {
		m.log.Printf("applying feature flag %s", api.FeatureFlagMTU3900)
		ignitionConfigs = append(ignitionConfigs, mtu3900.NewMTU3900())
	}

	// future patch additions go here...

	// Apply systemd unit replacements and add the files in
	for _, i := range ignitionConfigs {
		m.log.Printf("applying ignition config patch %#v", i)
		files, units, err := i.Files()
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to apply ignition config modification %v", i))
		}
		bootstrap.Common.Config.Systemd.Units = utilignition.MergeUnits(bootstrap.Common.Config.Systemd.Units, units)
		bootstrap.Config.Storage.Files = append(bootstrap.Config.Storage.Files, files...)
	}

	data, err := ignition.Marshal(bootstrap.Config)
	if err != nil {
		return errors.Wrap(err, "failed to Marshal Ignition config")
	}
	bootstrap.File.Data = data

	return nil
}
