package data

import (
	"context"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/db"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/pathbuilder"
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
	tenant            string
	dbVersionSupport  *db.VersionSupport
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

// WithTenant specifies tenant of referenced objects
func (rr *ReferenceReplacer) WithTenant(tenant string) *ReferenceReplacer {
	rr.tenant = tenant
	return rr
}

// Do replace the references of the in this builder specified data object
func (rr *ReferenceReplacer) Do(ctx context.Context) error {
	path := pathbuilder.References(pathbuilder.Components{
		ID:                rr.uuid,
		Class:             rr.className,
		DBVersion:         rr.dbVersionSupport,
		ReferenceProperty: rr.referenceProperty,
		ConsistencyLevel:  rr.consistencyLevel,
		Tenant:            rr.tenant,
	})
	responseData, responseErr := rr.connection.RunREST(ctx, path, http.MethodPut, *rr.referencePayload)
	return except.CheckResponseDataErrorAndStatusCode(responseData, responseErr, 200)
}
