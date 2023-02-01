package data

import (
	"context"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/except"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/util"
	"github.com/weaviate/weaviate/entities/models"
)

// ReferenceCreator builder to add a reference to the property of a data object
type ReferenceCreator struct {
	connection        *connection.Connection
	className         string
	uuid              string
	referenceProperty string
	referencePayload  *models.SingleRef
	dbVersionSupport  *util.DBVersionSupport
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

// Do add the reference specified by the set payload to the object and property specified in the builder.
func (rc *ReferenceCreator) Do(ctx context.Context) error {
	path := buildReferencesPath(rc.uuid, rc.className, rc.referenceProperty, rc.dbVersionSupport)
	responseData, responseErr := rc.connection.RunREST(ctx, path, http.MethodPost, *rc.referencePayload)
	return except.CheckResponseDataErrorAndStatusCode(responseData, responseErr, 200)
}
