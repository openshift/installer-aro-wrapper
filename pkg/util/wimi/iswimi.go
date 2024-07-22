package wimi

import "github.com/openshift/ARO-Installer/pkg/api"

// IsWimi checks whether a cluster is a Workload Identity cluster or a Service Principal cluster
func IsWimi(oc *api.OpenShiftCluster) bool {
	if oc.Properties.PlatformWorkloadIdentityProfile != nil && oc.Properties.ServicePrincipalProfile == nil {
		return true
	}
	return false
}
