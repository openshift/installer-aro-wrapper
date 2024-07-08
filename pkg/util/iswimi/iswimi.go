package iswimi

import "github.com/openshift/ARO-Installer/pkg/api"

// IsWimi checks whether a cluster is Workload Identity or classic
func IsWimi(cluster api.OpenShiftClusterProperties) bool {
	if cluster.PlatformWorkloadIdentityProfile == nil || cluster.ServicePrincipalProfile != nil {
		return false
	}
	return true
}
