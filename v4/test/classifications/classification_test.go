package classifications

import (
	"context"
	"fmt"
	"testing"

	"github.com/semi-technologies/weaviate-go-client/v4/test/testsuit"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/classifications"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/testenv"
	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/semi-technologies/weaviate/usecases/classification"
	"github.com/stretchr/testify/assert"
)

func TestClassifications_integration(t *testing.T) {

	t.Run("up", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			fmt.Printf(err.Error())
			t.Fail()
		}
	})

	t.Run("POST /classifications", func(t *testing.T) {
		client := testsuit.CreateTestClient()
		createClassificationClasses(t, client)

		classifyProperties := []string{"taged"}
		basedOnProperties := []string{"description"}
		classification, err := client.Classifications().Scheduler().WithType(classifications.Contextual).WithClassName("Pizza").WithClassifyProperties(classifyProperties).WithBasedOnProperties(basedOnProperties).Do(context.Background())
		assert.Nil(t, err)
		assert.Contains(t, classification.BasedOnProperties, "description")
		assert.Contains(t, classification.ClassifyProperties, "taged")
		classification, err = client.Classifications().Scheduler().WithType(classifications.Contextual).WithClassName("Pizza").WithClassifyProperties(classifyProperties).WithBasedOnProperties(basedOnProperties).WithWaitForCompletion().Do(context.Background())
		assert.Nil(t, err)
		assert.Contains(t, classification.BasedOnProperties, "description")
		assert.Contains(t, classification.ClassifyProperties, "taged")

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("GET /classifications/{id}", func(t *testing.T) {
		client := testsuit.CreateTestClient()
		createClassificationClasses(t, client)

		var k int32 = 3
		c, err := client.Classifications().Scheduler().
			WithType(classifications.KNN).
			WithSettings(&classification.ParamsKNN{K: &k}).
			WithClassName("Pizza").
			WithClassifyProperties([]string{"taged"}).
			WithBasedOnProperties([]string{"description"}).
			Do(context.Background())
		assert.Nil(t, err)

		getC, getErr := client.Classifications().Getter().WithID(string(c.ID)).Do(context.Background())
		assert.Nil(t, getErr)
		assert.Equal(t, c.ID, getC.ID)
		knn, ok := getC.Settings.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, float64(3), knn["k"])

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

func createClassificationClasses(t *testing.T, client *weaviate.Client) {
	testsuit.CreateWeaviateTestSchemaFood(t, client)

	// Create a class Tag
	schemaClassTag := &models.Class{
		Class:       "Tag",
		Description: "tag for a pizza",
		Properties: []*models.Property{
			{
				DataType:    []string{"string"},
				Description: "name",
				Name:        "name",
			},
		},
	}
	errT := client.Schema().ClassCreator().WithClass(schemaClassTag).Do(context.Background())
	assert.Nil(t, errT)
	// Create a reference property that allows to tag a pizza
	tagProperty := models.Property{
		DataType:    []string{"Tag"},
		Description: "tag of pizza",
		Name:        "taged",
	}
	addTagPropertyToPizzaErr := client.Schema().PropertyCreator().WithProperty(&tagProperty).WithClassName("Pizza").Do(context.Background())
	assert.Nil(t, addTagPropertyToPizzaErr)

	// Create two pizzas
	pizza1 := &models.Object{
		Class: "Pizza",
		ID:    "97fa5147-bdad-4d74-9a81-f8babc811b09",
		Properties: map[string]string{
			"name":        "Quattro Formaggi",
			"description": "Pizza quattro formaggi Italian: [ˈkwattro forˈmaddʒi] (four cheese pizza) is a variety of pizza in Italian cuisine that is topped with a combination of four kinds of cheese, usually melted together, with (rossa, red) or without (bianca, white) tomato sauce. It is popular worldwide, including in Italy,[1] and is one of the iconic items from pizzerias's menus.",
		},
	}
	pizza2 := &models.Object{
		Class: "Pizza",
		ID:    "97fa5147-bdad-4d74-9a81-f8babc811b09",
		Properties: map[string]string{
			"name":        "Frutti di Mare",
			"description": "Frutti di Mare is an Italian type of pizza that may be served with scampi, mussels or squid. It typically lacks cheese, with the seafood being served atop a tomato sauce.",
		},
	}
	_, batchErr := client.Batch().ObjectsBatcher().WithObjects(pizza1, pizza2).Do(context.Background())
	assert.Nil(t, batchErr)
	// Create two tags
	tag1 := &models.Object{
		Class: "Tag",
		Properties: map[string]string{
			"name": "vegetarian",
		},
	}
	tag2 := &models.Object{
		Class: "Tag",
		Properties: map[string]string{
			"name": "seafood",
		},
	}
	_, batchErr2 := client.Batch().ObjectsBatcher().WithObjects(tag1, tag2).Do(context.Background())
	assert.Nil(t, batchErr2)
}
