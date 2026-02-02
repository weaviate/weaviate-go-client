package api

import (
	"errors"
	"fmt"

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

	NearVector *NearVector
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

type (
	SearchTarget struct {
		CombinationMethod CombinationMethod
		Vectors           []TargetVector
	}
	TargetVector struct {
		Vector
		Weight *float32
	}
	NearVector struct {
		Target    SearchTarget
		Certainty *float64
		Distance  *float64
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

	var err error
	switch {
	case r.NearVector != nil:
		req.NearVector, err = marshalNearVector(r.NearVector)
	}
	if err != nil {
		return nil, err
	}
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

func marshalNearVector(req *NearVector) (*proto.NearVector, error) {
	nv := &proto.NearVector{
		Distance:  req.Distance,
		Certainty: req.Certainty,
	}

	switch len(req.Target.Vectors) {
	case 0:
		return nil, nil
	case 1:
		tv := req.Target.Vectors[0]
		v, err := marshalVector(&tv.Vector)
		if err != nil {
			return nil, fmt.Errorf("near vector: %w", err)
		}
		vectors := []*proto.Vectors{v}

		if tv.Name == "" {
			nv.Vectors = vectors
		} else {
			nv.VectorForTargets = append(nv.VectorForTargets, &proto.VectorForTarget{
				Name:    tv.Name,
				Vectors: vectors,
			})
		}
		return nv, nil
	}

	// Pre-allocate slices for vectors and targets.
	// Do not allocate WeightsForTarget, as targets may have no weights.
	nv.VectorForTargets = make([]*proto.VectorForTarget, len(req.Target.Vectors))
	nv.Targets = &proto.Targets{
		TargetVectors: make([]string, len(req.Target.Vectors)),
		Combination:   req.Target.CombinationMethod.proto(),
	}

	for i, tv := range req.Target.Vectors {
		v, err := marshalVector(&tv.Vector)
		if err != nil {
			return nil, fmt.Errorf("near vector: %w", err)
		}
		nv.Targets.TargetVectors[i] = tv.Name
		nv.VectorForTargets[i] = &proto.VectorForTarget{
			Name:    tv.Name,
			Vectors: []*proto.Vectors{v},
		}
		if tv.Weight != nil {
			nv.Targets.WeightsForTargets = append(nv.Targets.WeightsForTargets,
				&proto.WeightsForTarget{
					Target: tv.Name,
					Weight: *tv.Weight,
				})
		}
	}
	return nv, nil
}

// marshalVector marshals [Vector.Single] or [Vector.Multi] to bytes,
// depending on the presence. If neither is present it returns an error.
func marshalVector(v *Vector) (*proto.Vectors, error) {
	out := &proto.Vectors{Name: v.Name}
	switch {
	case v.Single != nil:
		out.Type = proto.Vectors_VECTOR_TYPE_SINGLE_FP32
		out.VectorBytes = marshalSingle(v.Single)
	case v.Multi != nil:
		out.Type = proto.Vectors_VECTOR_TYPE_MULTI_FP32
		out.VectorBytes = marshalMulti(v.Multi)
	default:
		return nil, errors.New("empty vector")
	}
	return out, nil
}

type CombinationMethod string

const (
	_                              CombinationMethod = ""
	CombinationMethodSum           CombinationMethod = "SUM"
	CombinationMethodMin           CombinationMethod = "MIN"
	CombinationMethodAverage       CombinationMethod = "AVERAGE"
	CombinationMethodManualWeights CombinationMethod = "MANUAL_WEIGHTS"
	CombinationMethodRelativeScore CombinationMethod = "RELATIVE_SCORE"
)

// proto converts CombinationMethod into a protobuf value.
func (cm CombinationMethod) proto() proto.CombinationMethod {
	switch cm {
	case CombinationMethodSum:
		return proto.CombinationMethod_COMBINATION_METHOD_TYPE_SUM
	case CombinationMethodMin:
		return proto.CombinationMethod_COMBINATION_METHOD_TYPE_MIN
	case CombinationMethodAverage:
		return proto.CombinationMethod_COMBINATION_METHOD_TYPE_AVERAGE
	case CombinationMethodManualWeights:
		return proto.CombinationMethod_COMBINATION_METHOD_TYPE_MANUAL
	case CombinationMethodRelativeScore:
		return proto.CombinationMethod_COMBINATION_METHOD_TYPE_RELATIVE_SCORE
	default:
		return proto.CombinationMethod_COMBINATION_METHOD_UNSPECIFIED
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
