package testsuit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/schema"
)

const TenantKey = "tenantName"

// ##### SCHEMA #####

func CreateSchemaPizza(t *testing.T, client *weaviate.Client) {
	createSchema(t, client, classPizza())
}

func CreateSchemaSoup(t *testing.T, client *weaviate.Client) {
	createSchema(t, client, classSoup())
}

func CreateSchemaRisotto(t *testing.T, client *weaviate.Client) {
	createSchema(t, client, classRisotto())
}

func CreateSchemaFood(t *testing.T, client *weaviate.Client) {
	CreateSchemaPizza(t, client)
	CreateSchemaSoup(t, client)
	CreateSchemaRisotto(t, client)
}

func CreateSchemaPizzaForTenants(t *testing.T, client *weaviate.Client) {
	createSchema(t, client, classPizzaForTenants(TenantKey))
}

func CreateSchemaSoupForTenants(t *testing.T, client *weaviate.Client) {
	createSchema(t, client, classSoupForTenants(TenantKey))
}

func CreateSchemaRisottoForTenants(t *testing.T, client *weaviate.Client) {
	createSchema(t, client, classRisottoForTenants(TenantKey))
}

func CreateSchemaFoodForTenants(t *testing.T, client *weaviate.Client) {
	CreateSchemaPizzaForTenants(t, client)
	CreateSchemaSoupForTenants(t, client)
	CreateSchemaRisottoForTenants(t, client)
}

func createSchema(t *testing.T, client *weaviate.Client, class *models.Class) {
	err := client.Schema().ClassCreator().
		WithClass(class).
		Do(context.Background())

	require.Nil(t, err)
}

// ##### CLASSES #####

func classPizza() *models.Class {
	return &models.Class{
		Class:               "Pizza",
		Description:         "A delicious religion like food and arguably the best export of Italy.",
		InvertedIndexConfig: &models.InvertedIndexConfig{IndexTimestamps: true},
		Properties:          classPropertiesFood(),
	}
}

func classPizzaForTenants(tenantKey string) *models.Class {
	class := classPizza()
	class.Properties = classPropertiesFoodForTenants(tenantKey)
	class.MultiTenancyConfig = &models.MultiTenancyConfig{
		Enabled:   true,
		TenantKey: tenantKey,
	}
	return class
}

func classSoup() *models.Class {
	return &models.Class{
		Class:       "Soup",
		Description: "Mostly water based brew of sustenance for humans.",
		Properties:  classPropertiesFood(),
	}
}

func classSoupForTenants(tenantKey string) *models.Class {
	class := classSoup()
	class.Properties = classPropertiesFoodForTenants(tenantKey)
	class.MultiTenancyConfig = &models.MultiTenancyConfig{
		Enabled:   true,
		TenantKey: tenantKey,
	}
	return class
}

func classRisotto() *models.Class {
	return &models.Class{
		Class:               "Risotto",
		Description:         "Risotto is a northern Italian rice dish cooked with broth.",
		InvertedIndexConfig: &models.InvertedIndexConfig{IndexTimestamps: true},
		Properties:          classPropertiesFood(),
	}
}

func classRisottoForTenants(tenantKey string) *models.Class {
	class := classRisotto()
	class.Properties = classPropertiesFoodForTenants(tenantKey)
	class.MultiTenancyConfig = &models.MultiTenancyConfig{
		Enabled:   true,
		TenantKey: tenantKey,
	}
	return class
}

func classPropertiesFood() []*models.Property {
	nameProperty := &models.Property{
		Name:         "name",
		Description:  "name",
		DataType:     schema.DataTypeText.PropString(),
		Tokenization: models.PropertyTokenizationField,
	}
	descriptionProperty := &models.Property{
		Name:        "description",
		Description: "description",
		DataType:    schema.DataTypeText.PropString(),
	}
	bestBeforeProperty := &models.Property{
		Name:        "best_before",
		Description: "You better eat this food before it expires",
		DataType:    schema.DataTypeDate.PropString(),
	}
	priceProperty := &models.Property{
		Name:        "price",
		Description: "price",
		DataType:    schema.DataTypeNumber.PropString(),
		ModuleConfig: map[string]interface{}{
			"text2vec-contextionary": map[string]interface{}{
				"skip": true,
			},
		},
	}

	return []*models.Property{
		nameProperty, descriptionProperty, bestBeforeProperty, priceProperty,
	}
}

func classPropertiesFoodForTenants(tenantKey string) []*models.Property {
	return append(classPropertiesFood(), &models.Property{
		Name:        tenantKey,
		Description: "property used as tenant key",
		DataType:    schema.DataTypeText.PropString(),
	})
}

