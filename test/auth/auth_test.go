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
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/auth"
)

const (
	oktaScope = "some_scope"
	wcsUser   = "oidc-test-user@weaviate.io"
	oktaUser  = "test@test.de"
)

func TestAuth_clientCredential(t *testing.T) {
	t.Skip("Skipping bc client secret keys are expired")
	tests := []struct {
		name   string
		envVar string
		scope  []string
		port   int
	}{
		{name: "Okta", envVar: "OKTA_CLIENT_SECRET", scope: []string{oktaScope}, port: testsuit.OktaCCPort},
		{name: "Azure", envVar: "AZURE_CLIENT_SECRET", scope: []string{"4706508f-30c2-469b-8b12-ad272b3de864/.default"}, port: testsuit.AzurePort},
		{name: "Azure (hardcoded scope)", envVar: "AZURE_CLIENT_SECRET", scope: nil, port: testsuit.AzurePort},
	}

	for _, tc := range tests {
		clientSecret := os.Getenv(tc.envVar)
		if clientSecret == "" {
			t.Skip("No client secret supplied for ", tc.name)
		}

		clientCredentialConf := auth.ClientCredentials{ClientSecret: clientSecret, Scopes: tc.scope}
		cfg := weaviate.Config{Host: fmt.Sprintf("localhost:%v", tc.port), Scheme: "http", StartupTimeout: 60 * time.Second, AuthConfig: clientCredentialConf}
		client, err := weaviate.NewClient(cfg)
		assert.Nil(t, err)

		authErr := client.Schema().AllDeleter().Do(context.TODO())
		assert.Nil(t, authErr)
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
		t.Run(tc.name, func(t *testing.T) {
			clientCredentialConf := auth.ClientCredentials{ClientSecret: tc.secret, Scopes: tc.scope}

			cfg := weaviate.Config{Host: fmt.Sprintf("localhost:%v", testsuit.OktaCCPort), Scheme: "http", StartupTimeout: 20 * time.Second, AuthConfig: clientCredentialConf}
			client, err := weaviate.NewClient(cfg)
			assert.NotNil(t, err)
			assert.Nil(t, client)
		})
	}
}

func TestAuth_UserPW(t *testing.T) {
	tests := []struct {
		name    string
		user    string
		envVar  string
		scope   []string
		port    int
		warning bool
	}{
		{name: "WCS", user: wcsUser, envVar: "WCS_DUMMY_CI_PW", port: testsuit.WCSPort, warning: false},
		{name: "Okta (no scope)", user: oktaUser, envVar: "OKTA_DUMMY_CI_PW", port: testsuit.OktaUsersPort, warning: false},
		{name: "Okta", user: oktaUser, envVar: "OKTA_DUMMY_CI_PW", port: testsuit.OktaUsersPort, scope: []string{"offline_access"}, warning: false},
		{name: "Okta (scope without refresh)", user: oktaUser, envVar: "OKTA_DUMMY_CI_PW", port: testsuit.OktaUsersPort, scope: []string{"offline_access"}, warning: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// write log to buffer
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer func() {
				log.SetOutput(os.Stderr)
			}()

			pw := os.Getenv(tc.envVar)
			if pw == "" {
				t.Skip("No password supplied for " + tc.name)
			} else {
				// This should be in a branch, so the GC can collect the client and with that shut down the background
				// routine that writes to the log. Otherwise, we'd have a data race between the this goroutine and the
				// test accessing the buffer.
				clientCredentialConf := auth.ResourceOwnerPasswordFlow{Username: tc.user, Password: pw, Scopes: tc.scope}

				cfg := weaviate.Config{Host: fmt.Sprintf("localhost:%v", tc.port), Scheme: "http", StartupTimeout: 60 * time.Second, AuthConfig: clientCredentialConf}
				client, err := weaviate.NewClient(cfg)
				assert.Nil(t, err)

				authErr := client.Schema().AllDeleter().Do(context.TODO())
				assert.Nil(t, authErr)
			}
			runtime.GC()

			if tc.warning {
				assert.True(t, strings.Contains(buf.String(), "Auth002"))
			}
		})
	}
}

