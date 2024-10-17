package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"path/filepath"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/releaseimage"
	"github.com/openshift/installer/pkg/asset/targets"

	"github.com/openshift/installer-aro-wrapper/pkg/api"
	"github.com/openshift/installer-aro-wrapper/pkg/aro"
	"github.com/openshift/installer-aro-wrapper/pkg/bootstraplogging"
	"github.com/openshift/installer-aro-wrapper/pkg/cluster/graph"
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

	imageRegistryConfig := &aro.AROImageRegistryConfig{
		AccountName:   m.oc.Properties.ImageRegistryStorageAccountName,
		ContainerName: "image-registry",
		HTTPSecret:    hex.EncodeToString(httpSecret),
	}

	dnsConfig := &aro.ARODNSConfig{
		APIIntIP:  m.oc.Properties.APIServerProfile.IntIP,
		IngressIP: m.oc.Properties.IngressProfiles[0].IP,
	}

	ingressConfig := &aro.AROIngressService{}
	workerRegistry := &aro.AROWorkerRegistries{}

	if m.oc.Properties.NetworkProfile.GatewayPrivateEndpointIP != "" {
		dnsConfig.GatewayPrivateEndpointIP = m.oc.Properties.NetworkProfile.GatewayPrivateEndpointIP
		dnsConfig.GatewayDomains = append(m.env.GatewayDomains(), m.oc.Properties.ImageRegistryStorageAccountName+".blob."+m.env.Environment().StorageEndpointSuffix)
	}

	fileFetcher := &aroFileFetcher{directory: "/"}

	aroManifests := &AROManifests{}
	aroManifestsExist, err := aroManifests.Load(fileFetcher)
	if err != nil {
		err = fmt.Errorf("error loading ARO manifests: %w", err)
		m.log.Error(err)
		return nil, err
	}

	boundSaSigningKey := &AROBoundSASigningKey{}
	_, err = boundSaSigningKey.Load(fileFetcher)
	if err != nil {
		err = fmt.Errorf("error loading boundSASigningKey: %w", err)
		m.log.Error(err)
		return nil, err
	}

	g := graph.Graph{}
	g.Set(installConfig, image, clusterID, bootstrapLoggingConfig, dnsConfig, imageRegistryConfig, &boundSaSigningKey.BoundSASigningKey, ingressConfig, workerRegistry)

	m.log.Print("resolving graph")
	for _, a := range targets.Cluster {
		err = g.Resolve(a)
		if err != nil {
			return nil, err
		}
	}

	// Handle MTU3900 feature flag
	if m.oc.Properties.NetworkProfile.MTUSize == api.MTU3900 {
		m.log.Printf("applying feature flag %s", api.FeatureFlagMTU3900)
		if err = m.overrideEthernetMTU(g); err != nil {
			return nil, err
		}
	}

	// Add ARO Manifests to bootstrap Files
	if aroManifestsExist {
		if err = appendFilesToBootstrap(aroManifests, g); err != nil {
			return nil, err
		}
	}

	return g, nil
}

func appendFilesToBootstrap(a asset.WritableAsset, g graph.Graph) error {
	bootstrap := g.Get(&bootstrap.Bootstrap{}).(*bootstrap.Bootstrap)
	for _, file := range a.Files() {
		manifest := ignition.FileFromBytes(filepath.Join(rootPath, file.Filename), "root", 0644, file.Data)
		bootstrap.Config.Storage.Files = append(bootstrap.Config.Storage.Files, manifest)
	}

	data, err := ignition.Marshal(bootstrap.Config)
	if err != nil {
		return err
	}
	bootstrap.File.Data = data
	return nil
}
