package cluster

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/testenv"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/verbosity"
)

func TestClusterNodes_integration(t *testing.T) {
	t.Skip("needs to be changed to use assert.EventuallyWithT() approach")
	const (
		expectedWeaviateVersion = "1.24.0-rc.0"
		expectedWeaviateGitHash = "b6e6111"
	)

	t.Run("up", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			fmt.Print(err.Error())
			t.Fail()
		}
	})

	t.Run("GET /nodes without data", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		nodesStatus, err := client.Cluster().NodesStatusGetter().
			WithOutput(verbosity.OutputVerbose).
			Do(context.Background())

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
		client := testsuit.CreateTestClient(false)
		testsuit.CreateTestSchemaAndData(t, client)
		defer testsuit.CleanUpWeaviate(t, client)

		nodesStatus, err := client.Cluster().NodesStatusGetter().
			WithOutput(verbosity.OutputVerbose).
			Do(context.Background())

		require.Nil(t, err)
		require.NotNil(t, nodesStatus)
		assert.Len(t, nodesStatus.Nodes, 1)

		nodeStatus := nodesStatus.Nodes[0]
		assert.NotEmpty(t, nodeStatus.Name)
		assert.Equal(t, expectedWeaviateVersion, nodeStatus.Version)
		assert.Equal(t, expectedWeaviateGitHash, nodeStatus.GitHash)
		assert.Equal(t, models.NodeStatusStatusHEALTHY, *nodeStatus.Status)
		require.NotNil(t, nodeStatus.Stats)
		assert.Equal(t, int64(9), nodeStatus.Stats.ObjectCount)
		assert.Equal(t, int64(3), nodeStatus.Stats.ShardCount)

		assert.Len(t, nodeStatus.Shards, 3)
		for _, shardStatus := range nodeStatus.Shards {
			assert.NotEmpty(t, shardStatus.Name)
			switch shardStatus.Class {
			case "Pizza":
				assert.Equal(t, int64(4), shardStatus.ObjectCount)
			case "Soup":
				assert.Equal(t, int64(2), shardStatus.ObjectCount)
			case "Risotto":
				assert.Equal(t, int64(3), shardStatus.ObjectCount)
			default:
				t.Fatalf("unexpected class name")
			}
		}

		// query only for one class
		nodesStatusSingleClass, err := client.Cluster().NodesStatusGetter().
			WithOutput(verbosity.OutputVerbose).
			WithClass("Pizza").Do(context.Background())
		require.Nil(t, err)
		assert.Len(t, nodesStatusSingleClass.Nodes, 1)
		nodeStatusSingleClass := nodesStatusSingleClass.Nodes[0]

		assert.Equal(t, int64(4), nodeStatusSingleClass.Stats.ObjectCount)
		assert.Equal(t, int64(1), nodeStatusSingleClass.Stats.ShardCount)
		// minimal output is default, which should not contain shards
		assert.Nil(t, nodeStatusSingleClass.Shards)
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			fmt.Print(err.Error())
			t.Fail()
		}
	})
}
