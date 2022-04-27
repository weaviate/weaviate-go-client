package batch

import (
	"context"
	"testing"

	"github.com/semi-technologies/weaviate-go-client/v4/test/testsuit"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/testenv"
	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/stretchr/testify/assert"
)

func TestBatchCreate_integration(t *testing.T) {
	t.Run("up", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatal(err.Error())
		}
	})

	t.Run("POST /batch/{type}", func(t *testing.T) {
		client := testsuit.CreateTestClient()
		testsuit.CreateWeaviateTestSchemaFood(t, client)

		// Create some classes to add in a batch
		propertySchemaT1 := map[string]string{
			"name":        "Hawaii",
			"description": "Universally accepted to be the best pizza ever created.",
		}
		classT1, errPayloadT := client.Data().Creator().WithClassName("Pizza").WithID("abefd256-8574-442b-9293-9205193737ee").WithProperties(propertySchemaT1).PayloadObject()
		assert.Nil(t, errPayloadT)
		classT2 := &models.Object{
			Class: "Pizza",
			ID:    "97fa5147-bdad-4d74-9a81-f8babc811b09",
			Properties: map[string]string{
				"name":        "Doener",
				"description": "A innovation, some say revolution, in the pizza industry.",
			},
		}
		propertySchemaA1 := map[string]string{
			"name":        "Chicken",
			"description": "Used by humans when their inferior genetics are attacked by microscopic organisms.",
		}
		classA1, errPayloadA := client.Data().Creator().WithClassName("Soup").WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").WithProperties(propertySchemaA1).PayloadObject()
		assert.Nil(t, errPayloadA)
		classA2 := &models.Object{
			Class: "Soup",
			ID:    "07473b34-0ab2-4120-882d-303d9e13f7af",
			Properties: map[string]string{
				"name":        "Beautiful",
				"description": "Putting the game of letter soups to a whole new level.",
			},
		}

		batchResultT, batchErrT := client.Batch().ObjectsBatcher().WithObject(classT1).WithObject(classT2).Do(context.Background())
		assert.Nil(t, batchErrT)
		assert.NotNil(t, batchResultT)
		assert.Equal(t, 2, len(batchResultT))
		batchResultA, batchErrA := client.Batch().ObjectsBatcher().WithObject(classA1).WithObject(classA2).Do(context.Background())
		assert.Nil(t, batchErrA)
		assert.NotNil(t, batchResultA)
		assert.Equal(t, 2, len(batchResultA))

		objectT1, objErrT1 := client.Data().ObjectsGetter().WithID("abefd256-8574-442b-9293-9205193737ee").Do(context.Background())
		assert.Nil(t, objErrT1)
		assert.NotNil(t, objectT1)
		objectT2, objErrT2 := client.Data().ObjectsGetter().WithID("97fa5147-bdad-4d74-9a81-f8babc811b09").Do(context.Background())
		assert.Nil(t, objErrT2)
		assert.NotNil(t, objectT2)
		objectA1, objErrA1 := client.Data().ObjectsGetter().WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").Do(context.Background())
		assert.Nil(t, objErrA1)
		assert.NotNil(t, objectA1)
		objectA2, objErrA2 := client.Data().ObjectsGetter().WithID("07473b34-0ab2-4120-882d-303d9e13f7af").Do(context.Background())
		assert.Nil(t, objErrA2)
		assert.NotNil(t, objectA2)

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("POST /batch/references", func(t *testing.T) {
		client := testsuit.CreateTestClient()
		testsuit.CreateWeaviateTestSchemaFoodWithReferenceProperty(t, client)

		// Create some objects
		classT := &models.Object{
			Class: "Pizza",
			ID:    "97fa5147-bdad-4d74-9a81-f8babc811b09",
			Properties: map[string]string{
				"name":        "Doener",
				"description": "A innovation, some say revolution, in the pizza industry.",
			},
		}
		batchResultT, batchErrT := client.Batch().ObjectsBatcher().WithObject(classT).Do(context.Background())
		assert.Nil(t, batchErrT)
		assert.NotNil(t, batchResultT)
		classA := &models.Object{
			Class: "Soup",
			ID:    "07473b34-0ab2-4120-882d-303d9e13f7af",
			Properties: map[string]string{
				"name":        "Beautiful",
				"description": "Putting the game of letter soups to a whole new level.",
			},
		}
		batchResultA, batchErrA := client.Batch().ObjectsBatcher().WithObject(classA).Do(context.Background())
		assert.Nil(t, batchErrA)
		assert.NotNil(t, batchResultA)

		// Define references
		refTtoA := &models.BatchReference{
			From: "weaviate://localhost/Pizza/97fa5147-bdad-4d74-9a81-f8babc811b09/otherFoods",
			To:   "weaviate://localhost/07473b34-0ab2-4120-882d-303d9e13f7af",
		}
		refTtoT := client.Batch().ReferencePayloadBuilder().WithFromClassName("Pizza").WithFromRefProp("otherFoods").WithFromID("97fa5147-bdad-4d74-9a81-f8babc811b09").WithToID("97fa5147-bdad-4d74-9a81-f8babc811b09").Payload()

		refAtoT := &models.BatchReference{
			From: "weaviate://localhost/Soup/07473b34-0ab2-4120-882d-303d9e13f7af/otherFoods",
			To:   "weaviate://localhost/97fa5147-bdad-4d74-9a81-f8babc811b09",
		}
		refAtoA := client.Batch().ReferencePayloadBuilder().WithFromClassName("Soup").WithFromRefProp("otherFoods").WithFromID("07473b34-0ab2-4120-882d-303d9e13f7af").WithToID("07473b34-0ab2-4120-882d-303d9e13f7af").Payload()

		// Add references in batch
		referenceBatchResult, err := client.Batch().ReferencesBatcher().WithReference(refTtoA).WithReference(refTtoT).WithReference(refAtoT).WithReference(refAtoA).Do(context.Background())
		assert.Nil(t, err)
		assert.NotNil(t, referenceBatchResult)

		// Assert
		objectT, objErrT := client.Data().ObjectsGetter().WithID("97fa5147-bdad-4d74-9a81-f8babc811b09").Do(context.Background())
		assert.Nil(t, objErrT)
		valuesT := objectT[0].Properties.(map[string]interface{})
		referencesT := testsuit.ParseReferenceResponseToStruct(t, valuesT["otherFoods"])
		assert.Equal(t, 2, len(referencesT))
		beaconsT := []string{string(referencesT[0].Beacon), string(referencesT[1].Beacon)}
		assert.Contains(t, beaconsT, "weaviate://localhost/07473b34-0ab2-4120-882d-303d9e13f7af")
		assert.Contains(t, beaconsT, "weaviate://localhost/97fa5147-bdad-4d74-9a81-f8babc811b09")

		objectA, objErrA := client.Data().ObjectsGetter().WithID("07473b34-0ab2-4120-882d-303d9e13f7af").Do(context.Background())
		assert.Nil(t, objErrA)
		valuesA := objectA[0].Properties.(map[string]interface{})
		referencesA := testsuit.ParseReferenceResponseToStruct(t, valuesA["otherFoods"])
		assert.Equal(t, 2, len(referencesA))
		beaconsA := []string{string(referencesA[0].Beacon), string(referencesA[1].Beacon)}
		assert.Contains(t, beaconsA, "weaviate://localhost/07473b34-0ab2-4120-882d-303d9e13f7af")
		assert.Contains(t, beaconsA, "weaviate://localhost/97fa5147-bdad-4d74-9a81-f8babc811b09")

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			t.Fatal(err.Error())
		}
	})
}
