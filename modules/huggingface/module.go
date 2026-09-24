package huggingface

import "github.com/weaviate/weaviate-go-client/v6/modules"

func init() {
	modules.Register(*new(Text2Vec))
}

// Text2Vec is a vectorizer for text properties based on the text2vec-huggingface module.
// Unset fields inherit the server defaults.
//
// Model and PassageModel are mutually exclusive; the server rejects a config that sets both.
//
// See https://docs.weaviate.io/weaviate/model-providers/huggingface/embeddings.
type Text2Vec struct {
	// Model defaults to sentence-transformers/msmarco-bert-base-dot-v5.
	Model string `json:"model,omitempty"`
	// PassageModel is an alias the server reads when Model is unset.
	// Prefer Model; this field keeps configs created by other clients intact.
	PassageModel string `json:"passageModel,omitempty"`
	// EndpointURL points to a dedicated inference endpoint; when set, the server skips model checks.
	EndpointURL string `json:"endpointURL,omitempty"`
	// WaitForModel waits for the model to be loaded. Defaults to false.
	WaitForModel *bool `json:"waitForModel,omitempty" nest:"options"`
	// UseGPU runs inference on a GPU. Defaults to false.
	UseGPU *bool `json:"useGPU,omitempty" nest:"options"`
	// UseCache enables the Inference API cache. Defaults to true.
	UseCache *bool `json:"useCache,omitempty" nest:"options"`
	// Properties limits vectorization to these properties.
	// By default, all text properties are vectorized.
	Properties []string `json:"properties,omitempty"`
}

func (Text2Vec) Name() string { return "text2vec-huggingface" }
