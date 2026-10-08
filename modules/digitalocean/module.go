package digitalocean

import "github.com/weaviate/weaviate-go-client/v6/modules"

func init() {
	modules.Register(*new(Text2Vec))
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
