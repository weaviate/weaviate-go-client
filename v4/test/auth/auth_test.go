package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/semi-technologies/weaviate-go-client/v4/test/testsuit"

	"github.com/semi-technologies/weaviate-go-client/v4/weaviate"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/auth"
	"github.com/stretchr/testify/assert"
)

const (
	OktaScope = "some_scope"
	WcsUser   = "ms_2d0e007e7136de11d5f29fce7a53dae219a51458@existiert.net"
	OktaUser  = "test@test.de"
)

func TestAuth_clientCredential(t *testing.T) {
	tests := []struct {
		name   string
		envVar string
		scope  []string
		port   int
	}{
		{name: "Okta", envVar: "OKTA_CLIENT_SECRET", scope: []string{OktaScope}, port: testsuit.OktaCCPort},
		{name: "Azure", envVar: "AZURE_CLIENT_SECRET", scope: []string{"4706508f-30c2-469b-8b12-ad272b3de864/.default"}, port: testsuit.AzurePort},
		{name: "Azure (hardcoded scope)", envVar: "AZURE_CLIENT_SECRET", scope: nil, port: testsuit.AzurePort},
	}

	for _, tc := range tests {
		clientSecret := os.Getenv(tc.envVar)
		if clientSecret == "" {
			t.Skip("No client secret supplied for ", tc.name)
		}

		clientCredentialConf := auth.ClientCredentials{ClientSecret: clientSecret, Scopes: tc.scope}
		cfg, err := weaviate.NewConfig("localhost:"+fmt.Sprint(tc.port), "http", clientCredentialConf, nil)
		assert.Nil(t, err)
		client := weaviate.New(*cfg)
		AuthErr := client.Schema().AllDeleter().Do(context.TODO())
		assert.Nil(t, AuthErr)
	}
}

func TestAuth_clientCredential_WrongParameters(t *testing.T) {
	clientSecret := os.Getenv("OKTA_CLIENT_SECRET")
	if clientSecret == "" {
		t.Skip("No client secret supplied for okta")
	}

	tests := []struct {
		name   string
		secret string
		scope  []string
	}{
		{name: "Wrong credential", secret: "ImNotaRealSecret", scope: []string{"OktaScope"}},
		{name: "Wrong scope", secret: clientSecret, scope: []string{"MadeUpScope"}},
	}

	for _, tc := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			clientCredentialConf := auth.ClientCredentials{ClientSecret: tc.secret, Scopes: tc.scope}
			cfg, err := weaviate.NewConfig("localhost:"+fmt.Sprint(testsuit.OktaCCPort), "http", clientCredentialConf, nil)
			assert.Nil(t, err)
			client := weaviate.New(*cfg)
			AuthErr := client.Schema().AllDeleter().Do(context.TODO())
			assert.NotNil(t, AuthErr)
		})
	}
}

func TestAuth_UserPW_WCS(t *testing.T) {
	tests := []struct {
		name    string
		user    string
		envVar  string
		scope   []string
		port    int
		warning bool
	}{
		{name: "WCS", user: WcsUser, envVar: "WCS_DUMMY_CI_PW", port: testsuit.WCSPort, warning: false},
		{name: "Okta (no scope)", user: OktaUser, envVar: "OKTA_DUMMY_CI_PW", port: testsuit.OktaUsersPort, warning: false},
		{name: "Okta", user: OktaUser, envVar: "OKTA_DUMMY_CI_PW", port: testsuit.OktaUsersPort, scope: []string{"offline_access"}, warning: false},
		{name: "Okta (scope without refresh)", user: OktaUser, envVar: "OKTA_DUMMY_CI_PW", port: testsuit.OktaUsersPort, scope: []string{"offline_access"}, warning: true},
	}
	for _, tc := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			// write log to buffer
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer func() {
				log.SetOutput(os.Stderr)
			}()

			pw := os.Getenv(tc.envVar)
			if pw == "" {
				t.Skip("No password supplied for " + tc.name)
			}

			clientCredentialConf := auth.ResourceOwnerPasswordFlow{Username: tc.user, Password: pw, Scopes: tc.scope}
			cfg, err := weaviate.NewConfig("localhost:"+fmt.Sprint(tc.port), "http", clientCredentialConf, nil)
			assert.Nil(t, err)
			client := weaviate.New(*cfg)
			AuthErr := client.Schema().AllDeleter().Do(context.TODO())
			assert.Nil(t, AuthErr)

			if tc.warning {
				assert.True(t, strings.Contains(buf.String(), "Auth001"))
			}
		})
	}
}

func TestAuth_UserPW_wrongPW(t *testing.T) {
	clientCredentialConf := auth.ResourceOwnerPasswordFlow{Username: "SomeUsername", Password: "IamWrong"}
	_, err := weaviate.NewConfig("localhost:"+fmt.Sprint(testsuit.WCSPort), "http", clientCredentialConf, nil)
	assert.NotNil(t, err)
}

