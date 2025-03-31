package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net"
	"net/url"
	"path/filepath"
	"text/template"

	"github.com/coreos/ignition/v2/config/util"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/releaseimage"
	targetassets "github.com/openshift/installer/pkg/asset/targets"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"

	"github.com/openshift/installer-aro-wrapper/pkg/api"
	"github.com/openshift/installer-aro-wrapper/pkg/cluster/graph"
	bootstrapfiles "github.com/openshift/installer-aro-wrapper/pkg/data/bootstrap"
	"github.com/openshift/installer-aro-wrapper/pkg/data/manifests"
	"github.com/openshift/installer-aro-wrapper/pkg/installer/dnsmasq"
	"github.com/openshift/installer-aro-wrapper/pkg/installer/etchost"
	"github.com/openshift/installer-aro-wrapper/pkg/installer/mdsd"
)

const (
	cvoOverridesFilename = "manifests/cvo-overrides.yaml"
)

var (
	userDataTmpl = template.Must(template.New("user-data").Parse(`apiVersion: v1
kind: Secret
metadata:
  name: {{.name}}
  namespace: openshift-machine-api
type: Opaque
data:
  disableTemplating: "dHJ1ZQo="
  userData: {{.content}}
`))
	dnsCfgFilename = filepath.Join(rootPath, "manifests", "cluster-dns-02-config.yml")
)

// applyInstallConfigCustomisations modifies the InstallConfig and creates
// parent assets, then regenerates the InstallConfig for use for Ignition
// generation, etc.
func (m *manager) applyInstallConfigCustomisations(ctx context.Context, installConfig *installconfig.InstallConfig, image *releaseimage.Image) (graph.Graph, error) {
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

	imageRegistryConfig := struct {
		AccountName   string
		ContainerName string
		HTTPSecret    string
	}{
		AccountName:   m.oc.Properties.ImageRegistryStorageAccountName,
		ContainerName: "image-registry",
		HTTPSecret:    hex.EncodeToString(httpSecret),
	}

	localdnsConfig := dnsmasq.DNSConfig{
		APIIntIP:  m.oc.Properties.APIServerProfile.IntIP,
		IngressIP: m.oc.Properties.IngressProfiles[0].IP,
	}

	if m.oc.Properties.NetworkProfile.GatewayPrivateEndpointIP != "" {
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
	g.Set(installConfig, image, clusterID, &boundSaSigningKey.BoundSASigningKey)

	m.log.Print("resolving graph")
	for _, a := range targetassets.IgnitionConfigs {
		err = g.Resolve(ctx, a)
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
		AROIngressIP:        localdnsConfig.IngressIP,
	}
	err = manifests.AppendManifestsFilesToBootstrap(bootstrapAsset, config)
	if err != nil {
		return nil, err
	}
	err = etchost.AppendEtcHostFiles(bootstrapAsset, *installConfig, localdnsConfig)
	if err != nil {
		return nil, err
	}
	err = removeDNSConfigData(bootstrapAsset, *installConfig)
	if err != nil {
		return nil, err
	}
	// Update Master and Worker Pointer Ignition with ARO API-Int IP
	if err = replacePointerIgnition(bootstrapAsset, g, &localdnsConfig); err != nil {
		return nil, err
	}
	// Update machine-config-server cert to allow connecting with API-Int LB IP
	if err = updateMCSCertKey(g, installConfig, &localdnsConfig); err != nil {
		return nil, err
	}
	data, err := ignition.Marshal(bootstrapAsset.Config)
	if err != nil {
		return nil, err
	}
	bootstrapAsset.File.Data = data
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

func removeDNSConfigData(bootstrap *bootstrap.Bootstrap, installConfig installconfig.InstallConfig) error {
	dns := &configv1.DNS{
		TypeMeta: metav1.TypeMeta{
			APIVersion: configv1.SchemeGroupVersion.String(),
			Kind:       "DNS",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster",
			// not namespaced
		},
		Spec: configv1.DNSSpec{
			BaseDomain: installConfig.Config.ClusterDomain(),
		},
	}
	data, err := yaml.Marshal(dns)
	if err != nil {
		return err
	}
	config := ignition.FileFromBytes(dnsCfgFilename, "root", 0644, data)
	bootstrap.Config.Storage.Files = bootstrapfiles.ReplaceOrAppend(bootstrap.Config.Storage.Files, []igntypes.File{config})
	return nil
}

// replacePointerIgnition performs the same functionality as the upstream
// installer's pointerIgnitionConfig() but with ARO specific DNS config
func replacePointerIgnition(a *bootstrap.Bootstrap, g graph.Graph, localdnsConfig *dnsmasq.DNSConfig) (err error) {
	masterPointerIgn := g.Get(&machine.Master{}).(*machine.Master)
	workerPointerIgn := g.Get(&machine.Worker{}).(*machine.Worker)
	ignitionHost := net.JoinHostPort(localdnsConfig.APIIntIP, "22623")

	masterPointerIgn.Config.Ignition.Config.Merge[0].Source = util.StrToPtr(func() *url.URL {
		return &url.URL{
			Scheme: "https",
			Host:   ignitionHost,
			Path:   "/config/master",
		}
	}().String())

	workerPointerIgn.Config.Ignition.Config.Merge[0].Source = util.StrToPtr(func() *url.URL {
		return &url.URL{
			Scheme: "https",
			Host:   ignitionHost,
			Path:   "/config/worker",
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

	// Update the user-data information for the machine.
	//
	// Note, we do not need to update asset/machines/.(Master/Worker)'s
	// UserDataSecret since it is only used at first generation (which has
	// already by the time we have got here) and we can just update the files it
	// would have generated directly. We also do not need to generate the
	// "override" MachineConfig which would amend the host.
	//
	// This is also done for masters since it will likely be used if Control
	// Plane Machine Sets creates a new control plane node.
	masterUserDataPath := filepath.Join("openshift", "99_openshift-cluster-api_master-user-data-secret.yaml")
	workerUserDataPath := filepath.Join("openshift", "99_openshift-cluster-api_worker-user-data-secret.yaml")

	masterData, err := userDataSecret("master-user-data", masterPointerIgn.File.Data)
	if err != nil {
		return errors.Wrap(err, "failed to create user-data secret for master machines")
	}
	workerData, err := userDataSecret("worker-user-data", workerPointerIgn.File.Data)
	if err != nil {
		return errors.Wrap(err, "failed to create user-data secret for worker machines")
	}

	a.Config.Storage.Files = bootstrapfiles.ReplaceOrAppend(a.Config.Storage.Files, []igntypes.File{
		ignition.FileFromBytes(filepath.Join(rootPath, masterUserDataPath), "root", 0644, masterData),
		ignition.FileFromBytes(filepath.Join(rootPath, workerUserDataPath), "root", 0644, workerData),
	})
	return nil
}

func userDataSecret(name string, content []byte) ([]byte, error) {
	encodedData := map[string]string{
		"name":    name,
		"content": base64.StdEncoding.EncodeToString(content),
	}
	buf := &bytes.Buffer{}
	if err := userDataTmpl.Execute(buf, encodedData); err != nil {
		return nil, errors.Wrap(err, "failed to execute user-data template")
	}
	return buf.Bytes(), nil
}
