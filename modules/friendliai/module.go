package friendliai

import (
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/modules"
	proto "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

func init() {
	modules.Register(*new(Generative))
}

// Generative configures the generative-friendliai module, either as a collection's
// default generative module or as a per-query provider in the generate package.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - generative-friendliai].
//
// [Weaviate Docs - generative-friendliai]: https://docs.weaviate.io/weaviate/model-providers/friendliai/generative
type Generative struct {
	// BaseURL overrides the default request URL, e.g. a proxy.
	BaseURL string `json:"baseURL,omitempty"`
	// Model name
	Model     string `json:"model,omitempty"`
	MaxTokens *int64 `json:"maxTokens,omitempty"`

	Temperature *float64 `json:"temperature,omitempty"`

	// Provider holds query-time options. It is ignored in collection configuration.
	*Provider `json:"-"`
}

func (Generative) Name() string { return "generative-friendliai" }

// Provider holds options that apply only when Generative is used as a per-query provider.
type Provider struct {
	N    *int64
	TopP *float64
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
		Kind: &proto.GenerativeProvider_Friendliai{
			Friendliai: &proto.GenerativeFriendliAI{
				BaseUrl: api.NilZero(g.BaseURL),
				Model:   api.NilZero(g.Model),

				MaxTokens:   g.MaxTokens,
				Temperature: g.Temperature,

				N:    p.N,
				TopP: p.TopP,
			},
		},
	}
}
