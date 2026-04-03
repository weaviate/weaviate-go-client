package auth_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/transport"
	"github.com/weaviate/weaviate-go-client/v6/internal/auth"
	"golang.org/x/oauth2"
)

//nolint:errcheck
func TestOIDC(t *testing.T) {
	// tokenJSON is the struct representing the HTTP response from OAuth2
	// providers returning a token or error in JSON form.
	// https://datatracker.ietf.org/doc/html/rfc6749#section-5.1
	//
	// Sorced from golang.org/x/oauth2/internal/token.go and trimmed.
	type tokenJSON struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}

	defaultScopes := []string{"profile"}

	for _, tt := range []struct {
		grant       string
		ex          transport.Exchanger
		verifyGrant func(*testing.T, url.Values)
		wantScopes  []string
		resp        tokenJSON
	}{
		{
			grant: "refresh_token",
			ex: auth.RefreshToken(oauth2.Token{
				AccessToken:  "expired-access-token",
				RefreshToken: "my-refresh",
				Expiry:       time.Now().Add(-time.Minute),
			}),
			verifyGrant: func(t *testing.T, v url.Values) {
				assert.Equal(t, "refresh_token", v.Get("grant_type"), "bad grant type")
				assert.Equal(t, "my-refresh", v.Get("refresh_token"), "bad grant")
			},
			resp: tokenJSON{
				AccessToken:  "fresh-access-token",
				RefreshToken: "fresh-refresh-token",
				ExpiresIn:    900,
			},
		},
		{
			grant: "refresh_token",
			ex: auth.RefreshToken(oauth2.Token{
				AccessToken:  "expired-access-token",
				RefreshToken: "my-refresh",
				ExpiresIn:    -60,
				// RefreshToken.Exchange should fill in the missing expiry
			}),
			verifyGrant: func(t *testing.T, v url.Values) {
				assert.Equal(t, "refresh_token", v.Get("grant_type"), "bad grant type")
				assert.Equal(t, "my-refresh", v.Get("refresh_token"), "bad grant")
			},
			resp: tokenJSON{
				AccessToken:  "fresh-access-token",
				RefreshToken: "fresh-refresh-token",
				ExpiresIn:    900,
			},
		},
		{
			grant: "client_credentials",
			ex: auth.ClientCredentials{
				ClientSecret: "my-secret",
				Scopes:       []string{"email"},
			},
			verifyGrant: func(t *testing.T, v url.Values) {
				assert.Equal(t, "client_credentials", v.Get("grant_type"), "bad grant type")
				assert.Equal(t, "my-client", v.Get("client_id"), "bad client_id")
				assert.Equal(t, "my-secret", v.Get("client_secret"), "bad client_secret")
			},
			wantScopes: append(defaultScopes, "offline_access", "email"),
			resp: tokenJSON{
				AccessToken: "fresh-access-token",
				ExpiresIn:   900,
			},
		},
		{
			grant: "password",
			ex: auth.ResourceOwnerPasswordCredentials{
				ClientSecret: "my-secret",
				Username:     "my-username",
				Password:     "my-password",
				Scopes:       []string{"email"},
			},
			verifyGrant: func(t *testing.T, v url.Values) {
				assert.Equal(t, "password", v.Get("grant_type"), "bad grant type")
				assert.Equal(t, "my-username", v.Get("username"), "bad username")
				assert.Equal(t, "my-password", v.Get("password"), "bad password")
			},
			wantScopes: append(defaultScopes, "offline_access", "email"),
			resp: tokenJSON{
				AccessToken:  "fresh-access-token",
				RefreshToken: "fresh-refresh-token",
				ExpiresIn:    900,
			},
		},
	} {
		t.Run(tt.grant, func(t *testing.T) {
			// Arrange
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer r.Body.Close()

				body, err := io.ReadAll(r.Body)
				require.NoError(t, err, "read request body")

				assert.Equal(t,
					"application/x-www-form-urlencoded",
					r.Header.Get("content-type"),
					"bad content-type",
				)

				v, err := url.ParseQuery(string(body))
				require.NoError(t, err, "parse query")

				tt.verifyGrant(t, v)

				var scopes []string
				if v.Has("scope") {
					scopes = strings.Split(v.Get("scope"), " ")
				}
				slices.Sort(scopes)
				slices.Sort(tt.wantScopes)
				assert.Equal(t, tt.wantScopes, scopes, "bad scopes")

				resp, err := json.Marshal(tt.resp)
				require.NoError(t, err, "marshal token response")

				w.Header().Set("content-type", "application/json")
				w.Write(resp)
			})

			srv := httptest.NewServer(handler)
			t.Cleanup(srv.Close)

			// Act
			ts, err := tt.ex.Exchange(t.Context(), oauth2.Config{
				ClientID: "my-client",
				Scopes:   defaultScopes,
				Endpoint: oauth2.Endpoint{TokenURL: srv.URL},
			})
			require.NoError(t, err, "exchange")
			require.NotNil(t, ts, "got nil token source")

			got, err := ts.Token()
			require.NoError(t, err, "get token")
			require.NotNil(t, got, "got nil token")

			// Assert
			assert.Equal(t, tt.resp.AccessToken, got.AccessToken, "bad access token")
			assert.Equal(t, tt.resp.RefreshToken, got.RefreshToken, "bad refresh token")
			assert.True(t, got.Valid(), "invalid token")
			if tt.resp.ExpiresIn > 0 {
				assert.NotZerof(t, got.Expiry, "token never expires, expected %ds", tt.resp.ExpiresIn)
			}
		})
	}
}
