package api

import (
	"encoding/json"
)

const DefaultVectorName = "default"

type Vector struct {
	Name   string
	Single []float32
	Multi  [][]float32
}

var _ json.Unmarshaler = (*Vector)(nil)

// UnmarshalJSON implements json.Unmarshaler.
func (v *Vector) UnmarshalJSON(data []byte) error {
	var single []float32
	if err := json.Unmarshal(data, &single); err == nil {
		*v = Vector{Single: single}
		return nil
	}

	var multi [][]float32
	if err := json.Unmarshal(data, &multi); err != nil {
		return err
	}
	*v = Vector{Multi: multi}
	return nil
}

// Vectors is a map of named vectors. An empty string is an alias for "default" vector.
type Vectors map[string]Vector

var _ json.Unmarshaler = (*Vector)(nil)

func (vs *Vectors) UnmarshalJSON(data []byte) error {
	var vectors map[string]json.RawMessage
	if err := json.Unmarshal(data, &vectors); err != nil {
		return err
	}

	*vs = Vectors{}
	for k, data := range vectors {
		var v Vector
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		v.Name = k
		(*vs)[k] = v
	}
	return nil
}
