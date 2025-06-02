package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/Azure/go-autorest/autorest/to"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/golang/mock/gomock"
	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	icazure "github.com/openshift/installer/pkg/asset/installconfig/azure"
	"github.com/openshift/installer/pkg/asset/installconfig/azure/mock"
	"github.com/openshift/installer/pkg/asset/releaseimage"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	"github.com/openshift/installer-aro-wrapper/pkg/api"
	"github.com/openshift/installer-aro-wrapper/pkg/bootstraplogging"
	"github.com/openshift/installer-aro-wrapper/pkg/env"
)

var expectedBootstrapStorageFileList = []string{"/etc/fluentbit/journal.conf",
	"/etc/sysconfig/fluentbit",

	"/etc/mdsd.d/mdsd.env",
	"/etc/mdsd.d/secret/mdsdcert.pem",
	"/etc/sysconfig/mdsd",

	"/etc/dnsmasq.conf",
	"/usr/local/bin/aro-dnsmasq-pre.sh",
	"/etc/NetworkManager/dispatcher.d/30-eth0-mtu-3900",

	"/etc/hosts.d/aro.conf",
	"/usr/local/bin/aro-etchosts-resolver.sh",

	"/opt/openshift/manifests/aro-imageregistry.yaml",
	"/opt/openshift/openshift/99_openshift-machineconfig_99-master-aro-dns.yaml",
	"/opt/openshift/openshift/99_openshift-machineconfig_99-master-aro-etc-hosts-gateway-domains.yaml",
	"/opt/openshift/openshift/99_openshift-machineconfig_99-worker-aro-dns.yaml",
	"/opt/openshift/openshift/99_openshift-machineconfig_99-worker-aro-etc-hosts-gateway-domains.yaml",

	"/opt/openshift/manifests/aro-ingress-namespace.yaml",
	"/opt/openshift/manifests/aro-ingress-service.yaml",
	"/opt/openshift/manifests/aro-worker-registries.yaml",
	"/opt/openshift/manifests/cluster-dns-02-config.yml",
	"/opt/openshift/openshift/99_openshift-cluster-api_master-user-data-secret.yaml",
	"/opt/openshift/openshift/99_openshift-cluster-api_worker-user-data-secret.yaml",
}

var expectedBootstrapSystemdFileList = []string{"fluentbit.service", "mdsd.service", "aro-etchosts-resolver.service", "dnsmasq.service"}

var apiIntIP = "203.0.113.1"
var expectedMasterIgnitionSource = "https://" + apiIntIP + ":22623/config/master"
var expectedWorkerIgnitionSource = "https://" + apiIntIP + ":22623/config/worker"
var expectedDNSConfigSource = `apiVersion: config.openshift.io/v1
kind: DNS
metadata:
  creationTimestamp: null
  name: cluster
spec:
  baseDomain: test-cluster.test.example.com
  platform:
    aws: null
    type: ""
status: {}
`

func fakeBootstrapLoggingConfig(_ env.Interface, _ *api.OpenShiftCluster) (*bootstraplogging.Config, error) {
	return &bootstraplogging.Config{
		Certificate:       "# This is not a real certificate",
		Key:               "# This is not a real private key", // notsecret
		Namespace:         "test-logging-namespace",
		Account:           "test-logging-account",
		Environment:       "test-logging-environment",
		ConfigVersion:     "42",
		Region:            "centralus",
		ResourceID:        "test-cluster-resource-id",
		SubscriptionID:    "test-subscription",
		ResourceName:      "test-logging-resource",
		ResourceGroupName: "test-resource-group",
		MdsdImage:         "registry.example.com/mdsd:latest",
		FluentbitImage:    "registry.example.com/fluentbit:latest",
	}, nil
}

func fakeGatewayDomains(_ env.Interface, _ *api.OpenShiftCluster) []string {
	return []string{
		"gateway.mock1.example.com",
		"gateway.mock2.example.com",
	}
}

