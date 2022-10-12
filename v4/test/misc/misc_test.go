package misc

import (
	"context"
	"fmt"
	"testing"

	"github.com/semi-technologies/weaviate-go-client/v4/test/testsuit"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/testenv"
	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMisc_integration(t *testing.T) {

	t.Run("up", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			fmt.Printf(err.Error())
			t.Fail()
		}
	})

	t.Run("GET /.well-known/ready", func(t *testing.T) {

		client := testsuit.CreateTestClient()
		isReady, err := client.Misc().ReadyChecker().Do(context.Background())

		assert.Nil(t, err)
		assert.True(t, isReady)
	})

	t.Run("GET /.well-known/live", func(t *testing.T) {

		client := testsuit.CreateTestClient()
		isLive, err := client.Misc().LiveChecker().Do(context.Background())

		assert.Nil(t, err)
		assert.True(t, isLive)
	})

	t.Run("GET /.well-known/openid-configuration", func(t *testing.T) {

		client := testsuit.CreateTestClient()
		openIDconfig, err := client.Misc().OpenIDConfigurationGetter().Do(context.Background())

		assert.Nil(t, err)
		assert.Nil(t, openIDconfig)
	})

	t.Run("GET /meta", func(t *testing.T) {
		client := testsuit.CreateTestClient()
		meta, err := client.Misc().MetaGetter().Do(context.Background())
		assert.Nil(t, err)
		assert.NotEmpty(t, meta.Version)
		modules, modulesOK := meta.Modules.(map[string]interface{})
		assert.True(t, modulesOK)
		text2vecContextionary := modules["text2vec-contextionary"]
		assert.NotEmpty(t, text2vecContextionary)
		text2vecContextionaryConfig, ok := text2vecContextionary.(map[string]interface{})
		assert.True(t, ok)
		text2vecContextionaryVersion := text2vecContextionaryConfig["version"]
		assert.NotEmpty(t, text2vecContextionaryVersion)
		text2vecContextionaryWordCount := text2vecContextionaryConfig["wordCount"]
		assert.NotEmpty(t, text2vecContextionaryWordCount)
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			fmt.Printf(err.Error())
			t.Fail()
		}
	})
}

func TestMiscNodes_integration(t *testing.T) {

	const (
		expectedWeaviateVersion = "1.15.4"
		expectedWeaviateGitHash = "d1ff58c"
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
		nodesStatus, err := client.Misc().NodesStatusGetter().Do(context.Background())

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

		nodesStatus, err := client.Misc().NodesStatusGetter().Do(context.Background())

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

func TestMisc_connection_error(t *testing.T) {
	t.Run("ready", func(t *testing.T) {
		cfg := weaviate.Config{
			Host:   "localhorst",
			Scheme: "http",
		}

		client := weaviate.New(cfg)
		isReady, err := client.Misc().ReadyChecker().Do(context.Background())

		assert.NotNil(t, err)
		assert.False(t, isReady)
	})

	t.Run("live", func(t *testing.T) {
		cfg := weaviate.Config{
			Host:   "localhorst",
			Scheme: "http",
		}

		client := weaviate.New(cfg)
		isReady, err := client.Misc().LiveChecker().Do(context.Background())

		assert.NotNil(t, err)
		assert.False(t, isReady)
	})
}
