package data

import (
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/paragons"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/testenv"
	"github.com/stretchr/testify/assert"
	"testing"
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



		payload, payloadErr := client.Data.ReferencePayloadBuilder().WithId("6bb06a43-e7f0-393e-9ecf-3c0f4e129064").Payload()
		assert.Nil(t, payloadErr)

		refErr := client.Data.ReferenceCreator().WithId().WithReferenceProperty("").WithReference(payload).Do(context.Background())
		assert.Nil(t, refErr)


		errCreateT := client.Data.Validator().WithClassName("Pizza").WithID("abefd256-8574-442b-9293-9205193737ee").WithSchema(propertySchemaT).Do(context.Background())
		assert.Nil(t, errCreateT)

		errCreateA := client.Data.Validator().WithClassName("Soup").WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").WithSchema(propertySchemaA).WithKind(paragons.SemanticKindActions).Do(context.Background())
		assert.Nil(t, errCreateA)



		cleanUpWeaviate(t, client)
	})

	t.Run("PUT /{type}/{id}/references/{propertyName}", func(t *testing.T) {
		t.Fail()
	})

	t.Run("DELETE /{type}/{id}/references/{propertyName}", func(t *testing.T) {
		t.Fail()
	})


}