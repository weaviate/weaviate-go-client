package databricks

import "github.com/weaviate/weaviate-go-client/v6/modules"

func init() {
	modules.Register(*new(Text2Vec))
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
