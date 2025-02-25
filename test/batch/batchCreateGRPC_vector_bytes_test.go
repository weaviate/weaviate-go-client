package batch

import (
	"context"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test/docker"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/grpc"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/schema"
)

func TestBatchCreate_gRPC_vector_bytes_field(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name          string
		weaviateImage string
	}{
		{
			name:          "without support for vector_bytes field",
			weaviateImage: "semitechnologies/weaviate:1.22.5",
		},
		{
			name:          "with support for vector_bytes field",
			weaviateImage: "semitechnologies/weaviate:1.22.6",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create container with a specific Weaviate image
			container, err := docker.StartWeaviate(ctx, tt.weaviateImage)
			require.NoError(t, err)
			defer func() {
				// terminate container when test ends
				require.NoError(t, container.Terminate(ctx))
			}()
			// create client
			cfg := weaviate.Config{
				Host:   container.Endpoint(docker.HTTP),
				Scheme: "http",
				GrpcConfig: &grpc.Config{
					Host: container.Endpoint(docker.GRPC),
				},
			}
			client, err := weaviate.NewClient(cfg)
			require.NoError(t, err)
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
				}
				err = client.Schema().ClassCreator().WithClass(class).Do(ctx)
				require.NoError(t, err)
				// perform batch import
				id := strfmt.UUID("00000000-0000-0000-0000-000000000001")
				vector := models.C11yVector{0.11, 0.22, 0.33, 0.123, -0.900009, -0.0000000001}
				object := &models.Object{
					ID:    id,
					Class: className,
					Properties: map[string]interface{}{
						"name": "some name",
					},
					Vector: vector,
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
				assert.Equal(t, vector, objs[0].Vector)
			})
		})
	}
}
