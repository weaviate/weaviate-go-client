package testsuit

import (
	"context"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/semi-technologies/weaviate-go-client/v5/weaviate"
	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/stretchr/testify/assert"
)

// CreateWeaviateTestSchemaFood creates a class for each semantic type (Pizza and Soup)
// and adds some primitive properties (name and description)
func CreateWeaviateTestSchemaFood(t *testing.T, client *weaviate.Client) {
	schemaClassThing := &models.Class{
		Class:               "Pizza",
		Description:         "A delicious religion like food and arguably the best export of Italy.",
		InvertedIndexConfig: &models.InvertedIndexConfig{IndexTimestamps: true},
	}
	schemaClassAction := &models.Class{
		Class:       "Soup",
		Description: "Mostly water based brew of sustenance for humans.",
	}
	errT := client.Schema().ClassCreator().WithClass(schemaClassThing).Do(context.Background())
	assert.Nil(t, errT)
	errA := client.Schema().ClassCreator().WithClass(schemaClassAction).Do(context.Background())
	assert.Nil(t, errA)
	nameProperty := &models.Property{
		DataType:    []string{"string"},
		Description: "name",
		Name:        "name",
	}
	bestBeforeProperty := &models.Property{
		DataType:    []string{"date"},
		Description: "You better eat this food before it expires",
		Name:        "best_before",
	}
	descriptionProperty := &models.Property{
		DataType:    []string{"text"},
		Description: "description",
		Name:        "description",
	}
	priceProperty := &models.Property{
		DataType:    []string{"number"},
		Description: "price",
		Name:        "price",
		ModuleConfig: map[string]interface{}{
			"text2vec-contextionary": map[string]interface{}{
				"skip": true,
			},
		},
	}

	propErrT1 := client.Schema().PropertyCreator().WithClassName("Pizza").WithProperty(nameProperty).Do(context.Background())
	assert.Nil(t, propErrT1)
	propErrA1 := client.Schema().PropertyCreator().WithClassName("Soup").WithProperty(nameProperty).Do(context.Background())
	assert.Nil(t, propErrA1)
	propErrT2 := client.Schema().PropertyCreator().WithClassName("Pizza").WithProperty(descriptionProperty).Do(context.Background())
	assert.Nil(t, propErrT2)
	propErrA2 := client.Schema().PropertyCreator().WithClassName("Soup").WithProperty(descriptionProperty).Do(context.Background())
	assert.Nil(t, propErrA2)
	propErrT3 := client.Schema().PropertyCreator().WithClassName("Pizza").WithProperty(priceProperty).Do(context.Background())
	assert.Nil(t, propErrT3)
	propErrA3 := client.Schema().PropertyCreator().WithClassName("Soup").WithProperty(priceProperty).Do(context.Background())
	assert.Nil(t, propErrA3)
	propErrT4 := client.Schema().PropertyCreator().WithClassName("Pizza").WithProperty(bestBeforeProperty).Do(context.Background())
	assert.Nil(t, propErrT4)
	propErrA4 := client.Schema().PropertyCreator().WithClassName("Soup").WithProperty(bestBeforeProperty).Do(context.Background())
	assert.Nil(t, propErrA4)
}

func CreateWeaviateTestSchemaWithVectorizorlessClass(t *testing.T, client *weaviate.Client) {
	vectorizorlessClass := &models.Class{
		Class:       "Donut",
		Description: "A type of leavened fried dough commonly covered with glaze and sprinkles.",
		Vectorizer:  "none",
	}

	err := client.Schema().ClassCreator().WithClass(vectorizorlessClass).Do(context.Background())
	assert.Nil(t, err)

	nameProperty := &models.Property{
		DataType:    []string{"string"},
		Description: "name",
		Name:        "name",
	}
	descriptionProperty := &models.Property{
		DataType:    []string{"text"},
		Description: "description",
		Name:        "description",
	}

	propErr1 := client.Schema().PropertyCreator().WithClassName("Donut").WithProperty(nameProperty).Do(context.Background())
	assert.Nil(t, propErr1)

	propErr2 := client.Schema().PropertyCreator().WithClassName("Donut").WithProperty(descriptionProperty).Do(context.Background())
	assert.Nil(t, propErr2)
}

