package weaviate_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6"
	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/transport"
	"github.com/weaviate/weaviate-go-client/v6/internal/auth"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"golang.org/x/oauth2"
)

// DO NOT enable t.Parallel() for this test as it modifies the global state.
func TestNewLocal(t *testing.T) {
	newFunc := transport.New
	t.Cleanup(func() { transport.New = newFunc })

	t.Run("default", func(t *testing.T) {
		var got transport.Config
		transport.New = func(_ context.Context, cfg transport.Config) (internal.Transport, error) {
			got = cfg
			return testkit.NopTransport, nil
		}

		c, err := weaviate.NewLocal(t.Context())
		assert.NoError(t, err)
		assert.NotNil(t, c, "nil client")

		assert.Equal(t, transport.Config{
			Scheme:   "http",
			RESTHost: "localhost",
			GRPCHost: "localhost",
			RESTPort: 8080,
			GRPCPort: 50051,
			Header: http.Header{
				"X-Weaviate-Client": {"weaviate-client-go" + "/" + weaviate.Version()},
			},
			Timeout: transport.Timeout{
				Read:  30 * time.Second,
				Write: 90 * time.Second,
			},
			Version: api.Version,
		}, got)
	})

	t.Run("with options", func(t *testing.T) {
		var got transport.Config
		transport.New = func(_ context.Context, cfg transport.Config) (internal.Transport, error) {
			got = cfg
			return testkit.NopTransport, nil
		}

		tokenSource := oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken:  "access_token",
			RefreshToken: "refresh_token",
		})

		c, err := weaviate.NewLocal(t.Context(),
			weaviate.WithScheme("https"),
			weaviate.WithHTTPPort(7070),
			weaviate.WithGRPCPort(54321),
			weaviate.WithHeader(http.Header{
				"X-Test": {"heads", "up"},
			}),
			weaviate.WithReadTimeout(20*time.Second),
			weaviate.WithBatchTimeout(100*time.Millisecond),
			weaviate.WithTokenSource(tokenSource),
		)
		assert.NoError(t, err)
		assert.NotNil(t, c, "nil client")

		assert.Equal(t, transport.Config{
			// Defaults
			RESTHost: "localhost",
			GRPCHost: "localhost",

			// Custom
			Scheme:   "https",
			RESTPort: 7070,
			GRPCPort: 54321,
			Header: http.Header{
				"X-Test":            {"heads", "up"},
				"X-Weaviate-Client": {"weaviate-client-go" + "/" + weaviate.Version()},
			},
			Auth: tokenSource,
			Timeout: transport.Timeout{
				Read:  20 * time.Second,
				Write: 90 * time.Second,
				Batch: 100 * time.Millisecond,
			},
			Version: api.Version,
		}, got)
	})
}

