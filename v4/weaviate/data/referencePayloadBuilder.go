package data

import (
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/connection"
	"github.com/semi-technologies/weaviate/entities/models"
)

// ReferencePayloadBuilder to create a payload that references a data object.
// The payload may be added to a reference property in another data object.
type ReferencePayloadBuilder struct {
	connection *connection.Connection
	uuid       string
}

// WithID specifies the uuid of the object to be referenced
func (rpb *ReferencePayloadBuilder) WithID(uuid string) *ReferencePayloadBuilder {
	rpb.uuid = uuid
	return rpb
}

// Payload to reference the in the builder specified data object
func (rpb *ReferencePayloadBuilder) Payload() *models.SingleRef {
	beacon := fmt.Sprintf("weaviate://localhost/%v", rpb.uuid)
	ref := &models.SingleRef{
		Beacon: strfmt.URI(beacon),
	}
	return ref
}
