package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
	Body() any
}

// StatusAccepter is an interface that response types can implement to
// control which HTTP status codes > 299 should not result in an error.
//
// Transport always treats codes < 299 as successful and only call
// AcceptStatus with codes > 299.
type StatusAccepter interface {
	// AcceptStatus returns true if a status code is acceptable.
	AcceptStatus(code int) bool
}

func (c *httpClient) do(ctx context.Context, req Endpoint, dest any) error {
	var body io.Reader
	if b := req.Body(); b != nil {
		marshaled, err := json.Marshal(b)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		log.Printf("%s", marshaled)
		body = bytes.NewReader(marshaled)
	}

	url := c.url(req)
	httpreq, err := http.NewRequestWithContext(ctx, req.Method(), url, body)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	// Clone default request headers.
	httpreq.Header = c.header.Clone()

	// Accept JSON even if dest is nil, as we don't want to spoof the request.
	httpreq.Header.Set("Accept", "application/json")
	if body != nil {
		httpreq.Header.Set("Content-Type", "application/json")
	}

	res, err := c.c.Do(httpreq)
	if err != nil {
		return fmt.Errorf("execute request: %q", err)
	}

	// Response body SHOULD always be read completely and closed
	// to allow the underlying [http.Transport] to re-use the TCP connection.
	// See: https://pkg.go.dev/net/http#Client.Do
	resBody, err := io.ReadAll(res.Body)
	res.Body.Close()

	// TODO(dyma): not sure if we should always report this error.
	// What if we don't need the body because dest=nil and status is OK?
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if res.StatusCode > 299 {
		// Some request types may want to "swallow" an bad status code,
		// and not return an error in that case. We will not try to unmarshal
		// the body in this case as it may not contain valid JSON, in which case
		// we'll have to return an error anyways.
		if acc, ok := dest.(StatusAccepter); ok {
			if acc.AcceptStatus(res.StatusCode) {
				return nil
			}
		}
		// TODO(dyma): better error handling?
		return fmt.Errorf("HTTP %d: %s", res.StatusCode, resBody)
	}

	if dest != nil {
		if err := json.Unmarshal(resBody, dest); err != nil {
			return fmt.Errorf("unmarshal response body: %w", err)
		}
	}
	return nil
}

type httpClient struct {
	c       *http.Client
	baseURL string
	header  http.Header
}

func newHTTP(opt Config) *httpClient {
	baseURL := fmt.Sprintf(
		"%s://%s:%d/%s/",
		opt.Scheme, opt.HTTPHost, opt.HTTPPort, opt.Version,
	)
	return &httpClient{
		c:       &http.Client{},
		baseURL: baseURL,
		header:  opt.Header,
	}
}

func (c *httpClient) url(req Endpoint) string {
	var url strings.Builder

	url.WriteString(c.baseURL)
	url.WriteString(strings.TrimLeft(req.Path(), "/"))

	if query := req.Query(); len(query) > 0 {
		url.WriteString("?")
		url.WriteString(query.Encode())
	}

	return url.String()
}

// BaseEndpoint implements [Endpoint] methods which may return nil.
// These values are usually optional in the request.
// BaseEndpoint can be embeded in a another request struct to reduce boilerplate.
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

// StaticEndpoint creates a new static endpoint with a method and a path.
func StaticEndpoint(method, path string) *staticEndpoint {
	return &staticEndpoint{method: method, path: path}
}

// staticEndpoint implements [Endpoint] for requests that has
// neither query or path parameters, nor the body. Each request
// to a static endpoint looks identical to all other ones, e.g. GET /live.
// Since such request is independent of its inputs (it has none of them),
// it can be created once and reused across the program's lifetime.
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
