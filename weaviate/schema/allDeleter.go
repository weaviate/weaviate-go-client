package schema

import (
	"context"
	"github.com/semi-technologies/weaviate-go-client/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviate/except"
	"github.com/semi-technologies/weaviate-go-client/weaviate/paragons"
)

// AllDeleter builder object to delete an entire schema
type AllDeleter struct {
	connection *connection.Connection
	schemaAPI  *API
}

// Do deletes all schema classes from weaviate
func (ad *AllDeleter) Do(ctx context.Context) error {
	schema, getSchemaErr := ad.schemaAPI.Getter().Do(ctx)
	if getSchemaErr != nil {
		return except.NewDerivedWeaviateClientError(getSchemaErr)
	}
	for _, class := range schema.Actions.Classes {
		delErr := ad.schemaAPI.ClassDeleter().WithClassName(class.Class).WithKind(paragons.SemanticKindActions).Do(ctx)
		if delErr != nil {
			return delErr
		}
	}
	for _, class := range schema.Things.Classes {
		delErr := ad.schemaAPI.ClassDeleter().WithClassName(class.Class).WithKind(paragons.SemanticKindThings).Do(ctx)
		if delErr != nil {
			return delErr
		}
	}
	return nil
}
