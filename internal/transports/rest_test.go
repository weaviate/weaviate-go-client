package transports_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/internal/transports"
)

func TestREST_Do(t *testing.T) {
	// 1. start a server
	// 2. send some request
	// 3. verify payload / query / body correct
	// 4. verify status code accepter
	// 5. verify unmarshaled JSON correctly
	// 6. verify errors correclty

	contentTypeHeader := http.Header{"Content-Type": {"application/json"}}
	defaultHeader := http.Header{"X-Default": {"a", "b"}}
	version := "v0"

	for _, tt := range []struct {
		name string
		req  transports.Endpoint

		respBody string // Set response body to return.
		respCode int    // Override returned status code (default: HTTP 200).

		// respBody will be read into dest; leave it unset.
		// want is the expected value after deserialization.
		dest, want  any
		assertError assert.ErrorAssertionFunc // Error assertion. (default: assert.NoError)
	}{
		{
			name: "no payload",
			req: &endpoint{
				method: http.MethodGet,
				path:   "/test",
			},
		},
		{
			name: "with query",
			req: &endpoint{
				method: http.MethodGet,
				path:   "/test",
				query:  url.Values{"limit": {"10"}, "format": {"short"}},
			},
		},
		{
			name: "with payload",
			req: &endpoint{
				method: http.MethodPost,
				path:   "/test",
				body:   testkit.Ptr(5),
			},
			respBody: "123",
			want:     float64(123), // "123" is unmarshaled into float64 by default.
		},
		{
			name: "with malformed payload",
			req: &endpoint{
				method: http.MethodPost,
				path:   "/test",
				body:   new(malformedBody),
			},
			assertError: func(t assert.TestingT, err error, msgAndArgs ...any) bool {
				return assert.ErrorIs(t, err, testkit.ErrWhaam, msgAndArgs...)
			},
		},
		{
			name: "malformed response body",
			req: &endpoint{
				method: http.MethodGet,
				path:   "/test",
			},
			respBody:    "{id: 123}",
			assertError: assert.Error,
		},
		{
			name: "error status code",
			req: &endpoint{
				method: http.MethodDelete,
				path:   "/test",
			},
			respCode: http.StatusPaymentRequired,
			respBody: "Payment Required",
			assertError: func(t assert.TestingT, err error, msgAndArgs ...any) (ok bool) {
				var httpErr *transports.HTTPError
				if assert.ErrorAs(t, err, &httpErr, msgAndArgs...) {
					ok = assert.Equal(t, http.StatusPaymentRequired, httpErr.Code, "status code") &&
						assert.Equal(t, "Payment Required", httpErr.Body, "response error body")
				}
				return
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			require.NotNil(t, tt.req, "bad test case: req is nil")

			// Arrange
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tt.req.Method(), r.Method, "request method")
				assert.Equal(t, "/"+version+tt.req.Path(), r.URL.Path, "request path")
				assert.Equal(t, tt.req.Query().Encode(), r.URL.Query().Encode(), "query parameters")
				assert.Subset(t, r.Header, defaultHeader, "default headers missing")

				// If request has no body, the r.Body is expected to be &http.noBody{}, never nil.
				require.NotNil(t, assert.NotNil(t, r.Body, "http.Request.Body is nil"))
				defer r.Body.Close()

				gotBody, err := io.ReadAll(r.Body)
				if assert.NoError(t, err, "read test request body") {
					if reqBody := tt.req.Body(); reqBody == nil {
						assert.Empty(t, gotBody, "request must have no body")
					} else {
						assert.Subset(t, r.Header, contentTypeHeader, "Content-Type header missing")

						wantBody, err := json.Marshal(reqBody)
						if assert.NoError(t, err, "marshal test request body") {
							assert.JSONEq(t, string(wantBody), string(gotBody), "bad request body")
						}
					}
				}

				code := http.StatusOK
				if tt.respCode != 0 {
					code = tt.respCode
				}
				w.WriteHeader(code)

				if tt.respBody != "" {
					io.WriteString(w, tt.respBody)
				}
			})

			srv := httptest.NewServer(handler)
			defer srv.Close()

			url, _ := url.Parse(srv.URL)
			port, _ := strconv.Atoi(url.Port())
			rest := transports.NewREST(transports.RESTConfig{
				Scheme:  "http",
				Host:    url.Hostname(),
				Port:    port,
				Version: version,
				Header:  defaultHeader, // todo: add default headers
			})

			assertError := assert.NoError
			if tt.assertError != nil {
				assertError = tt.assertError
			}

			// Act
			err := rest.Do(t.Context(), tt.req, &tt.dest)

			// Assert
			assertError(t, err, "request error")
			assert.Equal(t, tt.want, tt.dest, "bad response value")
		})
	}
}

