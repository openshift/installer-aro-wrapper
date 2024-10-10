package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"os"
	"path/filepath"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/ARO-Installer/pkg/cluster/graph"
)

const (
	aroBoundSASigningKeyDir       = "boundsasigningkey"
	installerBoundSASigningKeyDir = "tls"
)

// AROBoundSASigningKey contains a user provided key and public parts for the
// service account signing key used by kube-apiserver.
// This asset does not generate any new content and only loads these files from disk
// when provided by the user.
type AROBoundSASigningKey struct {
	tls.BoundSASigningKey
}

var _ asset.WritableAsset = (*AROBoundSASigningKey)(nil)

// Name returns a human friendly name for the asset.
func (*AROBoundSASigningKey) Name() string {
	return "ARO Service Account Signing key"
}

// Load reads the private key from the disk.
// It ensures that the key provided is a valid RSA key.
func (sk *AROBoundSASigningKey) Load(f asset.FileFetcher) (bool, error) {
	keyFile, err := f.FetchByName(filepath.Join(aroBoundSASigningKeyDir, "bound-service-account-signing-key.key"))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	rsaKey, err := tls.PemToPrivateKey(keyFile.Data)
	if err != nil {
		logrus.Debugf("Failed to load rsa.PrivateKey from file: %s", err)
		return false, errors.Wrap(err, "failed to load rsa.PrivateKey from the file")
	}
	pubData, err := tls.PublicKeyToPem(&rsaKey.PublicKey)
	if err != nil {
		return false, errors.Wrap(err, "failed to extract public key from the key")
	}
	sk.FileList = []*asset.File{
		{Filename: filepath.Join(installerBoundSASigningKeyDir, "bound-service-account-signing-key.key"), Data: keyFile.Data},
		{Filename: filepath.Join(installerBoundSASigningKeyDir, "bound-service-account-signing-key.pub"), Data: pubData},
	}
	return true, nil
}

// Append ARO boundSASigningKey files to the generated graph's bootstrap asset
func (sk *AROBoundSASigningKey) AppendFilesToBootstrap(g graph.Graph) error {
	bootstrap := g.Get(&bootstrap.Bootstrap{}).(*bootstrap.Bootstrap)
	for _, file := range sk.Files() {
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
