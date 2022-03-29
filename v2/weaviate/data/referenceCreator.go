package data

import (
	"context"
	"fmt"
	"net/http"

	"github.com/semi-technologies/weaviate-go-client/v2/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/v2/weaviate/except"
	"github.com/semi-technologies/weaviate/entities/models"
)

// ReferenceCreator builder to add a reference to the property of a data object
type ReferenceCreator struct {
	connection        *connection.Connection
	uuid              string
	referenceProperty string
	referencePayload  *models.SingleRef
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
	path := fmt.Sprintf("/objects/%v/references/%v", rc.uuid, rc.referenceProperty)
	responseData, responseErr := rc.connection.RunREST(ctx, path, http.MethodPost, *rc.referencePayload)
	return except.CheckResponnseDataErrorAndStatusCode(responseData, responseErr, 200)
}
