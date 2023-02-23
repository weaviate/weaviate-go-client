package data

import (
	"context"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/db"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/except"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/pathbuilder"
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
	dbVersionSupport  *db.VersionSupport
}

// WithClassName specifies the class name of the object on which the reference will be deleted
func (rr *ReferenceDeleter) WithClassName(className string) *ReferenceDeleter {
	rr.className = className
	return rr
}

// WithID specifies the uuid of the object on which the reference will be deleted
func (rr *ReferenceDeleter) WithID(uuid string) *ReferenceDeleter {
	rr.uuid = uuid
	return rr
}

// WithReferenceProperty specifies the property on which the reference should be deleted
func (rr *ReferenceDeleter) WithReferenceProperty(propertyName string) *ReferenceDeleter {
	rr.referenceProperty = propertyName
	return rr
}

// WithReference specifies reference payload of the reference about to be deleted
func (rr *ReferenceDeleter) WithReference(referencePayload *models.SingleRef) *ReferenceDeleter {
	rr.referencePayload = referencePayload
	return rr
}

// WithConsistencyLevel determines how many replicas must acknowledge a request
// before it is considered successful. Mutually exclusive with node_name param.
// Can be one of 'ALL', 'ONE', or 'QUORUM'.
func (rr *ReferenceDeleter) WithConsistencyLevel(cl string) *ReferenceDeleter {
	rr.consistencyLevel = cl
	return rr
}

// Do remove the reference defined by the payload set in this builder to the property and object defined in this builder
func (rr *ReferenceDeleter) Do(ctx context.Context) error {
	path := pathbuilder.References(pathbuilder.Components{
		ID:                rr.uuid,
		Class:             rr.className,
		DBVersion:         rr.dbVersionSupport,
		ReferenceProperty: rr.referenceProperty,
		ConsistencyLevel:  rr.consistencyLevel,
	})
	responseData, responseErr := rr.connection.RunREST(ctx, path, http.MethodDelete, *rr.referencePayload)
	return except.CheckResponseDataErrorAndStatusCode(responseData, responseErr, 204)
}
