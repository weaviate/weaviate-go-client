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

// Generative configures the generative-meta module, either as a collection's
// default generative module or as a per-query provider in the generate package.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - generative-meta].
//
// [Weaviate Docs - generative-meta]: https://docs.weaviate.io/weaviate/model-providers/meta/generative
type Generative struct {
	// BaseURL overrides the default request URL, e.g. a proxy.
	BaseURL string `json:"baseURL,omitempty"`
	// Model name
	Model     string `json:"model,omitempty"`
	MaxTokens *int64 `json:"maxTokens,omitempty"`

	Temperature      *float64 `json:"temperature,omitempty"`
	TopP             *float64 `json:"topP,omitempty"`
	FrequencyPenalty *float64 `json:"frequencyPenalty,omitempty"`
	PresencePenalty  *float64 `json:"presencePenalty,omitempty"`
	ReasoningEffort  string   `json:"reasoningEffort,omitempty"`

	// Provider holds query-time options. It is ignored in collection configuration.
	*Provider `json:"-"`
}

const (
	NoEffort      = "none"
	MinimalEffort = "minimal"
	LowEffort     = "low"
	MediumEffort  = "medium"
	HighEffort    = "high"
	XHighEffort   = "xhigh"
)

func (Generative) Name() string { return "generative-meta" }

// Provider holds options that apply only when Generative is used as a per-query provider.
type Provider struct {
	// Images are base64-encoded images passed to the model with the prompt.
	Images []string
	// ImageProperties lists the object properties that contain images for the prompt.
	ImageProperties []string
	// ReturnMetadata requests provider metadata, such as token usage, in the response.
	ReturnMetadata bool
}

func (g Generative) GenerativeProvider() *proto.GenerativeProvider {
	p := g.Provider
	if p == nil {
		p = new(Provider)
	}
	return &proto.GenerativeProvider{
		ReturnMetadata: p.ReturnMetadata,
		Kind: &proto.GenerativeProvider_Meta{
			Meta: &proto.GenerativeMeta{
				BaseUrl: api.NilZero(g.BaseURL),
				Model:   api.NilZero(g.Model),

				MaxTokens:        g.MaxTokens,
				Temperature:      g.Temperature,
				TopP:             g.TopP,
				FrequencyPenalty: g.FrequencyPenalty,
				PresencePenalty:  g.PresencePenalty,
				ReasoningEffort:  effort[g.ReasoningEffort],

				Images:          internal.TextArray(p.Images),
				ImageProperties: internal.TextArray(p.ImageProperties),
			},
		},
	}
}

var effort = map[string]*proto.GenerativeMeta_ReasoningEffort{
	NoEffort:      proto.GenerativeMeta_REASONING_EFFORT_NONE.Enum(),
	MinimalEffort: proto.GenerativeMeta_REASONING_EFFORT_MINIMAL.Enum(),
	LowEffort:     proto.GenerativeMeta_REASONING_EFFORT_LOW.Enum(),
	MediumEffort:  proto.GenerativeMeta_REASONING_EFFORT_MEDIUM.Enum(),
	HighEffort:    proto.GenerativeMeta_REASONING_EFFORT_HIGH.Enum(),
	XHighEffort:   proto.GenerativeMeta_REASONING_EFFORT_XHIGH.Enum(),
}
