package databricks

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

// Text2Vec is a vectorizer for text properties based on the text2vec-databricks module.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - text2vec-databricks].
//
// [Weaviate Docs - text2vec-databricks]: https://docs.weaviate.io/weaviate/model-providers/databricks/embeddings
type Text2Vec struct {
	// Endpoint is the URL of the model server.
	Endpoint string `json:"endpoint,omitempty"`
	// Properties limits vectorization to these properties.
	// By default, all text properties are vectorized.
	Properties []string `json:"properties,omitempty"`
	// Instruction is a system prompt passed to the embedding model.
	Instruction string `json:"instruction,omitempty"`
}

func (Text2Vec) Name() string { return "text2vec-databricks" }

// Generative configures the generative-databricks module, either as a collection's
// default generative module or as a per-query provider in the generate package.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - generative-databricks].
//
// [Weaviate Docs - generative-databricks]: https://docs.weaviate.io/weaviate/model-providers/databricks/generative
type Generative struct {
	// Endpoint is the Databricks serving endpoint URL.
	Endpoint  string `json:"endpoint,omitempty"`
	MaxTokens *int64 `json:"maxTokens,omitempty"`

	Temperature *float64 `json:"temperature,omitempty"`
	TopP        *float64 `json:"topP,omitempty"`

	// Provider holds query-time options. It is ignored in collection configuration.
	*Provider `json:"-"`
}

func (Generative) Name() string { return "generative-databricks" }

// Provider holds options that apply only when Generative is used as a per-query provider.
type Provider struct {
	// Model name
	Model            string
	N                *int64
	StopSequences    []string
	FrequencyPenalty *float64
	PresencePenalty  *float64
	LogProbs         *bool
	TopLogProbs      *int64
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
		Kind: &proto.GenerativeProvider_Databricks{
			Databricks: &proto.GenerativeDatabricks{
				Endpoint: api.NilZero(g.Endpoint),

				MaxTokens:   g.MaxTokens,
				Temperature: g.Temperature,
				TopP:        g.TopP,

				Model:            api.NilZero(p.Model),
				N:                p.N,
				Stop:             internal.TextArray(p.StopSequences),
				FrequencyPenalty: p.FrequencyPenalty,
				PresencePenalty:  p.PresencePenalty,
				LogProbs:         p.LogProbs,
				TopLogProbs:      p.TopLogProbs,
			},
		},
	}
}
