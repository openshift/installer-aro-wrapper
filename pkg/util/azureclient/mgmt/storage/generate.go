package storage

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

//go:generate rm -rf ../../../../util/mocks/azureclient/mgmt/$GOPACKAGE
//go:generate mockgen -destination=../../../../util/mocks/azureclient/mgmt/$GOPACKAGE/$GOPACKAGE.go github.com/openshift/installer-aro-wrapper/pkg/util/azureclient/mgmt/$GOPACKAGE AccountsClient
//go:generate goimports -local=github.com/openshift/installer-aro-wrapper -e -w ../../../../util/mocks/azureclient/mgmt/$GOPACKAGE/$GOPACKAGE.go
