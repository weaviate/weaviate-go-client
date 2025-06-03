package misc

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/testenv"
)

func TestMisc_integration(t *testing.T) {
	t.Run("up", func(t *testing.T) {
		err := testenv.SetupLocalWeaviateWaitForStartup(false)
		if err != nil {
			fmt.Print(err.Error())
			t.Fail()
		}
	})

	t.Run("GET /.well-known/ready", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		isReady, err := client.Misc().ReadyChecker().Do(context.Background())

		assert.Nil(t, err)
		assert.True(t, isReady)
	})

	t.Run("GET /.well-known/live", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		isLive, err := client.Misc().LiveChecker().Do(context.Background())

		assert.Nil(t, err)
		assert.True(t, isLive)
	})

	t.Run("GET /meta", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
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
			fmt.Print(err.Error())
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

		require.NotNil(t, err)
		assert.False(t, isReady)
		assert.Contains(t, err.Error(), "status code: -1, error: check the DerivedFromError field for more information: Get \"http://localhorst/v1/.well-known/ready\": dial tcp: lookup localhorst")
	})

	t.Run("live", func(t *testing.T) {
		cfg := weaviate.Config{
			Host:   "localhorst",
			Scheme: "http",
		}

		client := weaviate.New(cfg)
		isReady, err := client.Misc().LiveChecker().Do(context.Background())

		require.NotNil(t, err)
		assert.False(t, isReady)
		assert.Contains(t, err.Error(), "status code: -1, error: check the DerivedFromError field for more information: Get \"http://localhorst/v1/.well-known/live\": dial tcp: lookup localhorst")
	})
}
