package batch

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v4/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/data/replication"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/testenv"
	"github.com/weaviate/weaviate/entities/models"
)

func TestBatchCreate_integration(t *testing.T) {
	t.Run("up", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatal(err.Error())
		}
	})

	t.Run("POST /batch/objects", func(t *testing.T) {
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
		classASlice := []*models.Object{classA1, classA2}

		batchResultT, batchErrT := client.Batch().ObjectsBatcher().WithObject(classT1).WithObject(classT2).Do(context.Background())
		assert.Nil(t, batchErrT)
		assert.NotNil(t, batchResultT)
		assert.Equal(t, 2, len(batchResultT))
		batchResultA, batchErrA := client.Batch().ObjectsBatcher().WithObjects(classA1, classA2).Do(context.Background())
		assert.Nil(t, batchErrA)
		assert.NotNil(t, batchResultA)
		assert.Equal(t, 2, len(batchResultA))

		batchResultSlice, batchErrSlice := client.Batch().ObjectsBatcher().WithObjects(classASlice...).Do(context.Background())
		assert.Nil(t, batchErrSlice)
		assert.NotNil(t, batchResultSlice)
		assert.Equal(t, 2, len(batchResultSlice))

		objectT1, objErrT1 := client.Data().ObjectsGetter().WithClassName("Pizza").WithID("abefd256-8574-442b-9293-9205193737ee").Do(context.Background())
		assert.Nil(t, objErrT1)
		assert.NotNil(t, objectT1)
		objectT2, objErrT2 := client.Data().ObjectsGetter().WithClassName("Pizza").WithID("97fa5147-bdad-4d74-9a81-f8babc811b09").Do(context.Background())
		assert.Nil(t, objErrT2)
		assert.NotNil(t, objectT2)
		objectA1, objErrA1 := client.Data().ObjectsGetter().WithClassName("Soup").WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").Do(context.Background())
		assert.Nil(t, objErrA1)
		assert.NotNil(t, objectA1)
		objectA2, objErrA2 := client.Data().ObjectsGetter().WithClassName("Soup").WithID("07473b34-0ab2-4120-882d-303d9e13f7af").Do(context.Background())
		assert.Nil(t, objErrA2)
		assert.NotNil(t, objectA2)

		testsuit.CleanUpWeaviate(t, client)
	})

	// Testing batch object creation with tunable consistency
	t.Run("POST /batch/{objects}?consistency_level={level}", func(t *testing.T) {
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
		classASlice := []*models.Object{classA1, classA2}

		batchResultT, batchErrT := client.Batch().ObjectsBatcher().WithObject(classT1).WithConsistencyLevel(replication.ConsistencyLevel.QUORUM).WithObject(classT2).Do(context.Background())
		assert.Nil(t, batchErrT)
		assert.NotNil(t, batchResultT)
		assert.Equal(t, 2, len(batchResultT))
		batchResultA, batchErrA := client.Batch().ObjectsBatcher().WithConsistencyLevel(replication.ConsistencyLevel.ONE).WithObjects(classA1, classA2).Do(context.Background())
		assert.Nil(t, batchErrA)
		assert.NotNil(t, batchResultA)
		assert.Equal(t, 2, len(batchResultA))

		batchResultSlice, batchErrSlice := client.Batch().ObjectsBatcher().WithConsistencyLevel(replication.ConsistencyLevel.ALL).WithObjects(classASlice...).Do(context.Background())
		assert.Nil(t, batchErrSlice)
		assert.NotNil(t, batchResultSlice)
		assert.Equal(t, 2, len(batchResultSlice))

		objectT1, objErrT1 := client.Data().ObjectsGetter().WithClassName("Pizza").WithID("abefd256-8574-442b-9293-9205193737ee").Do(context.Background())
		assert.Nil(t, objErrT1)
		assert.NotNil(t, objectT1)
		objectT2, objErrT2 := client.Data().ObjectsGetter().WithClassName("Pizza").WithID("97fa5147-bdad-4d74-9a81-f8babc811b09").Do(context.Background())
		assert.Nil(t, objErrT2)
		assert.NotNil(t, objectT2)
		objectA1, objErrA1 := client.Data().ObjectsGetter().WithClassName("Soup").WithID("565da3b6-60b3-40e5-ba21-e6bfe5dbba91").Do(context.Background())
		assert.Nil(t, objErrA1)
		assert.NotNil(t, objectA1)
		objectA2, objErrA2 := client.Data().ObjectsGetter().WithClassName("Soup").WithID("07473b34-0ab2-4120-882d-303d9e13f7af").Do(context.Background())
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
			To:   "weaviate://localhost/Soup/07473b34-0ab2-4120-882d-303d9e13f7af",
		}
		refTtoT := client.Batch().ReferencePayloadBuilder().
			WithFromClassName("Pizza").WithFromRefProp("otherFoods").WithFromID("97fa5147-bdad-4d74-9a81-f8babc811b09").
			WithToClassName("Pizza").WithToID("97fa5147-bdad-4d74-9a81-f8babc811b09").Payload()

		refAtoT := &models.BatchReference{
			From: "weaviate://localhost/Soup/07473b34-0ab2-4120-882d-303d9e13f7af/otherFoods",
			To:   "weaviate://localhost/Pizza/97fa5147-bdad-4d74-9a81-f8babc811b09",
		}
		refAtoA := client.Batch().ReferencePayloadBuilder().
			WithFromClassName("Soup").WithFromRefProp("otherFoods").WithFromID("07473b34-0ab2-4120-882d-303d9e13f7af").
			WithToClassName("Soup").WithToID("07473b34-0ab2-4120-882d-303d9e13f7af").Payload()

		// Add references in batch
		referenceBatchResult, err := client.Batch().ReferencesBatcher().
			WithReference(refTtoA).WithReference(refTtoT).WithReferences(refAtoT, refAtoA).Do(context.Background())
		assert.Nil(t, err)
		assert.NotNil(t, referenceBatchResult)

		// Assert
		objectT, objErrT := client.Data().ObjectsGetter().
			WithClassName("Pizza").WithID("97fa5147-bdad-4d74-9a81-f8babc811b09").Do(context.Background())
		assert.Nil(t, objErrT)
		valuesT := objectT[0].Properties.(map[string]interface{})
		referencesT := testsuit.ParseReferenceResponseToStruct(t, valuesT["otherFoods"])
		assert.Equal(t, 2, len(referencesT))
		beaconsT := []string{string(referencesT[0].Beacon), string(referencesT[1].Beacon)}
		assert.Contains(t, beaconsT, "weaviate://localhost/Soup/07473b34-0ab2-4120-882d-303d9e13f7af")
		assert.Contains(t, beaconsT, "weaviate://localhost/Pizza/97fa5147-bdad-4d74-9a81-f8babc811b09")

		objectA, objErrA := client.Data().ObjectsGetter().
			WithClassName("Soup").WithID("07473b34-0ab2-4120-882d-303d9e13f7af").Do(context.Background())
		assert.Nil(t, objErrA)
		valuesA := objectA[0].Properties.(map[string]interface{})
		referencesA := testsuit.ParseReferenceResponseToStruct(t, valuesA["otherFoods"])
		assert.Equal(t, 2, len(referencesA))
		beaconsA := []string{string(referencesA[0].Beacon), string(referencesA[1].Beacon)}
		assert.Contains(t, beaconsA, "weaviate://localhost/Soup/07473b34-0ab2-4120-882d-303d9e13f7af")
		assert.Contains(t, beaconsA, "weaviate://localhost/Pizza/97fa5147-bdad-4d74-9a81-f8babc811b09")

		testsuit.CleanUpWeaviate(t, client)
	})

	// Testing batch reference creation with tunable consistency
	t.Run("POST /batch/references?consistency_level={level}", func(t *testing.T) {
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
		batchResultT, batchErrT := client.Batch().ObjectsBatcher().
			WithObject(classT).WithConsistencyLevel(replication.ConsistencyLevel.ONE).Do(context.Background())
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
		batchResultA, batchErrA := client.Batch().ObjectsBatcher().
			WithObject(classA).WithConsistencyLevel(replication.ConsistencyLevel.QUORUM).Do(context.Background())
		assert.Nil(t, batchErrA)
		assert.NotNil(t, batchResultA)

		// Define references
		refTtoA := &models.BatchReference{
			From: "weaviate://localhost/Pizza/97fa5147-bdad-4d74-9a81-f8babc811b09/otherFoods",
			To:   "weaviate://localhost/Soup/07473b34-0ab2-4120-882d-303d9e13f7af",
		}
		refTtoT := client.Batch().ReferencePayloadBuilder().
			WithFromClassName("Pizza").WithFromRefProp("otherFoods").WithFromID("97fa5147-bdad-4d74-9a81-f8babc811b09").
			WithToClassName("Pizza").WithToID("97fa5147-bdad-4d74-9a81-f8babc811b09").Payload()

		refAtoT := &models.BatchReference{
			From: "weaviate://localhost/Soup/07473b34-0ab2-4120-882d-303d9e13f7af/otherFoods",
			To:   "weaviate://localhost/Pizza/97fa5147-bdad-4d74-9a81-f8babc811b09",
		}
		refAtoA := client.Batch().ReferencePayloadBuilder().
			WithFromClassName("Soup").WithFromRefProp("otherFoods").WithFromID("07473b34-0ab2-4120-882d-303d9e13f7af").
			WithToClassName("Soup").WithToID("07473b34-0ab2-4120-882d-303d9e13f7af").Payload()

		// Add references in batch
		referenceBatchResult, err := client.Batch().ReferencesBatcher().WithConsistencyLevel(replication.ConsistencyLevel.ALL).
			WithReference(refTtoA).WithReference(refTtoT).WithReferences(refAtoT, refAtoA).Do(context.Background())
		assert.Nil(t, err)
		assert.NotNil(t, referenceBatchResult)

		// Assert
		objectT, objErrT := client.Data().ObjectsGetter().
			WithClassName("Pizza").WithID("97fa5147-bdad-4d74-9a81-f8babc811b09").Do(context.Background())
		assert.Nil(t, objErrT)
		valuesT := objectT[0].Properties.(map[string]interface{})
		referencesT := testsuit.ParseReferenceResponseToStruct(t, valuesT["otherFoods"])
		assert.Equal(t, 2, len(referencesT))
		beaconsT := []string{string(referencesT[0].Beacon), string(referencesT[1].Beacon)}
		assert.Contains(t, beaconsT, "weaviate://localhost/Soup/07473b34-0ab2-4120-882d-303d9e13f7af")
		assert.Contains(t, beaconsT, "weaviate://localhost/Pizza/97fa5147-bdad-4d74-9a81-f8babc811b09")

		objectA, objErrA := client.Data().ObjectsGetter().
			WithClassName("Soup").WithID("07473b34-0ab2-4120-882d-303d9e13f7af").Do(context.Background())
		assert.Nil(t, objErrA)
		valuesA := objectA[0].Properties.(map[string]interface{})
		referencesA := testsuit.ParseReferenceResponseToStruct(t, valuesA["otherFoods"])
		assert.Equal(t, 2, len(referencesA))
		beaconsA := []string{string(referencesA[0].Beacon), string(referencesA[1].Beacon)}
		assert.Contains(t, beaconsA, "weaviate://localhost/Soup/07473b34-0ab2-4120-882d-303d9e13f7af")
		assert.Contains(t, beaconsA, "weaviate://localhost/Pizza/97fa5147-bdad-4d74-9a81-f8babc811b09")

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			t.Fatal(err.Error())
		}
	})
}

