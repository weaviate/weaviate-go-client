package weaviateclient

import (
	"context"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	clientModels "github.com/semi-technologies/weaviate-go-client/weaviateclient/models"
	"github.com/semi-technologies/weaviate/entities/models"
	"net/http"
)

type SchemaAPI struct {
	connection *connection.Connection
}

func (schema *SchemaAPI) Getter() (*SchemaGetter) {
	return &SchemaGetter{connection: schema.connection}
}

type SchemaGetter struct {
	connection *connection.Connection
}

func (sg *SchemaGetter) Do (ctx context.Context) (*clientModels.SchemaDump, error){
	responseData, err := sg.connection.RunREST(ctx, "/schema", http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	if responseData.StatusCode == 200 {
		var fullSchema clientModels.SchemaDump
		decodeErr := responseData.DecodeBodyIntoTarget(&fullSchema)
		if decodeErr != nil {
			return nil, decodeErr
		}
		return &fullSchema, nil
	}
	return nil, NewUnexpectedStatusCodeError(responseData.StatusCode, string(responseData.Body))
}

func (schema *SchemaAPI) ClassCreator() *ClassCreator {
	return &ClassCreator{
		connection: schema.connection,
		semanticKind: SemanticKindThings, // Set the default
	}
}

//.classCreator()
//.withClass(classObj)
//.withKind(weaviate.KIND_THINGS)
//.do()

type ClassCreator struct {
	connection *connection.Connection
	class *models.Class
	semanticKind SemanticKind
}

func (cc *ClassCreator) WithClass(class *models.Class) *ClassCreator {
	cc.class = class
	return cc
}

func (cc *ClassCreator) WithKind(semanticKind SemanticKind) *ClassCreator {
	cc.semanticKind = semanticKind
	return cc
}

func (cc *ClassCreator) Do(ctx context.Context) error {
	path := "/schema/"+string(cc.semanticKind)
	responseData, err := cc.connection.RunREST(ctx, path, http.MethodPost, cc.class)
	if err != nil {
		return err
	}
	if responseData.StatusCode == 200 {
		return nil
	}
	return NewUnexpectedStatusCodeError(responseData.StatusCode, string(responseData.Body))
}