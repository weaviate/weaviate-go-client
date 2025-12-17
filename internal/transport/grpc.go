package transport

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/weaviate/weaviate-go-client/v6/internal/api/gen/proto/v1"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// Since gRPC client is generated and is essentialy a third-party dependency,
// we cannot guarantee the response to be always non-nil, so we return an error
// on nil replies instead of doing dev.Assert.
var errNilReply = errors.New("nil reply")

// Message describes a gRPC request.
type Message[Request any, Reply any] interface {
	// Marshal creates a protobuf message holding the body of the request.
	//
	// In practice, the request struct will need to marshal itself
	// into an appropriate protobuf stub from the "internal/api/gen/proto" package.
	Marshal() *Request

	// Unmarshal converts the gRPC reply into a response type
	// corresponding to the kind of request and assigns it to dest.
	Unmarshal(r *Reply, dest any) error
}

// do obtains a protobuf message from the request body and dispatches
// to the appropriate proto.WeaviateClient method based on its kind.
func (c *gRPCClient) do(ctx context.Context, req any, dest any) error {
	dev.Assert(req != nil, "nil gRPC request")

	var err error
	switch m := req.(type) {
	case Message[proto.SearchRequest, proto.SearchReply]:
		err = c.search(ctx, m, dest)
	default:
		dev.Assert(false, "unknown gRPC message type %T", m)
	}

	if err != nil {
		return fmt.Errorf("gRPC: %w", err)
	}
	return nil
}

func (c *gRPCClient) search(ctx context.Context, m Message[proto.SearchRequest, proto.SearchReply], dest any) error {
	reply, err := c.wc.Search(ctx, m.Marshal())
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}

	if reply == nil {
		// Since gRPC client is generated and is essentialy a third-party dependency,
		// we cannot guarantee the response to be always non-nil, so we do not dev.Assert.
		return fmt.Errorf("search: %w", errNilReply)
	}

	if err := m.Unmarshal(reply, dest); err != nil {
		return fmt.Errorf("search: unmarshal respose: %w", err)
	}
	return nil
}

func newGRPCClient(opt Config) (*gRPCClient, error) {
	// TODO(dyma): apply relevant gRPC options.
	channel, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", opt.GRPCHost, opt.GRPCPort),
		// TODO(dyma): pass correct credentials if authentication is enabled or scheme == "https"
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.Header((*metadata.MD)(&opt.Header)),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("create gRPC channel: %w", err)
	}
	return &gRPCClient{
		channel: channel,
		wc:      proto.NewWeaviateClient(channel),
	}, nil
}

// gRPCClient is a wrapper around proto.WeaviateClient that dispatches protobuf messages
// and manages related client resources, i.e. the gRPC channel.
type gRPCClient struct {
	channel *grpc.ClientConn
	wc      proto.WeaviateClient
}

var _ io.Closer = (*gRPCClient)(nil)

func (c *gRPCClient) Close() error {
	return c.channel.Close()
}
