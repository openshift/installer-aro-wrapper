package manifests

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"bytes"
	"embed"
	"fmt"

	"github.com/vincent-petithory/dataurl"

	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/types"

	bootstrapfiles "github.com/openshift/installer-aro-wrapper/pkg/data/bootstrap"
)

//go:embed opt/*
var assets embed.FS

type ManifestsConfig struct {
	AROWorkerRegistries string
	HTTPSecret          string
	AccountName         string
	ContainerName       string
	CloudName           string
	AROIngressInternal  bool
	AROIngressIP        string
}

func AppendManifestsFilesToBootstrap(bootstrap *bootstrap.Bootstrap, manifestsConfig ManifestsConfig) error {
	err := bootstrapfiles.AddStorageFiles(bootstrap.Config, "opt", "opt", manifestsConfig, assets)
	if err != nil {
		return err
	}
	return nil
}

func AroWorkerRegistries(idss []types.ImageDigestSource) string {
	b := &bytes.Buffer{}

	fmt.Fprintf(b, "unqualified-search-registries = [\"registry.access.redhat.com\", \"docker.io\"]\n")

	for _, ids := range idss {
		fmt.Fprintf(b, "\n[[registry]]\n  prefix = \"\"\n  location = \"%s\"\n  mirror-by-digest-only = true\n", ids.Source)

		for _, mirror := range ids.Mirrors {
			fmt.Fprintf(b, "\n  [[registry.mirror]]\n    location = \"%s\"\n", mirror)
		}
	}

	du := dataurl.New(b.Bytes(), "text/plain")
	du.Encoding = dataurl.EncodingASCII

	return du.String()
}
