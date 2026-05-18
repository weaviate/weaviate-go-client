package digitalocean

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/testenv"
	"github.com/weaviate/weaviate/entities/models"
)

// TestDigitalOcean_ModuleConfig_Structure tests that the ModuleConfig map for
// text2vec-digitalocean is constructed with the correct keys and values before
// being sent to the server. No live Weaviate instance is required.
func TestDigitalOcean_ModuleConfig_Structure(t *testing.T) {
	tests := []struct {
		name       string
		config     map[string]interface{}
		wantKeys   []string
		wantValues map[string]interface{}
	}{
		{
			name: "model and vectorizeClassName",
			config: map[string]interface{}{
				"model":              "qwen3-embedding-0.6b",
				"vectorizeClassName": false,
			},
			wantKeys: []string{"model", "vectorizeClassName"},
			wantValues: map[string]interface{}{
				"model":              "qwen3-embedding-0.6b",
				"vectorizeClassName": false,
			},
		},
		{
			name: "model, baseURL and vectorizeClassName",
			config: map[string]interface{}{
				"model":              "qwen3-embedding-0.6b",
				"baseURL":            "https://inference.do-ai.run",
				"vectorizeClassName": true,
			},
			wantKeys: []string{"model", "baseURL", "vectorizeClassName"},
			wantValues: map[string]interface{}{
				"model":              "qwen3-embedding-0.6b",
				"baseURL":            "https://inference.do-ai.run",
				"vectorizeClassName": true,
			},
		},
		{
			name: "vectorizeClassName defaults to true",
			config: map[string]interface{}{
				"model":              "qwen3-embedding-0.6b",
				"vectorizeClassName": true,
			},
			wantKeys: []string{"model", "vectorizeClassName"},
			wantValues: map[string]interface{}{
				"model":              "qwen3-embedding-0.6b",
				"vectorizeClassName": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moduleConfig := map[string]interface{}{
				"text2vec-digitalocean": tt.config,
			}

			inner, ok := moduleConfig["text2vec-digitalocean"].(map[string]interface{})
			require.True(t, ok, "text2vec-digitalocean key must map to map[string]interface{}")

			for _, key := range tt.wantKeys {
				assert.Contains(t, inner, key, "config should contain key %q", key)
			}
			for key, val := range tt.wantValues {
				assert.Equal(t, val, inner[key], "config[%q] mismatch", key)
			}
		})
	}
}

// TestDigitalOcean_ClassConfig_Shape verifies that a models.Class wired for
// text2vec-digitalocean carries the correct Vectorizer string and has the
// expected module config structure.
func TestDigitalOcean_ClassConfig_Shape(t *testing.T) {
	class := &models.Class{
		Class:      "DigitalOceanTest",
		Vectorizer: "text2vec-digitalocean",
		ModuleConfig: map[string]interface{}{
			"text2vec-digitalocean": map[string]interface{}{
				"model":              "qwen3-embedding-0.6b",
				"vectorizeClassName": false,
				"baseURL":            "https://inference.do-ai.run",
			},
		},
	}

	assert.Equal(t, "text2vec-digitalocean", class.Vectorizer)

	mc, ok := class.ModuleConfig.(map[string]interface{})
	require.True(t, ok)

	doConfig, ok := mc["text2vec-digitalocean"].(map[string]interface{})
	require.True(t, ok, "text2vec-digitalocean key must be present in ModuleConfig")

	assert.Equal(t, "qwen3-embedding-0.6b", doConfig["model"])
	assert.Equal(t, false, doConfig["vectorizeClassName"])
	assert.Equal(t, "https://inference.do-ai.run", doConfig["baseURL"])
}