// ##### DATA #####

func CreateDataPizza(t *testing.T, client *weaviate.Client) {
	createData(t, client, []*models.Object{
		objectPizzaQuattroFormaggi(),
		objectPizzaFruttiDiMare(),
		objectPizzaHawaii(),
		objectPizzaDoener(),
	})
}

func CreateDataSoup(t *testing.T, client *weaviate.Client) {
	createData(t, client, []*models.Object{
		objectSoupChicken(),
		objectSoupBeautiful(),
	})
}

func CreateDataRisotto(t *testing.T, client *weaviate.Client) {
	createData(t, client, []*models.Object{
		objectRisottoRisiEBisi(),
		objectRisottoAllaPilota(),
		objectRisottoAlNeroDiSeppia(),
	})
}

func CreateDataFood(t *testing.T, client *weaviate.Client) {
	createData(t, client, []*models.Object{
		objectPizzaQuattroFormaggi(),
		objectPizzaFruttiDiMare(),
		objectPizzaHawaii(),
		objectPizzaDoener(),

		objectSoupChicken(),
		objectSoupBeautiful(),

		objectRisottoRisiEBisi(),
		objectRisottoAllaPilota(),
		objectRisottoAlNeroDiSeppia(),
	})
}

func createData(t *testing.T, client *weaviate.Client, objects []*models.Object) {
	resp, err := client.Batch().ObjectsBatcher().
		WithObjects(objects...).
		Do(context.Background())

	require.Nil(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp, len(objects))
}

func CreateDataPizzaForTenants(t *testing.T, client *weaviate.Client, tenantNames ...string) {
	createDataForTenants(t, client, TenantKey, tenantNames, func() []*models.Object {
		return []*models.Object{
			objectPizzaQuattroFormaggi(),
			objectPizzaFruttiDiMare(),
			objectPizzaHawaii(),
			objectPizzaDoener(),
		}
	})
}

func CreateDataSoupForTenants(t *testing.T, client *weaviate.Client, tenantNames ...string) {
	createDataForTenants(t, client, TenantKey, tenantNames, func() []*models.Object {
		return []*models.Object{
			objectSoupChicken(),
			objectSoupBeautiful(),
		}
	})
}

func CreateDataRisottoForTenants(t *testing.T, client *weaviate.Client, tenantNames ...string) {
	createDataForTenants(t, client, TenantKey, tenantNames, func() []*models.Object {
		return []*models.Object{
			objectRisottoRisiEBisi(),
			objectRisottoAllaPilota(),
			objectRisottoAlNeroDiSeppia(),
		}
	})
}

func CreateDataFoodForTenants(t *testing.T, client *weaviate.Client, tenantNames ...string) {
	createDataForTenants(t, client, TenantKey, tenantNames, func() []*models.Object {
		return []*models.Object{
			objectPizzaQuattroFormaggi(),
			objectPizzaFruttiDiMare(),
			objectPizzaHawaii(),
			objectPizzaDoener(),

			objectSoupChicken(),
			objectSoupBeautiful(),

			objectRisottoRisiEBisi(),
			objectRisottoAllaPilota(),
			objectRisottoAlNeroDiSeppia(),
		}
	})
}

func createDataForTenants(t *testing.T, client *weaviate.Client, tenantKey string,
	tenantNames []string, objectsSupplier func() []*models.Object,
) {
	for _, name := range tenantNames {
		objects := objectsSupplier()
		for _, object := range objects {
			props := object.Properties.(map[string]interface{})
			props[tenantKey] = name
		}

		resp, err := client.Batch().ObjectsBatcher().
			WithObjects(objects...).
			WithTenantKey(name).
			Do(context.Background())

		require.Nil(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp, len(objects))
	}
}

// ##### OBJECTS #####

func objectPizzaQuattroFormaggi() *models.Object {
	return &models.Object{
		Class: "Pizza",
		ID:    "10523cdd-15a2-42f4-81fa-267fe92f7cd6",
		Properties: map[string]interface{}{
			"name":        "Quattro Formaggi",
			"description": "Pizza quattro formaggi Italian: [ˈkwattro forˈmaddʒi] (four cheese pizza) is a variety of pizza in Italian cuisine that is topped with a combination of four kinds of cheese, usually melted together, with (rossa, red) or without (bianca, white) tomato sauce. It is popular worldwide, including in Italy,[1] and is one of the iconic items from pizzerias's menus.",
			"price":       float32(1.1),
			"best_before": "2022-05-03T12:04:40+02:00",
		},
	}
}

