package armstorage

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"

	"github.com/openshift/installer-aro-wrapper/pkg/env"
)

type BlobStorageClient interface {
}

type AzBlobStorageClient struct {
	*azblob.Client
}

var _ BlobStorageClient = &AzBlobStorageClient{}

func NewAzBlobStorageClient(env env.Core, subscriptionID string, containerURL string, credential azcore.TokenCredential) (*AzBlobStorageClient, error) {
	client, err := azblob.NewClient(containerURL, credential, nil)
	return &AzBlobStorageClient{
		Client: client,
	}, err
}
