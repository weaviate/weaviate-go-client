package testkit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal"
)

// TransportFunc implements internal.Transport to allow mocking it
// in tests with a single function.
type TransportFunc func(context.Context, any, any) error

var _ internal.Transport = (*TransportFunc)(nil)

// Do implements internal.Transport.
func (f TransportFunc) Do(ctx context.Context, req, dest any) error {
	return f(ctx, req, dest)
}

// NopTransport is a simple spy implementation that sets a boolean flag
// when it's Do method is called.
//
// Example:
//
//	func TestUsesTransport(t *testing.T) {
//		var nop testkit.NopTransport
//		c := data.NewClient(&nop, api.RequestDefaults{})
//
//		err := c.Delete(ctx, uuid.New())
//
//		assert.NoError(t, err)
//		assert.True(t, nop.Used(), "must call transport.Do")
//	}
type NopTransport struct{ used bool }

var _ internal.Transport = (*NopTransport)(nil)

// Do implements internal.Transport.
func (t *NopTransport) Do(context.Context, any, any) error {
	t.used = true
	return nil
}

// Used returns the value of the t.used flag and resets it to zero.
func (t *NopTransport) Used() bool {
	defer func() { t.used = false }()
	return t.used
}

// NewTransport creates a MockTransoprt populated with request/response stubs.
// All requests must be consumed -- this is verified on test cleanup.
func NewTransport[Req, Resp any](t *testing.T, stubs []Stub[Req, Resp]) *MockTransport[Req, Resp] {
	mock := &MockTransport[Req, Resp]{t: t, stubs: stubs}
	t.Cleanup(func() {
		require.True(t, mock.Done(), "requests were not fully consumed")
	})
	return mock
}

// Stub describes request-response pair for a MockTransport request
// along with any error that should be returned.
type Stub[Req, Resp any] struct {
	Request  *Req  // Expected request value. Leave unset to skip request check.
	Response Resp  // Response value. Will not be used if request dest is nil.
	Err      error // Error returned from Do. If set, Response and dest are ignored.
}

// MockTransport uses responses one by one until the slice is exhausted.
// It will not cycle through the responses, so calling the transport more times
// than the number of responses it has will fail the associated test.
type MockTransport[Req, Resp any] struct {
	t     *testing.T
	stubs []Stub[Req, Resp]
}

var _ internal.Transport = (*MockTransport[any, any])(nil)

// Do type-asserts dest to ensure it matches the expected type T,
// and consumes the next response, assigning it to dest.
// Returns Response.Err.
func (t *MockTransport[Req, Resp]) Do(_ context.Context, req, dest any) error {
	t.t.Helper()

	if len(t.stubs) == 0 {
		require.Failf(t.t, "too many requests", "%#v is not expected", req)
	}

	var stub Stub[Req, Resp]
	stub, t.stubs = t.stubs[0], t.stubs[1:] // pop front

	if stub.Request != nil && assert.IsType(t.t, (*Req)(nil), req, "bad request") {
		assert.Equal(t.t, stub.Request, req, "bad request")
	}

	if stub.Err != nil {
		return stub.Err
	}

	if dest != nil {
		require.IsType(t.t, (*Resp)(nil), dest, "bad dest")
		*dest.(*Resp) = stub.Response
	}
	return nil
}

// Done returns true if all requests have been consumed.
func (t *MockTransport[Req, Resp]) Done() bool { return len(t.stubs) == 0 }
