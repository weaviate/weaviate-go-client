package contextualai

import (
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/modules"
	"github.com/weaviate/weaviate-go-client/v6/modules/internal"
	proto "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

func init() {
	modules.Register(*new(Generative))
}

// Generative configures the generative-contextualai module, either as a collection's
// default generative module or as a per-query provider in the generate package.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - generative-contextualai].
//
// [Weaviate Docs - generative-contextualai]: https://docs.weaviate.io/weaviate/model-providers/contextualai/generative
type Generative struct {
	// Model name
	Model        string `json:"model,omitempty"`
	MaxNewTokens *int64 `json:"maxNewTokens,omitempty"`

	Temperature *float64 `json:"temperature,omitempty"`
	TopP        *float64 `json:"topP,omitempty"`

	// SystemPrompt replaces the model's default system prompt.
	SystemPrompt string `json:"systemPrompt,omitempty"`
	// AvoidCommentary limits the response to grounded facts, without commentary.
	AvoidCommentary *bool `json:"avoidCommentary,omitempty"`
	// Knowledge lists extra text passed to the model with the retrieved objects.
	Knowledge []string `json:"knowledge,omitempty"`

	// Provider holds query-time options. It is ignored in collection configuration.
	*Provider `json:"-"`
}

func (Generative) Name() string { return "generative-contextualai" }

// Provider holds options that apply only when Generative is used as a per-query provider.
type Provider struct {
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
		Kind: &proto.GenerativeProvider_Contextualai{
			Contextualai: &proto.GenerativeContextualAI{
				Model: api.NilZero(g.Model),

				MaxNewTokens: g.MaxNewTokens,
				Temperature:  g.Temperature,
				TopP:         g.TopP,

				SystemPrompt:    api.NilZero(g.SystemPrompt),
				AvoidCommentary: g.AvoidCommentary,
				Knowledge:       internal.TextArray(g.Knowledge),
			},
		},
	}
}
