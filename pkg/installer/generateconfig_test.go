package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	mgmtcompute "github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/openshift/installer-aro-wrapper/pkg/api"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/cluster"
	targetassets "github.com/openshift/installer/pkg/asset/targets"
	aztypes "github.com/openshift/installer/pkg/types/azure"

	"github.com/sirupsen/logrus"
)

func TestVMNetworkingType(t *testing.T) {
	capabilityName := aztypes.AcceleratedNetworkingEnabled
	for _, tt := range []struct {
		name     string
		sku      *mgmtcompute.ResourceSku
		wantType string
	}{
		{
			name: "sku with support for accelerated networking",
			sku: &mgmtcompute.ResourceSku{
				Capabilities: &([]mgmtcompute.ResourceSkuCapabilities{
					{Name: &capabilityName, Value: to.StringPtr("True")},
				}),
			},
			wantType: "Accelerated",
		}, {
			name: "sku without support for accelerated networking",
			sku: &mgmtcompute.ResourceSku{
				Capabilities: &([]mgmtcompute.ResourceSkuCapabilities{
					{Name: &capabilityName, Value: to.StringPtr("False")},
				}),
			},
			wantType: "Basic",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			result := determineVMNetworkingType(tt.sku)

			if result != tt.wantType {
				t.Error(result)
			}
		})
	}
}

