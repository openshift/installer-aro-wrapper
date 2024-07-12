package aroign

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/Azure/go-autorest/autorest/to"
	ign3types "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/openshift/installer/pkg/asset/ignition"
	mcfgv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
	"github.com/pkg/errors"
	"github.com/vincent-petithory/dataurl"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var etcHostsTemplate = template.Must(template.New("etchosts").Parse(`127.0.0.1	localhost localhost.localdomain localhost4 localhost4.localdomain4
::1	localhost localhost.localdomain localhost6 localhost6.localdomain6
{{ .APIIntIP }}	api.{{ .ClusterDomain }} api-int.{{ .ClusterDomain }}
{{ $.GatewayPrivateEndpointIP }}	{{ range $GatewayDomain := .GatewayDomains }}{{ $GatewayDomain }} {{ end }}
`))

type etcHostsTemplateData struct {
	ClusterDomain            string
	APIIntIP                 string
	GatewayDomains           []string
	GatewayPrivateEndpointIP string
}

func GenerateEtcHostsAdditionalDomains(clusterDomain string, apiIntIP string, gatewayDomains []string, gatewayPrivateEndpointIP string) ([]byte, error) {
	buf := &bytes.Buffer{}
	templateData := etcHostsTemplateData{
		ClusterDomain:            clusterDomain,
		APIIntIP:                 apiIntIP,
		GatewayDomains:           gatewayDomains,
		GatewayPrivateEndpointIP: gatewayPrivateEndpointIP,
	}

	if err := etcHostsTemplate.Execute(buf, templateData); err != nil {
		return nil, errors.Wrap(err, "failed to execute etc hosts template")
	}

	return buf.Bytes(), nil
}

func EtcHostsIgnitionConfig(clusterDomain string, apiIntIP string, gatewayDomains []string, gatewayPrivateEndpointIP string) (*ign3types.Config, error) {
	data, err := GenerateEtcHostsAdditionalDomains(clusterDomain, apiIntIP, gatewayDomains, gatewayPrivateEndpointIP)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate addtional hosts for etc hosts")
	}

	ign := &ign3types.Config{
		Ignition: ign3types.Ignition{
			Version: ign3types.MaxVersion.String(),
		},
		Storage: ign3types.Storage{
			Files: []ign3types.File{
				{
					Node: ign3types.Node{
						Path:      "/etc/hosts",
						Overwrite: to.BoolPtr(true),
					},
					FileEmbedded1: ign3types.FileEmbedded1{
						Contents: ign3types.Resource{
							Source: to.StringPtr(dataurl.EncodeBytes(data)),
						},
						Mode: to.IntPtr(0644),
					},
				},
			},
		},
	}

	return ign, nil
}

func EtcHostsMachineConfig(clusterDomain string, apiIntIP string, gatewayDomains []string, gatewayPrivateEndpointIP string, role string) (*mcfgv1.MachineConfig, error) {
	ignConfig, err := EtcHostsIgnitionConfig(clusterDomain, apiIntIP, gatewayDomains, gatewayPrivateEndpointIP)
	if err != nil {
		return nil, err
	}

	// marshalled, err := json.Marshal(ignConfig)
	// fmt.Printf("marshalled: %s\n", marshalled)

	rawExt, err := ignition.ConvertToRawExtension(*ignConfig)
	if err != nil {
		return nil, err
	}

	return &mcfgv1.MachineConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: mcfgv1.SchemeGroupVersion.String(),
			Kind:       "MachineConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("99-%s-aro-etc-hosts-gateway-domains", role),
			Labels: map[string]string{
				"machineconfiguration.openshift.io/role": role,
			},
		},
		Spec: mcfgv1.MachineConfigSpec{
			Config: rawExt,
		},
	}, nil

}
