package types

import "github.com/weaviate/weaviate-go-client/v6/internal/api"

type (
	Vector api.Vector

	// Vectors is a map of named vectors.
	// An empty string key is an alias for "default" vector.
	Vectors api.Vectors
)

var (
	// Compile-time assertion that Vector implements NearVectorTarget.
	_ api.NearVectorTarget = (*Vector)(nil)
	// Compile-time assertion that Vector implements TargetVector.
	_ api.TargetVector = (*Vector)(nil)
)

// CombinationMethod implements NearVectorTarget.
func (v *Vector) CombinationMethod() api.CombinationMethod {
	return api.CombinationMethodUnspecified
}

// Targets implements api.NearVectorTarget.
func (v *Vector) Vectors() []api.TargetVector {
	return []api.TargetVector{v}
}

func (v *Vector) Weight() float32     { return 0 }
func (v *Vector) Vector() *api.Vector { return (*api.Vector)(v) }

func (vs Vectors) ToSlice() []Vector {
	out := make([]Vector, 0, len(vs))
	for _, v := range vs {
		out = append(out, Vector(v))
	}
	return out
}
