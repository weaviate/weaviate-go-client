package api

import (
	"github.com/google/uuid"
	proto "github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/proto/v1"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
)

type SearchRequest struct {
	RequestDefaults

	Limit            int
	Offset           int
	AutoLimit        int
	After            uuid.UUID
	ReturnProperties []ReturnProperty
	ReturnReferences []ReturnReference
	ReturnVectors    []string
	ReturnMetadata   ReturnMetadata
}

var (
	_ Message[proto.SearchRequest, proto.SearchReply] = (*SearchRequest)(nil)
	_ MessageMarshaler[proto.SearchRequest]           = (*SearchRequest)(nil)
)

func (r *SearchRequest) Method() MethodFunc[proto.SearchRequest, proto.SearchReply] {
	return proto.WeaviateClient.Search
}
func (r *SearchRequest) Body() MessageMarshaler[proto.SearchRequest] { return r }

type (
	ReturnMetadata struct {
		UUID         bool
		CreatedAt    bool
		LastUpdateAt bool
		Distance     bool
		Certainty    bool
		Score        bool
		ExplainScore bool
	}
	ReturnProperty struct {
		Name             string
		NestedProperties []ReturnProperty
	}
	ReturnReference struct {
		PropertyName     string
		TargetCollection string
		ReturnProperties []ReturnProperty
		ReturnReferences []ReturnReference
		ReturnVectors    []string
		ReturnMetadata   ReturnMetadata
	}
)

func (r *SearchRequest) MarshalMessage() (*proto.SearchRequest, error) {
	after := r.After.String()
	if r.After == uuid.Nil {
		after = ""
	}
	req := &proto.SearchRequest{
		Collection:       r.CollectionName,
		Tenant:           r.Tenant,
		ConsistencyLevel: r.ConsistencyLevel.proto(),
		Limit:            uint32(r.Limit),
		Offset:           uint32(r.Offset),
		Autocut:          uint32(r.AutoLimit),
		After:            after,
		Metadata: &proto.MetadataRequest{
			Uuid:               true,
			Distance:           r.ReturnMetadata.Distance,
			Certainty:          r.ReturnMetadata.Certainty,
			CreationTimeUnix:   r.ReturnMetadata.CreatedAt,
			LastUpdateTimeUnix: r.ReturnMetadata.LastUpdateAt,
			Score:              r.ReturnMetadata.Score,
			ExplainScore:       r.ReturnMetadata.ExplainScore,
		},
		Properties: new(proto.PropertiesRequest),
	}

	marshalReturnVectors(req.Metadata, r.ReturnVectors)
	marshalReturnProperties(req.Properties, r.ReturnProperties)
	marshalReturnReferences(req.Properties, r.ReturnReferences)

	return req, nil
}

func marshalReturnProperties(req *proto.PropertiesRequest, rps []ReturnProperty) {
	dev.AssertNotNil(req, "nil req")

	if len(rps) == 0 && rps != nil {
		// ReturnProperties were explicitly set to an empty slice, do not return any.
		return
	}

	if rps == nil {
		// ReturnProperties were not set at all, default to all properties.
		req.ReturnAllNonrefProperties = true
		return
	}

	// walk traverses the ReturnProperty tree and collects requested nested object properties.
	//
	// The reason we cannot recursively call marshalReturnProperties itself is
	// that PropertiesRequest has a different shape from ObjectPropertiesRequest.
	var walk func(*[]*proto.ObjectPropertiesRequest, *ReturnProperty)

	walk = func(os *[]*proto.ObjectPropertiesRequest, rp *ReturnProperty) {
		o := &proto.ObjectPropertiesRequest{PropName: rp.Name}
		for _, np := range rp.NestedProperties {
			if len(np.NestedProperties) == 0 {
				o.PrimitiveProperties = append(o.PrimitiveProperties, np.Name)
			} else {
				walk(&o.ObjectProperties, &np)
			}
		}
		*os = append(*os, o)
	}

	// Add all "primitive" and "nested" object properties to the request.
	for _, rp := range rps {
		if len(rp.NestedProperties) == 0 {
			req.NonRefProperties = append(req.NonRefProperties, rp.Name)
		} else {
			walk(&req.ObjectProperties, &rp)
		}
	}
}

// marshalReturnReferences traverses each ReturnReference tree in the slice
// and collects requested references and properties.
func marshalReturnReferences(req *proto.PropertiesRequest, rrs []ReturnReference) {
	dev.AssertNotNil(req, "nil req")

	for _, rr := range rrs {
		ref := &proto.RefPropertiesRequest{
			ReferenceProperty: rr.PropertyName,
			TargetCollection:  rr.TargetCollection,
			Metadata: &proto.MetadataRequest{
				Uuid:               true,
				CreationTimeUnix:   rr.ReturnMetadata.CreatedAt,
				LastUpdateTimeUnix: rr.ReturnMetadata.LastUpdateAt,
				Vectors:            rr.ReturnVectors,
			},
			Properties: new(proto.PropertiesRequest),
		}

		marshalReturnVectors(ref.Metadata, rr.ReturnVectors)
		marshalReturnProperties(ref.Properties, rr.ReturnProperties)
		marshalReturnReferences(ref.Properties, rr.ReturnReferences)

		req.RefProperties = append(req.RefProperties, ref)
	}
}

func marshalReturnVectors(req *proto.MetadataRequest, vectors []string) {
	if len(vectors) == 0 && vectors != nil {
		// ReturnVectors were explicitly set to an empty slice, include the "only" vector.
		req.Vector = true
		req.Vectors = nil
	} else {
		req.Vectors = vectors
	}
}

// proto converts ConsistencyLevel into a protobuf value.
func (cl ConsistencyLevel) proto() *proto.ConsistencyLevel {
	switch cl {
	case ConsistencyLevelOne:
		return ptr(proto.ConsistencyLevel_CONSISTENCY_LEVEL_ONE)
	case ConsistencyLevelQuorum:
		return ptr(proto.ConsistencyLevel_CONSISTENCY_LEVEL_QUORUM)
	case ConsistencyLevelAll:
		return ptr(proto.ConsistencyLevel_CONSISTENCY_LEVEL_ALL)
	default:
		return nil
	}
}

// ptr is a helper for passing pointers to constants.
func ptr[T any](v T) *T { return &v }
