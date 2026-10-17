package azureclient

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

type identityTransport func(*http.Request) (*http.Response, error)

func (f identityTransport) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestWorkloadIdentityCredentialOptions(t *testing.T) {
	for _, environment := range []AROEnvironment{PublicCloud, USGovernmentCloud} {
		for _, scope := range []string{environment.KeyVaultScope, environment.ResourceManagerScope} {
			t.Run(environment.Name+"/"+scope, func(t *testing.T) {
				tokenFile := filepath.Join(t.TempDir(), "token")
				if err := os.WriteFile(tokenFile, []byte("projected-service-account-token"), 0600); err != nil {
					t.Fatal(err)
				}
				t.Setenv("AZURE_CLIENT_ID", "installer-hosting-client")
				t.Setenv("AZURE_TENANT_ID", "hosting-tenant")
				t.Setenv("AZURE_FEDERATED_TOKEN_FILE", tokenFile)

				options := environment.WorkloadIdentityCredentialOptions()
				authority := strings.TrimRight(environment.Cloud.ActiveDirectoryAuthorityHost, "/")
				if options.Cloud.ActiveDirectoryAuthorityHost != environment.Cloud.ActiveDirectoryAuthorityHost {
					t.Fatal("credential options do not use the deployment cloud")
				}
				// Instance discovery is unrelated to the wrapper's token-exchange contract.
				options.DisableInstanceDiscovery = true
				tokenRequests := 0
				options.Transport = identityTransport(func(req *http.Request) (*http.Response, error) {
					if !strings.HasPrefix(req.URL.String(), authority+"/hosting-tenant/") {
						return nil, fmt.Errorf("unexpected identity endpoint %s", req.URL.Host)
					}
					var body string
					switch {
					case req.Method == http.MethodGet && strings.HasSuffix(req.URL.Path, "/.well-known/openid-configuration"):
						body = fmt.Sprintf(`{
							"authorization_endpoint": "%[1]s/hosting-tenant/oauth2/v2.0/authorize",
							"token_endpoint": "%[1]s/hosting-tenant/oauth2/v2.0/token",
							"issuer": "%[1]s/hosting-tenant/v2.0"
						}`, authority)
					case req.Method == http.MethodPost && req.URL.Path == "/hosting-tenant/oauth2/v2.0/token":
						tokenRequests++
						b, err := io.ReadAll(req.Body)
						if err != nil {
							return nil, err
						}
						form, err := url.ParseQuery(string(b))
						if err != nil {
							return nil, err
						}
						for key, want := range map[string]string{
							"client_id":             "installer-hosting-client",
							"client_assertion":      "projected-service-account-token",
							"client_assertion_type": "urn:ietf:params:oauth:client-assertion-type:jwt-bearer",
							"grant_type":            "client_credentials",
						} {
							if got := form.Get(key); got != want {
								t.Errorf("unexpected token request field %s", key)
							}
						}
						if !slices.Contains(strings.Fields(form.Get("scope")), scope) {
							t.Error("token request is missing the requested Azure service scope")
						}
						if form.Has("client_secret") {
							t.Error("hosting identity must not use a client secret")
						}
						body = `{"access_token":"hosting-access-token","token_type":"Bearer","expires_in":3600}`
					default:
						return nil, fmt.Errorf("unexpected identity request %s %s", req.Method, req.URL.Path)
					}
					return &http.Response{
						StatusCode: http.StatusOK,
						Header:     http.Header{"Content-Type": []string{"application/json"}},
						Body:       io.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				})
				credential, err := azidentity.NewWorkloadIdentityCredential(options)
				if err != nil {
					t.Fatal(err)
				}
				token, err := credential.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{scope}})
				if err != nil {
					t.Fatal(err)
				}
				if token.Token != "hosting-access-token" || tokenRequests != 1 {
					t.Fatal("expected one successful projected-token exchange")
				}
			})
		}
	}
}
