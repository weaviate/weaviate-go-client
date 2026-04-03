package transports

import (
	"context"
	"fmt"
	"io"

	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/metadata"
)

// Config for [GRPC] transport.
type GRPCConfig[Client any] struct {
	Host           string             // Hostname of the gRPC host.
	Port           int                // Port number of the gRPC host.
	Header         *metadata.MD       // Headers added with each request.
	MaxMessageSize int                // Maximum gRPC message size in bytes.
	TokenSource    oauth2.TokenSource // OAuth2 token provider.
	TLS            bool               // If true, channel will use TLS protocol.

	NewGRPCClient NewGRPCClientFunc[Client]
}

// NewGRPCClientFunc creates a new instance of the underlying GRPC transport.
type NewGRPCClientFunc[Client any] func(grpc.ClientConnInterface) Client

// RPC describes a gRPC request in the given Client.
type RPC[Client any] interface {
	Do(context.Context, Client) error
}

func (c *GRPC[Client]) Do(ctx context.Context, rpc RPC[Client]) error {
	dev.AssertNotNil(rpc, "rpc")

	if err := rpc.Do(ctx, c.client); err != nil {
		return fmt.Errorf("grpc: %w", err)
	}
	return nil
}

// GRPC is a wrapper around a protobuf client that dispatches messages
// and manages related client resources, i.e. the gRPC channel.
//
// Unline [REST], which also takes care of request execution, marshaling
// and unmarshaling of the request/response payloads, GRPC is only concerned
// with resource management. This is because the generated Client stub will
// already contain serialization code, response status handling, and such.
type GRPC[Client any] struct {
	channel *grpc.ClientConn
	client  Client
}

func NewGRPC[Client any](cfg GRPCConfig[Client]) (*GRPC[Client], error) {
	dev.AssertNotNil(cfg.NewGRPCClient, "cfg.NewGRPCClient")

	callOpts := []grpc.CallOption{
		grpc.Header(cfg.Header),
	}
	if cfg.MaxMessageSize > 0 {
		callOpts = append(callOpts,
			grpc.MaxCallSendMsgSize(cfg.MaxMessageSize),
			grpc.MaxCallRecvMsgSize(cfg.MaxMessageSize),
		)
	}

	dialOpts := []grpc.DialOption{
		grpc.WithDefaultCallOptions(callOpts...),
	}

	transportCreds := insecure.NewCredentials()
	if cfg.TLS {
		transportCreds = credentials.NewTLS(nil)
	}
	dialOpts = append(dialOpts, grpc.WithTransportCredentials(transportCreds))

	if cfg.TokenSource != nil {
		dialOpts = append(dialOpts, grpc.WithPerRPCCredentials(
			&oauth.TokenSource{TokenSource: cfg.TokenSource},
		))
	}

	target := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	channel, err := grpc.NewClient(target, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("create gRPC channel: %w", err)
	}
	dev.AssertNotNil(channel, "channel")

	client := cfg.NewGRPCClient(channel)
	dev.AssertNotNil(client, "client")

	return &GRPC[Client]{
		channel: channel,
		client:  client,
	}, nil
}

var _ io.Closer = (*GRPC[any])(nil)

func (c *GRPC[Client]) Close() error {
	return c.channel.Close()
}
