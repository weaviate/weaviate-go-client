package twelvelabs

import (
	"github.com/weaviate/weaviate-go-client/v6/modules"
	"github.com/weaviate/weaviate-go-client/v6/modules/internal"
)

func init() {
	modules.Register(*new(Multi2Vec))
}

// Multi2Vec is a vectorizer for text and image properties based on the multi2vec-twelvelabs module.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - multi2vec-twelvelabs].
//
// [Weaviate Docs - multi2vec-twelvelabs]: https://docs.weaviate.io/weaviate/model-providers/twelvelabs/embeddings-multimodal
type Multi2Vec struct {
	// BaseURL overrides the default request URL, e.g. a proxy.
	BaseURL string
	// Model is the embedding model name.
	Model string
	// TextFields and ImageFields name the properties to vectorize.
	TextFields  []modules.WeightedField `json:"-"`
	ImageFields []modules.WeightedField `json:"-"`
}

func (Multi2Vec) Name() string { return "multi2vec-twelvelabs" }

type multi2vecJSON struct {
	BaseURL string `json:"baseURL,omitempty"`
	Model   string `json:"model,omitempty"`
}

func (m2v Multi2Vec) EncodeMap() (map[string]any, error) {
	return internal.EncodeWithWeightedFields(multi2vecJSON{
		BaseURL: m2v.BaseURL,
		Model:   m2v.Model,
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
		TextFields:  wfs.TextFields,
		ImageFields: wfs.ImageFields,
	}
	return nil
}
