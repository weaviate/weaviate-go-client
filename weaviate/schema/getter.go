package schema

import (
	"context"
	"github.com/semi-technologies/weaviate-go-client/weaviate/except"
	"github.com/semi-technologies/weaviate-go-client/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviate/paragons"
	"net/http"
)

// Getter builder to get the current schema loaded in weaviate
type Getter struct {
	connection *connection.Connection
}

// Do get and return the weaviate schema
func (sg *Getter) Do(ctx context.Context) (*paragons.SchemaDump, error) {
	responseData, err := sg.connection.RunREST(ctx, "/schema", http.MethodGet, nil)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	if responseData.StatusCode == 200 {
		var fullSchema paragons.SchemaDump
		decodeErr := responseData.DecodeBodyIntoTarget(&fullSchema)
		return &fullSchema, decodeErr
	}
	return nil, except.NewWeaviateClientError(responseData.StatusCode, string(responseData.Body))
}
