package openai

import (
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/modules"
	proto "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

func init() {
	modules.Register(*new(Text2Vec))
	modules.Register(*new(Generative))
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
	BaseURL          string   `json:"baseURL,omitempty"`
	Model            string   `json:"model,omitempty"`
	APIVersion       string   `json:"apiVersion,omitempty"`
	Temperature      *float64 `json:"temperature,omitempty"`
	TopP             *float64 `json:"topP,omitempty"`
	MaxTokens        *int64   `json:"maxTokens,omitempty"`
	FrequencyPenalty *float64 `json:"frequencyPenalty,omitempty"`
	PresencePenalty  *float64 `json:"presencePenalty,omitempty"`
	ReasoningEffort  string   `json:"reasoningEffort,omitempty"`
	Verbosity        string   `json:"verbosity,omitempty"`

	ResourceName string `json:"resourceName,omitempty"`
	DeploymentID string `json:"deploymentId,omitempty"`

	*Provider `json:"-"`
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

type Provider struct {
	N               *int64
	StopSequences   []string
	Images          []string
	ImageProperties []string
	ReturnMetadata  bool
}

func (g Generative) GenerativeProvider() proto.GenerativeProvider {
	return proto.GenerativeProvider{
		ReturnMetadata: g.ReturnMetadata,
		Kind: &proto.GenerativeProvider_Openai{
			Openai: &proto.GenerativeOpenAI{
				BaseUrl:    api.NilZero(g.BaseURL),
				Model:      api.NilZero(g.Model),
				ApiVersion: api.NilZero(g.APIVersion),

				Temperature:      g.Temperature,
				TopP:             g.TopP,
				MaxTokens:        g.MaxTokens,
				FrequencyPenalty: g.FrequencyPenalty,
				PresencePenalty:  g.PresencePenalty,
				ReasoningEffort:  effort[g.ReasoningEffort],
				Verbosity:        verbosity[g.Verbosity],

				IsAzure:      new(g.DeploymentID != "" || g.ResourceName != ""),
				ResourceName: api.NilZero(g.ResourceName),
				DeploymentId: api.NilZero(g.DeploymentID),

				N:               g.N,
				Stop:            &proto.TextArray{Values: g.StopSequences},
				Images:          &proto.TextArray{Values: g.Images},
				ImageProperties: &proto.TextArray{Values: g.ImageProperties},
			},
		},
	}
}

var effort = map[string]*proto.GenerativeOpenAI_ReasoningEffort{
	MinimalEffort: new(proto.GenerativeOpenAI_REASONING_EFFORT_MINIMAL),
	LowEffort:     new(proto.GenerativeOpenAI_REASONING_EFFORT_LOW),
	MediumEffort:  new(proto.GenerativeOpenAI_REASONING_EFFORT_MEDIUM),
	HighEffort:    new(proto.GenerativeOpenAI_REASONING_EFFORT_HIGH),
}

var verbosity = map[string]*proto.GenerativeOpenAI_Verbosity{
	LowVerbosity:    new(proto.GenerativeOpenAI_VERBOSITY_LOW),
	MediumVerbosity: new(proto.GenerativeOpenAI_VERBOSITY_MEDIUM),
	HighVerbosity:   new(proto.GenerativeOpenAI_VERBOSITY_HIGH),
}
