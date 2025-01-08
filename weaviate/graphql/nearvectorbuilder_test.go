package graphql

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate/entities/models"
)

func TestNearMultiVectorBuilder_build(t *testing.T) {
	t.Run("Sum combination with vector", func(t *testing.T) {
		vector := NearVectorArgumentBuilder{}
		targets := MultiTargetArgumentBuilder{}
		str := vector.WithVector([]float32{1, 2, 3}).WithTargets(targets.Sum("one", "two")).build()
		fmt.Println(str)
		require.NotContains(t, str, "vectorPerTarget: ")
		require.Contains(t, str, "vector: [1,2,3]")
		require.Contains(t, str, "combinationMethod: sum")
		require.Contains(t, str, "targetVectors: [\"one\",\"two\"]")
		require.NotContains(t, str, "weights: ")
	})

	t.Run("Average combination with vector per target", func(t *testing.T) {
		vector := NearVectorArgumentBuilder{}
		targets := MultiTargetArgumentBuilder{}
		str := vector.WithVectorPerTarget(map[string]models.Vector{"one": []float32{1, 2, 3}, "two": []float32{4, 5, 6}}).WithTargets(targets.Average("one", "two")).build()
		require.Contains(t, str, "vectorPerTarget: ")
		require.NotContains(t, str, "vector: ")
		require.Contains(t, str, "one: [[1,2,3]]")
		require.Contains(t, str, "two: [[4,5,6]]")
		require.Contains(t, str, "combinationMethod: average")
		require.Contains(t, str, "targetVectors: [\"one\",\"two\"]")
		require.NotContains(t, str, "weights: ")
	})

	t.Run("Minimum combination with all", func(t *testing.T) {
		vector := NearVectorArgumentBuilder{}
		targets := MultiTargetArgumentBuilder{}
		str := vector.WithVector([]float32{1, 2, 3}).WithTargets(targets.Minimum("one", "two")).WithDistance(0.01).build()
		require.NotContains(t, str, "vectorPerTarget: ")
		require.Contains(t, str, "vector: [1,2,3]")
		require.Contains(t, str, "combinationMethod: minimum")
		require.Contains(t, str, "targetVectors: [\"one\",\"two\"]")
		require.Contains(t, str, "distance: 0.01")
		require.NotContains(t, str, "weights: ")
	})

	t.Run("ManualWeights combination with vector", func(t *testing.T) {
		vector := NearVectorArgumentBuilder{}
		targets := MultiTargetArgumentBuilder{}
		str := vector.WithVector([]float32{1, 2, 3}).WithTargets(targets.ManualWeights(map[string]float32{"one": 1, "two": 2})).build()
		require.NotContains(t, str, "vectorPerTarget: ")
		require.Contains(t, str, "vector: [1,2,3]")
		require.Contains(t, str, "combinationMethod: manualWeights")
		require.Contains(t, str, "targetVectors: ")
		require.Contains(t, str, "\"one\"")
		require.Contains(t, str, "\"two\"")
		require.Contains(t, str, "weights: ")
		require.Contains(t, str, "one: 1")
		require.Contains(t, str, "two: 2")
	})

	t.Run("RelativeScore combination with vector", func(t *testing.T) {
		vector := NearVectorArgumentBuilder{}
		targets := MultiTargetArgumentBuilder{}
		str := vector.WithVector([]float32{1, 2, 3}).WithTargets(targets.RelativeScore(map[string]float32{"one": 1, "two": 2})).build()
		require.NotContains(t, str, "vectorPerTarget: ")
		require.Contains(t, str, "vector: [1,2,3]")
		require.Contains(t, str, "combinationMethod: relativeScore")
		require.Contains(t, str, "targetVectors: ")
		require.Contains(t, str, "\"one\"")
		require.Contains(t, str, "\"two\"")
		// order is not reproducible
		require.Contains(t, str, "weights: {")
		require.Contains(t, str, "one: 1")
		require.Contains(t, str, "two: 2")
	})

	t.Run("No combination with vector per target", func(t *testing.T) {
		vector := NearVectorArgumentBuilder{}
		str := vector.WithVectorPerTarget(map[string]models.Vector{"one": []float32{1, 2, 3}, "two": []float32{4, 5, 6}}).build()
		require.Contains(t, str, "vectorPerTarget: ")
		require.Contains(t, str, "one: [[1,2,3]]")
		require.Contains(t, str, "two: [[4,5,6]]")
		require.NotContains(t, str, "vector:")
		require.NotContains(t, str, "combinationMethod:")
		require.NotContains(t, str, "weights:")
	})

	t.Run("No combination with multiple vectors per target", func(t *testing.T) {
		vector := NearVectorArgumentBuilder{}
		str := vector.WithVectorsPerTarget(map[string][]models.Vector{"one": {[]float32{1, 2, 3}, []float32{7, 8, 9}}, "two": {[]float32{4, 5, 6}}}).build()
		require.Contains(t, str, "vectorPerTarget: ")
		require.Contains(t, str, "one: [[1,2,3],[7,8,9]]")
		require.Contains(t, str, "two: [[4,5,6]]")
		require.NotContains(t, str, "vector:")
		require.NotContains(t, str, "combinationMethod:")
		require.NotContains(t, str, "weights:")
	})

	t.Run("No combination with vector per target and target vectors", func(t *testing.T) {
		vector := NearVectorArgumentBuilder{}
		str := vector.WithVectorsPerTarget(map[string][]models.Vector{
			"one": {[]float32{1, 2, 3}, []float32{7, 8, 9}},
			"two": {[]float32{4, 5, 6}},
		},
		).WithTargetVectors("one", "two").build()
		require.Contains(t, str, "vectorPerTarget: ")
		require.Contains(t, str, "one: [[1,2,3],[7,8,9]]")
		require.Contains(t, str, "two: [[4,5,6]]")
		require.Contains(t, str, "targetVectors: ")
		require.Contains(t, str, "\"one\",\"one\"")
		require.Contains(t, str, "\"two\"")

		require.NotContains(t, str, "vector:")
		require.NotContains(t, str, "combinationMethod:")
		require.NotContains(t, str, "weights:")
	})

	t.Run("multi vector with target vector", func(t *testing.T) {
		nearVector := NearVectorArgumentBuilder{}
		str := nearVector.WithVector([][]float32{{1, 2, 3}, {1, 2, 3}}).WithTargetVectors("colbert").build()
		require.Equal(t, str, "nearVector:{vector: [[1,2,3],[1,2,3]] targetVectors: [\"colbert\"]}")
	})

	t.Run("vector with target vector", func(t *testing.T) {
		nearVector := NearVectorArgumentBuilder{}
		str := nearVector.WithVector([]float32{1, 2, 3}).WithTargetVectors("regular").build()
		require.Equal(t, str, "nearVector:{vector: [1,2,3] targetVectors: [\"regular\"]}")
	})
}
