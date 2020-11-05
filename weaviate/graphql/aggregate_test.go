package graphql

import (
	"github.com/semi-technologies/weaviate-go-client/weaviate/paragons"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAggregateBuilder(t *testing.T) {

	t.Run("Simple Explore", func(t *testing.T) {
		conMock := &MockRunREST{}

		builder := AggregateBuilder{
			connection:   conMock,
			semanticKind: paragons.SemanticKindThings,
		}

		query := builder.WithClassName("Pizza").WithFields("meta {count}").build()

		expected := `{Aggregate{Things{Pizza{meta {count}}}}}`
		assert.Equal(t, expected, query)
	})

	t.Run("Missuse", func(t *testing.T) {
		conMock := &MockRunREST{}
		builder := AggregateBuilder{
			connection:   conMock,
			semanticKind: paragons.SemanticKindThings,
		}
		query := builder.build()
		assert.NotEmpty(t, query, "Check that there is no panic if query is not validly build")
	})

}