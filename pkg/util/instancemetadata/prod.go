package instancemetadata

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/form3tech-oss/jwt-go"

	"github.com/Azure/ARO-RP/pkg/util/azureclaim"
	"github.com/Azure/ARO-RP/pkg/util/azureclient"
)

type prod struct {
	instanceMetadata

	do func(*http.Request) (*http.Response, error)
}

func newProd(ctx context.Context) (InstanceMetadata, error) {
	p := &prod{
		do: http.DefaultClient.Do,
	}

	err := p.populateInstanceMetadata()
	if err != nil {
		return nil, err
	}

	err = p.populateTenantAndClientIDFromMSI(ctx)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *prod) getServicePrincipalTokenAndClientIdFromMSI() (string, string, error) {
	shard, err := getAksShardFromEnvironment()
	if err != nil {
		return "", "", err
	}

	msi_endpoint, err := url.Parse("http://169.254.169.254/metadata/identity/oauth2/token")
	if err != nil {
		return "", "", err
	}

	msi_parameters := msi_endpoint.Query()
	msi_parameters.Add("api-version", "2018-02-01")
	msi_parameters.Add("resource", p.instanceMetadata.environment.ResourceManagerEndpoint)
	msi_parameters.Add("mi_res_id", fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ManagedIdentity/userAssignedIdentities/aro-aks-cluster-%03d-agentpool",
		p.SubscriptionID(),
		p.ResourceGroup(),
		shard,
	))

	msi_endpoint.RawQuery = msi_parameters.Encode()
	req, err := http.NewRequest("GET", msi_endpoint.String(), nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Add("Metadata", "true")

	resp, err := p.do(req)
	if err != nil {
		return "", "", err
	}

	responseBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", "", err
	}

	var responseJson *struct {
		AccessToken string `json:"access_token"`
		ClientID    string `json:"client_id"`
	}

	err = json.Unmarshal(responseBytes, &responseJson)
	if err != nil {
		return "", "", err
	}

	return responseJson.AccessToken, responseJson.ClientID, nil
}
func (p *prod) populateTenantAndClientIDFromMSI(ctx context.Context) error {
	accessToken, clientId, err := p.getServicePrincipalTokenAndClientIdFromMSI()
	if err != nil {
		return err
	}

	parser := &jwt.Parser{}
	c := &azureclaim.AzureClaim{}
	_, _, err = parser.ParseUnverified(accessToken, c)
	if err != nil {
		return err
	}

	p.tenantID = c.TenantID
	p.aksMsiClientID = clientId

	return nil
}

func (p *prod) populateInstanceMetadata() error {
	req, err := http.NewRequest(http.MethodGet, "http://169.254.169.254/metadata/instance/compute?api-version=2019-03-11", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Metadata", "true")

	resp, err := p.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	if strings.SplitN(resp.Header.Get("Content-Type"), ";", 2)[0] != "application/json" {
		return fmt.Errorf("unexpected content type %q", resp.Header.Get("Content-Type"))
	}

	var m *struct {
		Location          string `json:"location,omitempty"`
		ResourceGroupName string `json:"resourceGroupName,omitempty"`
		SubscriptionID    string `json:"subscriptionId,omitempty"`
		AzEnvironment     string `json:"azEnvironment,omitempty"`
	}

	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return err
	}

	environment, err := azureclient.EnvironmentFromName(m.AzEnvironment)
	if err != nil {
		return err
	}
	p.environment = &environment
	p.subscriptionID = m.SubscriptionID
	p.location = m.Location
	p.resourceGroup = m.ResourceGroupName

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	p.hostname = hostname

	return nil
}
