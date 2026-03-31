package transports_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/internal/transports"
	"google.golang.org/grpc"
)

func TestNewGRPC(t *testing.T) {
	grpc, err := transports.NewGRPC(transports.GRPCConfig[any]{
		Host: "example.com",
		Port: 12345,

		NewGRPCClient: func(channel grpc.ClientConnInterface) any {
			require.NotNil(t, channel, "grpc channel")
			if assert.IsType(t, (*grpc.ClientConn)(nil), channel) {
				conn := channel.(*grpc.ClientConn)
				require.Equal(t, "dns:///example.com:12345", conn.CanonicalTarget(), "canonical target")
			}
			return struct{}{}
		},
	})

	require.NoError(t, err)
	require.NotNil(t, grpc, "grpc transport")
}

func TestGRPC_Do(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		grpc, err := transports.NewGRPC(transports.GRPCConfig[any]{
			NewGRPCClient: func(channel grpc.ClientConnInterface) any {
				return 92
			},
		})
		require.NoError(t, err, "create grpc transport")
		require.NotNil(t, grpc, "grpc transport")

		require.NoError(t, grpc.Do(t.Context(), rpcFunc(func(_ context.Context, client any) error {
			assert.Equal(t, 92, client, "injected client")
			return nil
		})), "request error")
	})

	t.Run("with error", func(t *testing.T) {
		grpc, err := transports.NewGRPC(transports.GRPCConfig[any]{
			NewGRPCClient: func(channel grpc.ClientConnInterface) any {
				return 92
			},
		})
		require.NoError(t, err, "create grpc transport")
		require.NotNil(t, grpc, "grpc transport")

		require.ErrorIs(t, grpc.Do(t.Context(), rpcFunc(func(_ context.Context, client any) error {
			assert.Equal(t, 92, client, "injected client")
			return testkit.ErrWhaam
		})), testkit.ErrWhaam, "request error")
	})
}

type rpcFunc func(ctx context.Context, client any) error

var _ transports.RPC[any] = (*rpcFunc)(nil)

func (f rpcFunc) Do(ctx context.Context, client any) error {
	return f(ctx, client)
}
