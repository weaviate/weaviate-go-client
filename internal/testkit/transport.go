package testkit

import (
	"context"

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
