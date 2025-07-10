package schema

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/fault"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/testenv"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/schema"
	"github.com/weaviate/weaviate/entities/storagestate"
)

func TestSchema_integration(t *testing.T) {
	t.Run("up", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to setup weaviate: %s", err)
		}
	})

	t.Run("POST /schema Test1", func(t *testing.T) {
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

		client.Schema().ClassDeleter().WithClassName(schemaClass.Class).Do(context.Background())
		err := client.Schema().ClassCreator().WithClass(schemaClass).Do(context.Background())
		assert.Nil(t, err)

		loadedSchema, getErr := client.Schema().Getter().Do(context.Background())
		assert.Nil(t, getErr)
		assert.Equal(t, 1, len(loadedSchema.Classes))

		schemaClass.MultiTenancyConfig = defaultMultiTenancyConfig
		assert.EqualValues(t, schemaClass, loadedSchema.Classes[0])

		// Clean up classes
		errRm := client.Schema().AllDeleter().Do(context.Background())
		assert.Nil(t, errRm)
	})

	t.Run("POST /schema - Test2", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)

		schemaClass := &models.Class{
			Class:               "Run",
			Description:         "Running from the fuzz",
			VectorIndexType:     "hnsw",
			Vectorizer:          "text2vec-contextionary",
			InvertedIndexConfig: defaultInvertedIndexConfig,
			ModuleConfig:        defaultModuleConfig,
			ShardingConfig:      defaultShardingConfig,
			VectorIndexConfig:   defaultVectorIndexConfig,
			ReplicationConfig:   defaultReplicationConfig,
		}

		err := client.Schema().ClassCreator().WithClass(schemaClass).Do(context.Background())
		assert.Nil(t, err)

		loadedSchema, getErr := client.Schema().Getter().Do(context.Background())
		assert.Nil(t, getErr)
		assert.Equal(t, 1, len(loadedSchema.Classes))

		schemaClass.MultiTenancyConfig = defaultMultiTenancyConfig
		assert.Equal(t, schemaClass, loadedSchema.Classes[0])

		// Clean up classes
		errRm := client.Schema().AllDeleter().Do(context.Background())
		assert.Nil(t, errRm)
	})

	t.Run("Delete /schema/{type}", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)

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
		assert.Equal(t, 2, len(loadedSchema.Classes), "There are classes in the schema that are not part of this test")

		class, err := client.Schema().ClassGetter().WithClassName(schemaClassAction.Class).Do(context.Background())
		require.NoError(t, err)
		assert.NotNil(t, class)
		assert.Equal(t, class.Class, schemaClassAction.Class)

		class, err = client.Schema().ClassGetter().WithClassName(schemaClassThing.Class).Do(context.Background())
		require.NoError(t, err)
		assert.NotNil(t, class)
		assert.Equal(t, class.Class, schemaClassThing.Class)

		errRm1 := client.Schema().ClassDeleter().WithClassName(schemaClassThing.Class).Do(context.Background())
		errRm2 := client.Schema().ClassDeleter().WithClassName(schemaClassAction.Class).Do(context.Background())
		assert.Nil(t, errRm1)
		assert.Nil(t, errRm2)

		loadedSchema, getErr = client.Schema().Getter().Do(context.Background())
		assert.Nil(t, getErr)
		assert.Equal(t, 0, len(loadedSchema.Classes))
	})

	t.Run("PUT /schema", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)

		schemaClass := &models.Class{
			Class:               "Run",
			Description:         "Running from the fuzz",
			VectorIndexType:     "hnsw",
			InvertedIndexConfig: defaultInvertedIndexConfig,
			ModuleConfig:        defaultModuleConfig,
			ShardingConfig:      defaultShardingConfig,
			VectorIndexConfig:   defaultVectorIndexConfig,
			ReplicationConfig:   defaultReplicationConfig,
		}

		err := client.Schema().ClassCreator().WithClass(schemaClass).Do(context.Background())
		assert.Nil(t, err)

		// Now update the class
		err = client.Schema().ClassUpdater().WithClass(&models.Class{
			Class: schemaClass.Class,
			VectorIndexConfig: map[string]interface{}{
				"ef": 42,
			},
		}).Do(context.Background())

		assert.Nil(t, err)

		loadedSchema, getErr := client.Schema().Getter().Do(context.Background())
		assert.Nil(t, getErr)
		assert.Equal(t, 1, len(loadedSchema.Classes))

		vectorIndexConfig := loadedSchema.Classes[0].VectorIndexConfig.(map[string]interface{})
		assert.Equal(t, float64(42), vectorIndexConfig["ef"].(float64))

		// With Class name missing
		err = client.Schema().ClassUpdater().WithClass(&models.Class{
			VectorIndexConfig: map[string]interface{}{
				"ef": 42,
			},
		}).Do(context.Background())
		assert.Error(t, err)

		// Without WithClass
		err = client.Schema().ClassUpdater().Do(context.Background())
		assert.Error(t, err)

		// Clean up classes
		errRm := client.Schema().AllDeleter().Do(context.Background())
		assert.Nil(t, errRm)
	})

	t.Run("Delete All schema", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)

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
		client := testsuit.CreateTestClient(false)

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
			DataType:    []string{"text"},
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
		assert.Equal(t, models.PropertyTokenizationWord, loadedSchema.Classes[0].Properties[0].Tokenization)
		assert.Equal(t, "name", loadedSchema.Classes[1].Properties[0].Name)
		assert.Equal(t, models.PropertyTokenizationWord, loadedSchema.Classes[1].Properties[0].Tokenization)

		// Clean up classes
		errRm := client.Schema().AllDeleter().Do(context.Background())
		assert.Nil(t, errRm)
	})

	t.Run("PUT /schema/{className} to add vectors", func(t *testing.T) {
		ctx := context.Background()
		client := testsuit.CreateTestClient(false)
		err := client.Schema().ClassCreator().WithClass(&models.Class{
			Class: "PizzaAddVector",
			VectorConfig: map[string]models.VectorConfig{
				"default": {
					VectorIndexType: "hnsw",
					Vectorizer: map[string]interface{}{
						"none": map[string]interface{}{},
					},
				},
			},
		}).Do(ctx)
		require.NoError(t, err, "create PizzaAddVector collection")

		err = client.Schema().VectorAdder().
			WithClassName("PizzaAddVector").
			WithVectors(map[string]models.VectorConfig{
				"vector-a": {
					VectorIndexType: "hnsw",
					Vectorizer: map[string]interface{}{
						"none": map[string]interface{}{},
					},
				},
				"vector-b": {
					VectorIndexType: "hnsw",
					Vectorizer: map[string]interface{}{
						"none": map[string]interface{}{},
					},
				},
			}).
			Do(ctx)
		require.NoError(t, err, "add vector-a and vector-b")

		pizza, err := client.Schema().ClassGetter().WithClassName("PizzaAddVector").Do(ctx)
		require.NoError(t, err, "get PizzaAddVector collection")

		require.Contains(t, pizza.VectorConfig, "default")
		require.Contains(t, pizza.VectorConfig, "vector-a")
		require.Contains(t, pizza.VectorConfig, "vector-b")
	})

	t.Run("GET /schema/{className}", func(t *testing.T) {
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

		nonExistantClass, getErr := client.Schema().ClassGetter().WithClassName("NonExistentClass").Do(context.Background())
		assert.NotNil(t, getErr)
		assert.Nil(t, nonExistantClass)

		// Clean up classes
		errRm := client.Schema().AllDeleter().Do(context.Background())
		assert.Nil(t, errRm)
	})

	t.Run("CHECK /schema/{className}", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)

		schemaClass := &models.Class{
			Class: "Band",
		}

		err := client.Schema().ClassCreator().WithClass(schemaClass).Do(context.Background())
		assert.Nil(t, err)

		ok, err := client.Schema().ClassExistenceChecker().WithClassName("Band").Do(context.Background())
		assert.Nil(t, err)
		assert.True(t, ok)

		ok, err = client.Schema().ClassExistenceChecker().WithClassName("NonExistentClass").Do(context.Background())
		assert.Nil(t, err)
		assert.False(t, ok)

		// Clean up classes
		errRm := client.Schema().AllDeleter().Do(context.Background())
		assert.Nil(t, errRm)
	})

	t.Run("GET /schema/{className}/shards", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)

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
		client := testsuit.CreateTestClient(false)

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
		client := testsuit.CreateTestClient(false)

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

	t.Run("Create class with BM25 config", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)

		schemaClass := &models.Class{
			Class:             "Run",
			Description:       "Running from the fuzz",
			VectorIndexType:   "hnsw",
			Vectorizer:        "text2vec-contextionary",
			ModuleConfig:      defaultModuleConfig,
			ShardingConfig:    defaultShardingConfig,
			VectorIndexConfig: defaultVectorIndexConfig,
			InvertedIndexConfig: &models.InvertedIndexConfig{
				Bm25: &models.BM25Config{
					K1: 1.11,
					B:  0.66,
				},
			},
		}

		err := client.
			Schema().
			ClassCreator().
			WithClass(schemaClass).
			Do(context.Background())
		require.Nil(t, err)

		loadedSchema, getErr := client.Schema().Getter().Do(context.Background())
		require.Nil(t, getErr)
		require.Equal(t, 1, len(loadedSchema.Classes))
		assert.Equal(t, schemaClass.InvertedIndexConfig.Bm25,
			loadedSchema.Classes[0].InvertedIndexConfig.Bm25)

		// Clean up classes
		errRm := client.Schema().AllDeleter().Do(context.Background())
		require.Nil(t, errRm)
	})

	t.Run("Create class with Stopword config", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)

		schemaClass := &models.Class{
			Class:             "SpaceThings",
			Description:       "Things about the universe",
			VectorIndexType:   "hnsw",
			Vectorizer:        "text2vec-contextionary",
			ModuleConfig:      defaultModuleConfig,
			ShardingConfig:    defaultShardingConfig,
			VectorIndexConfig: defaultVectorIndexConfig,
			InvertedIndexConfig: &models.InvertedIndexConfig{
				Stopwords: &models.StopwordConfig{
					Preset:    "en",
					Additions: []string{"star", "nebula"},
					Removals:  []string{"a", "the"},
				},
			},
		}

		err := client.
			Schema().
			ClassCreator().
			WithClass(schemaClass).
			Do(context.Background())
		require.Nil(t, err)

		loadedSchema, getErr := client.Schema().Getter().Do(context.Background())
		require.Nil(t, getErr)
		require.Equal(t, 1, len(loadedSchema.Classes))
		assert.Equal(t, schemaClass.InvertedIndexConfig.Stopwords,
			loadedSchema.Classes[0].InvertedIndexConfig.Stopwords)

		// Clean up classes
		errRm := client.Schema().AllDeleter().Do(context.Background())
		require.Nil(t, errRm)
	})

	t.Run("Create class with BM25 and Stopword config", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)

		schemaClass := &models.Class{
			Class:               "SpaceThings",
			Description:         "Things about the universe",
			VectorIndexType:     "hnsw",
			Vectorizer:          "text2vec-contextionary",
			ModuleConfig:        defaultModuleConfig,
			ShardingConfig:      defaultShardingConfig,
			VectorIndexConfig:   defaultVectorIndexConfig,
			InvertedIndexConfig: defaultInvertedIndexConfig,
		}
		schemaClass.InvertedIndexConfig.Bm25 = &models.BM25Config{
			K1: 1.777,
			B:  0.777,
		}
		schemaClass.InvertedIndexConfig.Stopwords = &models.StopwordConfig{
			Preset:    "en",
			Additions: []string{"star", "nebula"},
			Removals:  []string{"a", "the"},
		}

		err := client.
			Schema().
			ClassCreator().
			WithClass(schemaClass).
			Do(context.Background())
		require.Nil(t, err)

		loadedSchema, getErr := client.Schema().Getter().Do(context.Background())
		require.Nil(t, getErr)
		require.Equal(t, 1, len(loadedSchema.Classes))
		assert.Equal(t, schemaClass.InvertedIndexConfig.Bm25,
			loadedSchema.Classes[0].InvertedIndexConfig.Bm25)
		assert.Equal(t, schemaClass.InvertedIndexConfig.Stopwords,
			loadedSchema.Classes[0].InvertedIndexConfig.Stopwords)

		// Clean up classes
		errRm := client.Schema().AllDeleter().Do(context.Background())
		require.Nil(t, errRm)
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to tear down weaviate: %s", err)
		}
	})
}

