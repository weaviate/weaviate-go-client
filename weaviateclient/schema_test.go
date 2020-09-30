package weaviateclient

import (
	"context"
	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSchema_integration(t *testing.T) {
	// TODO up and down function

	t.Run("POST /schema/things", func(t *testing.T) {

		cfg := Config{
			Host:   "localhost:8080",
			Scheme: "http",
		}
		client := New(cfg)

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
	})

	t.Run("POST /schema/actions", func(t *testing.T) {

		cfg := Config{
			Host:   "localhost:8080",
			Scheme: "http",
		}
		client := New(cfg)

		schemaClass := &models.Class{
			Class:              "Run",
			Description:        "Running from the fuzz",
		}

		err := client.Schema.ClassCreator().WithClass(schemaClass).WithKind(SemanticKindActions).Do(context.Background())
		assert.Nil(t, err)

		loadedSchema, getErr := client.Schema.Getter().Do(context.Background())
		assert.Nil(t, getErr)
		assert.Equal(t, 1, len(loadedSchema.Actions.Classes))
		assert.Equal(t, schemaClass, loadedSchema.Actions.Classes[0])
	})



}

func TestSchema_unit(t *testing.T) {

	t.Run("Run Do withouth setting a class", func(t *testing.T) {
		// TODO NOT IMPLEMENTED
		t.Fail()

	})

}