package api

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
	proto "github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/proto/v1"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/internal/transports"
	"google.golang.org/grpc"
)

func TestVersionedTransport_Do(t *testing.T) {
	t.Run("rest endpoint", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			var got *endpoint
			vt := versionedTransport{
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
			err := vt.Do(t.Context(), req, &resp)
			require.NoError(t, err, "transport error")

			require.Equal(t, req, got, "bad request")
			require.Equal(t, 5, resp, "bad response")
		})

		t.Run("error", func(t *testing.T) {
			vt := versionedTransport{
				rest: restFunc(func(_ context.Context, req transports.Endpoint, dest any) error {
					return testkit.ErrWhaam
				}),
			}

			err := vt.Do(t.Context(), &endpoint{}, nil)

			require.ErrorIs(t, err, testkit.ErrWhaam, "REST transport error not propagated")
		})
	})

	t.Run("grpc message", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			vt := versionedTransport{gRPC: new(fakeGRPC)}

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
			err := vt.Do(t.Context(), req, &resp)
			require.NoError(t, err, "transport error")

			require.EqualExportedValues(t, req.req, req.capture.req, "bad request")
			require.True(t, resp.used, "response not used")
		})

		t.Run("error", func(t *testing.T) {
			vt := versionedTransport{gRPC: new(fakeGRPC)}

			req := &message[proto.SearchRequest, proto.SearchReply]{
				req: &proto.SearchRequest{},
				err: testkit.ErrWhaam,
			}
			err := vt.Do(t.Context(), req, nil)

			require.ErrorIs(t, err, testkit.ErrWhaam, "gRPC transport error not propagated")
		})
	})

	t.Run("bad inputs", func(t *testing.T) {
		var vt versionedTransport
		require.Panics(t, func() { vt.Do(t.Context(), nil, nil) }, "Do accepted nil request")
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
		wc      proto.WeaviateClient
		req     *In
		options []grpc.CallOption
	}
}

var _ Message[proto.SearchRequest, proto.SearchReply] = (*message[proto.SearchRequest, proto.SearchReply])(nil)

func (m *message[In, Out]) Method() methodFunc[In, Out]  { return m.captureReq }
func (m *message[In, Out]) MarshalMessage() (*In, error) { return m.req, nil }

// capture is a fake messageFunc[In, Out] that captures the functions arguments.
func (m *message[In, Out]) captureReq(wc proto.WeaviateClient, _ context.Context, req *In, opts ...grpc.CallOption) (*Out, error) {
	m.capture.wc = wc
	m.capture.req = req
	m.capture.options = opts

	return m.out, m.err
}

// reply implments [MessageUnmarshaler] for testing.
// It's a "spy": 'used' is true once it's UnmarshalMessage() has been called.
type reply struct {
	used bool
	err  error
}

var _ MessageUnmarshaler[proto.SearchReply] = (*reply)(nil)

func (r *reply) UnmarshalMessage(*proto.SearchReply) error {
	r.used = true
	return r.err
}
