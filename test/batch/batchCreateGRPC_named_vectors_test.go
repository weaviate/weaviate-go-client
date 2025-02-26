package batch

import (
	"context"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/testenv"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/schema"
)

func TestBatchCreate_gRPC_named_vectors(t *testing.T) {
	ctx := context.Background()
	err := testenv.SetupLocalWeaviate()
	if err != nil {
		require.Nil(t, err, "failed to start weaviate")
	}

	client := testsuit.CreateTestClient(true)

	t.Run("batch import", func(t *testing.T) {
		tests := []struct {
			name         string
			vectorConfig map[string]models.VectorConfig
			vectors      models.Vectors
		}{
			{
				name: "regular vectors",
				vectorConfig: map[string]models.VectorConfig{
					"none1": {
						Vectorizer: map[string]interface{}{
							"none": nil,
						},
						VectorIndexType: "hnsw",
					},
					"none2": {
						Vectorizer: map[string]interface{}{
							"none": nil,
						},
						VectorIndexType: "flat",
					},
				},
				vectors: models.Vectors{
					"none1": []float32{1, 2},
					"none2": []float32{0.11, 0.22, 0.33},
				},
			},
			{
				name: "regular and colbert vectors",
				vectorConfig: map[string]models.VectorConfig{
					"regular": {
						Vectorizer: map[string]interface{}{
							"none": nil,
						},
						VectorIndexType: "hnsw",
					},
					"colbert": {
						Vectorizer: map[string]interface{}{
							"none": nil,
						},
						VectorIndexConfig: map[string]interface{}{
							"multivector": map[string]interface{}{
								"enabled": true,
							},
						},
						VectorIndexType: "hnsw",
					},
				},
				vectors: models.Vectors{
					"regular": []float32{1, 2},
					"colbert": [][]float32{{0.09, 0.11}, {0.22, 0.33}, {0.33, 0.44}},
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// clean up DB
				err := client.Schema().AllDeleter().Do(context.Background())
				require.NoError(t, err)
				// create class
				className := "NoVectorizer"
				class := &models.Class{
					Class: className,
					Properties: []*models.Property{
						{
							Name:     "name",
							DataType: []string{schema.DataTypeText.String()},
						},
					},
					VectorConfig: tt.vectorConfig,
				}
				err = client.Schema().ClassCreator().WithClass(class).Do(ctx)
				require.NoError(t, err)
				// perform batch import
				id := strfmt.UUID("00000000-0000-0000-0000-000000000001")
				vectors := tt.vectors
				object := &models.Object{
					ID:    id,
					Class: className,
					Properties: map[string]interface{}{
						"name": "some name",
					},
					Vectors: vectors,
				}
				batchResponse, err := client.Batch().ObjectsBatcher().WithObjects(object).Do(ctx)
				require.NoError(t, err)
				assert.NotNil(t, batchResponse)
				assert.Equal(t, 1, len(batchResponse))
				// get object
				objs, err := client.Data().ObjectsGetter().
					WithClassName(className).WithID(id.String()).WithVector().
					Do(ctx)
				require.NoError(t, err)
				require.NotEmpty(t, objs)
				require.Len(t, objs, 1)
				assert.Equal(t, id, objs[0].ID)
				require.Len(t, objs[0].Vectors, len(vectors))
				for targetVector := range vectors {
					assert.Equal(t, vectors[targetVector], objs[0].Vectors[targetVector])
				}
			})
		}
	})

	err = testenv.TearDownLocalWeaviate()
	if err != nil {
		require.Nil(t, err, "failed to tear down weaviate")
	}
}
