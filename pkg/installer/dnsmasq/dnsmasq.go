package dnsmasq

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

	_ "embed"

	"github.com/coreos/go-semver/semver"
	"github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/vincent-petithory/dataurl"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/Azure/go-autorest/autorest/to"

	mcfgv1 "github.com/openshift/api/machineconfiguration/v1"
)

const (
	configFileName    = "dnsmasq.conf"
	unitFileName      = "dnsmasq.service"
	prescriptFileName = "aro-dnsmasq-pre.sh"
)

func config(clusterDomain, apiIntIP, ingressIP string, gatewayDomains []string, gatewayPrivateEndpointIP string) ([]byte, string, []byte, error) {
	t := template.Must(template.New(configFileName).Parse(configFile))
	config := &bytes.Buffer{}

	err := t.ExecuteTemplate(config, configFileName, &struct {
		ClusterDomain            string
		APIIntIP                 string
		IngressIP                string
		GatewayDomains           []string
		GatewayPrivateEndpointIP string
	}{
		ClusterDomain:            clusterDomain,
		APIIntIP:                 apiIntIP,
		IngressIP:                ingressIP,
		GatewayDomains:           gatewayDomains,
		GatewayPrivateEndpointIP: gatewayPrivateEndpointIP,
	})
	if err != nil {
		return nil, "", nil, err
	}
	t = template.Must(template.New(unitFileName).Parse(unitFile))
	service := &bytes.Buffer{}

	err = t.ExecuteTemplate(service, unitFileName, nil)
	if err != nil {
		return nil, "", nil, err
	}

	t = template.Must(template.New(prescriptFileName).Parse(preScriptFile))
	startpre := &bytes.Buffer{}

	err = t.ExecuteTemplate(startpre, prescriptFileName, nil)
	if err != nil {
		return nil, "", nil, err
	}

	return config.Bytes(), service.String(), startpre.Bytes(), nil
}

func Ignition3Config(clusterDomain, apiIntIP, ingressIP string, gatewayDomains []string, gatewayPrivateEndpointIP string) (*types.Config, error) {
	config, service, startpre, err := config(clusterDomain, apiIntIP, ingressIP, gatewayDomains, gatewayPrivateEndpointIP)
	if err != nil {
		return nil, err
	}

	ign := &types.Config{
		Ignition: types.Ignition{
			// This Ignition Config version should be kept up to date with the default
			// rendered Ignition Config version from the Machine Config Operator version
			// on the lowest OCP version we support (4.7).
			Version: semver.Version{
				Major: 3,
				Minor: 2,
			}.String(),
		},
		Storage: types.Storage{
			Files: []types.File{
				{
					Node: types.Node{
						Overwrite: to.BoolPtr(true),
						Path:      "/etc/" + configFileName,
						User: types.NodeUser{
							Name: to.StringPtr("root"),
						},
					},
					FileEmbedded1: types.FileEmbedded1{
						Contents: types.Resource{
							Source: to.StringPtr(dataurl.EncodeBytes(config)),
						},
						Mode: to.IntPtr(0644),
					},
				},
				{
					Node: types.Node{
						Overwrite: to.BoolPtr(true),
						Path:      "/usr/local/bin/" + prescriptFileName,
						User: types.NodeUser{
							Name: to.StringPtr("root"),
						},
					},
					FileEmbedded1: types.FileEmbedded1{
						Contents: types.Resource{
							Source: to.StringPtr(dataurl.EncodeBytes(startpre)),
						},
						Mode: to.IntPtr(0744),
					},
				},
			},
		},
		Systemd: types.Systemd{
			Units: []types.Unit{
				{
					Contents: &service,
					Enabled:  to.BoolPtr(true),
					Name:     unitFileName,
				},
			},
		},
	}

	return ign, nil
}

func MachineConfig(clusterDomain, apiIntIP, ingressIP, role string, gatewayDomains []string, gatewayPrivateEndpointIP string) (*mcfgv1.MachineConfig, error) {
	ignConfig, err := Ignition3Config(clusterDomain, apiIntIP, ingressIP, gatewayDomains, gatewayPrivateEndpointIP)
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(ignConfig)
	if err != nil {
		return nil, err
	}

	// canonicalise the machineconfig payload the same way as MCO
	var i interface{}
	err = json.Unmarshal(b, &i)
	if err != nil {
		return nil, err
	}

	rawExt := runtime.RawExtension{}
	rawExt.Raw, err = json.Marshal(i)
	if err != nil {
		return nil, err
	}

	return &mcfgv1.MachineConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: mcfgv1.SchemeGroupVersion.String(),
			Kind:       "MachineConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("99-%s-aro-dns", role),
			Labels: map[string]string{
				"machineconfiguration.openshift.io/role": role,
			},
		},
		Spec: mcfgv1.MachineConfigSpec{
			Config: rawExt,
		},
	}, nil
}
