package voyageai

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

// TestVoyageAI_Integration_Voyage4 tests real API calls with voyage-4 model
func TestVoyageAI_Integration_Voyage4(t *testing.T) {
	if os.Getenv("VOYAGEAI_APIKEY") == "" {
		t.Skip("VOYAGEAI_APIKEY not set, skipping integration test")
	}

	t.Run("setup weaviate", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to setup weaviate: %s", err)
		}
	})

	t.Run("Create, insert and query with voyage-4", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		ctx := context.Background()
		className := "Voyage4Integration"

		// Clean up before test
		_ = client.Schema().ClassDeleter().WithClassName(className).Do(ctx)

		// Create schema with voyage-4
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
					Name:     "content",
					DataType: []string{"text"},
				},
			},
		}

		err := client.Schema().ClassCreator().WithClass(schemaClass).Do(ctx)
		require.Nil(t, err, "Failed to create schema")

		// Insert test data
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

		// Verify objects have vectors
		result, err := client.Data().ObjectsGetter().
			WithClassName(className).
			WithID("00000000-0000-0000-0000-000000000001").
			WithVector().
			Do(ctx)
		require.Nil(t, err)
		require.Len(t, result, 1)
		require.NotEmpty(t, result[0].Vector, "Object should have a vector")

		// Expected dimension is 1024 for voyage-4
		assert.Equal(t, 1024, len(result[0].Vector), "voyage-4 default dimension should be 1024")

		// Perform semantic search
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

