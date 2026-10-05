package cohere

import "github.com/weaviate/weaviate-go-client/v6/modules"

func init() {
	modules.Register(*new(Text2Vec))
}

// Text2Vec is a vectorizer for text properties based on the text2vec-cohere module.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - text2vec-cohere].
//
// [Weaviate Docs - text2vec-cohere]: https://docs.weaviate.io/weaviate/model-providers/cohere/embeddings
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
	// Truncate sets how inputs longer than the model's context are cut.
	Truncate string `json:"truncate,omitempty"`
}

func (Text2Vec) Name() string { return "text2vec-cohere" }

const (
	TruncateNone  = "NONE"
	TruncateStart = "START"
	TruncateEnd   = "END"
	TruncateLeft  = "LEFT"
	TruncateRight = "RIGHT"
)
