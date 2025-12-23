package api

const DefaultVectorName = "default"

type Vector struct {
	Name   string
	Single []float32
	Multi  [][]float32
}

// Compile-time assertion that Vector implements NearVectorTarget.
var _ NearVectorTarget = (*Vector)(nil)

func (v *Vector) CombinationMethod() CombinationMethod {
	return CombinationMethodUnspecified
}

func (v Vector) Vectors() []TargetVector {
	return []TargetVector{targetVector{v: v}}
}

// targetVector implements TargetVector for a single Vector.
type targetVector struct{ v Vector }

var _ TargetVector = (*targetVector)(nil)

func (tv targetVector) Weight() float32 { return 0 }

func (tv targetVector) Vector() *Vector { return &tv.v }

// Vectors is a map of named vectors. An empty string is an alias for "default" vector.
type Vectors map[string]Vector

func (vs Vectors) ToSlice() []Vector {
	out := make([]Vector, 0, len(vs))
	for _, v := range vs {
		out = append(out, v)
	}
	return out
}
