package installer

import (
	"context"

	"github.com/openshift/installer-aro-wrapper/pkg/asset/targets"
	"github.com/openshift/installer-aro-wrapper/pkg/cluster/graph"
)

// Add new ARO assets
// Create the ignition assets (targets)
func (m *manager) InjectIgnitionFiles(ctx context.Context, g graph.Graph) (graph.Graph, error) {
	// Resolve the asset graph.
	m.log.Print("resolving graph")
	for _, a := range targets.IgnitionConfigs {
		err := g.Resolve(a)
		if err != nil {
			return nil, err
		}
	}
	return g, nil
}

// Create cluster assets.
func (m *manager) InjectClusterFiles(ctx context.Context, g graph.Graph) (graph.Graph, error) {
	m.log.Print("resolving graph")
	for _, a := range targets.Cluster {
		err := g.Resolve(a)
		if err != nil {
			return nil, err
		}
	}
	return g, nil
}
