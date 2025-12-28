package api

const DefaultVectorName = "default"

type Vector struct {
	Name   string
	Single []float32
	Multi  [][]float32
}

// Vectors is a map of named vectors. An empty string is an alias for "default" vector.
type Vectors map[string]Vector

func (vs Vectors) ToSlice() []Vector {
	out := make([]Vector, 0, len(vs))
	for _, v := range vs {
		out = append(out, v)
	}
	return out
}
