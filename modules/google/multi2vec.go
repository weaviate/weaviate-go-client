package google

import "github.com/weaviate/weaviate-go-client/v6/modules"

func init() {
	modules.Register(*new(Multi2Vec))
}

// Multi2Vec is a vectorizer for text, image, video and audio properties based on the multi2vec-google module.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - multi2vec-google].
//
// [Weaviate Docs - multi2vec-google]: https://docs.weaviate.io/weaviate/model-providers/google/embeddings-multimodal
type Multi2Vec struct {
	// APIEndpoint is the host name.
	APIEndpoint string `json:"apiEndpoint,omitempty"`
	// Model is the embedding model name.
	Model string `json:"model,omitempty"`
	// Dimensions is the size of the output vectors.
	Dimensions int `json:"dimensions,omitzero"`
	// TextFields, ImageFields, VideoFields and AudioFields name the properties to vectorize.
	TextFields  []string `json:"textFields,omitempty"`
	ImageFields []string `json:"imageFields,omitempty"`
	VideoFields []string `json:"videoFields,omitempty"`
	AudioFields []string `json:"audioFields,omitempty"`
	// Weights sets the contribution of each field to the combined vector.
	Weights Weights `json:"weights,omitzero"`
	// ProjectID is the Google Cloud project ID.
	ProjectID string `json:"projectId,omitempty"`
	// Location is the Vertex AI region to run the model in.
	Location string `json:"location,omitempty"`
	// VideoIntervalSeconds is the length of the video segments to embed.
	VideoIntervalSeconds int `json:"videoIntervalSeconds,omitzero"`
}

func (Multi2Vec) Name() string { return "multi2vec-google" }

// Weights lists one weight per field, in the same order as the TextFields, ImageFields, VideoFields and AudioFields lists.
type Weights struct {
	TextFields  []float32 `json:"textFields,omitempty"`
	ImageFields []float32 `json:"imageFields,omitempty"`
	VideoFields []float32 `json:"videoFields,omitempty"`
	AudioFields []float32 `json:"audioFields,omitempty"`
}
