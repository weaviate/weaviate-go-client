package voyageai

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/testenv"
	"github.com/weaviate/weaviate/entities/models"
)

// Voyage-4 family model specifications
var voyage4Models = []struct {
	name              string
	model             string
	defaultDimensions int
	maxInputTokens    int
	description       string
}{
	{
		name:              "voyage-4",
		model:             "voyage-4",
		defaultDimensions: 1024,
		maxInputTokens:    32000,
		description:       "Balanced general-purpose and multilingual retrieval",
	},
	{
		name:              "voyage-4-lite",
		model:             "voyage-4-lite",
		defaultDimensions: 1024,
		maxInputTokens:    32000,
		description:       "Best performance-to-cost ratio, highest batch throughput",
	},
	{
		name:              "voyage-4-large",
		model:             "voyage-4-large",
		defaultDimensions: 1024,
		maxInputTokens:    32000,
		description:       "Highest retrieval quality",
	},
}

// TestVoyageAI_SchemaConfiguration tests schema creation with VoyageAI vectorizer
func TestVoyageAI_SchemaConfiguration(t *testing.T) {
	t.Run("setup weaviate", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to setup weaviate: %s", err)
		}
	})

	for _, modelSpec := range voyage4Models {
		t.Run("Create schema with "+modelSpec.name, func(t *testing.T) {
			client := testsuit.CreateTestClient(false)
			ctx := context.Background()
			className := "VoyageTest" + modelSpec.name

			// Clean up before test
			_ = client.Schema().ClassDeleter().WithClassName(className).Do(ctx)

			schemaClass := &models.Class{
				Class:       className,
				Description: modelSpec.description,
				Vectorizer:  "text2vec-voyageai",
				ModuleConfig: map[string]interface{}{
					"text2vec-voyageai": map[string]interface{}{
						"model":              modelSpec.model,
						"vectorizeClassName": true,
					},
				},
				Properties: []*models.Property{
					{
						Name:        "content",
						DataType:    []string{"text"},
						Description: "The text content to vectorize",
					},
				},
			}

			err := client.Schema().ClassCreator().WithClass(schemaClass).Do(ctx)
			require.Nil(t, err, "Failed to create schema for %s", modelSpec.name)

			// Verify schema was created
			loadedClass, err := client.Schema().ClassGetter().WithClassName(className).Do(ctx)
			require.Nil(t, err)
			require.NotNil(t, loadedClass)

			assert.Equal(t, className, loadedClass.Class)
			assert.Equal(t, "text2vec-voyageai", loadedClass.Vectorizer)
			assert.NotNil(t, loadedClass.ModuleConfig)

			// Verify module config
			moduleConfig, ok := loadedClass.ModuleConfig.(map[string]interface{})
			require.True(t, ok, "ModuleConfig should be a map")

			voyageConfig, ok := moduleConfig["text2vec-voyageai"].(map[string]interface{})
			require.True(t, ok, "text2vec-voyageai config should exist")
			assert.Equal(t, modelSpec.model, voyageConfig["model"])

			// Clean up after test
			err = client.Schema().ClassDeleter().WithClassName(className).Do(ctx)
			assert.Nil(t, err)
		})
	}

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to tear down weaviate: %s", err)
		}
	})
}

// TestVoyageAI_FlexibleDimensions tests creating schemas with different output dimensions
func TestVoyageAI_FlexibleDimensions(t *testing.T) {
	t.Run("setup weaviate", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to setup weaviate: %s", err)
		}
	})

	flexibleDimensions := []int{256, 512, 1024, 2048}

	for _, dims := range flexibleDimensions {
		t.Run("Create schema with voyage-4 and dimensions "+string(rune(dims)), func(t *testing.T) {
			client := testsuit.CreateTestClient(false)
			ctx := context.Background()
			className := "VoyageFlexDims"

			// Clean up before test
			_ = client.Schema().ClassDeleter().WithClassName(className).Do(ctx)

			schemaClass := &models.Class{
				Class:      className,
				Vectorizer: "text2vec-voyageai",
				ModuleConfig: map[string]interface{}{
					"text2vec-voyageai": map[string]interface{}{
						"model":              "voyage-4",
						"dimensions":         dims,
						"vectorizeClassName": false,
					},
				},
				Properties: []*models.Property{
					{
						Name:     "content",
						DataType: []string{"text"},
					},
				},
			}

			err := client.Schema().ClassCreator().WithClass(schemaClass).Do(ctx)
			require.Nil(t, err, "Failed to create schema with dimensions %d", dims)

			// Verify schema was created
			loadedClass, err := client.Schema().ClassGetter().WithClassName(className).Do(ctx)
			require.Nil(t, err)
			require.NotNil(t, loadedClass)

			// Verify module config
			moduleConfig, ok := loadedClass.ModuleConfig.(map[string]interface{})
			require.True(t, ok)

			voyageConfig, ok := moduleConfig["text2vec-voyageai"].(map[string]interface{})
			require.True(t, ok)

			// Dimensions should be set (may be returned as float64 from JSON)
			if configDims, exists := voyageConfig["dimensions"]; exists {
				switch v := configDims.(type) {
				case float64:
					assert.Equal(t, float64(dims), v)
				case int:
					assert.Equal(t, dims, v)
				}
			}

			// Clean up
			err = client.Schema().ClassDeleter().WithClassName(className).Do(ctx)
			assert.Nil(t, err)
		})
	}

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to tear down weaviate: %s", err)
		}
	})
}

