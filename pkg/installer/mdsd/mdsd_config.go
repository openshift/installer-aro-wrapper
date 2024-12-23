package mdsd

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"github.com/openshift/installer/pkg/asset/bootstraplogging"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"

	bootstrapfiles "github.com/openshift/installer-aro-wrapper/pkg/data/bootstrap"
	mdsdfiles "github.com/openshift/installer-aro-wrapper/pkg/data/mdsd"
)

func AppendMdsdFiles(bootstrap *bootstrap.Bootstrap, bootstrapLoggingConfig *bootstraplogging.Config) error {
	err := bootstrapfiles.AddStorageFiles(bootstrap.Config, "etc", "etc", bootstrapLoggingConfig, mdsdfiles.Assets)
	if err != nil {
		return err
	}

	err = bootstrapfiles.AddSystemdUnits(bootstrap.Config, "systemd/units", bootstrapLoggingConfig, []string{"fluentbit.service", "mdsd.service"}, mdsdfiles.Assets)
	return err
}
