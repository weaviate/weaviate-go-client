package transports

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// Config for [GRPC] transport.
type GRPCConfig[C any] struct {
	Scheme         string // Scheme for request URLs, "http" or "https".
	Host           string // Hostname of the gRPC host.
	Port           int    // Port number of the gRPC host.
	Header         *metadata.MD
	Timeout        time.Duration
	MaxMessageSize int
	// TODO(dyma): Authentication

	// Ping forces [NewTransport] to try and connect to the gRPC server.
	// By default [grpc.Client] will only establish a connection on the first call
	// to one of its methods to avoid I/O on instantiation.
	Ping bool

	NewGRPCClient NewGRPCClientFunc[C]
}

type NewGRPCClientFunc[Client any] func(grpc.ClientConnInterface) Client

type RPC[Client any] interface {
	Do(context.Context, Client) error
}

func (g *GRPC[C]) Client() C {
	return g.client
}

func (c *GRPC[C]) Do(ctx context.Context, rpc RPC[C]) error {
	dev.Assert(rpc != nil, "nil rpc")

	if err := rpc.Do(ctx, c.client); err != nil {
		return fmt.Errorf("grpc: %w", err)
	}
	return nil
}

func NewGRPC[C any](cfg GRPCConfig[C]) (*GRPC[C], error) {
	callOpts := []grpc.CallOption{
		grpc.Header(cfg.Header),
	}
	if cfg.MaxMessageSize > 0 {
		callOpts = append(callOpts,
			grpc.MaxCallSendMsgSize(cfg.MaxMessageSize),
			grpc.MaxCallRecvMsgSize(cfg.MaxMessageSize),
		)
	}

	// TODO(dyma): apply relevant gRPC options.
	channel, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		// TODO(dyma): pass correct credentials if authentication is enabled or scheme == "https"
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(callOpts...),
	)
	if err != nil {
		return nil, fmt.Errorf("create gRPC channel: %w", err)
	}
	return &GRPC[C]{
		channel: channel,
		client:  cfg.NewGRPCClient(channel),
	}, nil
}

// GRPC is a wrapper around a protobuf client that dispatches messages
// and manages related client resources, i.e. the gRPC channel.
type GRPC[C any] struct {
	channel *grpc.ClientConn
	client  C
}

var _ io.Closer = (*GRPC[any])(nil)

func (c *GRPC[C]) Close() error {
	return c.channel.Close()
}
