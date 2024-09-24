package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"crypto/rand"
	"encoding/hex"
	"path/filepath"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/releaseimage"
	"github.com/openshift/installer/pkg/asset/targets"
	"github.com/openshift/installer/pkg/asset/templates/content/bootkube"
	"github.com/openshift/installer/pkg/asset/tls"

	"github.com/openshift/ARO-Installer/pkg/api"
	"github.com/openshift/ARO-Installer/pkg/bootstraplogging"
	"github.com/openshift/ARO-Installer/pkg/cluster/graph"
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

	block := []byte(*m.oc.Properties.ClusterProfile.BoundServiceAccountSigningKey)
	rsaKey, err := tls.PemToPrivateKey(block)
	if err != nil {
		return nil, err
	}
	pubData, err := tls.PublicKeyToPem(&rsaKey.PublicKey)
	if err != nil {
		return nil, err
	}

	boundSASigningKey := &tls.BoundSASigningKey{
		FileList: []*asset.File{
			{
				Filename: filepath.Join("tls", "bound-service-account-signing-key.key"),
				Data:     block,
			},
			{
				Filename: filepath.Join("tls", "bound-service-account-signing-key.pub"),
				Data:     pubData,
			},
		},
	}

	g := graph.Graph{}
	g.Set(installConfig, image, clusterID, bootstrapLoggingConfig, dnsConfig, imageRegistryConfig, boundSASigningKey)

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

	return g, nil
}
