package version

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"github.com/openshift/installer-aro-wrapper/pkg/api"
)

const InstallArchitectureVersion = api.ArchitectureVersionV2

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

var GitCommit = "unknown"

// InstallStream describes stream we are defaulting to for all new clusters
var InstallStream = &Stream{
	Version:  NewVersion(4, 11, 16),
	PullSpec: "quay.io/openshift-release-dev/ocp-release@sha256:97410a5db655a9d3017b735c2c0747c849d09ff551765e49d5272b80c024a844",
}

// UpgradeStreams describes list of streams we support for upgrades
var (
	UpgradeStreams = []*Stream{
		InstallStream,
		{
			Version:  NewVersion(4, 10, 20),
			PullSpec: "quay.io/openshift-release-dev/ocp-release@sha256:b89ada9261a1b257012469e90d7d4839d0d2f99654f5ce76394fa3f06522b600",
		},
		{
			Version:  NewVersion(4, 9, 28),
			PullSpec: "quay.io/openshift-release-dev/ocp-release@sha256:4084d94969b186e20189649b5affba7da59f7d1943e4e5bc7ef78b981eafb7a8",
		},
		{
			Version:  NewVersion(4, 8, 18),
			PullSpec: "quay.io/openshift-release-dev/ocp-release@sha256:321aae3d3748c589bc2011062cee9fd14e106f258807dc2d84ced3f7461160ea",
		},
		{
			Version:  NewVersion(4, 7, 30),
			PullSpec: "quay.io/openshift-release-dev/ocp-release@sha256:aba54b293dc151f5c0fd96d4353ced6ced3e7da6620c1c10714ab32d0577486f",
		},
		{
			Version:  NewVersion(4, 6, 44),
			PullSpec: "quay.io/openshift-release-dev/ocp-release@sha256:d042aa235b538721a39989b13d7d9d3537af9b57e9fd10f485dd04461932ec85",
		},
		{
			Version:  NewVersion(4, 5, 39),
			PullSpec: "quay.io/openshift-release-dev/ocp-release@sha256:c4b9eb565c64df97afe7841bbcc0469daec7973e46ae588739cc30ea9062172b",
		},
		{
			Version:  NewVersion(4, 4, 33),
			PullSpec: "quay.io/openshift-release-dev/ocp-release@sha256:a035dddd8a5e5c99484138951ef4aba021799b77eb9046f683a5466c23717738",
		},
	}
)

// FluentbitImage contains the location of the Fluentbit container image
func FluentbitImage(acrDomain string) string {
	return acrDomain + "/fluentbit:1.9.10-cm20241208@sha256:fa35a491542b1e531b73658da83e47f0f549786a186f00b0cdaffec86100c980"
}

// MdmImage contains the location of the MDM container image
// https://eng.ms/docs/products/geneva/collect/references/linuxcontainers
func MdmImage(acrDomain string) string {
	return acrDomain + "/distroless/genevamdm:2.2024.1115.1908-5b4aed-20241230t1713@sha256:249a57801d76244f722a739c8bb03cb519cbfbc3ca8356b7da36ffe9084afecd"
}

// MdsdImage contains the location of the MDSD container image
// https://eng.ms/docs/products/geneva/collect/references/linuxcontainers
func MdsdImage(acrDomain string) string {
	return acrDomain + "/distroless/genevamdsd:mariner_20241212.2@sha256:a7a71af5b631ea5a8ad587a09d8680b17719cae25b95de81e8a4d71f2cc55f45"
}

// MUOImage contains the location of the Managed Upgrade Operator container image
func MUOImage(acrDomain string) string {
	return acrDomain + "/managed-upgrade-operator:aro-b4"
}
