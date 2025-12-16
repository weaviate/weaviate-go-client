package api

import "net/url"

const Version = "v1"

type RequestDefaults struct {
	CollectionName   string
	Tenant           string
	ConsistencyLevel ConsistencyLevel
}

// endpoint implements [transport.Endpoint] methods which may return nil.
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
//	var _ transport.Endpoint = (*ListSongsRequest)(nil)
//
//	func (r *ListSongsRequest) Method() string { return http.MethodGet }
//	func (r *ListSongsRequest) Path() string   { return r.Artist + "/songs" }
type endpoint struct{}

func (*endpoint) Query() url.Values { return nil }
func (*endpoint) Body() any         { return nil }
