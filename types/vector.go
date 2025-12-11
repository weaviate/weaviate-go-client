package types

type Vector struct {
	Name   string
	Single []float32
	Multi  [][]float32
}

func (v Vector) ToProto() {}

// Vectors is a map of named vectors.
// The key is the vector name (empty string for default vector).
type Vectors map[string]Vector
