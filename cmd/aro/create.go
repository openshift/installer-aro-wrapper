package main

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"encoding/json"
	"os"

	"github.com/openshift/installer/pkg/asset"
	targetassets "github.com/openshift/installer/pkg/asset/targets"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/openshift/installer-aro-wrapper/pkg/api"
	"github.com/openshift/installer-aro-wrapper/pkg/installer"
)

type target struct {
	name    string
	command *cobra.Command
}

// each target is a variable to preserve the order when creating subcommands and still
// allow other functions to directly access each target individually.
var (
	manifestsTarget = target{
		name: "Manifests",
		command: &cobra.Command{
			Use:   "manifests",
			Short: "Generates the Kubernetes manifests",
			Run: func(cmd *cobra.Command, args []string) {
				ctx := context.Background()
				log := logrus.NewEntry(logrus.StandardLogger())
				i, err := _makeInstaller(ctx, log, rootOpts.dir)
				if err != nil {
					logrus.Error(err)
					logrus.Exit(1)
				}
				g, err := i.Manifests(ctx)
				if err != nil {
					logrus.Error(err)
					logrus.Exit(1)
				}

				runner := func(directory string, manifests []asset.WritableAsset) error {
					for _, m := range manifests {
						err = g.Resolve(ctx, m)
						if err != nil {
							err = errors.Wrapf(err, "failed to fetch %s", m.Name())
						}

						a := g.Get(m).(asset.WritableAsset)
						if err2 := asset.PersistToFile(a, directory); err2 != nil {
							err2 = errors.Wrapf(err2, "failed to write asset (%s) to disk", a.Name())
							if err != nil {
								logrus.Error(err2)
								return err
							}
							return err2
						}
					}
					return nil
				}

				err = runner(rootOpts.dir, targetassets.Manifests)
				if err != nil {
					logrus.Error(err)
					logrus.Exit(1)
				}

				err = runner(rootOpts.dir, targetassets.IgnitionConfigs)
				if err != nil {
					logrus.Error(err)
					logrus.Exit(1)
				}
			},
		},
	}

	ignitionConfigsTarget = target{
		name: "Ignition Configs",
		command: &cobra.Command{
			Use:   "ignition-configs",
			Short: "Generates the Ignition Config asset",
			// FIXME: add longer descriptions for our commands with examples for better UX.
			// Long:  "",
		},
	}
	clusterTarget = target{
		name: "Cluster",
		command: &cobra.Command{
			Use:   "cluster",
			Short: "Create an OpenShift cluster",
			Run: func(cmd *cobra.Command, args []string) {
				ctx := context.Background()
				log := logrus.NewEntry(logrus.StandardLogger())
				i, err := _makeInstaller(ctx, log, rootOpts.dir)
				if err != nil {
					logrus.Error(err)
					logrus.Exit(1)
				}

				err = i.Install(ctx)
				if err != nil {
					logrus.Error(err)
					logrus.Exit(1)
				}
			},
		},
	}

	targets = []target{manifestsTarget, ignitionConfigsTarget, clusterTarget}
)

func newCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create part of an OpenShift cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	for _, t := range targets {
		t.command.Args = cobra.ExactArgs(0)
		cmd.AddCommand(t.command)
	}

	return cmd
}

func _getOpenShiftCluster() (*api.OpenShiftCluster, error) {
	var oc api.OpenShiftCluster
	ocFile, err := os.ReadFile("/.azure/99_aro.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(ocFile, &oc)
	if err != nil {
		return nil, err
	}

	return &oc, nil
}

func _getSubscription() (*api.Subscription, error) {
	var sub api.Subscription
	subFile, err := os.ReadFile("/.azure/99_sub.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(subFile, &sub)
	if err != nil {
		return nil, err
	}

	return &sub, nil
}

func _makeInstaller(ctx context.Context, log *logrus.Entry, assetsDir string) (installer.Interface, error) {
	var err error

	// unmarshal oc
	var oc *api.OpenShiftCluster
	oc, err = _getOpenShiftCluster()
	if err != nil {
		return nil, err
	}

	// unmarshal subscription
	var sub *api.Subscription
	sub, err = _getSubscription()
	if err != nil {
		return nil, err
	}

	return installer.MakeInstaller(ctx, log, assetsDir, oc, sub)
}
