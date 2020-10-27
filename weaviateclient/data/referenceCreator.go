package data

import (
	"context"
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/clienterrors"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/paragons"
	"github.com/semi-technologies/weaviate/entities/models"
	"net/http"
)

// ReferenceCreator builder to add a reference to the property of a data object
type ReferenceCreator struct {
	connection        *connection.Connection
	semanticKind      paragons.SemanticKind
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

// WithKind specifies the semantic kind that is used for the data object
// If not called the builder defaults to `things`
func (rc *ReferenceCreator) WithKind(semanticKind paragons.SemanticKind) *ReferenceCreator {
	rc.semanticKind = semanticKind
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
	path := fmt.Sprintf("/%v/%v/references/%v", string(rc.semanticKind), rc.uuid, rc.referenceProperty)
	responseData, responseErr := rc.connection.RunREST(ctx, path, http.MethodPost, *rc.referencePayload)
	return clienterrors.CheckResponnseDataErrorAndStatusCode(responseData, responseErr, 200)
}
