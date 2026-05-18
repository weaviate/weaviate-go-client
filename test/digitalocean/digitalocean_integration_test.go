package digitalocean

import (
	"context"
	"os"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/graphql"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/testenv"
	"github.com/weaviate/weaviate/entities/models"
)

// TestDigitalOcean_Integration tests real API calls with the text2vec-digitalocean
// vectorizer. Requires DIGITALOCEAN_APIKEY to be set in the environment.
func TestDigitalOcean_Integration(t *testing.T) {
	if os.Getenv("DIGITALOCEAN_APIKEY") == "" {
		t.Skip("DIGITALOCEAN_APIKEY not set, skipping integration test")
	}

	t.Run("setup weaviate", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to setup weaviate: %s", err)
		}
	})

	t.Run("Create, insert and query with text2vec-digitalocean", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		ctx := context.Background()
		className := "DigitalOceanIntegration"

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
		require.Nil(t, err, "Failed to create schema")

		testObjects := []*models.Object{
			{
				Class: className,
				ID:    strfmt.UUID("00000000-0000-0000-0000-000000000001"),
				Properties: map[string]interface{}{
					"content": "Machine learning is a subset of artificial intelligence.",
				},
			},
			{
				Class: className,
				ID:    strfmt.UUID("00000000-0000-0000-0000-000000000002"),
				Properties: map[string]interface{}{
					"content": "Deep learning uses neural networks with many layers.",
				},
			},
			{
				Class: className,
				ID:    strfmt.UUID("00000000-0000-0000-0000-000000000003"),
				Properties: map[string]interface{}{
					"content": "Natural language processing helps computers understand text.",
				},
			},
		}

		batcher := client.Batch().ObjectsBatcher()
		for _, obj := range testObjects {
			batcher.WithObject(obj)
		}
		_, err = batcher.Do(ctx)
		require.Nil(t, err, "Failed to batch insert objects")

		result, err := client.Data().ObjectsGetter().
			WithClassName(className).
			WithID("00000000-0000-0000-0000-000000000001").
			WithVector().
			Do(ctx)
		require.Nil(t, err)
		require.Len(t, result, 1)
		require.NotEmpty(t, result[0].Vector, "Object should have a vector")

		searchResult, err := client.GraphQL().Get().
			WithClassName(className).
			WithFields(
				graphql.Field{Name: "content"},
				graphql.Field{Name: "_additional", Fields: []graphql.Field{{Name: "id"}, {Name: "distance"}}},
			).
			WithNearText(client.GraphQL().NearTextArgBuilder().WithConcepts([]string{"AI and neural networks"})).
			WithLimit(3).
			Do(ctx)
		require.Nil(t, err)
		require.NotNil(t, searchResult)
		require.Nil(t, searchResult.Errors)

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

// TestDigitalOcean_Integration_WithBaseURL tests specifying a custom baseURL
// alongside the model.
func TestDigitalOcean_Integration_WithBaseURL(t *testing.T) {
	if os.Getenv("DIGITALOCEAN_APIKEY") == "" {
		t.Skip("DIGITALOCEAN_APIKEY not set, skipping integration test")
	}

	t.Run("setup weaviate", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to setup weaviate: %s", err)
		}
	})

	t.Run("Create schema with explicit baseURL", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		ctx := context.Background()
		className := "DigitalOceanBaseURLIntegration"

		_ = client.Schema().ClassDeleter().WithClassName(className).Do(ctx)

		schemaClass := &models.Class{
			Class:      className,
			Vectorizer: "text2vec-digitalocean",
			ModuleConfig: map[string]interface{}{
				"text2vec-digitalocean": map[string]interface{}{
					"model":              "qwen3-embedding-0.6b",
					"baseURL":            "https://inference.do-ai.run",
					"vectorizeClassName": false,
				},
			},
			Properties: []*models.Property{
				{Name: "content", DataType: []string{"text"}},
			},
		}

		err := client.Schema().ClassCreator().WithClass(schemaClass).Do(ctx)
		require.Nil(t, err, "Failed to create schema with explicit baseURL")

		loadedClass, err := client.Schema().ClassGetter().WithClassName(className).Do(ctx)
		require.Nil(t, err)
		require.NotNil(t, loadedClass)
		assert.Equal(t, "text2vec-digitalocean", loadedClass.Vectorizer)

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
