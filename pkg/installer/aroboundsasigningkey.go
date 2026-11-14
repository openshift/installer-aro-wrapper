package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"crypto"
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

	privateKey, err := pemToPrivateKey(keyFile.Data)
	if err != nil {
		logrus.Debugf("Failed to load private key from file: %s", err)
		return false, errors.Wrap(err, "failed to load private key from the file")
	}
	signer, ok := privateKey.(crypto.Signer)
	if !ok {
		return false, errors.New("private key does not implement crypto.Signer")
	}
	pubData, err := publicKeyToPem(signer.Public())
	if err != nil {
		return false, errors.Wrap(err, "failed to extract public key from the key")
	}
	sk.FileList = []*asset.File{
		{Filename: filepath.Join(installerBoundSASigningKeyDir, "bound-service-account-signing-key.key"), Data: keyFile.Data},
		{Filename: filepath.Join(installerBoundSASigningKeyDir, "bound-service-account-signing-key.pub"), Data: pubData},
	}
	return true, nil
}