func fakeCluster() *api.OpenShiftCluster {
	return &api.OpenShiftCluster{
		ID:   "test-cluster-resource-id",
		Name: "test-cluster",
		Properties: api.OpenShiftClusterProperties{
			InfraID:                         "test-infra-id",
			ImageRegistryStorageAccountName: "test-image-registry-storage-acct",
			APIServerProfile: api.APIServerProfile{
				IntIP: apiIntIP,
			},
			IngressProfiles: []api.IngressProfile{
				{
					IP: "192.0.2.1",
				},
			},
			NetworkProfile: api.NetworkProfile{
				GatewayPrivateEndpointIP: "203.0.113.2",
				MTUSize:                  api.MTU3900,
			},
		},
	}
}

func fakeManager() *manager {
	return &manager{
		clusterUUID:               "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
		log:                       logrus.NewEntry(logrus.StandardLogger()),
		oc:                        fakeCluster(),
		getBootstrapLoggingConfig: fakeBootstrapLoggingConfig,
		getGatewayDomains:         fakeGatewayDomains,
	}
}

func makeInstallConfig() *installconfig.InstallConfig {
	mpPlatform := types.MachinePoolPlatform{
		Azure: &azuretypes.MachinePool{
			Zones:            []string{"1", "2"},
			InstanceType:     "Standard_D2s_v3",
			EncryptionAtHost: true,
			VMNetworkingType: "Basic",
			OSDisk: azuretypes.OSDisk{
				DiskSizeGB: 1024,
			},
			OSImage: azuretypes.OSImage{
				Publisher: "azureopenshift",
				Offer:     "aro4",
				SKU:       "aro_417",
				Version:   "417.00.20240517",
				Plan:      azuretypes.ImageNoPurchasePlan,
			},
		},
	}

	return &installconfig.InstallConfig{
		AssetBase: installconfig.AssetBase{
			Config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-cluster",
				},
				BaseDomain: "test.example.com",
				Networking: &types.Networking{
					MachineNetwork: []types.MachineNetworkEntry{
						{
							CIDR: *ipnet.MustParseCIDR("127.0.0.0/8"),
						},
					},
					NetworkType: string(api.SoftwareDefinedNetworkOVNKubernetes),
					ClusterNetwork: []types.ClusterNetworkEntry{
						{
							CIDR:       *ipnet.MustParseCIDR("10.128.0.0/14"),
							HostPrefix: 23,
						},
					},
					ServiceNetwork: []ipnet.IPNet{
						*ipnet.MustParseCIDR("172.30.0.0/16"),
					},
				},
				ControlPlane: &types.MachinePool{
					Name:           "master",
					Replicas:       to.Int64Ptr(3),
					Platform:       mpPlatform,
					Hyperthreading: "Enabled",
					Architecture:   types.ArchitectureAMD64,
				},
				Compute: []types.MachinePool{
					{
						Name:           "worker",
						Replicas:       to.Int64Ptr(2),
						Platform:       mpPlatform,
						Hyperthreading: "Enabled",
						Architecture:   types.ArchitectureAMD64,
					},
				},
				Platform: types.Platform{
					Azure: &azuretypes.Platform{
						Region:                   "centralus",
						NetworkResourceGroupName: "test-nrg",
						VirtualNetwork:           "test-net",
						ControlPlaneSubnet:       "test-cp-subnet",
						ComputeSubnet:            "test-worker-subnet",
						CloudName:                "AzurePublicCloud",
						OutboundType:             azuretypes.LoadbalancerOutboundType,
						ResourceGroupName:        "test-resource-group",
					},
				},
				PullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}",
				FIPS:       false,
				ImageDigestSources: []types.ImageDigestSource{
					{
						Source: "quay.io/openshift-release-dev/ocp-release",
						Mirrors: []string{
							"registry.example.com/openshift-release-dev/ocp-release",
						},
					},
					{
						Source: "quay.io/openshift-release-dev/ocp-release-nightly",
						Mirrors: []string{
							"registry.example.com/openshift-release-dev/ocp-release-nightly",
						},
					},
					{
						Source: "quay.io/openshift-release-dev/ocp-v4.0-art-dev",
						Mirrors: []string{
							"registry.example.com/openshift-release-dev/ocp-v4.0-art-dev",
						},
					},
				},
				Publish: types.ExternalPublishingStrategy,
				Capabilities: &types.Capabilities{
					BaselineCapabilitySet: configv1.ClusterVersionCapabilitySetNone,
					AdditionalEnabledCapabilities: []configv1.ClusterVersionCapability{
						configv1.ClusterVersionCapabilityBuild,
						configv1.ClusterVersionCapabilityCloudCredential,
						configv1.ClusterVersionCapabilityConsole,
						configv1.ClusterVersionCapabilityCSISnapshot,
						configv1.ClusterVersionCapabilityDeploymentConfig,
						configv1.ClusterVersionCapabilityImageRegistry,
						configv1.ClusterVersionCapabilityInsights,
						configv1.ClusterVersionCapabilityMachineAPI,
						configv1.ClusterVersionCapabilityMarketplace,
						configv1.ClusterVersionCapabilityNodeTuning,
						configv1.ClusterVersionCapabilityOpenShiftSamples,
						configv1.ClusterVersionCapabilityOperatorLifecycleManager,
						configv1.ClusterVersionCapabilityStorage,
					},
				},
			},
		},
		Azure: &icazure.Metadata{
			CloudName:   "AzurePublicCloud",
			ARMEndpoint: "arm.example.com",
			Credentials: &icazure.Credentials{
				TenantID:       "test-tenant",
				SubscriptionID: "test-subscription",
				ClientID:       "test-client-id",
				ClientSecret:   "c3VwZXItc2VjcmV0", // notsecret
			},
		},
	}
}

