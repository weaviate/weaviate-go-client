package query

import (
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
)

// VectorKind defines what "vector" can mean in the context of search.
// For NearVector, NearObject, NearText, and NearMedia search, 2 meanings
// are possible:
//
//   - User provides one or more _embeddings_, which should be compared to
//     embeddings stored in one or more _vector indexes_.
//
//   - User provides one or more _names of the vector indexes_, in which the
//     embeddings that should be used as search targets are stored.
//
// Each "vector" can be assigned a specific weight, which [WeightedVector]
// encodes. Several vector targets can be combined in a [MultiVectorTarget].
//
// At the time of writing Weaviate allows a query to contain multiple _embeddings_,
// but only one _string input_ (be that text, UUID or a base64-encoded image).
// For this reason, embedding and vector index name can be provided as a single
// value [types.Vector], but string input and vector index name must be passed
// separately in [NearText], [NearMedia], and [NearObject] queries.
type VectorKind interface{ Vector() api.Vector }

type VectorName string

var (
	_ VectorKind   = (*VectorName)(nil)
	_ VectorTarget = (*VectorName)(nil)
)

func (name VectorName) Vector() api.Vector {
	return api.Vector{Name: string(name)}
}

func (name VectorName) Vectors() []api.TargetVector {
	return []api.TargetVector{{Vector: name.Vector()}}
}

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

var _ VectorTarget = (*MultiVectorTarget)(nil)

// WeightedVector assigns a weight to a vector target used in a query.
// Use Target() to construct one.
type WeightedVector[V VectorKind] struct {
	vector V
	weight float32
}

func Weighted[V VectorKind](v V, weight float32) WeightedVector[V] {
	return WeightedVector[V]{vector: v, weight: weight}
}

func Sum[V VectorKind](vectors []V) *MultiVectorTarget {
	return zeroWeightTargets(api.CombinationMethodSum, vectors)
}

func Min[V VectorKind](vectors []V) *MultiVectorTarget {
	return zeroWeightTargets(api.CombinationMethodMin, vectors)
}

func Average[V VectorKind](vectors []V) *MultiVectorTarget {
	return zeroWeightTargets(api.CombinationMethodAverage, vectors)
}

func ManualWeights[V VectorKind](vectors []WeightedVector[V]) *MultiVectorTarget {
	targets := make([]api.TargetVector, len(vectors))
	for i, v := range vectors {
		targets[i] = api.TargetVector{Vector: v.vector.Vector(), Weight: &v.weight}
	}
	return &MultiVectorTarget{
		combinationMethod: api.CombinationMethodManualWeights,
		targets:           targets,
	}
}

func RelativeScore[V VectorKind](vectors []WeightedVector[V]) *MultiVectorTarget {
	targets := make([]api.TargetVector, len(vectors))
	for i, v := range vectors {
		targets[i] = api.TargetVector{Vector: v.vector.Vector(), Weight: &v.weight}
	}
	return &MultiVectorTarget{
		combinationMethod: api.CombinationMethodRelativeScore,
		targets:           targets,
	}
}

// Combine target vectors into a MultiVectorTarget keeping the weight unset.
// The server will determine combination method will determine the weights
func zeroWeightTargets[V VectorKind](cm api.CombinationMethod, vectors []V) *MultiVectorTarget {
	targets := make([]api.TargetVector, len(vectors))
	for i, v := range vectors {
		targets[i] = api.TargetVector{Vector: v.Vector()}
	}

	return &MultiVectorTarget{
		combinationMethod: cm,
		targets:           targets,
	}
}

func (m *MultiVectorTarget) CombinationMethod() api.CombinationMethod {
	return m.combinationMethod
}

func (m *MultiVectorTarget) Vectors() []api.TargetVector {
	return m.targets
}