// endpoint implements [transports.Endpoint] for testing.
type endpoint struct {
	method string
	path   string
	query  url.Values
	body   any
}

var _ transports.Endpoint = (*endpoint)(nil)

func (e *endpoint) Method() string    { return e.method }
func (e *endpoint) Path() string      { return e.path }
func (e *endpoint) Query() url.Values { return e.query }
func (e *endpoint) Body() any         { return e.body }

// malformedBody returns [testkit.ErrWhaam] when marshaled to JSON.
type malformedBody struct{}

var _ json.Marshaler = (*malformedBody)(nil)

func (*malformedBody) MarshalJSON() ([]byte, error) { return nil, testkit.ErrWhaam }

func TestBaseEndoint(t *testing.T) {
	var endpoint transports.BaseEndpoint

	assert.Nil(t, endpoint.Query(), "query")
	assert.Nil(t, endpoint.Body(), "body")
}

func TestStaticEndpoint(t *testing.T) {
	static := transports.StaticEndpoint(http.MethodGet, "/live")

	assert.Equal(t, static.Method(), http.MethodGet, "method")
	assert.Equal(t, static.Path(), "/live", "path")
	assert.Nil(t, static.Query(), "query")
	assert.Nil(t, static.Body(), "body")
}

func TestIdentityEndpoint(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		id := "test-id"
		pathFmt := "/string/%s"
		wantPath := fmt.Sprintf(pathFmt, id)

		factory := transports.IdentityEndpoint[string](http.MethodGet, pathFmt)
		req := factory(id)

		if assert.Implements(t, (*transports.Endpoint)(nil), req, "factory returns valid requests") {
			ep := req.(transports.Endpoint)

			assert.Equal(t, ep.Method(), http.MethodGet, "method")
			assert.Equal(t, ep.Path(), wantPath, "path")
			assert.Nil(t, ep.Query(), "query")
			assert.Nil(t, ep.Body(), "body")
		}
	})

	t.Run("int", func(t *testing.T) {
		id := 123
		pathFmt := "/int/%d"
		wantPath := fmt.Sprintf(pathFmt, id)

		factory := transports.IdentityEndpoint[int](http.MethodGet, pathFmt)
		req := factory(id)

		if assert.Implements(t, (*transports.Endpoint)(nil), req, "factory returns valid requests") {
			ep := req.(transports.Endpoint)

			assert.Equal(t, ep.Method(), http.MethodGet, "method")
			assert.Equal(t, ep.Path(), wantPath, "path")
			assert.Nil(t, ep.Query(), "query")
			assert.Nil(t, ep.Body(), "body")
		}
	})

	t.Run("uuid.UUID", func(t *testing.T) {
		id := uuid.New()
		pathFmt := "/uuid/%s"
		wantPath := fmt.Sprintf(pathFmt, id)

		factory := transports.IdentityEndpoint[uuid.UUID](http.MethodGet, pathFmt)
		req := factory(id)

		if assert.Implements(t, (*transports.Endpoint)(nil), req, "factory returns valid requests") {
			ep := req.(transports.Endpoint)

			assert.Equal(t, ep.Method(), http.MethodGet, "method")
			assert.Equal(t, ep.Path(), wantPath, "path")
			assert.Nil(t, ep.Query(), "query")
			assert.Nil(t, ep.Body(), "body")
		}
	})

	t.Run("invalid pathFmt", func(t *testing.T) {
		pathFmt := "/first/%v/second/%d"
		require.Panics(t, func() {
			transports.IdentityEndpoint[any](http.MethodGet, pathFmt)
		}, "must validate pathFmt on creation (%q has %d formatting directives)",
			pathFmt, strings.Count(pathFmt, "%"),
		)
	})
}
