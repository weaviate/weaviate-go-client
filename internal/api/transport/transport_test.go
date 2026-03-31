package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	proto "github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/proto/v1"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/internal/transports"
	"google.golang.org/grpc"
)

// Test that transport fetches instance's /meta information when created.
func TestNew(t *testing.T) {
	defaultHeader := http.Header{
		"X-Custom-Header": {"92"},
	}

	var fetchedMeta bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if assert.Equal(t, "/v0/meta", r.URL.Path, "request path") {
			fetchedMeta = true
			meta, err := json.Marshal(GetInstanceMetadataResponse{
				GRPCMaxMessageSize: 2048,
			})
			require.NoError(t, err, "marshal mock response")

			assert.Subset(t, r.Header, defaultHeader, "default headers missing")

			_, err = w.Write(meta)
			require.NoError(t, err, "write mock response")
		}
	}))
	t.Cleanup(srv.Close)

	url, err := url.Parse(srv.URL)
	require.NoError(t, err, "parse url")

	port, err := strconv.Atoi(url.Port())
	require.NoError(t, err, "parse port")

	tport, err := New(t.Context(), Config{
		Scheme:   url.Scheme,
		RESTHost: url.Hostname(),
		RESTPort: port,
		Header:   defaultHeader,
		Version:  "v0",
	})

	require.NoError(t, err)
	require.NotNil(t, tport, "nil transport")
	require.True(t, fetchedMeta, "transport must fetch instance metadata on startup")
}

func TestTransport_Do(t *testing.T) {
	t.Run("rest endpoint", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			var got *endpoint
			tport := transport{
				rest: restFunc(func(_ context.Context, req transports.Endpoint, dest any) error {
					got = req.(*endpoint) // capture actual request
					*(dest.(*int)) = 5    // write expected response
					return nil
				}),
			}

			req := &endpoint{
				method: http.MethodPost,
				path:   "/test",
				query:  url.Values{"limit": {"10"}},
				body:   92,
			}
			var resp int
			err := tport.Do(t.Context(), req, &resp)
			require.NoError(t, err, "transport error")

			require.Equal(t, req, got, "bad request")
			require.Equal(t, 5, resp, "bad response")
		})

		t.Run("error", func(t *testing.T) {
			tport := transport{
				rest: restFunc(func(_ context.Context, req transports.Endpoint, dest any) error {
					return testkit.ErrWhaam
				}),
			}

			err := tport.Do(t.Context(), &endpoint{}, nil)

			require.ErrorIs(t, err, testkit.ErrWhaam, "REST transport error not propagated")
		})

		t.Run("timeout", func(t *testing.T) {
			read, write := 10*time.Second, 30*time.Second

			for _, tt := range []struct {
				method  string
				timeout time.Duration
			}{
				{
					method:  http.MethodGet,
					timeout: read,
				},
				{
					method:  http.MethodHead,
					timeout: read,
				},
				{
					method:  http.MethodPost,
					timeout: write,
				},
				{
					method:  http.MethodPut,
					timeout: write,
				},
				{
					method:  http.MethodPatch,
					timeout: write,
				},
				{
					method:  http.MethodDelete,
					timeout: write,
				},
			} {
				t.Run(tt.method, func(t *testing.T) {
					var start time.Time

					tport := transport{
						timeout: Timeout{
							Read:  read,
							Write: write,
						},
						rest: restFunc(func(ctx context.Context, _ transports.Endpoint, _ any) error {
							d, ok := ctx.Deadline()
							require.True(t, ok, "context must have a deadline")

							timeout := d.Sub(start)
							require.InDelta(t, tt.timeout, timeout, float64(time.Millisecond))
							return nil
						}),
					}

					start = time.Now()
					err := tport.Do(t.Context(), &endpoint{method: tt.method}, nil)
					require.NoError(t, err, "request error")
				})
			}
		})
	})

	t.Run("grpc message", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			tport := transport{gRPC: new(fakeGRPC)}

			// Actual request is captured by message itself,
			// because unlike transports.Endpoint, each transports.RPC
			// provides the function that will execute this request.
			req := &message[proto.SearchRequest, proto.SearchReply]{
				req: &proto.SearchRequest{
					Limit:   1,
					Offset:  2,
					Autocut: 3,
				},
				out: &proto.SearchReply{},
			}

			var resp reply
			err := tport.Do(t.Context(), req, &resp)
			require.NoError(t, err, "transport error")

			require.EqualExportedValues(t, req.req, req.capture.req, "bad request")
			require.True(t, resp.used, "response not used")
		})

		t.Run("error", func(t *testing.T) {
			tport := transport{gRPC: new(fakeGRPC)}

			var resp reply
			req := &message[proto.SearchRequest, proto.SearchReply]{
				req: &proto.SearchRequest{},
				err: testkit.ErrWhaam,
			}
			err := tport.Do(t.Context(), req, &resp)

			require.ErrorIs(t, err, testkit.ErrWhaam, "gRPC transport error not propagated")
		})

		t.Run("timeout", func(t *testing.T) {
			read, write := 10*time.Second, 30*time.Second

			for _, tt := range []struct {
				name    string
				message any
				timeout time.Duration
			}{
				{
					name:    "search",
					message: &message[proto.SearchRequest, proto.SearchReply]{},
					timeout: read,
				},
				{
					name:    "aggregate",
					message: &message[proto.AggregateRequest, proto.AggregateReply]{},
					timeout: read,
				},
			} {
				t.Run(tt.name, func(t *testing.T) {
					var got context.Context

					tport := transport{
						timeout: Timeout{
							Read:  read,
							Write: write,
						},
						gRPC: gRPCFunc(func(ctx context.Context, rpc transports.RPC[proto.WeaviateClient]) error {
							got = ctx // Capture context for this request.
							return nil
						}),
					}

					start := time.Now()
					err := tport.Do(t.Context(), tt.message, new(reply))
					require.NoError(t, err, "request error")

					d, ok := got.Deadline()
					require.True(t, ok, "context must have a deadline")

					timeout := d.Sub(start)
					require.InDelta(t, tt.timeout, timeout, float64(time.Millisecond))
				})
			}
		})
	})
}

