package test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/auth"
)

const (
	AccessToken  = "HELLO!IamAnAccessToken"
	RefreshToken = "IAmARefreshToken"
)

// Test that the client warns when no refresh token is provided by the authentication provider
func TestAuthMock_NoRefreshToken(t *testing.T) {
	tests := []struct {
		name       string
		authConfig auth.Config
		scope      []string
	}{
		{name: "User/PW", authConfig: auth.ResourceOwnerPasswordFlow{Username: "SomeUsername", Password: "IamWrong"}},
		{name: "Bearer token", authConfig: auth.BearerToken{AccessToken: "NotAToken"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// write log to buffer
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer func() {
				log.SetOutput(os.Stderr)
			}()

			// endpoint for access tokens
			muxToken := http.NewServeMux()
			muxToken.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(fmt.Sprint(`{"access_token": "` + AccessToken + `", "expires_in": "5"}`)))
			})
			sToken := httptest.NewServer(muxToken)
			defer sToken.Close()

			// provides all endpoints
			muxEndpoints := http.NewServeMux()
			muxEndpoints.HandleFunc("/endpoints", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(fmt.Sprint(`{"token_endpoint": "` + sToken.URL + `/auth"}`)))
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
			mux.HandleFunc("/v1/.well-known/ready", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`{}`))
			})
			s := httptest.NewServer(mux)
			defer s.Close()
			cfg := weaviate.Config{Host: strings.TrimPrefix(s.URL, "http://"), Scheme: "http", StartupTimeout: 60 * time.Second, AuthConfig: tc.authConfig}

			client, err := weaviate.NewClient(cfg)
			require.Nil(t, err)
			assert.True(t, strings.Contains(buf.String(), "Auth002"))

			AuthErr := client.Schema().AllDeleter().Do(context.TODO())
			assert.Nil(t, AuthErr)
		})
	}
}

