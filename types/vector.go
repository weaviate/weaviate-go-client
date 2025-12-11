package types

type Vector struct {
	Name   string
	Single []float32
	Multi  [][]float32
}

func (v Vector) ToProto() {}