func TestBatchCreate_MultiTenancy(t *testing.T) {
	cleanup := func() {
		client := testsuit.CreateTestClient()
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
		t.Skip("to be enabled")
		defer cleanup()

		tenants := []string{"tenantNo1", "tenantNo2"}

		client := testsuit.CreateTestClient()
		testsuit.CreateSchemaPizzaForTenants(t, client)
		testsuit.CreateSchemaSoupForTenants(t, client)
		testsuit.CreateTenantsPizza(t, client, tenants...)
		testsuit.CreateTenantsSoup(t, client, tenants...)

		for _, tenant := range tenants {
			resp, err := client.Batch().ObjectsBatcher().
				WithObjects(
					&models.Object{
						Class: "Pizza",
						ID:    testsuit.PIZZA_QUATTRO_FORMAGGI_ID,
						Properties: map[string]interface{}{
							"name":        "Quattro Formaggi",
							"description": "Pizza quattro formaggi Italian: [ˈkwattro forˈmaddʒi] (four cheese pizza) is a variety of pizza in Italian cuisine that is topped with a combination of four kinds of cheese, usually melted together, with (rossa, red) or without (bianca, white) tomato sauce. It is popular worldwide, including in Italy,[1] and is one of the iconic items from pizzerias's menus.",
							"price":       float32(1.1),
							"best_before": "2022-05-03T12:04:40+02:00",
						},
						Tenant: tenant,
					},
					&models.Object{
						Class: "Pizza",
						ID:    testsuit.PIZZA_FRUTTI_DI_MARE_ID,
						Properties: map[string]interface{}{
							"name":        "Frutti di Mare",
							"description": "Frutti di Mare is an Italian type of pizza that may be served with scampi, mussels or squid. It typically lacks cheese, with the seafood being served atop a tomato sauce.",
							"price":       float32(1.2),
							"best_before": "2022-05-05T07:16:30+02:00",
						},
						Tenant: tenant,
					},
					&models.Object{
						Class: "Soup",
						ID:    testsuit.SOUP_CHICKENSOUP_ID,
						Properties: map[string]interface{}{
							"name":        "ChickenSoup",
							"description": "Used by humans when their inferior genetics are attacked by microscopic organisms.",
							"price":       float32(2.1),
						},
						Tenant: tenant,
					},
					&models.Object{
						Class: "Soup",
						ID:    testsuit.SOUP_BEAUTIFUL_ID,
						Properties: map[string]interface{}{
							"name":        "Beautiful",
							"description": "Putting the game of letter soups to a whole new level.",
							"price":       float32(2.2),
						},
						Tenant: tenant,
					}).
				Do(context.Background())

			require.Nil(t, err)
			require.NotNil(t, resp)
			require.Len(t, resp, 4)

			ids := make([]string, len(resp))
			for i := range resp {
				require.NotNil(t, resp[i])
				require.NotNil(t, resp[i].Result)
				// require.NotNil(t, resp[i].Result.Status)
				// assert.Equal(t, "SUCCESS", *resp[i].Result.Status)
				assert.Equal(t, tenant, resp[i].Tenant)

				ids[i] = resp[i].ID.String()
			}
			assert.ElementsMatch(t, ids, []string{
				testsuit.PIZZA_QUATTRO_FORMAGGI_ID,
				testsuit.PIZZA_FRUTTI_DI_MARE_ID,
				testsuit.SOUP_CHICKENSOUP_ID,
				testsuit.SOUP_BEAUTIFUL_ID,
			})
		}

		t.Run("verify created", func(t *testing.T) {
			for _, tenant := range tenants {
				for id, className := range map[string]string{
					testsuit.PIZZA_QUATTRO_FORMAGGI_ID: "Pizza",
					testsuit.PIZZA_FRUTTI_DI_MARE_ID:   "Pizza",
					testsuit.SOUP_CHICKENSOUP_ID:       "Soup",
					testsuit.SOUP_BEAUTIFUL_ID:         "Soup",
				} {
					exists, err := client.Data().Checker().
						WithID(id).
						WithClassName(className).
						WithTenant(tenant).
						Do(context.Background())

					require.Nil(t, err)
					require.True(t, exists)
				}
			}
		})
	})

	t.Run("fails creating objects of MT class without tenant", func(t *testing.T) {
		t.Skip("to be enabled")
		defer cleanup()

		tenants := []string{"tenantNo1", "tenantNo2"}

		client := testsuit.CreateTestClient()
		testsuit.CreateSchemaPizzaForTenants(t, client)
		testsuit.CreateSchemaSoupForTenants(t, client)
		testsuit.CreateTenantsPizza(t, client, tenants...)
		testsuit.CreateTenantsSoup(t, client, tenants...)

		resp, err := client.Batch().ObjectsBatcher().
			WithObjects(
				&models.Object{
					Class: "Pizza",
					ID:    testsuit.PIZZA_QUATTRO_FORMAGGI_ID,
					Properties: map[string]interface{}{
						"name":        "Quattro Formaggi",
						"description": "Pizza quattro formaggi Italian: [ˈkwattro forˈmaddʒi] (four cheese pizza) is a variety of pizza in Italian cuisine that is topped with a combination of four kinds of cheese, usually melted together, with (rossa, red) or without (bianca, white) tomato sauce. It is popular worldwide, including in Italy,[1] and is one of the iconic items from pizzerias's menus.",
						"price":       float32(1.1),
						"best_before": "2022-05-03T12:04:40+02:00",
					},
				},
				&models.Object{
					Class: "Pizza",
					ID:    testsuit.PIZZA_FRUTTI_DI_MARE_ID,
					Properties: map[string]interface{}{
						"name":        "Frutti di Mare",
						"description": "Frutti di Mare is an Italian type of pizza that may be served with scampi, mussels or squid. It typically lacks cheese, with the seafood being served atop a tomato sauce.",
						"price":       float32(1.2),
						"best_before": "2022-05-05T07:16:30+02:00",
					},
				},
				&models.Object{
					Class: "Soup",
					ID:    testsuit.SOUP_CHICKENSOUP_ID,
					Properties: map[string]interface{}{
						"name":        "ChickenSoup",
						"description": "Used by humans when their inferior genetics are attacked by microscopic organisms.",
						"price":       float32(2.1),
					},
				},
				&models.Object{
					Class: "Soup",
					ID:    testsuit.SOUP_BEAUTIFUL_ID,
					Properties: map[string]interface{}{
						"name":        "Beautiful",
						"description": "Putting the game of letter soups to a whole new level.",
						"price":       float32(2.2),
					},
				}).
			Do(context.Background())

		require.Nil(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp, 4)

		for i := range resp {
			require.NotNil(t, resp[i])
			require.NotNil(t, resp[i].Result)
			// require.NotNil(t, resp[i].Result.Status)
			// assert.Equal(t, "FAILED", *resp[i].Result.Status)
			require.NotNil(t, resp[i].Result.Errors)
			require.NotNil(t, resp[i].Result.Errors.Error)
			require.Len(t, resp[i].Result.Errors.Error, 1)
			assert.Contains(t, resp[i].Result.Errors.Error[0].Message, "has multi-tenancy enabled, but request was without tenant")
			assert.Empty(t, resp[i].Tenant)
		}

		t.Run("verify not created", func(t *testing.T) {
			for _, tenant := range tenants {
				for id, className := range map[string]string{
					testsuit.PIZZA_QUATTRO_FORMAGGI_ID: "Pizza",
					testsuit.PIZZA_FRUTTI_DI_MARE_ID:   "Pizza",
					testsuit.SOUP_CHICKENSOUP_ID:       "Soup",
					testsuit.SOUP_BEAUTIFUL_ID:         "Soup",
				} {
					exists, err := client.Data().Checker().
						WithID(id).
						WithClassName(className).
						WithTenant(tenant).
						Do(context.Background())

					require.Nil(t, err)
					require.False(t, exists)
				}
			}
		})
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to tear down weaviate: %s", err)
		}
	})
}

