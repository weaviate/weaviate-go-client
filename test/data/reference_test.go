package data

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/weaviate/weaviate-go-client/v4/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/testenv"
	"github.com/weaviate/weaviate/entities/models"
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

	t.Run("POST /{type}/{className}/{id}/references/{propertyName}?consistency_level={level}", func(t *testing.T) {
		client := testsuit.CreateTestClient(8080, nil)
		testsuit.CreateWeaviateTestSchemaFoodWithReferenceProperty(t, client)

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

		_, errCreate1 := client.Data().Creator().WithClassName("Pizza").WithID(id1).WithProperties(props1).Do(context.Background())
		assert.Nil(t, errCreate1)
		_, errCreate2 := client.Data().Creator().WithClassName("Soup").WithID(id2).WithProperties(props2).Do(context.Background())
		assert.Nil(t, errCreate2)
		_, errCreate3 := client.Data().Creator().WithClassName("Soup").WithID(id3).WithProperties(props3).Do(context.Background())
		assert.Nil(t, errCreate3)

		// Hawaii -> ChickenSoup
		// Payload to reference the ChickenSoup
		payload := client.Data().ReferencePayloadBuilder().WithClassName("Soup").WithID(id2).Payload()
		// Add the reference to the ChickenSoup to the Pizza OtherFoods reference
		refErr := client.Data().ReferenceCreator().
			WithClassName("Pizza").
			WithID(id1).
			WithReferenceProperty("otherFoods").
			WithReference(payload).
			WithConsistencyLevel("ONE").
			Do(context.Background())
		assert.Nil(t, refErr)

		// ChickenSoup -> Hawaii
		// Payload to reference the Hawaii
		payload = client.Data().ReferencePayloadBuilder().WithClassName("Pizza").WithID(id1).Payload()
		// Add the reference to the ChickenSoup to the Pizza OtherFoods reference
		refErr = client.Data().ReferenceCreator().
			WithClassName("Soup").WithID(id2).
			WithReferenceProperty("otherFoods").
			WithReference(payload).
			WithConsistencyLevel("ALL").
			Do(context.Background())
		assert.Nil(t, refErr)

		// Pozole -> Hawaii
		// Payload to reference the Hawaii
		payload = client.Data().ReferencePayloadBuilder().WithClassName("Pizza").WithID(id1).Payload()
		// Add the reference to the Pozole to the Pizza OtherFoods reference
		refErr = client.Data().ReferenceCreator().
			WithClassName("Soup").WithID(id3).
			WithReferenceProperty("otherFoods").
			WithReference(payload).
			WithConsistencyLevel("QUORUM").
			Do(context.Background())
		assert.Nil(t, refErr)

		found1, err := client.Data().ObjectsGetter().WithClassName("Pizza").WithID(id1).Do(context.Background())
		assert.Nil(t, err)
		found1Props := found1[0].Properties.(map[string]interface{})
		assert.Contains(t, found1Props, "otherFoods")
		ref1 := testsuit.ParseReferenceResponseToStruct(t, found1Props["otherFoods"])
		assert.Equal(t, strfmt.URI(fmt.Sprintf("weaviate://localhost/Soup/%s", id2)), ref1[0].Beacon)

		found2, err := client.Data().ObjectsGetter().WithClassName("Soup").WithID(id2).Do(context.Background())
		assert.Nil(t, err)
		found2Props := found2[0].Properties.(map[string]interface{})
		assert.Contains(t, found2Props, "otherFoods")
		ref2 := testsuit.ParseReferenceResponseToStruct(t, found2Props["otherFoods"])
		assert.Equal(t, strfmt.URI(fmt.Sprintf("weaviate://localhost/Pizza/%s", id1)), ref2[0].Beacon)

		found3, err := client.Data().ObjectsGetter().WithClassName("Soup").WithID(id3).Do(context.Background())
		assert.Nil(t, err)
		found3Props := found3[0].Properties.(map[string]interface{})
		assert.Contains(t, found3Props, "otherFoods")
		ref3 := testsuit.ParseReferenceResponseToStruct(t, found3Props["otherFoods"])
		assert.Equal(t, strfmt.URI(fmt.Sprintf("weaviate://localhost/Pizza/%s", id1)), ref3[0].Beacon)

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

	t.Run("PUT /{type}/{className}/{id}/references/{propertyName}?consistency_level={level}", func(t *testing.T) {
		client := testsuit.CreateTestClient(8080, nil)
		testsuit.CreateWeaviateTestSchemaFoodWithReferenceProperty(t, client)

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

		_, errCreate1 := client.Data().Creator().WithClassName("Pizza").WithID(id1).WithProperties(props1).Do(context.Background())
		assert.Nil(t, errCreate1)
		_, errCreate2 := client.Data().Creator().WithClassName("Soup").WithID(id2).WithProperties(props2).Do(context.Background())
		assert.Nil(t, errCreate2)
		_, errCreate3 := client.Data().Creator().WithClassName("Soup").WithID(id3).WithProperties(props3).Do(context.Background())
		assert.Nil(t, errCreate3)

		// Hawaii -> ChickenSoup
		// Payload to reference the ChickenSoup
		payload1 := client.Data().ReferencePayloadBuilder().WithClassName("Soup").WithID(id2).Payload()
		// Add the reference to the ChickenSoup to the Pizza OtherFoods reference
		refErr := client.Data().ReferenceCreator().
			WithClassName("Pizza").
			WithID(id1).
			WithReferenceProperty("otherFoods").
			WithReference(payload1).
			Do(context.Background())
		assert.Nil(t, refErr)

		// ChickenSoup -> Hawaii
		// Payload to reference the Hawaii
		payload2 := client.Data().ReferencePayloadBuilder().WithClassName("Pizza").WithID(id1).Payload()
		// Add the reference to the ChickenSoup to the Pizza OtherFoods reference
		refErr = client.Data().ReferenceCreator().
			WithClassName("Soup").
			WithID(id2).
			WithReferenceProperty("otherFoods").
			WithReference(payload2).
			Do(context.Background())
		assert.Nil(t, refErr)

		// Pozole -> Hawaii
		// Payload to reference the Hawaii
		payload3 := client.Data().ReferencePayloadBuilder().WithClassName("Pizza").WithID(id1).Payload()
		// Add the reference to the Pozole to the Pizza OtherFoods reference
		refErr = client.Data().ReferenceCreator().
			WithClassName("Soup").WithID(id3).
			WithReferenceProperty("otherFoods").
			WithReference(payload3).
			Do(context.Background())
		assert.Nil(t, refErr)

		replacePayload1 := client.Data().ReferencePayloadBuilder().WithClassName("Pizza").WithID(id1).Payload()
		replaceErr1 := client.Data().ReferenceReplacer().
			WithClassName("Pizza").
			WithID(id1).
			WithReferenceProperty("otherFoods").
			WithReferences(&models.MultipleRef{replacePayload1}).
			Do(context.Background())
		assert.Nil(t, replaceErr1)

		replacePayload2 := client.Data().ReferencePayloadBuilder().WithClassName("Soup").WithID(id2).Payload()
		replaceErr2 := client.Data().ReferenceReplacer().
			WithClassName("Soup").
			WithID(id2).
			WithReferenceProperty("otherFoods").
			WithReferences(&models.MultipleRef{replacePayload2}).
			Do(context.Background())
		assert.Nil(t, replaceErr2)

		replacePayload3 := client.Data().ReferencePayloadBuilder().WithClassName("Soup").WithID(id3).Payload()
		replaceErr3 := client.Data().ReferenceReplacer().
			WithClassName("Soup").
			WithID(id3).
			WithReferenceProperty("otherFoods").
			WithReferences(&models.MultipleRef{replacePayload3}).
			Do(context.Background())
		assert.Nil(t, replaceErr3)

		found1, err := client.Data().ObjectsGetter().WithClassName("Pizza").WithID(id1).Do(context.Background())
		assert.Nil(t, err)
		found1Props := found1[0].Properties.(map[string]interface{})
		assert.Contains(t, found1Props, "otherFoods")
		ref1 := testsuit.ParseReferenceResponseToStruct(t, found1Props["otherFoods"])
		assert.Equal(t, strfmt.URI(fmt.Sprintf("weaviate://localhost/Pizza/%s", id1)), ref1[0].Beacon)

		found2, err := client.Data().ObjectsGetter().WithClassName("Soup").WithID(id2).Do(context.Background())
		assert.Nil(t, err)
		found2Props := found2[0].Properties.(map[string]interface{})
		assert.Contains(t, found1Props, "otherFoods")
		ref2 := testsuit.ParseReferenceResponseToStruct(t, found2Props["otherFoods"])
		assert.Equal(t, strfmt.URI(fmt.Sprintf("weaviate://localhost/Soup/%s", id2)), ref2[0].Beacon)

		found3, err := client.Data().ObjectsGetter().WithClassName("Soup").WithID(id3).Do(context.Background())
		assert.Nil(t, err)
		found3Props := found3[0].Properties.(map[string]interface{})
		assert.Contains(t, found1Props, "otherFoods")
		ref3 := testsuit.ParseReferenceResponseToStruct(t, found3Props["otherFoods"])
		assert.Equal(t, strfmt.URI(fmt.Sprintf("weaviate://localhost/Soup/%s", id3)), ref3[0].Beacon)

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

	t.Run("DELETE /{type}/{className}/{id}/references/{propertyName}?consistency_level={level}", func(t *testing.T) {
		client := testsuit.CreateTestClient(8080, nil)
		testsuit.CreateWeaviateTestSchemaFoodWithReferenceProperty(t, client)

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

		_, errCreate1 := client.Data().Creator().WithClassName("Pizza").WithID(id1).WithProperties(props1).Do(context.Background())
		assert.Nil(t, errCreate1)
		_, errCreate2 := client.Data().Creator().WithClassName("Soup").WithID(id2).WithProperties(props2).Do(context.Background())
		assert.Nil(t, errCreate2)
		_, errCreate3 := client.Data().Creator().WithClassName("Soup").WithID(id3).WithProperties(props3).Do(context.Background())
		assert.Nil(t, errCreate3)

		// Hawaii -> ChickenSoup
		// Payload to reference the ChickenSoup
		payload1 := client.Data().ReferencePayloadBuilder().WithClassName("Soup").WithID(id2).Payload()
		// Add the reference to the ChickenSoup to the Pizza OtherFoods reference
		refErr := client.Data().ReferenceCreator().
			WithClassName("Pizza").
			WithID(id1).
			WithReferenceProperty("otherFoods").
			WithReference(payload1).
			Do(context.Background())
		assert.Nil(t, refErr)

		// ChickenSoup -> Hawaii
		// Payload to reference the Hawaii
		payload2 := client.Data().ReferencePayloadBuilder().WithClassName("Pizza").WithID(id1).Payload()
		// Add the reference to the ChickenSoup to the Pizza OtherFoods reference
		refErr = client.Data().ReferenceCreator().
			WithClassName("Soup").WithID(id2).
			WithReferenceProperty("otherFoods").
			WithReference(payload2).
			Do(context.Background())
		assert.Nil(t, refErr)

		// Pozole -> Hawaii
		// Payload to reference the Hawaii
		payload3 := client.Data().ReferencePayloadBuilder().WithClassName("Pizza").WithID(id1).Payload()
		// Add the reference to the Pozole to the Pizza OtherFoods reference
		refErr = client.Data().ReferenceCreator().
			WithClassName("Soup").WithID(id3).
			WithReferenceProperty("otherFoods").
			WithReference(payload3).
			Do(context.Background())
		assert.Nil(t, refErr)

		client.Data().ReferenceDeleter().
			WithClassName("Pizza").
			WithID(id1).
			WithReferenceProperty("otherFoods").
			WithReference(payload1).
			WithConsistencyLevel("ONE").
			Do(context.Background())
		client.Data().ReferenceDeleter().
			WithClassName("Soup").
			WithID(id2).
			WithReferenceProperty("otherFoods").
			WithReference(payload2).
			WithConsistencyLevel("ALL").
			Do(context.Background())
		client.Data().ReferenceDeleter().
			WithClassName("Soup").
			WithID(id3).
			WithReferenceProperty("otherFoods").
			WithReference(payload2).
			WithConsistencyLevel("QUORUM").
			Do(context.Background())

		found1, err := client.Data().ObjectsGetter().WithClassName("Pizza").WithID(id1).Do(context.Background())
		assert.Nil(t, err)
		found1Props := found1[0].Properties.(map[string]interface{})
		ref1 := testsuit.ParseReferenceResponseToStruct(t, found1Props["otherFoods"])
		assert.Equal(t, 0, len(ref1))

		found2, err := client.Data().ObjectsGetter().WithClassName("Soup").WithID(id2).Do(context.Background())
		assert.Nil(t, err)
		found2Props := found2[0].Properties.(map[string]interface{})
		ref2 := testsuit.ParseReferenceResponseToStruct(t, found2Props["otherFoods"])
		assert.Equal(t, 0, len(ref2))

		found3, err := client.Data().ObjectsGetter().WithClassName("Soup").WithID(id3).Do(context.Background())
		assert.Nil(t, err)
		found3Props := found3[0].Properties.(map[string]interface{})
		ref3 := testsuit.ParseReferenceResponseToStruct(t, found3Props["otherFoods"])
		assert.Equal(t, 0, len(ref3))

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
