package graphql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerativeSearch_build(t *testing.T) {
	t.Run("single result", func(t *testing.T) {
		gs := NewGSWithSingleResult("Describe this pizza : {name}")
		result := gs.build()
		require.Equal(t, `_additional{generate(singleResult:{ prompt: """Describe this pizza : {name}""" }){singleResult error}}`, result)
	})

	t.Run("grouped result", func(t *testing.T) {
		gs := NewGSWithGroupedResult("Why are these pizzas very popular?")
		result := gs.build()
		require.Equal(t, `_additional{generate(groupedResult:{ task: """Why are these pizzas very popular?""" }){groupedResult error}}`, result)
	})

}
