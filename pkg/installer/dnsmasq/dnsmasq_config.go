package dnsmasq

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
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
	dnsmasqIgnConfig, err := Ignition3Config(installConfig.Config.ClusterDomain(), dnsConfig.APIIntIP, dnsConfig.IngressIP, dnsConfig.GatewayDomains, dnsConfig.GatewayPrivateEndpointIP, true)
	if err != nil {
		return err
	}
	bootstrapAsset.Config.Storage.Files = bootstrapfiles.ReplaceOrAppend(bootstrapAsset.Config.Storage.Files, dnsmasqIgnConfig.Storage.Files)
	bootstrapAsset.Config.Systemd.Units = bootstrapfiles.ReplaceOrAppendSystemd(bootstrapAsset.Config.Systemd.Units, dnsmasqIgnConfig.Systemd.Units)

	dnsmasqMasterMachineConfig, err := MachineConfig(installConfig.Config.ClusterDomain(), dnsConfig.APIIntIP, dnsConfig.IngressIP, "master", dnsConfig.GatewayDomains, dnsConfig.GatewayPrivateEndpointIP, true)
	if err != nil {
		return err
	}
	dnsmasqWorkerMachineConfig, err := MachineConfig(installConfig.Config.ClusterDomain(), dnsConfig.APIIntIP, dnsConfig.IngressIP, "worker", dnsConfig.GatewayDomains, dnsConfig.GatewayPrivateEndpointIP, true)
	if err != nil {
		return err
	}
	bootstrapfiles.AppendMachineConfigToBootstrap(dnsmasqMasterMachineConfig, bootstrapAsset, "/opt/openshift/openshift/99_openshift-machineconfig_99-master-aro-dns.yaml")
	bootstrapfiles.AppendMachineConfigToBootstrap(dnsmasqWorkerMachineConfig, bootstrapAsset, "/opt/openshift/openshift/99_openshift-machineconfig_99-worker-aro-dns.yaml")
	return nil
}
