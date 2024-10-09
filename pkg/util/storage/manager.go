package storage

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"

	"github.com/openshift/installer-aro-wrapper/pkg/api"
	"github.com/openshift/installer-aro-wrapper/pkg/env"
	storagesdk "github.com/openshift/installer-aro-wrapper/pkg/util/azureclient/azuresdk/armstorage"
	"github.com/openshift/installer-aro-wrapper/pkg/util/azureclient/azuresdk/azblob"
)

type Manager interface {
	BlobService(ctx context.Context, resourceGroup, account string, p armstorage.Permissions, r armstorage.SignedResourceTypes) (azblob.BlobsClient, error)
}

type manager struct {
	env                  env.Core
	storageAccounts      storagesdk.AccountsClient
	credential           azcore.TokenCredential
	usesWorkloadIdentity bool
}

func NewManager(env env.Core, subscriptionID string, credential azcore.TokenCredential, usesWorkloadIdentity bool) (m Manager, err error) {
	var accountsClient storagesdk.AccountsClient
	if !usesWorkloadIdentity {
		accountsClient, err = storagesdk.NewAccountsClient(env.Environment(), subscriptionID, credential)
		if err != nil {
			return nil, err
		}
	}
	return &manager{
		env:                  env,
		storageAccounts:      accountsClient,
		usesWorkloadIdentity: usesWorkloadIdentity,
		credential:           credential,
	}, nil
}

func getCorrectErrWhenTooManyRequests(err error) error {
	responseError, ok := err.(*azcore.ResponseError)
	if !ok {
		return err
	}
	if responseError.StatusCode != http.StatusTooManyRequests {
		return err
	}
	msg := "Requests are being throttled due to Azure Storage limits being exceeded. Please visit https://learn.microsoft.com/en-us/azure/openshift/troubleshoot#exceeding-azure-storage-limits for more details."
	cloudError := &api.CloudError{
		StatusCode: http.StatusTooManyRequests,
		CloudErrorBody: &api.CloudErrorBody{
			Code:    api.CloudErrorCodeThrottlingLimitExceeded,
			Message: "ThrottlingLimitExceeded",
			Details: []api.CloudErrorBody{
				{
					Message: msg,
				},
			},
		},
	}
	return cloudError
}

func (m *manager) BlobService(ctx context.Context, resourceGroup, account string, p armstorage.Permissions, r armstorage.SignedResourceTypes) (blobClient azblob.BlobsClient, err error) {
	serviceURL := fmt.Sprintf("https://%s.blob.%s", account, m.env.Environment().StorageEndpointSuffix)
	if m.usesWorkloadIdentity {
		blobClient, err = azblob.NewBlobsClientUsingEntra(ctx, m.env.Environment(), serviceURL, m.credential)
		if err != nil {
			return nil, err
		}
	} else {
		t := time.Now().UTC().Truncate(time.Second)
		res, err := m.storageAccounts.ListAccountSAS(ctx, resourceGroup, account, armstorage.AccountSasParameters{
			Services:               to.Ptr(armstorage.ServicesB),
			ResourceTypes:          to.Ptr(r),
			Permissions:            to.Ptr(p),
			Protocols:              to.Ptr(armstorage.HTTPProtocolHTTPS),
			SharedAccessStartTime:  &t,
			SharedAccessExpiryTime: to.Ptr(t.Add(24 * time.Hour)),
		}, nil)
		if err != nil {
			return nil, getCorrectErrWhenTooManyRequests(err)
		}

		_, err = url.ParseQuery(*res.AccountSasToken)
		if err != nil {
			return nil, err
		}

		sasURL := fmt.Sprintf("%s/?%s", serviceURL, *res.AccountSasToken)
		blobClient, err = azblob.NewBlobsClientUsingSAS(ctx, sasURL, m.env.Environment())
		if err != nil {
			return nil, err
		}
	}
	return blobClient, nil
}