func TestReplication(t *testing.T) {
	t.Skip("skipping replication tests for now")
	t.Run("up", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to setup weaviate: %s", err)
		}
	})
	t.Run("Create class with implicit replication config", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		className := "ReplicationClass"

		schemaClass := &models.Class{
			Class: className,
		}

		err := client.
			Schema().
			ClassCreator().
			WithClass(schemaClass).
			Do(context.Background())
		require.Nil(t, err)

		loadedClass, getErr := client.Schema().ClassGetter().WithClassName(className).Do(context.Background())
		require.Nil(t, getErr)
		require.NotNil(t, loadedClass.ReplicationConfig)
		assert.Equal(t, int64(1), loadedClass.ReplicationConfig.Factor)

		// Clean up classes
		errRm := client.Schema().AllDeleter().Do(context.Background())
		require.Nil(t, errRm)
	})

	t.Run("Create class with explicit replication config", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		className := "ReplicationClass"

		schemaClass := &models.Class{
			Class: className,
			ReplicationConfig: &models.ReplicationConfig{
				Factor: 3,
			},
		}

		err := client.
			Schema().
			ClassCreator().
			WithClass(schemaClass).
			Do(context.Background())
		require.Nil(t, err)

		loadedClass, getErr := client.Schema().ClassGetter().WithClassName(className).Do(context.Background())
		require.Nil(t, getErr)
		require.NotNil(t, loadedClass.ReplicationConfig)
		assert.Equal(t, int64(3), loadedClass.ReplicationConfig.Factor)

		// Clean up classes
		errRm := client.Schema().AllDeleter().Do(context.Background())
		require.Nil(t, errRm)
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to tear down weaviate: %s", err)
		}
	})
}

