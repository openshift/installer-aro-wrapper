package data

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/openshift/installer/data"
	_ "github.com/openshift/installer/data/data"
)

//go:embed all:vendor/github.com/openshift/installer/data/data
var installerData embed.FS

func init() {
	dataDir := "vendor/github.com/openshift/installer/data/data"
	dataFS, err := fs.Sub(installerData, dataDir)
	if err != nil {
		panic(err)
	}
	// Propagate our locally-generated data back into the installer library
	data.Assets = http.FS(dataFS)
}
