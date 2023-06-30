package bootstrapimagepull

import (
	"embed"

	"github.com/Azure/ARO-RP/pkg/util/ignition"
	"github.com/coreos/ignition/v2/config/v3_2/types"
)

//go:embed staticresources
var staticFiles embed.FS

type Config struct {
}

func GetFiles() ([]types.Unit, error) {
	_, units, err := ignition.GetFiles(staticFiles, Config{}, map[string]bool{})
	return units, err
}
