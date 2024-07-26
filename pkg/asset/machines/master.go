package machines

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	machinev1 "github.com/openshift/api/machine/v1"
	machinev1beta1 "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	icazure "github.com/openshift/installer/pkg/asset/installconfig/azure"
	"github.com/openshift/installer/pkg/asset/machines/azure"
	"github.com/openshift/installer/pkg/asset/machines/machineconfig"
	"github.com/openshift/installer/pkg/asset/rhcos"
	"github.com/openshift/installer/pkg/asset/templates/content/bootkube"
	"github.com/openshift/installer/pkg/types"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/azure/defaults"
	mcv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	"github.com/openshift/installer-aro-wrapper/pkg/asset/ignition"
	"github.com/openshift/installer-aro-wrapper/pkg/dnsmasq"
)

// Master generates the machines for the `Master` machine pool.
type Master struct {
	UserDataFile       *asset.File
	MachineConfigFiles []*asset.File
	MachineFiles       []*asset.File
	MasterMachineSet   *asset.File

	// SecretFiles is used by the baremetal platform to register the
	// credential information for communicating with management
	// controllers on hosts.
	SecretFiles []*asset.File

	// NetworkConfigSecretFiles is used by the baremetal platform to
	// store the networking configuration per host
	NetworkConfigSecretFiles []*asset.File

	// HostFiles is the list of baremetal hosts provided in the
	// installer configuration.
	HostFiles []*asset.File
}

const (
	directory = "openshift"

	// secretFileName is the format string for constructing the Secret
	// filenames for baremetal clusters.
	secretFileName = "99_openshift-cluster-api_host-bmc-secrets-%s.yaml"

	// networkConfigSecretFileName is the format string for constructing
	// the networking configuration Secret filenames for baremetal
	// clusters.
	networkConfigSecretFileName = "99_openshift-cluster-api_host-network-config-secrets-%s.yaml"

	// hostFileName is the format string for constucting the Host
	// filenames for baremetal clusters.
	hostFileName = "99_openshift-cluster-api_hosts-%s.yaml"

	// MasterMachineFileName is the format string for constucting the
	// Master Machine filenames.
	MasterMachineFileName = "99_openshift-cluster-api_Master-machines-%s.yaml"

	// MasterUserDataFileName is the filename used for the Master
	// user-data secret.
	MasterUserDataFileName = "99_openshift-cluster-api_Master-user-data-secret.yaml"

	// MasterUserDataFileName is the filename used for the control plane machine sets.
	MasterMachineSetFileName = "99_openshift-machine-api_Master-control-plane-machine-set.yaml"
)

var (
	secretFileNamePattern              = fmt.Sprintf(secretFileName, "*")
	networkConfigSecretFileNamePattern = fmt.Sprintf(networkConfigSecretFileName, "*")
	hostFileNamePattern                = fmt.Sprintf(hostFileName, "*")
	MasterMachineFileNamePattern       = fmt.Sprintf(MasterMachineFileName, "*")

	_ asset.WritableAsset = (*Master)(nil)
)

// Name returns a human friendly name for the Master Asset.
func (m *Master) Name() string {
	return "Master Machines"
}
func defaultAzureMachinePoolPlatform() azuretypes.MachinePool {
	return azuretypes.MachinePool{
		OSDisk: azuretypes.OSDisk{
			DiskSizeGB: powerOfTwoRootVolumeSize,
			DiskType:   azuretypes.DefaultDiskType,
		},
	}
}

// Dependencies returns all of the dependencies directly needed by the
// Master asset
func (m *Master) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.ClusterID{},
		// PlatformCredsCheck just checks the creds (and asks, if needed)
		// We do not actually use it in this asset directly, hence
		// it is put in the dependencies but not fetched in Generate
		&installconfig.PlatformCredsCheck{},
		&installconfig.InstallConfig{},
		new(rhcos.Image),
		&ignition.Master{},
		&bootkube.ARODNSConfig{},
	}
}

