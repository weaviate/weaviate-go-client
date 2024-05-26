package search

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v4/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/grpc"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/testenv"
)

func TestSearch_all_properties(t *testing.T) {
	ctx := context.Background()
	require.NoError(t, testenv.SetupLocalWeaviate(), "failed to start weaviate")

	port, _, _ := testsuit.GetPortAndAuthPw()
	cfg := weaviate.Config{
		Host:   fmt.Sprintf("localhost:%v", port),
		Scheme: "http",
		GrpcConfig: &grpc.Config{
			Host: "localhost:50051",
		},
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		require.Nil(t, err)
	}

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

			err := client.Schema().AllDeleter().Do(ctx)
			require.Nil(t, err)

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
			results, err := client.Search().
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

			results, err := client.Search().
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
			results, err := client.Search().
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
				assert.NotNil(t, res.Vectors["author_and_colors"])
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
			results, err := client.Search().
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

			results, err := client.Search().
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

	require.Nil(t, testenv.TearDownLocalWeaviate(), "failed to tear down weaviate")
}
