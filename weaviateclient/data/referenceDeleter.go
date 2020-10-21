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

type ReferenceDelter struct {
	connection *connection.Connection
	semanticKind paragons.SemanticKind
	uuid string
	referenceProperty string
	referencePayload *models.SingleRef
}

func (rr *ReferenceDelter) WithID(uuid string) *ReferenceDelter {
	rr.uuid = uuid
	return rr
}

func (rr *ReferenceDelter) WithReferenceProperty(propertyName string) *ReferenceDelter {
	rr.referenceProperty = propertyName
	return rr
}

func (rr *ReferenceDelter) WithKind(semanticKind paragons.SemanticKind) *ReferenceDelter {
	rr.semanticKind = semanticKind
	return rr
}

func (rr *ReferenceDelter) WithReference(referencePayload *models.SingleRef) *ReferenceDelter {
	rr.referencePayload = referencePayload
	return rr
}

func (rr *ReferenceDelter) Do(ctx context.Context) error {
	path := fmt.Sprintf("/%v/%v/references/%v", string(rr.semanticKind), rr.uuid, rr.referenceProperty)
	responseData, responseErr := rr.connection.RunREST(ctx, path, http.MethodDelete, *rr.referencePayload)
	return clienterrors.CheckResponnseDataErrorAndStatusCode(responseData, responseErr, 204)
}
