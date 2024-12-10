package env

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

//go:generate rm -rf ../util/mocks/$GOPACKAGE
//go:generate mockgen -destination=../util/mocks/$GOPACKAGE/$GOPACKAGE.go github.com/openshift/installer-aro-wrapper/pkg/$GOPACKAGE Core,Interface,CertificateRefresher
//go:generate goimports -local=github.com/openshift/installer-aro-wrapper -e -w ../util/mocks/$GOPACKAGE/$GOPACKAGE.go
//go:generate enumer -type Feature -output zz_generated_feature_enumer.go