// TestVoyageAI_Integration_Voyage4Lite tests real API calls with voyage-4-lite model
func TestVoyageAI_Integration_Voyage4Lite(t *testing.T) {
	if os.Getenv("VOYAGEAI_APIKEY") == "" {
		t.Skip("VOYAGEAI_APIKEY not set, skipping integration test")
	}

	t.Run("setup weaviate", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to setup weaviate: %s", err)
		}
	})

	t.Run("Create, insert and query with voyage-4-lite", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		ctx := context.Background()
		className := "Voyage4LiteIntegration"

		// Clean up before test
		_ = client.Schema().ClassDeleter().WithClassName(className).Do(ctx)

		// Create schema with voyage-4-lite
		schemaClass := &models.Class{
			Class:      className,
			Vectorizer: "text2vec-voyageai",
			ModuleConfig: map[string]interface{}{
				"text2vec-voyageai": map[string]interface{}{
					"model":              "voyage-4-lite",
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

		// Insert test data
		testObjects := []*models.Object{
			{
				Class: className,
				ID:    strfmt.UUID("00000000-0000-0000-0000-000000000011"),
				Properties: map[string]interface{}{
					"content": "Python is a popular programming language for data science.",
				},
			},
			{
				Class: className,
				ID:    strfmt.UUID("00000000-0000-0000-0000-000000000012"),
				Properties: map[string]interface{}{
					"content": "JavaScript is widely used for web development.",
				},
			},
			{
				Class: className,
				ID:    strfmt.UUID("00000000-0000-0000-0000-000000000013"),
				Properties: map[string]interface{}{
					"content": "Go is known for its simplicity and performance.",
				},
			},
		}

		batcher := client.Batch().ObjectsBatcher()
		for _, obj := range testObjects {
			batcher.WithObject(obj)
		}
		_, err = batcher.Do(ctx)
		require.Nil(t, err, "Failed to batch insert objects")

		// Verify objects have vectors
		result, err := client.Data().ObjectsGetter().
			WithClassName(className).
			WithID("00000000-0000-0000-0000-000000000011").
			WithVector().
			Do(ctx)
		require.Nil(t, err)
		require.Len(t, result, 1)
		require.NotEmpty(t, result[0].Vector, "Object should have a vector")

		// Expected dimension is 1024 for voyage-4-lite
		assert.Equal(t, 1024, len(result[0].Vector), "voyage-4-lite default dimension should be 1024")

		// Perform semantic search
		searchResult, err := client.GraphQL().Get().
			WithClassName(className).
			WithFields(
				graphql.Field{Name: "content"},
				graphql.Field{Name: "_additional", Fields: []graphql.Field{{Name: "id"}, {Name: "distance"}}},
			).
			WithNearText(client.GraphQL().NearTextArgBuilder().WithConcepts([]string{"programming languages"})).
			WithLimit(3).
			Do(ctx)
		require.Nil(t, err)
		require.NotNil(t, searchResult)
		require.Nil(t, searchResult.Errors)

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

// TestVoyageAI_Integration_Voyage4Large tests real API calls with voyage-4-large model
func TestVoyageAI_Integration_Voyage4Large(t *testing.T) {
	if os.Getenv("VOYAGEAI_APIKEY") == "" {
		t.Skip("VOYAGEAI_APIKEY not set, skipping integration test")
	}

	t.Run("setup weaviate", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to setup weaviate: %s", err)
		}
	})

	t.Run("Create, insert and query with voyage-4-large", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		ctx := context.Background()
		className := "Voyage4LargeIntegration"

		// Clean up before test
		_ = client.Schema().ClassDeleter().WithClassName(className).Do(ctx)

		// Create schema with voyage-4-large
		schemaClass := &models.Class{
			Class:      className,
			Vectorizer: "text2vec-voyageai",
			ModuleConfig: map[string]interface{}{
				"text2vec-voyageai": map[string]interface{}{
					"model":              "voyage-4-large",
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

		// Insert test data
		testObjects := []*models.Object{
			{
				Class: className,
				ID:    strfmt.UUID("00000000-0000-0000-0000-000000000021"),
				Properties: map[string]interface{}{
					"content": "Climate change is affecting global weather patterns.",
				},
			},
			{
				Class: className,
				ID:    strfmt.UUID("00000000-0000-0000-0000-000000000022"),
				Properties: map[string]interface{}{
					"content": "Renewable energy sources are becoming more cost effective.",
				},
			},
			{
				Class: className,
				ID:    strfmt.UUID("00000000-0000-0000-0000-000000000023"),
				Properties: map[string]interface{}{
					"content": "Electric vehicles are reducing carbon emissions.",
				},
			},
		}

		batcher := client.Batch().ObjectsBatcher()
		for _, obj := range testObjects {
			batcher.WithObject(obj)
		}
		_, err = batcher.Do(ctx)
		require.Nil(t, err, "Failed to batch insert objects")

		// Verify objects have vectors
		result, err := client.Data().ObjectsGetter().
			WithClassName(className).
			WithID("00000000-0000-0000-0000-000000000021").
			WithVector().
			Do(ctx)
		require.Nil(t, err)
		require.Len(t, result, 1)
		require.NotEmpty(t, result[0].Vector, "Object should have a vector")

		// Expected dimension is 1024 for voyage-4-large
		assert.Equal(t, 1024, len(result[0].Vector), "voyage-4-large default dimension should be 1024")

		// Perform semantic search
		searchResult, err := client.GraphQL().Get().
			WithClassName(className).
			WithFields(
				graphql.Field{Name: "content"},
				graphql.Field{Name: "_additional", Fields: []graphql.Field{{Name: "id"}, {Name: "distance"}}},
			).
			WithNearText(client.GraphQL().NearTextArgBuilder().WithConcepts([]string{"environmental sustainability"})).
			WithLimit(3).
			Do(ctx)
		require.Nil(t, err)
		require.NotNil(t, searchResult)
		require.Nil(t, searchResult.Errors)

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

// TestVoyageAI_Integration_FlexibleDimensions tests real API calls with custom dimensions
func TestVoyageAI_Integration_FlexibleDimensions(t *testing.T) {
	if os.Getenv("VOYAGEAI_APIKEY") == "" {
		t.Skip("VOYAGEAI_APIKEY not set, skipping integration test")
	}

	t.Run("setup weaviate", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to setup weaviate: %s", err)
		}
	})

	testCases := []struct {
		name       string
		dimensions int
	}{
		{"256 dimensions", 256},
		{"512 dimensions", 512},
		{"2048 dimensions", 2048},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := testsuit.CreateTestClient(false)
			ctx := context.Background()
			className := "VoyageFlexDimsIntegration"

			// Clean up before test
			_ = client.Schema().ClassDeleter().WithClassName(className).Do(ctx)

			// Create schema with custom dimensions
			schemaClass := &models.Class{
				Class:      className,
				Vectorizer: "text2vec-voyageai",
				ModuleConfig: map[string]interface{}{
					"text2vec-voyageai": map[string]interface{}{
						"model":              "voyage-4",
						"dimensions":         tc.dimensions,
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
			require.Nil(t, err, "Failed to create schema with %d dimensions", tc.dimensions)

			// Insert test data
			obj := &models.Object{
				Class: className,
				ID:    strfmt.UUID("00000000-0000-0000-0000-000000000099"),
				Properties: map[string]interface{}{
					"content": "Testing flexible dimensions with VoyageAI embeddings.",
				},
			}

			_, err = client.Data().Creator().WithClassName(className).WithID(obj.ID.String()).WithProperties(obj.Properties).Do(ctx)
			require.Nil(t, err, "Failed to insert object")

			// Verify vector dimension
			result, err := client.Data().ObjectsGetter().
				WithClassName(className).
				WithID("00000000-0000-0000-0000-000000000099").
				WithVector().
				Do(ctx)
			require.Nil(t, err)
			require.Len(t, result, 1)
			require.NotEmpty(t, result[0].Vector)
			assert.Equal(t, tc.dimensions, len(result[0].Vector), "Vector should have %d dimensions", tc.dimensions)

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

// TestVoyageAI_Integration_NamedVectors tests real API calls with named vectors
func TestVoyageAI_Integration_NamedVectors(t *testing.T) {
	if os.Getenv("VOYAGEAI_APIKEY") == "" {
		t.Skip("VOYAGEAI_APIKEY not set, skipping integration test")
	}

	t.Run("setup weaviate", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to setup weaviate: %s", err)
		}
	})

	t.Run("Create and query with named vectors", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		ctx := context.Background()
		className := "VoyageNamedVectorsIntegration"

		// Clean up before test
		_ = client.Schema().ClassDeleter().WithClassName(className).Do(ctx)

		// Create schema with named vectors
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
							"model":      "voyage-4",
							"properties": []interface{}{"content"},
						},
					},
					VectorIndexType: "hnsw",
				},
			},
		}

		err := client.Schema().ClassCreator().WithClass(schemaClass).Do(ctx)
		require.Nil(t, err, "Failed to create schema with named vectors")

		// Insert test data
		obj := &models.Object{
			Class: className,
			ID:    strfmt.UUID("00000000-0000-0000-0000-000000000088"),
			Properties: map[string]interface{}{
				"title":   "Introduction to Vector Databases",
				"content": "Vector databases store and retrieve high-dimensional vectors for similarity search applications.",
			},
		}

		_, err = client.Data().Creator().
			WithClassName(className).
			WithID(obj.ID.String()).
			WithProperties(obj.Properties).
			Do(ctx)
		require.Nil(t, err, "Failed to insert object")

		// Verify object has named vectors
		result, err := client.Data().ObjectsGetter().
			WithClassName(className).
			WithID("00000000-0000-0000-0000-000000000088").
			WithVector().
			Do(ctx)
		require.Nil(t, err)
		require.Len(t, result, 1)

		// Named vectors should be present
		require.NotNil(t, result[0].Vectors, "Object should have named vectors")
		assert.Contains(t, result[0].Vectors, "title_vector")
		assert.Contains(t, result[0].Vectors, "content_vector")

		// Verify the vectors exist (interface types can't be checked for length directly)
		titleVector := result[0].Vectors["title_vector"]
		contentVector := result[0].Vectors["content_vector"]
		require.NotNil(t, titleVector, "title_vector should not be nil")
		require.NotNil(t, contentVector, "content_vector should not be nil")

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
