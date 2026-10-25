package internal

import (
	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/modules"
)

type WeightedFields struct {
	ImageFields []modules.WeightedField
	TextFields  []modules.WeightedField
}

func WithWeightedFields(module any, wfs WeightedFields) (map[string]any, error) {
	dest := make(map[string]any)
	if err := internal.Encode(module, dest); err != nil {
		return nil, err
	}
	addFields(dest, "textFields", wfs.TextFields)
	addFields(dest, "imageFields", wfs.ImageFields)
	return dest, nil
}

func addFields(m map[string]any, key string, wfs []modules.WeightedField) {
	if len(wfs) == 0 {
		return
	}
	var fields []string
	var weights []float32
	for _, wf := range wfs {
		fields = append(fields, wf.Property)
		weights = append(weights, wf.Weight)
	}

	m[key] = fields
	if _, ok := m["weights"]; !ok {
		m["weights"] = make(map[string]any)
	}
	m["weights"].(map[string]any)[key] = weights
}

func DecodeWithWeightedFields(m map[string]any, dest any) (wfs WeightedFields, err error) {
	if err = internal.Decode(m, dest); err != nil {
		return
	}

	readFields(m, "textFields", &wfs.TextFields)
	readFields(m, "imageFields", &wfs.ImageFields)
	return
}

func readFields(m map[string]any, key string, dest *[]modules.WeightedField) {
	var fields []string
	if f, ok := m[key]; ok {
		var ok bool
		if fields, ok = f.([]string); ok {
			var weights []float32
			if w, ok := m["weights"]; ok {
				if weightsMap, ok := w.(map[string]any); ok {
					weights = weightsMap[key].([]float32)
				}
			}
			for i, field := range fields {
				wf := modules.WeightedField{
					Property: field,
				}
				if i <= len(weights)-1 {
					wf.Weight = weights[i]
				}
				*dest = append(*dest, wf)
			}
		}
	}
}
