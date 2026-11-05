package meta

import (
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/modules"
	"github.com/weaviate/weaviate-go-client/v6/modules/internal"
	proto "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

func init() {
	modules.Register(*new(Generative))
}

// Generative is a generative module based on the generative-meta module.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - generative-meta].
//
// [Weaviate Docs - generative-meta]: https://docs.weaviate.io/weaviate/model-providers/meta/generative
type Generative struct {
	// BaseURL overrides the default request URL, e.g. a proxy.
	BaseURL string `json:"baseURL,omitempty"`
	// Model is the inference model name.
	Model string `json:"model,omitempty"`
	// Temperature controls the randomness of the output, between 0 and 2.
	Temperature *float64 `json:"temperature,omitempty"`
	// TopP is the nucleus sampling probability, between 0 and 1.
	TopP *float64 `json:"topP,omitempty"`
	// MaxTokens caps the number of generated tokens.
	MaxTokens *int64 `json:"maxTokens,omitempty"`
	// FrequencyPenalty penalizes tokens by their frequency in the output, between -2 and 2.
	FrequencyPenalty *float64 `json:"frequencyPenalty,omitempty"`
	// PresencePenalty penalizes tokens already present in the prompt, between -2 and 2.
	PresencePenalty *float64 `json:"presencePenalty,omitempty"`
	// ReasoningEffort selects how much reasoning the model spends on the request.
	// See the Effort constants in this package.
	ReasoningEffort string `json:"reasoningEffort,omitempty"`

	// Images and ImageProperties make images available to the model.
	// They apply to query-time generation only and are not part of the collection config.
	Images          []string `json:"-"`
	ImageProperties []string `json:"-"`

	// ReturnMetadata includes token usage in the response.
	// It applies to query-time generation only and is not part of the collection config.
	ReturnMetadata bool `json:"-"`
}

func (Generative) Name() string { return "generative-meta" }

func (g Generative) GenerativeProvider() *proto.GenerativeProvider {
	return &proto.GenerativeProvider{
		ReturnMetadata: g.ReturnMetadata,
		Kind: &proto.GenerativeProvider_Meta{
			Meta: &proto.GenerativeMeta{
				BaseUrl:          api.NilZero(g.BaseURL),
				Model:            api.NilZero(g.Model),
				Temperature:      g.Temperature,
				TopP:             g.TopP,
				MaxTokens:        g.MaxTokens,
				FrequencyPenalty: g.FrequencyPenalty,
				PresencePenalty:  g.PresencePenalty,
				ReasoningEffort:  effort[g.ReasoningEffort],
				Images:           internal.TextArray(g.Images),
				ImageProperties:  internal.TextArray(g.ImageProperties),
			},
		},
	}
}

const (
	NoEffort      = "none"
	MinimalEffort = "minimal"
	LowEffort     = "low"
	MediumEffort  = "medium"
	HighEffort    = "high"
	XHighEffort   = "xhigh"
)

var effort = map[string]*proto.GenerativeMeta_ReasoningEffort{
	NoEffort:      proto.GenerativeMeta_REASONING_EFFORT_NONE.Enum(),
	MinimalEffort: proto.GenerativeMeta_REASONING_EFFORT_MINIMAL.Enum(),
	LowEffort:     proto.GenerativeMeta_REASONING_EFFORT_LOW.Enum(),
	MediumEffort:  proto.GenerativeMeta_REASONING_EFFORT_MEDIUM.Enum(),
	HighEffort:    proto.GenerativeMeta_REASONING_EFFORT_HIGH.Enum(),
	XHighEffort:   proto.GenerativeMeta_REASONING_EFFORT_XHIGH.Enum(),
}
