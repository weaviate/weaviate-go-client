package aws

import "github.com/weaviate/weaviate-go-client/v6/modules"

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
	Model string `json:"model,omitempty"`
	// Dimensions is the size of the output vectors.
	Dimensions int `json:"dimensions,omitzero"`
	// TextFields and ImageFields name the properties to vectorize.
	TextFields  []string `json:"textFields,omitempty"`
	ImageFields []string `json:"imageFields,omitempty"`
	// Weights sets the contribution of each field to the combined vector.
	Weights Weights `json:"weights,omitzero"`
	// Region is the AWS region to run the model in.
	Region string `json:"region,omitempty"`
}

func (Multi2Vec) Name() string { return "multi2vec-aws" }

// Weights lists one weight per field, in the same order as the TextFields and ImageFields lists.
type Weights struct {
	TextFields  []float32 `json:"textFields,omitempty"`
	ImageFields []float32 `json:"imageFields,omitempty"`
}
