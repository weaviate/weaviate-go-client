package weaviate_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/weaviate/weaviate-go-client/v6"
	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/transport"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
)

// DO NOT enable t.Parallel() for this test as it messes with the global state.
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
		}, got)
	})

	t.Run("with options", func(t *testing.T) {
		var got transport.Config
		transport.New = func(_ context.Context, cfg transport.Config) (internal.Transport, error) {
			got = cfg
			return testkit.NopTransport, nil
		}

		c, err := weaviate.NewLocal(t.Context(),
			weaviate.WithScheme("https"),
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
			Timeout: transport.Timeout{
				Read:  20 * time.Second,
				Write: 90 * time.Second,
				Batch: 100 * time.Millisecond,
			},
		}, got)
	})
}

// DO NOT enable t.Parallel() for this test as it messes with the global state.
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
			Timeout: transport.Timeout{
				Read:  30 * time.Second,
				Write: 90 * time.Second,
			},
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
			Timeout: transport.Timeout{
				Read:  20 * time.Second,
				Write: 90 * time.Second,
				Batch: 100 * time.Millisecond,
			},
		}, got)
	})

	t.Run("namespaces", func(t *testing.T) {
		transport.New = func(_ context.Context, cfg transport.Config) (internal.Transport, error) {
			return testkit.NopTransport, nil
		}

		c, err := weaviate.NewClient(t.Context())
		assert.NoError(t, err)
		assert.NotNil(t, c, "nil client")

		assert.NotNil(t, c.Collections, "nil collections")
		assert.NotNil(t, c.Backup, "nil backup")
	})
}