// DO NOT enable t.Parallel() for this test as it modifies the global state.
func TestNewWeaviateCloud(t *testing.T) {
	newFunc := transport.New
	t.Cleanup(func() { transport.New = newFunc })

	t.Run("default", func(t *testing.T) {
		var got transport.Config
		transport.New = func(_ context.Context, cfg transport.Config) (internal.Transport, error) {
			got = cfg
			return testkit.NopTransport, nil
		}

		c, err := weaviate.NewWeaviateCloud(t.Context(), "example.com", "api-key")
		assert.NoError(t, err)
		assert.NotNil(t, c, "nil client")

		assert.Equal(t, transport.Config{
			Scheme:   "https",
			RESTHost: "example.com",
			GRPCHost: "grpc-example.com",
			RESTPort: 443,
			GRPCPort: 443,
			Header: http.Header{
				"X-Weaviate-Client": {"weaviate-client-go" + "/" + weaviate.Version()},
			},
			Auth: oauth2.StaticTokenSource(&oauth2.Token{
				AccessToken: "api-key",
			}),
			Timeout: transport.Timeout{
				Read:  30 * time.Second,
				Write: 90 * time.Second,
			},
			Version: api.Version,
		}, got)
	})

	t.Run("weaviate domain", func(t *testing.T) {
		for _, domain := range []string{
			"weaviate.io",
			"weaviate.cloud",
			"semi.technology",
		} {
			t.Run(domain, func(t *testing.T) {
				var got transport.Config
				transport.New = func(_ context.Context, cfg transport.Config) (internal.Transport, error) {
					got = cfg
					return testkit.NopTransport, nil
				}

				c, err := weaviate.NewWeaviateCloud(t.Context(), "my."+domain, "api-key")
				assert.NoError(t, err)
				assert.NotNil(t, c, "nil client")

				assert.Equal(t, got.Header, http.Header{
					"X-Weaviate-Client":      {"weaviate-client-go" + "/" + weaviate.Version()},
					"X-Weaviate-Cluster-Url": {"https://my." + domain + ":443"},
				})
			})
		}
	})

	t.Run("with options", func(t *testing.T) {
		var got transport.Config
		transport.New = func(_ context.Context, cfg transport.Config) (internal.Transport, error) {
			got = cfg
			return testkit.NopTransport, nil
		}

		c, err := weaviate.NewWeaviateCloud(t.Context(), "example.com", "api-key",
			weaviate.WithHTTPPort(7070),
			weaviate.WithGRPCPort(54321),
			weaviate.WithHeader(http.Header{
				"X-Test": {"heads", "up"},
			}),
			weaviate.WithReadTimeout(20*time.Second),
			weaviate.WithBatchTimeout(100*time.Millisecond),
		)
		assert.NoError(t, err)
		assert.NotNil(t, c, "nil client")

		assert.Equal(t, transport.Config{
			Scheme:   "https",
			RESTHost: "example.com",
			GRPCHost: "grpc-example.com",
			RESTPort: 7070,
			GRPCPort: 54321,
			Header: http.Header{
				"X-Test":            {"heads", "up"},
				"X-Weaviate-Client": {"weaviate-client-go" + "/" + weaviate.Version()},
			},
			Auth: oauth2.StaticTokenSource(&oauth2.Token{
				AccessToken: "api-key",
			}),
			Timeout: transport.Timeout{
				Read:  20 * time.Second,
				Write: 90 * time.Second,
				Batch: 100 * time.Millisecond,
			},
			Version: api.Version,
		}, got)
	})

	t.Run("namespaces", func(t *testing.T) {
		transport.New = func(_ context.Context, cfg transport.Config) (internal.Transport, error) {
			return testkit.NopTransport, nil
		}

		c, err := weaviate.NewClient(t.Context())
		require.NoError(t, err)
		require.NotNil(t, c, "nil client")

		assert.NotNil(t, c.Collections, "nil collections")
		assert.NotNil(t, c.Backup, "nil backup")
	})
}

// DO NOT enable t.Parallel() for this test as it modifies the global state.
func TestWithAPIKey(t *testing.T) {
	newFunc := transport.New
	t.Cleanup(func() { transport.New = newFunc })

	var got transport.Config
	transport.New = func(_ context.Context, cfg transport.Config) (internal.Transport, error) {
		got = cfg
		return testkit.NopTransport, nil
	}

	c, err := weaviate.NewClient(t.Context(), weaviate.WithAPIKey("api-key"))
	assert.NoError(t, err, "new client")
	assert.NotNil(t, c, "nil client")

	require.NotNil(t, got.Auth, "token source")
	if assert.Implements(t, (*oauth2.TokenSource)(nil), got.Auth, "auth provider") {
		src := got.Auth.(oauth2.TokenSource)
		tok, err := src.Token()
		assert.NoError(t, err, "token error")

		assert.Zero(t, tok.RefreshToken, "refresh token")
		assert.Zero(t, tok.ExpiresIn, "expires in")
		assert.Zero(t, tok.Expiry, "expires in")
		assert.Equal(t, "Bearer", tok.Type(), "token type")
	}
}

