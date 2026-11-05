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

// Generative is a generative module based on the generative-digitalocean module.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - generative-digitalocean].
//
// [Weaviate Docs - generative-digitalocean]: https://docs.weaviate.io/weaviate/model-providers/digitalocean/generative
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
	// Stop ends generation when the model emits any of these sequences.
	Stop []string `json:"stop,omitempty"`

	// ReturnMetadata includes token usage in the response.
	// It applies to query-time generation only and is not part of the collection config.
	ReturnMetadata bool `json:"-"`
}

func (Generative) Name() string { return "generative-digitalocean" }

func (g Generative) GenerativeProvider() *proto.GenerativeProvider {
	return &proto.GenerativeProvider{
		ReturnMetadata: g.ReturnMetadata,
		Kind: &proto.GenerativeProvider_Digitalocean{
			Digitalocean: &proto.GenerativeDigitalOcean{
				BaseUrl:          api.NilZero(g.BaseURL),
				Model:            api.NilZero(g.Model),
				Temperature:      g.Temperature,
				TopP:             g.TopP,
				MaxTokens:        g.MaxTokens,
				FrequencyPenalty: g.FrequencyPenalty,
				PresencePenalty:  g.PresencePenalty,
				Stop:             internal.TextArray(g.Stop),
			},
		},
	}
}
