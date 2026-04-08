package grpc

import "google.golang.org/grpc/keepalive"

type Config struct {
	// Secured set it to true if it's a secured connection
	Secured bool
	// Host of the Weaviate instance, this is a mandatory field.
	// If host is without a port number then the 80 port
	// for insecured and 443 port for secured connections will be used.
	Host string
	// Keepalive parameters for the gRPC connection.
	// If nil, defaults are used (Time: 30s, Timeout: 10s, PermitWithoutStream: true).
	Keepalive *keepalive.ClientParameters
}
