package machines

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	machinev1 "github.com/openshift/api/machine/v1"
	machinev1beta1 "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/aro/dnsmasq"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	icazure "github.com/openshift/installer/pkg/asset/installconfig/azure"
	"github.com/openshift/installer/pkg/asset/machines/azure"
	"github.com/openshift/installer/pkg/asset/machines/machineconfig"
	"github.com/openshift/installer/pkg/asset/rhcos"
	"github.com/openshift/installer/pkg/asset/templates/content/bootkube"
	"github.com/openshift/installer/pkg/types"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
	azuredefaults "github.com/openshift/installer/pkg/types/azure/defaults"
	mcv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	"github.com/openshift/installer-aro-wrapper/pkg/asset/ignition"
)

const (
	// workerMachineSetFileName is the format string for constructing the worker MachineSet filenames.
	workerMachineSetFileName = "99_openshift-cluster-api_worker-machineset-%s.yaml"

	// workerUserDataFileName is the filename used for the worker user-data secret.
	workerUserDataFileName = "99_openshift-cluster-api_worker-user-data-secret.yaml"

	// decimalRootVolumeSize is the size in GB we use for some platforms.
	// See below.
	decimalRootVolumeSize = 120

	// powerOfTwoRootVolumeSize is the size in GB we use for other platforms.
	// The reasons for the specific choices between these two may boil down
	// to which section of code the person adding a platform was copy-pasting from.
	// https://github.com/openshift/openshift-docs/blob/main/modules/installation-requirements-user-infra.adoc#minimum-resource-requirements
	powerOfTwoRootVolumeSize = 128
)

var (
	workerMachineSetFileNamePattern = fmt.Sprintf(workerMachineSetFileName, "*")

	_ asset.WritableAsset = (*Worker)(nil)
)

// Worker generates the machinesets for `worker` machine pool.
type Worker struct {
	UserDataFile       *asset.File
	MachineConfigFiles []*asset.File
	MachineSetFiles    []*asset.File
}

// Name returns a human friendly name for the Worker Asset.
func (w *Worker) Name() string {
	return "Worker Machines"
}

// Dependencies returns all of the dependencies directly needed by the
// Worker asset
func (w *Worker) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.ClusterID{},
		// PlatformCredsCheck just checks the creds (and asks, if needed)
		// We do not actually use it in this asset directly, hence
		// it is put in the dependencies but not fetched in Generate
		&installconfig.PlatformCredsCheck{},
		&installconfig.InstallConfig{},
		new(rhcos.Image),
		new(rhcos.Release),
		&ignition.Worker{},
		&bootkube.ARODNSConfig{},
	}
}