// DO NOT enable t.Parallel() for this test as it modifies the global state.
func TestOIDCAuthentication(t *testing.T) {
	newFunc := transport.New
	t.Cleanup(func() { transport.New = newFunc })

	for _, tt := range []struct {
		name string
		opt  weaviate.Option
		auth any
	}{
		{
			name: "bearer token",
			opt:  weaviate.WithBearerToken(oauth2.Token{}),
			auth: auth.RefreshToken(oauth2.Token{}),
		},
		{
			name: "client credentials",
			opt:  weaviate.WithClientCredentials("secret", []string{"email"}),
			auth: auth.ClientCredentials{
				ClientSecret: "secret",
				Scopes:       []string{"email"},
			},
		},
		{
			name: "ropc",
			opt:  weaviate.WithResourceOwnerPasswordCredentials("secret", "john_doe", "xxx", []string{"email"}),
			auth: auth.ResourceOwnerPasswordCredentials{
				ClientSecret: "secret",
				Username:     "john_doe",
				Password:     "xxx",
				Scopes:       []string{"email"},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var got transport.Config
			transport.New = func(_ context.Context, cfg transport.Config) (internal.Transport, error) {
				got = cfg
				return testkit.NopTransport, nil
			}

			c, err := weaviate.NewClient(t.Context(), tt.opt)
			assert.NoError(t, err, "new client")
			assert.NotNil(t, c, "nil client")

			assert.Equal(t, tt.auth, got.Auth, "bad auth provider")
		})
	}
}

// DO NOT enable t.Parallel() for this test as it modifies the global state.
func TestLiveReady(t *testing.T) {
	newFunc := transport.New
	t.Cleanup(func() { transport.New = newFunc })

	transportOK := testkit.NopTransport
	transportErr := testkit.TransportFunc(func(context.Context, any, any) error { return testkit.ErrWhaam })

	for _, tt := range []struct {
		name      string
		act       func(*weaviate.Client, context.Context) (bool, error)
		transport internal.Transport
		err       testkit.Error
		want      bool
	}{
		{
			name:      "live ok",
			act:       (*weaviate.Client).IsLive,
			transport: transportOK,
			want:      true,
		},
		{
			name:      "live err",
			act:       (*weaviate.Client).IsLive,
			transport: transportErr,
			err:       testkit.ExpectError,
		},
		{
			name:      "ready ok",
			act:       (*weaviate.Client).IsReady,
			transport: transportOK,
			want:      true,
		},
		{
			name:      "ready err",
			act:       (*weaviate.Client).IsReady,
			transport: transportErr,
			err:       testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport.New = func(_ context.Context, cfg transport.Config) (internal.Transport, error) {
				return tt.transport, nil
			}

			c, err := weaviate.NewClient(t.Context())
			require.NoError(t, err, "new client")
			require.NotNil(t, c, "nil client")

			got, err := tt.act(c, t.Context())
			tt.err.Assert(t, err, "request error")
			assert.Equal(t, tt.want, got, "bad result")
		})
	}
}

func TestMetadata(t *testing.T) {
	newFunc := transport.New
	t.Cleanup(func() { transport.New = newFunc })

	tport := testkit.NewTransport(t, []testkit.Stub[any, any]{{
		Request: testkit.Ptr[any](api.GetInstanceMetadataRequest),
		Response: api.GetInstanceMetadataResponse{
			Hostname: "example.com",
			Version:  "v1.37.0",
			Modules: map[string]any{
				"text2vec-weaviate": true,
				"backup-s3":         true,
			},
			GRPCMaxMessageSize: 4096,
		},
	}})
	transport.New = func(_ context.Context, cfg transport.Config) (internal.Transport, error) {
		return tport, nil
	}

	c, err := weaviate.NewClient(t.Context())
	require.NoError(t, err, "new client")
	require.NotNil(t, c, "nil client")

	got, err := c.Metadata(t.Context())
	assert.NoError(t, err, "request error")
	assert.EqualValues(t, &weaviate.InstanceMetadata{
		Hostname: "example.com",
		Version:  "v1.37.0",
		Modules: map[string]any{
			"text2vec-weaviate": true,
			"backup-s3":         true,
		},
		GRPCMaxMessageSize: 4096,
	}, got, "bad result")
}
