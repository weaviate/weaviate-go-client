package data

import (
	"context"
	"fmt"
	"net/http"

	"github.com/semi-technologies/weaviate-go-client/v2/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/v2/weaviate/except"
	"github.com/semi-technologies/weaviate/entities/models"
)

// ReferenceReplacer builder to replace reference(s) with new one(s)
type ReferenceReplacer struct {
	connection        *connection.Connection
	uuid              string
	referenceProperty string
	referencePayload  *models.MultipleRef
}

// WithID specifies the uuid of the object about to get its reference replaced
func (rr *ReferenceReplacer) WithID(uuid string) *ReferenceReplacer {
	rr.uuid = uuid
	return rr
}

// WithReferenceProperty specifies the property that should replace
func (rr *ReferenceReplacer) WithReferenceProperty(propertyName string) *ReferenceReplacer {
	rr.referenceProperty = propertyName
	return rr
}

// WithReferences the set of references that should replace the currently existing references
func (rr *ReferenceReplacer) WithReferences(referencePayload *models.MultipleRef) *ReferenceReplacer {
	rr.referencePayload = referencePayload
	return rr
}

// Do replace the references of the in this builder specified data object
func (rr *ReferenceReplacer) Do(ctx context.Context) error {
	path := fmt.Sprintf("/objects/%v/references/%v", rr.uuid, rr.referenceProperty)
	responseData, responseErr := rr.connection.RunREST(ctx, path, http.MethodPut, *rr.referencePayload)
	return except.CheckResponnseDataErrorAndStatusCode(responseData, responseErr, 200)
}
