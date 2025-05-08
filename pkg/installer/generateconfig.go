package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"os"
	"slices"
	"strings"

	mgmtcompute "github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/to"
	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset/installconfig"
	icazure "github.com/openshift/installer/pkg/asset/installconfig/azure"
	"github.com/openshift/installer/pkg/asset/releaseimage"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/validation"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer-aro-wrapper/pkg/api"
	"github.com/openshift/installer-aro-wrapper/pkg/util/computeskus"
	utilpem "github.com/openshift/installer-aro-wrapper/pkg/util/pem"
	"github.com/openshift/installer-aro-wrapper/pkg/util/pullsecret"
	"github.com/openshift/installer-aro-wrapper/pkg/util/stringutils"
	"github.com/openshift/installer-aro-wrapper/pkg/util/subnet"
)

const CONTROL_PLANE_MACHINE_COUNT = 3

func (m *manager) generateInstallConfig(ctx context.Context) (*installconfig.InstallConfig, *releaseimage.Image, error) {
	resourceGroup := stringutils.LastTokenByte(m.oc.Properties.ClusterProfile.ResourceGroupID, '/')

	pullSecret, err := pullsecret.Build(m.oc, string(m.oc.Properties.ClusterProfile.PullSecret))
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	for _, key := range []string{"cloud.openshift.com"} {
		pullSecret, err = pullsecret.RemoveKey(pullSecret, key)
		if err != nil {
			return nil, nil, errors.WithStack(err)
		}
	}

	r, err := azure.ParseResourceID(m.oc.ID)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	_, masterSubnetName, err := subnet.Split(m.oc.Properties.MasterProfile.SubnetID)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	vnetID, workerSubnetName, err := subnet.Split(m.oc.Properties.WorkerProfiles[0].SubnetID)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	vnetr, err := azure.ParseResourceID(vnetID)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(m.oc.Properties.SSHKey)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	sshkey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	domain := m.oc.Properties.ClusterProfile.Domain
	if !strings.ContainsRune(domain, '.') {
		domain += "." + m.env.Domain()
	}

	masterSKU, err := m.env.VMSku(string(m.oc.Properties.MasterProfile.VMSize))
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	masterVMNetworkingType := determineVMNetworkingType(masterSKU)

	workerSKU, err := m.env.VMSku(string(m.oc.Properties.WorkerProfiles[0].VMSize))
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	workerVMNetworkingType := determineVMNetworkingType(workerSKU)

	controlPlaneZones, workerZones, err := determineAvailabilityZones(masterSKU, workerSKU)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	// centraluseuap reports one zone, so we need to perform a non-zonal install in that region
	if strings.EqualFold(m.oc.Location, "centraluseuap") {
		workerZones = []string{""}
		controlPlaneZones = []string{""}
	}

	// Set NetworkType to OVNKubernetes by default
	softwareDefinedNetwork := string(api.SoftwareDefinedNetworkOVNKubernetes)
	if string(m.oc.Properties.NetworkProfile.SoftwareDefinedNetwork) != "" {
		softwareDefinedNetwork = string(m.oc.Properties.NetworkProfile.SoftwareDefinedNetwork)
	}

	// determine outbound type based on cluster visibility
	outboundType := azuretypes.LoadbalancerOutboundType
	if m.oc.Properties.NetworkProfile.OutboundType == api.OutboundTypeUserDefinedRouting {
		outboundType = azuretypes.UserDefinedRoutingOutboundType
	}

	var masterDiskEncryptionSet *azuretypes.DiskEncryptionSet
	var workerDiskEncryptionSet *azuretypes.DiskEncryptionSet

	if m.oc.Properties.MasterProfile.DiskEncryptionSetID != "" {
		masterDiskEncryptionSetResource, err := azure.ParseResourceID(m.oc.Properties.MasterProfile.DiskEncryptionSetID)
		if err != nil {
			return nil, nil, err
		}
		masterDiskEncryptionSet = &azuretypes.DiskEncryptionSet{
			SubscriptionID: masterDiskEncryptionSetResource.SubscriptionID,
			ResourceGroup:  masterDiskEncryptionSetResource.ResourceGroup,
			Name:           masterDiskEncryptionSetResource.ResourceName,
		}
	}

	if m.oc.Properties.WorkerProfiles[0].DiskEncryptionSetID != "" {
		workerDiskEncryptionSetResource, err := azure.ParseResourceID(m.oc.Properties.WorkerProfiles[0].DiskEncryptionSetID)
		if err != nil {
			return nil, nil, err
		}
		workerDiskEncryptionSet = &azuretypes.DiskEncryptionSet{
			SubscriptionID: workerDiskEncryptionSetResource.SubscriptionID,
			ResourceGroup:  workerDiskEncryptionSetResource.ResourceGroup,
			Name:           workerDiskEncryptionSetResource.ResourceName,
		}
	}

	// TODO: Load this from the OpenShiftCluster from the RP maybe, or get it
	// from a manifest so it can be specified in the RP's
	// OpenShiftClusterVersions?
	rhcosImage := &azuretypes.OSImage{
		Publisher: "azureopenshift",
		Offer:     "aro4",
		SKU:       "aro_417",         // "aro_4x"
		Version:   "417.94.20240701", // "4x.yy.2020zzzz"
		Plan:      azuretypes.ImageNoPurchasePlan,
	}

	installConfig := &installconfig.InstallConfig{
		AssetBase: installconfig.AssetBase{
			Config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: domain[:strings.IndexByte(domain, '.')],
				},
				SSHKey:     sshkey.Type() + " " + base64.StdEncoding.EncodeToString(sshkey.Marshal()),
				BaseDomain: domain[strings.IndexByte(domain, '.')+1:],
				Networking: &types.Networking{
					MachineNetwork: []types.MachineNetworkEntry{
						{
							CIDR: *ipnet.MustParseCIDR("127.0.0.0/8"), // dummy
						},
					},
					NetworkType: softwareDefinedNetwork,
					ClusterNetwork: []types.ClusterNetworkEntry{
						{
							CIDR:       *ipnet.MustParseCIDR(m.oc.Properties.NetworkProfile.PodCIDR),
							HostPrefix: 23,
						},
					},
					ServiceNetwork: []ipnet.IPNet{
						*ipnet.MustParseCIDR(m.oc.Properties.NetworkProfile.ServiceCIDR),
					},
				},
				ControlPlane: &types.MachinePool{
					Name:     "master",
					Replicas: to.Int64Ptr(CONTROL_PLANE_MACHINE_COUNT),
					Platform: types.MachinePoolPlatform{
						Azure: &azuretypes.MachinePool{
							Zones:            controlPlaneZones,
							InstanceType:     string(m.oc.Properties.MasterProfile.VMSize),
							EncryptionAtHost: m.oc.Properties.MasterProfile.EncryptionAtHost == api.EncryptionAtHostEnabled,
							VMNetworkingType: masterVMNetworkingType,
							OSDisk: azuretypes.OSDisk{
								DiskEncryptionSet: masterDiskEncryptionSet,
								DiskSizeGB:        1024,
							},
							OSImage: *rhcosImage,
						},
					},
					Hyperthreading: "Enabled",
					Architecture:   types.ArchitectureAMD64,
				},
				Compute: []types.MachinePool{
					{
						Name:     m.oc.Properties.WorkerProfiles[0].Name,
						Replicas: to.Int64Ptr(int64(m.oc.Properties.WorkerProfiles[0].Count)),
						Platform: types.MachinePoolPlatform{
							Azure: &azuretypes.MachinePool{
								Zones:            workerZones,
								InstanceType:     string(m.oc.Properties.WorkerProfiles[0].VMSize),
								EncryptionAtHost: m.oc.Properties.WorkerProfiles[0].EncryptionAtHost == api.EncryptionAtHostEnabled,
								VMNetworkingType: workerVMNetworkingType,
								OSDisk: azuretypes.OSDisk{
									DiskEncryptionSet: workerDiskEncryptionSet,
									DiskSizeGB:        int32(m.oc.Properties.WorkerProfiles[0].DiskSizeGB),
								},
								OSImage: *rhcosImage,
							},
						},
						Hyperthreading: "Enabled",
						Architecture:   types.ArchitectureAMD64,
					},
				},
				Platform: types.Platform{
					Azure: &azuretypes.Platform{
						Region:                   strings.ToLower(m.oc.Location), // Used in k8s object names, so must pass DNS-1123 validation
						NetworkResourceGroupName: vnetr.ResourceGroup,
						VirtualNetwork:           vnetr.ResourceName,
						ControlPlaneSubnet:       masterSubnetName,
						ComputeSubnet:            workerSubnetName,
						CloudName:                azuretypes.CloudEnvironment(m.env.Environment().Name),
						OutboundType:             outboundType,
						ResourceGroupName:        resourceGroup,
						// We specify BaseDomainResourceGroupName even though we
						// do not create Public DNS zones to pass validation.
						// See https://issues.redhat.com/browse/OCPSTRAT-991 for
						// a more permanent fix (disabling public DNS
						// provisioning).
						BaseDomainResourceGroupName: resourceGroup,
					},
				},
				PullSecret: pullSecret,
				FIPS:       m.oc.Properties.ClusterProfile.FipsValidatedModules == api.FipsValidatedModulesEnabled,
				ImageDigestSources: []types.ImageDigestSource{
					{
						Source: "quay.io/openshift-release-dev/ocp-release",
						Mirrors: []string{
							fmt.Sprintf("%s/openshift-release-dev/ocp-release", m.env.ACRDomain()),
						},
					},
					{
						Source: "quay.io/openshift-release-dev/ocp-release-nightly",
						Mirrors: []string{
							fmt.Sprintf("%s/openshift-release-dev/ocp-release-nightly", m.env.ACRDomain()),
						},
					},
					{
						Source: "quay.io/openshift-release-dev/ocp-v4.0-art-dev",
						Mirrors: []string{
							fmt.Sprintf("%s/openshift-release-dev/ocp-v4.0-art-dev", m.env.ACRDomain()),
						},
					},
				},
				Publish: types.ExternalPublishingStrategy,
				Capabilities: &types.Capabilities{
					// don't include the baremetal capability (in the baseline default)
					BaselineCapabilitySet: configv1.ClusterVersionCapabilitySetNone,
					AdditionalEnabledCapabilities: []configv1.ClusterVersionCapability{
						configv1.ClusterVersionCapabilityBuild,
						configv1.ClusterVersionCapabilityCloudControllerManager,
						configv1.ClusterVersionCapabilityCloudCredential,
						configv1.ClusterVersionCapabilityConsole,
						configv1.ClusterVersionCapabilityCSISnapshot,
						configv1.ClusterVersionCapabilityDeploymentConfig,
						configv1.ClusterVersionCapabilityImageRegistry,
						configv1.ClusterVersionCapabilityIngress,
						configv1.ClusterVersionCapabilityInsights,
						configv1.ClusterVersionCapabilityMachineAPI,
						configv1.ClusterVersionCapabilityMarketplace,
						configv1.ClusterVersionCapabilityNodeTuning,
						configv1.ClusterVersionCapabilityOpenShiftSamples,
						configv1.ClusterVersionCapabilityOperatorLifecycleManager,
						configv1.ClusterVersionCapabilityStorage,
					},
				},
			}},
	}

	if m.oc.Properties.IngressProfiles[0].Visibility == api.VisibilityPrivate {
		installConfig.Config.Publish = types.InternalPublishingStrategy
	}

	var credentials *icazure.Credentials

	if m.oc.UsesWorkloadIdentity() {
		installConfig.Config.CredentialsMode = types.ManualCredentialsMode

		credentials, err = m.newInstallConfigClientCertificateCredential(m.sub.Properties.TenantID, r.SubscriptionID)
		if err != nil {
			return nil, nil, err
		}
	} else {
		credentials = &icazure.Credentials{
			TenantID:       m.sub.Properties.TenantID,
			SubscriptionID: r.SubscriptionID,
			ClientID:       m.oc.Properties.ServicePrincipalProfile.ClientID,
			ClientSecret:   string(m.oc.Properties.ServicePrincipalProfile.ClientSecret),
		}
	}

	installConfig.Azure = icazure.NewMetadataWithCredentials(
		azuretypes.CloudEnvironment(m.env.Environment().Name),
		m.env.Environment().ResourceManagerEndpoint,
		credentials,
	)

	releaseImageOverride := os.Getenv("OPENSHIFT_INSTALL_RELEASE_IMAGE_OVERRIDE")
	if releaseImageOverride == "" {
		return nil, nil, fmt.Errorf("no release image in 'OPENSHIFT_INSTALL_RELEASE_IMAGE_OVERRIDE'")
	}

	image := &releaseimage.Image{
		PullSpec: releaseImageOverride,
	}

	err = validation.ValidateInstallConfig(installConfig.Config, false).ToAggregate()
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return installConfig, image, err
}

