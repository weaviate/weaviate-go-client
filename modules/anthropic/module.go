package anthropic

import (
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/modules"
	"github.com/weaviate/weaviate-go-client/v6/modules/internal"
	proto "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

func init() {
	modules.Register(*new(Generative))
}

// Generative configures the generative-anthropic module, either as a collection's
// default generative module or as a per-query provider in the generate package.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - generative-anthropic].
//
// [Weaviate Docs - generative-anthropic]: https://docs.weaviate.io/weaviate/model-providers/anthropic/generative
type Generative struct {
	// BaseURL overrides the default request URL, e.g. a proxy.
	BaseURL string `json:"baseURL,omitempty"`
	// Model name
	Model         string   `json:"model,omitempty"`
	MaxTokens     *int64   `json:"maxTokens,omitempty"`
	StopSequences []string `json:"stopSequences,omitempty"`

	Temperature *float64 `json:"temperature,omitempty"`
	TopK        *int64   `json:"topK,omitempty"`
	TopP        *float64 `json:"topP,omitempty"`

	// Provider holds query-time options. It is ignored in collection configuration.
	*Provider `json:"-"`
}

func (Generative) Name() string { return "generative-anthropic" }

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
		Kind: &proto.GenerativeProvider_Anthropic{
			Anthropic: &proto.GenerativeAnthropic{
				BaseUrl: api.NilZero(g.BaseURL),
				Model:   api.NilZero(g.Model),

				MaxTokens:     g.MaxTokens,
				Temperature:   g.Temperature,
				TopK:          g.TopK,
				TopP:          g.TopP,
				StopSequences: internal.TextArray(g.StopSequences),

				Images:          internal.TextArray(p.Images),
				ImageProperties: internal.TextArray(p.ImageProperties),
			},
		},
	}
}
