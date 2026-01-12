package weaviate_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weaviate/weaviate-go-client/v6"
	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
)

// DO NOT enable t.Parallel() for this test as it messes with the global state.
func TestNewLocal(t *testing.T) {
	tf := api.NewTransport
	t.Cleanup(func() { api.NewTransport = tf })

	var nop testkit.NopTransport

	t.Run("default", func(t *testing.T) {
		var got api.TransportConfig
		api.NewTransport = func(cfg api.TransportConfig) (internal.Transport, error) {
			got = cfg
			return &nop, nil
		}

		c, err := weaviate.NewLocal(t.Context(), nil)
		assert.NotNil(t, c, "nil client")
		assert.NoError(t, err)

		assert.Equal(t, api.TransportConfig{
			Scheme:   "http",
			RESTHost: "localhost",
			GRPCHost: "localhost",
			RESTPort: 8080,
			GRPCPort: 50051,
			Header: http.Header{
				"X-Weaviate-Client": {"weaviate-client-go" + "/" + weaviate.Version()},
			},
		}, got, "default local config")
	})

	t.Run("custom", func(t *testing.T) {
		var got api.TransportConfig
		api.NewTransport = func(cfg api.TransportConfig) (internal.Transport, error) {
			got = cfg
			return &nop, nil
		}

		c, err := weaviate.NewLocal(t.Context(), &weaviate.ConnectionConfig{
			Scheme:   "https",
			HTTPPort: 7070,
			GRPCPort: 54321,
			Header: http.Header{
				"X-Test": {"heads", "up"},
			},
		})
		assert.NotNil(t, c, "nil client")
		assert.NoError(t, err)

		assert.Equal(t, api.TransportConfig{
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
		}, got, "default local config")
	})
}