func TestSchema_errors(t *testing.T) {
	t.Run("up", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to setup weaviate: %s", err)
		}
	})

	t.Run("Run Do without setting a class", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)

		err := client.Schema().ClassCreator().Do(context.Background())
		assert.NotNil(t, err)
	})

	t.Run("Fail to add property having not supported tokenization", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)

		pizzaClass := &models.Class{
			Class:       "Pizza",
			Description: "A delicious religion like food and arguably the best export of Italy.",
		}

		err := client.Schema().ClassCreator().WithClass(pizzaClass).Do(context.Background())
		assert.Nil(t, err)

		notExistingTokenizationProperty := &models.Property{
			DataType:     []string{"text"},
			Description:  "name",
			Name:         "name",
			Tokenization: "not-existing",
		}

		err = client.Schema().PropertyCreator().WithClassName("Pizza").
			WithProperty(notExistingTokenizationProperty).Do(context.Background())
		assert.ErrorContains(t, err, `status code: 422, error: {"code":606,"message":"tokenization in body should be one of `)

		notSupportedTokenizationProperty1 := &models.Property{
			DataType:     []string{"text"},
			Description:  "description",
			Name:         "description",
			Tokenization: models.PropertyTokenizationField,
		}

		err = client.Schema().PropertyCreator().WithClassName("Pizza").
			WithProperty(notSupportedTokenizationProperty1).Do(context.Background())
		// Since v1.19 tokenization field is supported for data type text
		assert.Nil(t, err)

		notSupportedTokenizationProperty2 := &models.Property{
			DataType:     []string{"int[]"},
			Description:  "calories",
			Name:         "calories",
			Tokenization: models.PropertyTokenizationWord,
		}

		err = client.Schema().PropertyCreator().WithClassName("Pizza").
			WithProperty(notSupportedTokenizationProperty2).Do(context.Background())
		assert.EqualError(t, err, "status code: 422, error: {\"error\":[{\"message\":\"tokenization is not allowed for data type 'int[]'\"}]}\n")
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to tear down weaviate: %s", err)
		}
	})
}

