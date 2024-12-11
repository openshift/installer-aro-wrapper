package etchost

import (
	"github.com/openshift/installer-aro-wrapper/pkg/installer/dnsmasq"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/installconfig"
)

func AppendEtcHostFiles(bootstrapAsset *bootstrap.Bootstrap, installConfig installconfig.InstallConfig, dnsConfig dnsmasq.DNSConfig) error {
	etcHostIgnConfig, err := EtcHostsIgnitionConfig(installConfig.Config.ClusterDomain(), dnsConfig.APIIntIP, dnsConfig.GatewayDomains, dnsConfig.GatewayPrivateEndpointIP)
	if err != nil {
		return err
	}
	bootstrapAsset.Config.Storage.Files = append(bootstrapAsset.Config.Storage.Files, etcHostIgnConfig.Storage.Files...)
	bootstrapAsset.Config.Systemd.Units = append(bootstrapAsset.Config.Systemd.Units, etcHostIgnConfig.Systemd.Units...)
	etcHostMasterMachineConfig, err := EtcHostsMachineConfig(installConfig.Config.ClusterDomain(), dnsConfig.APIIntIP, dnsConfig.GatewayDomains, dnsConfig.GatewayPrivateEndpointIP, "master")
	if err != nil {
		return err
	}
	etcHostWorkerMachineConfig, err := EtcHostsMachineConfig(installConfig.Config.ClusterDomain(), dnsConfig.APIIntIP, dnsConfig.GatewayDomains, dnsConfig.GatewayPrivateEndpointIP, "worker")
	if err != nil {
		return err
	}

	err = dnsmasq.AppendMachineConfigToBootstrap(etcHostMasterMachineConfig, bootstrapAsset, "/opt/openshift/openshift/99_openshift-machineconfig_99-master-aro-etc-hosts-gateway-domains.yaml")
	if err != nil {
		return err
	}
	err = dnsmasq.AppendMachineConfigToBootstrap(etcHostWorkerMachineConfig, bootstrapAsset, "/opt/openshift/openshift/99_openshift-machineconfig_99-worker-aro-etc-hosts-gateway-domains.yaml")
	if err != nil {
		return err
	}
	return nil
}