// Generate generates the Master asset.
func (m *Master) Generate(dependencies asset.Parents) error {
	clusterID := &installconfig.ClusterID{}
	installConfig := &installconfig.InstallConfig{}
	rhcosImage := new(rhcos.Image)
	mign := &ignition.Master{}
	aroDNSConfig := &bootkube.ARODNSConfig{}
	dependencies.Get(clusterID, installConfig, rhcosImage, mign, aroDNSConfig)

	MasterUserDataSecretName := "Master-user-data"

	ic := installConfig.Config

	pool := *ic.ControlPlane
	var err error
	machines := []machinev1beta1.Machine{}
	var MasterMachineSet *machinev1.ControlPlaneMachineSet
	mpool := defaultAzureMachinePoolPlatform()
	mpool.InstanceType = defaults.ControlPlaneInstanceType(
		installConfig.Config.Platform.Azure.CloudName,
		installConfig.Config.Platform.Azure.Region,
		installConfig.Config.ControlPlane.Architecture,
	)
	mpool.OSDisk.DiskSizeGB = 1024
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
	useImageGallery := installConfig.Azure.CloudName != azuretypes.StackCloud

	machines, cpms, err := azure.Machines(clusterID.InfraID, ic, &pool, string(*rhcosImage), "Master", MasterUserDataSecretName, capabilities, useImageGallery)
	if err != nil {
		return errors.Wrap(err, "failed to create Master machine objects")
	}
	azure.ConfigMasters(machines, cpms, clusterID.InfraID)

	data, err := userDataSecret(MasterUserDataSecretName, mign.File.Data)
	if err != nil {
		return errors.Wrap(err, "failed to create user-data secret for Master machines")
	}

	m.UserDataFile = &asset.File{
		Filename: filepath.Join(directory, MasterUserDataFileName),
		Data:     data,
	}

	machineConfigs := []*mcv1.MachineConfig{}
	if pool.Hyperthreading == types.HyperthreadingDisabled {
		ignHT, err := machineconfig.ForHyperthreadingDisabled("Master")
		if err != nil {
			return errors.Wrap(err, "failed to create ignition for hyperthreading disabled for Master machines")
		}
		machineConfigs = append(machineConfigs, ignHT)
	}
	if ic.SSHKey != "" {
		ignSSH, err := machineconfig.ForAuthorizedKeys(ic.SSHKey, "Master")
		if err != nil {
			return errors.Wrap(err, "failed to create ignition for authorized SSH keys for Master machines")
		}
		machineConfigs = append(machineConfigs, ignSSH)
	}
	if ic.FIPS {
		ignFIPS, err := machineconfig.ForFIPSEnabled("Master")
		if err != nil {
			return errors.Wrap(err, "failed to create ignition for FIPS enabled for Master machines")
		}
		machineConfigs = append(machineConfigs, ignFIPS)
	}
	ignARODNS, err := dnsmasq.MachineConfig(installConfig.Config.ClusterDomain(), aroDNSConfig.APIIntIP, aroDNSConfig.IngressIP, "Master", aroDNSConfig.GatewayDomains, aroDNSConfig.GatewayPrivateEndpointIP, true)
	if err != nil {
		return errors.Wrap(err, "failed to create ignition for ARO DNS for Master machines")
	}
	machineConfigs = append(machineConfigs, ignARODNS)

	m.MachineConfigFiles, err = machineconfig.Manifests(machineConfigs, "Master", directory)
	if err != nil {
		return errors.Wrap(err, "failed to create MachineConfig manifests for Master machines")
	}

	m.MachineFiles = make([]*asset.File, len(machines))
	if MasterMachineSet != nil && *pool.Replicas > 1 {
		data, err := yaml.Marshal(MasterMachineSet)
		if err != nil {
			return errors.Wrapf(err, "marshal control plane machine set")
		}
		m.MasterMachineSet = &asset.File{
			Filename: filepath.Join(directory, MasterMachineSetFileName),
			Data:     data,
		}
	}
	padFormat := fmt.Sprintf("%%0%dd", len(fmt.Sprintf("%d", len(machines))))
	for i, machine := range machines {
		data, err := yaml.Marshal(machine)
		if err != nil {
			return errors.Wrapf(err, "marshal Master %d", i)
		}

		padded := fmt.Sprintf(padFormat, i)
		m.MachineFiles[i] = &asset.File{
			Filename: filepath.Join(directory, fmt.Sprintf(MasterMachineFileName, padded)),
			Data:     data,
		}
	}
	return nil
}