func objectPizzaFruttiDiMare() *models.Object {
	return &models.Object{
		Class: "Pizza",
		ID:    "927dd3ac-e012-4093-8007-7799cc7e81e4",
		Properties: map[string]interface{}{
			"name":        "Frutti di Mare",
			"description": "Frutti di Mare is an Italian type of pizza that may be served with scampi, mussels or squid. It typically lacks cheese, with the seafood being served atop a tomato sauce.",
			"price":       float32(1.2),
			"best_before": "2022-05-05T07:16:30+02:00",
		},
	}
}

func objectPizzaHawaii() *models.Object {
	return &models.Object{
		Class: "Pizza",
		// this uuid guarantees that it's the first for cursor tests (otherwise
		// they might be flaky if the randomly generated ids are sometimes higher
		// and sometimes lower
		ID: "00000000-0000-0000-0000-000000000000",
		Properties: map[string]interface{}{
			"name":        "Hawaii",
			"description": "Universally accepted to be the best pizza ever created.",
			"price":       float32(1.3),
		},
	}
}

func objectPizzaDoener() *models.Object {
	return &models.Object{
		Class: "Pizza",
		ID:    "5b6a08ba-1d46-43aa-89cc-8b070790c6f2",
		Properties: map[string]interface{}{
			"name":        "Doener",
			"description": "A innovation, some say revolution, in the pizza industry.",
			"price":       float32(1.4),
		},
	}
}

func objectSoupChicken() *models.Object {
	return &models.Object{
		Class: "Soup",
		ID:    "8c156d37-81aa-4ce9-a811-621e2702b825",
		Properties: map[string]interface{}{
			"name":        "ChickenSoup",
			"description": "Used by humans when their inferior genetics are attacked by microscopic organisms.",
			"price":       float32(2.1),
		},
	}
}

func objectSoupBeautiful() *models.Object {
	return &models.Object{
		Class: "Soup",
		ID:    "27351361-2898-4d1a-aad7-1ca48253eb0b",
		Properties: map[string]interface{}{
			"name":        "Beautiful",
			"description": "Putting the game of letter soups to a whole new level.",
			"price":       float32(2.2),
		},
	}
}

func objectRisottoRisiEBisi() *models.Object {
	return &models.Object{
		Class: "Risotto",
		ID:    "da751a25-f573-4715-a893-e607b2de0ba4",
		Properties: map[string]interface{}{
			"name":        "Risi e bisi",
			"description": "A Veneto spring dish.",
			"price":       float32(3.1),
			"best_before": "2022-05-03T12:04:40+02:00",
		},
	}
}

func objectRisottoAllaPilota() *models.Object {
	return &models.Object{
		Class: "Risotto",
		ID:    "10c2ee44-7d58-42be-9d64-5766883ca8cb",
		Properties: map[string]interface{}{
			"name":        "Risotto alla pilota",
			"description": "A specialty of Mantua, made with sausage, pork, and Parmesan cheese.",
			"price":       float32(3.2),
			"best_before": "2022-05-03T12:04:40+02:00",
		},
	}
}

func objectRisottoAlNeroDiSeppia() *models.Object {
	return &models.Object{
		Class: "Risotto",
		ID:    "696bf381-7f98-40a4-bcad-841780e00e0e",
		Properties: map[string]interface{}{
			"name":        "Risotto al nero di seppia",
			"description": "A specialty of the Veneto region, made with cuttlefish cooked with their ink-sacs intact, leaving the risotto black .",
			"price":       float32(3.3),
		},
	}
}

// ##### TENANTS #####

func CreateTenantsPizza(t *testing.T, client *weaviate.Client, tenantNames ...string) {
	createTenants(t, client, "Pizza", tenantNames)
}

func CreateTenantsSoup(t *testing.T, client *weaviate.Client, tenantNames ...string) {
	createTenants(t, client, "Soup", tenantNames)
}

func CreateTenantsRisotto(t *testing.T, client *weaviate.Client, tenantNames ...string) {
	createTenants(t, client, "Risotto", tenantNames)
}

func CreateTenantsFood(t *testing.T, client *weaviate.Client, tenantNames ...string) {
	CreateTenantsPizza(t, client, tenantNames...)
	CreateTenantsSoup(t, client, tenantNames...)
	CreateTenantsRisotto(t, client, tenantNames...)
}

func createTenants(t *testing.T, client *weaviate.Client, className string, tenantNames []string) {
	tenants := make([]models.Tenant, len(tenantNames))
	for i, name := range tenantNames {
		tenants[i] = models.Tenant{Name: name}
	}

	err := client.Schema().TenantCreator().
		WithClassName(className).
		WithTenants(tenants...).
		Do(context.Background())
	require.Nil(t, err)
}
