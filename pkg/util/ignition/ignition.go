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

func GetFiles(staticFiles embed.FS, templateData interface{}, enabledUnits map[string]bool) ([]types.File, []types.Unit, error) {
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

			if _, ok := enabledUnits[finalFilename]; ok {
				unit.Enabled = ignutil.BoolToPtr(true)
			}

			units = append(units, unit)
		} else {
			finalFilename := "/" + strings.TrimSuffix(cleanPath, ".template")

			var mode int
			if dirPath[len(dirPath)-1] == "bin" {
				mode = 0555
			} else {
				mode = 0600
			}
			files = append(files, ignition.FileFromBytes(finalFilename, "root", mode, data))
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return files, units, nil
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