func TestAuth_UserPW_wrongPW(t *testing.T) {
	clientCredentialConf := auth.ResourceOwnerPasswordFlow{Username: "SomeUsername", Password: "IamWrong"}
	cfg := weaviate.Config{Host: fmt.Sprintf("localhost:%v", testsuit.WCSPort), Scheme: "http", StartupTimeout: 60 * time.Second, AuthConfig: clientCredentialConf}
	_, err := weaviate.NewClient(cfg)
	assert.NotNil(t, err)
}

func TestNoAuthOnWeaviateWithoutAuth(t *testing.T) {
	cfg := weaviate.Config{Host: fmt.Sprintf("localhost:%v", testsuit.NoAuthPort), Scheme: "http", StartupTimeout: 60 * time.Second}
	client, err := weaviate.NewClient(cfg)
	assert.Nil(t, err)

	authErr := client.Schema().AllDeleter().Do(context.TODO())
	assert.Nil(t, authErr)
}

func TestNoAuthOnWeaviateWithAuth(t *testing.T) {
	cfg := weaviate.Config{Host: fmt.Sprintf("localhost:%v", testsuit.WCSPort), Scheme: "http", StartupTimeout: 60 * time.Second}
	client, err := weaviate.NewClient(cfg)
	assert.Nil(t, err)

	authErr := client.Schema().AllDeleter().Do(context.TODO())
	assert.NotNil(t, authErr)
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
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer func() {
				log.SetOutput(os.Stderr)
			}()
			cfg := weaviate.Config{Host: fmt.Sprintf("localhost:%v", testsuit.NoAuthPort), Scheme: "http", StartupTimeout: 60 * time.Second, AuthConfig: tc.authConfig}
			client, err := weaviate.NewClient(cfg)
			assert.Nil(t, err)
			assert.True(t, strings.Contains(buf.String(), "The client was configured to use authentication"))

			authErr := client.Schema().AllDeleter().Do(context.TODO())
			assert.Nil(t, authErr)
		})
	}
}

func TestAuthNoWeaviateOnPort(t *testing.T) {
	cfg := weaviate.Config{Host: "localhost:" + fmt.Sprint(testsuit.NoWeaviatePort), Scheme: "http", StartupTimeout: 0 * time.Second, AuthConfig: auth.ResourceOwnerPasswordFlow{Username: "SomeUsername", Password: "IamWrong"}}
	_, err := weaviate.NewClient(cfg)
	assert.NotNil(t, err)
}

func TestAuthBearerToken(t *testing.T) {
	tests := []struct {
		name   string
		user   string
		envVar string
		port   int
	}{
		{name: "WCS", user: wcsUser, envVar: "WCS_DUMMY_CI_PW", port: testsuit.WCSPort},
		{name: "Okta", user: oktaUser, envVar: "OKTA_DUMMY_CI_PW", port: testsuit.OktaUsersPort},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pw := os.Getenv(tc.envVar)
			if pw == "" {
				t.Skip("No password supplied for " + tc.name)
			}
			url := fmt.Sprintf("localhost:%v", tc.port)

			accessToken, refreshToken := getAccessToken(t, url, tc.user, pw)
			cfg := weaviate.Config{Host: url, Scheme: "http", StartupTimeout: 60 * time.Second, AuthConfig: auth.BearerToken{AccessToken: accessToken, RefreshToken: refreshToken}}
			client, err := weaviate.NewClient(cfg)
			assert.Nil(t, err)

			authErr := client.Schema().AllDeleter().Do(context.TODO())
			assert.Nil(t, authErr)
		})
	}
}

