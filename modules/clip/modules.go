package clip

import (
	"github.com/weaviate/weaviate-go-client/v6/modules"
	"github.com/weaviate/weaviate-go-client/v6/modules/internal"
)

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
	InferenceURL string
	// TextFields and ImageFields name the properties to vectorize.
	TextFields  []modules.WeightedField `json:"-"`
	ImageFields []modules.WeightedField `json:"-"`
}

func (Multi2Vec) Name() string { return "multi2vec-clip" }

type multi2vecJSON struct {
	InferenceURL string `json:"inferenceUrl,omitempty"`
}

func (m2v Multi2Vec) EncodeMap() (map[string]any, error) {
	return internal.EncodeWithWeightedFields(multi2vecJSON{
		InferenceURL: m2v.InferenceURL,
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
		InferenceURL: dest.InferenceURL,
		TextFields:   wfs.TextFields,
		ImageFields:  wfs.ImageFields,
	}
	return nil
}
