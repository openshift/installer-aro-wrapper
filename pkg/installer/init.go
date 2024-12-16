package installer

import "os"

func init() {
	os.Setenv("OPENSHIFT_INSTALL_INVOKER", "ARO")
}
