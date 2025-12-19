package types

import "github.com/weaviate/weaviate-go-client/v6/internal/api"

type (
	Vector api.Vector

	// Vectors is a map of named vectors.
	// An empty string key is an alias for "default" vector.
	Vectors api.Vectors
)

// Compile-time assertion that Vector implements NearVectorTarget.
var _ api.NearVectorTarget = (*Vector)(nil)

// CombinationMethod implements NearVectorTarget.
func (v Vector) CombinationMethod() api.CombinationMethod {
	return v.CombinationMethod()
}

// Targets implements api.NearVectorTarget.
func (v Vector) Vectors() []api.TargetVector {
	return v.Vectors()
}

func (vs Vectors) ToSlice() []*Vector {
	out := make([]*Vector, len(vs))
	for _, v := range vs {
		out = append(out, (*Vector)(v))
	}
	return out
}
