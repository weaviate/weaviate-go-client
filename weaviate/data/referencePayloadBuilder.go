package data

import (
	"github.com/go-openapi/strfmt"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/crossref"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/db"
	"github.com/weaviate/weaviate/entities/models"
)

// ReferencePayloadBuilder to create a payload that references a data object.
// The payload may be added to a reference property in another data object.
type ReferencePayloadBuilder struct {
	connection       *connection.Connection
	className        string
	uuid             string
	dbVersionSupport *db.VersionSupport
}

// WithClassName specifies the class name of the object to be referenced
func (rpb *ReferencePayloadBuilder) WithClassName(className string) *ReferencePayloadBuilder {
	rpb.className = className
	return rpb
}

// WithID specifies the uuid of the object to be referenced
func (rpb *ReferencePayloadBuilder) WithID(uuid string) *ReferencePayloadBuilder {
	rpb.uuid = uuid
	return rpb
}

// Payload to reference the in the builder specified data object
func (rpb *ReferencePayloadBuilder) Payload() *models.SingleRef {
	beacon := crossref.BuildBeacon(rpb.uuid, rpb.className, rpb.dbVersionSupport)
	ref := &models.SingleRef{
		Beacon: strfmt.URI(beacon),
	}
	return ref
}
