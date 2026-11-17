package ollama

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

// Text2Vec is a vectorizer for text properties based on the text2vec-ollama module.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - text2vec-ollama].
//
// [Weaviate Docs - text2vec-ollama]: https://docs.weaviate.io/weaviate/model-providers/ollama/embeddings
type Text2Vec struct {
	// APIEndpoint is the Ollama server URL.
	APIEndpoint string `json:"apiEndpoint,omitempty"`
	// Model is the embedding model name.
	Model string `json:"model,omitempty"`
	// Properties limits vectorization to these properties.
	// By default, all text properties are vectorized.
	Properties []string `json:"properties,omitempty"`
}

func (Text2Vec) Name() string { return "text2vec-ollama" }

// Generative configures the generative-ollama module, either as a collection's
// default generative module or as a per-query provider in the generate package.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - generative-ollama].
//
// [Weaviate Docs - generative-ollama]: https://docs.weaviate.io/weaviate/model-providers/ollama/generative
type Generative struct {
	// APIEndpoint is the Ollama server URL.
	APIEndpoint string `json:"apiEndpoint,omitempty"`
	// Model name
	Model string `json:"model,omitempty"`

	// Provider holds query-time options. It is ignored in collection configuration.
	*Provider `json:"-"`
}

func (Generative) Name() string { return "generative-ollama" }

// Provider holds options that apply only when Generative is used as a per-query provider.
type Provider struct {
	Temperature *float64
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
		Kind: &proto.GenerativeProvider_Ollama{
			Ollama: &proto.GenerativeOllama{
				ApiEndpoint: api.NilZero(g.APIEndpoint),
				Model:       api.NilZero(g.Model),

				Temperature:     p.Temperature,
				Images:          internal.TextArray(p.Images),
				ImageProperties: internal.TextArray(p.ImageProperties),
			},
		},
	}
}
