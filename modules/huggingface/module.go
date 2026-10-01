package huggingface

import (
	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/modules"
)

func init() {
	modules.Register(*new(Text2Vec))
}

// Text2Vec is a vectorizer for text properties based on the text2vec-huggingface module.
// Unset fields inherit the server defaults.
//
// Model and PassageModel are mutually exclusive; the server rejects a config that sets both.
//
// See [Weaviate Docs - text2vec-huggingface].
//
// [Weaviate Docs - text2vec-huggingface]: https://docs.weaviate.io/weaviate/model-providers/huggingface/embeddings
type Text2Vec struct {
	// Model is the Hugging Face model ID.
	Model string
	// PassageModel is an alias the server reads when Model is unset.
	// Prefer Model; this field keeps configs created by other clients intact.
	PassageModel string
	// EndpointURL points to a dedicated inference endpoint; when set, the server skips model checks.
	EndpointURL string
	// WaitForModel waits for the model to be loaded.
	WaitForModel *bool
	// UseGPU runs inference on a GPU.
	UseGPU *bool
	// UseCache enables the Inference API cache.
	UseCache *bool
	// Properties limits vectorization to these properties.
	// By default, all text properties are vectorized.
	Properties []string
}

func (Text2Vec) Name() string { return "text2vec-huggingface" }

var (
	_ internal.MapEncoder = (*Text2Vec)(nil)
	_ internal.MapDecoder = (*Text2Vec)(nil)
)

// text2VecJSON is the wire shape of [Text2Vec]: the inference flags are nested under "options".
type text2VecJSON struct {
	Model        string       `json:"model,omitempty"`
	PassageModel string       `json:"passageModel,omitempty"`
	EndpointURL  string       `json:"endpointURL,omitempty"`
	Options      *optionsJSON `json:"options,omitempty"`
	Properties   []string     `json:"properties,omitempty"`
}

type optionsJSON struct {
	WaitForModel *bool `json:"waitForModel,omitempty"`
	UseGPU       *bool `json:"useGPU,omitempty"`
	UseCache     *bool `json:"useCache,omitempty"`
}

func (t2v Text2Vec) EncodeMap() (map[string]any, error) {
	var options *optionsJSON
	if t2v.WaitForModel != nil || t2v.UseGPU != nil || t2v.UseCache != nil {
		options = &optionsJSON{
			WaitForModel: t2v.WaitForModel,
			UseGPU:       t2v.UseGPU,
			UseCache:     t2v.UseCache,
		}
	}

	dest := make(map[string]any)
	if err := internal.Encode(text2VecJSON{
		Model:        t2v.Model,
		PassageModel: t2v.PassageModel,
		EndpointURL:  t2v.EndpointURL,
		Options:      options,
		Properties:   t2v.Properties,
	}, dest); err != nil {
		return nil, err
	}
	return dest, nil
}

func (t2v *Text2Vec) DecodeMap(m map[string]any) error {
	var dest text2VecJSON
	if err := internal.Decode(m, &dest); err != nil {
		return err
	}

	*t2v = Text2Vec{
		Model:        dest.Model,
		PassageModel: dest.PassageModel,
		EndpointURL:  dest.EndpointURL,
		Properties:   dest.Properties,
	}
	if opt := dest.Options; opt != nil {
		t2v.WaitForModel = opt.WaitForModel
		t2v.UseGPU = opt.UseGPU
		t2v.UseCache = opt.UseCache
	}
	return nil
}
