package schema

import (
	"context"
	"fmt"
	"testing"

	"github.com/semi-technologies/weaviate-go-client/test/testsuit"
	"github.com/semi-technologies/weaviate-go-client/weaviate/testenv"
	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/stretchr/testify/assert"
)

func TestSchema_integration(t *testing.T) {

	t.Run("up", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			fmt.Printf(err.Error())
			t.Fail()
		}
	})

	t.Run("POST /schema", func(t *testing.T) {
		client := testsuit.CreateTestClient()

		schemaClass := &models.Class{
			Class:           "Band",
			Description:     "Band that plays and produces music",
			Properties:      nil,
			VectorIndexType: "hnsw",
			Vectorizer:      "text2vec-contextionary",
		}

		err := client.Schema().ClassCreator().WithClass(schemaClass).Do(context.Background())
		assert.Nil(t, err)

		loadedSchema, getErr := client.Schema().Getter().Do(context.Background())
		assert.Nil(t, getErr)
		assert.Equal(t, 1, len(loadedSchema.Classes))
		assert.Equal(t, schemaClass, loadedSchema.Classes[0])
		assert.Equal(t, schemaClass.Class, loadedSchema.Classes[0].Class)
		assert.Equal(t, schemaClass.Description, loadedSchema.Classes[0].Description)

		// Clean up classes
		errRm := client.Schema().AllDeleter().Do(context.Background())
		assert.Nil(t, errRm)
	})

	t.Run("POST /schema", func(t *testing.T) {
		client := testsuit.CreateTestClient()

		schemaClass := &models.Class{
			Class:           "Run",
			Description:     "Running from the fuzz",
			VectorIndexType: "hnsw",
			Vectorizer:      "text2vec-contextionary",
		}

		err := client.Schema().ClassCreator().WithClass(schemaClass).Do(context.Background())
		assert.Nil(t, err)

		loadedSchema, getErr := client.Schema().Getter().Do(context.Background())
		assert.Nil(t, getErr)
		assert.Equal(t, 1, len(loadedSchema.Classes))
		assert.Equal(t, schemaClass, loadedSchema.Classes[0])

		// Clean up classes
		errRm := client.Schema().AllDeleter().Do(context.Background())
		assert.Nil(t, errRm)
	})

	t.Run("Delete /schema/{type}", func(t *testing.T) {
		client := testsuit.CreateTestClient()

		schemaClassThing := &models.Class{
			Class:       "Pizza",
			Description: "A delicious religion like food and arguably the best export of Italy.",
		}
		schemaClassAction := &models.Class{
			Class:       "ChickenSoup",
			Description: "A soup made in part out of chicken, not for chicken.",
		}

		errT := client.Schema().ClassCreator().WithClass(schemaClassThing).Do(context.Background())
		assert.Nil(t, errT)
		errA := client.Schema().ClassCreator().WithClass(schemaClassAction).Do(context.Background())
		assert.Nil(t, errA)
		loadedSchema, getErr := client.Schema().Getter().Do(context.Background())
		assert.Nil(t, getErr)
		assert.Equal(t, loadedSchema.Classes[0].Class, schemaClassThing.Class)
		assert.Equal(t, loadedSchema.Classes[1].Class, schemaClassAction.Class)
		assert.Equal(t, 2, len(loadedSchema.Classes), "There are classes in the schema that are not part of this test")

		errRm1 := client.Schema().ClassDeleter().WithClassName(schemaClassThing.Class).Do(context.Background())
		errRm2 := client.Schema().ClassDeleter().WithClassName(schemaClassAction.Class).Do(context.Background())
		assert.Nil(t, errRm1)
		assert.Nil(t, errRm2)

		loadedSchema, getErr = client.Schema().Getter().Do(context.Background())
		assert.Nil(t, getErr)
		assert.Equal(t, 0, len(loadedSchema.Classes))
	})

	t.Run("Delete All schema", func(t *testing.T) {
		client := testsuit.CreateTestClient()

		schemaClassThing := &models.Class{
			Class:       "Pizza",
			Description: "A delicious religion like food and arguably the best export of Italy.",
		}
		schemaClassAction := &models.Class{
			Class:       "ChickenSoup",
			Description: "A soup made in part out of chicken, not for chicken.",
		}

		errT := client.Schema().ClassCreator().WithClass(schemaClassThing).Do(context.Background())
		assert.Nil(t, errT)
		errA := client.Schema().ClassCreator().WithClass(schemaClassAction).Do(context.Background())
		assert.Nil(t, errA)

		errRm1 := client.Schema().AllDeleter().Do(context.Background())
		assert.Nil(t, errRm1)

		loadedSchema, getErr := client.Schema().Getter().Do(context.Background())
		assert.Nil(t, getErr)
		assert.Equal(t, 0, len(loadedSchema.Classes))
	})

	t.Run("POST /schema/{type}/{className}/properties", func(t *testing.T) {
		client := testsuit.CreateTestClient()

		schemaClassThing := &models.Class{
			Class:       "Pizza",
			Description: "A delicious religion like food and arguably the best export of Italy.",
		}
		schemaClassAction := &models.Class{
			Class:       "ChickenSoup",
			Description: "A soup made in part out of chicken, not for chicken.",
		}

		errT := client.Schema().ClassCreator().WithClass(schemaClassThing).Do(context.Background())
		assert.Nil(t, errT)
		errA := client.Schema().ClassCreator().WithClass(schemaClassAction).Do(context.Background())
		assert.Nil(t, errA)

		newProperty := &models.Property{
			DataType:    []string{"string"},
			Description: "name",
			Name:        "name",
		}

		propErrT := client.Schema().PropertyCreator().WithClassName("Pizza").WithProperty(newProperty).Do(context.Background())
		assert.Nil(t, propErrT)
		propErrA := client.Schema().PropertyCreator().WithClassName("ChickenSoup").WithProperty(newProperty).Do(context.Background())
		assert.Nil(t, propErrA)

		loadedSchema, getErr := client.Schema().Getter().Do(context.Background())
		assert.Nil(t, getErr)
		assert.Equal(t, 2, len(loadedSchema.Classes))
		assert.Equal(t, "name", loadedSchema.Classes[0].Properties[0].Name)
		assert.Equal(t, "name", loadedSchema.Classes[1].Properties[0].Name)

		// Clean up classes
		errRm := client.Schema().AllDeleter().Do(context.Background())
		assert.Nil(t, errRm)
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			fmt.Printf(err.Error())
			t.Fail()
		}
	})
}

func TestSchema_errors(t *testing.T) {

	t.Run("Run Do withouth setting a class", func(t *testing.T) {
		client := testsuit.CreateTestClient()

		err := client.Schema().ClassCreator().Do(context.Background())
		assert.NotNil(t, err)
	})

}