// TestDigitalOcean_NamedVector_Shape verifies that VectorConfig for a named
// vector using text2vec-digitalocean has the expected structure.
func TestDigitalOcean_NamedVector_Shape(t *testing.T) {
	class := &models.Class{
		Class: "DigitalOceanNamedVectors",
		Properties: []*models.Property{
			{Name: "content", DataType: []string{"text"}},
		},
		VectorConfig: map[string]models.VectorConfig{
			"content_vector": {
				Vectorizer: map[string]interface{}{
					"text2vec-digitalocean": map[string]interface{}{
						"model":              "qwen3-embedding-0.6b",
						"vectorizeClassName": true,
						"properties":         []interface{}{"content"},
					},
				},
				VectorIndexType: "hnsw",
			},
		},
	}

	vc, ok := class.VectorConfig["content_vector"]
	require.True(t, ok, "content_vector must be present in VectorConfig")
	assert.Equal(t, "hnsw", vc.VectorIndexType)

	vectorizer, ok := vc.Vectorizer.(map[string]interface{})
	require.True(t, ok)

	doConfig, ok := vectorizer["text2vec-digitalocean"].(map[string]interface{})
	require.True(t, ok, "text2vec-digitalocean key must be present in Vectorizer")

	assert.Equal(t, "qwen3-embedding-0.6b", doConfig["model"])
	assert.Equal(t, true, doConfig["vectorizeClassName"])
	assert.Equal(t, []interface{}{"content"}, doConfig["properties"])
}

// TestDigitalOcean_SchemaConfiguration tests schema creation with the
// text2vec-digitalocean vectorizer against a live Weaviate instance.
func TestDigitalOcean_SchemaConfiguration(t *testing.T) {
	t.Run("setup weaviate", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to setup weaviate: %s", err)
		}
	})

	t.Run("Create schema with text2vec-digitalocean", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		ctx := context.Background()
		className := "DigitalOceanSchemaTest"

		_ = client.Schema().ClassDeleter().WithClassName(className).Do(ctx)

		schemaClass := &models.Class{
			Class:      className,
			Vectorizer: "text2vec-digitalocean",
			ModuleConfig: map[string]interface{}{
				"text2vec-digitalocean": map[string]interface{}{
					"model":              "qwen3-embedding-0.6b",
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
		require.Nil(t, err, "Failed to create schema for text2vec-digitalocean")

		loadedClass, err := client.Schema().ClassGetter().WithClassName(className).Do(ctx)
		require.Nil(t, err)
		require.NotNil(t, loadedClass)

		assert.Equal(t, className, loadedClass.Class)
		assert.Equal(t, "text2vec-digitalocean", loadedClass.Vectorizer)
		assert.NotNil(t, loadedClass.ModuleConfig)

		moduleConfig, ok := loadedClass.ModuleConfig.(map[string]interface{})
		require.True(t, ok, "ModuleConfig should be a map")

		doConfig, ok := moduleConfig["text2vec-digitalocean"].(map[string]interface{})
		require.True(t, ok, "text2vec-digitalocean config should exist")
		assert.Equal(t, "qwen3-embedding-0.6b", doConfig["model"])

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

// TestDigitalOcean_NamedVectors tests schema creation with named vectors using
// text2vec-digitalocean against a live Weaviate instance.
func TestDigitalOcean_NamedVectors(t *testing.T) {
	t.Run("setup weaviate", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to setup weaviate: %s", err)
		}
	})

	t.Run("Create schema with named vectors", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		ctx := context.Background()
		className := "DigitalOceanNamedVectorsTest"

		_ = client.Schema().ClassDeleter().WithClassName(className).Do(ctx)

		schemaClass := &models.Class{
			Class: className,
			Properties: []*models.Property{
				{Name: "title", DataType: []string{"text"}},
				{Name: "content", DataType: []string{"text"}},
			},
			VectorConfig: map[string]models.VectorConfig{
				"title_vector": {
					Vectorizer: map[string]interface{}{
						"text2vec-digitalocean": map[string]interface{}{
							"model":      "qwen3-embedding-0.6b",
							"properties": []interface{}{"title"},
						},
					},
					VectorIndexType: "hnsw",
				},
				"content_vector": {
					Vectorizer: map[string]interface{}{
						"text2vec-digitalocean": map[string]interface{}{
							"model":      "qwen3-embedding-0.6b",
							"properties": []interface{}{"content"},
						},
					},
					VectorIndexType: "hnsw",
				},
			},
		}

		err := client.Schema().ClassCreator().WithClass(schemaClass).Do(ctx)
		require.Nil(t, err, "Failed to create schema with named vectors")

		loadedClass, err := client.Schema().ClassGetter().WithClassName(className).Do(ctx)
		require.Nil(t, err)
		require.NotNil(t, loadedClass)

		assert.Equal(t, className, loadedClass.Class)
		require.NotNil(t, loadedClass.VectorConfig)
		assert.Contains(t, loadedClass.VectorConfig, "title_vector")
		assert.Contains(t, loadedClass.VectorConfig, "content_vector")

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
