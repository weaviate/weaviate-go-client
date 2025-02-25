package batch

import (
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/crossref"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/db"
	"github.com/weaviate/weaviate/entities/models"
)

// ReferencePayloadBuilder to create references that may be added in a batch
type ReferencePayloadBuilder struct {
	connection       *connection.Connection
	fromClassName    string
	fromPropertyName string
	fromUUID         string
	toClassName      string
	toUUID           string
	tenant           string
	dbVersion        *db.VersionSupport
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

// WithToClassName class name of the referenced object
func (rpb *ReferencePayloadBuilder) WithToClassName(className string) *ReferencePayloadBuilder {
	rpb.toClassName = className
	return rpb
}

// WithToID UUID of the referenced object
func (rpb *ReferencePayloadBuilder) WithToID(uuid string) *ReferencePayloadBuilder {
	rpb.toUUID = uuid
	return rpb
}

// WithTenant specifies tenant of referenced objects
func (rpb *ReferencePayloadBuilder) WithTenant(tenant string) *ReferencePayloadBuilder {
	rpb.tenant = tenant
	return rpb
}

// Payload to be used in a batch request
func (rpb *ReferencePayloadBuilder) Payload() *models.BatchReference {
	from := fmt.Sprintf("weaviate://localhost/%v/%v/%v", rpb.fromClassName, rpb.fromUUID, rpb.fromPropertyName)
	to := crossref.BuildBeacon(rpb.toUUID, rpb.toClassName, rpb.dbVersion)

	return &models.BatchReference{
		From:   strfmt.URI(from),
		To:     strfmt.URI(to),
		Tenant: rpb.tenant,
	}
}
