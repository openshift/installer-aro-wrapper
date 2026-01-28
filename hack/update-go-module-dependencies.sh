#!/bin/bash -e

export GOPRIVATE=github.com
export GONOPROXY="y"

RELEASE=release-4.20
VM_SKU=aro_4$(echo $RELEASE | sed 's/.*\.//')
declare -a pinned=(
	"github.com/openshift/assisted-service/api"
	"github.com/openshift/assisted-service/client"
	"github.com/openshift/assisted-service/models"
	"github.com/openshift/api"
	"github.com/openshift/client-go"
)

read -p "This will update to $RELEASE. Is this correct? " -n 1
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then

	echo "||| Updating installer"
	go get github.com/openshift/installer@$RELEASE

	echo "||| Updating pinned dependencies"
	for i in "${pinned[@]}"; do
		go mod edit -replace $i=$(go list -mod=mod -m $i@$RELEASE | sed -e 's/ /@/')
	done
	go mod edit -dropreplace github.com/openshift/hive/apis

	echo "||| Running go mod tidy"
	go mod tidy

	RHCOS_VERSION_RAW=$(az vm image list --publisher azureopenshift --offer aro4 --sku $VM_SKU --all --query "sort_by([], &version)[-1].version" -o tsv)

	# Extract version number from RELEASE (e.g., release-4.20 -> 420)
	VERSION_PREFIX=$(echo ${RELEASE#release-} | sed 's/\.//')

	# Transform version: if Azure returns "9.6.20251015", convert to "420.6.20251015"
	# This handles cases where Azure may use a shortened prefix
	RHCOS_VERSION=$(echo "$RHCOS_VERSION_RAW" | sed "s/^[0-9]\+/$VERSION_PREFIX/")

	echo "Update pkg/installer/generateconfig.go 's rhcosImage struct with:"
	echo "SKU: \"$VM_SKU\""
	echo "Version: \"$RHCOS_VERSION\","
fi