// Generate generates the Worker asset.
func (w *Worker) Generate(dependencies asset.Parents) error {
	clusterID := &installconfig.ClusterID{}
	installConfig := &installconfig.InstallConfig{}
	rhcosImage := new(rhcos.Image)
	rhcosRelease := new(rhcos.Release)
	wign := &ignition.Worker{}
	aroDNSConfig := &bootkube.ARODNSConfig{}
	dependencies.Get(clusterID, installConfig, rhcosImage, rhcosRelease, wign, aroDNSConfig)

	workerUserDataSecretName := "worker-user-data"

	machineConfigs := []*mcv1.MachineConfig{}
	machineSets := []runtime.Object{}
	var err error
	ic := installConfig.Config
	for _, pool := range ic.Compute {
		pool := pool // this makes golint happy... G601: Implicit memory aliasing in for loop. (gosec)
		if pool.Hyperthreading == types.HyperthreadingDisabled {
			ignHT, err := machineconfig.ForHyperthreadingDisabled("worker")
			if err != nil {
				return errors.Wrap(err, "failed to create ignition for hyperthreading disabled for worker machines")
			}
			machineConfigs = append(machineConfigs, ignHT)
		}
		if ic.SSHKey != "" {
			ignSSH, err := machineconfig.ForAuthorizedKeys(ic.SSHKey, "worker")
			if err != nil {
				return errors.Wrap(err, "failed to create ignition for authorized SSH keys for worker machines")
			}
			machineConfigs = append(machineConfigs, ignSSH)
		}
		if ic.FIPS {
			ignFIPS, err := machineconfig.ForFIPSEnabled("worker")
			if err != nil {
				return errors.Wrap(err, "failed to create ignition for FIPS enabled for worker machines")
			}
			machineConfigs = append(machineConfigs, ignFIPS)
		}
		ignARODNS, err := dnsmasq.MachineConfig(installConfig.Config.ClusterDomain(), aroDNSConfig.APIIntIP, aroDNSConfig.IngressIP, "worker", aroDNSConfig.GatewayDomains, aroDNSConfig.GatewayPrivateEndpointIP, true)
		if err != nil {
			return errors.Wrap(err, "failed to create ignition for ARO DNS for worker machines")
		}
		machineConfigs = append(machineConfigs, ignARODNS)
		mpool := defaultAzureMachinePoolPlatform()
		mpool.InstanceType = azuredefaults.ComputeInstanceType(
			installConfig.Config.Platform.Azure.CloudName,
			installConfig.Config.Platform.Azure.Region,
			pool.Architecture,
		)
		mpool.Set(ic.Platform.Azure.DefaultMachinePlatform)
		mpool.Set(pool.Platform.Azure)

		session, err := installConfig.Azure.Session()
		if err != nil {
			return errors.Wrap(err, "failed to fetch session")
		}

		// Default to current subscription if one was not specified
		if mpool.OSDisk.DiskEncryptionSet != nil && mpool.OSDisk.DiskEncryptionSet.SubscriptionID == "" {
			mpool.OSDisk.DiskEncryptionSet.SubscriptionID = session.Credentials.SubscriptionID
		}

		client := icazure.NewClient(session)
		if len(mpool.Zones) == 0 {
			azs, err := client.GetAvailabilityZones(context.TODO(), ic.Platform.Azure.Region, mpool.InstanceType)
			if err != nil {
				return errors.Wrap(err, "failed to fetch availability zones")
			}
			mpool.Zones = azs
			if len(azs) == 0 {
				// if no azs are given we set to []string{""} for convenience over later operations.
				// It means no-zoned for the machine API
				mpool.Zones = []string{""}
			}
		}

		pool.Platform.Azure = &mpool

		capabilities, err := client.GetVMCapabilities(context.TODO(), mpool.InstanceType, installConfig.Config.Platform.Azure.Region)
		if err != nil {
			return err
		}

		useImageGallery := ic.Platform.Azure.CloudName != azuretypes.StackCloud
		sets, err := azure.MachineSets(clusterID.InfraID, ic, &pool, string(*rhcosImage), "worker", workerUserDataSecretName, capabilities, useImageGallery)
		if err != nil {
			return errors.Wrap(err, "failed to create worker machine objects")
		}
		for _, set := range sets {
			machineSets = append(machineSets, set)
		}
	}

	data, err := userDataSecret(workerUserDataSecretName, wign.File.Data)
	if err != nil {
		return errors.Wrap(err, "failed to create user-data secret for worker machines")
	}
	w.UserDataFile = &asset.File{
		Filename: filepath.Join(directory, workerUserDataFileName),
		Data:     data,
	}

	w.MachineConfigFiles, err = machineconfig.Manifests(machineConfigs, "worker", directory)
	if err != nil {
		return errors.Wrap(err, "failed to create MachineConfig manifests for worker machines")
	}

	w.MachineSetFiles = make([]*asset.File, len(machineSets))
	padFormat := fmt.Sprintf("%%0%dd", len(fmt.Sprintf("%d", len(machineSets))))
	for i, machineSet := range machineSets {
		data, err := yaml.Marshal(machineSet)
		if err != nil {
			return errors.Wrapf(err, "marshal worker %d", i)
		}

		padded := fmt.Sprintf(padFormat, i)
		w.MachineSetFiles[i] = &asset.File{
			Filename: filepath.Join(directory, fmt.Sprintf(workerMachineSetFileName, padded)),
			Data:     data,
		}
	}
	return nil
}

// Files returns the files generated by the asset.
func (w *Worker) Files() []*asset.File {
	files := make([]*asset.File, 0, 1+len(w.MachineConfigFiles)+len(w.MachineSetFiles))
	if w.UserDataFile != nil {
		files = append(files, w.UserDataFile)
	}
	files = append(files, w.MachineConfigFiles...)
	files = append(files, w.MachineSetFiles...)
	return files
}

// Load reads the asset files from disk.
func (w *Worker) Load(f asset.FileFetcher) (found bool, err error) {
	file, err := f.FetchByName(filepath.Join(directory, workerUserDataFileName))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	w.UserDataFile = file

	w.MachineConfigFiles, err = machineconfig.Load(f, "worker", directory)
	if err != nil {
		return true, err
	}

	fileList, err := f.FetchByPattern(filepath.Join(directory, workerMachineSetFileNamePattern))
	if err != nil {
		return true, err
	}

	w.MachineSetFiles = fileList
	return true, nil
}

// MachineSets returns MachineSet manifest structures.
func (w *Worker) MachineSets() ([]machinev1beta1.MachineSet, error) {
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(machinev1beta1.SchemeGroupVersion,
		&machinev1beta1.AzureMachineProviderSpec{},
	)
	machinev1.Install(scheme)
	decoder := serializer.NewCodecFactory(scheme).UniversalDecoder(
		machinev1.GroupVersion,
	)

	machineSets := []machinev1beta1.MachineSet{}
	for i, file := range w.MachineSetFiles {
		machineSet := &machinev1beta1.MachineSet{}
		err := yaml.Unmarshal(file.Data, &machineSet)
		if err != nil {
			return machineSets, errors.Wrapf(err, "unmarshal worker %d", i)
		}

		obj, _, err := decoder.Decode(machineSet.Spec.Template.Spec.ProviderSpec.Value.Raw, nil, nil)
		if err != nil {
			return machineSets, errors.Wrapf(err, "unmarshal worker %d", i)
		}

		machineSet.Spec.Template.Spec.ProviderSpec.Value = &runtime.RawExtension{Object: obj}
		machineSets = append(machineSets, *machineSet)
	}

	return machineSets, nil
}
