package dnsmasq

import (
	"bytes"
	"text/template"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/vincent-petithory/dataurl"
)

const restartScriptFileName = "99-dnsmasq-restart"

func nmDispatcherRestartDnsmasq() ([]byte, error) {
	t := template.Must(template.New(restartScriptFileName).Parse(restartScript))
	buf := &bytes.Buffer{}

	err := t.ExecuteTemplate(buf, restartScriptFileName, nil)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func restartScriptIgnFile(data []byte) types.File {
	return types.File{
		Node: types.Node{
			Overwrite: to.BoolPtr(true),
			Path:      "/etc/NetworkManager/dispatcher.d/" + restartScriptFileName,
			User: types.NodeUser{
				Name: to.StringPtr("root"),
			},
		},
		FileEmbedded1: types.FileEmbedded1{
			Contents: types.Resource{
				Source: to.StringPtr(dataurl.EncodeBytes(data)),
			},
			Mode: to.IntPtr(0744),
		},
	}
}
