package selfprovided

import "github.com/weaviate/weaviate-go-client/v6/modules"

// Vectorizer does not configure any vectorization module and instead
// allows the user to supply custom vectors alongisde objects' properties.
var Vectorizer = *new(vectorizer)

func init() {
	modules.Register(Vectorizer)
}

type vectorizer struct{}

func (vectorizer) Name() string { return "none" }
