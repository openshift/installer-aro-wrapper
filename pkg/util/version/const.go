package version

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

const (
	DevClusterGenevaLoggingAccount       = "AROClusterLogs"
	DevClusterGenevaLoggingConfigVersion = "2.4"
	DevClusterGenevaLoggingNamespace     = "AROClusterLogs"
	DevClusterGenevaMetricsAccount       = "AzureRedHatOpenShiftCluster"
	DevGenevaLoggingEnvironment          = "Test"
	DevRPGenevaLoggingAccount            = "ARORPLogs"
	DevRPGenevaLoggingConfigVersion      = "4.3"
	DevRPGenevaLoggingNamespace          = "ARORPLogs"
	DevRPGenevaMetricsAccount            = "AzureRedHatOpenShiftRP"

	DevGatewayGenevaLoggingConfigVersion = "4.3"
)

// FluentbitImage contains the location of the Fluentbit container image
func FluentbitImage(acrDomain string) string {
	return acrDomain + "/fluentbit:1.9.10-cm20241208@sha256:fa35a491542b1e531b73658da83e47f0f549786a186f00b0cdaffec86100c980"
}

// MdmImage contains the location of the MDM container image
// https://eng.ms/docs/products/geneva/collect/references/linuxcontainers
func MdmImage(acrDomain string) string {
	return acrDomain + "/distroless/genevamdm:2.2025.404.1254-77220c-20250406t1133@sha256:8e89bfec19c81398afa0ec51a97d748cc6b7b85cf9440dd1c7ea75b24302fe55"
}

// MdsdImage contains the location of the MDSD container image
// https://eng.ms/docs/products/geneva/collect/references/linuxcontainers
func MdsdImage(acrDomain string) string {
	return acrDomain + "/distroless/genevamdsd:mariner_20250225.2@sha256:da59ef7cfe3b0b9b6b453930cc629605cf3528ed11dbb88cdc50a38633198add"
}

// MUOImage contains the location of the Managed Upgrade Operator container image
func MUOImage(acrDomain string) string {
	return acrDomain + "/managed-upgrade-operator:aro-b4"
}
