package data

import (
	"context"
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/test/testsuit"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/clienterrors"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/paragons"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/testenv"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestData_integration(t *testing.T) {

	t.Run("up", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			fmt.Printf(err.Error())
			t.Fail()
		}
	})

	t.Run("POST /{semanticType}", func(t *testing.T) {

		client := testsuit.CreateTestClient()

		testsuit.CreateWeaviateTestSchemaFood(t, client)

		propertySchemaT := map[string]string{
			"name": "Hawaii",
			"description": "Universally accepted to be the best pizza ever created.",
		}
		propertySchemaA := map[string]string{
			"name": "ChickenSoup",
			"description": "Used by humans when their inferior genetics are attacked by microscopic organisms.",
		}

		errCreateT := client.Data.Creator().WithClassName("Pizza").WithID("abefd256-8574-442b-9293-9205193737ee").WithSchema(propertySchemaT).Do(context.Background())
		assert.Nil(t, errCreateT)
		errCreateA := client.Data.Creator().WithClassName("Soup").WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").WithSchema(propertySchemaA).WithKind(paragons.SemanticKindActions).Do(context.Background())
		assert.Nil(t, errCreateA)

		time.Sleep(2.0 * time.Second) // Give weaviate time to update its index
		objectT, objErrT := client.Data.ThingsGetter().WithID("abefd256-8574-442b-9293-9205193737ee").Do(context.Background())
		assert.Nil(t, objErrT)
		objectA, objErrA := client.Data.ActionsGetter().WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").Do(context.Background())
		assert.Nil(t, objErrA)

		assert.Equal(t, "Pizza", objectT[0].Class)
		valuesT := objectT[0].Schema.(map[string]interface{})
		assert.Equal(t, "Hawaii", valuesT["name"])
		assert.Equal(t, "Soup", objectA[0].Class)
		valuesA := objectA[0].Schema.(map[string]interface{})
		assert.Equal(t, "ChickenSoup", valuesA["name"])

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("GET /actions /things", func(t *testing.T) {
		client := testsuit.CreateTestClient()
		testsuit.CreateWeaviateTestSchemaFood(t, client)

		errCreate := client.Data.Creator().WithClassName("Pizza").WithSchema(map[string]string{
			"name": "Margherita",
			"description": "plain",
		}).Do(context.Background())
		assert.Nil(t, errCreate)
		errCreate = client.Data.Creator().WithClassName("Pizza").WithSchema(map[string]string{
			"name": "Pepperoni",
			"description": "meat",
		}).Do(context.Background())
		assert.Nil(t, errCreate)
		errCreate = client.Data.Creator().WithClassName("Soup").WithKind(paragons.SemanticKindActions).WithSchema(map[string]string{
			"name": "Chicken",
			"description": "meat",
		}).Do(context.Background())
		assert.Nil(t, errCreate)
		errCreate = client.Data.Creator().WithClassName("Soup").WithKind(paragons.SemanticKindActions).WithSchema(map[string]string{
			"name": "Tofu",
			"description": "vegetarian",
		}).Do(context.Background())
		assert.Nil(t, errCreate)

		time.Sleep(2.0 * time.Second)
		objectT, objErrT := client.Data.ThingsGetter().Do(context.Background())
		assert.Nil(t, objErrT)
		objectA, objErrA := client.Data.ActionsGetter().Do(context.Background())
		assert.Nil(t, objErrA)

		assert.Equal(t, 2, len(objectT))
		assert.Equal(t, 2, len(objectA))
		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("GET underscore properties", func(t *testing.T) {
		client := testsuit.CreateTestClient()

		testsuit.CreateWeaviateTestSchemaFood(t, client)

		propertySchemaT := map[string]string{
			"name": "Hawaii",
			"description": "Universally accepted to be the best pizza ever created.",
		}
		propertySchemaA := map[string]string{
			"name": "ChickenSoup",
			"description": "Used by humans when their inferior genetics are attacked by microscopic organisms.",
		}
		errCreateT := client.Data.Creator().WithClassName("Pizza").WithID("abefd256-8574-442b-9293-9205193737ee").WithSchema(propertySchemaT).Do(context.Background())
		assert.Nil(t, errCreateT)
		errCreateA := client.Data.Creator().WithClassName("Soup").WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").WithSchema(propertySchemaA).WithKind(paragons.SemanticKindActions).Do(context.Background())
		assert.Nil(t, errCreateA)

		time.Sleep(2.0 * time.Second) // Give weaviate time to update its index
		// THINGS
		objectT, objErrT := client.Data.ThingsGetter().WithID("abefd256-8574-442b-9293-9205193737ee").Do(context.Background())
		assert.Nil(t, objErrT)
		assert.Nil(t, objectT[0].Classification)
		assert.Nil(t, objectT[0].NearestNeighbors)
		assert.Nil(t, objectT[0].FeatureProjection)
		assert.Nil(t, objectT[0].Vector)
		assert.Nil(t, objectT[0].Interpretation)

		objectT, objErrT = client.Data.ThingsGetter().WithID("abefd256-8574-442b-9293-9205193737ee").WithUnderscoreInterpretation().Do(context.Background())
		assert.Nil(t, objErrT)
		assert.Nil(t, objectT[0].Classification)
		assert.Nil(t, objectT[0].NearestNeighbors)
		assert.Nil(t, objectT[0].FeatureProjection)
		assert.Nil(t, objectT[0].Vector)
		assert.NotNil(t, objectT[0].Interpretation)

		objectT, objErrT = client.Data.ThingsGetter().WithID("abefd256-8574-442b-9293-9205193737ee").WithUnderscoreInterpretation().WithUnderscoreClassification().WithUnderscoreNearestNeighbors().WithUnderscoreVector().Do(context.Background())
		assert.Nil(t, objErrT)
		assert.Nil(t, objectT[0].Classification) // Is nil because no classifications was executed
		assert.NotNil(t, objectT[0].NearestNeighbors)
		assert.Nil(t, objectT[0].FeatureProjection) // Is nil because feature projection is not possible on non list request
		assert.NotNil(t, objectT[0].Vector)
		assert.NotNil(t, objectT[0].Interpretation)

		// ACTIONS
		objectA, objErrA := client.Data.ActionsGetter().WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").Do(context.Background())
		assert.Nil(t, objErrA)
		assert.Nil(t, objectA[0].Classification)
		assert.Nil(t, objectA[0].NearestNeighbors)
		assert.Nil(t, objectA[0].FeatureProjection)
		assert.Nil(t, objectA[0].Vector)
		assert.Nil(t, objectA[0].Interpretation)

		objectA, objErrA = client.Data.ActionsGetter().WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").WithUnderscoreInterpretation().Do(context.Background())
		assert.Nil(t, objErrA)
		assert.Nil(t, objectA[0].Classification)
		assert.Nil(t, objectA[0].NearestNeighbors)
		assert.Nil(t, objectA[0].FeatureProjection)
		assert.Nil(t, objectA[0].Vector)
		assert.NotNil(t, objectA[0].Interpretation)

		objectA, objErrA = client.Data.ActionsGetter().WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").WithUnderscoreInterpretation().WithUnderscoreClassification().WithUnderscoreNearestNeighbors().WithUnderscoreFeatureProjection().WithUnderscoreVector().Do(context.Background())
		assert.Nil(t, objErrT)
		assert.Nil(t, objectT[0].Classification) // Is nil because no classifications was executed
		assert.NotNil(t, objectT[0].NearestNeighbors)
		assert.Nil(t, objectT[0].FeatureProjection) // Is nil because feature projection is not possible on non list request
		assert.NotNil(t, objectT[0].Vector)
		assert.NotNil(t, objectT[0].Interpretation)

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("DELETE /{type}", func(t *testing.T) {
		client := testsuit.CreateTestClient()

		testsuit.CreateWeaviateTestSchemaFood(t, client)

		propertySchemaT := map[string]string{
			"name": "Hawaii",
			"description": "Universally accepted to be the best pizza ever created.",
		}
		propertySchemaA := map[string]string{
			"name": "ChickenSoup",
			"description": "Used by humans when their inferior genetics are attacked by microscopic organisms.",
		}
		errCreateT := client.Data.Creator().WithClassName("Pizza").WithID("abefd256-8574-442b-9293-9205193737ee").WithSchema(propertySchemaT).Do(context.Background())
		assert.Nil(t, errCreateT)
		errCreateA := client.Data.Creator().WithClassName("Soup").WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").WithSchema(propertySchemaA).WithKind(paragons.SemanticKindActions).Do(context.Background())
		assert.Nil(t, errCreateA)

		time.Sleep(2.0 * time.Second) // Give weaviate time to update its index
		// THINGS
		deleteErrT := client.Data.Deleter().WithID("abefd256-8574-442b-9293-9205193737ee").Do(context.Background())
		assert.Nil(t, deleteErrT)
		time.Sleep(2.0 * time.Second)
		_, getErrT := client.Data.ThingsGetter().WithID("abefd256-8574-442b-9293-9205193737ee").Do(context.Background())
		statusCodeErrorT := getErrT.(*clienterrors.UnexpectedStatusCodeError)
		assert.Equal(t, 404, statusCodeErrorT.StatusCode)

		deleteErrA := client.Data.Deleter().WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").WithKind(paragons.SemanticKindActions).Do(context.Background())
		assert.Nil(t, deleteErrA)
		_, getErrA := client.Data.ThingsGetter().WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").Do(context.Background())
		statusCodeErrorA := getErrA.(*clienterrors.UnexpectedStatusCodeError)
		assert.Equal(t, 404, statusCodeErrorA.StatusCode)

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("PUT /{type}/{id}", func(t *testing.T) {
		// PUT replaces the object fully
		client := testsuit.CreateTestClient()

		testsuit.CreateWeaviateTestSchemaFood(t, client)

		propertySchemaT := map[string]string{
			"name": "Random",
			"description": "Missing description",
		}
		propertySchemaA := map[string]string{
			"name": "water",
			"description": "missing description",
		}
		errCreateT := client.Data.Creator().WithClassName("Pizza").WithID("abefd256-8574-442b-9293-9205193737ee").WithSchema(propertySchemaT).Do(context.Background())
		assert.Nil(t, errCreateT)
		errCreateA := client.Data.Creator().WithClassName("Soup").WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").WithSchema(propertySchemaA).WithKind(paragons.SemanticKindActions).Do(context.Background())
		assert.Nil(t, errCreateA)

		time.Sleep(2.0 * time.Second) // Give weaviate time to update its index
		propertySchemaT = map[string]string{
			"name": "Hawaii",
			"description": "Universally accepted to be the best pizza ever created.",
		}
		updateErrT := client.Data.Updater().WithID("abefd256-8574-442b-9293-9205193737ee").WithClassName("Pizza").WithSchema(propertySchemaT).Do(context.Background())
		assert.Nil(t, updateErrT)

		propertySchemaA = map[string]string{
			"name": "ChickenSoup",
			"description": "Used by humans when their inferior genetics are attacked by microscopic organisms.",
		}
		updateErrA := client.Data.Updater().WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").WithClassName("Soup").WithSchema(propertySchemaA).WithKind(paragons.SemanticKindActions).Do(context.Background())
		assert.Nil(t, updateErrA)
		time.Sleep(2.0 * time.Second)

		things, getErrT := client.Data.ThingsGetter().WithID("abefd256-8574-442b-9293-9205193737ee").Do(context.Background())
		assert.Nil(t, getErrT)
		valuesT := things[0].Schema.(map[string]interface{})
		assert.Equal(t, propertySchemaT["description"], valuesT["description"])
		assert.Equal(t, propertySchemaT["name"], valuesT["name"])

		actions, getErrT := client.Data.ActionsGetter().WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").Do(context.Background())
		assert.Nil(t, getErrT)
		valuesA := actions[0].Schema.(map[string]interface{})
		assert.Equal(t, propertySchemaA["description"], valuesA["description"])
		assert.Equal(t, propertySchemaA["name"], valuesA["name"])

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("PATCH(merge) /{type}/{id}", func(t *testing.T) {
		// PATCH merges the new object with the existing object
		client := testsuit.CreateTestClient()

		testsuit.CreateWeaviateTestSchemaFood(t, client)

		propertySchemaT := map[string]string{
			"name": "Hawaii",
			"description": "Missing description",
		}
		propertySchemaA := map[string]string{
			"name": "ChickenSoup",
			"description": "missing description",
		}
		errCreateT := client.Data.Creator().WithClassName("Pizza").WithID("abefd256-8574-442b-9293-9205193737ee").WithSchema(propertySchemaT).Do(context.Background())
		assert.Nil(t, errCreateT)
		errCreateA := client.Data.Creator().WithClassName("Soup").WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").WithSchema(propertySchemaA).WithKind(paragons.SemanticKindActions).Do(context.Background())
		assert.Nil(t, errCreateA)

		time.Sleep(2.0 * time.Second) // Give weaviate time to update its index
		propertySchemaT = map[string]string{
			"description": "Universally accepted to be the best pizza ever created.",
		}
		updateErrT := client.Data.Updater().WithID("abefd256-8574-442b-9293-9205193737ee").WithClassName("Pizza").WithSchema(propertySchemaT).WithMerge().Do(context.Background())
		assert.Nil(t, updateErrT)

		propertySchemaA = map[string]string{
			"description": "Used by humans when their inferior genetics are attacked by microscopic organisms.",
		}
		updateErrA := client.Data.Updater().WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").WithClassName("Soup").WithSchema(propertySchemaA).WithKind(paragons.SemanticKindActions).WithMerge().Do(context.Background())
		assert.Nil(t, updateErrA)
		time.Sleep(2.0 * time.Second)

		things, getErrT := client.Data.ThingsGetter().WithID("abefd256-8574-442b-9293-9205193737ee").Do(context.Background())
		assert.Nil(t, getErrT)
		valuesT := things[0].Schema.(map[string]interface{})
		assert.Equal(t, propertySchemaT["description"], valuesT["description"])
		assert.Equal(t, "Hawaii", valuesT["name"])

		actions, getErrT := client.Data.ActionsGetter().WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").Do(context.Background())
		assert.Nil(t, getErrT)
		valuesA := actions[0].Schema.(map[string]interface{})
		assert.Equal(t, propertySchemaA["description"], valuesA["description"])
		assert.Equal(t, "ChickenSoup", valuesA["name"])

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("POST /{type}/validate", func(t *testing.T) {
		client := testsuit.CreateTestClient()

		testsuit.CreateWeaviateTestSchemaFood(t, client)

		propertySchemaT := map[string]string{
			"name": "Hawaii",
			"description": "Universally accepted to be the best pizza ever created.",
		}
		propertySchemaA := map[string]string{
			"name": "ChickenSoup",
			"description": "Used by humans when their inferior genetics are attacked by microscopic organisms.",
		}


		errValidateT := client.Data.Validator().WithClassName("Pizza").WithID("abefd256-8574-442b-9293-9205193737ee").WithSchema(propertySchemaT).Do(context.Background())
		assert.Nil(t, errValidateT)

		errValidateA := client.Data.Validator().WithClassName("Soup").WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").WithSchema(propertySchemaA).WithKind(paragons.SemanticKindActions).Do(context.Background())
		assert.Nil(t, errValidateA)


		propertySchemaT["test"] = "not existing property"
		errValidateT = client.Data.Validator().WithClassName("Pizza").WithID("abefd256-8574-442b-9293-9205193737ee").WithSchema(propertySchemaT).Do(context.Background())
		assert.NotNil(t, errValidateT)

		propertySchemaA["test"] = "not existing property"
		errValidateA = client.Data.Validator().WithClassName("Soup").WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").WithSchema(propertySchemaA).WithKind(paragons.SemanticKindActions).Do(context.Background())
		assert.NotNil(t, errValidateA)


		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			fmt.Printf(err.Error())
			t.Fail()
		}
	})
}

