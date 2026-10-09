package bind

import "github.com/weaviate/weaviate-go-client/v6/modules"

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
	TextFields    []string `json:"textFields,omitempty"`
	ImageFields   []string `json:"imageFields,omitempty"`
	AudioFields   []string `json:"audioFields,omitempty"`
	VideoFields   []string `json:"videoFields,omitempty"`
	IMUFields     []string `json:"imuFields,omitempty"`
	ThermalFields []string `json:"thermalFields,omitempty"`
	DepthFields   []string `json:"depthFields,omitempty"`
	// Weights sets the contribution of each field to the combined vector.
	Weights Weights `json:"weights,omitzero"`
}

func (Multi2Vec) Name() string { return "multi2vec-bind" }

// Weights lists one weight per field, in the same order as the TextFields, ImageFields, AudioFields, VideoFields, IMUFields, ThermalFields and DepthFields lists.
type Weights struct {
	TextFields    []float32 `json:"textFields,omitempty"`
	ImageFields   []float32 `json:"imageFields,omitempty"`
	AudioFields   []float32 `json:"audioFields,omitempty"`
	VideoFields   []float32 `json:"videoFields,omitempty"`
	IMUFields     []float32 `json:"imuFields,omitempty"`
	ThermalFields []float32 `json:"thermalFields,omitempty"`
	DepthFields   []float32 `json:"depthFields,omitempty"`
}
