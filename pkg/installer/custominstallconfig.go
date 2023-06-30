package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/releaseimage"
	"github.com/openshift/installer/pkg/asset/targets"
	"github.com/openshift/installer/pkg/asset/templates/content/bootkube"
	"github.com/pkg/errors"

	"github.com/Azure/ARO-RP/pkg/api"
	"github.com/Azure/ARO-RP/pkg/bootstraplogging"
	"github.com/Azure/ARO-RP/pkg/cluster/graph"
	"github.com/Azure/ARO-RP/pkg/patches/bootstrapimagepull"
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
	g.Set(installConfig, image, clusterID, dnsConfig, imageRegistryConfig)

	m.log.Print("resolving graph")
	for _, a := range targets.Cluster {
		err = g.Resolve(a)
		if err != nil {
			return nil, err
		}
	}

	// Get the bootstrap configuration
	bootstrap := g.Get(&bootstrap.Bootstrap{}).(*bootstrap.Bootstrap)

	// Handle MTU3900 feature flag
	if m.oc.Properties.NetworkProfile.MTUSize == api.MTU3900 {
		m.log.Printf("applying feature flag %s", api.FeatureFlagMTU3900)
		if err := m.overrideEthernetMTU(bootstrap); err != nil {
			return nil, err
		}
	}

	err = m.addBootstrapLogging(bootstrap, bootstrapLoggingConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate bootstrap logging config")
	}

	err = m.addBootstrapImagePullLogging(bootstrap)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate bootstrap image pull logging")
	}

	data, err := ignition.Marshal(bootstrap.Config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to Marshal Ignition config")
	}
	bootstrap.File.Data = data

	return g, nil
}

func (m *manager) addBootstrapLogging(bootstrap *bootstrap.Bootstrap, config *bootstraplogging.Config) error {
	files, units, err := bootstraplogging.Files(config)
	if err != nil {
		return err
	}

	bootstrap.Config.Storage.Files = append(bootstrap.Config.Storage.Files, files...)
	bootstrap.Config.Systemd.Units = append(bootstrap.Config.Systemd.Units, units...)
	return nil
}

func (m *manager) addBootstrapImagePullLogging(bootstrap *bootstrap.Bootstrap) error {
	units, err := bootstrapimagepull.GetFiles()
	if err != nil {
		return err
	}

	bootstrap.Config.Systemd.Units = append(bootstrap.Config.Systemd.Units, units...)
	return nil
}
