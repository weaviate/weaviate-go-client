//nolint:errcheck
package transports_test

import (
	"context"
	"net"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/internal/transports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
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
		gRPC, err := transports.NewGRPC(transports.GRPCConfig[any]{
			NewGRPCClient: func(channel grpc.ClientConnInterface) any {
				return 92
			},
		})
		require.NoError(t, err, "create grpc transport")
		require.NotNil(t, gRPC, "grpc transport")

		require.NoError(t, gRPC.Do(t.Context(), func(_ context.Context, client any) error {
			assert.Equal(t, 92, client, "injected client")
			return nil
		}), "request error")
	})

	t.Run("with error", func(t *testing.T) {
		grpc, err := transports.NewGRPC(transports.GRPCConfig[any]{
			NewGRPCClient: func(channel grpc.ClientConnInterface) any {
				return 92
			},
		})
		require.NoError(t, err, "create grpc transport")
		require.NotNil(t, grpc, "grpc transport")

		require.ErrorIs(t, grpc.Do(t.Context(), func(_ context.Context, client any) error {
			assert.Equal(t, 92, client, "injected client")
			return testkit.ErrWhaam
		}), testkit.ErrWhaam, "request error")
	})

	t.Run("default headers", func(t *testing.T) {
		// Arrange: start a local gRPC server and register a handler with assertions.
		ts := startTestService(t, func(_ any, ctx context.Context, _ func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
			md, ok := metadata.FromIncomingContext(ctx)
			assert.True(t, ok, "incoming context should contain metadata")
			assert.Subset(t, md, metadata.MD{"x-findme": {"foo"}}, "default headers not present in request metadata")
			return nil, nil
		})

		gRPC, err := transports.NewGRPC(transports.GRPCConfig[grpc.ClientConnInterface]{
			Host:          ts.Host(),
			Port:          ts.Port(),
			Header:        &metadata.MD{"X-FindMe": {"foo"}},
			NewGRPCClient: func(channel grpc.ClientConnInterface) grpc.ClientConnInterface { return channel },
		})
		require.NoError(t, err, "new grpc transport")

		// Act: our handled above will verify that the request included expected headers.
		gRPC.Do(t.Context(), func(ctx context.Context, client grpc.ClientConnInterface) error {
			var empty emptypb.Empty
			return client.Invoke(ctx, ts.MethodName(), nil, &empty)
		})
	})
}

type testService struct {
	lis  net.Listener
	srv  *grpc.Server
	host string
	port int
}

// startTestService starts a local TCP [net.Listener] and creates a [grpc.Server]
// using that listener. The mh handler can be used to make assertions about the
// request or control how the requeset is processed.
//
// All resources are freed via [testing.T.Cleanup] hook.
func startTestService(t *testing.T, mh grpc.MethodHandler) *testService {
	lis, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)
	t.Cleanup(func() { lis.Close() })

	addr := strings.Split(lis.Addr().String(), ":")
	port, _ := strconv.Atoi(addr[1])

	srv := grpc.NewServer()
	srv.RegisterService(&grpc.ServiceDesc{
		ServiceName: "testService",
		Methods: []grpc.MethodDesc{
			{MethodName: "Test", Handler: mh},
		},
	}, nil)

	go srv.Serve(lis)
	t.Cleanup(srv.Stop)

	return &testService{
		lis:  lis,
		srv:  srv,
		host: addr[0],
		port: port,
	}
}

func (ts *testService) Host() string       { return ts.host }
func (ts *testService) Port() int          { return ts.port }
func (ts *testService) MethodName() string { return "/testService/Test" }
