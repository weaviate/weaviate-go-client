package openai

import "github.com/weaviate/weaviate-go-client/v6/modules"

func init() {
	modules.Register(*new(Text2Vec))
}

// Text2Vec is a vectorizer for text properties based on the text2vec-openai module.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - text2vec-openai] and [Weaviate Docs - text2vec-openai (Azure)].
//
// [Weaviate Docs - text2vec-openai]: https://docs.weaviate.io/weaviate/model-providers/openai/embeddings
// [Weaviate Docs - text2vec-openai (Azure)]: https://docs.weaviate.io/weaviate/model-providers/openai-azure/embeddings
type Text2Vec struct {
	// BaseURL overrides where API requests go, e.g. a proxy.
	BaseURL string `json:"baseURL,omitempty"`
	// Model is the embedding model name.
	Model string `json:"model,omitempty"`
	// Dimensions is the size of the output vectors.
	Dimensions int `json:"dimensions,omitzero"`
	// Properties limits vectorization to these properties.
	// By default, all text properties are vectorized.
	Properties []string `json:"properties,omitempty"`
	// ModelVersion applies only to the legacy models ("ada" etc.).
	ModelVersion string `json:"modelVersion,omitempty"`
	// ModelType selects the model family.
	ModelType string `json:"type,omitempty"`
	// Endpoint is the API path appended to BaseURL, e.g. "/api/v3/embeddings".
	Endpoint string `json:"endpoint,omitempty"`
	// ResourceName is the Azure OpenAI resource name.
	ResourceName string `json:"resourceName,omitempty"`
	// DeploymentID is the Azure OpenAI deployment ID.
	DeploymentID string `json:"deploymentId,omitempty"`
}

func (Text2Vec) Name() string { return "text2vec-openai" }

const (
	TextModel = "text"
	CodeModel = "code"
)

type Generative struct {
	BaseURL          string  `json:"baseURL"`
	Model            string  `json:"model"`
	APIVersion       string  `json:"apiVersion"`
	Temperature      float32 `json:"temperature"`
	TopP             float32 `json:"topP"`
	MaxTokens        int     `json:"maxTokens"`
	FrequencyPenalty float32 `json:"frequencyPenalty"`
	PresencePenalty  float32 `json:"presencePenalty"`
	ReasoningEffort  string  `json:"reasoningEffort"`
	Verbosity        string  `json:"verbosity"`
}

const (
	MinimalEffort = "minimal"
	LowEffort     = "low"
	MediumEffort  = "medium"
	HighEffort    = "high"
)

const (
	LowVerbosity    = "low"
	MediumVerbosity = "medium"
	HighVerbosity   = "high"
)

func (Generative) Name() string { return "generative-openai" }
