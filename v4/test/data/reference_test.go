package data

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/semi-technologies/weaviate-go-client/v4/test/testsuit"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/testenv"
	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/stretchr/testify/assert"
)

func TestData_reference_integration(t *testing.T) {
	t.Run("up", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			fmt.Printf(err.Error())
			t.Fail()
		}
	})

	t.Run("POST /{type}/{className}/{id}/references/{propertyName}", func(t *testing.T) {
		client := testsuit.CreateTestClient(8080, nil)
		testsuit.CreateWeaviateTestSchemaFoodWithReferenceProperty(t, client)

		propertySchemaT := map[string]string{
			"name":        "Hawaii",
			"description": "Universally accepted to be the best pizza ever created.",
		}
		propertySchemaA := map[string]string{
			"name":        "ChickenSoup",
			"description": "Used by humans when their inferior genetics are attacked by microscopic organisms.",
		}
		_, errCreateT := client.Data().Creator().WithClassName("Pizza").WithID("abefd256-8574-442b-9293-9205193737ee").WithProperties(propertySchemaT).Do(context.Background())
		assert.Nil(t, errCreateT)
		_, errCreateA := client.Data().Creator().WithClassName("Soup").WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").WithProperties(propertySchemaA).Do(context.Background())
		assert.Nil(t, errCreateA)

		// Thing -> Action
		// Payload to reference the ChickenSoup
		payload := client.Data().ReferencePayloadBuilder().WithClassName("Soup").WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").Payload()
		// Add the reference to the ChickenSoup to the Pizza OtherFoods reference
		refErr := client.Data().ReferenceCreator().WithClassName("Pizza").WithID("abefd256-8574-442b-9293-9205193737ee").WithReferenceProperty("otherFoods").WithReference(payload).Do(context.Background())
		assert.Nil(t, refErr)

		// Action -> Thing
		// Payload to reference the ChickenSoup
		payload = client.Data().ReferencePayloadBuilder().WithClassName("Pizza").WithID("abefd256-8574-442b-9293-9205193737ee").Payload()
		// Add the reference to the ChickenSoup to the Pizza OtherFoods reference
		refErr = client.Data().ReferenceCreator().WithClassName("Soup").WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").WithReferenceProperty("otherFoods").WithReference(payload).Do(context.Background())
		assert.Nil(t, refErr)

		things, getErrT := client.Data().ObjectsGetter().WithClassName("Pizza").WithID("abefd256-8574-442b-9293-9205193737ee").Do(context.Background())
		assert.Nil(t, getErrT)
		valuesT := things[0].Properties.(map[string]interface{})
		assert.Contains(t, valuesT, "otherFoods")
		referencesT := testsuit.ParseReferenceResponseToStruct(t, valuesT["otherFoods"])
		assert.Equal(t, strfmt.URI("weaviate://localhost/Soup/565da3b6-60b3-40e5-ba21-e6bfe5dbba91"), referencesT[0].Beacon)

		actions, getErrA := client.Data().ObjectsGetter().WithClassName("Soup").WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").Do(context.Background())
		assert.Nil(t, getErrA)
		valuesA := actions[0].Properties.(map[string]interface{})
		assert.Contains(t, valuesA, "otherFoods")
		referencesA := testsuit.ParseReferenceResponseToStruct(t, valuesA["otherFoods"])
		assert.Equal(t, strfmt.URI("weaviate://localhost/Pizza/abefd256-8574-442b-9293-9205193737ee"), referencesA[0].Beacon)

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("PUT /{type}/{className}/{id}/references/{propertyName}", func(t *testing.T) {
		client := testsuit.CreateTestClient(8080, nil)
		testsuit.CreateWeaviateTestSchemaFoodWithReferenceProperty(t, client)

		// Create things with references
		propertySchemaT := map[string]string{
			"name":        "Hawaii",
			"description": "Universally accepted to be the best pizza ever created.",
		}
		propertySchemaA := map[string]string{
			"name":        "ChickenSoup",
			"description": "Used by humans when their inferior genetics are attacked by microscopic organisms.",
		}
		_, errCreateT := client.Data().Creator().WithClassName("Pizza").WithID("abefd256-8574-442b-9293-9205193737ee").WithProperties(propertySchemaT).Do(context.Background())
		assert.Nil(t, errCreateT)
		_, errCreateA := client.Data().Creator().WithClassName("Soup").WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").WithProperties(propertySchemaA).Do(context.Background())
		assert.Nil(t, errCreateA)
		// Thing -> Action
		// Payload to reference the ChickenSoup
		payload := client.Data().ReferencePayloadBuilder().WithClassName("Soup").WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").Payload()
		// Add the reference to the ChickenSoup to the Pizza OtherFoods reference
		refErr := client.Data().ReferenceCreator().WithClassName("Pizza").WithID("abefd256-8574-442b-9293-9205193737ee").WithReferenceProperty("otherFoods").WithReference(payload).Do(context.Background())
		assert.Nil(t, refErr)
		// Action -> Thing
		// Payload to reference the ChickenSoup
		payload = client.Data().ReferencePayloadBuilder().WithClassName("Pizza").WithID("abefd256-8574-442b-9293-9205193737ee").Payload()
		// Add the reference to the ChickenSoup to the Pizza OtherFoods reference
		refErr = client.Data().ReferenceCreator().WithClassName("Soup").WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").WithReferenceProperty("otherFoods").WithReference(payload).Do(context.Background())
		assert.Nil(t, refErr)

		// Replace the above reference with self references

		// Thing -> Thing
		payload2 := client.Data().ReferencePayloadBuilder().WithClassName("Pizza").WithID("abefd256-8574-442b-9293-9205193737ee").Payload()
		refList := models.MultipleRef{
			payload2,
		}
		refErr2 := client.Data().ReferenceReplacer().WithClassName("Pizza").WithID("abefd256-8574-442b-9293-9205193737ee").WithReferenceProperty("otherFoods").WithReferences(&refList).Do(context.Background())
		assert.Nil(t, refErr2)
		// Action -> Action
		payload2 = client.Data().ReferencePayloadBuilder().WithClassName("Soup").WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").Payload()
		refList = models.MultipleRef{
			payload2,
		}
		refErr = client.Data().ReferenceReplacer().WithClassName("Soup").WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").WithReferenceProperty("otherFoods").WithReferences(&refList).Do(context.Background())
		assert.Nil(t, refErr)

		things, getErrT := client.Data().ObjectsGetter().WithClassName("Pizza").WithID("abefd256-8574-442b-9293-9205193737ee").Do(context.Background())
		assert.Nil(t, getErrT)
		valuesT := things[0].Properties.(map[string]interface{})
		assert.Contains(t, valuesT, "otherFoods")
		referencesT := testsuit.ParseReferenceResponseToStruct(t, valuesT["otherFoods"])
		assert.Equal(t, strfmt.URI("weaviate://localhost/Pizza/abefd256-8574-442b-9293-9205193737ee"), referencesT[0].Beacon)

		actions, getErrA := client.Data().ObjectsGetter().WithClassName("Soup").WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").Do(context.Background())
		assert.Nil(t, getErrA)
		valuesA := actions[0].Properties.(map[string]interface{})
		assert.Contains(t, valuesA, "otherFoods")
		referencesA := testsuit.ParseReferenceResponseToStruct(t, valuesA["otherFoods"])
		assert.Equal(t, strfmt.URI("weaviate://localhost/Soup/565da3b6-60b3-40e5-ba21-e6bfe5dbba91"), referencesA[0].Beacon)

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("DELETE /{type}/{className}/{id}/references/{propertyName}", func(t *testing.T) {
		client := testsuit.CreateTestClient(8080, nil)
		testsuit.CreateWeaviateTestSchemaFoodWithReferenceProperty(t, client)

		// Create things with references
		propertySchemaT := map[string]string{
			"name":        "Hawaii",
			"description": "Universally accepted to be the best pizza ever created.",
		}
		propertySchemaA := map[string]string{
			"name":        "ChickenSoup",
			"description": "Used by humans when their inferior genetics are attacked by microscopic organisms.",
		}
		_, errCreateT := client.Data().Creator().WithClassName("Pizza").WithID("abefd256-8574-442b-9293-9205193737ee").WithProperties(propertySchemaT).Do(context.Background())
		assert.Nil(t, errCreateT)
		_, errCreateA := client.Data().Creator().WithClassName("Soup").WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").WithProperties(propertySchemaA).Do(context.Background())
		assert.Nil(t, errCreateA)
		// Thing -> Action
		// Payload to reference the ChickenSoup
		payload1 := client.Data().ReferencePayloadBuilder().WithClassName("Soup").WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").Payload()
		// Add the reference to the ChickenSoup to the Pizza OtherFoods reference
		refErr := client.Data().ReferenceCreator().WithClassName("Pizza").WithID("abefd256-8574-442b-9293-9205193737ee").WithReferenceProperty("otherFoods").WithReference(payload1).Do(context.Background())
		assert.Nil(t, refErr)
		// Action -> Thing
		// Payload to reference the ChickenSoup
		payload2 := client.Data().ReferencePayloadBuilder().WithClassName("Pizza").WithID("abefd256-8574-442b-9293-9205193737ee").Payload()
		// Add the reference to the ChickenSoup to the Pizza OtherFoods reference
		refErr = client.Data().ReferenceCreator().WithClassName("Soup").WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").WithReferenceProperty("otherFoods").WithReference(payload2).Do(context.Background())
		assert.Nil(t, refErr)

		client.Data().ReferenceDeleter().WithClassName("Pizza").WithID("abefd256-8574-442b-9293-9205193737ee").WithReferenceProperty("otherFoods").WithReference(payload1).Do(context.Background())
		client.Data().ReferenceDeleter().WithClassName("Soup").WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").WithReferenceProperty("otherFoods").WithReference(payload2).Do(context.Background())

		things, getErrT := client.Data().ObjectsGetter().WithClassName("Pizza").WithID("abefd256-8574-442b-9293-9205193737ee").Do(context.Background())
		assert.Nil(t, getErrT)
		valuesT := things[0].Properties.(map[string]interface{})
		referencesT := testsuit.ParseReferenceResponseToStruct(t, valuesT["otherFoods"])
		assert.Equal(t, 0, len(referencesT))

		actions, getErrA := client.Data().ObjectsGetter().WithClassName("Soup").WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").Do(context.Background())
		assert.Nil(t, getErrA)
		valuesA := actions[0].Properties.(map[string]interface{})
		referencesA := testsuit.ParseReferenceResponseToStruct(t, valuesA["otherFoods"])
		assert.Equal(t, 0, len(referencesA))

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
