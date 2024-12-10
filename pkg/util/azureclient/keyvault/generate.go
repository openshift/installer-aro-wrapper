package keyvault

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

//go:generate rm -rf ../../../util/mocks/azureclient/$GOPACKAGE
//go:generate mockgen -destination=../../../util/mocks/azureclient/$GOPACKAGE/$GOPACKAGE.go github.com/openshift/installer-aro-wrapper/pkg/util/azureclient/$GOPACKAGE BaseClient
//go:generate goimports -local=github.com/openshift/installer-aro-wrapper -e -w ../../../util/mocks/azureclient/$GOPACKAGE/$GOPACKAGE.go
