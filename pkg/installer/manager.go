package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"

	"github.com/openshift/installer-aro-wrapper/pkg/api"
	"github.com/openshift/installer-aro-wrapper/pkg/bootstraplogging"
	"github.com/openshift/installer-aro-wrapper/pkg/cluster/graph"
	"github.com/openshift/installer-aro-wrapper/pkg/env"
	"github.com/openshift/installer-aro-wrapper/pkg/util/azureclient/mgmt/features"
	"github.com/openshift/installer-aro-wrapper/pkg/util/refreshable"
)

type manager struct {
	log       *logrus.Entry
	env       env.Interface
	assetsDir string

	// clusterUUID is the UUID of the OpenShiftClusterDocument that contained
	// this OpenShiftCluster. It should be used where a unique ID for this
	// cluster is required.
	clusterUUID  string
	oc           *api.OpenShiftCluster
	sub          *api.Subscription
	fpAuthorizer refreshable.Authorizer

	deployments features.DeploymentsClient

	graph graph.Manager

	kubernetescli kubernetes.Interface

	getBootstrapLoggingConfig func(env.Interface, *api.OpenShiftCluster) (*bootstraplogging.Config, error)
	getGatewayDomains         func(env.Interface, *api.OpenShiftCluster) []string
}

type Interface interface {
	Install(ctx context.Context) error
	Manifests(ctx context.Context) (graph.Graph, error)
}

func gatewayDomains(env env.Interface, oc *api.OpenShiftCluster) []string {
	return append(env.GatewayDomains(), oc.Properties.ImageRegistryStorageAccountName+".blob."+env.Environment().StorageEndpointSuffix)
}

func NewInstaller(log *logrus.Entry, _env env.Interface, assetsDir string, clusterUUID string, oc *api.OpenShiftCluster, subscription *api.Subscription, fpAuthorizer refreshable.Authorizer, deployments features.DeploymentsClient, g graph.Manager) Interface {
	return &manager{
		log:                       log,
		env:                       _env,
		assetsDir:                 assetsDir,
		clusterUUID:               clusterUUID,
		oc:                        oc,
		sub:                       subscription,
		fpAuthorizer:              fpAuthorizer,
		deployments:               deployments,
		graph:                     g,
		getBootstrapLoggingConfig: bootstraplogging.GetConfig,
		getGatewayDomains:         gatewayDomains,
	}
}
