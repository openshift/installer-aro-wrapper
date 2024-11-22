package clusterapi

import (
	"context"

	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	ibmcloudtypes "github.com/openshift/installer/pkg/types/ibmcloud"
)

var _ clusterapi.PreProvider = (*Provider)(nil)
var _ clusterapi.Provider = (*Provider)(nil)

// Provider implements IBM Cloud CAPI installation.
type Provider struct{}

// Name returns the IBM Cloud provider name.
func (p Provider) Name() string {
	return ibmcloudtypes.Name
}

// PublicGatherEndpoint indicates that machine ready checks should NOT wait for an ExternalIP
// in the status when declaring machines ready.
func (Provider) PublicGatherEndpoint() clusterapi.GatherEndpoint { return clusterapi.InternalIP }

// PreProvision creates the IBM Cloud objects required prior to running capibmcloud.
func (p Provider) PreProvision(ctx context.Context, in clusterapi.PreProvisionInput) error {
	return nil
}
