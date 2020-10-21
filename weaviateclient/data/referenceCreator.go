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

type ReferenceCreator struct {
	connection *connection.Connection
	semanticKind paragons.SemanticKind
	uuid string
	referenceProperty string
	referencePayload *models.SingleRef
}

func (rc *ReferenceCreator) WithID(uuid string) *ReferenceCreator {
	rc.uuid = uuid
	return rc
}

func (rc *ReferenceCreator) WithReferenceProperty(propertyName string) *ReferenceCreator {
	rc.referenceProperty = propertyName
	return rc
}

func (rc *ReferenceCreator) WithKind(semanticKind paragons.SemanticKind) *ReferenceCreator {
	rc.semanticKind = semanticKind
	return rc
}

func (rc *ReferenceCreator) WithReference(referencePayload *models.SingleRef) *ReferenceCreator {
	rc.referencePayload = referencePayload
	return rc
}

func (rc *ReferenceCreator) Do(ctx context.Context) error {
	path := fmt.Sprintf("/%v/%v/references/%v", string(rc.semanticKind), rc.uuid, rc.referenceProperty)
	var responseData *connection.ResponseData
	var responseErr error
	if rc.semanticKind == paragons.SemanticKindThings {
		responseData, responseErr = rc.connection.RunREST(ctx, path, http.MethodPost, *rc.referencePayload)
	} else {
		responseData, responseErr = rc.connection.RunREST(ctx, path, http.MethodPost, *rc.referencePayload)
	}
	return clienterrors.CheckResponnseDataErrorAndStatusCode(responseData, responseErr, 200)
}

