package transformers

import "github.com/weaviate/weaviate-go-client/v6/modules"

func init() {
	modules.Register(*new(Text2Vec))
}

// Text2Vec is a vectorizer for text properties based on the text2vec-transformers module.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - text2vec-transformers].
//
// [Weaviate Docs - text2vec-transformers]: https://docs.weaviate.io/weaviate/model-providers/transformers/embeddings
type Text2Vec struct {
	// InferenceURL is the inference container URL for both passages and queries.
	InferenceURL string `json:"inferenceUrl,omitempty"`
	// Dimensions is the size of the output vectors.
	Dimensions int `json:"dimensions,omitzero"`
	// Properties limits vectorization to these properties.
	// By default, all text properties are vectorized.
	Properties []string `json:"properties,omitempty"`
	// PassageInferenceURL is the inference container URL for passages.
	PassageInferenceURL string `json:"passageInferenceUrl,omitempty"`
	// QueryInferenceURL is the inference container URL for queries.
	QueryInferenceURL string `json:"queryInferenceUrl,omitempty"`
	// PoolingStrategy selects how token embeddings are pooled into one vector.
	PoolingStrategy string `json:"poolingStrategy,omitempty"`
}

func (Text2Vec) Name() string { return "text2vec-transformers" }

const (
	MaskedMean = "masked_mean"
	CLS        = "cls"
)