// TestVoyageAI_NamedVectors tests creating schemas with named vectors using VoyageAI
func TestVoyageAI_NamedVectors(t *testing.T) {
	t.Run("setup weaviate", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to setup weaviate: %s", err)
		}
	})

	t.Run("Create schema with multiple named vectors", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		ctx := context.Background()
		className := "VoyageNamedVectors"

		// Clean up before test
		_ = client.Schema().ClassDeleter().WithClassName(className).Do(ctx)

		schemaClass := &models.Class{
			Class: className,
			Properties: []*models.Property{
				{
					Name:     "title",
					DataType: []string{"text"},
				},
				{
					Name:     "content",
					DataType: []string{"text"},
				},
			},
			VectorConfig: map[string]models.VectorConfig{
				"title_vector": {
					Vectorizer: map[string]interface{}{
						"text2vec-voyageai": map[string]interface{}{
							"model":      "voyage-4-lite",
							"properties": []interface{}{"title"},
						},
					},
					VectorIndexType: "hnsw",
				},
				"content_vector": {
					Vectorizer: map[string]interface{}{
						"text2vec-voyageai": map[string]interface{}{
							"model":      "voyage-4-large",
							"properties": []interface{}{"content"},
						},
					},
					VectorIndexType: "hnsw",
				},
			},
		}

		err := client.Schema().ClassCreator().WithClass(schemaClass).Do(ctx)
		require.Nil(t, err, "Failed to create schema with named vectors")

		// Verify schema was created
		loadedClass, err := client.Schema().ClassGetter().WithClassName(className).Do(ctx)
		require.Nil(t, err)
		require.NotNil(t, loadedClass)

		assert.Equal(t, className, loadedClass.Class)
		require.NotNil(t, loadedClass.VectorConfig)
		assert.Contains(t, loadedClass.VectorConfig, "title_vector")
		assert.Contains(t, loadedClass.VectorConfig, "content_vector")

		// Clean up
		err = client.Schema().ClassDeleter().WithClassName(className).Do(ctx)
		assert.Nil(t, err)
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to tear down weaviate: %s", err)
		}
	})
}

// TestVoyageAI_PropertyConfiguration tests property-level vectorization settings
func TestVoyageAI_PropertyConfiguration(t *testing.T) {
	t.Run("setup weaviate", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to setup weaviate: %s", err)
		}
	})

	t.Run("Create schema with property skip configuration", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		ctx := context.Background()
		className := "VoyagePropertyConfig"

		// Clean up before test
		_ = client.Schema().ClassDeleter().WithClassName(className).Do(ctx)

		schemaClass := &models.Class{
			Class:      className,
			Vectorizer: "text2vec-voyageai",
			ModuleConfig: map[string]interface{}{
				"text2vec-voyageai": map[string]interface{}{
					"model":              "voyage-4",
					"vectorizeClassName": false,
				},
			},
			Properties: []*models.Property{
				{
					Name:        "title",
					DataType:    []string{"text"},
					Description: "Title to vectorize",
				},
				{
					Name:        "content",
					DataType:    []string{"text"},
					Description: "Content to vectorize",
				},
				{
					Name:        "metadata",
					DataType:    []string{"text"},
					Description: "Metadata to skip",
					ModuleConfig: map[string]interface{}{
						"text2vec-voyageai": map[string]interface{}{
							"skip": true,
						},
					},
				},
			},
		}

		err := client.Schema().ClassCreator().WithClass(schemaClass).Do(ctx)
		require.Nil(t, err, "Failed to create schema with property config")

		// Verify schema was created
		loadedClass, err := client.Schema().ClassGetter().WithClassName(className).Do(ctx)
		require.Nil(t, err)
		require.NotNil(t, loadedClass)
		require.Len(t, loadedClass.Properties, 3)

		// Clean up
		err = client.Schema().ClassDeleter().WithClassName(className).Do(ctx)
		assert.Nil(t, err)
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to tear down weaviate: %s", err)
		}
	})
}
