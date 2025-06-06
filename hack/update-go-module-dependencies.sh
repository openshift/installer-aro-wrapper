#!/bin/bash -ex

# Background: https://groups.google.com/forum/#!topic/golang-nuts/51-D_YFC78k
#
# TLDR: OCP consists of many repos where for each release a new release branch gets created (release-X.Y).
# When we update vendors we want to get latest changes from the release branch for all of the dependencies.
# With Go modules we can't easily do it, but there is a workaround which consists of multiple steps:
# 	1. Get the latest commit from the branch using `go list -mod=mod -m MODULE@release-x.y`.
# 	2. Using `sed`, transform output of the above command into format accepted by `go mod edit -replace`.
#	3. Modify `go.mod` by calling `go mod edit -replace`.
#
# This needs to happen for each module that uses this branching strategy: all these repos
# live under github.com/openshift organisation.
#
# There are however, some exceptions:
# 	* Some repos under github.com/openshift do not use this strategy.
#     We should skip them in this script and manage directly with `go mod`.
# 	* Some dependencies pin their own dependencies to older commits.
#     For example, dependency Foo from release-4.7 branch requires
#	  dependency Bar at older commit which is
#     not compatible with Bar@release-4.7.
#
# Note that github.com/openshift org also contains forks of K8s upstream repos and we
# use these forks (indirectly in most cases). This means that
# we also must take care of replacing modules such as  github.com/metal3-io/baremetal-operator
# with github.com/openshift/baremetal-operator (just an example, there are more).

RELEASE=release-4.15
K8S_RELEASE=v0.28.3
GO_VERSION=1.20

for x in vendor/github.com/openshift/*; do
	case $x in
		# Review the list of special cases on each release.

		# Do not update Hive: it is not part of OCP
		vendor/github.com/openshift/hive)
			;;

		# Don't use replace directive for the installer.
		vendor/github.com/openshift/installer)
			;;

		# It is only used indirectly and intermediate dependencies pin to different incompatible commits.
		# We force a specific commit here to make all dependencies happy.
		vendor/github.com/openshift/cloud-credential-operator)
			go mod edit -replace github.com/openshift/cloud-credential-operator=github.com/openshift/cloud-credential-operator@v0.0.0-20240422222427-55199c9b5870
			;;

		# We can't use MCO 4.15 yet because it doesn't contain MCO's API anymore, it moved to openshift/api.
		# Upstream installer still uses that API from MCO repo though
		# Pin to release 4.14
		vendor/github.com/openshift/machine-config-operator)
			go mod edit -replace github.com/openshift/machine-config-operator=github.com/openshift/machine-config-operator@release-4.14
			;;

		# This repo doesn't follow release-x.y branching strategy
		# We skip it and let go mod resolve it
		vendor/github.com/openshift/custom-resource-status)
			;;

		*)
			go mod edit -replace "${x##vendor/}"="$(go list -mod=mod -m ${x##vendor/}@$RELEASE | sed -e 's/ /@/')"
			;;
	esac
done

for x in vendor/k8s.io/*; do
  case $x in
    # skip, it's replaced by openshift
    vendor/k8s.io/cloud-provider-vsphere)
      ;;
    # skip, they don't follow k8s versioning schema
    vendor/k8s.io/gengo|vendor/k8s.io/klog|vendor/k8s.io/kube-openapi|vendor/k8s.io/utils)
      ;;
    *)
      go mod edit -replace "${x##vendor/}"="$(go list -mod=mod -m ${x##vendor/}@$K8S_RELEASE | sed -e 's/ /@/')"
      ;;
  esac
done

# From installer(-aro), they don't use forks anymore!
go mod edit -replace sigs.k8s.io/cluster-api=sigs.k8s.io/cluster-api@v1.5.3
go mod edit -replace sigs.k8s.io/cluster-api-provider-aws/v2=sigs.k8s.io/cluster-api-provider-aws/v2@v2.0.0-20231024062453-0bf78b04b305
go mod edit -replace sigs.k8s.io/cluster-api-provider-azure=sigs.k8s.io/cluster-api-provider-azure@v1.11.1-0.20231026140308-a3f4914170d9

for x in baremetal-operator baremetal-operator/apis baremetal-operator/pkg/hardwareutils cluster-api-provider-baremetal cluster-api-provider-metal3 cluster-api-provider-metal3/api; do
  go mod edit -replace github.com/metal3-io/$x="$(go list -mod=mod -m github.com/openshift/$x@$RELEASE | sed -e 's/ /@/')"
done

go get github.com/openshift/installer@$RELEASE

go get ./...

go mod tidy -compat=$GO_VERSION
go mod vendor
