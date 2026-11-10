package mistral

import (
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/modules"
	proto "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

func init() {
	modules.Register(*new(Text2Vec))
	modules.Register(*new(Generative))
}

// Text2Vec is a vectorizer for text properties based on the text2vec-mistral module.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - text2vec-mistral].
//
// [Weaviate Docs - text2vec-mistral]: https://docs.weaviate.io/weaviate/model-providers/mistral/embeddings
type Text2Vec struct {
	// BaseURL overrides the default request URL, e.g. a proxy.
	BaseURL string `json:"baseURL,omitempty"`
	// Model is the embedding model name.
	Model string `json:"model,omitempty"`
	// Properties limits vectorization to these properties.
	// By default, all text properties are vectorized.
	Properties []string `json:"properties,omitempty"`
}

func (Text2Vec) Name() string { return "text2vec-mistral" }

// Generative configures the generative-mistral module, either as a collection's
// default generative module or as a per-query provider in the generate package.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - generative-mistral].
//
// [Weaviate Docs - generative-mistral]: https://docs.weaviate.io/weaviate/model-providers/mistral/generative
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

func (Generative) Name() string { return "generative-mistral" }

// Provider holds options that apply only when Generative is used as a per-query provider.
type Provider struct {
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
		Kind: &proto.GenerativeProvider_Mistral{
			Mistral: &proto.GenerativeMistral{
				BaseUrl: api.NilZero(g.BaseURL),
				Model:   api.NilZero(g.Model),

				MaxTokens:   g.MaxTokens,
				Temperature: g.Temperature,

				TopP: p.TopP,
			},
		},
	}
}
