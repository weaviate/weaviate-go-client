package internal

import (
	"context"
	"errors"

	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
	"github.com/weaviate/weaviate-go-client/v6/internal/gen/proto/v1"
)

type Transport interface {
	// Do executes a request and populates the response object.
	// Response dest SHOULD be nil if no response is expected
	// and MUST be a non-nil pointer otherwise.
	//
	// To keep execution transparent to the caller, the request type
	// does not enforce any explicit constraints. E.g. were request
	// an interface with a method like Type() "rest" | "grpc", the
	// caller would have to be aware of the execution details.
	//
	// Instead, "internal/api" package defines structs for all
	// supported requests. The contract is that Transport is
	// able to execute any one of those. Similarly,

	// Both `req`  and `dest` MUST be pointers to an "internal/api" struct.
	Do(_ context.Context, req api.Request, dest any) error
}

func NewTransport() Transport {
	// TODO(dyma): initialize correctly
	return &transport{}
}

type transport struct {
	gRPC proto.WeaviateClient
}

// Compile-time assertion that transport implements Transport.
var _ Transport = (*transport)(nil)

// Do switches dispatches to the appropriate execution method depending on the request type.
func (t *transport) Do(ctx context.Context, req api.Request, dest any) error {
	switch req := req.(type) {
	case *api.SearchRequest:
		return t.search(ctx, req, dev.AssertType[*api.SearchResponse](dest))
	}
	return nil
}

func (t *transport) search(ctx context.Context, req *api.SearchRequest, dest *api.SearchResponse) error {
	reply, err := t.gRPC.Search(ctx, api.MarshalSearchRequest(req))
	if err != nil || dest == nil {
		return err
	}
	if reply == nil {
		// Since gRPC client is generated and is essentialy a third-party dependency,
		// we cannot guarantee the response to be always non-nil, so we do not dev.Assert.
		return errors.New("nil response")
	}
	*dest = *api.UnmarshalSearchReply(reply)
	return nil
}
