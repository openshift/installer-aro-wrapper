package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"path/filepath"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/releaseimage"
	"github.com/openshift/installer/pkg/asset/targets"
	"github.com/openshift/installer/pkg/asset/templates/content/bootkube"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"

	"github.com/openshift/installer-aro-wrapper/pkg/api"
	"github.com/openshift/installer-aro-wrapper/pkg/cluster/graph"
)

const (
	cvoOverridesFilename = "manifests/cvo-overrides.yaml"
)

// applyInstallConfigCustomisations modifies the InstallConfig and creates
// parent assets, then regenerates the InstallConfig for use for Ignition
// generation, etc.
func (m *manager) applyInstallConfigCustomisations(installConfig *installconfig.InstallConfig, image *releaseimage.Image) (graph.Graph, error) {
	clusterID := &installconfig.ClusterID{
		UUID:    m.clusterUUID,
		InfraID: m.oc.Properties.InfraID,
	}

	bootstrapLoggingConfig, err := m.getBootstrapLoggingConfig(m.env, m.oc)
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
		dnsConfig.GatewayDomains = m.getGatewayDomains(m.env, m.oc)
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
	g.Set(installConfig, image, clusterID, bootstrapLoggingConfig, dnsConfig, imageRegistryConfig, &boundSaSigningKey.BoundSASigningKey)

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

	// Add ARO Manifests to bootstrap Files and CVO Overrides
	if aroManifestsExist {
		if err = appendFilesToCvoOverrides(aroManifests, g); err != nil {
			return nil, err
		}

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

// appendFilesToCvoOverides performs the same functionality as the upstream
// installer's CVOIgnore asset (pkg/asset/ignition/bootstrap/cvoignore.go),
// but for our custom AROManifests asset.
func appendFilesToCvoOverrides(a asset.WritableAsset, g graph.Graph) (err error) {
	cvoIgnore := g.Get(&bootstrap.CVOIgnore{}).(*bootstrap.CVOIgnore)
	bootstrap := g.Get(&bootstrap.Bootstrap{}).(*bootstrap.Bootstrap)

	var ignoredResources []configv1.ComponentOverride
	files := a.Files()
	seen := make(map[string]string, len(files))

	for _, file := range files {
		u := &unstructured.Unstructured{}
		if err := yaml.Unmarshal(file.Data, u); err != nil {
			return errors.Wrapf(err, "could not unmarshal %q", file.Filename)
		}

		group := u.GetObjectKind().GroupVersionKind().Group
		kind := u.GetKind()
		namespace := u.GetNamespace()
		name := u.GetName()

		key := fmt.Sprintf("%s |! %s |! %s |! %s", group, kind, namespace, name)
		if previousFile, ok := seen[key]; ok {
			return fmt.Errorf("multiple manifests for group %s kind %s namespace %s name %s: %s, %s", group, kind, namespace, name, previousFile, file.Filename)
		}
		seen[key] = file.Filename

		ignoredResources = append(ignoredResources,
			configv1.ComponentOverride{
				Kind:      kind,
				Group:     group,
				Namespace: namespace,
				Name:      name,
				Unmanaged: true,
			})
	}

	clusterVersion := &configv1.ClusterVersion{}
	var cvData []byte
	for i, file := range cvoIgnore.Files() {
		if file.Filename != cvoOverridesFilename {
			continue
		}

		if err := yaml.Unmarshal(file.Data, clusterVersion); err != nil {
			return errors.Wrapf(err, "could not unmarshal %q", file.Filename)
		}

		clusterVersion.Spec.Overrides = append(clusterVersion.Spec.Overrides, ignoredResources...)

		cvData, err = yaml.Marshal(clusterVersion)
		if err != nil {
			return errors.Wrap(err, "error marshalling clusterversion")
		}
		cvoIgnore.FileList[i] = &asset.File{
			Filename: file.Filename,
			Data:     cvData,
		}
	}

	ignPath := filepath.Join(rootPath, cvoOverridesFilename)
	for i, file := range bootstrap.Config.Storage.Files {
		if file.Path != ignPath {
			continue
		}

		bootstrap.Config.Storage.Files[i] = ignition.FileFromBytes(ignPath, "root", 0420, cvData)
	}

	return nil
}
