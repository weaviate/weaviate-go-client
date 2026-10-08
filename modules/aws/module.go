package aws

import "github.com/weaviate/weaviate-go-client/v6/modules"

func init() {
	modules.Register(*new(Text2Vec))
}

// Text2Vec is a vectorizer for text properties based on the text2vec-aws module.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - text2vec-aws].
//
// [Weaviate Docs - text2vec-aws]: https://docs.weaviate.io/weaviate/model-providers/aws/embeddings
type Text2Vec struct {
	// Model is the Bedrock model ID.
	Model string `json:"model,omitempty"`
	// Dimensions is the size of the output vectors.
	Dimensions int `json:"dimensions,omitzero"`
	// Properties limits vectorization to these properties.
	// By default, all text properties are vectorized.
	Properties []string `json:"properties,omitempty"`
	// Service selects the AWS service that hosts the model.
	Service string `json:"service,omitempty"`
	// Region is the AWS region to run the model in.
	Region string `json:"region,omitempty"`
	// Endpoint is the SageMaker endpoint name.
	Endpoint string `json:"endpoint,omitempty"`
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