// Test that client using CC automatically get a new token after expiration
func TestAuthMock_RefreshCC(t *testing.T) {
	i := 0
	// endpoint for access tokens
	muxToken := http.NewServeMux()
	muxToken.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		i += 1 // record how often this was called
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprint(`{"access_token": "` + AccessToken + `", "expires_in": "1"}`)))
	})
	sToken := httptest.NewServer(muxToken)
	defer sToken.Close()

	// provides all endpoints
	muxEndpoints := http.NewServeMux()
	muxEndpoints.HandleFunc("/endpoints", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprint(`{"token_endpoint": "` + sToken.URL + `/auth"}`)))
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
	mux.HandleFunc("/v1/.well-known/ready", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{}`))
	})
	s := httptest.NewServer(mux)
	defer s.Close()
	cfg := weaviate.Config{Host: strings.TrimPrefix(s.URL, "http://"), Scheme: "http", StartupTimeout: 60 * time.Second, AuthConfig: auth.ClientCredentials{ClientSecret: "SecretValue"}}
	client, err := weaviate.NewClient(cfg)
	assert.Nil(t, err)

	AuthErr := client.Schema().AllDeleter().Do(context.TODO())
	assert.Nil(t, AuthErr)
	assert.Equal(t, i, 4) // client does 4 initial calls to token endpoint

	time.Sleep(time.Second * 5)
	// current token expires, so the oauth client needs to get a new one
	AuthErr2 := client.Schema().AllDeleter().Do(context.TODO())
	assert.Equal(t, i, 5)
	assert.Nil(t, AuthErr2)
}

// Test that client uses refresh tokens to get new access/refresh tokens before their expiration, including during idle
// times.
func TestAuthMock_RefreshUserPWAndToken(t *testing.T) {
	expirationTimeRefreshToken := 3
	initialExpirationTimeToken := uint(2)
	// higher expiration time to check if the go routine correctly sleeps until the token almost expires
	expirationTimeToken := uint(12)
	tests := []struct {
		name       string
		authConfig auth.Config
		scope      []string
	}{
		{name: "User/PW", authConfig: auth.ResourceOwnerPasswordFlow{Username: "SomeUsername", Password: "IamWrong"}},
		{name: "Bearer token", authConfig: auth.BearerToken{
			AccessToken: AccessToken, ExpiresIn: initialExpirationTimeToken, RefreshToken: RefreshToken,
		}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tokenRefreshTime := time.Now()
			// endpoint for access tokens
			muxToken := http.NewServeMux()
			muxToken.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
				// refresh token cannot be expired
				assert.True(t, time.Since(tokenRefreshTime).Seconds() < float64(expirationTimeRefreshToken))

				tokenRefreshTime = time.Now() // update time when the tokens where refreshed the last time
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(
					fmt.Sprintf(`{"access_token": "%v", "expires_in": %v, "refresh_token": "%v", "refresh_expires_in" :  %v}`,
						AccessToken, expirationTimeToken, RefreshToken, expirationTimeRefreshToken)))
			})
			sToken := httptest.NewServer(muxToken)
			defer sToken.Close()

			// provides all endpoints
			muxEndpoints := http.NewServeMux()
			muxEndpoints.HandleFunc("/endpoints", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(fmt.Sprint(`{"token_endpoint": "` + sToken.URL + `/auth"}`)))
			})
			sEndpoints := httptest.NewServer(muxEndpoints)
			defer sEndpoints.Close()

			// Returns the address of the auth server
			mux := http.NewServeMux()
			mux.HandleFunc("/v1/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`{"href": "` + sEndpoints.URL + `/endpoints", "clientId": "DoesNotMatter"}`))
			})
			mux.HandleFunc("/v1/schema", func(w http.ResponseWriter, r *http.Request) {
				// Access Token cannot be expired
				assert.True(t, time.Since(tokenRefreshTime).Seconds() < float64(expirationTimeToken))
				w.Write([]byte(`{}`))
			})
			mux.HandleFunc("/v1/.well-known/ready", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`{}`))
			})
			s := httptest.NewServer(mux)
			defer s.Close()

			cfg := weaviate.Config{Host: strings.TrimPrefix(s.URL, "http://"), Scheme: "http", StartupTimeout: 60 * time.Second, AuthConfig: tc.authConfig}
			client, err := weaviate.NewClient(cfg)
			assert.Nil(t, err)

			AuthErr := client.Schema().AllDeleter().Do(context.TODO())
			assert.Nil(t, AuthErr)

			// access and refresh token expired, so the client needs to refresh automatically in the background
			time.Sleep(time.Second * 5)
			AuthErr2 := client.Schema().AllDeleter().Do(context.TODO())
			assert.Nil(t, AuthErr2)
		})
	}
}

// Test that the client can handle situations in which a proxy returns a catchall page for all requests
func TestAuthMock_CatchAllProxy(t *testing.T) {
	// write log to buffer
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	// Simulate a proxy that returns something if a page is not available => no valid json
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`NotAValidJsonResponse`))
	})
	mux.HandleFunc("/v1/schema", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{}`))
	})
	mux.HandleFunc("/v1/.well-known/ready", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{}`))
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	cfg := weaviate.Config{Host: strings.TrimPrefix(s.URL, "http://"), Scheme: "http", StartupTimeout: 60 * time.Second, AuthConfig: nil}
	client, err := weaviate.NewClient(cfg)
	assert.Nil(t, err)

	AuthErr := client.Schema().AllDeleter().Do(context.TODO())
	assert.Nil(t, AuthErr)
}

// Test that client using CC automatically get a new token after expiration
func TestAuthMock_CheckDefaultScopes(t *testing.T) {
	// endpoint for access tokens
	muxToken := http.NewServeMux()
	muxToken.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, _ := io.ReadAll(r.Body)
		bodyS := string(body)
		assert.Equal(t, bodyS[len(bodyS)-15:], "something+extra") // scopes are in the body

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprint(`{"access_token": "` + AccessToken + `", "expires_in": "1"}`)))
	})
	sToken := httptest.NewServer(muxToken)
	defer sToken.Close()

	// provides all endpoints
	muxEndpoints := http.NewServeMux()
	muxEndpoints.HandleFunc("/endpoints", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprint(`{"token_endpoint": "` + sToken.URL + `/auth"}`)))
	})
	sEndpoints := httptest.NewServer(muxEndpoints)
	defer sEndpoints.Close()

	// Returns the address of the auth server
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"href": "` + sEndpoints.URL + `/endpoints", "clientId": "DoesNotMatter", "scopes": ["something", "extra"]}`))
	})
	mux.HandleFunc("/v1/schema", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{}`))
	})
	mux.HandleFunc("/v1/.well-known/ready", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{}`))
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	cfg := weaviate.Config{Host: strings.TrimPrefix(s.URL, "http://"), Scheme: "http", StartupTimeout: 60 * time.Second, AuthConfig: auth.ClientCredentials{ClientSecret: "SecretValue"}}
	client, err := weaviate.NewClient(cfg)
	assert.Nil(t, err)

	AuthErr := client.Schema().AllDeleter().Do(context.TODO())
	assert.Nil(t, AuthErr)
}

// Test that the correct header is set when using an API Key to authenticate.
func TestAuthMock_WithSimpleAuthNoOidcViaApiKey(t *testing.T) {
	token := "super-secret-key"
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/schema", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		require.Equal(t, authHeader, "Bearer "+token)
		w.Write([]byte(`{}`))
	})
	mux.HandleFunc("/v1/.well-known/ready", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{}`))
	})
	s := httptest.NewServer(mux)
	defer s.Close()
	url := strings.TrimPrefix(s.URL, "http://")
	authConf := auth.ApiKey{Value: token}

	t.Run("NewClient", func(t *testing.T) {
		cfg := weaviate.Config{Host: strings.TrimPrefix(s.URL, "http://"), Scheme: "http", StartupTimeout: 60 * time.Second, AuthConfig: authConf}
		client, err := weaviate.NewClient(cfg)
		assert.Nil(t, err)

		AuthErr := client.Schema().AllDeleter().Do(context.TODO())
		assert.Nil(t, AuthErr)
	})

	t.Run("deprecated new Config", func(t *testing.T) {
		cfg, err := weaviate.NewConfig(url, "http", authConf, nil)
		assert.Nil(t, err)
		client := weaviate.New(*cfg)
		AuthErr := client.Schema().AllDeleter().Do(context.TODO())
		assert.Nil(t, AuthErr)
	})
}
