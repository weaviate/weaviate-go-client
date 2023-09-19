package batch

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v4/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/grpc"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/testenv"
)

func TestBatchCreate_gRPC_integration(t *testing.T) {
	err := testenv.SetupLocalWeaviate()
	if err != nil {
		require.Nil(t, err, "failed to start weaviate")
	}

	port, _, _ := testsuit.GetPortAndAuthPw()
	cfg := weaviate.Config{
		Host:   fmt.Sprintf("localhost:%v", port),
		Scheme: "http",
		GrpcConfig: grpc.Config{
			Enabled: true,
			Host:    "localhost:50051",
		},
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		require.Nil(t, err)
	}

	t.Run("gRPC batch import", func(t *testing.T) {
		err := client.Schema().AllDeleter().Do(context.Background())
		require.Nil(t, err)

		className := "AllProperties"
		testsuit.AllPropertiesSchemaCreate(t, client, className)
		objects := testsuit.AllPropertiesObjects(className)

		batchResultSlice, batchErrSlice := client.Batch().ObjectsBatcher().WithObjects(objects...).Do(context.Background())
		assert.Nil(t, batchErrSlice)
		assert.NotNil(t, batchResultSlice)
		assert.Equal(t, 3, len(batchResultSlice))
	})

	err = testenv.TearDownLocalWeaviate()
	if err != nil {
		require.Nil(t, err, "failed to tear down weaviate")
	}
}
