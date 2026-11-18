package deepseek

import (
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/modules"
	"github.com/weaviate/weaviate-go-client/v6/modules/internal"
	proto "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

func init() {
	modules.Register(*new(Generative))
}

// Generative configures the generative-deepseek module, either as a collection's
// default generative module or as a per-query provider in the generate package.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - generative-deepseek].
//
// [Weaviate Docs - generative-deepseek]: https://docs.weaviate.io/weaviate/model-providers/deepseek/generative
type Generative struct {
	// BaseURL overrides the default request URL, e.g. a proxy.
	BaseURL string `json:"baseURL,omitempty"`
	// Model name
	Model         string   `json:"model,omitempty"`
	MaxTokens     *int64   `json:"maxTokens,omitempty"`
	StopSequences []string `json:"stop,omitempty"`

	Temperature      *float64 `json:"temperature,omitempty"`
	TopP             *float64 `json:"topP,omitempty"`
	FrequencyPenalty *float64 `json:"frequencyPenalty,omitempty"`
	PresencePenalty  *float64 `json:"presencePenalty,omitempty"`

	// Provider holds query-time options. It is ignored in collection configuration.
	*Provider `json:"-"`
}

func (Generative) Name() string { return "generative-deepseek" }

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
		Kind: &proto.GenerativeProvider_Deepseek{
			Deepseek: &proto.GenerativeDeepseek{
				BaseUrl: api.NilZero(g.BaseURL),
				Model:   api.NilZero(g.Model),

				MaxTokens:        g.MaxTokens,
				Temperature:      g.Temperature,
				TopP:             g.TopP,
				FrequencyPenalty: g.FrequencyPenalty,
				PresencePenalty:  g.PresencePenalty,
				Stop:             internal.TextArray(g.StopSequences),
			},
		},
	}
}
