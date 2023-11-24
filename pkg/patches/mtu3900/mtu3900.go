package mtu3900

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	_ "embed"
	"fmt"

	"github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/machines/machineconfig"
	mcv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const ignFilePath = "/etc/NetworkManager/dispatcher.d/30-eth0-mtu-3900"

//go:embed files/30-eth0-mtu-3900
var ignFileData []byte

func newMTUIgnitionFile() types.File {
	return ignition.FileFromBytes(ignFilePath, "root", 0555, ignFileData)
}

func newMTUMachineConfigIgnitionFile(role string) (types.File, error) {
	mtuIgnitionConfig := types.Config{
		Ignition: types.Ignition{
			Version: types.MaxVersion.String(),
		},
		Storage: types.Storage{
			Files: []types.File{
				newMTUIgnitionFile(),
			},
		},
	}

	rawExt, err := ignition.ConvertToRawExtension(mtuIgnitionConfig)
	if err != nil {
		return types.File{}, err
	}

	mtuMachineConfig := &mcv1.MachineConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: mcv1.SchemeGroupVersion.String(),
			Kind:       "MachineConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("99-%s-mtu", role),
			Labels: map[string]string{
				"machineconfiguration.openshift.io/role": role,
			},
		},
		Spec: mcv1.MachineConfigSpec{
			Config: rawExt,
		},
	}

	configs := []*mcv1.MachineConfig{mtuMachineConfig}
	manifests, err := machineconfig.Manifests(configs, role, "/opt/openshift/openshift")
	if err != nil {
		return types.File{}, err
	}

	return ignition.FileFromBytes(manifests[0].Filename, "root", 0644, manifests[0].Data), nil
}
