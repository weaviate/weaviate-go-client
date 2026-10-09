package jinaai

import "github.com/weaviate/weaviate-go-client/v6/modules"

func init() {
	modules.Register(*new(Multi2Vec))
}

// Multi2Vec is a vectorizer for text and image properties based on the multi2vec-jinaai module.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - multi2vec-jinaai].
//
// [Weaviate Docs - multi2vec-jinaai]: https://docs.weaviate.io/weaviate/model-providers/jinaai/embeddings-multimodal
type Multi2Vec struct {
	// BaseURL overrides the default request URL, e.g. a proxy.
	BaseURL string `json:"baseURL,omitempty"`
	// Model is the embedding model name.
	Model string `json:"model,omitempty"`
	// Dimensions is the size of the output vectors.
	Dimensions int `json:"dimensions,omitzero"`
	// TextFields and ImageFields name the properties to vectorize.
	TextFields  []string `json:"textFields,omitempty"`
	ImageFields []string `json:"imageFields,omitempty"`
	// Weights sets the contribution of each field to the combined vector.
	Weights Weights `json:"weights,omitzero"`
}

func (Multi2Vec) Name() string { return "multi2vec-jinaai" }

// Weights lists one weight per field, in the same order as the TextFields and ImageFields lists.
type Weights struct {
	TextFields  []float32 `json:"textFields,omitempty"`
	ImageFields []float32 `json:"imageFields,omitempty"`
}
