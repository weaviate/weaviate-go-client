package data

import (
	"context"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/except"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/util"
	"github.com/weaviate/weaviate/entities/models"
)

// ReferenceReplacer builder to replace reference(s) with new one(s)
type ReferenceReplacer struct {
	connection        *connection.Connection
	className         string
	uuid              string
	referenceProperty string
	referencePayload  *models.MultipleRef
	consistencyLevel  string
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

// WithConsistencyLevel determines how many replicas must acknowledge a request
// before it is considered successful. Mutually exclusive with node_name param.
// Can be one of 'ALL', 'ONE', or 'QUORUM'.
func (rr *ReferenceReplacer) WithConsistencyLevel(cl string) *ReferenceReplacer {
	rr.consistencyLevel = cl
	return rr
}

// Do replace the references of the in this builder specified data object
func (rr *ReferenceReplacer) Do(ctx context.Context) error {
	path := buildReferencesPath(pathComponents{
		id:                rr.uuid,
		class:             rr.className,
		dbVersion:         rr.dbVersionSupport,
		referenceProperty: rr.referenceProperty,
		consistencyLevel:  rr.consistencyLevel,
	})
	responseData, responseErr := rr.connection.RunREST(ctx, path, http.MethodPut, *rr.referencePayload)
	return except.CheckResponseDataErrorAndStatusCode(responseData, responseErr, 200)
}
