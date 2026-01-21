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
