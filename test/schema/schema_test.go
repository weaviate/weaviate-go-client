package schema

import (
	"context"
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient"
	clientModels "github.com/semi-technologies/weaviate-go-client/weaviateclient/models"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/testenv"
	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSchema_integration(t *testing.T) {

	t.Run("up", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			fmt.Printf(err.Error())
			t.Fail()
		}
	})

	t.Run("POST /schema/things", func(t *testing.T) {

		cfg := weaviateclient.Config{
			Host:   "localhost:8080",
			Scheme: "http",
		}
		client := weaviateclient.New(cfg)

		schemaClass := &models.Class{
			Class:              "Band",
			Description:        "Band that plays and produces music",
			Keywords:           nil,
			Properties:         nil,
			VectorizeClassName: nil,
		}

		err := client.Schema.ClassCreator().WithClass(schemaClass).Do(context.Background())
		assert.Nil(t, err)

		loadedSchema, getErr := client.Schema.Getter().Do(context.Background())
		assert.Nil(t, getErr)
		assert.Equal(t, 1, len(loadedSchema.Things.Classes))
		assert.Equal(t, schemaClass, loadedSchema.Things.Classes[0])
		assert.Equal(t, schemaClass.Class, loadedSchema.Things.Classes[0].Class)
		assert.Equal(t, schemaClass.Description, loadedSchema.Things.Classes[0].Description)

		// Clean up classes
		errRm := client.Schema.AllDeleter().Do(context.Background())
		assert.Nil(t, errRm)
	})

	t.Run("POST /schema/actions", func(t *testing.T) {

		cfg := weaviateclient.Config{
			Host:   "localhost:8080",
			Scheme: "http",
		}
		client := weaviateclient.New(cfg)

		schemaClass := &models.Class{
			Class:              "Run",
			Description:        "Running from the fuzz",
		}

		err := client.Schema.ClassCreator().WithClass(schemaClass).WithKind(clientModels.SemanticKindActions).Do(context.Background())
		assert.Nil(t, err)

		loadedSchema, getErr := client.Schema.Getter().Do(context.Background())
		assert.Nil(t, getErr)
		assert.Equal(t, 1, len(loadedSchema.Actions.Classes))
		assert.Equal(t, schemaClass, loadedSchema.Actions.Classes[0])

		// Clean up classes
		errRm := client.Schema.AllDeleter().Do(context.Background())
		assert.Nil(t, errRm)
	})

	t.Run("Delete /schema/{type}", func(t *testing.T) {

		cfg := weaviateclient.Config{
			Host:   "localhost:8080",
			Scheme: "http",
		}
		client := weaviateclient.New(cfg)

		schemaClassThing := &models.Class{
			Class:              "Pizza",
			Description:        "A delicious religion like food and arguably the best export of Italy.",
		}
		schemaClassAction := &models.Class{
			Class:              "ChickenSoup",
			Description:        "A soup made in part out of chicken, not for chicken.",
		}

		errT := client.Schema.ClassCreator().WithClass(schemaClassThing).Do(context.Background())
		assert.Nil(t, errT)
		errA := client.Schema.ClassCreator().WithClass(schemaClassAction).WithKind(clientModels.SemanticKindActions).Do(context.Background())
		assert.Nil(t, errA)
		loadedSchema, getErr := client.Schema.Getter().Do(context.Background())
		assert.Nil(t, getErr)
		assert.Equal(t, loadedSchema.Things.Classes[0].Class, schemaClassThing.Class)
		assert.Equal(t, loadedSchema.Actions.Classes[0].Class, schemaClassAction.Class)
		assert.Equal(t, 1, len(loadedSchema.Things.Classes), "There are classes in the schema that are not part of this test")
		assert.Equal(t, 1, len(loadedSchema.Actions.Classes), "There are classes in the schema that are not part of this test")

		errRm1 := client.Schema.ClassDeleter().WithClassName(schemaClassThing.Class).Do(context.Background())
		errRm2 := client.Schema.ClassDeleter().WithClassName(schemaClassAction.Class).WithKind(clientModels.SemanticKindActions).Do(context.Background())
		assert.Nil(t, errRm1)
		assert.Nil(t, errRm2)

		loadedSchema, getErr = client.Schema.Getter().Do(context.Background())
		assert.Nil(t, getErr)
		assert.Equal(t, 0, len(loadedSchema.Things.Classes))
		assert.Equal(t, 0, len(loadedSchema.Actions.Classes))
	})

	t.Run("Delete All schema", func(t *testing.T) {

		cfg := weaviateclient.Config{
			Host:   "localhost:8080",
			Scheme: "http",
		}
		client := weaviateclient.New(cfg)

		schemaClassThing := &models.Class{
			Class:              "Pizza",
			Description:        "A delicious religion like food and arguably the best export of Italy.",
		}
		schemaClassAction := &models.Class{
			Class:              "ChickenSoup",
			Description:        "A soup made in part out of chicken, not for chicken.",
		}

		errT := client.Schema.ClassCreator().WithClass(schemaClassThing).Do(context.Background())
		assert.Nil(t, errT)
		errA := client.Schema.ClassCreator().WithClass(schemaClassAction).WithKind(clientModels.SemanticKindActions).Do(context.Background())
		assert.Nil(t, errA)

		errRm1 := client.Schema.AllDeleter().Do(context.Background())
		assert.Nil(t, errRm1)

		loadedSchema, getErr := client.Schema.Getter().Do(context.Background())
		assert.Nil(t, getErr)
		assert.Equal(t, 0, len(loadedSchema.Things.Classes))
		assert.Equal(t, 0, len(loadedSchema.Actions.Classes))
	})

	t.Run("POST /schema/{type}/{className}/properties", func(t *testing.T) {

		cfg := weaviateclient.Config{
			Host:   "localhost:8080",
			Scheme: "http",
		}
		client := weaviateclient.New(cfg)

		schemaClassThing := &models.Class{
			Class:              "Pizza",
			Description:        "A delicious religion like food and arguably the best export of Italy.",
		}
		schemaClassAction := &models.Class{
			Class:              "ChickenSoup",
			Description:        "A soup made in part out of chicken, not for chicken.",
		}

		errT := client.Schema.ClassCreator().WithClass(schemaClassThing).Do(context.Background())
		assert.Nil(t, errT)
		errA := client.Schema.ClassCreator().WithClass(schemaClassAction).WithKind(clientModels.SemanticKindActions).Do(context.Background())
		assert.Nil(t, errA)

		newProperty := models.Property{
			DataType:              []string{"string"},
			Description:           "name",
			Name:                  "name",
		}

		propErrT := client.Schema.PropertyCreator().WithClassName("Pizza").WithProperty(newProperty).Do(context.Background())
		assert.Nil(t, propErrT)
		propErrA := client.Schema.PropertyCreator().WithClassName("ChickenSoup").WithProperty(newProperty).WithKind(clientModels.SemanticKindActions).Do(context.Background())
		assert.Nil(t, propErrA)

		loadedSchema, getErr := client.Schema.Getter().Do(context.Background())
		assert.Nil(t, getErr)
		assert.Equal(t, 1, len(loadedSchema.Things.Classes))
		assert.Equal(t, 1, len(loadedSchema.Actions.Classes))
		assert.Equal(t, "name", loadedSchema.Things.Classes[0].Properties[0].Name)
		assert.Equal(t, "name", loadedSchema.Actions.Classes[0].Properties[0].Name)

		// Clean up classes
		errRm := client.Schema.AllDeleter().Do(context.Background())
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

func TestSchema_unit(t *testing.T) {

	t.Run("Run Do withouth setting a class", func(t *testing.T) {
		cfg := weaviateclient.Config{
			Host:   "localhost:8080",
			Scheme: "http",
		}
		client := weaviateclient.New(cfg)

		err := client.Schema.ClassCreator().Do(context.Background())
		assert.NotNil(t, err)


	})

}