func makeImage() *releaseimage.Image {
	return &releaseimage.Image{
		PullSpec: "quay.io/openshift-release-dev/ocp-release:4.17.0-x86_64",
	}
}

func mockClientCalls(client *mock.MockAPI) {
	client.EXPECT().GetVMCapabilities(gomock.Any(), "Standard_D2s_v3", "centralus").
		Return(map[string]string{
			"vCPUsAvailable":               "4",
			"MemoryGB":                     "16",
			"PremiumIO":                    "True",
			"HyperVGenerations":            "V1,V2",
			"AcceleratedNetworkingEnabled": "True",
			"CPUArchitectureType":          "x64",
		}, nil).
		AnyTimes()
	client.EXPECT().GetMarketplaceImage(gomock.Any(), "centralus", "azureopenshift", "aro4", "aro_417", "417.00.20240517").
		Return(compute.VirtualMachineImage{
			VirtualMachineImageProperties: &compute.VirtualMachineImageProperties{
				HyperVGeneration: compute.HyperVGenerationTypesV2,
			},
			Name:     to.StringPtr("aro_417"),
			Location: to.StringPtr("centralus"),
		}, nil).
		AnyTimes()
}

func TestApplyInstallConfigCustomisations(t *testing.T) {
	ctx := context.Background()
	m := fakeManager()
	inInstallConfig := makeInstallConfig()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockAPI(mockCtrl)
	inInstallConfig.Azure.UseMockClient(mockClient)
	mockClientCalls(mockClient)

	graph, err := m.applyInstallConfigCustomisations(ctx, inInstallConfig, makeImage())
	if err != nil {
		t.Fatal(err)
	}

	bootstrapAsset := graph.Get(&bootstrap.Bootstrap{}).(*bootstrap.Bootstrap)
	var temp map[string]any
	err = json.Unmarshal(bootstrapAsset.Files()[0].Data, &temp)
	if err != nil {
		t.Fatal(err)
	}
	verifyIgnitionFiles(t, temp, expectedBootstrapStorageFileList, expectedBootstrapSystemdFileList, bootstrapAsset.Files()[0].Filename)

	masterAsset := graph.Get(&machine.Master{}).(*machine.Master)
	workerAsset := graph.Get(&machine.Worker{}).(*machine.Worker)
	verifyMasterPointerIgnition(t, masterAsset.File.Data)
	verifyWorkerPointerIgnition(t, workerAsset.File.Data)
	verifyUpdateMCSCertKey(t, bootstrapAsset)
	verifyDNSPointerIgnition(t, bootstrapAsset)
}

