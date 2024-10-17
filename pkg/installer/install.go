package installer

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"time"

	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/kubeconfig"
	"github.com/openshift/installer/pkg/asset/releaseimage"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/openshift/installer-aro-wrapper/pkg/cluster/graph"
	"github.com/openshift/installer-aro-wrapper/pkg/util/restconfig"
	"github.com/openshift/installer-aro-wrapper/pkg/util/steps"
	"github.com/openshift/installer-aro-wrapper/pkg/util/stringutils"
)

func (m *manager) Manifests(ctx context.Context) (graph.Graph, error) {
	var (
		installConfig *installconfig.InstallConfig
		image         *releaseimage.Image
		g             graph.Graph
	)

	s := []steps.Step{
		steps.Action(func(ctx context.Context) error {
			var err error
			installConfig, image, err = m.generateInstallConfig(ctx)
			return err
		}),

		steps.Action(func(ctx context.Context) error {
			var err error
			// Applies ARO-specific customisations to the InstallConfig
			g, err = m.applyInstallConfigCustomisations(installConfig, image)
			return err
		}),
		steps.Action(func(ctx context.Context) error {
			return m.persistGraph(ctx, g)
		}),
	}

	err := steps.Run(ctx, m.log, 10*time.Second, s)
	return g, err
}

func (m *manager) Ignition(ctx context.Context) (graph.Graph, error) {
	var (
		installConfig *installconfig.InstallConfig
		image         *releaseimage.Image
		g             graph.Graph
	)

	s := []steps.Step{
		steps.Action(func(ctx context.Context) error {
			var err error
			installConfig, image, err = m.generateInstallConfig(ctx)
			return err
		}),

		steps.Action(func(ctx context.Context) error {
			var err error
			// Applies ARO-specific customisations to the InstallConfig
			g, err = m.applyInstallConfigCustomisations(installConfig, image)
			return err
		}),

		steps.Action(func(ctx context.Context) error {
			var err error
			// Applies ARO-specific customisations to the InstallConfig
			g, err = m.applyIgnitionConfigCustomisations(g)
			return err
		}),

		steps.Action(func(ctx context.Context) error {
			return m.persistGraph(ctx, g)
		}),
	}

	err := steps.Run(ctx, m.log, 10*time.Second, s)
	return g, err
}

func (m *manager) Install(ctx context.Context) error {
	s := []steps.Step{
		steps.AuthorizationRetryingAction(m.fpAuthorizer, m.deployResourceTemplate),
		steps.Action(m.initializeKubernetesClients),
		steps.Condition(m.bootstrapConfigMapReady, 60*time.Minute, true),
	}

	err := steps.Run(ctx, m.log, 10*time.Second, s)
	return err
}

// initializeKubernetesClients initializes clients using the Installer-generated
// kubeconfig.
func (m *manager) initializeKubernetesClients(ctx context.Context) error {
	resourceGroup := stringutils.LastTokenByte(m.oc.Properties.ClusterProfile.ResourceGroupID, '/')
	account := "cluster" + m.oc.Properties.StorageSuffix

	// Load the installer's generated kubeconfig from the graph
	pg, err := m.graph.LoadPersisted(ctx, resourceGroup, account)
	if err != nil {
		return err
	}

	var adminInternalClient *kubeconfig.AdminInternalClient
	err = pg.Get(&adminInternalClient)
	if err != nil {
		return err
	}

	// must not proceed if PrivateEndpointIP is not set.  In
	// k8s.io/client-go/transport/cache.go, k8s caches our transport, and it
	// can't tell if data in the restconfig.Dial closure has changed.  We don't
	// want it to cache a transport that can never work.
	if m.oc.Properties.NetworkProfile.APIServerPrivateEndpointIP == "" {
		return errors.New("privateEndpointIP is empty")
	}

	config, err := clientcmd.Load(adminInternalClient.File.Data)
	if err != nil {
		return err
	}

	r, err := clientcmd.NewDefaultClientConfig(*config, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return err
	}
	r.Dial = restconfig.DialContext(m.env, m.oc)

	m.kubernetescli, err = kubernetes.NewForConfig(r)
	return err
}
