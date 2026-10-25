package cohere

import (
	"github.com/weaviate/weaviate-go-client/v6/modules"
	"github.com/weaviate/weaviate-go-client/v6/modules/internal"
)

func init() {
	modules.Register(*new(Text2Vec))
	modules.Register(*new(Multi2Vec))
}

// Text2Vec is a vectorizer for text properties based on the text2vec-cohere module.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - text2vec-cohere].
//
// [Weaviate Docs - text2vec-cohere]: https://docs.weaviate.io/weaviate/model-providers/cohere/embeddings
type Text2Vec struct {
	// BaseURL overrides the default request URL, e.g. a proxy.
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

// Multi2Vec is a vectorizer for text and image properties based on the multi2vec-cohere module.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - multi2vec-cohere].
//
// [Weaviate Docs - multi2vec-cohere]: https://docs.weaviate.io/weaviate/model-providers/cohere/embeddings-multimodal
type Multi2Vec struct {
	// BaseURL overrides the default request URL, e.g. a proxy.
	BaseURL string
	// Model is the embedding model name.
	Model string
	// Dimensions is the size of the output vectors.
	Dimensions int
	// Truncate sets how inputs longer than the model's context are cut.
	Truncate string
	// TextFields and ImageFields name the properties to vectorize.
	TextFields  []modules.WeightedField `json:"-"`
	ImageFields []modules.WeightedField `json:"-"`
}

func (Multi2Vec) Name() string { return "multi2vec-cohere" }

type multi2vecJSON struct {
	BaseURL    string `json:"baseURL,omitempty"`
	Model      string `json:"model,omitempty"`
	Dimensions int    `json:"dimensions,omitzero"`
	Truncate   string `json:"truncate,omitempty"`
}

func (m2v Multi2Vec) EncodeMap() (map[string]any, error) {
	return internal.EncodeWithWeightedFields(multi2vecJSON{
		BaseURL:    m2v.BaseURL,
		Model:      m2v.Model,
		Dimensions: m2v.Dimensions,
		Truncate:   m2v.Truncate,
	}, internal.WeightedFields{
		TextFields:  m2v.TextFields,
		ImageFields: m2v.ImageFields,
	})
}

func (m2v *Multi2Vec) DecodeMap(m map[string]any) error {
	var dest multi2vecJSON
	wfs, err := internal.DecodeWithWeightedFields(m, &dest)
	if err != nil {
		return err
	}

	*m2v = Multi2Vec{
		BaseURL:     dest.BaseURL,
		Model:       dest.Model,
		Dimensions:  dest.Dimensions,
		Truncate:    dest.Truncate,
		TextFields:  wfs.TextFields,
		ImageFields: wfs.ImageFields,
	}
	return nil
}
