package keyvault

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"fmt"
	"os"

	"github.com/openshift/ARO-Installer/pkg/util/instancemetadata"
)

func URI(instancemetadata instancemetadata.InstanceMetadata, suffix string) (string, error) {
	for _, key := range []string{
		"ARO_KEYVAULT_PREFIX",
	} {
		if _, found := os.LookupEnv(key); !found {
			return "", fmt.Errorf("environment variable %q unset", key)
		}
	}

	return fmt.Sprintf("https://%s%s.%s/", os.Getenv("ARO_KEYVAULT_PREFIX"), suffix, instancemetadata.Environment().KeyVaultDNSSuffix), nil
}