func determineVMNetworkingType(vmSku *mgmtcompute.ResourceSku) string {
	var vmNetworkingType azuretypes.VMNetworkingCapability

	if computeskus.HasCapability(vmSku, azuretypes.AcceleratedNetworkingEnabled) {
		vmNetworkingType = azuretypes.VMnetworkingTypeAccelerated
	} else {
		vmNetworkingType = azuretypes.VMNetworkingTypeBasic
	}
	return string(vmNetworkingType)
}

func (m *manager) newInstallConfigClientCertificateCredential(tenantId, subscriptionId string) (*icazure.Credentials, error) {
	fpPrivateKey, fpCertificates := m.env.FPCertificates()

	clientCertificateFile, err := os.CreateTemp("/tmp", "fpClientCertificate-*.pem")
	if err != nil {
		return nil, err
	}

	defer clientCertificateFile.Close()

	if err = utilpem.Encode(clientCertificateFile, fpCertificates...); err != nil {
		return nil, err
	}
	if err = utilpem.Encode(clientCertificateFile, fpPrivateKey); err != nil {
		return nil, err
	}

	return &icazure.Credentials{
		TenantID:              tenantId,
		SubscriptionID:        subscriptionId,
		ClientID:              m.env.FPClientID(),
		ClientCertificatePath: clientCertificateFile.Name(),
	}, nil
}

