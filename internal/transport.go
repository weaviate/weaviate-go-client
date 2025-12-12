package internal

import "context"

type Transport interface {
	// Do executes a request and populates the response object.
	// Response dest SHOULD be nil if no response is expected and MUST
	// be a non-nil pointer otherwise.
	Do(_ context.Context, _ Request, dest any) error
}

// TODO: Transport defines its request types which map to protobuf / REST types, but in a separate `internal/request` package so that the execution is transparent to the caller.
// Something like:
//
//	type SearchRequest struct {
//		NearText   request.NearText
//		NearVector request.NearVector
//		BM25       request.BM25
//	}
//
// Used like so:
//
//	func (c *Client) Insert(ctx, ...) {
//		c.transport.Do(ctx, request.Insert{
//			Object: 	object,
//			Defaults: 	c.defaults, /* request.Defaults */
//		})
//	}
type Request any
