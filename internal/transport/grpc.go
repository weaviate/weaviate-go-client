package transport

import (
	"fmt"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// GRPCConfig options for [grpc.ClientConn].
type GRPCConfig struct {
	Scheme  string // Scheme for request URLs, "http" or "https".
	Host    string // Hostname of the GRPC host.
	Port    int    // Port number of the GRPC host
	Header  http.Header
	Timeout time.Duration
	// TODO: Authentication, Timeout

	// Ping forces [NewTransport] to try and connect to the gRPC server.
	// By default [grpc.Client] will only establish a connection on the
	// first call to one of its methods to avoid I/O on instantiation.
	Ping bool
}

func NewChannel(cfg GRPCConfig) (*grpc.ClientConn, error) {
	// TODO(dyma): apply relevant gRPC options.
	return grpc.NewClient(
		fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		// TODO(dyma): pass correct credentials if authentication is enabled or scheme == "https"
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.Header((*metadata.MD)(&cfg.Header)),
		),
	)
}
