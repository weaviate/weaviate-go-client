package data

import (
	"context"
	"net/http"

	"github.com/semi-technologies/weaviate-go-client/v5/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/v5/weaviate/except"
	"github.com/semi-technologies/weaviate-go-client/v5/weaviate/util"
	"github.com/semi-technologies/weaviate/entities/models"
)

// ReferenceReplacer builder to replace reference(s) with new one(s)
type ReferenceReplacer struct {
	connection        *connection.Connection
	className         string
	uuid              string
	referenceProperty string
	referencePayload  *models.MultipleRef
	dbVersionSupport  *util.DBVersionSupport
}

// WithClassName specifies the class name of the object about to get its reference replaced
func (rr *ReferenceReplacer) WithClassName(className string) *ReferenceReplacer {
	rr.className = className
	return rr
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
	path := buildReferencesPath(rr.uuid, rr.className, rr.referenceProperty, rr.dbVersionSupport)
	responseData, responseErr := rr.connection.RunREST(ctx, path, http.MethodPut, *rr.referencePayload)
	return except.CheckResponseDataErrorAndStatusCode(responseData, responseErr, 200)
}