func TestSchema_MultiTenancyConfig(t *testing.T) {
	cleanup := func() {
		client := testsuit.CreateTestClient(false)
		err := client.Schema().AllDeleter().Do(context.Background())
		require.Nil(t, err)
	}

	t.Run("setup weaviate", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to setup weaviate: %s", err)
		}
	})

	t.Run("class with MT config - MT enabled", func(t *testing.T) {
		defer cleanup()

		client := testsuit.CreateTestClient(false)
		className := "MultiTenantClass"
		schemaClass := &models.Class{
			Class: className,
			Properties: []*models.Property{
				{
					Name:     "someProperty",
					DataType: schema.DataTypeText.PropString(),
				},
			},
			MultiTenancyConfig: &models.MultiTenancyConfig{
				Enabled: true,
			},
		}

		err := client.Schema().ClassCreator().
			WithClass(schemaClass).
			Do(context.Background())
		require.Nil(t, err)

		t.Run("verify class created", func(t *testing.T) {
			loadedClass, err := client.Schema().ClassGetter().WithClassName(className).Do(context.Background())
			require.Nil(t, err)
			require.NotNil(t, loadedClass.MultiTenancyConfig)
			assert.Equal(t, true, loadedClass.MultiTenancyConfig.Enabled)
		})
	})

	t.Run("class with MT config - MT disabled", func(t *testing.T) {
		defer cleanup()

		client := testsuit.CreateTestClient(false)
		className := "MultiTenantClassDisabled"
		schemaClass := &models.Class{
			Class: className,
			Properties: []*models.Property{
				{
					Name:     "someProperty",
					DataType: schema.DataTypeText.PropString(),
				},
			},
			MultiTenancyConfig: &models.MultiTenancyConfig{
				Enabled: false,
			},
		}

		err := client.Schema().ClassCreator().
			WithClass(schemaClass).
			Do(context.Background())
		require.Nil(t, err)

		t.Run("verify class created", func(t *testing.T) {
			loadedClass, err := client.Schema().ClassGetter().WithClassName(className).Do(context.Background())
			require.Nil(t, err)
			require.NotNil(t, loadedClass.MultiTenancyConfig)
			assert.Equal(t, false, loadedClass.MultiTenancyConfig.Enabled)
		})
	})

	t.Run("class without MT config", func(t *testing.T) {
		defer cleanup()

		client := testsuit.CreateTestClient(false)
		className := "NonMultiTenantClass"
		schemaClass := &models.Class{
			Class: className,
			Properties: []*models.Property{
				{
					Name:     "someProperty",
					DataType: schema.DataTypeText.PropString(),
				},
			},
		}

		err := client.Schema().ClassCreator().
			WithClass(schemaClass).
			Do(context.Background())
		require.Nil(t, err)

		t.Run("verify class created", func(t *testing.T) {
			loadedClass, err := client.Schema().ClassGetter().WithClassName(className).Do(context.Background())
			require.Nil(t, err)
			require.NotNil(t, loadedClass.MultiTenancyConfig)
			assert.Equal(t, false, loadedClass.MultiTenancyConfig.Enabled)
		})
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to tear down weaviate: %s", err)
		}
	})
}

