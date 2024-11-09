package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"net/url"
	"path/filepath"

	"github.com/coreos/ignition/v2/config/util"
	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/cluster"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/kubeconfig"
	"github.com/openshift/installer/pkg/asset/password"
	"github.com/openshift/installer/pkg/asset/releaseimage"
	"github.com/openshift/installer/pkg/asset/templates/content/bootkube"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"

	"github.com/openshift/installer-aro-wrapper/pkg/api"
	"github.com/openshift/installer-aro-wrapper/pkg/cluster/graph"
	"github.com/openshift/installer-aro-wrapper/pkg/data/manifests"
	"github.com/openshift/installer-aro-wrapper/pkg/installer/dnsmasq"
	"github.com/openshift/installer-aro-wrapper/pkg/installer/etchost"
	"github.com/openshift/installer-aro-wrapper/pkg/installer/mdsd"
)

const (
	cvoOverridesFilename = "manifests/cvo-overrides.yaml"
)

var (
	targetAssets = []asset.WritableAsset{
		&cluster.Metadata{},
		&machine.MasterIgnitionCustomizations{},
		&machine.WorkerIgnitionCustomizations{},
		&cluster.TerraformVariables{},
		&kubeconfig.AdminClient{},
		&password.KubeadminPassword{},
		&tls.JournalCertKey{},
		&tls.RootCA{},
	}
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

	localdnsConfig := dnsmasq.DNSConfig{
		APIIntIP:  m.oc.Properties.APIServerProfile.IntIP,
		IngressIP: m.oc.Properties.IngressProfiles[0].IP,
	}

	dnsConfig := &bootkube.ARODNSConfig{
		APIIntIP:  m.oc.Properties.APIServerProfile.IntIP,
		IngressIP: m.oc.Properties.IngressProfiles[0].IP,
	}

	if m.oc.Properties.NetworkProfile.GatewayPrivateEndpointIP != "" {
		dnsConfig.GatewayPrivateEndpointIP = m.oc.Properties.NetworkProfile.GatewayPrivateEndpointIP
		dnsConfig.GatewayDomains = m.getGatewayDomains(m.env, m.oc)
		localdnsConfig.GatewayPrivateEndpointIP = m.oc.Properties.NetworkProfile.GatewayPrivateEndpointIP
		localdnsConfig.GatewayDomains = m.getGatewayDomains(m.env, m.oc)
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
	g.Set(installConfig, image, clusterID, dnsConfig, imageRegistryConfig, &boundSaSigningKey.BoundSASigningKey)

	m.log.Print("resolving graph")
	for _, a := range targetAssets {
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

	bootstrapAsset := g.Get(&bootstrap.Bootstrap{}).(*bootstrap.Bootstrap)
	err = dnsmasq.CreatednsmasqIgnitionFiles(bootstrapAsset, installConfig, localdnsConfig)
	if err != nil {
		return nil, err
	}
	err = mdsd.AppendMdsdFiles(bootstrapAsset, bootstrapLoggingConfig)
	if err != nil {
		return nil, err
	}
	config := manifests.ManifestsConfig{
		AROWorkerRegistries: manifests.AroWorkerRegistries(installConfig.Config.ImageDigestSources),
		HTTPSecret:          imageRegistryConfig.HTTPSecret,
		AccountName:         imageRegistryConfig.AccountName,
		ContainerName:       imageRegistryConfig.ContainerName,
		CloudName:           installConfig.Config.Azure.CloudName.Name(),
		AROIngressInternal:  installConfig.Config.Publish == "Internal",
		AROIngressIP:        dnsConfig.IngressIP,
	}
	err = manifests.AppendManifestsFilesToBootstrap(bootstrapAsset, config)
	if err != nil {
		return nil, err
	}
	err = etchost.AppendEtcHostFiles(bootstrapAsset, *installConfig, localdnsConfig)
	if err != nil {
		return nil, err
	}
	// Update Master and Worker Pointer Ignition with ARO API-Int IP
	if err = replacePointerIgnition(aroManifests, g, &localdnsConfig); err != nil {
		return nil, err
	}
	// Update machine-confog-server cert to allow connecting with API-Int LB IP
	if err = updateMCSCertKey(g, installConfig, &localdnsConfig); err != nil {
		return nil, err
	}
	data, err := ignition.Marshal(bootstrapAsset.Config)
	if err != nil {
		return nil, err
	}
	bootstrapAsset.File.Data = data
	g.Set(bootstrapAsset)

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

// replacePointerIgnition performs the same functionality as the upstream
// installer's pointerIgnitionConfig() but with ARO specific DNS config
func replacePointerIgnition(a asset.WritableAsset, g graph.Graph, localdnsConfig *dnsmasq.DNSConfig) (err error) {
	masterPointerIgn := g.Get(&machine.Master{}).(*machine.Master)
	workerPointerIgn := g.Get(&machine.Worker{}).(*machine.Worker)
	ignitionHost := net.JoinHostPort(localdnsConfig.APIIntIP, "22623")
	masterRole := "master"
	workerRole := "worker"

	masterPointerIgn.Config.Ignition.Config.Merge[0].Source = util.StrToPtr(func() *url.URL {
		return &url.URL{
			Scheme: "https",
			Host:   ignitionHost,
			Path:   fmt.Sprintf("/config/%s", masterRole),
		}
	}().String())

	workerPointerIgn.Config.Ignition.Config.Merge[0].Source = util.StrToPtr(func() *url.URL {
		return &url.URL{
			Scheme: "https",
			Host:   ignitionHost,
			Path:   fmt.Sprintf("/config/%s", workerRole),
		}
	}().String())

	data, err := ignition.Marshal(masterPointerIgn.Config)
	if err != nil {
		return errors.Wrap(err, "failed to marshal updated master pointer Ignition config")
	}

	masterPointerIgn.File.Data = data

	data, err = ignition.Marshal(workerPointerIgn.Config)
	if err != nil {
		return errors.Wrap(err, "failed to marshal updated worker pointer Ignition config")
	}
	workerPointerIgn.File.Data = data

	return nil
}
