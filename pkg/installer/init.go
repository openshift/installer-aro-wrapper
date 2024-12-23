package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import "os"

func init() {
	os.Setenv("OPENSHIFT_INSTALL_INVOKER", "ARO")
}
