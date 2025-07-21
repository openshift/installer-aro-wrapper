package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/tls"
)

const (
	aroBoundSASigningKeyDir       = "boundsasigningkey"
	installerBoundSASigningKeyDir = "tls"
)

// AROBoundSASigningKey is a custom wrapper of tls.BoundSASigningKey, to read the
// filepath expected in the ARO Installer wrapper's context
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
