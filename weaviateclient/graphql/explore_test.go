package graphql

import (
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/paragons"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExploreBuilder(t *testing.T) {

	t.Run("Simple Explore", func(t *testing.T) {
		conMock := &MockRunREST{}

		builder := Explore{
			connection: conMock,
		}

		query := builder.WithFields([]paragons.ExploreFields{paragons.Certainty, paragons.Beacon}).WithConcepts([]string{"Cheese", "pineapple"}).build()

		expected := `{Explore(concepts: ["Cheese","pineapple"]){certainty beacon }}`
		assert.Equal(t, expected, query)
	})

}