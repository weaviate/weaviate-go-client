package api

const DefaultVectorName = "default"

type Vector struct {
	Name   string
	Single []float32
	Multi  [][]float32
}

// Vectors is a map of named vectors. An empty string is an alias for "default" vector.
type Vectors map[string]Vector
