package client

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/auth"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/grpc"
)

func TestStartupTimeout_REST(t *testing.T) {
	_, grpcPort, authEnabled := testsuit.GetPortAndAuthPw()
	cfg := weaviate.Config{
		Host:   fmt.Sprintf("localhost:%v", testsuit.NoWeaviatePort),
		Scheme: "http",
		GrpcConfig: &grpc.Config{
			Host: fmt.Sprintf("localhost:%v", grpcPort),
		},
		StartupTimeout: time.Second,
	}
	if authEnabled {
		cfg.AuthConfig = auth.ApiKey{Value: "my-secret-key"}
	}
	_, err := weaviate.NewClient(cfg)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "weaviate did not start up in")
}

func TestStartupTimeout_GRPC(t *testing.T) {
	port, _, authEnabled := testsuit.GetPortAndAuthPw()
	cfg := weaviate.Config{
		Host:   fmt.Sprintf("localhost:%v", port),
		Scheme: "http",
		GrpcConfig: &grpc.Config{
			Host: fmt.Sprintf("localhost:%v", testsuit.NoWeaviateGRPCPort),
		},
		StartupTimeout: time.Second,
	}
	if authEnabled {
		cfg.AuthConfig = auth.ApiKey{Value: "my-secret-key"}
	}
	_, err := weaviate.NewClient(cfg)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "create grpc client")
}