func verifyIgnitionFiles(t *testing.T, temp map[string]any, storageFiles []string, systemdFiles []string, fileName string) {
	files := (temp["storage"].(map[string]any))["files"].([]any)
	systemd := (temp["systemd"].(map[string]any))["units"].([]any)
	storageFileList := map[string]string{}
	for _, file := range files {
		contents, found := file.(map[string]any)["contents"]
		if !found {
			contents = file.(map[string]any)["append"].([]any)[0]
		}
		storageFileList[file.(map[string]any)["path"].(string)] = contents.(map[string]any)["source"].(string)
	}
	systemdFileList := map[string]string{}
	for _, file := range systemd {
		contents, found := file.(map[string]any)["contents"]
		if !found {
			contents = file.(map[string]any)["dropins"].([]any)[0].(map[string]any)["contents"]
		}
		systemdFileList[file.(map[string]any)["name"].(string)] = contents.(string)
	}
	for _, file := range storageFiles {
		content, isFound := storageFileList[file]
		assert.True(t, isFound, fmt.Sprintf("file %v missing in storage file list in ignition file %s", file, fileName))
		if isFound {
			fileContents, err := base64.StdEncoding.DecodeString(strings.Split(content, "base64")[1][1:])
			if err != nil {
				t.Fatal(err)
			}
			if file == "/opt/openshift/manifests/aro-imageregistry.yaml" {
				content := string(fileContents)
				re := regexp.MustCompile(`httpSecret: "[A-Za-z0-9]+"`)
				fileContents = []byte(re.ReplaceAllString(content, `httpSecret: "test"`))
			} else if strings.Contains(file, "-user-data-secret.yaml") {
				content := string(fileContents)
				re := regexp.MustCompile(`userData: .*`)
				userData := re.FindString(string(fileContents))

				innerContent, err := base64.StdEncoding.DecodeString(strings.Split(userData, "userData: ")[1])
				if err != nil {
					t.Fatal(err)
				}

				assert.Contains(t, string(innerContent), "https://203.0.113.1:22623/config/")
				fileContents = []byte(re.ReplaceAllString(content, `userData: test`))
			}
			assert.EqualValues(t, expectedIgnitionFileContents[file], string(fileContents), fmt.Sprintf("missing storage data in file %v", file))
		}
	}
	for _, file := range systemdFiles {
		content, isFound := systemdFileList[file]
		assert.True(t, isFound, fmt.Sprintf("file %v missing from systemd file list in ignition file %s", file, fileName))
		if isFound {
			assert.EqualValues(t, expectedIgnitionServiceContents[file], content, fmt.Sprintf("missing systemd data in file %v", file))
		}
	}
	installConfigMap, err := base64.StdEncoding.DecodeString(strings.Split(storageFileList["/opt/openshift/openshift/openshift-install-manifests.yaml"], ",")[1])
	if err != nil {
		t.Fatal(err)
	}
	var config corev1.ConfigMap
	err = yaml.Unmarshal(installConfigMap, &config)
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValues(t, "ARO", config.Data["invoker"])
}

func verifyMasterPointerIgnition(t *testing.T, ignData []byte) {
	ignContents := &igntypes.Config{}
	err := json.Unmarshal(ignData, &ignContents)
	if err != nil {
		t.Fatal(err)
	}

	actualSource := *ignContents.Ignition.Config.Merge[0].Source
	assert.EqualValues(t, expectedMasterIgnitionSource, actualSource, fmt.Sprintf("expected master pointer ignition to be %s but found %s", expectedMasterIgnitionSource, actualSource))
}

