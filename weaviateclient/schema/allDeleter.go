package schema

import (
	"context"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/models"
)

type AllDeleter struct {
	connection *connection.Connection
	schemaAPI *SchemaAPI
}

func (ad *AllDeleter) Do(ctx context.Context) error {
	schema, getSchemaErr := ad.schemaAPI.Getter().Do(ctx)
	if getSchemaErr != nil {
		return getSchemaErr
	}
	for _, class := range schema.Actions.Classes {
		delErr := ad.schemaAPI.ClassDeleter().WithClassName(class.Class).WithKind(models.SemanticKindActions).Do(ctx)
		if delErr != nil {
			return delErr
		}
	}
	for _, class := range schema.Things.Classes {
		delErr := ad.schemaAPI.ClassDeleter().WithClassName(class.Class).WithKind(models.SemanticKindThings).Do(ctx)
		if delErr != nil {
			return delErr
		}
	}
	return nil
}