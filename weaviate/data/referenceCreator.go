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

// ReferenceCreator builder to add a reference to the property of a data object
type ReferenceCreator struct {
	connection        *connection.Connection
	className         string
	uuid              string
	referenceProperty string
	referencePayload  *models.SingleRef
	consistencyLevel  string
	tenant            string
	dbVersionSupport  *db.VersionSupport
}

// WithClassName specifies the class name of the object on which to add the reference
func (rc *ReferenceCreator) WithClassName(className string) *ReferenceCreator {
	rc.className = className
	return rc
}

// WithID specifies the uuid of the object on which to add the reference
func (rc *ReferenceCreator) WithID(uuid string) *ReferenceCreator {
	rc.uuid = uuid
	return rc
}

// WithTenant specifies tenant of referenced objects
func (rc *ReferenceCreator) WithTenant(tenant string) *ReferenceCreator {
	rc.tenant = tenant
	return rc
}

// WithReferenceProperty specifies the property that should hold the reference
func (rc *ReferenceCreator) WithReferenceProperty(propertyName string) *ReferenceCreator {
	rc.referenceProperty = propertyName
	return rc
}

// WithReference specifies the data object that should be referenced by the in this object specified reference property
// The payload may be created using the ReferencePayloadBuilder
func (rc *ReferenceCreator) WithReference(referencePayload *models.SingleRef) *ReferenceCreator {
	rc.referencePayload = referencePayload
	return rc
}

// WithConsistencyLevel determines how many replicas must acknowledge a request
// before it is considered successful. Mutually exclusive with node_name param.
// Can be one of 'ALL', 'ONE', or 'QUORUM'.
func (rc *ReferenceCreator) WithConsistencyLevel(cl string) *ReferenceCreator {
	rc.consistencyLevel = cl
	return rc
}

// Do add the reference specified by the set payload to the object and property specified in the builder.
func (rc *ReferenceCreator) Do(ctx context.Context) error {
	path := pathbuilder.References(pathbuilder.Components{
		ID:                rc.uuid,
		Class:             rc.className,
		ReferenceProperty: rc.referenceProperty,
		DBVersion:         rc.dbVersionSupport,
		ConsistencyLevel:  rc.consistencyLevel,
		Tenant:            rc.tenant,
	})
	responseData, responseErr := rc.connection.RunREST(ctx, path, http.MethodPost, *rc.referencePayload)
	return except.CheckResponseDataErrorAndStatusCode(responseData, responseErr, 200)
}
