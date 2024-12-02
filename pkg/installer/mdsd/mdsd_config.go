package mdsd

import (
	"github.com/openshift/installer-aro-wrapper/pkg/bootstraplogging"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
)

func AppendMdsdFiles(bootstrap *bootstrap.Bootstrap, bootstrapLoggingConfig *bootstraplogging.Config) error {
	err := AddStorageFiles(bootstrap.Config, "etc", "etc", bootstrapLoggingConfig)
	if err != nil {
		return err
	}

	err = AddSystemdUnits(bootstrap.Config, "systemd/units", bootstrapLoggingConfig, []string{"fluentbit.service", "mdsd.service"})
	return err
}
