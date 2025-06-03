package data

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/data/replication"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/fault"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/testenv"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/schema"
)

func createWeaviateTestSchemaUuid(t *testing.T, client *weaviate.Client) {
	schemaClassThing := &models.Class{
		Class:               "UserUUIDTest",
		Description:         "A useruuid",
		InvertedIndexConfig: &models.InvertedIndexConfig{IndexTimestamps: true},
	}

	errT := client.Schema().ClassCreator().WithClass(schemaClassThing).Do(context.Background())
	assert.Nil(t, errT)
	uuidPropertyDataType := []string{"uuid"}
	uuidArrayPropertyDataType := []string{"uuid[]"}
	UserUUIDTest := &models.Property{
		DataType:    uuidPropertyDataType,
		Description: "uuid",
		Name:        "uuid",
	}
	uuidArrayProperty := &models.Property{
		DataType:    uuidArrayPropertyDataType,
		Description: "uuid array",
		Name:        "uuidArray",
	}

	propErrT1 := client.Schema().PropertyCreator().WithClassName("UserUUIDTest").WithProperty(UserUUIDTest).Do(context.Background())
	assert.Nil(t, propErrT1)
	propErrA1 := client.Schema().PropertyCreator().WithClassName("UserUUIDTest").WithProperty(uuidArrayProperty).Do(context.Background())
	assert.Nil(t, propErrA1)
}

func TestData_uuid(t *testing.T) {
	t.Run("up", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			fmt.Print(err.Error())
			t.Fail()
		}
	})

	t.Run("POST /{semanticType}", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)

		createWeaviateTestSchemaUuid(t, client)

		propertySchemaT := map[string]string{
			"uuid": "565da3b6-60b3-40e5-ba21-e6bfe5dbba91",
		}
		propertySchemaA := map[string][]string{
			"uuidArray": {"565da3b6-60b3-40e5-ba21-e6bfe5dbba92", "565da3b6-60b3-40e5-ba21-e6bfe5dbba93"},
		}

		wrapperT, errCreateT := client.Data().Creator().
			WithClassName("UserUUIDTest").
			WithID("abefd256-8574-442b-9293-9205193737ee").
			WithProperties(propertySchemaT).
			Do(context.Background())
		assert.Nil(t, errCreateT)
		assert.NotNil(t, wrapperT.Object)
		wrapperA, errCreateA := client.Data().Creator().
			WithClassName("UserUUIDTest").
			WithID("abefd256-8574-442b-9293-9205193737ef").
			WithProperties(propertySchemaA).
			Do(context.Background())
		assert.Nil(t, errCreateA)
		assert.NotNil(t, wrapperA.Object)

		objectT, objErrT := client.Data().ObjectsGetter().
			WithClassName("UserUUIDTest").
			WithID("abefd256-8574-442b-9293-9205193737ee").
			Do(context.Background())
		assert.Nil(t, objErrT)
		objectA, objErrA := client.Data().ObjectsGetter().
			WithClassName("UserUUIDTest").
			WithID("abefd256-8574-442b-9293-9205193737ef").
			Do(context.Background())
		assert.Nil(t, objErrA)

		assert.Equal(t, "UserUUIDTest", objectT[0].Class)
		valuesT := objectT[0].Properties.(map[string]interface{})
		assert.Equal(t, "565da3b6-60b3-40e5-ba21-e6bfe5dbba91", valuesT["uuid"])
		assert.Equal(t, "UserUUIDTest", objectA[0].Class)
		valuesA := objectA[0].Properties.(map[string]interface{})
		var compData []interface{}
		compData = append(compData, "565da3b6-60b3-40e5-ba21-e6bfe5dbba92")
		compData = append(compData, "565da3b6-60b3-40e5-ba21-e6bfe5dbba93")
		assert.Equal(t, compData, valuesA["uuidArray"])
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			fmt.Print(err.Error())
			t.Fail()
		}
	})
}

