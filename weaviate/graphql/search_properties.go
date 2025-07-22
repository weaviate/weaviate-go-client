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

// This method is lacking support for json object properties
// it only supports non ref and reference properties
func (p *Properties) togrpc() *pb.PropertiesRequest {
	props := &pb.PropertiesRequest{
		NonRefProperties:          p.withProperties,
		ReturnAllNonrefProperties: len(p.withProperties) == 0,
	}
	if len(p.withReferences) > 0 {
		refProperties := make([]*pb.RefPropertiesRequest, len(p.withReferences))
		for i := range p.withReferences {
			refProperties[i] = p.withReferences[i].togrpc()
		}
		props.RefProperties = refProperties
	}
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
