package openai

import "github.com/weaviate/weaviate-go-client/v6/modules"

func init() {
	modules.Register(*new(Text2Vec))
}

// Text2Vec is a vectorizer for text properties
// based on the text2vec-openai module.
// Unset fields inherit the server defaults.
//
// To use Azure OpenAI, set ResourceName and DeploymentID together:
// the server detects Azure from the presence of both.
//
// See https://docs.weaviate.io/weaviate/model-providers/openai/embeddings
// and https://docs.weaviate.io/weaviate/model-providers/openai-azure/embeddings.
type Text2Vec struct {
	// Model defaults to text-embedding-3-small on the server,
	// which validates the name.
	Model string `json:"model,omitempty"`
	// Dimensions is the size of the output vectors.
	// Only the v3 models (text-embedding-3-*) support it.
	Dimensions int `json:"dimensions,omitzero"`
	// ModelType selects the model family.
	ModelType ModelType `json:"type,omitempty"`
	// ModelVersion applies only to the legacy models ("ada" etc.).
	ModelVersion string `json:"modelVersion,omitempty"`
	// BaseURL overrides where API requests go, e.g. a proxy.
	BaseURL string `json:"baseURL,omitempty"`
	// Endpoint is the API path appended to BaseURL, e.g. "/api/v3/embeddings".
	Endpoint string `json:"endpoint,omitempty"`
	// ResourceName is the Azure OpenAI resource name; required for Azure.
	ResourceName string `json:"resourceName,omitempty"`
	// DeploymentID is the Azure OpenAI deployment ID; required for Azure.
	DeploymentID string `json:"deploymentId,omitempty"`
	// Properties limits vectorization to these properties.
	// By default, all text properties are vectorized.
	Properties []string `json:"properties,omitempty"`
}

func (Text2Vec) Name() string { return "text2vec-openai" }

// ModelType selects the OpenAI model family.
type ModelType string

const (
	ModelTypeText ModelType = "text"
	ModelTypeCode ModelType = "code"
)
