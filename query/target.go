package query

import (
	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/types"
)

// MultiVectorTarget comprises multiple target vectors and their respective weights.
// Construct an appropriate MultiVectorTarget using one of these functions:
// - Sum
// - Average
// - Max
// - ManualWeights
// - RelativeScore
type MultiVectorTarget struct {
	combinationMethod internal.CombinationMethod
	targets           []WeightedTarget
}

var _ NearVectorTarget = (*MultiVectorTarget)(nil)

// WeightedTarget assigns a weight to a vector target used in a query.
// Use Target() to construct one.
type WeightedTarget struct {
	types.Vector
	weight float64
}

// Assign a weight to a target vector.
func Target(v types.Vector, weight float64) WeightedTarget {
	return WeightedTarget{Vector: v, weight: weight}
}

func Sum(vectors ...types.Vector) MultiVectorTarget {
	return zeroWeightTargets(internal.CombinationMethodSum, vectors)
}

func Max(vectors ...types.Vector) MultiVectorTarget {
	return zeroWeightTargets(internal.CombinationMethodMax, vectors)
}

func Average(vectors ...types.Vector) MultiVectorTarget {
	return zeroWeightTargets(internal.CombinationMethodAverage, vectors)
}

func ManualWeights(vectors ...WeightedTarget) MultiVectorTarget {
	return MultiVectorTarget{
		combinationMethod: internal.CombinationMethodManualWeights,
		targets:           vectors,
	}
}

func RelativeScore(vectors ...WeightedTarget) MultiVectorTarget {
	return MultiVectorTarget{
		combinationMethod: internal.CombinationMethodRelativeScore,
		targets:           vectors,
	}
}

// Combine target vectors into a MultiVectorTarget keeping the weight unset.
func zeroWeightTargets(cm internal.CombinationMethod, vectors []types.Vector) MultiVectorTarget {
	targets := make([]WeightedTarget, len(vectors))
	for _, v := range vectors {
		targets = append(targets, WeightedTarget{Vector: v})
	}

	return MultiVectorTarget{
		combinationMethod: cm,
		targets:           targets,
	}
}

// toProto implements NearVectorTarget.
func (m MultiVectorTarget) ToProto() {}
