package internal

import (
	"context"
)

type Transport interface {
	// Do executes a request and populates the response object.
	// Response dest SHOULD be nil if no response is expected
	// and MUST be a non-nil pointer otherwise.
	//
	// To keep execution transparent to the caller, the request type
	// only enforces a minimal constraint -- a request is anything
	// that MAY have a body.
	//
	// The "internal/api" package defines structs for all
	// supported requests, which in turn implement api.Request.
	// The contract is that Transport is able to execute any
	// one of those requests.
	//
	// The transport is also able to execute any custom [api.Endpoint].
	Do(ctx context.Context, req Request, dest any) error
}

// Request is anything that can have a body.
type Request interface {
	Body() any
}
