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

// RequestMessage enumerates all gRPC requests accepted by this transport.
type RequestMessage interface {
	proto.SearchRequest |
		proto.AggregateRequest |
		proto.TenantsGetRequest |
		proto.BatchDeleteRequest |
		proto.BatchObjectsRequest |
		proto.BatchReferencesRequest
}

// ReplyMessage enumerates gRPC replies this transport can produce.
type ReplyMessage interface {
	proto.SearchReply |
		proto.AggregateReply |
		proto.TenantsGetReply |
		proto.BatchDeleteReply |
		proto.BatchObjectsReply |
		proto.BatchReferencesReply
}

// MessageMarshaler marshals the body of the request into a protobuf message.
type MessageMarshaler[R RequestMessage] interface {
	MarshalMessage() *R
}

// UnmarshalMessage unmarshals a protobuf message into the response object.
type MessageUnmarshaler[R ReplyMessage] interface {
	UnmarshalMessage(*R) error
}

// do dispatches to the appropriate proto.WeaviateClient method based on
// the request type. req MUST implement MessageMarshaler for one of RequestMessage types,
// and dest MUST implement MessageUnmarshaler for the corresponding reply.
func (c *gRPCClient) do(ctx context.Context, req any, dest any) error {
	dev.Assert(req != nil, "nil gRPC request")

	var err error
	switch m := req.(type) {
	case MessageMarshaler[proto.SearchRequest]:
		err = c.search(ctx, m, dest)
	case MessageMarshaler[proto.AggregateRequest]:
		err = c.aggregate(ctx, m, dest)
	case MessageMarshaler[proto.BatchDeleteRequest]:
		err = c.batchDelete(ctx, m, dest)
	case MessageMarshaler[proto.BatchObjectsRequest]:
		err = c.batchObjects(ctx, m, dest)
	case MessageMarshaler[proto.BatchReferencesRequest]:
		err = c.batchReferences(ctx, m, dest)
	default:
		dev.Assert(false, "%T does not implement MessageMarshaler for any of the supported request types", m)
	}

	if err != nil {
		return fmt.Errorf("gRPC: %w", err)
	}
	return nil
}

// Unmarshal unmarshals reply R into dest. A nil dest means the reply can be ignored,
// which returns with a nil error immediately. A nil reply returns an non-nil error.
// A dest that does not implement MessageUnmarshaler[R] returns a non-nil error.
// Otherwise UnmarshalMessage() is called with reply *R and the unmarshaling error is returned.
func unmarshal[R ReplyMessage](reply *R, dest any) error {
	if dest == nil {
		return nil
	}
	if reply == nil {
		// Since gRPC client is generated and is essentially a third-party dependency,
		// we cannot guarantee the response to be always non-nil, so we return an error
		// on nil replies instead of doing dev.Assert.
		return errors.New("nil reply")
	}
	if u, ok := dest.(MessageUnmarshaler[R]); ok {
		return u.UnmarshalMessage(reply)
	}
	return fmt.Errorf(
		"cannot unmarshal %T into %T: dest does not implement %T",
		reply, dest, *new(MessageUnmarshaler[R]),
	)
}

func (c *gRPCClient) search(ctx context.Context, m MessageMarshaler[proto.SearchRequest], dest any) error {
	reply, err := c.wc.Search(ctx, m.MarshalMessage())
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}

	if err := unmarshal(reply, dest); err != nil {
		return fmt.Errorf("search: %w", err)
	}
	return nil
}

func (c *gRPCClient) aggregate(ctx context.Context, m MessageMarshaler[proto.AggregateRequest], dest any) error {
	reply, err := c.wc.Aggregate(ctx, m.MarshalMessage())
	if err != nil {
		return fmt.Errorf("aggregate: %w", err)
	}

	if err := unmarshal(reply, dest); err != nil {
		return fmt.Errorf("aggregate: %w", err)
	}
	return nil
}

func (c *gRPCClient) batchDelete(ctx context.Context, m MessageMarshaler[proto.BatchDeleteRequest], dest any) error {
	reply, err := c.wc.BatchDelete(ctx, m.MarshalMessage())
	if err != nil {
		return fmt.Errorf("batchDelete: %w", err)
	}

	if err := unmarshal(reply, dest); err != nil {
		return fmt.Errorf("batchDelete: %w", err)
	}
	return nil
}

func (c *gRPCClient) batchObjects(ctx context.Context, m MessageMarshaler[proto.BatchObjectsRequest], dest any) error {
	reply, err := c.wc.BatchObjects(ctx, m.MarshalMessage())
	if err != nil {
		return fmt.Errorf("batchObjects: %w", err)
	}

	if err := unmarshal(reply, dest); err != nil {
		return fmt.Errorf("batchObjects: %w", err)
	}
	return nil
}

func (c *gRPCClient) batchReferences(ctx context.Context, m MessageMarshaler[proto.BatchReferencesRequest], dest any) error {
	reply, err := c.wc.BatchReferences(ctx, m.MarshalMessage())
	if err != nil {
		return fmt.Errorf("batchReferences: %w", err)
	}

	if err := unmarshal(reply, dest); err != nil {
		return fmt.Errorf("batchReferences: %w", err)
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
