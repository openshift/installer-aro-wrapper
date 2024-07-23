package env

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

//go:generate rm -rf ../util/mocks/$GOPACKAGE
//go:generate go run ../../vendor/github.com/golang/mock/mockgen -destination=../util/mocks/$GOPACKAGE/$GOPACKAGE.go github.com/openshift/installer-aro-wrapper/pkg/$GOPACKAGE Core,Interface,CertificateRefresher
//go:generate go run ../../vendor/golang.org/x/tools/cmd/goimports -local=github.com/openshift/installer-aro-wrapper -e -w ../util/mocks/$GOPACKAGE/$GOPACKAGE.go
//go:generate go run ../../vendor/github.com/alvaroloes/enumer -type Feature -output zz_generated_feature_enumer.go
