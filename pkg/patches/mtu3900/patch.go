package mtu3900

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import "github.com/coreos/ignition/v2/config/v3_2/types"

type patch struct {
}

func NewMTU3900() *patch {
	return &patch{}
}

func (p *patch) Files() ([]types.File, []types.Unit, error) {
	f := make([]types.File, 0, 3)
	// Override MTU on the bootstrap node itself, so cluster-network-operator
	// gets an appropriate default MTU for OpenshiftSDN or OVNKubernetes when
	// it first starts up on the bootstrap node.
	f = append(f, newMTUIgnitionFile())

	// Then add the following MachineConfig manifest files to the bootstrap
	// node's Ignition config:
	//
	// /opt/openshift/openshift/99_openshift-machineconfig_99-master-mtu.yaml
	// /opt/openshift/openshift/99_openshift-machineconfig_99-worker-mtu.yaml
	ignitionFile, err := newMTUMachineConfigIgnitionFile("master")
	if err != nil {
		return nil, nil, err
	}
	f = append(f, ignitionFile)

	ignitionFile, err = newMTUMachineConfigIgnitionFile("worker")
	if err != nil {
		return nil, nil, err
	}
	f = append(f, ignitionFile)

	return f, nil, nil
}