func TestNoAuthOnWeaviateWithoutAuth(t *testing.T) {
	cfg, err := weaviate.NewConfig("localhost:"+fmt.Sprint(testsuit.NoAuthPort), "http", nil, nil)
	assert.Nil(t, err)
	client := weaviate.New(*cfg)

	AuthErr := client.Schema().AllDeleter().Do(context.TODO())
	assert.Nil(t, AuthErr)
}

func TestNoAuthOnWeaviateWithAuth(t *testing.T) {
	cfg, err := weaviate.NewConfig("localhost:"+fmt.Sprint(testsuit.WCSPort), "http", nil, nil)
	assert.Nil(t, err)
	client := weaviate.New(*cfg)

	AuthErr := client.Schema().AllDeleter().Do(context.TODO())
	assert.NotNil(t, AuthErr)
}

// Test that log contains a warning when configuring the client with authentication, but weaviate is configured without
// authentication. Otherwise, the client is working normally
func TestAuthOnWeaviateWithoutAuth(t *testing.T) {
	tests := []struct {
		name       string
		authConfig auth.Config
		scope      []string
	}{
		{name: "User/PW", authConfig: auth.ResourceOwnerPasswordFlow{Username: "SomeUsername", Password: "IamWrong"}},
		{name: "Client credentials", authConfig: auth.ClientCredentials{ClientSecret: "NotASecret", Scopes: []string{"No scope"}}},
		{name: "Bearer token", authConfig: auth.BearerToken{AccessToken: "NotAToken"}},
	}

	for _, tc := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer func() {
				log.SetOutput(os.Stderr)
			}()

			cfg, err := weaviate.NewConfig("localhost:"+fmt.Sprint(testsuit.NoAuthPort), "http", tc.authConfig, nil)
			assert.Nil(t, err)
			assert.True(t, strings.Contains(buf.String(), "The client was configured to use authentication"))

			client := weaviate.New(*cfg)
			AuthErr := client.Schema().AllDeleter().Do(context.TODO())
			assert.Nil(t, AuthErr)
		})
	}
}

func TestAuthNoWeaviateOnPort(t *testing.T) {
	_, err := weaviate.NewConfig("localhost:"+fmt.Sprint(testsuit.NoWeaviatePort), "http", auth.ResourceOwnerPasswordFlow{Username: "SomeUsername", Password: "IamWrong"}, nil)
	assert.NotNil(t, err)
}

func TestAuthBearerToken(t *testing.T) {
	tests := []struct {
		name   string
		user   string
		envVar string
		port   int
	}{
		{name: "WCS", user: WcsUser, envVar: "WCS_DUMMY_CI_PW", port: testsuit.WCSPort},
		{name: "Okta", user: OktaUser, envVar: "OKTA_DUMMY_CI_PW", port: testsuit.OktaUsersPort},
	}
	for _, tc := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			pw := os.Getenv(tc.envVar)
			if pw == "" {
				t.Skip("No password supplied for " + tc.name)
			}
			url := "localhost:" + fmt.Sprint(tc.port)

			AccessToken, RefreshToken := get_access_token(t, url, tc.user, pw)
			cfg, err := weaviate.NewConfig(url, "http", auth.BearerToken{AccessToken: AccessToken, RefreshToken: RefreshToken}, nil)
			assert.Nil(t, err)

			client := weaviate.New(*cfg)
			AuthErr := client.Schema().AllDeleter().Do(context.TODO())
			assert.Nil(t, AuthErr)
		})
	}
}

func get_access_token(t *testing.T, weavUrl, user, pw string) (string, string) {
	resp, err := http.Get("http://" + weavUrl + "/v1/.well-known/openid-configuration")
	if err != nil {
		t.Fail()
	}
	body, _ := io.ReadAll(resp.Body)
	cfg := struct {
		Href     string `json:"href"`
		ClientID string `json:"clientId"`
	}{}
	json.Unmarshal(body, &cfg)
	resp.Body.Close()
	respAuth, err := http.Get(cfg.Href)
	bodyAuth, _ := io.ReadAll(respAuth.Body)
	endpoint := struct {
		TokenEndpoint string `json:"token_endpoint"`
	}{}
	json.Unmarshal(bodyAuth, &endpoint)
	respAuth.Body.Close()
	respToken, _ := http.PostForm(endpoint.TokenEndpoint, url.Values{
		"grant_type": []string{"password"}, "client_id": []string{cfg.ClientID}, "username": []string{user}, "password": []string{pw},
	})
	bodyTokens, _ := io.ReadAll(respToken.Body)
	respToken.Body.Close()

	tokens := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{}
	json.Unmarshal(bodyTokens, &tokens)
	return tokens.AccessToken, tokens.RefreshToken
}
