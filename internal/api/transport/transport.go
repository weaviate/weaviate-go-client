package transport

import (
	"net/http"
	"time"

	"github.com/weaviate/weaviate-go-client/v6/internal"
)

type Config struct {
	Scheme   string      // Scheme for request URLs, "http" or "https".
	RESTHost string      // Hostname of the REST host.
	RESTPort int         // Port number of the REST host
	GRPCHost string      // Hostname of the gRPC host.
	GRPCPort int         // Port number of the gRPC host.
	Header   http.Header // Request headers.
	Timeout  Timeout     // Request timeout options.
}

// Timeout sets client-side timeouts.
type Timeout struct {
	// Timeout for REST requests using HTTP GET or HEAD methods,
	// and gRPC requests using [WeaviateClient.Search],
	// [WeaviateClient.Aggregate], or [WeaviateClient.TenantsGet] methods.
	Read time.Duration

	// Timeout for REST requests using HTTP POST, PUT, PATCH, or DELETE methods,
	// and gRPC requests using [WeaviateClient.BatchDelete],
	// [WeaviateClient.BatchObjects] or [WeaviateClient.BatchReferences] methods.
	Write time.Duration // Timeout for insert requests.
	Batch time.Duration // Timeout for batch insert requests.
}

// NewFunc returns an [internal.Transport] instance for [transport.Config].
type NewFunc func(Config) (internal.Transport, error)

var New NewFunc = newTransport

func newTransport(Config) (internal.Transport, error) {
	return nil, nil
}
