package nvidia

import (
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/modules"
	"github.com/weaviate/weaviate-go-client/v6/modules/internal"
	proto "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

func init() {
	modules.Register(*new(Text2Vec))
	modules.Register(*new(Multi2Vec))
	modules.Register(*new(Generative))
}

// Text2Vec is a vectorizer for text properties based on the text2vec-nvidia module.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - text2vec-nvidia].
//
// [Weaviate Docs - text2vec-nvidia]: https://docs.weaviate.io/weaviate/model-providers/nvidia/embeddings
type Text2Vec struct {
	// BaseURL overrides the default request URL, e.g. a proxy.
	BaseURL string `json:"baseURL,omitempty"`
	// Model is the embedding model name.
	Model string `json:"model,omitempty"`
	// Properties limits vectorization to these properties.
	// By default, all text properties are vectorized.
	Properties []string `json:"properties,omitempty"`
	// Truncate sets how inputs longer than the model's context are cut.
	Truncate string `json:"truncate,omitempty"`
}

func (Text2Vec) Name() string { return "text2vec-nvidia" }

const (
	TruncateNone  = "NONE"
	TruncateStart = "START"
	TruncateEnd   = "END"
)

// Multi2Vec is a vectorizer for text and image properties based on the multi2vec-nvidia module.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - multi2vec-nvidia].
//
// [Weaviate Docs - multi2vec-nvidia]: https://docs.weaviate.io/weaviate/model-providers/nvidia/embeddings-multimodal
type Multi2Vec struct {
	// BaseURL overrides the default request URL, e.g. a proxy.
	BaseURL string
	// Model is the embedding model name.
	Model string
	// TextFields and ImageFields name the properties to vectorize.
	TextFields  []modules.WeightedField `json:"-"`
	ImageFields []modules.WeightedField `json:"-"`
}

func (Multi2Vec) Name() string { return "multi2vec-nvidia" }

type multi2vecJSON struct {
	BaseURL string `json:"baseURL,omitempty"`
	Model   string `json:"model,omitempty"`
}

func (m2v Multi2Vec) EncodeMap() (map[string]any, error) {
	return internal.EncodeWithWeightedFields(multi2vecJSON{
		BaseURL: m2v.BaseURL,
		Model:   m2v.Model,
	}, internal.WeightedFields{
		TextFields:  m2v.TextFields,
		ImageFields: m2v.ImageFields,
	})
}

func (m2v *Multi2Vec) DecodeMap(m map[string]any) error {
	var dest multi2vecJSON
	wfs, err := internal.DecodeWithWeightedFields(m, &dest)
	if err != nil {
		return err
	}

	*m2v = Multi2Vec{
		BaseURL:     dest.BaseURL,
		Model:       dest.Model,
		TextFields:  wfs.TextFields,
		ImageFields: wfs.ImageFields,
	}
	return nil
}

// Generative configures the generative-nvidia module, either as a collection's
// default generative module or as a per-query provider in the generate package.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - generative-nvidia].
//
// [Weaviate Docs - generative-nvidia]: https://docs.weaviate.io/weaviate/model-providers/nvidia/generative
type Generative struct {
	// BaseURL overrides the default request URL, e.g. a proxy.
	BaseURL string `json:"baseURL,omitempty"`
	// Model name
	Model     string `json:"model,omitempty"`
	MaxTokens *int64 `json:"maxTokens,omitempty"`

	Temperature *float64 `json:"temperature,omitempty"`
	TopP        *float64 `json:"topP,omitempty"`

	// Provider holds query-time options. It is ignored in collection configuration.
	*Provider `json:"-"`
}

func (Generative) Name() string { return "generative-nvidia" }

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
		Kind: &proto.GenerativeProvider_Nvidia{
			Nvidia: &proto.GenerativeNvidia{
				BaseUrl: api.NilZero(g.BaseURL),
				Model:   api.NilZero(g.Model),

				MaxTokens:   g.MaxTokens,
				Temperature: g.Temperature,
				TopP:        g.TopP,
			},
		},
	}
}
