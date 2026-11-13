package digitalocean

import (
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/modules"
	"github.com/weaviate/weaviate-go-client/v6/modules/internal"
	proto "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

func init() {
	modules.Register(*new(Text2Vec))
	modules.Register(*new(Generative))
}

// Text2Vec is a vectorizer for text properties based on the text2vec-digitalocean module.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - text2vec-digitalocean].
//
// [Weaviate Docs - text2vec-digitalocean]: https://docs.weaviate.io/weaviate/model-providers/digitalocean/embeddings
type Text2Vec struct {
	// BaseURL overrides the default request URL, e.g. a proxy.
	BaseURL string `json:"baseURL,omitempty"`
	// Model is the embedding model name.
	Model string `json:"model,omitempty"`
	// Properties limits vectorization to these properties.
	// By default, all text properties are vectorized.
	Properties []string `json:"properties,omitempty"`
}

func (Text2Vec) Name() string { return "text2vec-digitalocean" }

// Generative configures the generative-digitalocean module, either as a collection's
// default generative module or as a per-query provider in the generate package.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - generative-digitalocean].
//
// [Weaviate Docs - generative-digitalocean]: https://docs.weaviate.io/weaviate/model-providers/digitalocean/generative
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

func (Generative) Name() string { return "generative-digitalocean" }

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
		Kind: &proto.GenerativeProvider_Digitalocean{
			Digitalocean: &proto.GenerativeDigitalOcean{
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
