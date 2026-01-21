package transports

import "net/url"

// Endpoint describes a REST request.
type Endpoint interface {
	// Method returns an HTTP method appropriate for the request.
	Method() string

	// Path returns endpoint URL with path parameters populated.
	Path() string

	// Query returns query string, if the request supports query parameters.
	// A request which does not have query parameters can safely return nil.
	Query() url.Values

	// Body returns the request body, which will be marshaled to JSON.
	// Requests that do not have a payload should return nil.
	Body() any
}

// BaseEndpoint implements [Endpoint] methods which may return nil.
// These values are usually optional in the request.
// BaseEndpoint can be embedded in a another request struct to reduce boilerplate.
//
// Example:
//
//	type DeleteSongRequest struct {
//		transport.BaseEndpoint
//		SongID string
//	}
//
//	// DeleteSongRequest implements transport.Endpoint.
//	var _ transport.Endpoint = (*DeleteSongRequest)(nil)
//
//	func (r *DeleteSongRequest) Method() string { return http.MethodDelete }
//	func (r *DeleteSongRequest) Path() string   { return "/songs/" + r.SongID }
type BaseEndpoint struct{}

func (*BaseEndpoint) Query() url.Values { return nil }
func (*BaseEndpoint) Body() any         { return nil }
