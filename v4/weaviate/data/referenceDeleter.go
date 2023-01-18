package data

import (
	"context"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/except"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/util"
	"github.com/weaviate/weaviate/entities/models"
)

// ReferenceDeleter builder to remove a reference from a data object
type ReferenceDeleter struct {
	connection        *connection.Connection
	className         string
	uuid              string
	referenceProperty string
	referencePayload  *models.SingleRef
	dbVersionSupport  *util.DBVersionSupport
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

// Do remove the reference defined by the payload set in this builder to the property and object defined in this builder
func (rr *ReferenceDeleter) Do(ctx context.Context) error {
	path := buildReferencesPath(rr.uuid, rr.className, rr.referenceProperty, rr.dbVersionSupport)
	responseData, responseErr := rr.connection.RunREST(ctx, path, http.MethodDelete, *rr.referencePayload)
	return except.CheckResponseDataErrorAndStatusCode(responseData, responseErr, 204)
}
