package dnsmasq

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"

	"github.com/Azure/go-autorest/autorest/to"

	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/installconfig"

	bootstrapfiles "github.com/openshift/installer-aro-wrapper/pkg/data/bootstrap"
)

type DNSConfig struct {
	APIIntIP                 string
	IngressIP                string
	GatewayDomains           []string
	GatewayPrivateEndpointIP string
}

func CreatednsmasqIgnitionFiles(bootstrapAsset *bootstrap.Bootstrap, installConfig *installconfig.InstallConfig, dnsConfig DNSConfig) error {
	dnsmasqIgnConfig, err := Ignition3Config(installConfig.Config.ClusterDomain(), dnsConfig.APIIntIP, dnsConfig.IngressIP, dnsConfig.GatewayDomains, dnsConfig.GatewayPrivateEndpointIP)
	if err != nil {
		return err
	}
	bootstrapAsset.Config.Storage.Files = bootstrapfiles.ReplaceOrAppend(bootstrapAsset.Config.Storage.Files, dnsmasqIgnConfig.Storage.Files)
	bootstrapAsset.Config.Systemd.Units = bootstrapfiles.ReplaceOrAppendSystemd(bootstrapAsset.Config.Systemd.Units, dnsmasqIgnConfig.Systemd.Units)

	// Add a drop-in for node-image-pull.service to ensure it starts after dnsmasq
	// and aro-etchosts-resolver have configured DNS and /etc/hosts. Required for
	// 4.19+ where node-image-pull.service (from the base OS image) pulls the
	// OCP release image from ACR during bootstrap.
	nodeImagePullDropin := "[Unit]\nAfter=dnsmasq.service\nAfter=aro-etchosts-resolver.service\n"
	bootstrapAsset.Config.Systemd.Units = bootstrapfiles.ReplaceOrAppendSystemd(
		bootstrapAsset.Config.Systemd.Units,
		[]igntypes.Unit{
			{
				Name: "node-image-pull.service",
				Dropins: []igntypes.Dropin{
					{
						Name:     "10-aro-dns-ordering.conf",
						Contents: to.StringPtr(nodeImagePullDropin),
					},
				},
			},
		},
	)

	dnsmasqMasterMachineConfig, err := MachineConfig(installConfig.Config.ClusterDomain(), dnsConfig.APIIntIP, dnsConfig.IngressIP, "master", dnsConfig.GatewayDomains, dnsConfig.GatewayPrivateEndpointIP)
	if err != nil {
		return err
	}
	dnsmasqWorkerMachineConfig, err := MachineConfig(installConfig.Config.ClusterDomain(), dnsConfig.APIIntIP, dnsConfig.IngressIP, "worker", dnsConfig.GatewayDomains, dnsConfig.GatewayPrivateEndpointIP)
	if err != nil {
		return err
	}
	bootstrapfiles.AppendMachineConfigToBootstrap(dnsmasqMasterMachineConfig, bootstrapAsset, "/opt/openshift/openshift/99_openshift-machineconfig_99-master-aro-dns.yaml")
	bootstrapfiles.AppendMachineConfigToBootstrap(dnsmasqWorkerMachineConfig, bootstrapAsset, "/opt/openshift/openshift/99_openshift-machineconfig_99-worker-aro-dns.yaml")
	return nil
}
