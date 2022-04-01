package schema

import (
	"context"
	"net/http"

	"github.com/semi-technologies/weaviate-go-client/v3/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/v3/weaviate/except"
	"github.com/semi-technologies/weaviate/entities/models"
)

// ClassCreator builder object to create a schema class
type ClassCreator struct {
	connection *connection.Connection
	class      *models.Class
}

// WithClass specifies the class that will be added to the schema
func (cc *ClassCreator) WithClass(class *models.Class) *ClassCreator {
	cc.class = class
	return cc
}

// WithBM25Config specifies the tuning parameters to be used in BM25 search
func (cc *ClassCreator) WithBM25Config(config *models.BM25Config) *ClassCreator {
	if cc.class.InvertedIndexConfig == nil {
		cc.class.InvertedIndexConfig = &models.InvertedIndexConfig{}
	}

	cc.class.InvertedIndexConfig.Bm25 = config
	return cc
}

// WithStopwordConfig specifies the stopwords to be used during search
func (cc *ClassCreator) WithStopwordConfig(config *models.StopwordConfig) *ClassCreator {
	if cc.class.InvertedIndexConfig == nil {
		cc.class.InvertedIndexConfig = &models.InvertedIndexConfig{}
	}

	cc.class.InvertedIndexConfig.Stopwords = config
	return cc
}

// Do create a class in the schema as specified in the builder
func (cc *ClassCreator) Do(ctx context.Context) error {
	responseData, err := cc.connection.RunREST(ctx, "/schema", http.MethodPost, cc.class)
	return except.CheckResponseDataErrorAndStatusCode(responseData, err, 200)
}
