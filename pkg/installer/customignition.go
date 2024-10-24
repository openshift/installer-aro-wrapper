package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"

	"github.com/coreos/ignition/v2/config/util"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/targets"

	"github.com/openshift/installer-aro-wrapper/pkg/aro"
	"github.com/openshift/installer-aro-wrapper/pkg/bootstraplogging"
	"github.com/openshift/installer-aro-wrapper/pkg/cluster/graph"
	"github.com/openshift/installer-aro-wrapper/pkg/dnsmasq"
)

// applyInstallConfigCustomisations modifies the InstallConfig and creates
// parent assets, then regenerates the InstallConfig for use for Ignition
// generation, etc.
func (m *manager) applyIgnitionConfigCustomisations(g graph.Graph) (graph.Graph, error) {
	master := &machine.Master{}
	worker := &machine.Worker{}
	bootstrapFile := &bootstrap.Bootstrap{}
	installConfigFile := &installconfig.InstallConfig{}
	aroDNSConfig := g.Get(&aro.ARODNSConfig{}).(*aro.ARODNSConfig)

	assetList := []asset.Asset{master, worker, bootstrapFile, installConfigFile}
	for _, asset := range assetList {
		err := g.Resolve(asset)
		if err != nil {
			return nil, err
		}
	}
	master.Config.Ignition.Config.Merge[0].Source = setMachineHost("master", aroDNSConfig.APIIntIP)
	worker.Config.Ignition.Config.Merge[0].Source = setMachineHost("worker", aroDNSConfig.APIIntIP)
	var temp map[string]any
	err := json.Unmarshal(bootstrapFile.Common.File.Data, &temp)
	if err != nil {
		return nil, err
	}
	bootstrapLoggingConfig := &bootstraplogging.Config{}
	err = g.Resolve(bootstrapLoggingConfig)
	if err != nil {
		return nil, err
	}
	temp["LoggingConfig"] = bootstrapLoggingConfig
	bootstrapFile.Common.File.Data, err = json.Marshal(temp)
	if err != nil {
		return nil, err
	}
	g.Set(master, worker, bootstrapFile)
	m.log.Print("resolving graph")
	for _, a := range targets.IgnitionConfigs {
		err := g.Resolve(a)
		if err != nil {
			return nil, err
		}
	}
	newServices := []string{
		"fluentbit.service",
		"mdsd.service",
	}
	err = bootstrap.AddStorageFiles(bootstrapFile.Config, "/", "bootstrap/files", temp)
	if err != nil {
		return nil, err
	}
	err = bootstrap.AddSystemdUnits(bootstrapFile.Config, "bootstrap/systemd/common/units", temp, newServices)
	if err != nil {
		return nil, err
	}
	dnsmasqIgnConfig, err := dnsmasq.Ignition3Config(installConfigFile.Config.ClusterDomain(), aroDNSConfig.APIIntIP, aroDNSConfig.IngressIP, aroDNSConfig.GatewayDomains, aroDNSConfig.GatewayPrivateEndpointIP, true)
	if err != nil {
		return nil, err
	}

	bootstrapFile.Config.Storage.Files = append(bootstrapFile.Config.Storage.Files, dnsmasqIgnConfig.Storage.Files...)
	bootstrapFile.Config.Systemd.Units = append(bootstrapFile.Config.Systemd.Units, dnsmasqIgnConfig.Systemd.Units...)
	return g, nil
}

func setMachineHost(role, ip string) *string {
	return util.StrToPtr(func() *url.URL {
		return &url.URL{
			Scheme: "https",
			Host:   net.JoinHostPort(ip, "22623"),
			Path:   fmt.Sprintf("/config/%s", role),
		}
	}().String())
}
