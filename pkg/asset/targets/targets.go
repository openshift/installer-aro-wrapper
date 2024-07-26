package targets

import (
	aroassets "github.com/openshift/installer-aro-wrapper/pkg/asset/ignition"
	aroign "github.com/openshift/installer-aro-wrapper/pkg/asset/ignition/bootstrap"
	aro "github.com/openshift/installer-aro-wrapper/pkg/template/content/bootkube"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/cluster"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/kubeconfig"
	"github.com/openshift/installer/pkg/asset/machines"
	"github.com/openshift/installer/pkg/asset/manifests"
	"github.com/openshift/installer/pkg/asset/password"
	"github.com/openshift/installer/pkg/asset/templates/content/bootkube"
	"github.com/openshift/installer/pkg/asset/templates/content/openshift"
	"github.com/openshift/installer/pkg/asset/tls"
)

var (
	// InstallConfig are the install-config targeted assets.
	InstallConfig = []asset.WritableAsset{
		&installconfig.InstallConfig{},
	}

	// Manifests are the manifests targeted assets.
	Manifests = []asset.WritableAsset{
		&aroign.Bootstrap{},
		&machines.Master{},
		&machines.Worker{},
		&manifests.Manifests{},
		&manifests.Openshift{},
	}

	// ManifestTemplates are the manifest-templates targeted assets.
	ManifestTemplates = []asset.WritableAsset{
		&bootkube.KubeCloudConfig{},
		&bootkube.MachineConfigServerTLSSecret{},
		&bootkube.CVOOverrides{},
		&bootkube.KubeSystemConfigmapRootCA{},
		&bootkube.OpenshiftConfigSecretPullSecret{},
		&aro.AROWorkerRegistries{},
		&aro.AROIngressService{},
		&aro.AROImageRegistry{},
		&openshift.CloudCredsSecret{},
		&openshift.KubeadminPasswordSecret{},
		&openshift.RoleCloudCredsSecretReader{},
		&openshift.AzureCloudProviderSecret{},
	}

	// IgnitionConfigs are the ignition-configs targeted assets.
	IgnitionConfigs = []asset.WritableAsset{
		&kubeconfig.AdminClient{},
		&password.KubeadminPassword{},
		&aroassets.Master{},
		&aroassets.Worker{},
		&aroign.Bootstrap{},
		&cluster.Metadata{},
	}

	// Cluster are the cluster targeted assets.
	Cluster = []asset.WritableAsset{
		&cluster.Metadata{},
		&machine.MasterIgnitionCustomizations{},
		&machine.WorkerIgnitionCustomizations{},
		&cluster.TerraformVariables{},
		&kubeconfig.AdminClient{},
		&password.KubeadminPassword{},
		&tls.JournalCertKey{},
	}
)
