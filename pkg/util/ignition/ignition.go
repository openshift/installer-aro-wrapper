package ignition

import (
	"bytes"
	"embed"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"

	ignutil "github.com/coreos/ignition/v2/config/util"
	"github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/openshift/installer/pkg/asset/ignition"
)

// GetFiles pulls files from a embed.FS and templates them. SystemD units
// included are loaded as Units, and enabled if enabledUnits["name"] is true.
// AdditionalFileModes can be used to mark files with special file modes, since
// go:embed does not store that information.
func GetFiles(staticFiles embed.FS, templateData interface{}, enabledUnits map[string]bool, additionalFileModes map[string]int) ([]types.File, []types.Unit, error) {
	files := make([]types.File, 0)
	units := make([]types.Unit, 0)

	err := fs.WalkDir(staticFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if d == nil || !d.Type().IsRegular() {
			return nil
		}

		dirPath := strings.Split(path, string(filepath.Separator))[1:]
		cleanPath := filepath.Join(dirPath...)

		var data []byte

		if filepath.Ext(path) == ".template" {
			tmpl, err := template.ParseFS(staticFiles, path)
			if err != nil {
				return err
			}
			buf := &bytes.Buffer{}
			if err := tmpl.Execute(buf, templateData); err != nil {
				return err
			}
			data = buf.Bytes()
		} else {
			file, err := staticFiles.Open(path)
			if err != nil {
				return err
			}
			data, err = io.ReadAll(file)
			if err != nil {
				return err
			}
		}

		if dirPath[0] == "systemd" {
			finalFilename := strings.TrimSuffix(filepath.Base(path), ".template")

			unit := types.Unit{
				Name:     finalFilename,
				Contents: ignutil.StrToPtr(string(data)),
			}

			if got, ok := enabledUnits[finalFilename]; ok {
				unit.Enabled = ignutil.BoolToPtr(got)
			}

			units = append(units, unit)
		} else {
			finalFilename := "/" + strings.TrimSuffix(cleanPath, ".template")
			f := ignition.FileFromBytes(finalFilename, "root", 0555, data)

			// if we have a special file mode for this file, set it
			if got, ok := additionalFileModes[finalFilename]; ok {
				f.Mode = ignutil.IntToPtr(got)
			}

			files = append(files, f)
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return files, units, err
}

func MergeUnits(units []types.Unit, newUnits []types.Unit) []types.Unit {
	for _, new := range newUnits {
		found := false
		for i, old := range units {
			if old.Name == new.Name {
				found = true
				units[i] = new
				break
			}
		}
		if !found {
			units = append(units, new)
		}
	}
	return units
}
