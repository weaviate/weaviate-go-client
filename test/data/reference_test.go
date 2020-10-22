package data

import (
	"context"
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/semi-technologies/weaviate-go-client/test/testsuit"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/paragons"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/testenv"
	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestData_reference_integration(t *testing.T) {

	t.Run("up", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			fmt.Printf(err.Error())
			t.Fail()
		}
	})

	t.Run("POST /{type}/{id}/references/{propertyName}", func(t *testing.T) {
		client := testsuit.CreateTestClient()
		testsuit.CreateWeaviateTestSchemaFoodWithReferenceProperty(t, client)

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

		time.Sleep(2.0 * time.Second)
		// Thing -> Action
		// Payload to reference the ChickenSoup
		payload := client.Data.ReferencePayloadBuilder().WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").WithKind(paragons.SemanticKindActions).Payload()
		// Add the reference to the ChickenSoup to the Pizza OtherFoods reference
		refErr := client.Data.ReferenceCreator().WithID("abefd256-8574-442b-9293-9205193737ee").WithReferenceProperty("otherFoods").WithReference(payload).Do(context.Background())
		assert.Nil(t, refErr)

		// Action -> Thing
		// Payload to reference the ChickenSoup
		payload = client.Data.ReferencePayloadBuilder().WithID("abefd256-8574-442b-9293-9205193737ee").Payload()
		// Add the reference to the ChickenSoup to the Pizza OtherFoods reference
		refErr = client.Data.ReferenceCreator().WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").WithKind(paragons.SemanticKindActions).WithReferenceProperty("otherFoods").WithReference(payload).Do(context.Background())
		assert.Nil(t, refErr)

		time.Sleep(2.0 * time.Second)


		things, getErrT := client.Data.ThingGetter().WithID("abefd256-8574-442b-9293-9205193737ee").Do(context.Background())
		assert.Nil(t, getErrT)
		valuesT := things[0].Schema.(map[string]interface{})
		assert.Contains(t, valuesT, "otherFoods")
		referencesT := testsuit.ParseReferenceResponseToStruct(t, valuesT["otherFoods"])
		assert.Equal(t, strfmt.URI("weaviate://localhost/actions/565da3b6-60b3-40e5-ba21-e6bfe5dbba91"), referencesT[0].Beacon)

		actions, getErrA := client.Data.ActionGetter().WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").Do(context.Background())
		assert.Nil(t, getErrA)
		valuesA := actions[0].Schema.(map[string]interface{})
		assert.Contains(t, valuesA, "otherFoods")
		referencesA := testsuit.ParseReferenceResponseToStruct(t, valuesA["otherFoods"])
		assert.Equal(t, strfmt.URI("weaviate://localhost/things/abefd256-8574-442b-9293-9205193737ee"), referencesA[0].Beacon)


		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("PUT /{type}/{id}/references/{propertyName}", func(t *testing.T) {
		client := testsuit.CreateTestClient()
		testsuit.CreateWeaviateTestSchemaFoodWithReferenceProperty(t, client)

		// Create things with references
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
		time.Sleep(2.0 * time.Second)
		// Thing -> Action
		// Payload to reference the ChickenSoup
		payload := client.Data.ReferencePayloadBuilder().WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").WithKind(paragons.SemanticKindActions).Payload()
		// Add the reference to the ChickenSoup to the Pizza OtherFoods reference
		refErr := client.Data.ReferenceCreator().WithID("abefd256-8574-442b-9293-9205193737ee").WithReferenceProperty("otherFoods").WithReference(payload).Do(context.Background())
		assert.Nil(t, refErr)
		// Action -> Thing
		// Payload to reference the ChickenSoup
		payload = client.Data.ReferencePayloadBuilder().WithID("abefd256-8574-442b-9293-9205193737ee").Payload()
		// Add the reference to the ChickenSoup to the Pizza OtherFoods reference
		refErr = client.Data.ReferenceCreator().WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").WithKind(paragons.SemanticKindActions).WithReferenceProperty("otherFoods").WithReference(payload).Do(context.Background())
		assert.Nil(t, refErr)
		time.Sleep(2.0 * time.Second)

		// Replace the above reference with self references

		// Thing -> Thing
		payload2 := client.Data.ReferencePayloadBuilder().WithID("abefd256-8574-442b-9293-9205193737ee").Payload()
		refList := models.MultipleRef{
			payload2,
		}
		refErr2 := client.Data.ReferenceReplacer().WithID("abefd256-8574-442b-9293-9205193737ee").WithReferenceProperty("otherFoods").WithReferences(&refList).Do(context.Background())
		assert.Nil(t, refErr2)
		// Action -> Action
		payload2 = client.Data.ReferencePayloadBuilder().WithKind(paragons.SemanticKindActions).WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").Payload()
		refList = models.MultipleRef{
			payload2,
		}
		refErr = client.Data.ReferenceReplacer().WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").WithKind(paragons.SemanticKindActions).WithReferenceProperty("otherFoods").WithReferences(&refList).Do(context.Background())
		assert.Nil(t, refErr)
		time.Sleep(2.0 * time.Second)


		things, getErrT := client.Data.ThingGetter().WithID("abefd256-8574-442b-9293-9205193737ee").Do(context.Background())
		assert.Nil(t, getErrT)
		valuesT := things[0].Schema.(map[string]interface{})
		assert.Contains(t, valuesT, "otherFoods")
		referencesT := testsuit.ParseReferenceResponseToStruct(t, valuesT["otherFoods"])
		assert.Equal(t, strfmt.URI("weaviate://localhost/things/abefd256-8574-442b-9293-9205193737ee"), referencesT[0].Beacon)

		actions, getErrA := client.Data.ActionGetter().WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").Do(context.Background())
		assert.Nil(t, getErrA)
		valuesA := actions[0].Schema.(map[string]interface{})
		assert.Contains(t, valuesA, "otherFoods")
		referencesA := testsuit.ParseReferenceResponseToStruct(t, valuesA["otherFoods"])
		assert.Equal(t, strfmt.URI("weaviate://localhost/actions/565da3b6-60b3-40e5-ba21-e6bfe5dbba91"), referencesA[0].Beacon)



		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("DELETE /{type}/{id}/references/{propertyName}", func(t *testing.T) {
		client := testsuit.CreateTestClient()
		testsuit.CreateWeaviateTestSchemaFoodWithReferenceProperty(t, client)

		// Create things with references
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
		time.Sleep(2.0 * time.Second)
		// Thing -> Action
		// Payload to reference the ChickenSoup
		payload1 := client.Data.ReferencePayloadBuilder().WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").WithKind(paragons.SemanticKindActions).Payload()
		// Add the reference to the ChickenSoup to the Pizza OtherFoods reference
		refErr := client.Data.ReferenceCreator().WithID("abefd256-8574-442b-9293-9205193737ee").WithReferenceProperty("otherFoods").WithReference(payload1).Do(context.Background())
		assert.Nil(t, refErr)
		// Action -> Thing
		// Payload to reference the ChickenSoup
		payload2 := client.Data.ReferencePayloadBuilder().WithID("abefd256-8574-442b-9293-9205193737ee").Payload()
		// Add the reference to the ChickenSoup to the Pizza OtherFoods reference
		refErr = client.Data.ReferenceCreator().WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").WithKind(paragons.SemanticKindActions).WithReferenceProperty("otherFoods").WithReference(payload2).Do(context.Background())
		assert.Nil(t, refErr)
		time.Sleep(2.0 * time.Second)

		client.Data.ReferenceDeleter().WithID("abefd256-8574-442b-9293-9205193737ee").WithReferenceProperty("otherFoods").WithReference(payload1).Do(context.Background())
		client.Data.ReferenceDeleter().WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").WithReferenceProperty("otherFoods").WithReference(payload2).WithKind(paragons.SemanticKindActions).Do(context.Background())
		time.Sleep(2.0 * time.Second)

		things, getErrT := client.Data.ThingGetter().WithID("abefd256-8574-442b-9293-9205193737ee").Do(context.Background())
		assert.Nil(t, getErrT)
		valuesT := things[0].Schema.(map[string]interface{})
		referencesT := testsuit.ParseReferenceResponseToStruct(t, valuesT["otherFoods"])
		assert.Equal(t, 0, len(referencesT))

		actions, getErrA := client.Data.ActionGetter().WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").Do(context.Background())
		assert.Nil(t, getErrA)
		valuesA := actions[0].Schema.(map[string]interface{})
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