func verifyWorkerPointerIgnition(t *testing.T, ignData []byte) {
	ignContents := &igntypes.Config{}
	err := json.Unmarshal(ignData, &ignContents)
	if err != nil {
		t.Fatal(err)
	}

	actualSource := *ignContents.Ignition.Config.Merge[0].Source
	assert.EqualValues(t, expectedWorkerIgnitionSource, actualSource, fmt.Sprintf("expected worker pointer ignition to be %s but found %s", expectedWorkerIgnitionSource, actualSource))
}

func verifyDNSPointerIgnition(t *testing.T, bootstrap *bootstrap.Bootstrap) {
	filesNeeded := []string{dnsCfgFilename}
	mapFiles := map[string]*string{}
	for _, key := range filesNeeded {
		mapFiles[key] = to.StringPtr("")
	}
	for _, file := range bootstrap.Config.Storage.Files {
		if _, ok := mapFiles[file.Path]; ok {
			mapFiles[file.Path] = file.Contents.Source
		}
	}
	for _, key := range filesNeeded {
		if *mapFiles[key] == "" {
			assert.Failf(t, "file %s content missing from bootstrap data", key)
		}
	}
	dnsConfigFile := mapFiles[dnsCfgFilename]
	var config configv1.DNS
	b, err := base64.StdEncoding.DecodeString(strings.Split(*dnsConfigFile, ",")[1])
	if err != nil {
		t.Fatal(err)
	}
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValues(t, expectedDNSConfigSource, string(b), fmt.Sprintf("expected dns config to be %s but found %s", expectedDNSConfigSource, string(b)))
}

func verifyUpdateMCSCertKey(t *testing.T, bootstrap *bootstrap.Bootstrap) {
	config := &igntypes.Config{}
	config = bootstrap.Config

	cert := &x509.Certificate{}
	var rawCert, rawKey []byte

	for i, fileData := range config.Storage.Files {
		if fileData.Path == mcsCertFile {
			contents := strings.Split(*config.Storage.Files[i].Contents.Source, ",")
			decodedCert, err := base64.StdEncoding.DecodeString(contents[1])
			if err != nil {
				t.Fatal(err)
			}
			rawCert = decodedCert
			assert.NotNil(t, rawCert)
			cert, err = tls.PemToCertificate(decodedCert)
			if err != nil {
				t.Fatal(err)
			}
			certPool := x509.NewCertPool()
			if !certPool.AppendCertsFromPEM(rawCert) {
				t.Error("failed to append certs from PEM")
			}
			opts := x509.VerifyOptions{
				Roots:   certPool,
				DNSName: apiIntIP,
			}
			_, err = cert.Verify(opts)
			assert.NoError(t, err, "verifyUpdateMCSCertKey")
		}
		if fileData.Path == mcsKeyFile {
			contents := strings.Split(*config.Storage.Files[i].Contents.Source, ",")
			decodedKey, err := base64.StdEncoding.DecodeString(contents[1])
			if err != nil {
				t.Fatal(err)
			}
			rawKey = decodedKey
			assert.NotNil(t, rawKey)
		}
	}
	for i, fileData := range config.Storage.Files {
		if fileData.Path == mcsCertKeyFilepath {
			contents := strings.Split(*config.Storage.Files[i].Contents.Source, ",")
			rawDecodedText, err := base64.StdEncoding.DecodeString(contents[1])
			if err != nil {
				t.Fatal(err)
			}
			mcsSecret := &corev1.Secret{}
			if err := yaml.Unmarshal(rawDecodedText, mcsSecret); err != nil {
				t.Fatal(err)
			}
			assert.EqualValues(t, rawCert, mcsSecret.Data[corev1.TLSCertKey], fmt.Sprintf("mismatched raw certs in %s and %s", mcsCertKeyFilepath, mcsCertFile))
			assert.EqualValues(t, rawKey, mcsSecret.Data[corev1.TLSPrivateKeyKey], fmt.Sprintf("mismatched raw private key in %s and %s", mcsCertKeyFilepath, mcsKeyFile))
		}
	}
}
