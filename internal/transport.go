package internal

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v6/internal/api"
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
	// able to execute any one of those. Similarly, `dest` must
	// be a pointer to one of the response structs in "internal/api".
	Do(_ context.Context, req api.Request, dest any) error
}
