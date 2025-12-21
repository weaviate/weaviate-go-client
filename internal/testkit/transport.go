package testkit

import (
	"context"
	"testing"

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

// NewResponder creates ResponderTransport populated with mock responses.
// All responses must be consumed -- this is verified on test cleanup.
func NewResponder[T any](t *testing.T, responses []Response[T]) *ResponderTransport[T] {
	tt := &ResponderTransport[T]{
		t:         t,
		responses: responses,
	}
	t.Cleanup(func() {
		require.True(t, tt.Done(), "requests were not fully consumed")
	})
	return tt
}

type Response[T any] struct {
	Value T
	Err   error
}

// ResponderTransport uses responses one by one until the slice is exhausted.
// It will not cycle through the responses, so calling the transport more times
// than the number of responses it has will fail the associated test.
type ResponderTransport[T any] struct {
	t         *testing.T
	responses []Response[T]
}

var _ internal.Transport = (*ResponderTransport[any])(nil)

// Do type-asserts dest to ensure it matches the expected type T,
// and consumes the next response, assigning it to dest.
// Returns Response.Err.
func (t *ResponderTransport[T]) Do(_ context.Context, req, dest any) error {
	t.t.Helper()

	require.IsType(t.t, (*T)(nil), dest)
	if len(t.responses) == 0 {
		require.Failf(t.t, "too many requests", "%#v is not expected", req)
	}

	var resp Response[T]
	resp, t.responses = t.responses[0], t.responses[1:] // pop front
	*dest.(*T) = resp.Value
	return resp.Err
}

func (t *ResponderTransport[T]) Done() bool { return len(t.responses) == 0 }
