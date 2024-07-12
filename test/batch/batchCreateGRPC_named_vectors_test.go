package batch

import (
	"context"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v4/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/testenv"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/schema"
)

func TestBatchCreate_gRPC_named_vectors(t *testing.T) {
	ctx := context.Background()
	err := testenv.SetupLocalWeaviate()
	if err != nil {
		require.Nil(t, err, "failed to start weaviate")
	}

	client := testsuit.CreateTestClient()

	t.Run("multiple vectors", func(t *testing.T) {
		// run batch import
		t.Run("gRPC batch import", func(t *testing.T) {
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
				VectorConfig: map[string]models.VectorConfig{
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
			}
			err = client.Schema().ClassCreator().WithClass(class).Do(ctx)
			require.NoError(t, err)
			// perform batch import
			id := strfmt.UUID("00000000-0000-0000-0000-000000000001")
			vectors := models.Vectors{
				"none1": []float32{1, 2},
				"none2": []float32{0.11, 0.22, 0.33},
			}
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
			require.Len(t, objs[0].Vectors, 2)
			assert.Equal(t, vectors["none1"], objs[0].Vectors["none1"])
			assert.Equal(t, vectors["none2"], objs[0].Vectors["none2"])
		})
	})

	err = testenv.TearDownLocalWeaviate()
	if err != nil {
		require.Nil(t, err, "failed to tear down weaviate")
	}
}
