package data

import (
	"context"
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/paragons"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/testenv"
	"github.com/semi-technologies/weaviate/entities/models"
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

	t.Run("POST /actions /things", func(t *testing.T) {

		client := createTestClient()

		createWeaviateTestSchemaFood(t, client)

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
		objectT, objErrT := client.Data.ThingGetter().WithID("abefd256-8574-442b-9293-9205193737ee").Do(context.Background())
		assert.Nil(t, objErrT)
		objectA, objErrA := client.Data.ActionGetter().WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").Do(context.Background())
		assert.Nil(t, objErrA)

		assert.Equal(t, "Pizza", objectT.Class)
		valuesT := objectT.Schema.(map[string]interface{})
		assert.Equal(t, "Hawaii", valuesT["name"])
		assert.Equal(t, "Soup", objectA.Class)
		valuesA := objectA.Schema.(map[string]interface{})
		assert.Equal(t, "ChickenSoup", valuesA["name"])

		// Clean up classes
		errRm := client.Schema.AllDeleter().Do(context.Background())
		assert.Nil(t, errRm)
	})

	t.Run("GET /actions /things", func(t *testing.T) {
		// Not implemented to only get thigns or actions without uuid yet

		//client := createTestClient()
		//createWeaviateTestSchemaFood(t, client)
		//
		//errCreate := client.Data.Creator().WithClassName("Pizza").WithSchema(map[string]string{
		//	"name": "Margherita",
		//	"description": "plain",
		//}).Do(context.Background())
		//assert.Nil(t, errCreate)
		//errCreate = client.Data.Creator().WithClassName("Pizza").WithSchema(map[string]string{
		//	"name": "Pepperoni",
		//	"description": "meat",
		//}).Do(context.Background())
		//assert.Nil(t, errCreate)
		//errCreate = client.Data.Creator().WithClassName("Soup").WithKind(paragons.SemanticKindActions).WithSchema(map[string]string{
		//	"name": "Chicken",
		//	"description": "meat",
		//}).Do(context.Background())
		//assert.Nil(t, errCreate)
		//errCreate = client.Data.Creator().WithClassName("Soup").WithKind(paragons.SemanticKindActions).WithSchema(map[string]string{
		//	"name": "Tofu",
		//	"description": "vegetarian",
		//}).Do(context.Background())
		//assert.Nil(t, errCreate)
		//
		//time.Sleep(2.0 * time.Second)
		//objectT, objErrT := client.Data.ThingGetter().Do(context.Background())
		//assert.Nil(t, objErrT)
		//objectA, objErrA := client.Data.ActionGetter().Do(context.Background())
		//assert.Nil(t, objErrA)

	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			fmt.Printf(err.Error())
			t.Fail()
		}
	})
}

func createTestClient() *weaviateclient.WeaviateClient {
	cfg := weaviateclient.Config{
		Host:   "localhost:8080",
		Scheme: "http",
	}
	client := weaviateclient.New(cfg)
	return client
}

func createWeaviateTestSchemaFood(t *testing.T, client *weaviateclient.WeaviateClient) {
	schemaClassThing := &models.Class{
		Class:              "Pizza",
		Description:        "A delicious religion like food and arguably the best export of Italy.",
	}
	schemaClassAction := &models.Class{
		Class:              "Soup",
		Description:        "Mostly water based brew of sustenance for humans.",
	}
	errT := client.Schema.ClassCreator().WithClass(schemaClassThing).Do(context.Background())
	assert.Nil(t, errT)
	errA := client.Schema.ClassCreator().WithClass(schemaClassAction).WithKind(paragons.SemanticKindActions).Do(context.Background())
	assert.Nil(t, errA)
	nameProperty := models.Property{
		DataType:              []string{"string"},
		Description:           "name",
		Name:                  "name",
	}
	descriptionProperty := models.Property{
		DataType:              []string{"string"},
		Description:           "description",
		Name:                  "description",
	}
	propErrT1 := client.Schema.PropertyCreator().WithClassName("Pizza").WithProperty(nameProperty).Do(context.Background())
	assert.Nil(t, propErrT1)
	propErrA1 := client.Schema.PropertyCreator().WithClassName("Soup").WithProperty(nameProperty).WithKind(paragons.SemanticKindActions).Do(context.Background())
	assert.Nil(t, propErrA1)
	propErrT2 := client.Schema.PropertyCreator().WithClassName("Pizza").WithProperty(descriptionProperty).Do(context.Background())
	assert.Nil(t, propErrT2)
	propErrA2 := client.Schema.PropertyCreator().WithClassName("Soup").WithProperty(descriptionProperty).WithKind(paragons.SemanticKindActions).Do(context.Background())
	assert.Nil(t, propErrA2)
}