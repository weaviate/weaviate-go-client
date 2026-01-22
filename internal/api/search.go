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
		ReturnReference  []ReturnReference
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
			Vectors:            r.ReturnVectors,
		},
	}

	// ReturnVectors were explicitly set to an empty slice, include the "only" vector.
	if len(r.ReturnVectors) == 0 && r.ReturnVectors != nil {
		req.Metadata.Vector = true
		req.Metadata.Vectors = nil
	}

	marshalReturnProperties(&req.Properties, r.ReturnProperties)

	for _, rr := range r.ReturnReferences {
		ref := &proto.RefPropertiesRequest{
			ReferenceProperty: rr.PropertyName,
			TargetCollection:  rr.TargetCollection,
			Metadata: &proto.MetadataRequest{
				Uuid:               true,
				CreationTimeUnix:   rr.ReturnMetadata.CreatedAt,
				LastUpdateTimeUnix: rr.ReturnMetadata.LastUpdateAt,
				Vectors:            rr.ReturnVectors,
			},
		}
		// ReturnVectors were explicitly set to an empty slice, include the "only" vector.
		if len(rr.ReturnVectors) == 0 && rr.ReturnVectors != nil {
			ref.Metadata.Vector = true
			ref.Metadata.Vectors = nil
		}

		marshalReturnProperties(&ref.Properties, rr.ReturnProperties)
		req.Properties.RefProperties = append(req.Properties.RefProperties, ref)
	}
	return req, nil
}

func marshalReturnProperties(preq **proto.PropertiesRequest, rps []ReturnProperty) {
	dev.AssertNotNil(preq, "nil **req")

	if len(rps) == 0 && rps != nil {
		// ReturnProperties were explicitly set to an empty slice, do not return any.
		return
	}

	if *preq == nil {
		*preq = &proto.PropertiesRequest{}
	}
	req := *preq

	if rps == nil {
		// ReturnProperties were not set at all, default to all properties.
		req.ReturnAllNonrefProperties = true
		return
	}

	// visit traverses the ReturnProperty tree and
	// collects requested nested object properties.
	var visit func(*[]*proto.ObjectPropertiesRequest, *ReturnProperty)

	visit = func(os *[]*proto.ObjectPropertiesRequest, rp *ReturnProperty) {
		o := &proto.ObjectPropertiesRequest{PropName: rp.Name}
		for _, np := range rp.NestedProperties {
			if len(np.NestedProperties) == 0 {
				o.PrimitiveProperties = append(o.PrimitiveProperties, np.Name)
			} else {
				visit(&o.ObjectProperties, &np)
			}
		}
		*os = append(*os, o)
	}

	// Add all "primitive" and "nested" object properties to the request.
	for _, rp := range rps {
		if len(rp.NestedProperties) == 0 {
			req.NonRefProperties = append(req.NonRefProperties, rp.Name)
		} else {
			visit(&req.ObjectProperties, &rp)
		}
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
