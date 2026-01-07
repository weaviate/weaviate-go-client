package internal

import (
	"context"
)

type Transport[Req any, Dest any] interface {
	// Do executes a request and populates the response object.
	// Response dest SHOULD be nil if no response is expected
	// and MUST be a non-nil pointer otherwise.
	Do(ctx context.Context, req Req, dest Dest) error
}

type TransportFunc[Req any, Dest any] func(context.Context, Req, Dest) error
