package internal

import (
	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/modules"
)

type WeightedFields struct {
	TextFields    []modules.WeightedField
	ImageFields   []modules.WeightedField
	AudioFields   []modules.WeightedField
	VideoFields   []modules.WeightedField
	DepthFields   []modules.WeightedField
	ThermalFields []modules.WeightedField
	IMUFields     []modules.WeightedField
}

func EncodeWithWeightedFields(module any, wfs WeightedFields) (map[string]any, error) {
	dest := make(map[string]any)
	if module != nil {
		if err := internal.Encode(module, dest); err != nil {
			return nil, err
		}
	}
	addFields(dest, "textFields", wfs.TextFields)
	addFields(dest, "imageFields", wfs.ImageFields)
	addFields(dest, "audioFields", wfs.AudioFields)
	addFields(dest, "videoFields", wfs.VideoFields)
	addFields(dest, "depthFields", wfs.DepthFields)
	addFields(dest, "thermalFields", wfs.ThermalFields)
	addFields(dest, "imuFields", wfs.IMUFields)
	return dest, nil
}

func addFields(m map[string]any, key string, wfs []modules.WeightedField) {
	if len(wfs) == 0 {
		return
	}
	var fields []string
	var weights []float32
	for _, wf := range wfs {
		fields = append(fields, wf.Name)
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
	readFields(m, "audioFields", &wfs.AudioFields)
	readFields(m, "videoFields", &wfs.VideoFields)
	readFields(m, "depthFields", &wfs.DepthFields)
	readFields(m, "thermalFields", &wfs.ThermalFields)
	readFields(m, "imuFields", &wfs.IMUFields)
	return
}

func readFields(m map[string]any, key string, dest *[]modules.WeightedField) {
	if f, ok := m[key]; ok {
		fields, ok := f.([]string)
		if !ok {
			return
		}
		var weights []float32
		if w, ok := m["weights"]; ok {
			if weightsMap, ok := w.(map[string]any); ok {
				weights = weightsMap[key].([]float32)
			}
		}
		for i, field := range fields {
			wf := modules.WeightedField{
				Name: field,
			}
			if i <= len(weights)-1 {
				wf.Weight = weights[i]
			}
			*dest = append(*dest, wf)
		}
	}
}
