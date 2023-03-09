package graphql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerativeSearch_build(t *testing.T) {
	t.Run("single result", func(t *testing.T) {
		gs := NewGenerativeSearch().SingleResult("Describe this pizza : {name}")
		result := gs.build()

		assert.Equal(t, `generate(singleResult:{prompt:"""Describe this pizza : {name}"""})`, result.Name)
		assert.ElementsMatch(t, []Field{{Name: "singleResult"}, {Name: "error"}}, result.Fields)
	})

	t.Run("grouped result", func(t *testing.T) {
		gs := NewGenerativeSearch().GroupedResult("Why are these pizzas very popular?")
		result := gs.build()

		assert.Equal(t, `generate(groupedResult:{task:"""Why are these pizzas very popular?"""})`, result.Name)
		assert.ElementsMatch(t, []Field{{Name: "groupedResult"}, {Name: "error"}}, result.Fields)
	})

	t.Run("with single result and grouped result", func(t *testing.T) {
		gs := NewGenerativeSearch().SingleResult("Describe this pizza : {name}").GroupedResult("Why are these pizzas very popular?")
		result := gs.build()

		assert.Equal(t, `generate(singleResult:{prompt:"""Describe this pizza : {name}"""} groupedResult:{task:"""Why are these pizzas very popular?"""})`, result.Name)
		assert.ElementsMatch(t, []Field{{Name: "singleResult"}, {Name: "groupedResult"}, {Name: "error"}}, result.Fields)
	})
}
