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
type RPC[Client any] func(context.Context, Client) error

func (c *GRPC[Client]) Do(ctx context.Context, rpc RPC[Client]) error {
	dev.AssertNotNil(rpc, "rpc")

	if err := rpc(ctx, c.client); err != nil {
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

	var callOpts []grpc.CallOption
	if cfg.MaxMessageSize > 0 {
		callOpts = append(callOpts,
			grpc.MaxCallSendMsgSize(cfg.MaxMessageSize),
			grpc.MaxCallRecvMsgSize(cfg.MaxMessageSize),
		)
	}

	dialOpts := []grpc.DialOption{
		grpc.WithDefaultCallOptions(callOpts...),
	}

	if cfg.Header != nil {
		dialOpts = append(dialOpts, withDefaultHeader(*cfg.Header))
	}

	transportCreds := insecure.NewCredentials()
	if cfg.TLS {
		transportCreds = credentials.NewTLS(nil)
	}
	dialOpts = append(dialOpts, grpc.WithTransportCredentials(transportCreds))

	if cfg.TokenSource != nil {
		dialOpts = append(dialOpts, grpc.WithPerRPCCredentials(
			&tokenSource{
				TokenSource: cfg.TokenSource,
				tls:         cfg.TLS,
			},
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

// Client returns the underlying client instance.
func (c *GRPC[Client]) Client() Client {
	return c.client
}

var _ io.Closer = (*GRPC[any])(nil)

func (c *GRPC[Client]) Close() error {
	return c.channel.Close()
}

// withDefaultHeader creates an interceptor that adds md headers to the request context.
func withDefaultHeader(md metadata.MD) grpc.DialOption {
	return grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var pairs []string
		for k, v := range md {
			if len(v) == 0 {
				continue
			}
			pairs = append(pairs, k, v[0])
		}
		return invoker(metadata.AppendToOutgoingContext(ctx, pairs...), method, req, reply, cc, opts...)
	})
}

type tokenSource struct {
	oauth2.TokenSource
	tls bool // Enable transport level security.
}

var _ credentials.PerRPCCredentials = (*tokenSource)(nil)

// GetRequestMetadata gets the request metadata as a map from a [tokenSource].
// This behaves exactly like [oauth.TokenSource.GetRequestMetadata] but omits
// security level check if TLS is disabled.
//
// We disable this precaution knowingly to allow the users to use authenticated
// connections on an unprotected (possibly private) network.
//
// See also: [credentials.CheckSecurityLevel].
func (ts *tokenSource) GetRequestMetadata(ctx context.Context, _ ...string) (map[string]string, error) {
	token, err := ts.Token()
	if err != nil {
		return nil, err
	}
	if ts.tls {
		ri, _ := credentials.RequestInfoFromContext(ctx)
		if err = credentials.CheckSecurityLevel(ri.AuthInfo, credentials.PrivacyAndIntegrity); err != nil {
			return nil, fmt.Errorf("unable to transfer TokenSource PerRPCCredentials: %v", err)
		}
	}
	return map[string]string{
		"authorization": token.Type() + " " + token.AccessToken,
	}, nil
}

func (ts *tokenSource) RequireTransportSecurity() bool { return ts.tls }
