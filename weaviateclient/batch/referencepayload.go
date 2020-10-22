package batch

import (
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/paragons"
	"github.com/semi-technologies/weaviate/entities/models"
)

// ReferencePayloadBuilder to create references that may be added in a batch
type ReferencePayloadBuilder struct {
	connection *connection.Connection
	fromSemanticKind paragons.SemanticKind
	toSemanticKind paragons.SemanticKind
	fromClassName string
	fromPropertyName string
	fromUUID string
	toUUID string
}

func (rpb *ReferencePayloadBuilder) WithFromClassName(className string) *ReferencePayloadBuilder {
	rpb.fromClassName = className
	return rpb
}

func (rpb *ReferencePayloadBuilder) WithFromRefProp(propertyName string) *ReferencePayloadBuilder {
	rpb.fromPropertyName = propertyName
	return rpb
}

func (rpb *ReferencePayloadBuilder) WithFromId(uuid string) *ReferencePayloadBuilder {
	rpb.fromUUID = uuid
	return rpb
}

func (rpb *ReferencePayloadBuilder) WithToId(uuid string) *ReferencePayloadBuilder {
	rpb.toUUID = uuid
	return rpb
}

func (rpb *ReferencePayloadBuilder) WithFromKind(semanticKind paragons.SemanticKind) *ReferencePayloadBuilder {
	rpb.fromSemanticKind = semanticKind
	return rpb
}

func (rpb *ReferencePayloadBuilder) WithToKind(semanticKind paragons.SemanticKind) *ReferencePayloadBuilder {
	rpb.toSemanticKind = semanticKind
	return rpb
}

func (rpb *ReferencePayloadBuilder) Payload() *models.BatchReference {
	from := fmt.Sprintf("weaviate://localhost/%v/%v/%v/%v", string(rpb.fromSemanticKind), rpb.fromClassName, rpb.fromUUID, rpb.fromPropertyName)
	to := fmt.Sprintf("weaviate://localhost/%v/%v", string(rpb.toSemanticKind), rpb.toUUID)

	return &models.BatchReference{
		From: strfmt.URI(from),
		To:   strfmt.URI(to),
	}
}