func TestSchema_Tenants(t *testing.T) {
	cleanup := func() {
		client := testsuit.CreateTestClient(false)
		err := client.Schema().AllDeleter().Do(context.Background())
		require.Nil(t, err)
	}

	t.Run("setup weaviate", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to setup weaviate: %s", err)
		}
	})

	className := "Pizza"
	ctx := context.Background()

	t.Run("adds tenants to MT class", func(t *testing.T) {
		defer cleanup()

		client := testsuit.CreateTestClient(false)
		testsuit.CreateSchemaPizzaForTenants(t, client)

		t.Run("adds single tenant", func(t *testing.T) {
			tenant := models.Tenant{
				Name: "tenantNo1",
			}

			err := client.Schema().TenantsCreator().
				WithClassName(className).
				WithTenants(tenant).
				Do(ctx)

			require.Nil(t, err)
		})

		t.Run("adds multiple tenants", func(t *testing.T) {
			tenants := []models.Tenant{
				{Name: "tenantNo2"},
				{Name: "tenantNo3"},
			}

			err := client.Schema().TenantsCreator().
				WithClassName(className).
				WithTenants(tenants...).
				Do(ctx)

			require.Nil(t, err)
		})
	})

	t.Run("fails adding tenants to non-MT class", func(t *testing.T) {
		defer cleanup()

		tenants := []models.Tenant{
			{Name: "tenantNo1"},
			{Name: "tenantNo2"},
		}

		client := testsuit.CreateTestClient(false)
		testsuit.CreateSchemaPizza(t, client)

		err := client.Schema().TenantsCreator().
			WithClassName(className).
			WithTenants(tenants...).
			Do(ctx)

		require.NotNil(t, err)
		clientErr := err.(*fault.WeaviateClientError)
		assert.Equal(t, 422, clientErr.StatusCode)
		assert.Contains(t, clientErr.Msg, "multi-tenancy is not enabled for class")
	})

	t.Run("gets tenants of MT class", func(t *testing.T) {
		defer cleanup()

		tenants := testsuit.Tenants{
			{Name: "tenantNo1"},
			{Name: "tenantNo2"},
		}

		client := testsuit.CreateTestClient(false)
		testsuit.CreateSchemaPizzaForTenants(t, client)
		testsuit.CreateTenantsPizza(t, client, tenants...)

		gotTenants, err := client.Schema().TenantsGetter().
			WithClassName(className).
			Do(ctx)

		require.Nil(t, err)
		require.Len(t, gotTenants, len(tenants))

		assert.ElementsMatch(t, tenants.Names(), testsuit.Tenants(gotTenants).Names())
	})

	t.Run("fails getting tenants from non-MT class", func(t *testing.T) {
		defer cleanup()

		client := testsuit.CreateTestClient(false)
		testsuit.CreateSchemaPizza(t, client)

		gotTenants, err := client.Schema().TenantsGetter().
			WithClassName(className).
			Do(ctx)

		require.NotNil(t, err)
		clientErr := err.(*fault.WeaviateClientError)
		assert.Equal(t, 422, clientErr.StatusCode)
		assert.Contains(t, clientErr.Msg, "multi-tenancy is not enabled for class")
		require.Nil(t, gotTenants)
	})

	t.Run("deletes tenants from MT class", func(t *testing.T) {
		defer cleanup()

		tenants := []models.Tenant{
			{Name: "tenantNo1"},
			{Name: "tenantNo2"},
			{Name: "tenantNo3"},
		}

		client := testsuit.CreateTestClient(false)
		testsuit.CreateSchemaPizzaForTenants(t, client)
		testsuit.CreateTenantsPizza(t, client, tenants...)

		t.Run("does not error on deleting non existent tenant", func(t *testing.T) {
			err := client.Schema().TenantsDeleter().
				WithClassName(className).
				WithTenants(tenants[0].Name, "nonExistentTenant").
				Do(ctx)

			require.Nil(t, err)
		})

		t.Run("deletes multiple tenants", func(t *testing.T) {
			err := client.Schema().TenantsDeleter().
				WithClassName(className).
				WithTenants(tenants[1].Name, tenants[2].Name).
				Do(ctx)

			require.Nil(t, err)
		})
	})

	t.Run("fails deleting tenants from non-MT class", func(t *testing.T) {
		defer cleanup()

		tenants := []string{"tenantNo1", "tenantNo2"}

		client := testsuit.CreateTestClient(false)
		testsuit.CreateSchemaPizza(t, client)

		err := client.Schema().TenantsDeleter().
			WithClassName(className).
			WithTenants(tenants...).
			Do(ctx)

		require.NotNil(t, err)
		clientErr := err.(*fault.WeaviateClientError)
		assert.Equal(t, 422, clientErr.StatusCode)
		assert.Contains(t, clientErr.Msg, "multi-tenancy is not enabled for class")
	})

	t.Run("updates tenants of MT class", func(t *testing.T) {
		defer cleanup()

		tenants := []models.Tenant{
			{Name: "tenantNo1"},
			{Name: "tenantNo2"},
		}

		client := testsuit.CreateTestClient(false)
		testsuit.CreateSchemaPizzaForTenants(t, client)
		testsuit.CreateTenantsPizza(t, client, tenants...)

		t.Run("fails updating not existent tenant", func(t *testing.T) {
			err := client.Schema().TenantsUpdater().
				WithClassName(className).
				WithTenants(models.Tenant{
					Name:           "nonExistentTenant",
					ActivityStatus: models.TenantActivityStatusCOLD,
				}).Do(ctx)

			require.NotNil(t, err)
			clientErr := err.(*fault.WeaviateClientError)
			assert.Equal(t, 422, clientErr.StatusCode)
			assert.Contains(t, clientErr.Msg, "not found")
		})

		t.Run("updates existent tenants", func(t *testing.T) {
			err := client.Schema().TenantsUpdater().
				WithClassName(className).
				WithTenants(
					models.Tenant{
						Name:           tenants[0].Name,
						ActivityStatus: models.TenantActivityStatusCOLD,
					},
					models.Tenant{
						Name:           tenants[1].Name,
						ActivityStatus: models.TenantActivityStatusCOLD,
					},
				).Do(ctx)

			require.Nil(t, err)
		})
	})

	t.Run("fails updating tenants of non-MT class", func(t *testing.T) {
		defer cleanup()

		tenants := []models.Tenant{
			{
				Name:           "tenantNo1",
				ActivityStatus: models.TenantActivityStatusCOLD,
			},
			{
				Name:           "tenantNo2",
				ActivityStatus: models.TenantActivityStatusCOLD,
			},
		}

		client := testsuit.CreateTestClient(false)
		testsuit.CreateSchemaPizza(t, client)

		err := client.Schema().TenantsUpdater().
			WithClassName(className).
			WithTenants(tenants...).
			Do(ctx)

		require.NotNil(t, err)
		clientErr := err.(*fault.WeaviateClientError)
		assert.Equal(t, 422, clientErr.StatusCode)
		assert.Contains(t, clientErr.Msg, "multi-tenancy is not enabled for class")
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to tear down weaviate: %s", err)
		}
	})
}