func getAccessToken(t *testing.T, weaviateUrl, user, pw string) (string, string) {
	resp, err := http.Get(fmt.Sprintf("http://%s/v1/.well-known/openid-configuration", weaviateUrl))
	require.Nil(t, err)
	body, err := io.ReadAll(resp.Body)
	require.Nil(t, err)
	cfg := struct {
		Href     string `json:"href"`
		ClientID string `json:"clientId"`
	}{}
	err = json.Unmarshal(body, &cfg)
	require.Nil(t, err)
	if err := resp.Body.Close(); err != nil {
		t.Error(err)
	}
	respAuth, err := http.Get(cfg.Href)
	require.Nil(t, err)
	bodyAuth, err := io.ReadAll(respAuth.Body)
	require.Nil(t, err)
	endpoint := struct {
		TokenEndpoint string `json:"token_endpoint"`
	}{}
	err = json.Unmarshal(bodyAuth, &endpoint)
	require.Nil(t, err)
	err = respAuth.Body.Close()
	require.Nil(t, err)
	respToken, err := http.PostForm(endpoint.TokenEndpoint, url.Values{
		"grant_type": []string{"password"}, "client_id": []string{cfg.ClientID}, "username": []string{user}, "password": []string{pw},
	})
	require.Nil(t, err)
	bodyTokens, err := io.ReadAll(respToken.Body)
	require.Nil(t, err)
	err = respToken.Body.Close()
	require.Nil(t, err)
	// get tokens
	tokens := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{}
	err = json.Unmarshal(bodyTokens, &tokens)
	require.Nil(t, err)
	return tokens.AccessToken, tokens.RefreshToken
}

func TestAuth_ApiKey(t *testing.T) {
	tests := []struct {
		headers map[string]string
		newConf bool
	}{
		{headers: nil, newConf: false},
		{headers: map[string]string{}, newConf: false},
		{headers: map[string]string{"Key": "Val"}, newConf: false},
		{headers: nil, newConf: true},
		{headers: map[string]string{}, newConf: true},
		{headers: map[string]string{"Key": "Val"}, newConf: true},
	}
	for _, tc := range tests {
		t.Run("headers", func(t *testing.T) {
			url := fmt.Sprintf("127.0.0.1:%v", testsuit.WCSPort)
			authConf := auth.ApiKey{Value: "my-secret-key"}
			var client *weaviate.Client

			if tc.newConf {
				cfg, err := weaviate.NewConfig(url, "http", authConf, nil)
				assert.Nil(t, err)
				client = weaviate.New(*cfg)

			} else {
				cfg := weaviate.Config{Host: url, Scheme: "http", StartupTimeout: 60 * time.Second, AuthConfig: authConf, Headers: tc.headers}
				var err error
				client, err = weaviate.NewClient(cfg)
				assert.Nil(t, err)
			}

			authErr := client.Schema().AllDeleter().Do(context.TODO())
			assert.Nil(t, authErr)
		})
	}
}

func TestAuth_WrongApiKey(t *testing.T) {
	url := fmt.Sprintf("127.0.0.1:%v", testsuit.WCSPort)
	tests := []struct {
		headers map[string]string
	}{
		{headers: nil}, {headers: map[string]string{}}, {headers: map[string]string{"Key": "Val"}},
	}
	for _, tc := range tests {
		t.Run("headers", func(t *testing.T) {
			cfg := weaviate.Config{Host: url, Scheme: "http", StartupTimeout: 60 * time.Second, AuthConfig: auth.ApiKey{Value: "wrong_key"}, Headers: tc.headers}
			client, err := weaviate.NewClient(cfg)
			assert.Nil(t, err)

			authErr := client.Schema().AllDeleter().Do(context.TODO())
			assert.NotNil(t, authErr)
		})
	}
}

func TestAuth_ConfigCheck(t *testing.T) {
	tests := []struct {
		port   int
		exists bool
	}{
		{port: testsuit.WCSPort, exists: true}, {port: testsuit.NoAuthPort, exists: false},
	}
	for _, tc := range tests {
		t.Run("oidc config", func(t *testing.T) {
			url := fmt.Sprintf("127.0.0.1:%v", tc.port)

			cfg := weaviate.Config{Host: url, Scheme: "http"}
			client := weaviate.New(cfg)
			openIDconfig, err := client.Misc().OpenIDConfigurationGetter().Do(context.Background())
			assert.Nil(t, err)

			if tc.exists {
				assert.NotNil(t, openIDconfig)
			} else {
				assert.Nil(t, openIDconfig)
			}
		})
	}
}
