package transport

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/gen/proto/v1"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// do obtains a protobuf message from the request body and dispatches
// to the appropriate proto.WeaviateClient method based on its kind.
func (c *gRPCClient) do(ctx context.Context, req internal.Message, dest any) error {
	dev.Assert(req != nil, "nil gRPC request")

	msg := req.NewMessage()
	dev.Assert(msg != nil, "nil gRPC message")

	switch msg := msg.(type) {
	case *proto.SearchRequest:
		return c.search(ctx, msg, dev.AssertType[*api.SearchResponse](dest))
	}
	return nil
}

func newGRPCClient(opt internal.TransportOptions) (*gRPCClient, error) {
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

func (c *gRPCClient) search(ctx context.Context, req *proto.SearchRequest, dest *api.SearchResponse) error {
	reply, err := c.wc.Search(ctx, req)
	if err != nil || dest == nil {
		return err
	}
	if reply == nil {
		// Since gRPC client is generated and is essentialy a third-party dependency,
		// we cannot guarantee the response to be always non-nil, so we do not dev.Assert.
		return errors.New("nil reply")
	}
	*dest = *api.NewSearchResponse(reply)
	return nil
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
