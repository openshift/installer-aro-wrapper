package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"

	"github.com/openshift/ARO-Installer/pkg/api"
	"github.com/openshift/ARO-Installer/pkg/cluster/graph"
	"github.com/openshift/ARO-Installer/pkg/env"
	"github.com/openshift/ARO-Installer/pkg/util/azureclient/mgmt/features"
	"github.com/openshift/ARO-Installer/pkg/util/refreshable"
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
}

type Interface interface {
	Install(ctx context.Context) error
	Manifests(ctx context.Context) (graph.Graph, error)
}

func NewInstaller(log *logrus.Entry, _env env.Interface, assetsDir string, clusterUUID string, oc *api.OpenShiftCluster, subscription *api.Subscription, fpAuthorizer refreshable.Authorizer, deployments features.DeploymentsClient, g graph.Manager) Interface {
	return &manager{
		log:          log,
		env:          _env,
		assetsDir:    assetsDir,
		clusterUUID:  clusterUUID,
		oc:           oc,
		sub:          subscription,
		fpAuthorizer: fpAuthorizer,
		deployments:  deployments,
		graph:        g,
	}
}
