package aroign

import (
	"bytes"
	"text/template"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/blang/semver"
	ign3types "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/vincent-petithory/dataurl"
)

const (
	configFileName = "etchosts.conf"
)

func config(clusterDomain string, apiIntIP string, gatewayDomains []string, gatewayPrivateEndpointIP string) ([]byte, error) {
	t := template.Must(template.New(configFileName).Parse(`{{ .APIIntIP }}	api.{{ .ClusterDomain }} api-int.{{ .ClusterDomain }}
	{{ $.GatewayPrivateEndpointIP }}	{{- range $GatewayDomain := .GatewayDomains }}{{ $GatewayDomain }} {{- end }}`))
	buf := &bytes.Buffer{}

	err := t.ExecuteTemplate(buf, configFileName, &struct {
		ClusterDomain            string
		APIIntIP                 string
		GatewayDomains           []string
		GatewayPrivateEndpointIP string
	}{
		ClusterDomain:            clusterDomain,
		APIIntIP:                 apiIntIP,
		GatewayDomains:           gatewayDomains,
		GatewayPrivateEndpointIP: gatewayPrivateEndpointIP,
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func EtcHostsIgnition(clusterDomain string, apiIntIP string, gatewayDomains []string, gatewayPrivateEndpointIP string) (*ign3types.Config, error) {

	ign := &ign3types.Config{
		Ignition: ign3types.Ignition{
			// This Ignition Config version should be kept up to date with the default
			// rendered Ignition Config version from the Machine Config Operator version
			// on the lowest OCP version we support (4.7).
			Version: semver.Version{
				Major: 3,
				Minor: 2,
			}.String(),
		},
		Storage: ign3types.Storage{
			Files: []ign3types.File{
				{
					Node: ign3types.Node{
						Path: "/etc/hosts",
						User: ign3types.NodeUser{
							Name: to.StringPtr("root"),
						},
					},
					FileEmbedded1: ign3types.FileEmbedded1{
						Append: []ign3types.Resource{
							Source: to.StringPtr(dataurl.EncodeBytes(config)),
						},
						Mode: to.IntPtr(0644),
					},
				},
			}
		}
	}

	return ign, nil
}
