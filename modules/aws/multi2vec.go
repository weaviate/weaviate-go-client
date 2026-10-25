package aws

import (
	"github.com/weaviate/weaviate-go-client/v6/modules"
	modinternal "github.com/weaviate/weaviate-go-client/v6/modules/internal"
)

func init() {
	modules.Register(*new(Multi2Vec))
}

// Multi2Vec is a vectorizer for text and image properties based on the multi2vec-aws module.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - AWS embeddings].
//
// [Weaviate Docs - AWS embeddings]: https://docs.weaviate.io/weaviate/model-providers/aws/embeddings
type Multi2Vec struct {
	// Model is the Bedrock model ID.
	Model string
	// Dimensions is the size of the output vectors.
	Dimensions int
	// Region is the AWS region to run the model in.
	Region string

	ImageFields []modules.WeightedField
	TextFields  []modules.WeightedField
}

type multi2vecJSON struct {
	// Model is the Bedrock model ID.
	Model string `json:"model,omitempty"`
	// Dimensions is the size of the output vectors.
	Dimensions int `json:"dimensions,omitzero"`
	// Region is the AWS region to run the model in.
	Region string `json:"region,omitempty"`
}

func (Multi2Vec) Name() string { return "multi2vec-aws" }

func (m2v Multi2Vec) EncodeMap() (map[string]any, error) {
	return modinternal.WithWeightedFields(multi2vecJSON{
		Model:      m2v.Model,
		Dimensions: m2v.Dimensions,
		Region:     m2v.Region,
	}, modinternal.WeightedFields{
		ImageFields: m2v.ImageFields,
		TextFields:  m2v.TextFields,
	})
}

func (m2v *Multi2Vec) DecodeMap(m map[string]any) error {
	var dest multi2vecJSON
	wfs, err := modinternal.DecodeWithWeightedFields(m, &dest)
	if err != nil {
		return err
	}

	*m2v = Multi2Vec{
		Model:       dest.Model,
		Dimensions:  dest.Dimensions,
		Region:      dest.Region,
		ImageFields: wfs.ImageFields,
		TextFields:  wfs.TextFields,
	}
	return nil
}
