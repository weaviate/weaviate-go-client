package schema

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/db"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/internal"
	"github.com/weaviate/weaviate/entities/models"
)

// ClassUpdater builder object to update a schema class
type ClassUpdater struct {
	connection        *connection.Connection
	dbVersionProvider *db.VersionProvider
	class             *models.Class
}

// WithClass specifies the class properties that will be added to the schema
func (cu *ClassUpdater) WithClass(class *models.Class) *ClassUpdater {
	cu.class = class
	return cu
}

// Do create a class in the schema as specified in the builder
func (cu *ClassUpdater) Do(ctx context.Context) error {
	if cu.class == nil || cu.class.Class == "" {
		return except.NewWeaviateClientError(0, "A class must be provided")
	}
	if err := internal.CheckTextAnalyzerSupport(cu.dbVersionProvider, cu.class); err != nil {
		return err
	}
	path := fmt.Sprintf("/schema/%v", cu.class.Class)
	responseData, err := cu.connection.RunREST(ctx, path, http.MethodPut, cu.class)
	return except.CheckResponseDataErrorAndStatusCode(responseData, err, 200)
}
