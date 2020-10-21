package data

import (
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/paragons"
	"github.com/semi-technologies/weaviate/entities/models"
)

type ReferencePayloadBuilder struct {
	connection *connection.Connection
	uuid string
	semanticKind paragons.SemanticKind
}

func (rpb *ReferencePayloadBuilder) WithID(uuid string) *ReferencePayloadBuilder {
	rpb.uuid = uuid
	return rpb
}

func (rpb *ReferencePayloadBuilder) WithKind(semanticKind paragons.SemanticKind) *ReferencePayloadBuilder {
	rpb.semanticKind = semanticKind
	return rpb
}

func (rpb *ReferencePayloadBuilder) Payload() *models.SingleRef {
	beacon := fmt.Sprintf("weaviate://localhost/%v/%v", string(rpb.semanticKind), rpb.uuid)
	ref := &models.SingleRef{
		Beacon: strfmt.URI(beacon),
	}
	return ref
}