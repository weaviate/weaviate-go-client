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

// ReferenceDeleter builder to remove a reference from a data object
type ReferenceDeleter struct {
	connection        *connection.Connection
	className         string
	uuid              string
	referenceProperty string
	referencePayload  *models.SingleRef
	consistencyLevel  string
	tenant            string
	dbVersionSupport  *db.VersionSupport
}

// WithClassName specifies the class name of the object on which the reference will be deleted
func (rd *ReferenceDeleter) WithClassName(className string) *ReferenceDeleter {
	rd.className = className
	return rd
}

// WithID specifies the uuid of the object on which the reference will be deleted
func (rd *ReferenceDeleter) WithID(uuid string) *ReferenceDeleter {
	rd.uuid = uuid
	return rd
}

// WithReferenceProperty specifies the property on which the reference should be deleted
func (rd *ReferenceDeleter) WithReferenceProperty(propertyName string) *ReferenceDeleter {
	rd.referenceProperty = propertyName
	return rd
}

// WithReference specifies reference payload of the reference about to be deleted
func (rd *ReferenceDeleter) WithReference(referencePayload *models.SingleRef) *ReferenceDeleter {
	rd.referencePayload = referencePayload
	return rd
}

// WithConsistencyLevel determines how many replicas must acknowledge a request
// before it is considered successful. Mutually exclusive with node_name param.
// Can be one of 'ALL', 'ONE', or 'QUORUM'.
func (rd *ReferenceDeleter) WithConsistencyLevel(cl string) *ReferenceDeleter {
	rd.consistencyLevel = cl
	return rd
}

// WithTenant specifies tenant of referenced objects
func (rd *ReferenceDeleter) WithTenant(tenant string) *ReferenceDeleter {
	rd.tenant = tenant
	return rd
}

// Do remove the reference defined by the payload set in this builder to the property and object defined in this builder
func (rd *ReferenceDeleter) Do(ctx context.Context) error {
	path := pathbuilder.References(pathbuilder.Components{
		ID:                rd.uuid,
		Class:             rd.className,
		DBVersion:         rd.dbVersionSupport,
		ReferenceProperty: rd.referenceProperty,
		ConsistencyLevel:  rd.consistencyLevel,
		Tenant:            rd.tenant,
	})
	responseData, responseErr := rd.connection.RunREST(ctx, path, http.MethodDelete, *rd.referencePayload)
	return except.CheckResponseDataErrorAndStatusCode(responseData, responseErr, 204)
}
