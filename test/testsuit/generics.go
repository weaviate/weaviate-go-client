package testsuit

import (
	"context"
	"github.com/go-openapi/strfmt"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/paragons"
	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

// CreateWeaviateTestSchemaFood creates a class for each semantic type (Pizza and Soup)
// and adds some primitive properties (name and description)
func CreateWeaviateTestSchemaFood(t *testing.T, client *weaviateclient.WeaviateClient) {
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
	errA := client.Schema.ClassCreator().WithClass(schemaClassAction).WithKind(paragons.SemanticKindActions).Do(context.Background())
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
	propErrA1 := client.Schema.PropertyCreator().WithClassName("Soup").WithProperty(nameProperty).WithKind(paragons.SemanticKindActions).Do(context.Background())
	assert.Nil(t, propErrA1)
	propErrT2 := client.Schema.PropertyCreator().WithClassName("Pizza").WithProperty(descriptionProperty).Do(context.Background())
	assert.Nil(t, propErrT2)
	propErrA2 := client.Schema.PropertyCreator().WithClassName("Soup").WithProperty(descriptionProperty).WithKind(paragons.SemanticKindActions).Do(context.Background())
	assert.Nil(t, propErrA2)
}

// CreateWeaviateTestSchemaFoodWithReferenceProperty create the testing schema with a reference field otherFoods on both classes
func CreateWeaviateTestSchemaFoodWithReferenceProperty(t *testing.T, client *weaviateclient.WeaviateClient) {
	CreateWeaviateTestSchemaFood(t, client)
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

// CleanUpWeaviate removes the schema and thereby all data
func CleanUpWeaviate(t *testing.T, client *weaviateclient.WeaviateClient) {
	// Clean up all classes and by that also all data
	errRm := client.Schema.AllDeleter().Do(context.Background())
	assert.Nil(t, errRm)
}

// CreateTestClient running on local host 8080
func CreateTestClient() *weaviateclient.WeaviateClient {
	cfg := weaviateclient.Config{
		Host:   "localhost:8080",
		Scheme: "http",
	}
	client := weaviateclient.New(cfg)
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



