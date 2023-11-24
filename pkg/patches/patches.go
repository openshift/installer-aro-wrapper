package patches

import "github.com/coreos/ignition/v2/config/v3_2/types"

type IgnitionPatch interface {
	Files() ([]types.File, []types.Unit, error)
}