// restFunc implements REST transport as a function.
type restFunc func(context.Context, transports.Endpoint, any) error

func (f restFunc) Do(ctx context.Context, req transports.Endpoint, dest any) error {
	return f(ctx, req, dest)
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

type gRPCFunc func(ctx context.Context, rpc transports.RPC[proto.WeaviateClient]) error

func (f gRPCFunc) Do(ctx context.Context, rpc transports.RPC[proto.WeaviateClient]) error {
	return f(ctx, rpc)
}

// fakeGRPC calls rpc.Do with nil [proto.WeaviateClient].
// It's a dummy that should be used together with [message].
type fakeGRPC struct{}

func (*fakeGRPC) Do(ctx context.Context, rpc transports.RPC[proto.WeaviateClient]) error {
	return rpc.Do(ctx, nil)
}

// message implements [Message] for testing.
type message[In RequestMessage, Out ReplyMessage] struct {
	req *In   // Expected request.
	out *Out  // Response returned by fake methodFunc.
	err error // Error returned by fake methodFunc.

	// Use values passed to methodFunc in assertions.
	capture struct {
		ctx     context.Context
		wc      proto.WeaviateClient
		req     *In
		options []grpc.CallOption
	}
}

var _ Message[proto.SearchRequest, proto.SearchReply] = (*message[proto.SearchRequest, proto.SearchReply])(nil)

func (m *message[In, Out]) Method() MethodFunc[In, Out]  { return m.captureReq }
func (m *message[In, Out]) Body() MessageMarshaler[In]   { return m }
func (m *message[In, Out]) MarshalMessage() (*In, error) { return m.req, nil }

// capture is a fake messageFunc[In, Out] that captures the functions arguments.
func (m *message[In, Out]) captureReq(wc proto.WeaviateClient, ctx context.Context, req *In, opts ...grpc.CallOption) (*Out, error) {
	m.capture.ctx = ctx
	m.capture.wc = wc
	m.capture.req = req
	m.capture.options = opts

	return m.out, m.err
}

// reply implments [MessageUnmarshaler] for testing.
// Its used value is true once it's UnmarshalMessage() has been called.
type reply struct {
	used bool
	err  error
}

var _ MessageUnmarshaler[proto.SearchReply] = (*reply)(nil)

func (r *reply) UnmarshalMessage(*proto.SearchReply) error {
	r.used = true
	return r.err
}
