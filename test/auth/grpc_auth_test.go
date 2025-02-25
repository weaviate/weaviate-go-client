package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/auth"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/grpc"
	"github.com/weaviate/weaviate/entities/models"
)

func TestAuth_ApiKey_gRPC(t *testing.T) {
	ctx := context.TODO()
	cfg := weaviate.Config{
		Host:           fmt.Sprintf("127.0.0.1:%v", testsuit.WCSPort),
		Scheme:         "http",
		StartupTimeout: 60 * time.Second,
		AuthConfig:     auth.ApiKey{Value: "my-secret-key"},
		GrpcConfig:     &grpc.Config{Host: fmt.Sprintf("127.0.0.1:%v", testsuit.WCSGRPCPort)},
	}
	client, err := weaviate.NewClient(cfg)
	require.NoError(t, err)
	require.NotNil(t, client)

	err = client.Schema().AllDeleter().Do(ctx)
	require.NoError(t, err)

	id := "d3ae0320-a532-4e7b-9a89-6064263a182a"
	className := "Wine"

	t.Run("perform gRPC batch import (API Key is passed using gRPC header)", func(t *testing.T) {
		obj := &models.Object{
			ID:    strfmt.UUID(id),
			Class: className,
			Properties: map[string]interface{}{
				"name": "Classic Wine",
			},
		}
		resp, err := client.Batch().ObjectsBatcher().WithObjects(obj).Do(ctx)
		require.NoError(t, err)
		require.Len(t, resp, 1)
		assert.NotNil(t, resp[0])
	})

	t.Run("verify that object exists", func(t *testing.T) {
		exists, err := client.Data().Checker().WithClassName(className).WithID(id).Do(ctx)
		require.NoError(t, err)
		assert.True(t, exists)
	})
}
