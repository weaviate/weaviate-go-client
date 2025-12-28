package types

import "github.com/weaviate/weaviate-go-client/v6/internal/api"

type (
	Vector api.Vector

	// Vectors is a map of named vectors.
	// An empty string key is an alias for "default" vector.
	Vectors api.Vectors
)

// Vector implements [query.VectorKind].
func (v Vector) Vector() api.Vector { return api.Vector(v) }

// Vectors implements [query.VectorTarget].
func (v Vector) Vectors() []api.TargetVector { return []api.TargetVector{{Vector: v.Vector()}} }

func (vs Vectors) ToSlice() []Vector {
	out := make([]Vector, 0, len(vs))
	for _, v := range vs {
		out = append(out, Vector(v))
	}
	return out
}
