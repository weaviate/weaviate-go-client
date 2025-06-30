package alias

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/testenv"
	"github.com/weaviate/weaviate/entities/models"
)

func TestAlias_integration(t *testing.T) {
	// 1. Create alias for non-existing class. - Fail
	// 1. Create alias for existing class. - Pass
	// 1. Create same alias for same existing class. - ???
	// 1. Create same alias for two existing class. - ??? shouldn't allow. alias is globally unique

	t.Run("up", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to setup weaviate: %s", err)
		}
	})

	t.Run("POST /alias", func(t *testing.T) {
		ctx := context.Background()
		client := testsuit.CreateTestClient(false)

		schemaClass := &models.Class{
			Class:               "Band",
			Description:         "Band that plays and produces music",
			Properties:          nil,
			VectorIndexType:     "hnsw",
			Vectorizer:          "text2vec-contextionary",
			InvertedIndexConfig: defaultInvertedIndexConfig,
			ModuleConfig:        defaultModuleConfig,
			ShardingConfig:      defaultShardingConfig,
			VectorIndexConfig:   defaultVectorIndexConfig,
			ReplicationConfig:   defaultReplicationConfig,
		}

		alias := &models.Alias{
			Alias: "Band-Alias",
			Class: schemaClass.Class,
		}

		err := client.Schema().ClassDeleter().WithClassName(schemaClass.Class).Do(ctx)
		require.NoError(t, err)

		err = client.Schema().ClassCreator().WithClass(schemaClass).Do(ctx)
		require.NoError(t, err)

		err = client.Alias().AliasCreator().WithAlias(alias).Do(ctx)
		require.NoError(t, err)

		errRm := client.Schema().AllDeleter().Do(context.Background())
		assert.Nil(t, errRm)
	})
}
