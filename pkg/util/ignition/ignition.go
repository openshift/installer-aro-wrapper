package ignition

import (
	"bytes"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"

	ignutil "github.com/coreos/ignition/v2/config/util"
	"github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/openshift/installer/pkg/asset/ignition"
)

func GetFiles(staticFiles fs.FS, templateData interface{}, enabledUnits map[string]bool) ([]types.File, []types.Unit, error) {
	files := make([]types.File, 0)
	units := make([]types.Unit, 0)

	err := fs.WalkDir(staticFiles, "/", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		dirPath := filepath.SplitList(filepath.Dir(path))

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
			finalFilename := strings.TrimSuffix(path, ".template")

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
