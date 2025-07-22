package search

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/graphql"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/testenv"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/schema"
)

func TestSearch_all_properties(t *testing.T) {
	ctx := context.Background()

	err := testenv.SetupLocalWeaviate()
	if err != nil {
		require.Nil(t, err, "failed to start weaviate")
	}
	client := testsuit.CreateTestClient(true)

	t.Run("clean DB", func(t *testing.T) {
		err := client.Schema().AllDeleter().Do(ctx)
		require.NoError(t, err)
	})

	t.Run("search with all properties", func(t *testing.T) {
		className := "AllPropertiesWithCrossRefsAndMultipleVectorizers"
		t.Run("gRPC batch import", func(t *testing.T) {
			tests := []struct {
				name                string
				className           string
				properties          []map[string]interface{}
				withCrossRefs       bool
				withMultipleVectors bool
			}{
				{
					name:                "all primitive properties with cross references (single and multi ref types) with nested with nested array objects and with multiple vectorizers configuration",
					className:           className,
					properties:          testsuit.AllPropertiesDataWithCrossReferencesWithNestedArrayObjectsAsMap(),
					withCrossRefs:       true,
					withMultipleVectors: true,
				},
			}
			for _, tt := range tests {
				className := tt.className
				objects := testsuit.AllPropertiesObjectsWithData(className, tt.properties)
				data := tt.properties

				testsuit.AllPropertiesSchemaCreate(t, client, className, tt.withCrossRefs, tt.withMultipleVectors)

				batchResultSlice, batchErrSlice := client.Batch().ObjectsBatcher().WithObjects(objects...).Do(ctx)
				assert.Nil(t, batchErrSlice)
				assert.NotNil(t, batchResultSlice)
				assert.Equal(t, 3, len(batchResultSlice))

				for i := range objects {
					objs, err := client.Data().ObjectsGetter().
						WithID(objects[i].ID.String()).
						WithClassName(objects[i].Class).
						WithVector().
						Do(ctx)
					require.NoError(t, err)
					require.Len(t, objs, 1)
					obj := objs[0]
					assert.Equal(t, className, obj.Class)
					props, ok := obj.Properties.(map[string]interface{})
					require.True(t, ok)
					require.NotNil(t, props)
					properties := data[i]
					require.Equal(t, len(props), len(properties))
					for propName := range properties {
						assert.NotNil(t, props[propName])
					}
					if tt.withMultipleVectors {
						assert.Len(t, obj.Vectors, 2)
					}
				}
			}
		})

		t.Run("search", func(t *testing.T) {
			t.Run("find all primitive and array types", func(t *testing.T) {
				props := []string{
					"color", "colors", "author", "authors", "number", "numbers", "int", "ints",
					"uuid", "uuids", "date", "dates", "bool", "bools",
				}
				results, err := client.Experimental().Search().
					WithCollection(className).
					WithProperties(props...).
					Do(ctx)
				require.NoError(t, err)
				assert.Len(t, results, 3)
				for _, res := range results {
					assert.NotEmpty(t, res.ID)
					assert.Equal(t, className, res.Collection)
					require.Len(t, res.Properties, len(props))
					for _, prop := range props {
						assert.NotNil(t, res.Properties[prop])
					}
				}
			})
			t.Run("find all primitive and array and reference types", func(t *testing.T) {
				props := []string{
					"color", "colors", "author", "authors", "number", "numbers", "int", "ints",
					"uuid", "uuids", "date", "dates", "bool", "bools",
				}

				results, err := client.Experimental().Search().
					WithCollection(className).
					WithProperties(props...).
					WithReferences(&graphql.Reference{
						ReferenceProperty: "hasRefClass",
						TargetCollection:  "RefClass",
						Properties:        []string{"category"},
						Metadata:          &graphql.Metadata{ID: true},
					}).
					Do(ctx)
				require.NoError(t, err)
				assert.Len(t, results, 3)
				for _, res := range results {
					assert.NotEmpty(t, res.ID)
					assert.Equal(t, className, res.Collection)
					require.Len(t, res.Properties, len(props))
					for _, prop := range props {
						assert.NotNil(t, res.Properties[prop])
					}
					require.NotEmpty(t, res.References)
					for _, ref := range res.References {
						assert.Equal(t, "hasRefClass", ref.Name)
						require.Len(t, ref.ReferenceProperties, 1)
						for _, refProps := range ref.ReferenceProperties {
							require.Len(t, refProps.Properties, 1)
							assert.NotNil(t, refProps.Properties["category"])
							assert.NotEmpty(t, refProps.Metadata.ID)
						}
					}
				}
			})
			t.Run("find all primitive and array types along with metadata", func(t *testing.T) {
				props := []string{
					"color", "colors", "author", "authors", "number", "numbers", "int", "ints",
					"uuid", "uuids", "date", "dates", "bool", "bools",
				}
				results, err := client.Experimental().Search().
					WithCollection(className).
					WithProperties(props...).
					WithMetadata(&graphql.Metadata{
						ID: true, CreationTimeUnix: true, LastUpdateTimeUnix: true, Vector: true, Vectors: []string{"author_and_colors"},
					}).
					Do(ctx)
				require.NoError(t, err)
				assert.Len(t, results, 3)
				for _, res := range results {
					assert.NotEmpty(t, res.ID)
					assert.Equal(t, className, res.Collection)
					require.Len(t, res.Properties, len(props))
					for _, prop := range props {
						assert.NotNil(t, res.Properties[prop])
					}
					assert.True(t, res.Metadata.CreationTimeUnix > 0)
					assert.True(t, res.Metadata.LastUpdateTimeUnix > 0)
					assert.Empty(t, res.Vector)
					require.NotEmpty(t, res.Vectors)
					require.NotNil(t, res.Vectors["author_and_colors"])
					assert.NotEmpty(t, res.Vectors["author_and_colors"].GetVector())
				}
			})

			t.Run("find all primitive and array types along with metadata with nearText", func(t *testing.T) {
				props := []string{
					"color", "colors", "author", "authors", "number", "numbers", "int", "ints",
					"uuid", "uuids", "date", "dates", "bool", "bools",
				}
				nearText := client.GraphQL().NearTextArgBuilder().
					WithConcepts([]string{"Jenny"}).
					WithTargetVectors("author_and_colors").
					WithCertainty(0.8)
				results, err := client.Experimental().Search().
					WithNearText(nearText).
					WithCollection(className).
					WithProperties(props...).
					WithMetadata(&graphql.Metadata{
						ID: true, Certainty: true, CreationTimeUnix: true, LastUpdateTimeUnix: true, Vector: true, Vectors: []string{"author_and_colors"},
					}).
					Do(ctx)
				require.NoError(t, err)
				assert.Len(t, results, 1)
				for _, res := range results {
					assert.NotEmpty(t, res.ID)
					assert.Equal(t, className, res.Collection)
					require.Len(t, res.Properties, len(props))
					for _, prop := range props {
						assert.NotNil(t, res.Properties[prop])
					}
					assert.True(t, res.Metadata.CreationTimeUnix > 0)
					assert.True(t, res.Metadata.LastUpdateTimeUnix > 0)
					assert.True(t, res.Metadata.Certainty > 0.8)
					assert.Empty(t, res.Vector)
					require.NotEmpty(t, res.Vectors)
					assert.NotNil(t, res.Vectors["author_and_colors"])
				}
			})
			t.Run("find all primitive and array types along with metadata with hybrid", func(t *testing.T) {
				t.Skip("it fails needs investigation")
				props := []string{
					"color", "colors", "author", "authors", "number", "numbers", "int", "ints",
					"uuid", "uuids", "date", "dates", "bool", "bools",
				}

				nearText := client.GraphQL().NearTextArgBuilder().WithConcepts([]string{"Jenny"})
				searches := client.GraphQL().HybridSearchesArgumentBuilder().WithNearText(nearText)
				hybrid := client.GraphQL().HybridArgumentBuilder().
					WithQuery("Jenny").
					WithSearches(searches).
					WithProperties([]string{"author"}).
					WithTargetVectors("author_and_colors")

				results, err := client.Experimental().Search().
					WithHybrid(hybrid).
					WithCollection(className).
					WithProperties(props...).
					WithMetadata(&graphql.Metadata{
						ID: true, Score: true, ExplainScore: true, Vectors: []string{"author_and_colors"},
					}).
					Do(ctx)
				require.NoError(t, err)
				assert.Len(t, results, 3)
				for i, res := range results {
					assert.NotEmpty(t, res.ID)
					assert.Equal(t, className, res.Collection)
					require.Len(t, res.Properties, len(props))
					for _, prop := range props {
						assert.NotNil(t, res.Properties[prop])
					}
					if i == 0 {
						assert.Equal(t, res.Metadata.Score, 1)
					}
					if i == 1 {
						assert.True(t, res.Metadata.Score > 0)
					}
					if i == 2 {
						assert.Equal(t, res.Metadata.Score, 0)
					}
					assert.NotEmpty(t, res.Metadata.ExplainScore)
					assert.Empty(t, res.Vector)
					require.NotEmpty(t, res.Vectors)
					assert.NotNil(t, res.Vectors["author_and_colors"])
				}
			})
		})
	})

	t.Run("multivectors", func(t *testing.T) {
		noVectorizerClass := "NoVectorizerWithMultiVectors"
		t.Run("create class", func(t *testing.T) {
			class := &models.Class{
				Class: noVectorizerClass,
				Properties: []*models.Property{
					{
						Name:     "name",
						DataType: []string{schema.DataTypeText.String()},
					},
				},
				VectorConfig: map[string]models.VectorConfig{
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
			}
			err := client.Schema().ClassCreator().WithClass(class).Do(ctx)
			require.NoError(t, err)
		})
		t.Run("import data", func(t *testing.T) {
			vectors := []models.Vectors{
				{
					"regular": []float32{1, 2},
					"colbert": [][]float32{{0.09, 0.11}, {0.22, 0.33}, {0.33, 0.44}},
				},
				{
					"regular": []float32{3, 4},
					"colbert": [][]float32{{0.19, 0.21}, {0.333, 0.33}, {0.33, 0.33333}},
				},
				{
					"regular": []float32{1, 1},
					"colbert": [][]float32{{0.1111111, 0.1}, {0.1, 0.1}, {0.1, 0.1}},
				},
			}
			objs := make([]*models.Object, len(vectors))
			for i, v := range vectors {
				id := strfmt.UUID(fmt.Sprintf("00000000-0000-0000-0000-00000000000%v", i))
				objs[i] = &models.Object{
					ID:    id,
					Class: noVectorizerClass,
					Properties: map[string]any{
						"name": fmt.Sprintf("name %v", i),
					},
					Vectors: v,
				}
			}
			batchResponse, err := client.Batch().ObjectsBatcher().WithObjects(objs...).Do(ctx)
			require.NoError(t, err)
			assert.NotNil(t, batchResponse)
			assert.Equal(t, len(objs), len(batchResponse))
			// get object
			for i := range vectors {
				id := fmt.Sprintf("00000000-0000-0000-0000-00000000000%v", i)
				res, err := client.Data().ObjectsGetter().
					WithClassName(noVectorizerClass).WithID(id).
					WithVector().
					Do(ctx)
				require.NoError(t, err)
				require.NotEmpty(t, res)
				require.Len(t, res, 1)
				assert.Equal(t, id, res[0].ID.String())
				require.Len(t, res[0].Vectors, 2)
				assert.NotNil(t, res[0].Vectors["regular"])
				assert.NotNil(t, res[0].Vectors["colbert"])
			}
		})
		t.Run("search", func(t *testing.T) {
			bm25 := client.GraphQL().Bm25ArgBuilder().WithQuery("name 1").WithProperties("name")
			results, err := client.Experimental().Search().
				WithBM25(bm25).
				WithCollection(noVectorizerClass).
				WithMetadata(&graphql.Metadata{
					ID: true, Vectors: []string{"regular", "colbert"},
				}).
				Do(ctx)
			require.NoError(t, err)
			assert.Len(t, results, 3)
			for _, res := range results {
				assert.NotEmpty(t, res.ID)
				assert.Equal(t, noVectorizerClass, res.Collection)
				require.Len(t, res.Vectors, 2)
				assert.NotNil(t, res.Vectors["regular"].Vector)
				assert.NotNil(t, res.Vectors["regular"].GetVector())
				assert.NotNil(t, res.Vectors["colbert"].Vector)
				assert.NotNil(t, res.Vectors["colbert"].GetMultiVector())
			}
		})
		t.Run("multi target search", func(t *testing.T) {
			nearVector := client.GraphQL().NearVectorArgBuilder().
				WithTargets(client.GraphQL().MultiTargetArgumentBuilder().Average("regular", "colbert")).
				WithVectorsPerTarget(map[string][]models.Vector{
					"regular": {[]float32{1, 2}},
					"colbert": {[][]float32{{0.09, 0.11}, {0.22, 0.33}, {0.33, 0.44}}},
				})
			results, err := client.Experimental().Search().
				WithNearVector(nearVector).
				WithCollection(noVectorizerClass).
				WithMetadata(&graphql.Metadata{
					ID: true, Vectors: []string{"regular", "colbert"},
				}).
				Do(ctx)
			require.NoError(t, err)
			assert.Len(t, results, 3)
			for _, res := range results {
				assert.NotEmpty(t, res.ID)
				assert.Equal(t, noVectorizerClass, res.Collection)
				require.Len(t, res.Vectors, 2)
				assert.NotNil(t, res.Vectors["regular"].Vector)
				assert.NotNil(t, res.Vectors["regular"].GetVector())
				assert.NotNil(t, res.Vectors["colbert"].Vector)
				assert.NotNil(t, res.Vectors["colbert"].GetMultiVector())
			}
		})
	})

	require.Nil(t, testenv.TearDownLocalWeaviate(), "failed to tear down weaviate")
}
