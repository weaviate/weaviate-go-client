package graphql

import (
	pb "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

type Properties struct {
	withProperties []string
	withReferences []*ReferenceProperties
}

func (p *Properties) WithProperties(properties ...string) *Properties {
	p.withProperties = properties
	return p
}

func (p *Properties) WithReferences(references ...*ReferenceProperties) *Properties {
	p.withReferences = references
	return p
}

func (p *Properties) togrpc() *pb.PropertiesRequest {
	props := &pb.PropertiesRequest{
		NonRefProperties:          p.withProperties,
		ReturnAllNonrefProperties: false,
	}
	objectProperties := make([]*pb.ObjectPropertiesRequest, len(p.withProperties))
	for i := range p.withProperties {
		objectProperties[i] = &pb.ObjectPropertiesRequest{PropName: p.withProperties[i]}
	}
	props.ObjectProperties = objectProperties
	if len(p.withReferences) > 0 {
		refProperties := make([]*pb.RefPropertiesRequest, len(p.withReferences))
		for i := range p.withReferences {
			refProperties[i] = p.withReferences[i].togrpc()
		}
		props.RefProperties = refProperties
	}
	// TODO: object properties
	return props
}

type ReferenceProperties struct {
	targetCollection  string
	referenceProperty string
	withProperties    []string
	withMetadata      *Metadata
}

func NewReferenceProperties() *ReferenceProperties {
	return &ReferenceProperties{}
}

func (p *ReferenceProperties) WithTargetCollection(targetCollection string) *ReferenceProperties {
	p.targetCollection = targetCollection
	return p
}

func (p *ReferenceProperties) WithReferenceProperty(referenceProperty string) *ReferenceProperties {
	p.referenceProperty = referenceProperty
	return p
}

func (p *ReferenceProperties) WithProperties(properties ...string) *ReferenceProperties {
	p.withProperties = properties
	return p
}

func (p *ReferenceProperties) WithMetadata(metadata *Metadata) *ReferenceProperties {
	p.withMetadata = metadata
	return p
}

func (p *ReferenceProperties) togrpc() *pb.RefPropertiesRequest {
	refProps := &pb.RefPropertiesRequest{
		TargetCollection:  p.targetCollection,
		ReferenceProperty: p.referenceProperty,
		Metadata:          p.withMetadata.togrpc(),
	}
	if len(p.withProperties) > 0 {
		props := &Properties{withProperties: p.withProperties}
		refProps.Properties = props.togrpc()
	}
	return refProps
}
