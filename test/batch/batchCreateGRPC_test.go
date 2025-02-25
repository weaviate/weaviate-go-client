package batch

import (
	"context"
	"testing"

	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/testenv"
)

func TestBatchCreate_gRPC_integration(t *testing.T) {
	err := testenv.SetupLocalWeaviate()
	if err != nil {
		require.Nil(t, err, "failed to start weaviate")
	}

	client := testsuit.CreateTestClient(true)

	t.Run("gRPC batch import", func(t *testing.T) {
		tests := []struct {
			name                string
			className           string
			properties          []map[string]interface{}
			withCrossRefs       bool
			withMultipleVectors bool
		}{
			{
				name:       "all primitive properties",
				className:  "AllProperties",
				properties: testsuit.AllPropertiesDataAsMap(),
			},
			{
				name:       "all primitive properties with nested objects",
				className:  "AllPropertiesWithNested",
				properties: testsuit.AllPropertiesDataWithNestedObjectsAsMap(),
			},
			{
				name:       "all primitive properties with nested array objects",
				className:  "AllPropertiesWithNestedArray",
				properties: testsuit.AllPropertiesDataWithNestedArrayObjectsAsMap(),
			},
			{
				name:          "all primitive properties with cross references (single and multi ref types)",
				className:     "AllPropertiesWithCrossRefs",
				properties:    testsuit.AllPropertiesDataWithCrossReferencesAsMap(),
				withCrossRefs: true,
			},
			{
				name:          "all primitive properties with cross references (single and multi ref types) with nested and nested array objects",
				className:     "AllPropertiesWithCrossRefs",
				properties:    testsuit.AllPropertiesDataWithCrossReferencesWithNestedArrayObjectsAsMap(),
				withCrossRefs: true,
			},
			{
				name:                "all primitive properties with cross references (single and multi ref types) with nested with nested array objects and with multiple vectorizers configuration",
				className:           "AllPropertiesWithCrossRefsAndMultipleVectorizers",
				properties:          testsuit.AllPropertiesDataWithCrossReferencesWithNestedArrayObjectsAsMap(),
				withCrossRefs:       true,
				withMultipleVectors: true,
			},
		}
		for _, tt := range tests {
			className := tt.className
			objects := testsuit.AllPropertiesObjectsWithData(className, tt.properties)
			data := tt.properties

			err := client.Schema().AllDeleter().Do(context.Background())
			require.Nil(t, err)

			testsuit.AllPropertiesSchemaCreate(t, client, className, tt.withCrossRefs, tt.withMultipleVectors)

			batchResultSlice, batchErrSlice := client.Batch().ObjectsBatcher().WithObjects(objects...).Do(context.Background())
			assert.Nil(t, batchErrSlice)
			assert.NotNil(t, batchResultSlice)
			assert.Equal(t, 3, len(batchResultSlice))

			for i := range objects {
				objs, err := client.Data().ObjectsGetter().
					WithID(objects[i].ID.String()).
					WithClassName(objects[i].Class).
					WithVector().
					Do(context.Background())
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

	t.Run("no uuids", func(t *testing.T) {
		className := "NoUUIDs"

		require.Nil(t, client.Schema().ClassDeleter().WithClassName(className).Do(context.Background()))

		class := &models.Class{
			Class: className,
			Properties: []*models.Property{
				{Name: "text", DataType: []string{schema.DataTypeText.String()}},
			},
		}
		require.Nil(t, client.Schema().ClassCreator().WithClass(class).Do(context.Background()))
		objects := []*models.Object{
			{
				Class: className, Properties: map[string]interface{}{"text": "text1"},
			},
		}
		_, batchErrSlice := client.Batch().ObjectsBatcher().WithObjects(objects...).Do(context.Background())
		assert.Nil(t, batchErrSlice)

		objs, err := client.Data().ObjectsGetter().WithClassName(className).Do(context.Background())
		require.NoError(t, err)
		require.Len(t, objs, 1)
		assert.NotEmpty(t, objs[0].ID)
		assert.Equal(t, "text1", objs[0].Properties.(map[string]interface{})["text"])

		require.Nil(t, client.Schema().ClassDeleter().WithClassName(className).Do(context.Background()))
	})

	err = testenv.TearDownLocalWeaviate()
	if err != nil {
		require.Nil(t, err, "failed to tear down weaviate")
	}
}