func TestApplyInstallConfigCustomisations(t *testing.T) {
	ctx := context.Background()
	log := logrus.NewEntry(logrus.StandardLogger())

	log.Info("Verifying environment variables are set")
	for _, key := range []string{
		"ARO_AZURE_SUBSCRIPTION_ID",
		"ARO_AZURE_TENANT_ID",
		"ARO_RESOURCEGROUP",
		"ARO_RESOURCEGROUP_ID",
		"ARO_LOCATION",
		"PULL_SECRET",
		"DOMAIN_NAME",
		"ARO_AZURE_RP_CLIENT_ID",
		"ARO_AZURE_RP_CLIENT_SECRET",
		"ARO_BASE_PATH",
		"ARO_ACCOUNT_OWNER",
		"ARO_SSH_PRIVATE_KEY_PATH",
		"ARO_VNET_NAME",
		"ARO_MASTER_SUBNET_NAME",
		"ARO_WORKER_SUBNET_NAME",
	} {
		if _, found := os.LookupEnv(key); !found {
			t.Fatalf("environment variable %q unset", key)
		}
	}

	log.Info("Getting values from environment variables")
	subscriptionID := os.Getenv("ARO_AZURE_SUBSCRIPTION_ID")
	tenantID := os.Getenv("ARO_AZURE_TENANT_ID")
	resourceGroup := os.Getenv("ARO_RESOURCEGROUP")
	resourceGroupID := os.Getenv("ARO_RESOURCEGROUP_ID")
	location := os.Getenv("ARO_LOCATION")
	pullSecret := os.Getenv("PULL_SECRET")
	domain := os.Getenv("DOMAIN_NAME")
	clientID := os.Getenv("ARO_AZURE_RP_CLIENT_ID")
	clientSecret := os.Getenv("ARO_AZURE_RP_CLIENT_SECRET")
	aroBasePath := os.Getenv("ARO_BASE_PATH")
	email := os.Getenv("ARO_ACCOUNT_OWNER")
	sshPrivateKeyPath := os.Getenv("ARO_SSH_PRIVATE_KEY_PATH")
	vnet := os.Getenv("ARO_VNET_NAME")
	masterSubnet := os.Getenv("ARO_MASTER_SUBNET_NAME")
	workerSubnet := os.Getenv("ARO_WORKER_SUBNET_NAME")

	assetsDirectory := filepath.Join(aroBasePath, "assets")

	log.Info("Setting up resources")
	ocResource := azure.Resource{
		SubscriptionID: subscriptionID,
		ResourceGroup:  resourceGroup,
		Provider:       "test-provider",
		ResourceType:   "test-resource-type",
		ResourceName:   "test-resource-name",
	}
	masterSubnetResource := azure.Resource{
		SubscriptionID: subscriptionID,
		ResourceGroup:  resourceGroup,
		Provider:       "Microsoft.Network",
		ResourceType:   fmt.Sprintf("virtualNetworks/%s/subnets", vnet),
		ResourceName:   masterSubnet,
	}
	workerSubnetResource := azure.Resource{
		SubscriptionID: subscriptionID,
		ResourceGroup:  resourceGroup,
		Provider:       "Microsoft.Network",
		ResourceType:   fmt.Sprintf("virtualNetworks/%s/subnets", vnet),
		ResourceName:   workerSubnet,
	}

	log.Info("Reading SSH private key")
	sshKeyPem, err := os.ReadFile(sshPrivateKeyPath)
	if err != nil {
		t.Fatalf("failed to read ssh private key: %v", err)
	}

	log.Info("Decoding SSH private key")
	sshKeyDer, _ := pem.Decode(sshKeyPem)
	if sshKeyDer == nil {
		t.Fatal("failed to decode ssh private key")
	}

	log.Info("Setting up OpenShiftCluster")
	var oc api.OpenShiftCluster = api.OpenShiftCluster{
		ID:       ocResource.String(),
		Name:     "test-cluster-name",
		Type:     "test-cluster-type",
		Location: location,
		Properties: api.OpenShiftClusterProperties{
			InfraID: fmt.Sprintf("%s-infra-id", os.Getenv("USER")),
			ClusterProfile: api.ClusterProfile{
				ResourceGroupID: resourceGroupID,
				PullSecret:      api.SecureString(pullSecret),
				Domain:          domain,
			},
			ServicePrincipalProfile: &api.ServicePrincipalProfile{
				ClientID:     clientID,
				ClientSecret: api.SecureString(clientSecret),
				SPObjectID:   "test-spobject-id",
			},
			NetworkProfile: api.NetworkProfile{
				PodCIDR:                "192.168.100.0/23",
				ServiceCIDR:            "192.168.200.0/23",
				SoftwareDefinedNetwork: api.SoftwareDefinedNetworkOVNKubernetes,
			},
			MasterProfile: api.MasterProfile{
				VMSize:   api.VMSizeStandardD8asV4,
				SubnetID: masterSubnetResource.String(),
			},
			WorkerProfiles: []api.WorkerProfile{
				{
					Name:     "worker",
					VMSize:   api.VMSizeStandardD8asV4,
					SubnetID: workerSubnetResource.String(),
				},
			},
			APIServerProfile: api.APIServerProfile{
				IntIP: "192.168.0.100",
			},
			IngressProfiles: []api.IngressProfile{
				{
					IP: "192.168.0.1",
				},
			},
			SSHKey: sshKeyDer.Bytes,
		},
		/*
			Identity: &api.Identity{
				Type:                   "",
				UserAssignedIdentities: "",
				IdentityURL:            "",
			},
		*/
	}

	log.Info("Setting up Subscription")
	var sub api.Subscription = api.Subscription{
		State: api.SubscriptionStateRegistered,
		Properties: &api.SubscriptionProperties{
			TenantID: tenantID,
			AccountOwner: &api.AccountOwnerProfile{
				Email: email,
			},
		},
	}

	log.Info("Setting InstallDir")
	cluster.InstallDir = os.Getenv("ARO_BASE_PATH")

	log.Info("Making Installer")
	i, err := MakeInstaller(ctx, log, assetsDirectory, &oc, &sub)
	if err != nil {
		t.Fatalf("failed to make installer: %v", err)
	}

	log.Info("Generating InstallConfig")
	ic, image, err := i.GenerateInstallConfig(ctx)
	if err != nil {
		t.Fatalf("failed to generate install config: %v", err)
	}

	log.Info("Applying customizations to InstallConfig")
	g, err := i.ApplyInstallConfigCustomisations(ic, image)
	if err != nil {
		t.Fatalf("failed to apply install config customizations: %v", err)
	}

	log.Info("Rendering Manifests")
	for _, m := range targetassets.Manifests {
		log.Infof("-> %v", reflect.TypeOf(m))
		err = g.Resolve(m)
		if err != nil {
			t.Fatalf("failed to resolve asset: %v", err)
		}

		a := g.Get(m).(asset.WritableAsset)
		err = asset.PersistToFile(a, assetsDirectory)
		if err != nil {
			t.Fatalf("failed to persist asset to file: %v", err)
		}
	}

	log.Info("Rendering Ignition Configs")
	for _, m := range targetassets.IgnitionConfigs {
		log.Infof("-> %v", reflect.TypeOf(m))
		err = g.Resolve(m)
		if err != nil {
			t.Fatalf("failed to resolve asset: %v", err)
		}

		a := g.Get(m).(asset.WritableAsset)
		err = asset.PersistToFile(a, assetsDirectory)
		if err != nil {
			t.Fatalf("failed to persist asset to file: %v", err)
		}
	}

	log.Info("Done")
}
