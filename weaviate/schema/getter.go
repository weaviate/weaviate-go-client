package schema

import (
	"context"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

// Dump Contains all semantic types and respective classes of the schema
type Dump struct {
	models.Schema
}

// Getter builder to get the current schema loaded in weaviate
type Getter struct {
	connection *connection.Connection
}

// Do get and return the weaviate schema
func (sg *Getter) Do(ctx context.Context) (*Dump, error) {
	responseData, err := sg.connection.RunREST(ctx, "/schema", http.MethodGet, nil)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	if responseData.StatusCode == 200 {
		var fullSchema Dump
		decodeErr := responseData.DecodeBodyIntoTarget(&fullSchema)
		return &fullSchema, decodeErr
	}
	return nil, except.NewWeaviateClientError(responseData.StatusCode, string(responseData.Body))
}
