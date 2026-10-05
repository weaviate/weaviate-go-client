package voyageai

import "github.com/weaviate/weaviate-go-client/v6/modules"

func init() {
	modules.Register(*new(Text2Vec))
}

// Text2Vec is a vectorizer for text properties based on the text2vec-voyageai module.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - text2vec-voyageai].
//
// [Weaviate Docs - text2vec-voyageai]: https://docs.weaviate.io/weaviate/model-providers/voyageai/embeddings
type Text2Vec struct {
	// BaseURL overrides where API requests go, e.g. a proxy.
	BaseURL string `json:"baseURL,omitempty"`
	// Model is the embedding model name.
	Model string `json:"model,omitempty"`
	// Dimensions is the size of the output vectors.
	Dimensions int `json:"dimensions,omitzero"`
	// Properties limits vectorization to these properties.
	// By default, all text properties are vectorized.
	Properties []string `json:"properties,omitempty"`
	// Truncate cuts input texts to fit within the model's context length.
	Truncate *bool `json:"truncate,omitempty"`
}

func (Text2Vec) Name() string { return "text2vec-voyageai" }
