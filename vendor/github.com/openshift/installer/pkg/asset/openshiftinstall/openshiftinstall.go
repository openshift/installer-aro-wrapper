package openshiftinstall

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/version"
)

var (
	configPath = filepath.Join("openshift", "openshift-install-manifests.yaml")
)

// Config generates the openshift-install ConfigMap.
type Config struct {
	File *asset.File
}

var _ asset.WritableAsset = (*Config)(nil)

// Name returns a human friendly name for the asset.
func (*Config) Name() string {
	return "OpenShift Install (Manifests)"
}

// Dependencies returns all of the dependencies directly needed to generate
// the asset.
func (*Config) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates the openshift-install ConfigMap.
func (i *Config) Generate(dependencies asset.Parents) error {
	cm, err := CreateInstallConfigMap("openshift-install-manifests")
	if err != nil {
		return err
	}

	i.File = &asset.File{
		Filename: configPath,
		Data:     []byte(cm),
	}

	return nil
}

// Files returns the files generated by the asset.
func (i *Config) Files() []*asset.File {
	if i.File != nil {
		return []*asset.File{i.File}
	}
	return []*asset.File{}
}

// Load loads the already-rendered files back from disk.
func (i *Config) Load(f asset.FileFetcher) (bool, error) {
	file, err := f.FetchByName(configPath)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	i.File = file
	return true, nil
}

// CreateInstallConfigMap creates an openshift-install ConfigMap from the
// OPENSHIFT_INSTALL_INVOKER environment variable and the given name for the
// ConfigMap. This returns an error if marshalling to YAML fails.
func CreateInstallConfigMap(name string) (string, error) {
	invoker := "ARO"

	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "openshift-config",
			Name:      name,
		},
		Data: map[string]string{
			"version": version.Raw,
			"invoker": invoker,
		},
	}

	cmData, err := yaml.Marshal(cm)
	if err != nil {
		return "", errors.Wrapf(err, "failed to create %q ConfigMap", name)
	}

	return string(cmData), nil
}
