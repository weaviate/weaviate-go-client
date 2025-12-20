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
