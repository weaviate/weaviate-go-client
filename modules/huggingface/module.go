package huggingface

import "github.com/weaviate/weaviate-go-client/v6/modules"

func init() {
	modules.Register(*new(Text2Vec))
}

// Text2Vec is a vectorizer for text properties
// based on the text2vec-huggingface module.
// Unset fields inherit the server defaults.
//
// Set either Model, or PassageModel and QueryModel as a pair;
// the server rejects Model combined with the other two.
//
// See https://docs.weaviate.io/weaviate/model-providers/huggingface/embeddings.
type Text2Vec struct {
	// Model defaults to sentence-transformers/msmarco-bert-base-dot-v5.
	Model string `json:"model,omitempty"`
	// PassageModel vectorizes objects at import; pair it with QueryModel.
	PassageModel string `json:"passageModel,omitempty"`
	// QueryModel vectorizes search queries; pair it with PassageModel.
	QueryModel string `json:"queryModel,omitempty"`
	// EndpointURL points to a dedicated inference endpoint;
	// when set, the server skips model checks.
	EndpointURL string `json:"endpointURL,omitempty"`
	// Options control the Hugging Face Inference API behavior.
	Options *Options `json:"options,omitempty"`
	// Properties limits vectorization to these properties.
	// By default, all text properties are vectorized.
	Properties []string `json:"properties,omitempty"`
}

func (Text2Vec) Name() string { return "text2vec-huggingface" }

// Options control the Hugging Face Inference API behavior.
type Options struct {
	// WaitForModel waits for the model to be loaded. Defaults to false.
	WaitForModel bool `json:"waitForModel,omitzero"`
	// UseGPU runs inference on a GPU. Defaults to false.
	UseGPU bool `json:"useGPU,omitzero"`
	// UseCache enables the Inference API cache. It is a pointer
	// because the server defaults to true.
	UseCache *bool `json:"useCache,omitempty"`
}
