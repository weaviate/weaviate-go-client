package transports

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
)

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

// StatusAccepter is an interface that response types can implement to
// control which HTTP status codes > 299 should not result in an error.
//
// The [REST] transport always treats codes < 299 as successful
// and only calls AcceptStatus for codes > 299.
type StatusAccepter interface {
	// AcceptStatus returns true if a status code is acceptable.
	AcceptStatus(code int) bool
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

// IdentityEndpoint returns a factory for requests to identity endpoints.
//
// Example:
//
//	// Assign request creator func to a variable.
//	var DeleteSongRequest = transport.IdentityEndpoint[int](http.MethodDelete, "/songs/%v")
//
//	// Use it to create new requests
//	req := DeleteSongRequest(123)
//
// The pathFmt MUST ONLY contain a single formatting directive. Callers are free to use
// the formatting directive most appropriate to the identity type, e.g. %s for strings,
// but %d for numbers.
//
// IdentityEndpoint will panic on invalid pathFmt _before_ the factory is created.
//
//	// BAD: panics because pathFmt contains 2 formatting directives.
//	var DeleteSongsRequest = transport.IdentityEndpoint[uuid.UUID](http.MethodGet, "/albums/%v/songs/%v")
func IdentityEndpoint[I any](method, pathFmt string) func(I) any {
	dev.Assert(strings.Count(pathFmt, "%") == 1, "%s must have a single formatting directive", pathFmt)

	return func(id I) any {
		return &identityEndpoint[I]{
			method:  method,
			pathFmt: pathFmt,
			id:      id,
		}
	}
}

// identityEndpoint implements [Endpoint] for endpoints that use some ID parameter
// in the request path, e.g. 'DELETE /users/<user-id>' or 'GET /artist/<artist-id>/songs'.
// See [IdentityEndpoint] for more details.
type identityEndpoint[I any] struct {
	BaseEndpoint
	method, pathFmt string
	id              I
}

var _ Endpoint = (*identityEndpoint[any])(nil)

func (r *identityEndpoint[T]) Method() string { return r.method }
func (r *identityEndpoint[T]) Path() string   { return fmt.Sprintf(r.pathFmt, r.id) }

// StaticEndpoint creates a new static endpoint with a method and a path.
func StaticEndpoint(method, path string) *staticEndpoint {
	return &staticEndpoint{method: method, path: path}
}

// staticEndpoint implements [Endpoint] for requests that has
// neither query or path parameters, nor the body. Each request
// to a static endpoint looks identical to all other ones, e.g. GET /live.
// Since such request is independent of its inputs (it has none of them),
// it can be created once and re-used throughout the program's lifetime.
//
// Example:
//
//	var ListSongsRequest transport.Endpoint = transport.StaticEndpoint(http.MethodGet, "/songs")
type staticEndpoint struct {
	BaseEndpoint
	method, path string
}

var _ Endpoint = (*staticEndpoint)(nil)

func (e *staticEndpoint) Method() string { return e.method }
func (e *staticEndpoint) Path() string   { return e.path }
