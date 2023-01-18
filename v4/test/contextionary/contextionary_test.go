package contextionary

import (
	"context"
	"testing"

	"github.com/weaviate/weaviate-go-client/v4/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/testenv"
	"github.com/stretchr/testify/assert"
)

func TestContextionary_integration(t *testing.T) {
	t.Run("up", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatal(err.Error())
		}
	})

	t.Run("GET /modules/text2vec-contextionary/concepts/{concept}", func(t *testing.T) {
		client := testsuit.CreateTestClient()

		concepts, err := client.C11y().ConceptsGetter().WithConcept("pizzaHawaii").Do(context.Background())
		assert.Nil(t, err)
		if assert.NotNil(t, concepts) {
			assert.NotNil(t, concepts.ConcatenatedWord)
			assert.NotNil(t, concepts.IndividualWords)
		}
	})

	t.Run("POST /modules/text2vec-contextionary/extensions", func(t *testing.T) {
		client := testsuit.CreateTestClient()

		err1 := client.C11y().ExtensionCreator().WithConcept("xoxo").WithDefinition("Hugs and kisses").WithWeight(1.0).Do(context.Background())
		assert.Nil(t, err1)

		err2 := client.C11y().ExtensionCreator().WithConcept("xoxo").WithDefinition("Hugs and kisses").WithWeight(2.0).Do(context.Background())
		assert.NotNil(t, err2, "Weight must be between 0 and 1")
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			t.Fatal(err.Error())
		}
	})
}
