package bootstrap

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"bytes"
	"embed"
	"html/template"
	"io"
	"path"
	"path/filepath"
	"strings"

	ignutil "github.com/coreos/ignition/v2/config/util"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	mcfgv1 "github.com/openshift/api/machineconfiguration/v1"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"sigs.k8s.io/yaml"
)

func AddStorageFiles(config *igntypes.Config, base string, uri string, templateData interface{}, assets embed.FS) (err error) {
	file, err := assets.Open(uri)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	if info.IsDir() {
		children, err := assets.ReadDir(uri)
		if err != nil {
			return err
		}
		if err = file.Close(); err != nil {
			return err
		}

		for _, childInfo := range children {
			name := childInfo.Name()
			err = AddStorageFiles(config, path.Join(base, name), path.Join(uri, name), templateData, assets)
			if err != nil {
				return err
			}
		}
		return nil
	}

	name := info.Name()
	_, data, err := readFile(name, file, templateData)
	if err != nil {
		return err
	}

	filename := path.Base(uri)
	parentDir := path.Base(path.Dir(uri))

	var mode int
	appendToFile := false
	if parentDir == "bin" || parentDir == "dispatcher.d" {
		mode = 0555
	} else if filename == "motd" || filename == "containers.conf" {
		mode = 0644
		appendToFile = true
	} else if filename == "registries.conf" {
		// Having the mode be private breaks rpm-ostree, xref
		// https://github.com/openshift/installer/pull/6789
		mode = 0644
	} else {
		mode = 0600
	}
	ign := ignition.FileFromBytes(strings.TrimSuffix("/"+base, ".template"), "root", mode, data)
	if appendToFile {
		ignition.ConvertToAppendix(&ign)
	}

	// Replace files that already exist in the slice with ones added later, otherwise append them
	config.Storage.Files = replaceOrAppend(config.Storage.Files, ign)
	return nil
}

func AppendMachineConfigToBootstrap(machineConfig *mcfgv1.MachineConfig, bootstrapAsset *bootstrap.Bootstrap, path string) error {
	data, err := yaml.Marshal(machineConfig)
	if err != nil {
		return err
	}
	config := ignition.FileFromBytes(path, "root", 0644, data)
	bootstrapAsset.Config.Storage.Files = ReplaceOrAppend(bootstrapAsset.Config.Storage.Files, []igntypes.File{config})
	return nil
}

func ReplaceOrAppend(bootstrapFiles []igntypes.File, file []igntypes.File) []igntypes.File {
	for _, iff := range file {
		flag := false
		for i, f := range bootstrapFiles {
			if f.Node.Path == iff.Node.Path {
				bootstrapFiles[i] = iff
				flag = true
				break
			}
		}
		if !flag {
			bootstrapFiles = append(bootstrapFiles, iff)
		}
	}
	return bootstrapFiles
}

func ReplaceOrAppendSystemd(bootstrapFiles []igntypes.Unit, file []igntypes.Unit) []igntypes.Unit {
	for _, iff := range file {
		flag := false
		for i, f := range bootstrapFiles {
			if f.Name == iff.Name {
				bootstrapFiles[i] = iff
				flag = true
				break
			}
		}
		if !flag {
			bootstrapFiles = append(bootstrapFiles, iff)
		}
	}
	return bootstrapFiles
}

func AddSystemdUnits(config *igntypes.Config, uri string, templateData interface{}, enabledServices []string, assets embed.FS) (err error) {
	enabled := make(map[string]struct{}, len(enabledServices))
	for _, s := range enabledServices {
		enabled[s] = struct{}{}
	}

	directory, err := assets.Open(uri)
	if err != nil {
		return err
	}
	defer directory.Close()

	children, err := assets.ReadDir(uri)
	if err != nil {
		return err
	}

	for _, childInfo := range children {
		dir := path.Join(uri, childInfo.Name())
		file, err := assets.Open(dir)
		if err != nil {
			return err
		}
		defer file.Close()

		info, err := file.Stat()
		if err != nil {
			return err
		}

		if info.IsDir() {
			if dir := info.Name(); !strings.HasSuffix(dir, ".d") {
				continue
			}

			children, err := assets.ReadDir(uri)
			if err != nil {
				return err
			}
			if err = file.Close(); err != nil {
				return err
			}

			dropins := []igntypes.Dropin{}
			for _, childInfo := range children {
				file, err := assets.Open(path.Join(dir, childInfo.Name()))
				if err != nil {
					return err
				}
				defer file.Close()

				childName, contents, err := readFile(childInfo.Name(), file, templateData)
				if err != nil {
					return err
				}

				dropins = append(dropins, igntypes.Dropin{
					Name:     childName,
					Contents: ignutil.StrToPtr(string(contents)),
				})
			}

			name := strings.TrimSuffix(childInfo.Name(), ".d")
			unit := igntypes.Unit{
				Name:    name,
				Dropins: dropins,
			}
			if _, ok := enabled[name]; ok {
				unit.Enabled = ignutil.BoolToPtr(true)
			}
			config.Systemd.Units = append(config.Systemd.Units, unit)
		} else {
			name, contents, err := readFile(childInfo.Name(), file, templateData)
			if err != nil {
				return err
			}

			unit := igntypes.Unit{
				Name:     name,
				Contents: ignutil.StrToPtr(string(contents)),
			}
			if _, ok := enabled[name]; ok {
				unit.Enabled = ignutil.BoolToPtr(true)
			}
			config.Systemd.Units = append(config.Systemd.Units, unit)
		}
	}

	return nil
}

func replaceOrAppend(files []igntypes.File, file igntypes.File) []igntypes.File {
	for i, f := range files {
		if f.Node.Path == file.Node.Path {
			files[i] = file
			return files
		}
	}
	files = append(files, file)
	return files
}

func readFile(name string, reader io.Reader, templateData interface{}) (finalName string, data []byte, err error) {
	data, err = io.ReadAll(reader)
	if err != nil {
		return name, []byte{}, err
	}

	if filepath.Ext(name) == ".template" {
		name = strings.TrimSuffix(name, ".template")
		tmpl := template.New(name).Funcs(template.FuncMap{"replace": replace})
		tmpl, err := tmpl.Parse(string(data))
		if err != nil {
			return name, data, err
		}
		stringData := applyTemplateData(tmpl, templateData)
		data = []byte(stringData)
	}

	return name, data, nil
}

func applyTemplateData(template *template.Template, templateData interface{}) string {
	buf := &bytes.Buffer{}
	if err := template.Execute(buf, templateData); err != nil {
		panic(err)
	}
	return buf.String()
}

func replace(input, from, to string) string {
	return strings.ReplaceAll(input, from, to)
}
