package graphql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Neartextbuilder(t *testing.T) {
	t.Run("concepts with escaping", func(t *testing.T) {
		nt := NearTextArgumentBuilder{}

		str := nt.WithConcepts([]string{"\"I'm a complex concept\" says the string", "simple concept"}).build()
		expected := `nearText:{concepts: ["\"I'm a complex concept\" says the string","simple concept"]}`
		require.Equal(t, expected, str)
	})

	t.Run("move parameters with escaping", func(t *testing.T) {
		nt := NearTextArgumentBuilder{}
		mp := MoveParameters{Concepts: []string{"Extra quotes: \" ':", "no quotes"}}

		str := nt.WithConcepts([]string{"\"I'm a complex concept\" says the string", "simple concept"}).WithMoveAwayFrom(&mp).build()
		expected := `nearText:{concepts: ["\"I'm a complex concept\" says the string","simple concept"] moveAwayFrom: {concepts: ["Extra quotes: \" ':","no quotes"] force: 0}}`
		require.Equal(t, expected, str)
	})

	t.Run("concepts with escaping with targetVectors", func(t *testing.T) {
		nt := NearTextArgumentBuilder{}

		str := nt.WithConcepts([]string{"\"I'm a complex concept\" says the string", "simple concept"}).
			WithTargetVectors("targetVector").
			build()
		expected := `nearText:{concepts: ["\"I'm a complex concept\" says the string","simple concept"] targetVectors: ["targetVector"]}`
		require.Equal(t, expected, str)
	})
}
