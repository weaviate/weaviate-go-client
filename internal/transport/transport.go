package transport

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
)

// Do dispatches to the appropriate underlying transport depending on the request type.
// [Endpoint] is executed as REST requests. [Message] is executed via gRPC.
func (t *T) Do(ctx context.Context, req internal.Request, dest any) error {
	switch req := req.(type) {
	case internal.Endpoint:
		return t.http.do(ctx, req, dest)
	case internal.Message:
		return t.gRPC.do(ctx, req, dest)
	default:
		dev.Assert(false, "unknown request type %T", req)
	}
	return nil
}

// Config options for [T].
type Config struct {
	Scheme   string // Scheme for request URLs, "http" or "https".
	HTTPHost string // Hostname of the REST host.
	HTTPPort int    // Port number of the REST host
	GRPCHost string // Hostname of the gRPC host.
	GRPCPort int    // Port number of the gRPC host.
	Header   http.Header
	Timeout  time.Duration
	Version  string // Version of the REST API.
	// TODO: Authentication, Timeout

	// Ping forces [NewTransport] to try and connect to the gRPC server.
	// By default [grpc.Client] will only establish a connection on the first call
	// to one of its methods to avoid I/O on instantiation.
	Ping bool
}

func New(opt Config) (*T, error) {
	gRPC, err := newGRPCClient(opt)
	if err != nil {
		return nil, err
	}
	return &T{
		gRPC: gRPC,
		http: newHTTP(opt),
	}, nil
}

// T is an implementation of the [internal.Transport] interface.
// It executes REST requests via an [http.Client] and
// uses proto.WeaviateClient to perform gRPC requests.
//
// Use [New] to create a new transport. Call Close()
// when the T is not longer in use to free resources.
type T struct {
	gRPC *gRPCClient
	http *httpClient
}

// Compile-time assertion that transport implements Transport.
var _ io.Closer = (*T)(nil)

// Close closes the underlying gRPC channel.
func (t *T) Close() error {
	return t.gRPC.Close()
}
