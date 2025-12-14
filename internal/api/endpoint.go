package api

import (
	"net/url"
)

type Endpoint interface {
	Method() string
	Path() string
	Query() url.Values
	Body() any
}

// endpoint implements Endpoint methods which may return nil.
// These values are usually optional in the request.
// endpoint can be embeded in a another request struct to reduce boilerplate.
//
// Example:
//
//	type ListSongsRequest struct {
//		endpoint
//		Artist string
//	}
//
//	// ListSongsRequest implements Endpoint.
//	var _ Endpoint = (*ListSongsRequest)(nil)
//
//	func (r *ListSongsRequest) Method() string { return http.MethodGet }
//	func (r *ListSongsRequest) Path() string   { return r.Artist + "/songs" }
type endpoint struct{}

func (*endpoint) Query() url.Values { return nil }
func (*endpoint) Body() any         { return nil }
