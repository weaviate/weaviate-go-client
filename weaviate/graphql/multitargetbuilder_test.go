package graphql

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pb "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

func TestMultiTargetArgumentBuilder(t *testing.T) {
	t.Run("Sum combination", func(t *testing.T) {
		builder := &MultiTargetArgumentBuilder{}
		builder.Sum("one", "two")
		out := builder.build()
		assert.Equal(t, "combinationMethod: sum, targetVectors: [\"one\",\"two\"]", out)
		targets := builder.togrpc()
		assert.Equal(t, targets.Combination, pb.CombinationMethod_COMBINATION_METHOD_TYPE_SUM)
		assert.Len(t, targets.TargetVectors, 2)
		assert.Contains(t, targets.TargetVectors, "one", "two")
		assert.Nil(t, targets.WeightsForTargets)
	})

	t.Run("Average combination", func(t *testing.T) {
		builder := &MultiTargetArgumentBuilder{}
		builder.Average("one", "two")
		out := builder.build()
		assert.Equal(t, "combinationMethod: average, targetVectors: [\"one\",\"two\"]", out)
		targets := builder.togrpc()
		assert.Equal(t, targets.Combination, pb.CombinationMethod_COMBINATION_METHOD_TYPE_AVERAGE)
		assert.Len(t, targets.TargetVectors, 2)
		assert.Contains(t, targets.TargetVectors, "one", "two")
		assert.Nil(t, targets.WeightsForTargets)
	})

	t.Run("Minimum combination", func(t *testing.T) {
		builder := &MultiTargetArgumentBuilder{}
		builder.Minimum("one", "two")
		out := builder.build()
		assert.Equal(t, "combinationMethod: minimum, targetVectors: [\"one\",\"two\"]", out)
		targets := builder.togrpc()
		assert.Equal(t, targets.Combination, pb.CombinationMethod_COMBINATION_METHOD_TYPE_MIN)
		assert.Len(t, targets.TargetVectors, 2)
		assert.Contains(t, targets.TargetVectors, "one", "two")
		assert.Nil(t, targets.WeightsForTargets)
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
		targets := builder.togrpc()
		assert.Equal(t, targets.Combination, pb.CombinationMethod_COMBINATION_METHOD_TYPE_MANUAL)
		assert.Len(t, targets.TargetVectors, 2)
		assert.Contains(t, targets.TargetVectors, "one", "two")
		assert.Len(t, targets.WeightsForTargets, 2)
		for _, w := range targets.WeightsForTargets {
			if w.Target == "one" {
				assert.Equal(t, w.Weight, float32(1))
			}
			if w.Target == "two" {
				assert.Equal(t, w.Weight, float32(2))
			}
		}
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
		targets := builder.togrpc()
		assert.Equal(t, targets.Combination, pb.CombinationMethod_COMBINATION_METHOD_TYPE_RELATIVE_SCORE)
		assert.Len(t, targets.TargetVectors, 2)
		assert.Contains(t, targets.TargetVectors, "one", "two")
		assert.Len(t, targets.WeightsForTargets, 2)
		for _, w := range targets.WeightsForTargets {
			if w.Target == "one" {
				assert.Equal(t, w.Weight, float32(1))
			}
			if w.Target == "two" {
				assert.Equal(t, w.Weight, float32(2))
			}
		}
	})

	t.Run("RelativeScoreMulti combination", func(t *testing.T) {
		builder := &MultiTargetArgumentBuilder{}
		builder.RelativeScoreMulti(map[string][]float32{"one": {1}, "two": {2, 3}})
		out := builder.build()
		// Have to use Contains because the order of the keys in the map is not guaranteed
		require.Contains(t, out, "combinationMethod: relativeScore")
		require.Contains(t, out, "targetVectors: ")
		require.Contains(t, out, "\"one\"")
		require.Contains(t, out, "\"two\",\"two\"")
		require.Contains(t, out, "weights: ")
		require.Contains(t, out, "one: 1")
		require.Contains(t, out, "two: [2,3]")
		targets := builder.togrpc()
		assert.Equal(t, targets.Combination, pb.CombinationMethod_COMBINATION_METHOD_TYPE_RELATIVE_SCORE)
		assert.Len(t, targets.TargetVectors, 3)
		assert.Contains(t, targets.TargetVectors, "one", "two")
		assert.Len(t, targets.WeightsForTargets, 3)
		for _, w := range targets.WeightsForTargets {
			if w.Target == "one" {
				assert.Equal(t, w.Weight, float32(1))
			}
		}
	})

	t.Run("ManualWeightsMulti combination", func(t *testing.T) {
		builder := &MultiTargetArgumentBuilder{}
		builder.ManualWeightsMulti(map[string][]float32{"one": {1}, "two": {2, 3}})
		out := builder.build()
		// Have to use Contains because the order of the keys in the map is not guaranteed
		require.Contains(t, out, "combinationMethod: manualWeights")
		require.Contains(t, out, "targetVectors: ")
		require.Contains(t, out, "\"one\"")
		require.Contains(t, out, "\"two\",\"two\"")
		require.Contains(t, out, "weights: ")
		require.Contains(t, out, "one: 1")
		require.Contains(t, out, "two: [2,3]")
		targets := builder.togrpc()
		assert.Equal(t, targets.Combination, pb.CombinationMethod_COMBINATION_METHOD_TYPE_MANUAL)
		assert.Len(t, targets.TargetVectors, 3)
		assert.Contains(t, targets.TargetVectors, "one", "two")
		assert.Len(t, targets.WeightsForTargets, 3)
		for _, w := range targets.WeightsForTargets {
			if w.Target == "one" {
				assert.Equal(t, w.Weight, float32(1))
			}
		}
	})
}
