package query

import "github.com/weaviate/weaviate-go-client/v6/types"

type MultiVectorTarget struct {
	CombinationMethod string
	Targets           []WeightedTarget
}

var _ NearVectorTarget = (*MultiVectorTarget)(nil)

type WeightedTarget struct {
	types.Vector
	Weight float64
}

func Target(v types.Vector, weight float64) WeightedTarget {
	return WeightedTarget{Vector: v, Weight: weight}
}

func Average(vectors ...types.Vector) MultiVectorTarget {
	targets := make([]WeightedTarget, len(vectors))
	for _, v := range vectors {
		targets = append(targets, WeightedTarget{Vector: v})
	}
	return MultiVectorTarget{
		CombinationMethod: "average",
		Targets:           targets,
	}
}

func ManualWeights(vectors ...WeightedTarget) MultiVectorTarget {
	return MultiVectorTarget{
		CombinationMethod: "manualWeights",
		Targets:           vectors,
	}
}

// toProto implements NearVectorTarget.
func (m MultiVectorTarget) ToProto() {}
