package morph

import "github.com/weaviate/weaviate-go-client/v6/modules"

func init() {
	modules.Register(*new(Text2Vec))
}

// Text2Vec is a vectorizer for text properties based on the text2vec-morph module.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - text2vec-morph].
//
// [Weaviate Docs - text2vec-morph]: https://docs.weaviate.io/weaviate/model-providers/morph/embeddings
type Text2Vec struct {
	// BaseURL overrides where API requests go, e.g. a proxy.
	BaseURL string `json:"baseURL,omitempty"`
	// Model is the embedding model name.
	Model string `json:"model,omitempty"`
	// Properties limits vectorization to these properties.
	// By default, all text properties are vectorized.
	Properties []string `json:"properties,omitempty"`
	// Endpoint is the API path appended to BaseURL.
	Endpoint string `json:"endpoint,omitempty"`
}

func (Text2Vec) Name() string { return "text2vec-morph" }
