package api

const DefaultVectorName = "default"

type Vector struct {
	Name   string
	Single []float32
	Multi  [][]float32
}

var (
	_ NearVectorTarget = (*Vector)(nil)
	_ TargetVector     = (*Vector)(nil)
)

func (v *Vector) CombinationMethod() CombinationMethod {
	return CombinationMethodUnspecified
}
func (v *Vector) Vectors() []TargetVector { return []TargetVector{v} }
func (v *Vector) Vector() *Vector         { return v }
func (v *Vector) Weight() float32         { return 0 }

// Vectors is a map of named vectors. An empty string is an alias for "default" vector.
type Vectors map[string]Vector

func (vs Vectors) ToSlice() []Vector {
	out := make([]Vector, 0, len(vs))
	for _, v := range vs {
		out = append(out, v)
	}
	return out
}
