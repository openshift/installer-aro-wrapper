package graph

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"encoding/json"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/sirupsen/logrus"

	"github.com/openshift/ARO-Installer/pkg/util/encryption"
	"github.com/openshift/ARO-Installer/pkg/util/storage"
)

type Manager interface {
	Exists(ctx context.Context, resourceGroup, account string) (bool, error)
	Save(ctx context.Context, resourceGroup, account string, g Graph) error
	LoadPersisted(ctx context.Context, resourceGroup, account string) (PersistedGraph, error)
}

type manager struct {
	log *logrus.Entry

	aead    encryption.AEAD
	storage storage.Manager
}

func NewManager(log *logrus.Entry, aead encryption.AEAD, storage storage.Manager) Manager {
	return &manager{
		log: log,

		aead:    aead,
		storage: storage,
	}
}

func (m *manager) Exists(ctx context.Context, resourceGroup, account string) (bool, error) {
	m.log.Print("checking if graph exists")

	blobService, err := m.storage.BlobService(ctx, resourceGroup, account, armstorage.Permissions("r"), armstorage.SignedResourceTypesO)
	if err != nil {
		return false, err
	}

	return blobService.BlobExists(ctx, "aro", "graph")
}

// Load() should not be implemented: use LoadPersisted

func (m *manager) Save(ctx context.Context, resourceGroup, account string, g Graph) error {
	m.log.Print("save graph")

	blobService, err := m.storage.BlobService(ctx, resourceGroup, account, armstorage.Permissions("cw"), armstorage.SignedResourceTypesO)
	if err != nil {
		return err
	}

	bootstrap := g.Get(&bootstrap.Bootstrap{}).(*bootstrap.Bootstrap)

	_, err = blobService.UploadBuffer(ctx, "ignition", "bootstrap.ign", bootstrap.File.Data, nil)
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(g, "", "    ")
	if err != nil {
		return err
	}

	b, err = m.aead.Seal(b)
	if err != nil {
		return err
	}
	_, err = blobService.UploadBuffer(ctx, "aro", "graph", b, nil)
	return err
}

func (m *manager) LoadPersisted(ctx context.Context, resourceGroup, account string) (PersistedGraph, error) {
	m.log.Print("load persisted graph")

	blobService, err := m.storage.BlobService(ctx, resourceGroup, account, armstorage.Permissions("r"), armstorage.SignedResourceTypesO)
	if err != nil {
		return nil, err
	}

	rc, err := blobService.DownloadStream(ctx, "aro", "graph", nil)
	if err != nil {
		return nil, err
	}
	defer rc.Body.Close()

	b, err := io.ReadAll(rc.Body)
	if err != nil {
		return nil, err
	}

	b, err = m.aead.Open(b)
	if err != nil {
		return nil, err
	}

	var pg PersistedGraph
	err = json.Unmarshal(b, &pg)
	if err != nil {
		return nil, err
	}

	return pg, nil
}

// SavePersistedGraph could be implemented and used with care if needed, but
// currently we don't need it (and it's better that way)