// Files returns the files generated by the asset.
func (m *Master) Files() []*asset.File {
	files := make([]*asset.File, 0, 1+len(m.MachineConfigFiles)+len(m.MachineFiles))
	if m.UserDataFile != nil {
		files = append(files, m.UserDataFile)
	}
	files = append(files, m.MachineConfigFiles...)
	// Hosts refer to secrets, so place the secrets before the hosts
	// to avoid unnecessary reconciliation errors.
	files = append(files, m.SecretFiles...)
	files = append(files, m.NetworkConfigSecretFiles...)
	// Machines are linked to hosts via the machineRef, so we create
	// the hosts first to ensure if the operator starts trying to
	// reconcile a machine it can pick up the related host.
	files = append(files, m.HostFiles...)
	files = append(files, m.MachineFiles...)
	if m.MasterMachineSet != nil {
		files = append(files, m.MasterMachineSet)
	}
	return files
}

// Load reads the asset files from disk.
func (m *Master) Load(f asset.FileFetcher) (found bool, err error) {
	file, err := f.FetchByName(filepath.Join(directory, MasterUserDataFileName))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	m.UserDataFile = file

	m.MachineConfigFiles, err = machineconfig.Load(f, "Master", directory)
	if err != nil {
		return true, err
	}

	var fileList []*asset.File

	fileList, err = f.FetchByPattern(filepath.Join(directory, secretFileNamePattern))
	if err != nil {
		return true, err
	}
	m.SecretFiles = fileList

	fileList, err = f.FetchByPattern(filepath.Join(directory, networkConfigSecretFileNamePattern))
	if err != nil {
		return true, err
	}
	m.NetworkConfigSecretFiles = fileList

	fileList, err = f.FetchByPattern(filepath.Join(directory, hostFileNamePattern))
	if err != nil {
		return true, err
	}
	m.HostFiles = fileList

	fileList, err = f.FetchByPattern(filepath.Join(directory, MasterMachineFileNamePattern))
	if err != nil {
		return true, err
	}
	m.MachineFiles = fileList

	file, err = f.FetchByName(filepath.Join(directory, MasterMachineSetFileName))
	if err != nil {
		if os.IsNotExist(err) {
			// Choosing to ignore the CPMS file if it does not exist since UPI does not need it.
			logrus.Debugf("CPMS file missing. Ignoring it while loading machine asset.")
			return true, nil
		}
		return true, err
	}
	m.MasterMachineSet = file

	return true, nil
}

// Machines returns Master Machine manifest structures.
func (m *Master) Machines() ([]machinev1beta1.Machine, error) {
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(machinev1beta1.SchemeGroupVersion,
		&machinev1beta1.AzureMachineProviderSpec{},
	)
	scheme.AddKnownTypes(machinev1.GroupVersion,
		&machinev1.ControlPlaneMachineSet{},
	)

	machinev1beta1.AddToScheme(scheme)
	machinev1.Install(scheme)
	decoder := serializer.NewCodecFactory(scheme).UniversalDecoder(
		machinev1.GroupVersion,
	)

	machines := []machinev1beta1.Machine{}
	for i, file := range m.MachineFiles {
		machine := &machinev1beta1.Machine{}
		err := yaml.Unmarshal(file.Data, &machine)
		if err != nil {
			return machines, errors.Wrapf(err, "unmarshal Master %d", i)
		}

		obj, _, err := decoder.Decode(machine.Spec.ProviderSpec.Value.Raw, nil, nil)
		if err != nil {
			return machines, errors.Wrapf(err, "unmarshal Master %d", i)
		}

		machine.Spec.ProviderSpec.Value = &runtime.RawExtension{Object: obj}
		machines = append(machines, *machine)
	}

	return machines, nil
}

// IsMachineManifest tests whether a file is a manifest that belongs to the
// Master Machines or Worker Machines asset.
func IsMachineManifest(file *asset.File) bool {
	if filepath.Dir(file.Filename) != directory {
		return false
	}
	filename := filepath.Base(file.Filename)
	if filename == MasterUserDataFileName || filename == workerUserDataFileName || filename == MasterMachineSetFileName {
		return true
	}
	if matched, err := machineconfig.IsManifest(filename); err != nil {
		panic(err)
	} else if matched {
		return true
	}
	if matched, err := filepath.Match(MasterMachineFileNamePattern, filename); err != nil {
		panic("bad format for Master machine file name pattern")
	} else if matched {
		return true
	}
	if matched, err := filepath.Match(workerMachineSetFileNamePattern, filename); err != nil {
		panic("bad format for worker machine file name pattern")
	} else {
		return matched
	}
}
