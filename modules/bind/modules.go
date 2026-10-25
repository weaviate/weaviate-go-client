package bind

import (
	"github.com/weaviate/weaviate-go-client/v6/modules"
	"github.com/weaviate/weaviate-go-client/v6/modules/internal"
)

func init() {
	modules.Register(*new(Multi2Vec))
}

// Multi2Vec is a vectorizer for text, image, audio, video, IMU, thermal and depth properties based on the multi2vec-bind module.
//
// See [Weaviate Docs - multi2vec-bind].
//
// [Weaviate Docs - multi2vec-bind]: https://docs.weaviate.io/weaviate/model-providers/imagebind/embeddings-multimodal
type Multi2Vec struct {
	// Each *Fields list names the properties of one modality to vectorize.
	TextFields    []modules.WeightedField `json:"-"`
	ImageFields   []modules.WeightedField `json:"-"`
	AudioFields   []modules.WeightedField `json:"-"`
	VideoFields   []modules.WeightedField `json:"-"`
	IMUFields     []modules.WeightedField `json:"-"`
	ThermalFields []modules.WeightedField `json:"-"`
	DepthFields   []modules.WeightedField `json:"-"`
}

func (Multi2Vec) Name() string { return "multi2vec-bind" }

func (m2v Multi2Vec) EncodeMap() (map[string]any, error) {
	return internal.EncodeWithWeightedFields(nil, internal.WeightedFields{
		ImageFields:   m2v.ImageFields,
		TextFields:    m2v.TextFields,
		AudioFields:   m2v.AudioFields,
		VideoFields:   m2v.VideoFields,
		DepthFields:   m2v.DepthFields,
		ThermalFields: m2v.ThermalFields,
		IMUFields:     m2v.IMUFields,
	})
}

func (m2v *Multi2Vec) DecodeMap(m map[string]any) error {
	type multi2vecJSON Multi2Vec // Avoid recursive calls to DecodeMap

	var dest multi2vecJSON
	wfs, err := internal.DecodeWithWeightedFields(m, &dest)
	if err != nil {
		return err
	}

	*m2v = Multi2Vec{
		TextFields:    wfs.TextFields,
		ImageFields:   wfs.ImageFields,
		AudioFields:   wfs.AudioFields,
		VideoFields:   wfs.VideoFields,
		DepthFields:   wfs.DepthFields,
		ThermalFields: wfs.ThermalFields,
		IMUFields:     wfs.IMUFields,
	}
	return nil
}
