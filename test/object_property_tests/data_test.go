//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2023 Weaviate B.V. All rights reserved.
//
//  CONTACT: hello@weaviate.io
//

package object_property_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/testenv"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/schema"
)

func TestObjectProperty_Data(t *testing.T) {
	ctx := context.Background()
	vTrue := true
	vFalse := false

	className := "ObjectPropertyClass"
	id1 := "00000000-0000-0000-0000-000000000001"

	properties := map[string]interface{}{
		"text":          "red",
		"texts":         []string{"red", "blue"},
		"number":        float64(1.1),
		"numbers":       []float64{1.1, 2.2},
		"int":           int64(1),
		"ints":          []int64{1, 2},
		"uuid":          id1,
		"uuids":         []string{id1},
		"date":          "2009-11-01T23:00:00Z",
		"dates":         []string{"2009-11-01T23:00:00Z"},
		"bool":          true,
		"bools":         []bool{true, false},
		"nested_int":    int64(11),
		"nested_number": float64(11.11),
		"nested_text":   "nested text",
		"nested_objects": map[string]interface{}{
			"nested_bool_lvl2":    true,
			"nested_numbers_lvl2": []float64{11.1, 22.1},
		},
		"nested_array_objects": []interface{}{
			map[string]interface{}{
				"nested_bool_lvl2":    true,
				"nested_numbers_lvl2": []float64{111.1, 222.1},
			},
			map[string]interface{}{
				"nested_bool_lvl2":    false,
				"nested_numbers_lvl2": []float64{112.1, 222.1},
			},
		},
	}

	t.Run("up", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			fmt.Print(err.Error())
			t.Fail()
		}
	})

	client := testsuit.CreateTestClient(false)

	t.Run("clean up DB", func(t *testing.T) {
		// clean up DB
		err := client.Schema().AllDeleter().Do(context.Background())
		require.Nil(t, err)
	})

	t.Run("create schema and import object", func(t *testing.T) {
		class := &models.Class{
			Class: className,
			Properties: []*models.Property{
				{
					Name:            "objectProperty",
					DataType:        schema.DataTypeObject.PropString(),
					IndexFilterable: &vTrue,
					IndexSearchable: &vFalse,
					Tokenization:    "",
					NestedProperties: []*models.NestedProperty{
						{
							Name:     "text",
							DataType: schema.DataTypeText.PropString(),
						},
						{
							Name:     "texts",
							DataType: schema.DataTypeTextArray.PropString(),
						},
						{
							Name:     "number",
							DataType: schema.DataTypeNumber.PropString(),
						},
						{
							Name:     "numbers",
							DataType: schema.DataTypeNumberArray.PropString(),
						},
						{
							Name:     "int",
							DataType: schema.DataTypeInt.PropString(),
						},
						{
							Name:     "ints",
							DataType: schema.DataTypeIntArray.PropString(),
						},
						{
							Name:     "date",
							DataType: schema.DataTypeDate.PropString(),
						},
						{
							Name:     "dates",
							DataType: schema.DataTypeDateArray.PropString(),
						},
						{
							Name:     "bool",
							DataType: schema.DataTypeBoolean.PropString(),
						},
						{
							Name:     "bools",
							DataType: schema.DataTypeBooleanArray.PropString(),
						},
						{
							Name:     "uuid",
							DataType: schema.DataTypeUUID.PropString(),
						},
						{
							Name:     "uuids",
							DataType: schema.DataTypeUUIDArray.PropString(),
						},
						{
							Name:            "nested_int",
							DataType:        schema.DataTypeInt.PropString(),
							IndexFilterable: &vTrue,
							IndexSearchable: &vFalse,
							Tokenization:    "",
						},
						{
							Name:            "nested_number",
							DataType:        schema.DataTypeNumber.PropString(),
							IndexFilterable: &vTrue,
							IndexSearchable: &vFalse,
							Tokenization:    "",
						},
						{
							Name:            "nested_text",
							DataType:        schema.DataTypeText.PropString(),
							IndexFilterable: &vTrue,
							IndexSearchable: &vTrue,
							Tokenization:    models.PropertyTokenizationWord,
						},
						{
							Name:            "nested_objects",
							DataType:        schema.DataTypeObject.PropString(),
							IndexFilterable: &vTrue,
							IndexSearchable: &vFalse,
							Tokenization:    "",
							NestedProperties: []*models.NestedProperty{
								{
									Name:            "nested_bool_lvl2",
									DataType:        schema.DataTypeBoolean.PropString(),
									IndexFilterable: &vTrue,
									IndexSearchable: &vFalse,
									Tokenization:    "",
								},
								{
									Name:            "nested_numbers_lvl2",
									DataType:        schema.DataTypeNumberArray.PropString(),
									IndexFilterable: &vTrue,
									IndexSearchable: &vFalse,
									Tokenization:    "",
								},
							},
						},
						{
							Name:            "nested_array_objects",
							DataType:        schema.DataTypeObjectArray.PropString(),
							IndexFilterable: &vTrue,
							IndexSearchable: &vFalse,
							Tokenization:    "",
							NestedProperties: []*models.NestedProperty{
								{
									Name:            "nested_bool_lvl2",
									DataType:        schema.DataTypeBoolean.PropString(),
									IndexFilterable: &vTrue,
									IndexSearchable: &vFalse,
									Tokenization:    "",
								},
								{
									Name:            "nested_numbers_lvl2",
									DataType:        schema.DataTypeNumberArray.PropString(),
									IndexFilterable: &vTrue,
									IndexSearchable: &vFalse,
									Tokenization:    "",
								},
							},
						},
					},
				},
			},
		}
		err := client.Schema().ClassCreator().WithClass(class).Do(ctx)
		require.NoError(t, err)
		// add object
		_, err = client.Data().Creator().
			WithClassName(className).
			WithID(id1).
			WithProperties(properties).
			Do(context.TODO())
		require.Nil(t, err)
		// check if object exists
		res, err := client.Data().ObjectsGetter().WithID(id1).Do(ctx)
		require.Nil(t, err)
		require.Len(t, res, 1)
		props, ok := res[0].Properties.(map[string]interface{})
		require.True(t, ok)
		assert.NotNil(t, props)
		assert.Equal(t, 17, len(props))
		nestedProps, ok := props["nested_objects"].(map[string]interface{})
		require.True(t, ok)
		assert.NotNil(t, nestedProps)
		assert.Equal(t, 2, len(nestedProps))
		err = client.Schema().ClassDeleter().WithClassName(className).Do(ctx)
		require.NoError(t, err)
	})

	t.Run("create complicated object using auto schema", func(t *testing.T) {
		className := "NestedObject"
		_, err := client.Data().Creator().
			WithClassName(className).
			WithID(id1).
			WithProperties(map[string]interface{}{
				"nestedProp": map[string]interface{}{
					"a": int64(1),
					"nested": map[string]interface{}{
						"b": int64(2),
						"nested": map[string]interface{}{
							"c": int64(3),
							"nested": map[string]interface{}{
								"d": int64(4),
								"nested": map[string]interface{}{
									"e": int64(5),
									"nested": map[string]interface{}{
										"f": int64(6),
										"nested": map[string]interface{}{
											"g": int64(7),
											"nested": map[string]interface{}{
												"h": int64(7),
												"nested": map[string]interface{}{
													"i": int64(7),
													"nested": map[string]interface{}{
														"j": int64(8),
														"nested": map[string]interface{}{
															"k": int64(9),
															"nested": map[string]interface{}{
																"l": int64(10),
																"nested": map[string]interface{}{
																	"m": int64(10),
																	"nested": map[string]interface{}{
																		"n": int64(11),
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			}).
			Do(context.TODO())
		require.Nil(t, err)
		// check if object exists
		res, err := client.Data().ObjectsGetter().WithClassName(className).WithID(id1).Do(ctx)
		require.Nil(t, err)
		require.Len(t, res, 1)
		err = client.Schema().ClassDeleter().WithClassName(className).Do(ctx)
		require.NoError(t, err)
	})

	t.Run("load complicated json using auto schema", func(t *testing.T) {
		jsonFile, err := os.Open("./example.json")
		require.NoError(t, err)
		require.NotNil(t, jsonFile)
		byteValue, err := io.ReadAll(jsonFile)
		require.NoError(t, err)
		require.NotNil(t, byteValue)
		var result map[string]interface{}
		json.Unmarshal([]byte(byteValue), &result)
		className := "ComplicatedJson"
		_, err = client.Data().Creator().
			WithClassName(className).
			WithID(id1).
			WithProperties(map[string]interface{}{
				"complicated": result,
			}).
			Do(context.TODO())
		require.Nil(t, err)
		res, err := client.Data().ObjectsGetter().WithClassName(className).WithID(id1).Do(ctx)
		require.Nil(t, err)
		require.Len(t, res, 1)
		props, ok := res[0].Properties.(map[string]interface{})
		require.True(t, ok)
		assert.NotNil(t, props)
		assert.Equal(t, 1, len(props))
		complicated, ok := props["complicated"].(map[string]interface{})
		require.True(t, ok)
		assert.NotNil(t, complicated)
		assert.Equal(t, 2, len(complicated))
		objects, ok := complicated["objects"].([]interface{})
		require.True(t, ok)
		assert.NotNil(t, objects)
		assert.Equal(t, 2, len(objects))
		data, ok := objects[0].(map[string]interface{})
		require.True(t, ok)
		assert.NotNil(t, data)
		assert.Equal(t, 6, len(data))
		properties, ok := data["properties"].(map[string]interface{})
		require.True(t, ok)
		assert.NotNil(t, properties)
		err = client.Schema().ClassDeleter().WithClassName(className).Do(ctx)
		require.NoError(t, err)
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to tear down weaviate: %s", err)
		}
	})
}
