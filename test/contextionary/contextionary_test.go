package contextionary

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/testenv"
)

func TestContextionary_integration(t *testing.T) {
	if err := testenv.SetupLocalWeaviate(); err != nil {
		t.Fatalf("failed to setup weaviate: %s", err)
	}
	defer func() {
		if err := testenv.TearDownLocalWeaviate(); err != nil {
			t.Fatalf("failed to tear down weaviate: %s", err)
		}
	}()

	t.Run("GET /modules/text2vec-contextionary/concepts/{concept}", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)

		concepts, err := client.C11y().ConceptsGetter().WithConcept("pizzaHawaii").Do(context.Background())
		assert.Nil(t, err)
		if assert.NotNil(t, concepts) {
			assert.NotNil(t, concepts.ConcatenatedWord)
			assert.NotNil(t, concepts.IndividualWords)
		}
	})

	t.Run("POST /modules/text2vec-contextionary/extensions", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)

		err1 := client.C11y().ExtensionCreator().WithConcept("xoxo").WithDefinition("Hugs and kisses").WithWeight(1.0).Do(context.Background())
		assert.Nil(t, err1)

		err2 := client.C11y().ExtensionCreator().WithConcept("xoxo").WithDefinition("Hugs and kisses").WithWeight(2.0).Do(context.Background())
		assert.NotNil(t, err2, "Weight must be between 0 and 1")
	})
}
