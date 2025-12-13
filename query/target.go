package query

import (
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/types"
)

// MultiVectorTarget comprises multiple target vectors and their respective weights.
// Construct an appropriate MultiVectorTarget using one of these functions:
// - Sum
// - Average
// - Min
// - ManualWeights
// - RelativeScore
type MultiVectorTarget struct {
	combinationMethod api.CombinationMethod
	targets           []api.TargetVector
}

// Compile-time assertion that MultiVectorTarget implements api.NearVectorTarget.
var _ api.NearVectorTarget = (*MultiVectorTarget)(nil)

// WeightedTarget assigns a weight to a vector target used in a query.
// Use Target() to construct one.
type WeightedTarget struct {
	v      api.Vector
	weight float32
}

// Compile-time assertion that WeighteTarget implements api.TargetVector
var _ api.TargetVector = (*WeightedTarget)(nil)

func (wt WeightedTarget) Weight() float32 {
	return wt.weight
}

func (wt WeightedTarget) Vector() *api.Vector {
	return &wt.v
}

func Target(v types.Vector, weight float32) WeightedTarget {
	return WeightedTarget{v: api.Vector(v), weight: weight}
}

func Sum(vectors ...types.Vector) MultiVectorTarget {
	return zeroWeightTargets(api.CombinationMethodSum, vectors)
}

func Min(vectors ...types.Vector) MultiVectorTarget {
	return zeroWeightTargets(api.CombinationMethodMin, vectors)
}

func Average(vectors ...types.Vector) MultiVectorTarget {
	return zeroWeightTargets(api.CombinationMethodAverage, vectors)
}

func ManualWeights(vectors ...WeightedTarget) MultiVectorTarget {
	// Explicitly cast []WeightedTarget to []api.TargetVector
	targets := make([]api.TargetVector, len(vectors))
	for _, v := range vectors {
		targets = append(targets, v)
	}
	return MultiVectorTarget{
		combinationMethod: api.CombinationMethodManualWeights,
		targets:           targets,
	}
}

func RelativeScore(vectors ...WeightedTarget) MultiVectorTarget {
	// Explicitly cast []WeightedTarget to []api.TargetVector
	targets := make([]api.TargetVector, len(vectors))
	for _, v := range vectors {
		targets = append(targets, v)
	}
	return MultiVectorTarget{
		combinationMethod: api.CombinationMethodRelativeScore,
		targets:           targets,
	}
}

// Combine target vectors into a MultiVectorTarget keeping the weight unset.
// The server will determine combination method will determine the weights
func zeroWeightTargets(cm api.CombinationMethod, vectors []types.Vector) MultiVectorTarget {
	targets := make([]api.TargetVector, len(vectors))
	for _, v := range vectors {
		targets = append(targets, WeightedTarget{v: api.Vector(v)})
	}

	return MultiVectorTarget{
		combinationMethod: cm,
		targets:           targets,
	}
}

func (m MultiVectorTarget) CombinationMethod() api.CombinationMethod {
	return m.combinationMethod
}

func (m MultiVectorTarget) Vectors() []api.TargetVector {
	return m.targets
}
