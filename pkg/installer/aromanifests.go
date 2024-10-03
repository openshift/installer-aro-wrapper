package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"os"
	"path/filepath"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/pkg/errors"

	"github.com/openshift/ARO-Installer/pkg/cluster/graph"
)

const (
	aroManifestDir = "manifests"
	rootPath       = "/opt/openshift"
)

// Custom ARO asset to add custom manifests to install graph in installer-wrapper similar to installer's manifests.Manifests
type AROManifests struct {
	FileList []*asset.File
}

// ARO File Fetcher to read manifests
type aroFileFetcher struct {
	directory string
}

var (
	_ asset.WritableAsset = (*AROManifests)(nil)
	_ asset.FileFetcher   = (*aroFileFetcher)(nil)
)

func (am *AROManifests) Name() string {
	return "ARO Manifests"
}

func (am *AROManifests) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

func (am *AROManifests) Generate(dependencies asset.Parents) error {
	return nil
}

func (am *AROManifests) Files() []*asset.File {
	return am.FileList
}

func (am *AROManifests) Load(f asset.FileFetcher) (found bool, err error) {
	yamlFileList, err := f.FetchByPattern(filepath.Join(aroManifestDir, "*.yaml"))
	if err != nil {
		return false, errors.Wrap(err, "failed to load *.yaml files")
	}
	ymlFileList, err := f.FetchByPattern(filepath.Join(aroManifestDir, "*.yml"))
	if err != nil {
		return false, errors.Wrap(err, "failed to load *.yml files")
	}

	am.FileList = append(am.FileList, yamlFileList...)
	am.FileList = append(am.FileList, ymlFileList...)
	asset.SortFiles(am.FileList)

	return len(am.FileList) > 0, nil
}

// Append ARO manifest files to the generated graph's bootstrap asset
func (am *AROManifests) AppendFilesToBootstrap(g graph.Graph) error {
	bootstrap := g.Get(&bootstrap.Bootstrap{}).(*bootstrap.Bootstrap)
	for _, file := range am.Files() {
		manifest := ignition.FileFromBytes(filepath.Join(rootPath, file.Filename), "root", 0644, file.Data)
		bootstrap.Config.Storage.Files = append(bootstrap.Config.Storage.Files, manifest)
	}

	data, err := ignition.Marshal(bootstrap.Config)
	if err != nil {
		return err
	}
	bootstrap.File.Data = data
	return nil
}

func (f *aroFileFetcher) FetchByName(name string) (*asset.File, error) {
	data, err := os.ReadFile(filepath.Join(f.directory, name))
	if err != nil {
		return nil, err
	}
	return &asset.File{Filename: name, Data: data}, nil
}

func (f *aroFileFetcher) FetchByPattern(pattern string) (files []*asset.File, err error) {
	matches, err := filepath.Glob(filepath.Join(f.directory, pattern))
	if err != nil {
		return nil, err
	}

	files = make([]*asset.File, 0, len(matches))
	for _, path := range matches {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		files = append(files, &asset.File{
			Filename: path,
			Data:     data,
		})
	}

	return files, nil
}
