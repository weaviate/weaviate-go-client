package voyageai

import "github.com/weaviate/weaviate-go-client/v6/modules"

func init() {
	modules.Register(*new(Multi2Vec))
}

// Multi2Vec is a vectorizer for text, image and video properties based on the multi2vec-voyageai module.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - multi2vec-voyageai].
//
// [Weaviate Docs - multi2vec-voyageai]: https://docs.weaviate.io/weaviate/model-providers/voyageai/embeddings-multimodal
type Multi2Vec struct {
	// BaseURL overrides the default request URL, e.g. a proxy.
	BaseURL string `json:"baseURL,omitempty"`
	// Model is the embedding model name.
	Model string `json:"model,omitempty"`
	// TextFields, ImageFields and VideoFields name the properties to vectorize.
	TextFields  []string `json:"textFields,omitempty"`
	ImageFields []string `json:"imageFields,omitempty"`
	VideoFields []string `json:"videoFields,omitempty"`
	// Weights sets the contribution of each field to the combined vector.
	Weights Weights `json:"weights,omitzero"`
	// Truncate cuts inputs to fit within the model's context length.
	Truncate *bool `json:"truncate,omitempty"`
}

func (Multi2Vec) Name() string { return "multi2vec-voyageai" }

// Weights lists one weight per field, in the same order as the TextFields, ImageFields and VideoFields lists.
type Weights struct {
	TextFields  []float32 `json:"textFields,omitempty"`
	ImageFields []float32 `json:"imageFields,omitempty"`
	VideoFields []float32 `json:"videoFields,omitempty"`
}
