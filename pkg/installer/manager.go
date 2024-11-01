package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"os"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"

	"github.com/openshift/installer-aro-wrapper/pkg/api"
	"github.com/openshift/installer-aro-wrapper/pkg/cluster/graph"
	"github.com/openshift/installer-aro-wrapper/pkg/env"
	"github.com/openshift/installer-aro-wrapper/pkg/util/azureclient/mgmt/features"
	"github.com/openshift/installer-aro-wrapper/pkg/util/encryption"
	"github.com/openshift/installer-aro-wrapper/pkg/util/refreshable"
	"github.com/openshift/installer-aro-wrapper/pkg/util/storage"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/releaseimage"
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
	GenerateInstallConfig(ctx context.Context) (*installconfig.InstallConfig, *releaseimage.Image, error)
	ApplyInstallConfigCustomisations(installConfig *installconfig.InstallConfig, image *releaseimage.Image) (graph.Graph, error)
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

func MakeInstaller(ctx context.Context, log *logrus.Entry, assetsDir string, oc *api.OpenShiftCluster, sub *api.Subscription) (Interface, error) {
	_env, err := env.NewEnv(ctx, log)
	if err != nil {
		return nil, err
	}

	fpAuthorizer, err := refreshable.NewAuthorizer(_env, sub.Properties.TenantID)
	if err != nil {
		return nil, err
	}

	fpCredClusterTenant, err := _env.FPNewClientCertificateCredential(sub.Properties.TenantID)
	if err != nil {
		return nil, err
	}

	r, err := azure.ParseResourceID(oc.ID)
	if err != nil {
		return nil, err
	}

	storage, err := storage.NewManager(r.SubscriptionID, _env.Environment().StorageEndpointSuffix, fpCredClusterTenant, oc.UsesWorkloadIdentity(), _env.Environment().ArmClientOptions())
	if err != nil {
		return nil, err
	}
	deployments := features.NewDeploymentsClient(_env.Environment(), r.SubscriptionID, fpAuthorizer)

	aead, err := encryption.NewMulti(ctx, _env.ServiceKeyvault(), env.EncryptionSecretV2Name, env.EncryptionSecretName)
	if err != nil {
		return nil, err
	}

	graph := graph.NewManager(log, aead, storage)

	return NewInstaller(log, _env, assetsDir, os.Getenv("ARO_UUID"), oc, sub, fpAuthorizer, deployments, graph), nil

}
