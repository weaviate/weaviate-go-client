package api

import (
	"context"

	proto "github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/proto/v1"
	"google.golang.org/grpc"
)

// Message describes a gRPC request that can be executed via [proto.WeaviateClient].
//
// gRPC requests are modeled after the REST [Endpoint]; being comfortable with this
// analogy is helpful in reasoning about the interfaces used here.
type Message[In RequestMessage, Out ReplyMessage] interface {
	// Method informs the transport which method to call on the underlying client.
	Method() MethodFunc[In, Out]
	// Body returns an object that can marshal itself into a valid request message.
	Body() MessageMarshaler[In]
}

// MessageMarshaler marshals itself into a protobuf request message.
type MessageMarshaler[In RequestMessage] interface {
	MarshalMessage() (*In, error)
}

// UnmarshalMessage unmarshals a protobuf reply message.
type MessageUnmarshaler[Out ReplyMessage] interface {
	UnmarshalMessage(*Out) error
}

// RequestMessage enumerates all gRPC requests supported by [proto.WeaviateClient].
type RequestMessage interface {
	proto.SearchRequest
}

// ReplyMessage enumerates gRPC replies supported by [proto.WeaviateClient].
type ReplyMessage interface {
	proto.SearchReply
}

// MethodFunc is a method of the proto.WeaviateClient interface
// that accepts request In and returns reply Out.
type MethodFunc[In RequestMessage, Out ReplyMessage] func(proto.WeaviateClient, context.Context, *In, ...grpc.CallOption) (*Out, error)
