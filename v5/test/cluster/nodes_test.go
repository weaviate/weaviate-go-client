package cluster

import (
	"context"
	"fmt"
	"testing"

	"github.com/semi-technologies/weaviate-go-client/v5/test/testsuit"
	"github.com/semi-technologies/weaviate-go-client/v5/weaviate/testenv"
	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClusterNodes_integration(t *testing.T) {

	const (
		expectedWeaviateVersion = "1.16.0"
		expectedWeaviateGitHash = "9e74add"
	)

	t.Run("up", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			fmt.Printf(err.Error())
			t.Fail()
		}
	})

	t.Run("GET /nodes without data", func(t *testing.T) {
		client := testsuit.CreateTestClient()
		nodesStatus, err := client.Cluster().NodesStatusGetter().Do(context.Background())

		require.Nil(t, err)
		require.NotNil(t, nodesStatus)
		assert.Len(t, nodesStatus.Nodes, 1)

		nodeStatus := nodesStatus.Nodes[0]
		assert.NotEmpty(t, nodeStatus.Name)
		assert.Equal(t, expectedWeaviateVersion, nodeStatus.Version)
		assert.Equal(t, expectedWeaviateGitHash, nodeStatus.GitHash)
		assert.Equal(t, models.NodeStatusStatusHEALTHY, *nodeStatus.Status)
		assert.Len(t, nodeStatus.Shards, 0)
		require.NotNil(t, nodeStatus.Stats)
		assert.Equal(t, int64(0), nodeStatus.Stats.ObjectCount)
		assert.Equal(t, int64(0), nodeStatus.Stats.ShardCount)
	})

	t.Run("GET /nodes with data", func(t *testing.T) {
		client := testsuit.CreateTestClient()
		testsuit.CreateTestSchemaAndData(t, client)
		defer testsuit.CleanUpWeaviate(t, client)

		nodesStatus, err := client.Cluster().NodesStatusGetter().Do(context.Background())

		require.Nil(t, err)
		require.NotNil(t, nodesStatus)
		assert.Len(t, nodesStatus.Nodes, 1)

		nodeStatus := nodesStatus.Nodes[0]
		assert.NotEmpty(t, nodeStatus.Name)
		assert.Equal(t, expectedWeaviateVersion, nodeStatus.Version)
		assert.Equal(t, expectedWeaviateGitHash, nodeStatus.GitHash)
		assert.Equal(t, models.NodeStatusStatusHEALTHY, *nodeStatus.Status)
		require.NotNil(t, nodeStatus.Stats)
		assert.Equal(t, int64(6), nodeStatus.Stats.ObjectCount)
		assert.Equal(t, int64(2), nodeStatus.Stats.ShardCount)

		assert.Len(t, nodeStatus.Shards, 2)
		for _, shardStatus := range nodeStatus.Shards {
			assert.NotEmpty(t, shardStatus.Name)
			switch shardStatus.Class {
			case "Pizza":
				assert.Equal(t, int64(4), shardStatus.ObjectCount)
			case "Soup":
				assert.Equal(t, int64(2), shardStatus.ObjectCount)
			default:
				t.Fatalf("unexpected class name")
			}
		}
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			fmt.Printf(err.Error())
			t.Fail()
		}
	})
}
