package installer

import (
	"os"
	"path/filepath"

	"github.com/openshift/installer/pkg/asset"
	"github.com/pkg/errors"
)

const (
	aroManifestDir = "manifests"
)

type AROManifests struct {
	FileList []*asset.File
}

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

		filename, err := filepath.Rel(f.directory, path)
		if err != nil {
			return nil, err
		}

		files = append(files, &asset.File{
			Filename: filename,
			Data:     data,
		})
	}

	return files, nil
}
