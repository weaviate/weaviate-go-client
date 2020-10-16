package data

import (
	"context"
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient"
	clientModels "github.com/semi-technologies/weaviate-go-client/weaviateclient/models"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/testenv"
	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/stretchr/testify/assert"
	"testing"
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

		cfg := weaviateclient.Config{
			Host:   "localhost:8080",
			Scheme: "http",
		}
		client := weaviateclient.New(cfg)

		createWeaviateTestSchemaFood(t, client)

		propertySchema := map[string]string{
			"name": "Hawaii",
			"description": "Universally accepted to be the best pizza ever created.",
		}

		errCreateT := client.Data.Creator().WithClassName("Pizza").WithID("abefd256-8574-442b-9293-9205193737ee").WithSchema(propertySchema).Do(context.Background())
		assert.Nil(t, errCreateT)
		errCreateA := client.Data.Creator().WithClassName("Soup").WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").WithSchema(propertySchema).WithKind(clientModels.SemanticKindActions).Do(context.Background())
		assert.Nil(t, errCreateA)

		objectT, objErrT := client.Data.Getter.WithID("abefd256-8574-442b-9293-9205193737ee").Do(context.Background())
		assert.Nil(t, objErrT)
		objectA, objErrA := client.Data.Getter.WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").Do(context.Background())
		assert.Nil(t, objErrA)


		// TODO also assert object values
		t.Fail()

		// Clean up classes
		errRm := client.Schema.AllDeleter().Do(context.Background())
		assert.Nil(t, errRm)
	})

	t.Run("GET /actions /things", func(t *testing.T) {
		t.Fail()

	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			fmt.Printf(err.Error())
			t.Fail()
		}
	})
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
	errA := client.Schema.ClassCreator().WithClass(schemaClassAction).WithKind(clientModels.SemanticKindActions).Do(context.Background())
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
	propErrA1 := client.Schema.PropertyCreator().WithClassName("ChickenSoup").WithProperty(nameProperty).WithKind(clientModels.SemanticKindActions).Do(context.Background())
	assert.Nil(t, propErrA1)
	propErrT2 := client.Schema.PropertyCreator().WithClassName("Pizza").WithProperty(descriptionProperty).Do(context.Background())
	assert.Nil(t, propErrT2)
	propErrA2 := client.Schema.PropertyCreator().WithClassName("ChickenSoup").WithProperty(descriptionProperty).WithKind(clientModels.SemanticKindActions).Do(context.Background())
	assert.Nil(t, propErrA2)
}