func TestBatchReferenceCreate_MultiTenancy(t *testing.T) {
	cleanup := func() {
		client := testsuit.CreateTestClient()
		err := client.Schema().AllDeleter().Do(context.Background())
		require.Nil(t, err)
	}

	t.Run("setup weaviate", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to setup weaviate: %s", err)
		}
	})

	t.Run("creates references between MT classes", func(t *testing.T) {
		defer cleanup()

		tenants := []string{"tenantNo1", "tenantNo2"}
		soupIds := testsuit.IdsByClass["Soup"]
		pizzaIds := testsuit.IdsByClass["Pizza"]

		client := testsuit.CreateTestClient()
		testsuit.CreateSchemaSoupForTenants(t, client)
		testsuit.CreateTenantsSoup(t, client, tenants...)
		testsuit.CreateDataSoupForTenants(t, client, tenants...)

		testsuit.CreateSchemaPizzaForTenants(t, client)
		testsuit.CreateTenantsPizza(t, client, tenants...)
		testsuit.CreateDataPizzaForTenants(t, client, tenants...)

		t.Run("create ref property", func(t *testing.T) {
			err := client.Schema().PropertyCreator().
				WithClassName("Soup").
				WithProperty(&models.Property{
					Name:     "relatedToPizza",
					DataType: []string{"Pizza"},
				}).
				Do(context.Background())

			require.Nil(t, err)
		})

		references := []*models.BatchReference{}
		for _, tenant := range tenants {
			for _, soupId := range soupIds {
				for _, pizzaId := range pizzaIds {
					rpb := client.Batch().ReferencePayloadBuilder().
						WithFromClassName("Soup").
						WithFromRefProp("relatedToPizza").
						WithFromID(soupId).
						WithToClassName("Pizza").
						WithToID(pizzaId).
						WithTenant(tenant)

					references = append(references, rpb.Payload())
				}
			}
		}

		resp, err := client.Batch().ReferencesBatcher().
			WithReferences(references...).
			Do(context.Background())

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Len(t, resp, len(references))
		for i := range resp {
			require.NotNil(t, resp[i].Result)
			require.NotNil(t, resp[i].Result.Status)
			assert.Equal(t, "SUCCESS", *resp[i].Result.Status)
			assert.Nil(t, resp[i].Result.Errors)
		}

		t.Run("verify created", func(t *testing.T) {
			for _, tenant := range tenants {
				for _, soupId := range soupIds {
					objects, err := client.Data().ObjectsGetter().
						WithClassName("Soup").
						WithID(soupId).
						WithTenant(tenant).
						Do(context.Background())

					require.Nil(t, err)
					require.NotNil(t, objects)
					require.Len(t, objects, 1)
					assert.Len(t, objects[0].Properties.(map[string]interface{})["relatedToPizza"].([]interface{}),
						len(pizzaIds))
				}
			}
		})
	})

	t.Run("fails creating references between MT classes without tenant", func(t *testing.T) {
		defer cleanup()

		tenants := []string{"tenantNo1", "tenantNo2"}
		soupIds := testsuit.IdsByClass["Soup"]
		pizzaIds := testsuit.IdsByClass["Pizza"]

		client := testsuit.CreateTestClient()
		testsuit.CreateSchemaSoupForTenants(t, client)
		testsuit.CreateTenantsSoup(t, client, tenants...)
		testsuit.CreateDataSoupForTenants(t, client, tenants...)

		testsuit.CreateSchemaPizzaForTenants(t, client)
		testsuit.CreateTenantsPizza(t, client, tenants...)
		testsuit.CreateDataPizzaForTenants(t, client, tenants...)

		t.Run("create ref property", func(t *testing.T) {
			err := client.Schema().PropertyCreator().
				WithClassName("Soup").
				WithProperty(&models.Property{
					Name:     "relatedToPizza",
					DataType: []string{"Pizza"},
				}).
				Do(context.Background())

			require.Nil(t, err)
		})

		references := []*models.BatchReference{}
		for _, soupId := range soupIds {
			for _, pizzaId := range pizzaIds {
				rpb := client.Batch().ReferencePayloadBuilder().
					WithFromClassName("Soup").
					WithFromRefProp("relatedToPizza").
					WithFromID(soupId).
					WithToClassName("Pizza").
					WithToID(pizzaId)

				references = append(references, rpb.Payload())
			}
		}

		resp, err := client.Batch().ReferencesBatcher().
			WithReferences(references...).
			Do(context.Background())

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Len(t, resp, len(references))
		for i := range resp {
			require.NotNil(t, resp[i].Result)
			require.NotNil(t, resp[i].Result.Status)
			assert.Equal(t, "FAILED", *resp[i].Result.Status)
			require.NotNil(t, resp[i].Result.Errors)
			require.Len(t, resp[i].Result.Errors.Error, 1)
			assert.Contains(t, resp[i].Result.Errors.Error[0].Message, "has multi-tenancy enabled, but request was without tenant")
		}

		t.Run("verify not created", func(t *testing.T) {
			for _, tenant := range tenants {
				for _, soupId := range soupIds {
					objects, err := client.Data().ObjectsGetter().
						WithClassName("Soup").
						WithID(soupId).
						WithTenant(tenant).
						Do(context.Background())

					require.Nil(t, err)
					require.NotNil(t, objects)
					require.Len(t, objects, 1)
					assert.Nil(t, objects[0].Properties.(map[string]interface{})["relatedToPizza"])
				}
			}
		})
	})

	t.Run("creates references between MT and non-MT classes", func(t *testing.T) {
		defer cleanup()

		tenantSoup := "tenantSoup"
		pizzaIds := testsuit.IdsByClass["Pizza"]
		soupIds := testsuit.IdsByClass["Soup"]
		client := testsuit.CreateTestClient()

		testsuit.CreateSchemaSoupForTenants(t, client)
		testsuit.CreateTenantsSoup(t, client, tenantSoup)
		testsuit.CreateDataSoupForTenants(t, client, tenantSoup)

		testsuit.CreateSchemaPizza(t, client)
		testsuit.CreateDataPizza(t, client)

		t.Run("create ref property", func(t *testing.T) {
			err := client.Schema().PropertyCreator().
				WithClassName("Soup").
				WithProperty(&models.Property{
					Name:     "relatedToPizza",
					DataType: []string{"Pizza"},
				}).
				Do(context.Background())

			require.Nil(t, err)
		})

		references := []*models.BatchReference{}
		for _, soupId := range soupIds {
			for _, pizzaId := range pizzaIds {
				rpb := client.Batch().ReferencePayloadBuilder().
					WithFromClassName("Soup").
					WithFromRefProp("relatedToPizza").
					WithFromID(soupId).
					WithToClassName("Pizza").
					WithToID(pizzaId).
					WithTenant(tenantSoup)

				references = append(references, rpb.Payload())
			}
		}

		resp, err := client.Batch().ReferencesBatcher().
			WithReferences(references...).
			Do(context.Background())

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Len(t, resp, len(references))
		for i := range resp {
			require.NotNil(t, resp[i].Result)
			require.NotNil(t, resp[i].Result.Status)
			assert.Equal(t, "SUCCESS", *resp[i].Result.Status)
			assert.Nil(t, resp[i].Result.Errors)
		}

		t.Run("verify created", func(t *testing.T) {
			for _, soupId := range soupIds {
				objects, err := client.Data().ObjectsGetter().
					WithClassName("Soup").
					WithID(soupId).
					WithTenant(tenantSoup).
					Do(context.Background())

				require.Nil(t, err)
				require.NotNil(t, objects)
				require.Len(t, objects, 1)
				assert.Len(t, objects[0].Properties.(map[string]interface{})["relatedToPizza"].([]interface{}),
					len(pizzaIds))
			}
		})
	})

	t.Run("fails creating references between MT and non-MT classes without tenant", func(t *testing.T) {
		defer cleanup()

		tenantSoup := "tenantSoup"
		pizzaIds := testsuit.IdsByClass["Pizza"]
		soupIds := testsuit.IdsByClass["Soup"]
		client := testsuit.CreateTestClient()

		testsuit.CreateSchemaSoupForTenants(t, client)
		testsuit.CreateTenantsSoup(t, client, tenantSoup)
		testsuit.CreateDataSoupForTenants(t, client, tenantSoup)

		testsuit.CreateSchemaPizza(t, client)
		testsuit.CreateDataPizza(t, client)

		t.Run("create ref property", func(t *testing.T) {
			err := client.Schema().PropertyCreator().
				WithClassName("Soup").
				WithProperty(&models.Property{
					Name:     "relatedToPizza",
					DataType: []string{"Pizza"},
				}).
				Do(context.Background())

			require.Nil(t, err)
		})

		references := []*models.BatchReference{}
		for _, soupId := range soupIds {
			for _, pizzaId := range pizzaIds {
				rpb := client.Batch().ReferencePayloadBuilder().
					WithFromClassName("Soup").
					WithFromRefProp("relatedToPizza").
					WithFromID(soupId).
					WithToClassName("Pizza").
					WithToID(pizzaId)

				references = append(references, rpb.Payload())
			}
		}

		resp, err := client.Batch().ReferencesBatcher().
			WithReferences(references...).
			Do(context.Background())

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Len(t, resp, len(references))
		for i := range resp {
			require.NotNil(t, resp[i].Result)
			require.NotNil(t, resp[i].Result.Status)
			assert.Equal(t, "FAILED", *resp[i].Result.Status)
			require.NotNil(t, resp[i].Result.Errors)
			require.Len(t, resp[i].Result.Errors.Error, 1)
			assert.Contains(t, resp[i].Result.Errors.Error[0].Message, "has multi-tenancy enabled, but request was without tenant")
		}

		t.Run("verify not created", func(t *testing.T) {
			for _, soupId := range soupIds {
				objects, err := client.Data().ObjectsGetter().
					WithClassName("Soup").
					WithID(soupId).
					WithTenant(tenantSoup).
					Do(context.Background())

				require.Nil(t, err)
				require.NotNil(t, objects)
				require.Len(t, objects, 1)
				assert.Nil(t, objects[0].Properties.(map[string]interface{})["relatedToPizza"])
			}
		})

		t.Run("verify new objects not created", func(t *testing.T) {
			for _, soupId := range soupIds {
				exists, err := client.Data().Checker().
					WithClassName("Soup").
					WithID(soupId).
					Do(context.Background())

				require.NotNil(t, err)
				assert.False(t, exists)
			}

			for _, pizzaId := range pizzaIds {
				exists, err := client.Data().Checker().
					WithClassName("Pizza").
					WithID(pizzaId).
					WithTenant(tenantSoup).
					Do(context.Background())

				require.NotNil(t, err)
				assert.False(t, exists)
			}
		})
	})

	t.Run("fails creating references between non-MT and MT classes", func(t *testing.T) {
		defer cleanup()

		tenantPizza := "tenantPizza"
		pizzaIds := testsuit.IdsByClass["Pizza"]
		soupIds := testsuit.IdsByClass["Soup"]
		client := testsuit.CreateTestClient()

		testsuit.CreateSchemaSoup(t, client)
		testsuit.CreateDataSoup(t, client)

		testsuit.CreateSchemaPizzaForTenants(t, client)
		testsuit.CreateTenantsPizza(t, client, tenantPizza)
		testsuit.CreateDataPizzaForTenants(t, client, tenantPizza)

		t.Run("create ref property", func(t *testing.T) {
			err := client.Schema().PropertyCreator().
				WithClassName("Soup").
				WithProperty(&models.Property{
					Name:     "relatedToPizza",
					DataType: []string{"Pizza"},
				}).
				Do(context.Background())

			require.Nil(t, err)
		})

		references := []*models.BatchReference{}
		for _, soupId := range soupIds {
			for _, pizzaId := range pizzaIds {
				rpb := client.Batch().ReferencePayloadBuilder().
					WithFromClassName("Soup").
					WithFromRefProp("relatedToPizza").
					WithFromID(soupId).
					WithToClassName("Pizza").
					WithToID(pizzaId).
					WithTenant(tenantPizza)

				references = append(references, rpb.Payload())
			}
		}

		resp, err := client.Batch().ReferencesBatcher().
			WithReferences(references...).
			Do(context.Background())

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Len(t, resp, len(references))
		for i := range resp {
			require.NotNil(t, resp[i].Result)
			require.NotNil(t, resp[i].Result.Status)
			assert.Equal(t, "FAILED", *resp[i].Result.Status)
			require.NotNil(t, resp[i].Result.Errors)
			require.Len(t, resp[i].Result.Errors.Error, 1)
			assert.Contains(t, resp[i].Result.Errors.Error[0].Message, "cannot reference a multi-tenant enabled class from a non multi-tenant enabled class")
		}

		t.Run("verify not created", func(t *testing.T) {
			for _, soupId := range soupIds {
				objects, err := client.Data().ObjectsGetter().
					WithClassName("Soup").
					WithID(soupId).
					Do(context.Background())

				require.Nil(t, err)
				require.NotNil(t, objects)
				require.Len(t, objects, 1)
				assert.Nil(t, objects[0].Properties.(map[string]interface{})["relatedToPizza"])
			}
		})

		t.Run("verify new objects not created", func(t *testing.T) {
			for _, soupId := range soupIds {
				exists, err := client.Data().Checker().
					WithClassName("Soup").
					WithID(soupId).
					WithTenant(tenantPizza).
					Do(context.Background())

				require.NotNil(t, err)
				assert.False(t, exists)
			}
		})
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to tear down weaviate: %s", err)
		}
	})
}
