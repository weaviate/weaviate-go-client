package client

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v4/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/grpc"
)

func TestStartupTimeout_REST(t *testing.T) {
	_, grpcPort, authEnabled := testsuit.GetPortAndAuthPw()
	cfg := weaviate.Config{
		Host:   fmt.Sprintf("localhost:%v", testsuit.NoWeaviatePort),
		Scheme: "http",
		GrpcConfig: &grpc.Config{
			Host: fmt.Sprintf("localhost:%v", grpcPort),
		},
		StartupTimeout: 1000 * time.Millisecond,
	}
	if authEnabled {
		cfg.AuthConfig = auth.ApiKey{Value: "my-secret-key"}
	}
	_, err := weaviate.NewClient(cfg)
	require.NotNil(t, err)
}

func TestStartupTimeout_GRPC(t *testing.T) {
	port, _, authEnabled := testsuit.GetPortAndAuthPw()
	cfg := weaviate.Config{
		Host:   fmt.Sprintf("localhost:%v", port),
		Scheme: "http",
		GrpcConfig: &grpc.Config{
			Host: fmt.Sprintf("localhost:%v", testsuit.NoWeaviateGRPCPort),
		},
		StartupTimeout: 1000 * time.Millisecond,
	}
	if authEnabled {
		cfg.AuthConfig = auth.ApiKey{Value: "my-secret-key"}
	}
	_, err := weaviate.NewClient(cfg)
	require.NotNil(t, err)
}
