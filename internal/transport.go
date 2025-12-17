package internal

import (
	"context"
)

type Transport interface {
	// Do executes a request and populates the response object.
	// Response dest SHOULD be nil if no response is expected
	// and MUST be a non-nil pointer otherwise.
	//
	// To give the caller maximum flexibility in defining the requests
	// while keeping the execution transparent, [Transport] does not
	// enforce any constraints on the request type.
	// A request is anything that has a body, i.e. anything can be a request.
	Do(ctx context.Context, req any, dest any) error
}
