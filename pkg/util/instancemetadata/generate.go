package instancemetadata

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

//go:generate rm -rf ../mocks/$GOPACKAGE
//go:generate mockgen -destination=../mocks/$GOPACKAGE/$GOPACKAGE.go github.com/openshift/installer-aro-wrapper/pkg/util/$GOPACKAGE InstanceMetadata
//go:generate goimports -local=github.com/openshift/installer-aro-wrapper -e -w ../mocks/$GOPACKAGE/$GOPACKAGE.go
