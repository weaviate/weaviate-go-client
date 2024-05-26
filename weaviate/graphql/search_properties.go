package graphql

import (
	pb "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

type Properties struct {
	withProperties []string
	withReferences []*Reference
}

func (p *Properties) WithProperties(properties ...string) *Properties {
	p.withProperties = properties
	return p
}

func (p *Properties) WithReferences(references ...*Reference) *Properties {
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
	// TODO: add support for json object properties
	return props
}

type Reference struct {
	TargetCollection  string
	ReferenceProperty string
	Properties        []string
	Metadata          *Metadata
}

func (p *Reference) togrpc() *pb.RefPropertiesRequest {
	refProps := &pb.RefPropertiesRequest{
		TargetCollection:  p.TargetCollection,
		ReferenceProperty: p.ReferenceProperty,
	}
	if len(p.Properties) > 0 {
		props := &Properties{withProperties: p.Properties}
		refProps.Properties = props.togrpc()
	}
	if p.Metadata != nil {
		refProps.Metadata = p.Metadata.togrpc()
	}
	return refProps
}