func TestData_integration(t *testing.T) {
	t.Run("up", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			fmt.Print(err.Error())
			t.Fail()
		}
	})

	t.Run("POST /{semanticType}", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)

		testsuit.CreateWeaviateTestSchemaFood(t, client)

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
			WithClassName("Pizza").
			WithID("abefd256-8574-442b-9293-9205193737ee").
			Do(context.Background())
		assert.Nil(t, objErrT)
		objectA, objErrA := client.Data().ObjectsGetter().
			WithClassName("Soup").
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
		client := testsuit.CreateTestClient(false)

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
			WithClassName("Donut").
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
		client := testsuit.CreateTestClient(false)
		testsuit.CreateWeaviateTestSchemaFood(t, client)
		// create two pizzas and two soups
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

		pizzas, pizzasErr := client.Data().ObjectsGetter().WithClassName("Pizza").Do(context.Background())
		assert.Nil(t, pizzasErr)
		assert.Equal(t, 2, len(pizzas))
		soups, soupsErr := client.Data().ObjectsGetter().WithClassName("Soup").Do(context.Background())
		assert.Nil(t, soupsErr)
		assert.Equal(t, 2, len(soups))

		var firstPizzaID string // save for reuse with cursor
		pizza, pizzaErr := client.Data().ObjectsGetter().WithClassName("Pizza").WithLimit(1).Do(context.Background())
		assert.Nil(t, pizzaErr)
		assert.Equal(t, 1, len(pizza))
		firstPizzaID = pizza[0].ID.String()
		soup, soupErr := client.Data().ObjectsGetter().WithClassName("Soup").WithLimit(1).Do(context.Background())
		assert.Nil(t, soupErr)
		assert.Equal(t, 1, len(soup))

		pizzas, pizzasErr = client.Data().ObjectsGetter().WithClassName("Pizza").
			WithLimit(10).WithAfter(firstPizzaID).Do(context.Background())
		assert.Nil(t, pizzasErr)
		assert.Equal(t, 1, len(pizzas)) // only the other pizza should be left

		secondPizzaID := pizzas[0].ID.String() // save the id of the second pizza

		// WithOffset(0) should work the same as not using an offset
		pizzas_offset, pizzaErr := client.Data().ObjectsGetter().WithClassName("Pizza").WithOffset(0).Do(context.Background())
		assert.Nil(t, pizzaErr)
		assert.Equal(t, 2, len(pizzas_offset))

		// WithOffset(1) should only return the second pizza
		pizzas_offset, pizzaErr = client.Data().ObjectsGetter().WithClassName("Pizza").WithOffset(1).Do(context.Background())
		assert.Nil(t, pizzaErr)
		assert.Equal(t, 1, len(pizzas_offset))
		assert.Equal(t, secondPizzaID, pizzas_offset[0].ID.String())

		// WithOffset(5) should not return any pizzas
		pizzas_offset, pizzaErr = client.Data().ObjectsGetter().WithClassName("Pizza").WithOffset(5).Do(context.Background())
		assert.Nil(t, pizzaErr)
		assert.Equal(t, 0, len(pizzas_offset))

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("GET underscore properties", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)

		testsuit.CreateWeaviateTestSchemaFood(t, client)

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
			WithClassName("Pizza").
			WithID("abefd256-8574-442b-9293-9205193737ee").Do(context.Background())
		assert.Nil(t, objErrT)
		assert.Nil(t, objectT[0].Additional["classification"])
		assert.Nil(t, objectT[0].Additional["nearestNeighbors"])
		assert.Nil(t, objectT[0].Additional["featureProjection"])
		assert.Nil(t, objectT[0].Vector)
		assert.Nil(t, objectT[0].Additional["interpretation"])

		objectT, objErrT = client.Data().ObjectsGetter().
			WithClassName("Pizza").
			WithID("abefd256-8574-442b-9293-9205193737ee").
			WithAdditional("interpretation").Do(context.Background())
		assert.Nil(t, objErrT)
		assert.Nil(t, objectT[0].Additional["classification"])
		assert.Nil(t, objectT[0].Additional["nearestNeighbors"])
		assert.Nil(t, objectT[0].Additional["featureProjection"])
		assert.Nil(t, objectT[0].Vector)
		assert.NotNil(t, objectT[0].Additional["interpretation"])

		objectT, objErrT = client.Data().ObjectsGetter().
			WithClassName("Pizza").
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
			WithClassName("Soup").
			WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").Do(context.Background())
		assert.Nil(t, objErrA)
		assert.Nil(t, objectA[0].Additional["classification"])
		assert.Nil(t, objectA[0].Additional["nearestNeighbors"])
		assert.Nil(t, objectA[0].Additional["featureProjection"])
		assert.Nil(t, objectA[0].Vector)
		assert.Nil(t, objectA[0].Additional["interpretation"])

		objectA, objErrA = client.Data().ObjectsGetter().
			WithClassName("Soup").
			WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").
			WithAdditional("interpretation").Do(context.Background())
		assert.Nil(t, objErrA)
		assert.Nil(t, objectA[0].Additional["classification"])
		assert.Nil(t, objectA[0].Additional["nearestNeighbors"])
		assert.Nil(t, objectA[0].Additional["featureProjection"])
		assert.Nil(t, objectA[0].Vector)
		assert.NotNil(t, objectA[0].Additional["interpretation"])

		objectA, objErrA = client.Data().ObjectsGetter().
			WithClassName("Soup").
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

	t.Run("GET /objects/{className}/{id}", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		testsuit.CreateTestSchemaAndData(t, client)

		resp, err := client.Data().ObjectsGetter().
			WithClassName("Pizza").
			WithID("5b6a08ba-1d46-43aa-89cc-8b070790c6f2").
			Do(context.Background())
		assert.Nil(t, err)
		assert.Len(t, resp, 1)
		assert.Equal(t, "Doener", resp[0].Properties.(map[string]interface{})["name"])

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("DELETE /{type}", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)

		testsuit.CreateWeaviateTestSchemaFood(t, client)

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
			WithClassName("Pizza").
			WithID("abefd256-8574-442b-9293-9205193737ee").
			Do(context.Background())
		assert.Nil(t, deleteErrT)
		_, getErrT := client.Data().ObjectsGetter().
			WithClassName("Pizza").
			WithID("abefd256-8574-442b-9293-9205193737ee").
			Do(context.Background())
		statusCodeErrorT := getErrT.(*fault.WeaviateClientError)
		assert.Equal(t, 404, statusCodeErrorT.StatusCode)

		deleteErrA := client.Data().Deleter().
			WithClassName("Soup").
			WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").
			Do(context.Background())
		assert.Nil(t, deleteErrA)
		_, getErrA := client.Data().ObjectsGetter().
			WithClassName("Soup").
			WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").
			Do(context.Background())
		statusCodeErrorA := getErrA.(*fault.WeaviateClientError)
		assert.Equal(t, 404, statusCodeErrorA.StatusCode)

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("DELETE /objects/{className}/{id}", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		testsuit.CreateTestSchemaAndData(t, client)

		err := client.Data().Deleter().
			WithClassName("Pizza").
			WithID("5b6a08ba-1d46-43aa-89cc-8b070790c6f2").
			Do(context.Background())
		assert.Nil(t, err)

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("HEAD /{id}", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)

		testsuit.CreateWeaviateTestSchemaFood(t, client)

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
			WithClassName("Pizza").
			WithID(objTID).
			Do(context.Background())
		assert.Nil(t, checkErrT)
		assert.True(t, exists)
		// Double check that it actually exists in DB
		objT, afterCheckErrT := client.Data().ObjectsGetter().
			WithClassName("Pizza").
			WithID(objTID).
			Do(context.Background())
		assert.NotNil(t, objT)
		assert.Nil(t, afterCheckErrT)

		// Delete object
		deleteErrT := client.Data().Deleter().
			WithClassName("Pizza").
			WithID(objTID).
			Do(context.Background())
		assert.Nil(t, deleteErrT)
		// Verify that the object has been actually deleted
		_, getErrT := client.Data().ObjectsGetter().
			WithClassName("Pizza").
			WithID(objTID).
			Do(context.Background())
		statusCodeErrorT := getErrT.(*fault.WeaviateClientError)
		assert.Equal(t, 404, statusCodeErrorT.StatusCode)

		// Check object which doesn't exist
		exists, checkErrT = client.Data().Checker().
			WithClassName("Pizza").
			WithID(objTID).
			Do(context.Background())
		assert.Nil(t, checkErrT)
		assert.False(t, exists)
		// Double check that it really doesn't exits in DB
		_, getErrT = client.Data().ObjectsGetter().
			WithClassName("Pizza").
			WithID(objTID).
			Do(context.Background())
		statusCodeErrorT = getErrT.(*fault.WeaviateClientError)
		assert.Equal(t, 404, statusCodeErrorT.StatusCode)

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("HEAD /objects/{className}/{id}", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		testsuit.CreateTestSchemaAndData(t, client)

		exists, err := client.Data().Checker().
			WithClassName("Pizza").
			WithID("5b6a08ba-1d46-43aa-89cc-8b070790c6f2").
			Do(context.Background())
		assert.Nil(t, err)
		assert.True(t, exists)

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("PUT /{type}/{id}", func(t *testing.T) {
		// PUT replaces the object fully
		client := testsuit.CreateTestClient(false)

		testsuit.CreateWeaviateTestSchemaFood(t, client)

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
			WithClassName("Pizza").
			WithID("abefd256-8574-442b-9293-9205193737ee").
			Do(context.Background())
		assert.Nil(t, getErrT)
		valuesT := things[0].Properties.(map[string]interface{})
		assert.Equal(t, propertySchemaT["description"], valuesT["description"])
		assert.Equal(t, propertySchemaT["name"], valuesT["name"])

		actions, getErrT := client.Data().ObjectsGetter().
			WithClassName("Soup").
			WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").
			Do(context.Background())
		assert.Nil(t, getErrT)
		valuesA := actions[0].Properties.(map[string]interface{})
		assert.Equal(t, propertySchemaA["description"], valuesA["description"])
		assert.Equal(t, propertySchemaA["name"], valuesA["name"])

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("PUT /objects/{className}/{id}", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		testsuit.CreateTestSchemaAndData(t, client)

		id := "5b6a08ba-1d46-43aa-89cc-8b070790c6f2"
		props := map[string]interface{}{
			"name":        "Margherita",
			"description": "Invented in honor of the Queen of Italy, Margherita of Savoy",
			"price":       5.12,
		}

		err := client.Data().Updater().
			WithClassName("Pizza").
			WithID(id).
			WithProperties(props).
			Do(context.Background())
		assert.Nil(t, err)

		resp, err := client.Data().ObjectsGetter().
			WithClassName("Pizza").
			WithID(id).
			Do(context.Background())
		assert.Nil(t, err)
		assert.Len(t, resp, 1)
		assert.EqualValues(t, props, resp[0].Properties)

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("PATCH(merge) /{type}/{id}", func(t *testing.T) {
		// PATCH merges the new object with the existing object
		client := testsuit.CreateTestClient(false)

		testsuit.CreateWeaviateTestSchemaFood(t, client)

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
			WithClassName("Pizza").
			WithID("abefd256-8574-442b-9293-9205193737ee").
			Do(context.Background())
		assert.Nil(t, getErrT)
		valuesT := things[0].Properties.(map[string]interface{})
		assert.Equal(t, propertySchemaT["description"], valuesT["description"])
		assert.Equal(t, "Hawaii", valuesT["name"])

		actions, getErrT := client.Data().ObjectsGetter().
			WithClassName("Soup").
			WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").
			Do(context.Background())
		assert.Nil(t, getErrT)
		valuesA := actions[0].Properties.(map[string]interface{})
		assert.Equal(t, propertySchemaA["description"], valuesA["description"])
		assert.Equal(t, "ChickenSoup", valuesA["name"])

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("PATCH(with vector) /object/{clasName}/{id}", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)

		testsuit.CreateWeaviateTestSchemaWithVectorizorlessClass(t, client)

		id := "66411b32-5c3e-11ec-bf63-0242ac130002"
		propertySchema := map[string]string{
			"name":        "Glazed",
			"description": "The original, and most loved donut covering.",
		}

		vecA := []float32{
			0.09271229058504105, 0.16972236335277557, 0.06719677150249481, 0.001922651077620685,
			0.026900049299001694, 0.13298650085926056, 0.02028157375752926, -0.039743948727846146,
			-0.012937345542013645, 0.013409551233053207, -0.10988715291023254, -0.04618397727608681,
			-0.024261055514216423, 0.0663847103714943, 0.004502191673964262, 0.035319264978170395,
			0.10632412880659103, 0.08058158308267593, 0.08017968386411667, -0.02905050292611122,
			0.11437326669692993, 0.00924021378159523, -0.02222306653857231, 0.047553546726703644,
		}
		wrapper, errCreate := client.Data().Creator().
			WithClassName("Donut").
			WithID("66411b32-5c3e-11ec-bf63-0242ac130002").
			WithProperties(propertySchema).
			WithVector(vecA).
			Do(context.Background())
		assert.Nil(t, errCreate)
		assert.NotNil(t, wrapper.Object)

		vecT := []float32{
			0.11437326669692993, 0.16972236335277557, 0.06719677150249481, 0.001922651077620685,
			0.026900049299001694, 0.13298650085926056, 0.02028157375752926, -0.039743948727846146,
			-0.012937345542013645, 0.013409551233053207, -0.10988715291023254, -0.04618397727608681,
			-0.024261055514216423, 0.0663847103714943, 0.004502191673964262, 0.035319264978170395,
			0.10632412880659103, 0.08058158308267593, 0.08017968386411667, -0.02905050292611122,
			0.00924021378159523, 0.11437326669692993, -0.02222306653857231, 0.047553546726703644,
		}
		errUpdate := client.Data().Updater().WithClassName("Donut").
			WithID("66411b32-5c3e-11ec-bf63-0242ac130002").
			WithProperties(propertySchema).
			WithVector(vecT).
			Do(context.Background())
		assert.Nil(t, errUpdate)

		object, objErr := client.Data().ObjectsGetter().
			WithClassName("Donut").
			WithID(id).
			WithAdditional("vector").
			Do(context.Background())
		assert.Nil(t, objErr)
		assert.Len(t, object, 1)
		assert.Equal(t, vecT, []float32(object[0].Vector))
	})

	t.Run("PATCH /objects/{className}/{id}", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		testsuit.CreateTestSchemaAndData(t, client)

		id := "5b6a08ba-1d46-43aa-89cc-8b070790c6f2"
		inputProps := map[string]interface{}{
			"description": "Kebap, in pizza form",
		}
		expectedProps := map[string]interface{}{
			"name":        "Doener",
			"description": "Kebap, in pizza form",
			"price":       1.4,
		}

		err := client.Data().Updater().
			WithClassName("Pizza").
			WithID(id).
			WithProperties(inputProps).
			WithMerge().
			Do(context.Background())
		assert.Nil(t, err)

		resp, err := client.Data().ObjectsGetter().
			WithClassName("Pizza").
			WithID(id).
			Do(context.Background())
		assert.Nil(t, err)
		assert.Len(t, resp, 1)
		assert.EqualValues(t, expectedProps, resp[0].Properties)

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("POST /{type}/validate", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)

		testsuit.CreateWeaviateTestSchemaFood(t, client)

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

	t.Run("POST /objects?consistency_level={level}", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)

		testsuit.CreateWeaviateTestSchemaFood(t, client)

		var (
			id1 = "abefd256-8574-442b-9293-9205193737ee"
			id2 = "565da3b6-60b3-40e5-ba21-e6bfe5dbba91"
			id3 = "07f15e48-f819-48b3-86e8-12fd8a73546d"
		)

		var (
			props1 = map[string]string{
				"name":        "Hawaii",
				"description": "Universally accepted to be the best pizza ever created.",
			}
			props2 = map[string]string{
				"name":        "ChickenSoup",
				"description": "Used by humans when their inferior genetics are attacked by microscopic organisms.",
			}
			props3 = map[string]string{
				"name":        "Pozole",
				"description": "Means “hominy” and it is basically a cross between soup and stew. It is a popular and beloved dish throughout Mexico.",
			}
		)

		resp1, err1 := client.Data().Creator().
			WithClassName("Pizza").
			WithID(id1).
			WithProperties(props1).
			WithConsistencyLevel(replication.ConsistencyLevel.ONE).
			Do(context.Background())
		assert.Nil(t, err1)
		assert.NotNil(t, resp1.Object)
		resp2, err2 := client.Data().Creator().
			WithClassName("Soup").
			WithID(id2).
			WithProperties(props2).
			WithConsistencyLevel(replication.ConsistencyLevel.ALL).
			Do(context.Background())
		assert.Nil(t, err2)
		assert.NotNil(t, resp2.Object)
		resp3, err3 := client.Data().Creator().
			WithClassName("Soup").
			WithID(id3).
			WithProperties(props3).
			WithConsistencyLevel(replication.ConsistencyLevel.QUORUM).
			Do(context.Background())
		assert.Nil(t, err3)
		assert.NotNil(t, resp3.Object)

		found1, err1 := client.Data().ObjectsGetter().
			WithClassName("Pizza").
			WithID(id1).
			Do(context.Background())
		assert.Nil(t, err1)
		assert.Equal(t, "Pizza", found1[0].Class)
		assert.Equal(t, "Hawaii", found1[0].Properties.(map[string]interface{})["name"])

		found2, err2 := client.Data().ObjectsGetter().
			WithClassName("Soup").
			WithID(id2).
			Do(context.Background())
		assert.Nil(t, err2)
		assert.Equal(t, "Soup", found2[0].Class)
		assert.Equal(t, "ChickenSoup", found2[0].Properties.(map[string]interface{})["name"])

		found3, err3 := client.Data().ObjectsGetter().
			WithClassName("Soup").
			WithID(id3).
			Do(context.Background())
		assert.Nil(t, err3)
		assert.Equal(t, "Soup", found3[0].Class)
		assert.Equal(t, "Pozole", found3[0].Properties.(map[string]interface{})["name"])

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("PUT /objects/{className}/{id}?consistency_level={level}", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)

		testsuit.CreateWeaviateTestSchemaFood(t, client)

		var (
			id    = "abefd256-8574-442b-9293-9205193737ee"
			props = map[string]string{
				"name":        "Hawaii",
				"description": "Universally accepted to be the best pizza ever created.",
			}
		)

		var (
			newProps1 = map[string]string{
				"name":        "Double Pepperoni",
				"description": "There is no such thing as too much pepperoni.",
			}
			newProps2 = map[string]string{
				"name":        "Four Cheese",
				"description": "Mozzarella, Aged Havarti, Gorgonzola, Parmigiano-Reggiano. Enough said.",
			}
			newProps3 = map[string]string{
				"name":        "Philly Cheesesteak",
				"description": "Sliced steak, peppers, onions.",
			}
		)

		createResp, createErr := client.Data().Creator().
			WithClassName("Pizza").
			WithID(id).
			WithProperties(props).
			Do(context.Background())
		assert.Nil(t, createErr)
		assert.NotNil(t, createResp.Object)

		err := client.Data().Updater().
			WithClassName("Pizza").
			WithID(id).
			WithProperties(newProps1).
			WithConsistencyLevel(replication.ConsistencyLevel.ONE).
			Do(context.Background())
		assert.Nil(t, err)

		updated1, err := client.Data().ObjectsGetter().
			WithClassName("Pizza").
			WithID(id).
			Do(context.Background())
		assert.Nil(t, err)
		assert.Equal(t, "Pizza", updated1[0].Class)
		assert.Equal(t, newProps1["name"], updated1[0].Properties.(map[string]interface{})["name"])

		err = client.Data().Updater().
			WithClassName("Pizza").
			WithID(id).
			WithProperties(newProps2).
			WithConsistencyLevel(replication.ConsistencyLevel.ALL).
			Do(context.Background())
		assert.Nil(t, err)

		updated2, err := client.Data().ObjectsGetter().
			WithClassName("Pizza").
			WithID(id).
			Do(context.Background())
		assert.Nil(t, err)
		assert.Equal(t, "Pizza", updated2[0].Class)
		assert.Equal(t, newProps2["name"], updated2[0].Properties.(map[string]interface{})["name"])

		err = client.Data().Updater().
			WithClassName("Pizza").
			WithID(id).
			WithProperties(newProps3).
			WithConsistencyLevel(replication.ConsistencyLevel.QUORUM).
			Do(context.Background())
		assert.Nil(t, err)

		updated3, err := client.Data().ObjectsGetter().
			WithClassName("Pizza").
			WithID(id).
			Do(context.Background())
		assert.Nil(t, err)
		assert.Equal(t, "Pizza", updated2[0].Class)
		assert.Equal(t, newProps3["name"], updated3[0].Properties.(map[string]interface{})["name"])

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("PATCH /objects/{className}/{id}?consistency_level={level}", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)

		testsuit.CreateWeaviateTestSchemaFood(t, client)

		var (
			id    = "abefd256-8574-442b-9293-9205193737ee"
			props = map[string]string{
				"name":        "Hawaii",
				"description": "Universally accepted to be the best pizza ever created.",
			}
		)

		var (
			newProps1 = map[string]string{
				"name":        "Double Pepperoni",
				"description": "There is no such thing as too much pepperoni.",
			}
			newProps2 = map[string]string{
				"name":        "Four Cheese",
				"description": "Mozzarella, Aged Havarti, Gorgonzola, Parmigiano-Reggiano. Enough said.",
			}
			newProps3 = map[string]string{
				"name":        "Philly Cheesesteak",
				"description": "Sliced steak, peppers, onions.",
			}
		)

		createResp, createErr := client.Data().Creator().
			WithClassName("Pizza").
			WithID(id).
			WithProperties(props).
			Do(context.Background())
		assert.Nil(t, createErr)
		assert.NotNil(t, createResp.Object)

		err := client.Data().Updater().
			WithClassName("Pizza").
			WithID(id).
			WithProperties(newProps1).
			WithConsistencyLevel(replication.ConsistencyLevel.ONE).
			WithMerge().
			Do(context.Background())
		assert.Nil(t, err)

		updated1, err := client.Data().ObjectsGetter().
			WithClassName("Pizza").
			WithID(id).
			Do(context.Background())
		assert.Nil(t, err)
		assert.Equal(t, "Pizza", updated1[0].Class)
		assert.Equal(t, newProps1["name"], updated1[0].Properties.(map[string]interface{})["name"])

		err = client.Data().Updater().
			WithClassName("Pizza").
			WithID(id).
			WithProperties(newProps2).
			WithConsistencyLevel(replication.ConsistencyLevel.ALL).
			WithMerge().
			Do(context.Background())
		assert.Nil(t, err)

		updated2, err := client.Data().ObjectsGetter().
			WithClassName("Pizza").
			WithID(id).
			Do(context.Background())
		assert.Nil(t, err)
		assert.Equal(t, "Pizza", updated2[0].Class)
		assert.Equal(t, newProps2["name"], updated2[0].Properties.(map[string]interface{})["name"])

		err = client.Data().Updater().
			WithClassName("Pizza").
			WithID(id).
			WithProperties(newProps3).
			WithConsistencyLevel(replication.ConsistencyLevel.QUORUM).
			WithMerge().
			Do(context.Background())
		assert.Nil(t, err)

		updated3, err := client.Data().ObjectsGetter().
			WithClassName("Pizza").
			WithID(id).
			Do(context.Background())
		assert.Nil(t, err)
		assert.Equal(t, "Pizza", updated2[0].Class)
		assert.Equal(t, newProps3["name"], updated3[0].Properties.(map[string]interface{})["name"])

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("DELETE /objects/{className}/{id}?consistency_level={level}", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)

		testsuit.CreateWeaviateTestSchemaFood(t, client)

		var (
			id1 = "abefd256-8574-442b-9293-9205193737ee"
			id2 = "565da3b6-60b3-40e5-ba21-e6bfe5dbba91"
			id3 = "07f15e48-f819-48b3-86e8-12fd8a73546d"
		)

		var (
			props1 = map[string]string{
				"name":        "Hawaii",
				"description": "Universally accepted to be the best pizza ever created.",
			}
			props2 = map[string]string{
				"name":        "ChickenSoup",
				"description": "Used by humans when their inferior genetics are attacked by microscopic organisms.",
			}
			props3 = map[string]string{
				"name":        "Pozole",
				"description": "Means “hominy” and it is basically a cross between soup and stew. It is a popular and beloved dish throughout Mexico.",
			}
		)

		resp1, err1 := client.Data().Creator().
			WithClassName("Pizza").
			WithID(id1).
			WithProperties(props1).
			WithConsistencyLevel(replication.ConsistencyLevel.ONE).
			Do(context.Background())
		assert.Nil(t, err1)
		assert.NotNil(t, resp1.Object)
		resp2, err2 := client.Data().Creator().
			WithClassName("Soup").
			WithID(id2).
			WithProperties(props2).
			WithConsistencyLevel(replication.ConsistencyLevel.ALL).
			Do(context.Background())
		assert.Nil(t, err2)
		assert.NotNil(t, resp2.Object)
		resp3, err3 := client.Data().Creator().
			WithClassName("Soup").
			WithID(id3).
			WithProperties(props3).
			WithConsistencyLevel(replication.ConsistencyLevel.QUORUM).
			Do(context.Background())
		assert.Nil(t, err3)
		assert.NotNil(t, resp3.Object)

		expectedErr := &fault.WeaviateClientError{
			IsUnexpectedStatusCode: true,
			StatusCode:             http.StatusNotFound,
		}

		err := client.Data().Deleter().
			WithClassName(resp1.Object.Class).
			WithID(resp1.Object.ID.String()).
			WithConsistencyLevel(replication.ConsistencyLevel.ONE).
			Do(context.Background())
		assert.Nil(t, err)

		found1, err1 := client.Data().ObjectsGetter().
			WithClassName("Pizza").
			WithID(id1).
			Do(context.Background())
		assert.EqualValues(t, expectedErr, err1)
		assert.Nil(t, found1)

		err = client.Data().Deleter().
			WithClassName(resp2.Object.Class).
			WithID(resp2.Object.ID.String()).
			WithConsistencyLevel(replication.ConsistencyLevel.ALL).
			Do(context.Background())
		assert.Nil(t, err)

		found2, err2 := client.Data().ObjectsGetter().
			WithClassName("Pizza").
			WithID(id2).
			Do(context.Background())
		assert.EqualValues(t, expectedErr, err2)
		assert.Nil(t, found2)

		err = client.Data().Deleter().
			WithClassName(resp3.Object.Class).
			WithID(resp3.Object.ID.String()).
			WithConsistencyLevel(replication.ConsistencyLevel.QUORUM).
			Do(context.Background())
		assert.Nil(t, err)

		found3, err3 := client.Data().ObjectsGetter().
			WithClassName("Pizza").
			WithID(id2).
			Do(context.Background())
		assert.EqualValues(t, expectedErr, err3)
		assert.Nil(t, found3)

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			fmt.Print(err.Error())
			t.Fail()
		}
	})
}

func TestData_MultiTenancy(t *testing.T) {
	cleanup := func() {
		client := testsuit.CreateTestClient(false)
		err := client.Schema().AllDeleter().Do(context.Background())
		require.Nil(t, err)
	}

	t.Run("setup weaviate", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to setup weaviate: %s", err)
		}
	})

	t.Run("creates objects of MT class", func(t *testing.T) {
		defer cleanup()

		className := "Pizza"
		tenants := testsuit.Tenants{
			{Name: "tenantNo1"},
			{Name: "tenantNo2"},
		}

		client := testsuit.CreateTestClient(false)
		testsuit.CreateSchemaPizzaForTenants(t, client)
		testsuit.CreateTenantsPizza(t, client, tenants...)

		for _, tenant := range tenants {
			wrap, err := client.Data().Creator().
				WithClassName(className).
				WithID(testsuit.PIZZA_QUATTRO_FORMAGGI_ID).
				WithProperties(map[string]interface{}{
					"name":        "Quattro Formaggi",
					"description": "Pizza quattro formaggi Italian: [ˈkwattro forˈmaddʒi] (four cheese pizza) is a variety of pizza in Italian cuisine that is topped with a combination of four kinds of cheese, usually melted together, with (rossa, red) or without (bianca, white) tomato sauce. It is popular worldwide, including in Italy,[1] and is one of the iconic items from pizzerias's menus.",
					"price":       float32(1.1),
					"best_before": "2022-05-03T12:04:40+02:00",
				}).
				WithTenant(tenant.Name).
				Do(context.Background())

			require.Nil(t, err)
			require.NotNil(t, wrap)
			require.NotNil(t, wrap.Object)
			assert.Equal(t, strfmt.UUID(testsuit.PIZZA_QUATTRO_FORMAGGI_ID), wrap.Object.ID)
			assert.Equal(t, "Quattro Formaggi", wrap.Object.Properties.(map[string]interface{})["name"])
			assert.Equal(t, tenant.Name, wrap.Object.Tenant)

			wrap, err = client.Data().Creator().
				WithClassName(className).
				WithID(testsuit.PIZZA_FRUTTI_DI_MARE_ID).
				WithProperties(map[string]interface{}{
					"name":        "Frutti di Mare",
					"description": "Frutti di Mare is an Italian type of pizza that may be served with scampi, mussels or squid. It typically lacks cheese, with the seafood being served atop a tomato sauce.",
					"price":       float32(1.2),
					"best_before": "2022-05-05T07:16:30+02:00",
				}).
				WithTenant(tenant.Name).
				Do(context.Background())

			require.Nil(t, err)
			require.NotNil(t, wrap)
			require.NotNil(t, wrap.Object)
			assert.Equal(t, strfmt.UUID(testsuit.PIZZA_FRUTTI_DI_MARE_ID), wrap.Object.ID)
			assert.Equal(t, "Frutti di Mare", wrap.Object.Properties.(map[string]interface{})["name"])
			assert.Equal(t, tenant.Name, wrap.Object.Tenant)
		}

		t.Run("verify created", func(t *testing.T) {
			for _, tenant := range tenants {
				for _, id := range []string{
					testsuit.PIZZA_QUATTRO_FORMAGGI_ID,
					testsuit.PIZZA_FRUTTI_DI_MARE_ID,
				} {
					exists, err := client.Data().Checker().
						WithID(id).
						WithClassName(className).
						WithTenant(tenant.Name).
						Do(context.Background())

					require.Nil(t, err)
					require.True(t, exists)
				}
			}
		})
	})

	t.Run("fails creating objects of MT class without tenant", func(t *testing.T) {
		defer cleanup()

		className := "Pizza"
		tenants := testsuit.Tenants{
			{Name: "tenantNo1"},
			{Name: "tenantNo2"},
		}

		client := testsuit.CreateTestClient(false)
		testsuit.CreateSchemaPizzaForTenants(t, client)
		testsuit.CreateTenantsPizza(t, client, tenants...)

		wrap, err := client.Data().Creator().
			WithClassName(className).
			WithID(testsuit.PIZZA_QUATTRO_FORMAGGI_ID).
			WithProperties(map[string]interface{}{
				"name":        "Quattro Formaggi",
				"description": "Pizza quattro formaggi Italian: [ˈkwattro forˈmaddʒi] (four cheese pizza) is a variety of pizza in Italian cuisine that is topped with a combination of four kinds of cheese, usually melted together, with (rossa, red) or without (bianca, white) tomato sauce. It is popular worldwide, including in Italy,[1] and is one of the iconic items from pizzerias's menus.",
				"price":       float32(1.1),
				"best_before": "2022-05-03T12:04:40+02:00",
			}).
			Do(context.Background())

		require.NotNil(t, err)
		clientErr := err.(*fault.WeaviateClientError)
		assert.Equal(t, 422, clientErr.StatusCode)
		assert.Contains(t, clientErr.Msg, "has multi-tenancy enabled, but request was without tenant")
		require.Nil(t, wrap)

		wrap, err = client.Data().Creator().
			WithClassName(className).
			WithID(testsuit.PIZZA_FRUTTI_DI_MARE_ID).
			WithProperties(map[string]interface{}{
				"name":        "Frutti di Mare",
				"description": "Frutti di Mare is an Italian type of pizza that may be served with scampi, mussels or squid. It typically lacks cheese, with the seafood being served atop a tomato sauce.",
				"price":       float32(1.2),
				"best_before": "2022-05-05T07:16:30+02:00",
			}).
			Do(context.Background())

		require.NotNil(t, err)
		clientErr = err.(*fault.WeaviateClientError)
		assert.Equal(t, 422, clientErr.StatusCode)
		assert.Contains(t, clientErr.Msg, "has multi-tenancy enabled, but request was without tenant")
		require.Nil(t, wrap)

		t.Run("verify not created", func(t *testing.T) {
			for _, tenant := range tenants {
				for _, id := range []string{
					testsuit.PIZZA_QUATTRO_FORMAGGI_ID,
					testsuit.PIZZA_FRUTTI_DI_MARE_ID,
				} {
					exists, err := client.Data().Checker().
						WithID(id).
						WithClassName(className).
						WithTenant(tenant.Name).
						Do(context.Background())

					require.Nil(t, err)
					require.False(t, exists)
				}
			}
		})
	})

	t.Run("gets objects of MT class", func(t *testing.T) {
		defer cleanup()

		tenants := testsuit.Tenants{
			{Name: "tenantNo1"},
			{Name: "tenantNo2"},
		}

		client := testsuit.CreateTestClient(false)
		testsuit.CreateSchemaFoodForTenants(t, client)
		testsuit.CreateTenantsFood(t, client, tenants...)
		testsuit.CreateDataFoodForTenants(t, client, tenants.Names()...)

		extractIds := func(objs []*models.Object) []string {
			ids := make([]string, len(objs))
			for i, obj := range objs {
				ids[i] = obj.ID.String()
			}
			return ids
		}

		for _, tenant := range tenants {
			for className, ids := range testsuit.IdsByClass {
				for _, id := range ids {
					t.Run("single object by class+id", func(t *testing.T) {
						objects, err := client.Data().ObjectsGetter().
							WithID(id).
							WithClassName(className).
							WithTenant(tenant.Name).
							Do(context.Background())

						require.Nil(t, err)
						require.NotNil(t, objects)
						require.Len(t, objects, 1)
						assert.Equal(t, strfmt.UUID(id), objects[0].ID)
						assert.Equal(t, tenant.Name, objects[0].Tenant)
					})
				}

				t.Run("list objects by class", func(t *testing.T) {
					objects, err := client.Data().ObjectsGetter().
						WithClassName(className).
						WithTenant(tenant.Name).
						Do(context.Background())

					require.Nil(t, err)
					require.NotNil(t, objects)
					require.Len(t, objects, len(ids))
					assert.ElementsMatch(t, ids, extractIds(objects))
				})
			}

			t.Run("list all objects", func(t *testing.T) {
				objects, err := client.Data().ObjectsGetter().
					WithTenant(tenant.Name).
					Do(context.Background())

				require.Nil(t, err)
				require.NotNil(t, objects)
				require.Len(t, objects, len(testsuit.AllIds))
				assert.ElementsMatch(t, testsuit.AllIds, extractIds(objects))
			})
		}
	})

	t.Run("fails getting objects of MT class without tenant", func(t *testing.T) {
		defer cleanup()

		tenants := testsuit.Tenants{
			{Name: "tenantNo1"},
			{Name: "tenantNo2"},
		}

		client := testsuit.CreateTestClient(false)
		testsuit.CreateSchemaFoodForTenants(t, client)
		testsuit.CreateTenantsFood(t, client, tenants...)
		testsuit.CreateDataFoodForTenants(t, client, tenants.Names()...)

		for className, ids := range testsuit.IdsByClass {
			for _, id := range ids {
				t.Run("single object by class+id", func(t *testing.T) {
					objects, err := client.Data().ObjectsGetter().
						WithID(id).
						WithClassName(className).
						Do(context.Background())

					require.NotNil(t, err)
					clientErr := err.(*fault.WeaviateClientError)
					assert.Equal(t, 422, clientErr.StatusCode)
					assert.Contains(t, clientErr.Msg, "has multi-tenancy enabled, but request was without tenant")
					assert.Nil(t, objects)
				})
			}

			t.Run("list objects by class", func(t *testing.T) {
				objects, err := client.Data().ObjectsGetter().
					WithClassName(className).
					Do(context.Background())

				require.NotNil(t, err)
				clientErr := err.(*fault.WeaviateClientError)
				assert.Equal(t, 422, clientErr.StatusCode)
				assert.Contains(t, clientErr.Msg, "has multi-tenancy enabled, but request was without tenant")
				assert.Nil(t, objects)
			})
		}

		t.Run("list all objects", func(t *testing.T) {
			objects, err := client.Data().ObjectsGetter().
				Do(context.Background())

			require.Nil(t, err)
			require.NotNil(t, objects)
			assert.Len(t, objects, 0)
		})
	})

	t.Run("checks objects of MT class", func(t *testing.T) {
		defer cleanup()

		tenants := testsuit.Tenants{
			{Name: "tenantNo1"},
			{Name: "tenantNo2"},
		}

		client := testsuit.CreateTestClient(false)
		testsuit.CreateSchemaFoodForTenants(t, client)
		testsuit.CreateTenantsFood(t, client, tenants...)
		testsuit.CreateDataFoodForTenants(t, client, tenants.Names()...)

		for _, tenant := range tenants {
			for className, ids := range testsuit.IdsByClass {
				for _, id := range ids {
					exists, err := client.Data().Checker().
						WithID(id).
						WithClassName(className).
						WithTenant(tenant.Name).
						Do(context.Background())

					require.Nil(t, err)
					require.True(t, exists)
				}
			}
		}
	})

	t.Run("fails checking objects of MT class without tenant", func(t *testing.T) {
		defer cleanup()

		tenants := testsuit.Tenants{
			{Name: "tenantNo1"},
			{Name: "tenantNo2"},
		}

		client := testsuit.CreateTestClient(false)
		testsuit.CreateSchemaFoodForTenants(t, client)
		testsuit.CreateTenantsFood(t, client, tenants...)
		testsuit.CreateDataFoodForTenants(t, client, tenants.Names()...)

		for className, ids := range testsuit.IdsByClass {
			for _, id := range ids {
				exists, err := client.Data().Checker().
					WithID(id).
					WithClassName(className).
					Do(context.Background())

				require.NotNil(t, err)
				clientErr := err.(*fault.WeaviateClientError)
				assert.Equal(t, 422, clientErr.StatusCode)
				assert.Empty(t, clientErr.Msg) // no body in HEAD
				assert.False(t, exists)
			}
		}
	})

	t.Run("deletes objects from MT class", func(t *testing.T) {
		defer cleanup()

		tenants := testsuit.Tenants{
			{Name: "tenantNo1"},
			{Name: "tenantNo2"},
		}

		client := testsuit.CreateTestClient(false)
		testsuit.CreateSchemaFoodForTenants(t, client)
		testsuit.CreateTenantsFood(t, client, tenants...)
		testsuit.CreateDataFoodForTenants(t, client, tenants.Names()...)

		for _, tenant := range tenants {
			for className, ids := range testsuit.IdsByClass {
				expectedLeft := len(ids)

				for _, id := range ids {
					err := client.Data().Deleter().
						WithID(id).
						WithClassName(className).
						WithTenant(tenant.Name).
						Do(context.Background())

					require.Nil(t, err)
					expectedLeft--

					t.Run("verify deleted", func(t *testing.T) {
						exists, err := client.Data().Checker().
							WithID(id).
							WithClassName(className).
							WithTenant(tenant.Name).
							Do(context.Background())

						require.Nil(t, err)
						require.False(t, exists)
					})

					t.Run("verify left", func(t *testing.T) {
						objects, err := client.Data().ObjectsGetter().
							WithClassName(className).
							WithTenant(tenant.Name).
							Do(context.Background())

						require.Nil(t, err)
						require.NotNil(t, objects)
						assert.Len(t, objects, expectedLeft)
					})
				}
			}
		}
	})

	t.Run("fails deleting objects from MT class without tenant", func(t *testing.T) {
		defer cleanup()

		tenants := testsuit.Tenants{
			{Name: "tenantNo1"},
			{Name: "tenantNo2"},
		}

		client := testsuit.CreateTestClient(false)
		testsuit.CreateSchemaFoodForTenants(t, client)
		testsuit.CreateTenantsFood(t, client, tenants...)
		testsuit.CreateDataFoodForTenants(t, client, tenants.Names()...)

		for className, ids := range testsuit.IdsByClass {
			for _, id := range ids {
				err := client.Data().Deleter().
					WithID(id).
					WithClassName(className).
					Do(context.Background())

				require.NotNil(t, err)
				clientErr := err.(*fault.WeaviateClientError)
				assert.Equal(t, 422, clientErr.StatusCode)
				assert.Contains(t, clientErr.Msg, "has multi-tenancy enabled, but request was without tenant")

				t.Run("verify not deleted", func(t *testing.T) {
					for _, tenant := range tenants {
						exists, err := client.Data().Checker().
							WithID(id).
							WithClassName(className).
							WithTenant(tenant.Name).
							Do(context.Background())

						require.Nil(t, err)
						require.True(t, exists)
					}
				})
			}
		}

		t.Run("verify not deleted", func(t *testing.T) {
			for _, tenant := range tenants {
				objects, err := client.Data().ObjectsGetter().
					WithTenant(tenant.Name).
					Do(context.Background())

				require.Nil(t, err)
				require.NotNil(t, objects)
				assert.Len(t, objects, len(testsuit.AllIds))
			}
		})
	})

	t.Run("updates objects of MT class", func(t *testing.T) {
		defer cleanup()

		className := "Soup"
		tenants := testsuit.Tenants{
			{Name: "tenantNo1"},
			{Name: "tenantNo2"},
		}

		client := testsuit.CreateTestClient(false)
		testsuit.CreateSchemaSoupForTenants(t, client)
		testsuit.CreateTenantsSoup(t, client, tenants...)
		testsuit.CreateDataSoupForTenants(t, client, tenants.Names()...)

		for _, tenant := range tenants {
			err := client.Data().Updater().
				WithClassName(className).
				WithID(testsuit.SOUP_CHICKENSOUP_ID).
				WithProperties(map[string]interface{}{
					"name":        "ChickenSoup",
					"description": fmt.Sprintf("updated ChickenSoup description [%s]", tenant),
					"price":       float32(2.1),
				}).
				WithTenant(tenant.Name).
				Do(context.Background())

			require.Nil(t, err)

			err = client.Data().Updater().
				WithClassName(className).
				WithID(testsuit.SOUP_BEAUTIFUL_ID).
				WithProperties(map[string]interface{}{
					"name":        "Beautiful",
					"description": fmt.Sprintf("updated Beautiful description [%s]", tenant),
					"price":       float32(2.2),
				}).
				WithTenant(tenant.Name).
				Do(context.Background())

			require.Nil(t, err)
		}

		t.Run("verify updated", func(t *testing.T) {
			for _, tenant := range tenants {
				objects, err := client.Data().ObjectsGetter().
					WithID(testsuit.SOUP_CHICKENSOUP_ID).
					WithClassName(className).
					WithTenant(tenant.Name).
					Do(context.Background())

				require.Nil(t, err)
				require.NotNil(t, objects)
				require.Len(t, objects, 1)
				assert.Equal(t, tenant.Name, objects[0].Tenant)
				assert.Equal(t, fmt.Sprintf("updated ChickenSoup description [%s]", tenant),
					objects[0].Properties.(map[string]interface{})["description"])

				objects, err = client.Data().ObjectsGetter().
					WithID(testsuit.SOUP_BEAUTIFUL_ID).
					WithClassName(className).
					WithTenant(tenant.Name).
					Do(context.Background())

				require.Nil(t, err)
				require.NotNil(t, objects)
				require.Len(t, objects, 1)
				assert.Equal(t, tenant.Name, objects[0].Tenant)
				assert.Equal(t, fmt.Sprintf("updated Beautiful description [%s]", tenant),
					objects[0].Properties.(map[string]interface{})["description"])
			}
		})
	})

	t.Run("fails updating objects of MT class without tenant", func(t *testing.T) {
		defer cleanup()

		className := "Soup"
		tenants := testsuit.Tenants{
			{Name: "tenantNo1"},
			{Name: "tenantNo2"},
		}

		client := testsuit.CreateTestClient(false)
		testsuit.CreateSchemaSoupForTenants(t, client)
		testsuit.CreateTenantsSoup(t, client, tenants...)
		testsuit.CreateDataSoupForTenants(t, client, tenants.Names()...)

		err := client.Data().Updater().
			WithClassName(className).
			WithID(testsuit.SOUP_CHICKENSOUP_ID).
			WithProperties(map[string]interface{}{
				"name":        "ChickenSoup",
				"description": "updated ChickenSoup description",
				"price":       float32(2.1),
			}).
			Do(context.Background())

		require.NotNil(t, err)
		clientErr := err.(*fault.WeaviateClientError)
		assert.Equal(t, 422, clientErr.StatusCode)
		assert.Contains(t, clientErr.Msg, "has multi-tenancy enabled, but request was without tenant")

		err = client.Data().Updater().
			WithClassName(className).
			WithID(testsuit.SOUP_BEAUTIFUL_ID).
			WithProperties(map[string]interface{}{
				"name":        "Beautiful",
				"description": "updated Beautiful description",
				"price":       float32(2.2),
			}).
			Do(context.Background())

		require.NotNil(t, err)
		clientErr = err.(*fault.WeaviateClientError)
		assert.Equal(t, 422, clientErr.StatusCode)
		assert.Contains(t, clientErr.Msg, "has multi-tenancy enabled, but request was without tenant")

		t.Run("verify not updated", func(t *testing.T) {
			for _, tenant := range tenants {
				objects, err := client.Data().ObjectsGetter().
					WithID(testsuit.SOUP_CHICKENSOUP_ID).
					WithClassName(className).
					WithTenant(tenant.Name).
					Do(context.Background())

				require.Nil(t, err)
				require.NotNil(t, objects)
				require.Len(t, objects, 1)
				assert.Equal(t, tenant.Name, objects[0].Tenant)
				assert.Equal(t, "Used by humans when their inferior genetics are attacked by microscopic organisms.",
					objects[0].Properties.(map[string]interface{})["description"])

				objects, err = client.Data().ObjectsGetter().
					WithID(testsuit.SOUP_BEAUTIFUL_ID).
					WithClassName(className).
					WithTenant(tenant.Name).
					Do(context.Background())

				require.Nil(t, err)
				require.NotNil(t, objects)
				require.Len(t, objects, 1)
				assert.Equal(t, tenant.Name, objects[0].Tenant)
				assert.Equal(t, "Putting the game of letter soups to a whole new level.",
					objects[0].Properties.(map[string]interface{})["description"])
			}
		})
	})

	t.Run("merges objects of MT class", func(t *testing.T) {
		defer cleanup()

		className := "Soup"
		tenants := testsuit.Tenants{
			{Name: "tenantNo1"},
			{Name: "tenantNo2"},
		}

		client := testsuit.CreateTestClient(false)
		testsuit.CreateSchemaSoupForTenants(t, client)
		testsuit.CreateTenantsSoup(t, client, tenants...)
		testsuit.CreateDataSoupForTenants(t, client, tenants.Names()...)

		for _, tenant := range tenants {
			err := client.Data().Updater().
				WithClassName(className).
				WithID(testsuit.SOUP_CHICKENSOUP_ID).
				WithProperties(map[string]interface{}{
					"description": fmt.Sprintf("merged ChickenSoup description [%s]", tenant),
				}).
				WithTenant(tenant.Name).
				WithMerge().
				Do(context.Background())

			require.Nil(t, err)

			err = client.Data().Updater().
				WithClassName(className).
				WithID(testsuit.SOUP_BEAUTIFUL_ID).
				WithProperties(map[string]interface{}{
					"description": fmt.Sprintf("merged Beautiful description [%s]", tenant),
				}).
				WithTenant(tenant.Name).
				WithMerge().
				Do(context.Background())

			require.Nil(t, err)
		}

		t.Run("verify merged", func(t *testing.T) {
			for _, tenant := range tenants {
				objects, err := client.Data().ObjectsGetter().
					WithID(testsuit.SOUP_CHICKENSOUP_ID).
					WithClassName(className).
					WithTenant(tenant.Name).
					Do(context.Background())

				require.Nil(t, err)
				require.NotNil(t, objects)
				require.Len(t, objects, 1)
				assert.Equal(t, tenant.Name, objects[0].Tenant)
				assert.Equal(t, fmt.Sprintf("merged ChickenSoup description [%s]", tenant),
					objects[0].Properties.(map[string]interface{})["description"])

				objects, err = client.Data().ObjectsGetter().
					WithID(testsuit.SOUP_BEAUTIFUL_ID).
					WithClassName(className).
					WithTenant(tenant.Name).
					Do(context.Background())

				require.Nil(t, err)
				require.NotNil(t, objects)
				require.Len(t, objects, 1)
				assert.Equal(t, tenant.Name, objects[0].Tenant)
				assert.Equal(t, fmt.Sprintf("merged Beautiful description [%s]", tenant),
					objects[0].Properties.(map[string]interface{})["description"])
			}
		})
	})

	t.Run("fails merging objects of MT class without tenant", func(t *testing.T) {
		defer cleanup()

		className := "Soup"
		tenants := testsuit.Tenants{
			{Name: "tenantNo1"},
			{Name: "tenantNo2"},
		}

		client := testsuit.CreateTestClient(false)
		testsuit.CreateSchemaSoupForTenants(t, client)
		testsuit.CreateTenantsSoup(t, client, tenants...)
		testsuit.CreateDataSoupForTenants(t, client, tenants.Names()...)

		err := client.Data().Updater().
			WithClassName(className).
			WithID(testsuit.SOUP_CHICKENSOUP_ID).
			WithProperties(map[string]interface{}{
				"description": "merged ChickenSoup description",
			}).
			WithMerge().
			Do(context.Background())

		require.NotNil(t, err)
		clientErr := err.(*fault.WeaviateClientError)
		assert.Equal(t, 422, clientErr.StatusCode)
		assert.Contains(t, clientErr.Msg, "has multi-tenancy enabled, but request was without tenant")

		err = client.Data().Updater().
			WithClassName(className).
			WithID(testsuit.SOUP_BEAUTIFUL_ID).
			WithProperties(map[string]interface{}{
				"description": "merged Beautiful description",
			}).
			WithMerge().
			Do(context.Background())

		require.NotNil(t, err)
		clientErr = err.(*fault.WeaviateClientError)
		assert.Equal(t, 422, clientErr.StatusCode)
		assert.Contains(t, clientErr.Msg, "has multi-tenancy enabled, but request was without tenant")

		t.Run("verify not merged", func(t *testing.T) {
			for _, tenant := range tenants {
				objects, err := client.Data().ObjectsGetter().
					WithID(testsuit.SOUP_CHICKENSOUP_ID).
					WithClassName(className).
					WithTenant(tenant.Name).
					Do(context.Background())

				require.Nil(t, err)
				require.NotNil(t, objects)
				require.Len(t, objects, 1)
				assert.Equal(t, tenant.Name, objects[0].Tenant)
				assert.Equal(t, "Used by humans when their inferior genetics are attacked by microscopic organisms.",
					objects[0].Properties.(map[string]interface{})["description"])

				objects, err = client.Data().ObjectsGetter().
					WithID(testsuit.SOUP_BEAUTIFUL_ID).
					WithClassName(className).
					WithTenant(tenant.Name).
					Do(context.Background())

				require.Nil(t, err)
				require.NotNil(t, objects)
				require.Len(t, objects, 1)
				assert.Equal(t, tenant.Name, objects[0].Tenant)
				assert.Equal(t, "Putting the game of letter soups to a whole new level.",
					objects[0].Properties.(map[string]interface{})["description"])
			}
		})

		t.Run("auto tenant creation and tenant exists", func(t *testing.T) {
			ctx := context.TODO()
			className := "MultiAutoCreateTenantClass"
			autoTenant := "AutoTenant"
			id := strfmt.UUID("10000000-0000-0000-0000-000000000000")
			class := &models.Class{
				Class: className,
				Properties: []*models.Property{
					{
						Name: "name", DataType: []string{schema.DataTypeText.String()},
					},
				},
				MultiTenancyConfig: &models.MultiTenancyConfig{
					Enabled:            true,
					AutoTenantCreation: true,
				},
			}
			err := client.Schema().ClassCreator().WithClass(class).Do(ctx)
			require.NoError(t, err)

			exists, err := client.Schema().TenantsExists().WithClassName(className).WithTenant(autoTenant).Do(ctx)
			require.NoError(t, err)
			assert.False(t, exists)

			obj := &models.Object{
				ID:    id,
				Class: className,
				Properties: map[string]interface{}{
					"name": "some name",
				},
				Tenant: autoTenant,
			}

			resp, err := client.Batch().ObjectsBatcher().WithObjects(obj).Do(ctx)
			require.NoError(t, err)
			require.NotEmpty(t, resp)
			assert.Len(t, resp, 1)

			exists, err = client.Schema().TenantsExists().WithClassName(className).WithTenant(autoTenant).Do(ctx)
			require.NoError(t, err)
			assert.True(t, exists)
		})
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to tear down weaviate: %s", err)
		}
	})
}
