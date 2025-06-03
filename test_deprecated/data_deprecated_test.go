package test_deprecated

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/fault"
)

func TestData_integration_deprecated(t *testing.T) {
	t.Run("up", func(t *testing.T) {
		err := SetupLocalWeaviateDeprecated()
		if err != nil {
			fmt.Print(err.Error())
			t.Fail()
		}
	})

	t.Run("POST /{semanticType}", func(t *testing.T) {
		client := CreateTestClient()

		testsuit.CreateWeaviateTestSchemaFoodDeprecated(t, client)

		propertySchemaT := map[string]string{
			"name":        "Hawaii",
			"description": "Universally accepted to be the best pizza ever created.",
		}
		propertySchemaA := map[string]string{
			"name":        "ChickenSoup",
			"description": "Used by humans when their inferior genetics are attacked by microscopic organisms.",
		}

		wrapperT, errCreateT := client.Data().Creator().
			WithClassName("Pizza").
			WithID("abefd256-8574-442b-9293-9205193737ee").
			WithProperties(propertySchemaT).
			Do(context.Background())
		assert.Nil(t, errCreateT)
		assert.NotNil(t, wrapperT.Object)
		wrapperA, errCreateA := client.Data().Creator().
			WithClassName("Soup").
			WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").
			WithProperties(propertySchemaA).
			Do(context.Background())
		assert.Nil(t, errCreateA)
		assert.NotNil(t, wrapperA.Object)

		objectT, objErrT := client.Data().ObjectsGetter().
			WithID("abefd256-8574-442b-9293-9205193737ee").
			Do(context.Background())
		assert.Nil(t, objErrT)
		objectA, objErrA := client.Data().ObjectsGetter().
			WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").
			Do(context.Background())
		assert.Nil(t, objErrA)

		assert.Equal(t, "Pizza", objectT[0].Class)
		valuesT := objectT[0].Properties.(map[string]interface{})
		assert.Equal(t, "Hawaii", valuesT["name"])
		assert.Equal(t, "Soup", objectA[0].Class)
		valuesA := objectA[0].Properties.(map[string]interface{})
		assert.Equal(t, "ChickenSoup", valuesA["name"])

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("POST vectorizorless class", func(t *testing.T) {
		client := CreateTestClient()

		testsuit.CreateWeaviateTestSchemaWithVectorizorlessClass(t, client)

		propertySchema := map[string]string{
			"name":        "Glazed",
			"description": "The original, and most loved donut covering.",
		}

		vec := []float32{
			0.09271229058504105, 0.16972236335277557, 0.06719677150249481, 0.001922651077620685,
			0.026900049299001694, 0.13298650085926056, 0.02028157375752926, -0.039743948727846146,
			-0.012937345542013645, 0.013409551233053207, -0.10988715291023254, -0.04618397727608681,
			-0.024261055514216423, 0.0663847103714943, 0.004502191673964262, 0.035319264978170395,
			0.10632412880659103, 0.08058158308267593, 0.08017968386411667, -0.02905050292611122,
			0.11437326669692993, 0.00924021378159523, -0.02222306653857231, 0.047553546726703644,
			-0.002701037796214223, 0.15383613109588623, -0.02949700690805912, 0.06906864047050476,
			-0.09484986960887909, 0.06478357315063477, 0.11193037033081055, 0.01517763826996088,
		}

		wrapper, errCreate := client.Data().Creator().
			WithClassName("Donut").
			WithID("66411b32-5c3e-11ec-bf63-0242ac130002").
			WithProperties(propertySchema).
			WithVector(vec).
			Do(context.Background())
		assert.Nil(t, errCreate)
		assert.NotNil(t, wrapper.Object)

		object, objErr := client.Data().ObjectsGetter().
			WithID("66411b32-5c3e-11ec-bf63-0242ac130002").
			WithAdditional("vector").
			Do(context.Background())
		assert.Nil(t, objErr)

		require.True(t, len(object) > 0) // protect against index OOB in next asserstions
		assert.Equal(t, "Donut", object[0].Class)
		valuesV := object[0].Properties.(map[string]interface{})
		assert.Equal(t, "Glazed", valuesV["name"])

		assert.Equal(t, vec, []float32(object[0].Vector))

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("GET /actions /things", func(t *testing.T) {
		client := CreateTestClient()
		testsuit.CreateWeaviateTestSchemaFoodDeprecated(t, client)

		_, errCreate := client.Data().Creator().WithClassName("Pizza").WithProperties(map[string]string{
			"name":        "Margherita",
			"description": "plain",
		}).Do(context.Background())
		assert.Nil(t, errCreate)
		_, errCreate = client.Data().Creator().WithClassName("Pizza").WithProperties(map[string]string{
			"name":        "Pepperoni",
			"description": "meat",
		}).Do(context.Background())
		assert.Nil(t, errCreate)
		_, errCreate = client.Data().Creator().WithClassName("Soup").WithProperties(map[string]string{
			"name":        "Chicken",
			"description": "meat",
		}).Do(context.Background())
		assert.Nil(t, errCreate)
		_, errCreate = client.Data().Creator().WithClassName("Soup").WithProperties(map[string]string{
			"name":        "Tofu",
			"description": "vegetarian",
		}).Do(context.Background())
		assert.Nil(t, errCreate)

		objectT, objErrT := client.Data().ObjectsGetter().Do(context.Background())
		assert.Nil(t, objErrT)
		assert.Equal(t, 4, len(objectT))

		objectT2, objectErrT2 := client.Data().ObjectsGetter().WithAdditional("vector").WithLimit(1).Do(context.Background())
		assert.Nil(t, objectErrT2)
		assert.Equal(t, 1, len(objectT2))
		objectA2, objErrA2 := client.Data().ObjectsGetter().WithLimit(1).Do(context.Background())
		assert.Nil(t, objErrA2)
		assert.Equal(t, 1, len(objectA2))

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("GET underscore properties", func(t *testing.T) {
		client := CreateTestClient()

		testsuit.CreateWeaviateTestSchemaFoodDeprecated(t, client)

		propertySchemaT := map[string]string{
			"name":        "Hawaii",
			"description": "Universally accepted to be the best pizza ever created.",
		}
		propertySchemaA := map[string]string{
			"name":        "ChickenSoup",
			"description": "Used by humans when their inferior genetics are attacked by microscopic organisms.",
		}
		_, errCreateT := client.Data().Creator().WithClassName("Pizza").
			WithID("abefd256-8574-442b-9293-9205193737ee").
			WithProperties(propertySchemaT).Do(context.Background())
		assert.Nil(t, errCreateT)
		_, errCreateA := client.Data().Creator().WithClassName("Soup").
			WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").
			WithProperties(propertySchemaA).Do(context.Background())
		assert.Nil(t, errCreateA)

		// THINGS
		objectT, objErrT := client.Data().ObjectsGetter().
			WithID("abefd256-8574-442b-9293-9205193737ee").Do(context.Background())
		assert.Nil(t, objErrT)
		assert.Nil(t, objectT[0].Additional["classification"])
		assert.Nil(t, objectT[0].Additional["nearestNeighbors"])
		assert.Nil(t, objectT[0].Additional["featureProjection"])
		assert.Nil(t, objectT[0].Vector)
		assert.Nil(t, objectT[0].Additional["interpretation"])

		objectT, objErrT = client.Data().ObjectsGetter().
			WithID("abefd256-8574-442b-9293-9205193737ee").
			WithAdditional("interpretation").Do(context.Background())
		assert.Nil(t, objErrT)
		assert.Nil(t, objectT[0].Additional["classification"])
		assert.Nil(t, objectT[0].Additional["nearestNeighbors"])
		assert.Nil(t, objectT[0].Additional["featureProjection"])
		assert.Nil(t, objectT[0].Vector)
		assert.NotNil(t, objectT[0].Additional["interpretation"])

		objectT, objErrT = client.Data().ObjectsGetter().
			WithID("abefd256-8574-442b-9293-9205193737ee").
			WithAdditional("interpretation").
			WithAdditional("classification").
			WithAdditional("nearestNeighbors").
			WithVector().
			Do(context.Background())
		assert.Nil(t, objErrT)
		assert.Nil(t, objectT[0].Additional["classification"]) // Is nil because no classifications was executed
		assert.NotNil(t, objectT[0].Additional["nearestNeighbors"])
		assert.Nil(t, objectT[0].Additional["featureProjection"]) // Is nil because feature projection is not possible on non list request
		assert.NotNil(t, objectT[0].Vector)
		assert.NotNil(t, objectT[0].Additional["interpretation"])

		// ACTIONS
		objectA, objErrA := client.Data().ObjectsGetter().
			WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").Do(context.Background())
		assert.Nil(t, objErrA)
		assert.Nil(t, objectA[0].Additional["classification"])
		assert.Nil(t, objectA[0].Additional["nearestNeighbors"])
		assert.Nil(t, objectA[0].Additional["featureProjection"])
		assert.Nil(t, objectA[0].Vector)
		assert.Nil(t, objectA[0].Additional["interpretation"])

		objectA, objErrA = client.Data().ObjectsGetter().
			WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").
			WithAdditional("interpretation").Do(context.Background())
		assert.Nil(t, objErrA)
		assert.Nil(t, objectA[0].Additional["classification"])
		assert.Nil(t, objectA[0].Additional["nearestNeighbors"])
		assert.Nil(t, objectA[0].Additional["featureProjection"])
		assert.Nil(t, objectA[0].Vector)
		assert.NotNil(t, objectA[0].Additional["interpretation"])

		objectA, objErrA = client.Data().ObjectsGetter().
			WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").
			WithAdditional("interpretation").
			WithAdditional("classification").
			WithAdditional("nearestNeighbors").
			WithAdditional("featureProjection").
			WithVector().
			Do(context.Background())
		assert.Nil(t, objErrT)
		assert.Nil(t, objectT[0].Additional["classification"]) // Is nil because no classifications was executed
		assert.NotNil(t, objectT[0].Additional["nearestNeighbors"])
		assert.Nil(t, objectT[0].Additional["featureProjection"]) // Is nil because feature projection is not possible on non list request
		assert.NotNil(t, objectT[0].Vector)
		assert.NotNil(t, objectT[0].Additional["interpretation"])

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("DELETE /{type}", func(t *testing.T) {
		client := CreateTestClient()

		testsuit.CreateWeaviateTestSchemaFoodDeprecated(t, client)

		propertySchemaT := map[string]string{
			"name":        "Hawaii",
			"description": "Universally accepted to be the best pizza ever created.",
		}
		propertySchemaA := map[string]string{
			"name":        "ChickenSoup",
			"description": "Used by humans when their inferior genetics are attacked by microscopic organisms.",
		}
		_, errCreateT := client.Data().Creator().
			WithClassName("Pizza").
			WithID("abefd256-8574-442b-9293-9205193737ee").
			WithProperties(propertySchemaT).
			Do(context.Background())
		assert.Nil(t, errCreateT)
		_, errCreateA := client.Data().Creator().
			WithClassName("Soup").WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").
			WithProperties(propertySchemaA).
			Do(context.Background())
		assert.Nil(t, errCreateA)

		// THINGS
		deleteErrT := client.Data().Deleter().
			WithID("abefd256-8574-442b-9293-9205193737ee").
			Do(context.Background())
		assert.Nil(t, deleteErrT)
		_, getErrT := client.Data().ObjectsGetter().
			WithID("abefd256-8574-442b-9293-9205193737ee").
			Do(context.Background())
		statusCodeErrorT := getErrT.(*fault.WeaviateClientError)
		assert.Equal(t, 404, statusCodeErrorT.StatusCode)

		deleteErrA := client.Data().Deleter().
			WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").
			Do(context.Background())
		assert.Nil(t, deleteErrA)
		_, getErrA := client.Data().ObjectsGetter().
			WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").
			Do(context.Background())
		statusCodeErrorA := getErrA.(*fault.WeaviateClientError)
		assert.Equal(t, 404, statusCodeErrorA.StatusCode)

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("HEAD /{id}", func(t *testing.T) {
		client := CreateTestClient()

		testsuit.CreateWeaviateTestSchemaFoodDeprecated(t, client)

		objTID := "abefd256-8574-442b-9293-9205193737ee"
		propertySchemaT := map[string]string{
			"name":        "Hawaii",
			"description": "Universally accepted to be the best pizza ever created.",
		}
		_, errCreateT := client.Data().Creator().
			WithClassName("Pizza").
			WithID(objTID).
			WithProperties(propertySchemaT).
			Do(context.Background())
		assert.Nil(t, errCreateT)

		// Check object which exists
		exists, checkErrT := client.Data().Checker().
			WithID(objTID).
			Do(context.Background())
		assert.Nil(t, checkErrT)
		assert.True(t, exists)
		// Double check that it actually exists in DB
		objT, afterCheckErrT := client.Data().ObjectsGetter().
			WithID(objTID).
			Do(context.Background())
		assert.NotNil(t, objT)
		assert.Nil(t, afterCheckErrT)

		// Delete object
		deleteErrT := client.Data().Deleter().
			WithID(objTID).
			Do(context.Background())
		assert.Nil(t, deleteErrT)
		// Verify that the object has been actually deleted
		_, getErrT := client.Data().ObjectsGetter().
			WithID(objTID).
			Do(context.Background())
		statusCodeErrorT := getErrT.(*fault.WeaviateClientError)
		assert.Equal(t, 404, statusCodeErrorT.StatusCode)

		// Check object which doesn't exist
		exists, checkErrT = client.Data().Checker().
			WithID(objTID).
			Do(context.Background())
		assert.Nil(t, checkErrT)
		assert.False(t, exists)
		// Double check that it really doesn't exits in DB
		_, getErrT = client.Data().ObjectsGetter().
			WithID(objTID).
			Do(context.Background())
		statusCodeErrorT = getErrT.(*fault.WeaviateClientError)
		assert.Equal(t, 404, statusCodeErrorT.StatusCode)

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("PUT /{type}/{id}", func(t *testing.T) {
		// PUT replaces the object fully
		client := CreateTestClient()

		testsuit.CreateWeaviateTestSchemaFoodDeprecated(t, client)

		propertySchemaT := map[string]string{
			"name":        "Random",
			"description": "Missing description",
		}
		propertySchemaA := map[string]string{
			"name":        "water",
			"description": "missing description",
		}
		_, errCreateT := client.Data().Creator().
			WithClassName("Pizza").
			WithID("abefd256-8574-442b-9293-9205193737ee").
			WithProperties(propertySchemaT).
			Do(context.Background())
		assert.Nil(t, errCreateT)
		_, errCreateA := client.Data().Creator().
			WithClassName("Soup").
			WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").
			WithProperties(propertySchemaA).
			Do(context.Background())
		assert.Nil(t, errCreateA)

		propertySchemaT = map[string]string{
			"name":        "Hawaii",
			"description": "Universally accepted to be the best pizza ever created.",
		}
		updateErrT := client.Data().Updater().
			WithID("abefd256-8574-442b-9293-9205193737ee").
			WithClassName("Pizza").
			WithProperties(propertySchemaT).
			Do(context.Background())
		assert.Nil(t, updateErrT)

		propertySchemaA = map[string]string{
			"name":        "ChickenSoup",
			"description": "Used by humans when their inferior genetics are attacked by microscopic organisms.",
		}
		updateErrA := client.Data().Updater().
			WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").
			WithClassName("Soup").
			WithProperties(propertySchemaA).
			Do(context.Background())
		assert.Nil(t, updateErrA)

		things, getErrT := client.Data().ObjectsGetter().
			WithID("abefd256-8574-442b-9293-9205193737ee").
			Do(context.Background())
		assert.Nil(t, getErrT)
		valuesT := things[0].Properties.(map[string]interface{})
		assert.Equal(t, propertySchemaT["description"], valuesT["description"])
		assert.Equal(t, propertySchemaT["name"], valuesT["name"])

		actions, getErrT := client.Data().ObjectsGetter().
			WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").
			Do(context.Background())
		assert.Nil(t, getErrT)
		valuesA := actions[0].Properties.(map[string]interface{})
		assert.Equal(t, propertySchemaA["description"], valuesA["description"])
		assert.Equal(t, propertySchemaA["name"], valuesA["name"])

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("PATCH(merge) /{type}/{id}", func(t *testing.T) {
		// PATCH merges the new object with the existing object
		client := CreateTestClient()

		testsuit.CreateWeaviateTestSchemaFoodDeprecated(t, client)

		propertySchemaT := map[string]string{
			"name":        "Hawaii",
			"description": "Missing description",
		}
		propertySchemaA := map[string]string{
			"name":        "ChickenSoup",
			"description": "missing description",
		}
		_, errCreateT := client.Data().Creator().
			WithClassName("Pizza").
			WithID("abefd256-8574-442b-9293-9205193737ee").
			WithProperties(propertySchemaT).
			Do(context.Background())
		assert.Nil(t, errCreateT)
		_, errCreateA := client.Data().Creator().
			WithClassName("Soup").
			WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").
			WithProperties(propertySchemaA).
			Do(context.Background())
		assert.Nil(t, errCreateA)

		propertySchemaT = map[string]string{
			"description": "Universally accepted to be the best pizza ever created.",
		}
		updateErrT := client.Data().Updater().
			WithID("abefd256-8574-442b-9293-9205193737ee").
			WithClassName("Pizza").
			WithProperties(propertySchemaT).
			WithMerge().
			Do(context.Background())
		assert.Nil(t, updateErrT)

		propertySchemaA = map[string]string{
			"description": "Used by humans when their inferior genetics are attacked by microscopic organisms.",
		}
		updateErrA := client.Data().Updater().
			WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").
			WithClassName("Soup").
			WithProperties(propertySchemaA).
			WithMerge().
			Do(context.Background())
		assert.Nil(t, updateErrA)

		things, getErrT := client.Data().ObjectsGetter().
			WithID("abefd256-8574-442b-9293-9205193737ee").
			Do(context.Background())
		assert.Nil(t, getErrT)
		valuesT := things[0].Properties.(map[string]interface{})
		assert.Equal(t, propertySchemaT["description"], valuesT["description"])
		assert.Equal(t, "Hawaii", valuesT["name"])

		actions, getErrT := client.Data().ObjectsGetter().
			WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").
			Do(context.Background())
		assert.Nil(t, getErrT)
		valuesA := actions[0].Properties.(map[string]interface{})
		assert.Equal(t, propertySchemaA["description"], valuesA["description"])
		assert.Equal(t, "ChickenSoup", valuesA["name"])

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("POST /{type}/validate", func(t *testing.T) {
		client := CreateTestClient()

		testsuit.CreateWeaviateTestSchemaFoodDeprecated(t, client)

		propertySchemaT := map[string]string{
			"name":        "Hawaii",
			"description": "Universally accepted to be the best pizza ever created.",
		}
		propertySchemaA := map[string]string{
			"name":        "ChickenSoup",
			"description": "Used by humans when their inferior genetics are attacked by microscopic organisms.",
		}

		errValidateT := client.Data().Validator().
			WithClassName("Pizza").
			WithID("abefd256-8574-442b-9293-9205193737ee").
			WithProperties(propertySchemaT).
			Do(context.Background())
		assert.Nil(t, errValidateT)

		errValidateA := client.Data().Validator().
			WithClassName("Soup").
			WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").
			WithProperties(propertySchemaA).
			Do(context.Background())
		assert.Nil(t, errValidateA)

		propertySchemaT["test"] = "not existing property"
		errValidateT = client.Data().Validator().
			WithClassName("Pizza").
			WithID("abefd256-8574-442b-9293-9205193737ee").
			WithProperties(propertySchemaT).
			Do(context.Background())
		assert.NotNil(t, errValidateT)

		propertySchemaA["test"] = "not existing property"
		errValidateA = client.Data().Validator().
			WithClassName("Soup").
			WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").
			WithProperties(propertySchemaA).
			Do(context.Background())
		assert.NotNil(t, errValidateA)

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := TearDownLocalWeaviateDeprecated()
		if err != nil {
			fmt.Print(err.Error())
			t.Fail()
		}
	})
}
