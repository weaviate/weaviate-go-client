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

type ReferenceReplacer struct {
	connection *connection.Connection
	semanticKind paragons.SemanticKind
	uuid string
	referenceProperty string
	referencePayload *models.MultipleRef
}

func (rr *ReferenceReplacer) WithID(uuid string) *ReferenceReplacer {
	rr.uuid = uuid
	return rr
}

func (rr *ReferenceReplacer) WithReferenceProperty(propertyName string) *ReferenceReplacer {
	rr.referenceProperty = propertyName
	return rr
}

func (rr *ReferenceReplacer) WithKind(semanticKind paragons.SemanticKind) *ReferenceReplacer {
	rr.semanticKind = semanticKind
	return rr
}

func (rr *ReferenceReplacer) WithReferences(referencePayload *models.MultipleRef) *ReferenceReplacer {
	rr.referencePayload = referencePayload
	return rr
}

func (rr *ReferenceReplacer) Do(ctx context.Context) error {
	path := fmt.Sprintf("/%v/%v/references/%v", string(rr.semanticKind), rr.uuid, rr.referenceProperty)
	responseData, responseErr := rr.connection.RunREST(ctx, path, http.MethodPut, *rr.referencePayload)
	return clienterrors.CheckResponnseDataErrorAndStatusCode(responseData, responseErr, 200)
}