// CreateWeaviateTestSchemaFoodWithReferenceProperty create the testing schema with a reference field otherFoods on both classes
func CreateWeaviateTestSchemaFoodWithReferenceProperty(t *testing.T, client *weaviate.Client) {
	CreateWeaviateTestSchemaFood(t, client)
	referenceProperty := &models.Property{
		DataType:    []string{"Pizza", "Soup"},
		Description: "reference to other foods",
		Name:        "otherFoods",
	}
	err := client.Schema().PropertyCreator().WithClassName("Pizza").WithProperty(referenceProperty).Do(context.Background())
	assert.Nil(t, err)
	err = client.Schema().PropertyCreator().WithClassName("Soup").WithProperty(referenceProperty).Do(context.Background())
	assert.Nil(t, err)
}

// CleanUpWeaviate removes the schema and thereby all data
func CleanUpWeaviate(t *testing.T, client *weaviate.Client) {
	// Clean up all classes and by that also all data
	errRm := client.Schema().AllDeleter().Do(context.Background())
	assert.Nil(t, errRm)
}

// CreateTestClient running on local host 8080
func CreateTestClient() *weaviate.Client {
	cfg := weaviate.Config{
		Host:   "localhost:8080",
		Scheme: "http",
	}
	client := weaviate.New(cfg)
	return client
}

// ParseReferenceResponseToStruct from the interface typed property schema returned by the client
func ParseReferenceResponseToStruct(t *testing.T, reference interface{}) models.MultipleRef {
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

// CreateTestSchemaAndData with a few pizzas and soups
func CreateTestSchemaAndData(t *testing.T, client *weaviate.Client) {
	CreateWeaviateTestSchemaFood(t, client)

	// Create pizzas
	menuPizza := []*models.Object{
		{
			Class: "Pizza",
			Properties: map[string]interface{}{
				"name":        "Quattro Formaggi",
				"description": "Pizza quattro formaggi Italian: [ˈkwattro forˈmaddʒi] (four cheese pizza) is a variety of pizza in Italian cuisine that is topped with a combination of four kinds of cheese, usually melted together, with (rossa, red) or without (bianca, white) tomato sauce. It is popular worldwide, including in Italy,[1] and is one of the iconic items from pizzerias's menus.",
				"price":       float32(1.1),
				"best_before": "2022-05-03T12:04:40+02:00",
			},
		},
		{
			Class: "Pizza",
			Properties: map[string]interface{}{
				"name":        "Frutti di Mare",
				"description": "Frutti di Mare is an Italian type of pizza that may be served with scampi, mussels or squid. It typically lacks cheese, with the seafood being served atop a tomato sauce.",
				"price":       float32(1.2),
				"best_before": "2022-05-05T07:16:30+02:00",
			},
		},
		{
			Class: "Pizza",
			ID:    "5b6a08ba-1d46-43aa-89cc-8b070790c6f1",
			Properties: map[string]interface{}{
				"name":        "Hawaii",
				"description": "Universally accepted to be the best pizza ever created.",
				"price":       float32(1.3),
			},
		},
		{
			ID:    "5b6a08ba-1d46-43aa-89cc-8b070790c6f2",
			Class: "Pizza",
			Properties: map[string]interface{}{
				"name":        "Doener",
				"description": "A innovation, some say revolution, in the pizza industry.",
				"price":       float32(1.4),
			},
		},
	}
	menuSoup := []*models.Object{
		{
			Class: "Soup",
			Properties: map[string]interface{}{
				"name":        "ChickenSoup",
				"description": "Used by humans when their inferior genetics are attacked by microscopic organisms.",
				"price":       float32(2.1),
			},
		},
		{
			Class: "Soup",
			Properties: map[string]interface{}{
				"name":        "Beautiful",
				"description": "Putting the game of letter soups to a whole new level.",
				"price":       float32(2.2),
			},
		},
	}

	thingsBatcher := client.Batch().ObjectsBatcher()
	for _, pizza := range menuPizza {
		thingsBatcher.WithObject(pizza)
	}
	_, thingsErr := thingsBatcher.Do(context.Background())
	assert.Nil(t, thingsErr)

	_, actionsErr := client.Batch().ObjectsBatcher().WithObject(menuSoup[0]).WithObject(menuSoup[1]).Do(context.Background())
	assert.Nil(t, actionsErr)
}
