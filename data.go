package data

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/openshift/installer/data"
	_ "github.com/openshift/installer/data/data"
)

//go:embed all:vendor/github.com/openshift/installer/data
var installerData embed.FS

func init() {
	dataDir := "vendor/github.com/openshift/installer/data/data"
	if _, err := installerData.ReadDir(dataDir); err != nil {
		// openshift/installer-aro does not contain data in the filesystem, but
		// instead has the generated data compiled in
		return
	}

	dataFS, err := fs.Sub(installerData, dataDir)
	if err != nil {
		panic(err)
	}
	// Propagate our locally-generated data back into the installer library
	data.Assets = http.FS(dataFS)
}
