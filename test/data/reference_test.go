package data

import (
	"context"
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient"
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
		client := createTestClient()
		createWeaviateTestSchemaFoodWithReferenceProperty(t, client)

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
		referencesT := parseReferenceResponseToStruct(t, valuesT["otherFoods"])
		assert.Equal(t, strfmt.URI("weaviate://localhost/actions/565da3b6-60b3-40e5-ba21-e6bfe5dbba91"), referencesT[0].Beacon)

		actions, getErrA := client.Data.ActionGetter().WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").Do(context.Background())
		assert.Nil(t, getErrA)
		valuesA := actions[0].Schema.(map[string]interface{})
		assert.Contains(t, valuesA, "otherFoods")
		referencesA := parseReferenceResponseToStruct(t, valuesA["otherFoods"])
		assert.Equal(t, strfmt.URI("weaviate://localhost/things/abefd256-8574-442b-9293-9205193737ee"), referencesA[0].Beacon)


		cleanUpWeaviate(t, client)
	})

	t.Run("PUT /{type}/{id}/references/{propertyName}", func(t *testing.T) {
		client := createTestClient()
		createWeaviateTestSchemaFoodWithReferenceProperty(t, client)

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
		referencesT := parseReferenceResponseToStruct(t, valuesT["otherFoods"])
		assert.Equal(t, strfmt.URI("weaviate://localhost/things/abefd256-8574-442b-9293-9205193737ee"), referencesT[0].Beacon)

		actions, getErrA := client.Data.ActionGetter().WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").Do(context.Background())
		assert.Nil(t, getErrA)
		valuesA := actions[0].Schema.(map[string]interface{})
		assert.Contains(t, valuesA, "otherFoods")
		referencesA := parseReferenceResponseToStruct(t, valuesA["otherFoods"])
		assert.Equal(t, strfmt.URI("weaviate://localhost/actions/565da3b6-60b3-40e5-ba21-e6bfe5dbba91"), referencesA[0].Beacon)



		cleanUpWeaviate(t, client)
	})

	t.Run("DELETE /{type}/{id}/references/{propertyName}", func(t *testing.T) {
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

func parseReferenceResponseToStruct(t *testing.T, reference interface{}) models.MultipleRef {
	referenceList := reference.([]interface{})
	out := make(models.MultipleRef, len(referenceList))
	for i, untyped := range referenceList {
		asMap, ok := untyped.(map[string]interface{})
		assert.True(t, ok)
		beacon, ok := asMap["beacon"]
		assert.True(t, ok)
		beaconString, ok := beacon.(string)
		assert.True(t, ok)
		out[i] = &models.SingleRef{
			Beacon: strfmt.URI(beaconString),
		}
	}
	return out
}

func createWeaviateTestSchemaFoodWithReferenceProperty(t *testing.T, client *weaviateclient.WeaviateClient) {
	createWeaviateTestSchemaFood(t, client)
	referenceProperty := models.Property {
		DataType: []string{"Pizza", "Soup"},
		Description: "reference to other foods",
		Name: "otherFoods",
	}
	err := client.Schema.PropertyCreator().WithClassName("Pizza").WithProperty(referenceProperty).Do(context.Background())
	assert.Nil(t, err)
	err = client.Schema.PropertyCreator().WithClassName("Soup").WithProperty(referenceProperty).WithKind(paragons.SemanticKindActions).Do(context.Background())
	assert.Nil(t, err)
}