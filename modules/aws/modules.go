package aws

import (
	"github.com/weaviate/weaviate-go-client/v6/modules"
	"github.com/weaviate/weaviate-go-client/v6/modules/internal"
)

func init() {
	modules.Register(*new(Text2Vec))
	modules.Register(*new(Multi2Vec))
}

// Text2Vec is a vectorizer for text properties based on the text2vec-aws module.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - text2vec-aws].
//
// [Weaviate Docs - text2vec-aws]: https://docs.weaviate.io/weaviate/model-providers/aws/embeddings
type Text2Vec struct {
	// BedrockModel is the Bedrock model ID.
	BedrockModel string `json:"model,omitempty"`
	// Dimensions is the size of the output vectors.
	Dimensions int `json:"dimensions,omitzero"`
	// Properties limits vectorization to these properties.
	// By default, all text properties are vectorized.
	Properties []string `json:"properties,omitempty"`
	// Service selects the AWS service that hosts the model.
	Service string `json:"service,omitempty"`
	// Region is the AWS region to run the model in.
	Region string `json:"region,omitempty"`
	// SageMakerEndpoint is the endpoint name.
	SageMakerEndpoint string `json:"endpoint,omitempty"`
	// SageMakerModel for multi-model endpoints.
	SageMakerModel string `json:"targetModel,omitempty"`
	// SageMakerVariant is the production variant of the [SageMakerModel].
	SageMakerVariant string `json:"targetVariant,omitempty"`
}

func (Text2Vec) Name() string { return "text2vec-aws" }

const (
	Bedrock   = "bedrock"
	SageMaker = "sagemaker"
)

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

	TextFields  []modules.WeightedField
	ImageFields []modules.WeightedField
}

type multi2vecJSON struct {
	Model      string `json:"model,omitempty"`
	Dimensions int    `json:"dimensions,omitzero"`
	Region     string `json:"region,omitempty"`
}

func (Multi2Vec) Name() string { return "multi2vec-aws" }

func (m2v Multi2Vec) EncodeMap() (map[string]any, error) {
	return internal.EncodeWithWeightedFields(multi2vecJSON{
		Model:      m2v.Model,
		Dimensions: m2v.Dimensions,
		Region:     m2v.Region,
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
		Model:       dest.Model,
		Dimensions:  dest.Dimensions,
		Region:      dest.Region,
		TextFields:  wfs.TextFields,
		ImageFields: wfs.ImageFields,
	}
	return nil
}
