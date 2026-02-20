package transport

import (
	"net/http"

	"github.com/weaviate/weaviate-go-client/v6/internal"
)

type Config struct {
	Scheme   string      // Scheme for request URLs, "http" or "https".
	RESTHost string      // Hostname of the REST host.
	RESTPort int         // Port number of the REST host
	GRPCHost string      // Hostname of the gRPC host.
	GRPCPort int         // Port number of the gRPC host.
	Header   http.Header // Request headers.
}

// NewFunc returns an [internal.Transport] instance for [transport.Config].
type NewFunc func(Config) (internal.Transport, error)

var New NewFunc = newTransport

func newTransport(Config) (internal.Transport, error) {
	return nil, nil
}
