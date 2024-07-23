package graphql

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMultiTargetArgumentBuilder(t *testing.T) {
	t.Run("Sum combination", func(t *testing.T) {
		builder := &MultiTargetArgumentBuilder{}
		builder.Sum("one", "two")
		out := builder.build()
		assert.Equal(t, "combinationMethod: sum, targetVectors: [\"one\",\"two\"]", out)
	})

	t.Run("Average combination", func(t *testing.T) {
		builder := &MultiTargetArgumentBuilder{}
		builder.Average("one", "two")
		out := builder.build()
		assert.Equal(t, "combinationMethod: average, targetVectors: [\"one\",\"two\"]", out)
	})

	t.Run("Minimum combination", func(t *testing.T) {
		builder := &MultiTargetArgumentBuilder{}
		builder.Minimum("one", "two")
		out := builder.build()
		assert.Equal(t, "combinationMethod: minimum, targetVectors: [\"one\",\"two\"]", out)
	})

	t.Run("ManualWeights combination", func(t *testing.T) {
		builder := &MultiTargetArgumentBuilder{}
		builder.ManualWeights(map[string]float32{"one": 1, "two": 2})
		out := builder.build()
		// Have to use Contains because the order of the keys in the map is not guaranteed
		require.Contains(t, out, "combinationMethod: manualWeights")
		require.Contains(t, out, "targetVectors: ")
		require.Contains(t, out, "\"one\"")
		require.Contains(t, out, "\"two\"")
		require.Contains(t, out, "weights: ")
		require.Contains(t, out, "one: 1")
		require.Contains(t, out, "two: 2")
	})

	t.Run("RelativeScore combination", func(t *testing.T) {
		builder := &MultiTargetArgumentBuilder{}
		builder.RelativeScore(map[string]float32{"one": 1, "two": 2})
		out := builder.build()
		// Have to use Contains because the order of the keys in the map is not guaranteed
		require.Contains(t, out, "combinationMethod: relativeScore")
		require.Contains(t, out, "targetVectors: ")
		require.Contains(t, out, "\"one\"")
		require.Contains(t, out, "\"two\"")
		require.Contains(t, out, "weights: ")
		require.Contains(t, out, "one: 1")
		require.Contains(t, out, "two: 2")
	})
}
