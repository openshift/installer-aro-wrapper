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
	block := []byte("-----BEGIN RSA PRIVATE KEY-----\n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\n-----END RSA PRIVATE KEY-----")
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

	if err = m.testSecrets(g); err != nil {
		return nil, err
	}

	return g, nil
}
