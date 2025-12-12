package types

type Vector struct {
	Name   string
	Single []float32
	Multi  [][]float32
}

// Vector implements query.NearVectorTarget
func (v Vector) ToProto() {}

// Vectors is a map of named vectors. An empty string is an alias for "default" vector.
type Vectors map[string]Vector

func (vs Vectors) ToSlice() []Vector {
	out := make([]Vector, len(vs))
	for _, v := range vs {
		out = append(out, v)
	}
	return out
}
