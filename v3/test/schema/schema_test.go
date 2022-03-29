package schema

import (
	"context"
	"testing"

	"github.com/semi-technologies/weaviate-go-client/v3/test/testsuit"
	"github.com/semi-technologies/weaviate-go-client/v3/weaviate/testenv"
	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/semi-technologies/weaviate/entities/storagestate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSchema_integration(t *testing.T) {

	t.Run("up", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to setup weaviate: %s", err)
		}
	})

	t.Run("POST /schema", func(t *testing.T) {
		client := testsuit.CreateTestClient()

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
		}

		err := client.Schema().ClassCreator().WithClass(schemaClass).Do(context.Background())
		assert.Nil(t, err)

		loadedSchema, getErr := client.Schema().Getter().Do(context.Background())
		assert.Nil(t, getErr)
		assert.Equal(t, 1, len(loadedSchema.Classes))
		assert.Equal(t, schemaClass, loadedSchema.Classes[0])
		assert.Equal(t, schemaClass.Class, loadedSchema.Classes[0].Class)
		assert.Equal(t, schemaClass.Description, loadedSchema.Classes[0].Description)

		// Clean up classes
		errRm := client.Schema().AllDeleter().Do(context.Background())
		assert.Nil(t, errRm)
	})

	t.Run("POST /schema", func(t *testing.T) {
		client := testsuit.CreateTestClient()

		schemaClass := &models.Class{
			Class:               "Run",
			Description:         "Running from the fuzz",
			VectorIndexType:     "hnsw",
			Vectorizer:          "text2vec-contextionary",
			InvertedIndexConfig: defaultInvertedIndexConfig,
			ModuleConfig:        defaultModuleConfig,
			ShardingConfig:      defaultShardingConfig,
			VectorIndexConfig:   defaultVectorIndexConfig,
		}

		err := client.Schema().ClassCreator().WithClass(schemaClass).Do(context.Background())
		assert.Nil(t, err)

		loadedSchema, getErr := client.Schema().Getter().Do(context.Background())
		assert.Nil(t, getErr)
		assert.Equal(t, 1, len(loadedSchema.Classes))
		assert.Equal(t, schemaClass, loadedSchema.Classes[0])

		// Clean up classes
		errRm := client.Schema().AllDeleter().Do(context.Background())
		assert.Nil(t, errRm)
	})

	t.Run("Delete /schema/{type}", func(t *testing.T) {
		client := testsuit.CreateTestClient()

		schemaClassThing := &models.Class{
			Class:       "Pizza",
			Description: "A delicious religion like food and arguably the best export of Italy.",
		}
		schemaClassAction := &models.Class{
			Class:       "ChickenSoup",
			Description: "A soup made in part out of chicken, not for chicken.",
		}

		errT := client.Schema().ClassCreator().WithClass(schemaClassThing).Do(context.Background())
		assert.Nil(t, errT)
		errA := client.Schema().ClassCreator().WithClass(schemaClassAction).Do(context.Background())
		assert.Nil(t, errA)
		loadedSchema, getErr := client.Schema().Getter().Do(context.Background())
		assert.Nil(t, getErr)
		assert.Equal(t, loadedSchema.Classes[0].Class, schemaClassThing.Class)
		assert.Equal(t, loadedSchema.Classes[1].Class, schemaClassAction.Class)
		assert.Equal(t, 2, len(loadedSchema.Classes), "There are classes in the schema that are not part of this test")

		errRm1 := client.Schema().ClassDeleter().WithClassName(schemaClassThing.Class).Do(context.Background())
		errRm2 := client.Schema().ClassDeleter().WithClassName(schemaClassAction.Class).Do(context.Background())
		assert.Nil(t, errRm1)
		assert.Nil(t, errRm2)

		loadedSchema, getErr = client.Schema().Getter().Do(context.Background())
		assert.Nil(t, getErr)
		assert.Equal(t, 0, len(loadedSchema.Classes))
	})

	t.Run("Delete All schema", func(t *testing.T) {
		client := testsuit.CreateTestClient()

		schemaClassThing := &models.Class{
			Class:       "Pizza",
			Description: "A delicious religion like food and arguably the best export of Italy.",
		}
		schemaClassAction := &models.Class{
			Class:       "ChickenSoup",
			Description: "A soup made in part out of chicken, not for chicken.",
		}

		errT := client.Schema().ClassCreator().WithClass(schemaClassThing).Do(context.Background())
		assert.Nil(t, errT)
		errA := client.Schema().ClassCreator().WithClass(schemaClassAction).Do(context.Background())
		assert.Nil(t, errA)

		errRm1 := client.Schema().AllDeleter().Do(context.Background())
		assert.Nil(t, errRm1)

		loadedSchema, getErr := client.Schema().Getter().Do(context.Background())
		assert.Nil(t, getErr)
		assert.Equal(t, 0, len(loadedSchema.Classes))
	})

	t.Run("POST /schema/{type}/{className}/properties", func(t *testing.T) {
		client := testsuit.CreateTestClient()

		schemaClassThing := &models.Class{
			Class:       "Pizza",
			Description: "A delicious religion like food and arguably the best export of Italy.",
		}
		schemaClassAction := &models.Class{
			Class:       "ChickenSoup",
			Description: "A soup made in part out of chicken, not for chicken.",
		}

		errT := client.Schema().ClassCreator().WithClass(schemaClassThing).Do(context.Background())
		assert.Nil(t, errT)
		errA := client.Schema().ClassCreator().WithClass(schemaClassAction).Do(context.Background())
		assert.Nil(t, errA)

		newProperty := &models.Property{
			DataType:    []string{"string"},
			Description: "name",
			Name:        "name",
		}

		propErrT := client.Schema().PropertyCreator().WithClassName("Pizza").WithProperty(newProperty).Do(context.Background())
		assert.Nil(t, propErrT)
		propErrA := client.Schema().PropertyCreator().WithClassName("ChickenSoup").WithProperty(newProperty).Do(context.Background())
		assert.Nil(t, propErrA)

		loadedSchema, getErr := client.Schema().Getter().Do(context.Background())
		assert.Nil(t, getErr)
		assert.Equal(t, 2, len(loadedSchema.Classes))
		assert.Equal(t, "name", loadedSchema.Classes[0].Properties[0].Name)
		assert.Equal(t, "name", loadedSchema.Classes[1].Properties[0].Name)

		// Clean up classes
		errRm := client.Schema().AllDeleter().Do(context.Background())
		assert.Nil(t, errRm)
	})

	t.Run("GET /schema/{className}", func(t *testing.T) {
		client := testsuit.CreateTestClient()

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
		}

		err := client.Schema().ClassCreator().WithClass(schemaClass).Do(context.Background())
		assert.Nil(t, err)

		bandClass, getErr := client.Schema().ClassGetter().WithClassName("Band").Do(context.Background())
		assert.Nil(t, getErr)
		assert.NotNil(t, bandClass)
		assert.Equal(t, "Band", bandClass.Class)
		assert.Equal(t, "Band that plays and produces music", bandClass.Description)
		assert.Equal(t, "hnsw", bandClass.VectorIndexType)
		assert.Equal(t, "text2vec-contextionary", bandClass.Vectorizer)
		assert.Nil(t, bandClass.Properties)
		assert.NotNil(t, bandClass.ModuleConfig)
		assert.NotNil(t, bandClass.VectorIndexConfig)

		nonExistantClass, getErr := client.Schema().ClassGetter().WithClassName("NonExistantClass").Do(context.Background())
		assert.NotNil(t, getErr)
		assert.Nil(t, nonExistantClass)

		// Clean up classes
		errRm := client.Schema().AllDeleter().Do(context.Background())
		assert.Nil(t, errRm)
	})

	t.Run("GET /schema/{className}/shards", func(t *testing.T) {
		client := testsuit.CreateTestClient()

		class := &models.Class{
			Class:               "Article",
			Description:         "Archived news article",
			Properties:          nil,
			VectorIndexType:     "hnsw",
			Vectorizer:          "text2vec-contextionary",
			InvertedIndexConfig: defaultInvertedIndexConfig,
			ModuleConfig:        defaultModuleConfig,
			ShardingConfig:      defaultShardingConfig,
			VectorIndexConfig:   defaultVectorIndexConfig,
		}

		err := client.Schema().ClassCreator().WithClass(class).Do(context.Background())
		assert.Nil(t, err)

		shards, err := client.Schema().
			ShardsGetter().
			WithClassName(class.Class).
			Do(context.Background())
		assert.Nil(t, err)
		assert.NotEmpty(t, shards)
		assert.Equal(t, storagestate.StatusReady.String(), shards[0].Status)

		// clean up classes
		errRm := client.Schema().AllDeleter().Do(context.Background())
		assert.Nil(t, errRm)
	})

	t.Run("PUT /schema/{className}/shards/{shardName}", func(t *testing.T) {
		client := testsuit.CreateTestClient()

		class := &models.Class{
			Class:               "ClassOne",
			Properties:          nil,
			VectorIndexType:     "hnsw",
			Vectorizer:          "text2vec-contextionary",
			InvertedIndexConfig: defaultInvertedIndexConfig,
			ModuleConfig:        defaultModuleConfig,
			ShardingConfig:      defaultShardingConfig,
			VectorIndexConfig:   defaultVectorIndexConfig,
		}

		// create test class
		err := client.Schema().ClassCreator().WithClass(class).Do(context.Background())
		assert.Nil(t, err)

		// ensure shard exists for test class with correct status
		shards, err := client.Schema().ShardsGetter().WithClassName(class.Class).Do(context.Background())
		assert.Nil(t, err)
		require.NotEmpty(t, shards)
		assert.Equal(t, storagestate.StatusReady.String(), shards[0].Status)

		// set shard to readonly
		status, err := client.Schema().
			ShardUpdater().
			WithClassName(class.Class).
			WithShardName(shards[0].Name).
			WithStatus("READONLY").
			Do(context.Background())
		assert.Nil(t, err)
		require.NotNil(t, status)
		assert.Equal(t, storagestate.StatusReadOnly.String(), status.Status)

		// set shard to ready
		status, err = client.Schema().
			ShardUpdater().
			WithClassName(class.Class).
			WithShardName(shards[0].Name).
			WithStatus("READY").
			Do(context.Background())
		assert.Nil(t, err)
		require.NotNil(t, status)
		assert.Equal(t, storagestate.StatusReady.String(), status.Status)

		// clean up classes
		errRm := client.Schema().AllDeleter().Do(context.Background())
		assert.Nil(t, errRm)
	})

	t.Run("Update all shards convenience method", func(t *testing.T) {
		client := testsuit.CreateTestClient()

		shardCount := 3

		class := &models.Class{
			Class:               "ClassOne",
			Properties:          nil,
			VectorIndexType:     "hnsw",
			Vectorizer:          "text2vec-contextionary",
			InvertedIndexConfig: defaultInvertedIndexConfig,
			ModuleConfig:        defaultModuleConfig,
			ShardingConfig: map[string]interface{}{
				"actualCount":         float64(shardCount),
				"actualVirtualCount":  float64(128),
				"desiredCount":        float64(shardCount),
				"desiredVirtualCount": float64(128),
				"function":            "murmur3",
				"key":                 "_id",
				"strategy":            "hash",
				"virtualPerPhysical":  float64(128),
			},
			VectorIndexConfig: defaultVectorIndexConfig,
		}

		// create test class
		err := client.Schema().
			ClassCreator().
			WithClass(class).
			Do(context.Background())
		assert.Nil(t, err)

		resp, err := client.Schema().ShardsGetter().WithClassName(class.Class).Do(context.Background())
		assert.Nil(t, err)
		require.NotEmpty(t, resp)
		assert.Equal(t, shardCount, len(resp))

		// set all shards to readonly
		shards, err := client.Schema().
			ShardsUpdater().
			WithClassName(class.Class).
			WithStatus("READONLY").
			Do(context.Background())
		assert.Nil(t, err)
		require.NotEmpty(t, resp)
		assert.Equal(t, shardCount, len(resp))

		for _, shard := range shards {
			assert.Equal(t, storagestate.StatusReadOnly.String(), shard.Status)
		}

		// set all shards to ready
		shards, err = client.Schema().
			ShardsUpdater().
			WithClassName(class.Class).
			WithStatus("READY").
			Do(context.Background())
		assert.Nil(t, err)
		require.NotEmpty(t, resp)
		assert.Equal(t, shardCount, len(resp))

		for _, shard := range shards {
			assert.Equal(t, storagestate.StatusReady.String(), shard.Status)
		}

		// clean up classes
		errRm := client.Schema().AllDeleter().Do(context.Background())
		assert.Nil(t, errRm)
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to tear down weaviate: %s", err)
		}
	})
}

func TestSchema_errors(t *testing.T) {

	t.Run("Run Do without setting a class", func(t *testing.T) {
		client := testsuit.CreateTestClient()

		err := client.Schema().ClassCreator().Do(context.Background())
		assert.NotNil(t, err)
	})

}