func determineAvailabilityZones(controlPlaneSKU, workerSKU *mgmtcompute.ResourceSku) ([]string, []string, error) {
	// We handle the case where regions have no zones or >= zones than replicas,
	// but not when replicas > zones. We (currently) only support 3 control
	// plane replicas and Azure AZs will always be a minimum of 3, see
	// https://azure.microsoft.com/en-us/blog/our-commitment-to-expand-azure-availability-zones-to-more-regions/
	controlPlaneZones := computeskus.Zones(controlPlaneSKU)
	if len(controlPlaneZones) == 0 {
		controlPlaneZones = []string{""}
	} else if len(controlPlaneZones) < CONTROL_PLANE_MACHINE_COUNT {
		return nil, nil, fmt.Errorf("cluster creation with %d zones and %d control plane replicas is unsupported", len(controlPlaneZones), CONTROL_PLANE_MACHINE_COUNT)
	} else if len(controlPlaneZones) >= CONTROL_PLANE_MACHINE_COUNT {
		// Some regions may have more availability zones than control plane
		// replicas. In that case, sort them (just in case) and pick the first
		// ones in the list to satisfy the replica requirements. We may want to
		// be smarter in future, since if we have 3 replicas we will never take
		// advantage of a 4th zone if all 4 have SKU availability, but this will
		// cover the case where a SKU is available in zones 1, 2, 4 correctly.
		slices.Sort(controlPlaneZones)
		controlPlaneZones = controlPlaneZones[:CONTROL_PLANE_MACHINE_COUNT]
	}

	// Unlike above, we don't particularly mind if we pass the Installer more
	// zones than the usual 3 in a zonal region, since it automatically balances
	// them across the available zones. However, if a SKU is available in less
	// than 3 regions we will fail, since taints on cluster components like
	// Prometheus will prevent the eventual install from turning healthy.
	workerZones := computeskus.Zones(workerSKU)
	if len(workerZones) == 0 {
		workerZones = []string{""}
	} else if len(workerZones) < 3 {
		return nil, nil, fmt.Errorf("cluster creation with a worker SKU available on less than 3 zones is unsupported (available: %d)", len(workerZones))
	}

	return controlPlaneZones, workerZones, nil
}
