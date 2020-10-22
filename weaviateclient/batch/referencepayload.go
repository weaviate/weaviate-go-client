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

// WithFromClassName name of the class that the reference is added to
func (rpb *ReferencePayloadBuilder) WithFromClassName(className string) *ReferencePayloadBuilder {
	rpb.fromClassName = className
	return rpb
}

// WithFromRefProp name of the property that the reference is added to
func (rpb *ReferencePayloadBuilder) WithFromRefProp(propertyName string) *ReferencePayloadBuilder {
	rpb.fromPropertyName = propertyName
	return rpb
}

// WithFromID UUID of the object that the reference is added to
func (rpb *ReferencePayloadBuilder) WithFromID(uuid string) *ReferencePayloadBuilder {
	rpb.fromUUID = uuid
	return rpb
}

// WithToID UUID of the referenced object
func (rpb *ReferencePayloadBuilder) WithToID(uuid string) *ReferencePayloadBuilder {
	rpb.toUUID = uuid
	return rpb
}

// WithFromKind semantic kind of the object that the reference is added to
func (rpb *ReferencePayloadBuilder) WithFromKind(semanticKind paragons.SemanticKind) *ReferencePayloadBuilder {
	rpb.fromSemanticKind = semanticKind
	return rpb
}

// WithToKind semantic kind of the referenced object
func (rpb *ReferencePayloadBuilder) WithToKind(semanticKind paragons.SemanticKind) *ReferencePayloadBuilder {
	rpb.toSemanticKind = semanticKind
	return rpb
}

// Payload to be used in a batch request
func (rpb *ReferencePayloadBuilder) Payload() *models.BatchReference {
	from := fmt.Sprintf("weaviate://localhost/%v/%v/%v/%v", string(rpb.fromSemanticKind), rpb.fromClassName, rpb.fromUUID, rpb.fromPropertyName)
	to := fmt.Sprintf("weaviate://localhost/%v/%v", string(rpb.toSemanticKind), rpb.toUUID)

	return &models.BatchReference{
		From: strfmt.URI(from),
		To:   strfmt.URI(to),
	}
}