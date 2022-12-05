package test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/v4/test/testsuit"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/semi-technologies/weaviate-go-client/v4/weaviate"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/auth"
	"github.com/stretchr/testify/assert"
)

const OktaScope = "some_scope"

func TestAuth_clientCredential(t *testing.T) {
	tests := []struct {
		name   string
		envVar string
		scope  []string
		port   int
	}{
		{name: "Okta", envVar: "OKTA_CLIENT_SECRET", scope: []string{OktaScope}, port: testsuit.OktaPort},
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
			cfg, err := weaviate.NewConfig("localhost:"+fmt.Sprint(testsuit.OktaPort), "http", clientCredentialConf, nil)
			assert.Nil(t, err)
			client := weaviate.New(*cfg)
			AuthErr := client.Schema().AllDeleter().Do(context.TODO())
			assert.NotNil(t, AuthErr)
		})
	}
}

func TestAuth_UserPW_WCS(t *testing.T) {
	wcsPw := os.Getenv("WCS_DUMMY_CI_PW")
	if wcsPw == "" {
		t.Skip("No password supplied for WCS")
	}

	clientCredentialConf := auth.ResourceOwnerPasswordFlow{Username: "ms_2d0e007e7136de11d5f29fce7a53dae219a51458@existiert.net", Password: wcsPw}
	cfg, err := weaviate.NewConfig("localhost:"+fmt.Sprint(testsuit.WCSPort), "http", clientCredentialConf, nil)
	assert.Nil(t, err)
	client := weaviate.New(*cfg)
	AuthErr := client.Schema().AllDeleter().Do(context.TODO())
	assert.Nil(t, AuthErr)
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
// authentication. Otherwise the client is working normally
func TestAuthOnWeaviateWithoutAuth(t *testing.T) {
	tests := []struct {
		name       string
		authConfig auth.Config
		scope      []string
	}{
		{name: "User/PW", authConfig: auth.ResourceOwnerPasswordFlow{Username: "SomeUsername", Password: "IamWrong"}},
		{name: "Client credentials", authConfig: auth.ClientCredentials{ClientSecret: "NotASecret", Scopes: []string{"No scope"}}},
		{name: "Bearer token", authConfig: auth.BearerToken{Token: "NotAToken"}},
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

func TestUserPWNoRefreshToken(t *testing.T) {
	// write log to buffer
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	AccessToken := "HELLO!IamAnAccessToken"

	// endpoint for access tokens
	muxToken := http.NewServeMux()
	muxToken.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprint(`{"access_token": "` + AccessToken + `", "expires_in": "1"}`)))
	})
	sToken := httptest.NewServer(muxToken)
	defer sToken.Close()

	// provides all endpoints
	muxEndpoints := http.NewServeMux()
	muxEndpoints.HandleFunc("/endpoints", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf(`{"token_endpoint": "` + sToken.URL + `/auth"}`)))
	})
	sEndpoints := httptest.NewServer(muxEndpoints)
	defer sEndpoints.Close()

	// Returns the address of the auth server
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"href": "` + sEndpoints.URL + `/endpoints", "clientId": "DoesNotMatter"}`))
	})
	mux.HandleFunc("/v1/schema", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{}`))
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	cfg, err := weaviate.NewConfig(strings.TrimPrefix(s.URL, "http://"), "http", auth.ResourceOwnerPasswordFlow{Username: "SomeUsername", Password: "IamWrong"}, nil)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(buf.String(), "Your access token is valid for"))

	client := weaviate.New(*cfg)
	AuthErr := client.Schema().AllDeleter().Do(context.TODO())
	assert.Nil(t, AuthErr)
}
