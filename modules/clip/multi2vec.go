package clip

import "github.com/weaviate/weaviate-go-client/v6/modules"

func init() {
	modules.Register(*new(Multi2Vec))
}

// Multi2Vec is a vectorizer for text and image properties based on the multi2vec-clip module.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - multi2vec-clip].
//
// [Weaviate Docs - multi2vec-clip]: https://docs.weaviate.io/weaviate/model-providers/transformers/embeddings-multimodal
type Multi2Vec struct {
	// InferenceURL is the inference container URL.
	InferenceURL string `json:"inferenceUrl,omitempty"`
	// TextFields and ImageFields name the properties to vectorize.
	TextFields  []string `json:"textFields,omitempty"`
	ImageFields []string `json:"imageFields,omitempty"`
	// Weights sets the contribution of each field to the combined vector.
	Weights Weights `json:"weights,omitzero"`
}

func (Multi2Vec) Name() string { return "multi2vec-clip" }

// Weights lists one weight per field, in the same order as the TextFields and ImageFields lists.
type Weights struct {
	TextFields  []float32 `json:"textFields,omitempty"`
	ImageFields []float32 `json:"imageFields,omitempty"`
}
