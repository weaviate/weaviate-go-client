package model2vec

import "github.com/weaviate/weaviate-go-client/v6/modules"

func init() {
	modules.Register(*new(Text2Vec))
}

// Text2Vec is a vectorizer for text properties
// based on the text2vec-model2vec module.
type Text2Vec struct {
	URL        string   `json:"inferenceURL"`
	Properties []string `json:"properties"`
}

func (Text2Vec) Name() string { return "text2vec-model2vec" }