func TestSchema_TenantsActivationDeactivation(t *testing.T) {
	cleanup := func() {
		client := testsuit.CreateTestClient(false)
		err := client.Schema().AllDeleter().Do(context.Background())
		require.Nil(t, err)
	}

	t.Run("setup weaviate", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to setup weaviate: %s", err)
		}
	})

	t.Run("deactivate / activate journey", func(t *testing.T) {
		defer cleanup()

		tenants := testsuit.Tenants{
			{
				Name: "tenantNo1",
				// default status HOT
			},
			{
				Name:           "tenantNo2",
				ActivityStatus: models.TenantActivityStatusHOT,
			},
			{
				Name:           "tenantNo3",
				ActivityStatus: models.TenantActivityStatusCOLD,
			},
		}
		className := "Pizza"
		ctx := context.Background()
		ids := testsuit.IdsByClass[className]

		client := testsuit.CreateTestClient(false)

		assertTenantActive := func(t *testing.T, tenantName string) {
			gotTenants, err := client.Schema().TenantsGetter().
				WithClassName(className).
				Do(ctx)
			require.Nil(t, err)
			require.NotEmpty(t, gotTenants)

			byName := testsuit.Tenants(gotTenants).ByName(tenantName)
			require.NotNil(t, byName)
			require.Equal(t, models.TenantActivityStatusHOT, byName.ActivityStatus)

			objects, err := client.Data().ObjectsGetter().
				WithClassName(className).
				WithTenant(tenantName).
				Do(ctx)

			require.Nil(t, err)
			require.NotNil(t, objects)
			require.Len(t, objects, len(ids))
		}
		assertTenantInactive := func(t *testing.T, tenantName string) {
			gotTenants, err := client.Schema().TenantsGetter().
				WithClassName(className).
				Do(ctx)
			require.Nil(t, err)
			require.NotEmpty(t, gotTenants)

			byName := testsuit.Tenants(gotTenants).ByName(tenantName)
			require.NotNil(t, byName)
			require.Equal(t, models.TenantActivityStatusCOLD, byName.ActivityStatus)

			objects, err := client.Data().ObjectsGetter().
				WithClassName(className).
				WithTenant(tenantName).
				Do(ctx)

			require.NotNil(t, err)
			require.Nil(t, objects)
			clientErr := err.(*fault.WeaviateClientError)
			assert.Equal(t, 422, clientErr.StatusCode)
			assert.Contains(t, clientErr.Msg, "tenant not active")
		}

		t.Run("create tenants (1,2,3), populate active tenants (1,2)", func(t *testing.T) {
			testsuit.CreateSchemaPizzaForTenants(t, client)
			testsuit.CreateTenantsPizza(t, client, tenants...)
			testsuit.CreateDataPizzaForTenants(t, client, testsuit.Tenants(tenants[:2]).Names()...)

			assertTenantActive(t, tenants[0].Name)
			assertTenantActive(t, tenants[1].Name)
			assertTenantInactive(t, tenants[2].Name)
		})

		t.Run("deactivate tenant (1)", func(t *testing.T) {
			err := client.Schema().TenantsUpdater().
				WithClassName(className).
				WithTenants(models.Tenant{
					Name:           tenants[0].Name,
					ActivityStatus: models.TenantActivityStatusCOLD,
				}).
				Do(ctx)
			require.Nil(t, err)

			assertTenantInactive(t, tenants[0].Name)
			assertTenantActive(t, tenants[1].Name)
			assertTenantInactive(t, tenants[2].Name)
		})

		t.Run("activate and populate tenant (3)", func(t *testing.T) {
			err := client.Schema().TenantsUpdater().
				WithClassName(className).
				WithTenants(models.Tenant{
					Name:           tenants[2].Name,
					ActivityStatus: models.TenantActivityStatusHOT,
				}).
				Do(ctx)
			require.Nil(t, err)

			testsuit.CreateDataPizzaForTenants(t, client, tenants[2].Name)

			assertTenantInactive(t, tenants[0].Name)
			assertTenantActive(t, tenants[1].Name)
			assertTenantActive(t, tenants[2].Name)
		})

		t.Run("activate tenant (1)", func(t *testing.T) {
			err := client.Schema().TenantsUpdater().
				WithClassName(className).
				WithTenants(models.Tenant{
					Name:           tenants[0].Name,
					ActivityStatus: models.TenantActivityStatusHOT,
				}).
				Do(ctx)
			require.Nil(t, err)

			assertTenantActive(t, tenants[0].Name)
			assertTenantActive(t, tenants[1].Name)
			assertTenantActive(t, tenants[2].Name)
		})

		t.Run("deactivate tenant (2)", func(t *testing.T) {
			err := client.Schema().TenantsUpdater().
				WithClassName(className).
				WithTenants(models.Tenant{
					Name:           tenants[1].Name,
					ActivityStatus: models.TenantActivityStatusCOLD,
				}).
				Do(ctx)
			require.Nil(t, err)

			assertTenantActive(t, tenants[0].Name)
			assertTenantInactive(t, tenants[1].Name)
			assertTenantActive(t, tenants[2].Name)
		})

		t.Run("delete tenants", func(t *testing.T) {
			err := client.Schema().TenantsDeleter().
				WithClassName(className).
				WithTenants(tenants.Names()...).
				Do(ctx)

			require.Nil(t, err)
		})
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to tear down weaviate: %s", err)
		}
	})
}
