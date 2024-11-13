package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/2018-03-01/resources/mgmt/resources"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/golang/mock/gomock"
	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset/bootstraplogging"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/installconfig"
	icazure "github.com/openshift/installer/pkg/asset/installconfig/azure"
	"github.com/openshift/installer/pkg/asset/installconfig/azure/mock"
	"github.com/openshift/installer/pkg/asset/releaseimage"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer-aro-wrapper/pkg/api"
	"github.com/openshift/installer-aro-wrapper/pkg/env"
)

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
				IntIP: "203.0.113.1",
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
				SKU:       "aro_416",
				Version:   "416.00.20240517",
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
		PullSpec: "quay.io/openshift-release-dev/ocp-release:4.16.0-x86_64",
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
	client.EXPECT().GetMarketplaceImage(gomock.Any(), "centralus", "azureopenshift", "aro4", "aro_416", "416.00.20240517").
		Return(compute.VirtualMachineImage{
			VirtualMachineImageProperties: &compute.VirtualMachineImageProperties{
				HyperVGeneration: compute.HyperVGenerationTypesV2,
			},
			Name:     to.StringPtr("aro_416"),
			Location: to.StringPtr("centralus"),
		}, nil).
		AnyTimes()
	client.EXPECT().GetGroup(gomock.Any(), "test-resource-group").
		Return(&resources.Group{
			ID:       to.StringPtr("test-resource-group"),
			Location: to.StringPtr("centralus"),
		}, nil)
	client.EXPECT().GetHyperVGenerationVersion(gomock.Any(), "Standard_D2s_v3", "centralus", "").
		Return("V2", nil)
}

func TestApplyInstallConfigCustomisations(t *testing.T) {
	m := fakeManager()
	inInstallConfig := makeInstallConfig()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockClient := mock.NewMockAPI(mockCtrl)
	inInstallConfig.Azure.UseMockClient(mockClient)
	mockClientCalls(mockClient)

	graph, err := m.applyInstallConfigCustomisations(inInstallConfig, makeImage())
	if err != nil {
		t.Fatal(err)
	}

	bootstrapAsset := graph.Get(&bootstrap.Bootstrap{}).(*bootstrap.Bootstrap)
	bootstrapIgnition := string(bootstrapAsset.Files()[0].Data)

	t.Log(bootstrapIgnition)
}
