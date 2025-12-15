package api

import (
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
	"github.com/weaviate/weaviate-go-client/v6/internal/gen/proto/v1"
)

type (
	SearchRequest struct {
		RequestDefaults

		Limit            *int
		Offset           *int
		AutoLimit        *int
		After            *string
		ReturnProperties []ReturnProperty
		ReturnReferences []ReturnReference
		ReturnVectors    []string
		ReturnMetadata   Set[MetadataRequest]
		NearVector       *NearVector
	}
	ReturnProperty struct {
		Name             string
		NestedProperties []string
	}
	ReturnReference struct {
		PropertyName     string
		TargetCollection string
		ReturnProperties []ReturnProperty
		ReturnMetadata   Set[MetadataRequest]
	}
	NearVector struct {
		Target    NearVectorTarget
		Certainty *float64
		Distance  *float64
	}
	NearVectorTarget interface {
		CombinationMethod() CombinationMethod
		Vectors() []TargetVector
	}
	TargetVector interface {
		Vector() *Vector
		Weight() float32
	}
	SearchResponse struct {
		TookSeconds    float32
		Results        []Object
		GroupByResults map[string]Group
	}
	Object struct {
		Metadata   ObjectMetadata
		Properties map[string]any
		References map[string][]Object
	}
	Group struct {
		Name                     string
		MinDistance, MaxDistance float32
		Size                     int64
		Objects                  []GroupByObject
	}
	GroupByObject struct {
		Object
		BelongsToGroup string
	}
	ObjectMetadata struct {
		UUID               string
		CreationTimeUnix   *int64
		LastUpdateTimeUnix *int64
		Distance           *float32
		Certainty          *float32
		Score              *float32
		ExplainScore       *string
		UnnamedVector      Vector
		NamedVectors       Vectors
	}
)

// Compile-time assertions that SearchRequest implements Request.
var _ Request = (*SearchRequest)(nil)

func (r *SearchRequest) Body() any { return r }

// ReturnOnlyVector is a sentinel value the caller can pass in SearchRequest.ReturnVectors
// to request the single vector in collection to be returned.
// This is different from passing api.DefaultVectorName, as the "only" vector might have a different name.
var ReturnOnlyVector *[]string

type MetadataRequest string

const (
	MetadataCreationTimeUnix   MetadataRequest = "CreationTimeUnix"
	MetadataLastUpdateTimeUnix MetadataRequest = "LastUpdateTimeUnix"
	MetadataDistance           MetadataRequest = "Distance"
	MetadataCertainty          MetadataRequest = "Certainty"
	MetadataScore              MetadataRequest = "Score"
	MetadataExplainScore       MetadataRequest = "ExplainScore"
)

type ConsistencyLevel string

const (
	_                      ConsistencyLevel = "" // unspecified
	ConsistencyLevelOne    ConsistencyLevel = "ONE"
	ConsistencyLevelQuorum ConsistencyLevel = "QUORUM"
	ConsistencyLevelAll    ConsistencyLevel = "ALL"
)

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
		return ptr(proto.ConsistencyLevel_CONSISTENCY_LEVEL_UNSPECIFIED)
	}
}

type CombinationMethod string

const (
	// Return from NearVectorTarget implementations which represent a single vector.
	combinationMethodUnspecified   CombinationMethod = ""
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

// MarshalSearchRequest() constructs a proto.SearchRequest.
func MarshalSearchRequest(req *SearchRequest) *proto.SearchRequest {
	sr := &proto.SearchRequest{
		Collection:       req.CollectionName,
		Tenant:           req.Tenant,
		ConsistencyLevel: req.ConsistencyLevel.proto(),
		Limit:            uint32(nilZero(req.Limit)),
		Offset:           uint32(nilZero(req.Offset)),
		Autocut:          uint32(nilZero(req.AutoLimit)),
		After:            nilZero(req.After),
	}

	var properties proto.PropertiesRequest
	if req.ReturnProperties == nil {
		// ReturnProperties were not set at all, default to all properties
		properties.ReturnAllNonrefProperties = true
	} else if len(req.ReturnProperties) > 0 {
		// Only return selected properties
		var nonRef []string
		var nested []*proto.ObjectPropertiesRequest

		for _, p := range req.ReturnProperties {
			if len(p.NestedProperties) == 0 {
				nonRef = append(nonRef, p.Name)
			} else {
				nested = append(nested, &proto.ObjectPropertiesRequest{
					PropName:            p.Name,
					PrimitiveProperties: p.NestedProperties,
					// TODO(dyma): add deeply-nested properties
				})
			}
		}

		properties.NonRefProperties = nonRef
		properties.ObjectProperties = nested
	} else {
		// ReturnProperties were explicitly set to an empty slice, do not return any.
	}

	// ReturnVectors were explicitly set to an empty slice, include the "only" vector.
	returnTheOnlyVector := len(req.ReturnVectors) == 0 && req.ReturnVectors != nil

	sr.Metadata = &proto.MetadataRequest{
		Uuid:               true,
		CreationTimeUnix:   req.ReturnMetadata.Contains(MetadataCreationTimeUnix),
		LastUpdateTimeUnix: req.ReturnMetadata.Contains(MetadataLastUpdateTimeUnix),
		Distance:           req.ReturnMetadata.Contains(MetadataDistance),
		Certainty:          req.ReturnMetadata.Contains(MetadataCertainty),
		Score:              req.ReturnMetadata.Contains(MetadataScore),
		ExplainScore:       req.ReturnMetadata.Contains(MetadataExplainScore),
		Vector:             returnTheOnlyVector,
		Vectors:            req.ReturnVectors,
	}

	switch {
	case req.NearVector != nil:
		sr.NearVector = &proto.NearVector{
			Distance:  nilZero(&req.NearVector.Distance),
			Certainty: nilZero(&req.NearVector.Certainty),
		}

		targets := req.NearVector.Target.Vectors()
		if len(targets) == 0 {
			break
		}

		if len(targets) == 1 {
			v := targets[0].Vector()
			dev.Assert(v != nil, "nil target vector")
			sr.NearVector.Vectors = []*proto.Vectors{
				marshalVector(v),
			}
		} else {
			combination := req.NearVector.Target.CombinationMethod().proto()

			// Pre-allocate slices for vectors, targets, and target weights.
			sr.NearVector.VectorForTargets = make([]*proto.VectorForTarget, len(targets))
			sr.NearVector.Targets = &proto.Targets{
				TargetVectors:     make([]string, len(targets)),
				WeightsForTargets: make([]*proto.WeightsForTarget, len(targets)),
				Combination:       combination,
			}

			for _, target := range targets {
				v := target.Vector()
				sr.NearVector.Targets.TargetVectors = append(sr.NearVector.Targets.TargetVectors, v.Name)
				sr.NearVector.Targets.WeightsForTargets = append(
					sr.NearVector.Targets.WeightsForTargets, &proto.WeightsForTarget{
						Target: v.Name,
						Weight: target.Weight(),
					})
				sr.NearVector.VectorForTargets = append(
					sr.NearVector.VectorForTargets, &proto.VectorForTarget{
						Name: v.Name,
						Vectors: []*proto.Vectors{
							marshalVector(v),
						},
					})
			}
		}
	default:
		// It is not a mistake to leave search method unset.
		// This would be the case when fetch objects with a conventional filter.
	}

	return sr
}

func marshalVector(v *Vector) *proto.Vectors {
	out := &proto.Vectors{Name: v.Name}
	switch {
	case v.Single != nil:
		out.Type = proto.Vectors_VECTOR_TYPE_SINGLE_FP32
		out.VectorBytes = marshalSingle(v.Single)
	case v.Multi != nil:
		out.Type = proto.Vectors_VECTOR_TYPE_MULTI_FP32
		out.VectorBytes = marshalMulti(v.Multi)
	default:
		return nil
	}
	return out
}

func UnmarshalSearchReply(reply *proto.SearchReply) *SearchResponse {
	dev.Assert(reply != nil, "search reply is nil")

	objects := make([]Object, len(reply.Results))
	for _, r := range reply.Results {
		if r == nil {
			continue
		}

		// At this point proto.SearchResult should not be nil; otherwise,
		// unmarshaling it is pointless. This also lets us access its fields
		// (.Metadata, .Properties) safely.
		dev.Assert(r != nil, "result object is nil")
		objects = append(objects, unmarshalObject(r.Properties, r.Metadata))
	}

	groups := make(map[string]Group, len(reply.GroupByResults))
	for _, r := range reply.GroupByResults {
		if r == nil {
			continue
		}

		// At this point proto.GroupByResult should not be nil; otherwise,
		// unmarshaling it is pointless. This also lets us access its fields
		// (.Metadata, .Properties) safely.
		dev.Assert(r != nil, "result group is nil")

		objects := make([]GroupByObject, len(r.Objects))
		for _, obj := range r.Objects {
			if obj == nil {
				continue
			}
			dev.Assert(obj != nil, "group object is nil")

			objects = append(objects, GroupByObject{
				BelongsToGroup: r.Name,
				Object:         unmarshalObject(obj.Properties, obj.Metadata),
			})
		}

		groups[r.Name] = Group{
			Name:        r.Name,
			MinDistance: r.MinDistance,
			MaxDistance: r.MaxDistance,
			Size:        r.NumberOfObjects,
			Objects:     objects,
		}
	}

	return &SearchResponse{
		TookSeconds:    reply.GetTook(),
		Results:        objects,
		GroupByResults: groups,
	}
}

func unmarshalVectors(vectors []*proto.Vectors) Vectors {
	out := make(Vectors, len(vectors))
	for _, vector := range vectors {
		v := Vector{Name: vector.Name}
		bytes := vector.GetVectorBytes()
		switch vector.Type {
		case proto.Vectors_VECTOR_TYPE_SINGLE_FP32:
			v.Single = unmarshalSingle(bytes)
		case proto.Vectors_VECTOR_TYPE_MULTI_FP32:
			v.Multi = unmarshalMulti(bytes)
		}
		out[v.Name] = v
	}
	return out
}

func unmarshalObject(pr *proto.PropertiesResult, mr *proto.MetadataResult) Object {
	properties := make(map[string]any, len(pr.GetNonRefProps().GetFields()))
	for name, property := range pr.GetNonRefProps().GetFields() {
		var v any
		switch property.GetKind().(type) {
		case *proto.Value_NullValue:
			v = nil
		case *proto.Value_TextValue:
			v = property.GetTextValue()
		case *proto.Value_IntValue:
			v = property.GetIntValue()
		case *proto.Value_NumberValue:
			v = property.GetNumberValue()
		case *proto.Value_DateValue:
			v = property.GetDateValue()
		case *proto.Value_BoolValue:
			v = property.GetBoolValue()
		case *proto.Value_BlobValue:
			v = property.GetBlobValue()
		default:
			// TODO(dyma): support other types
		}
		properties[name] = v
	}

	references := make(map[string][]Object, len(pr.GetRefProps()))
	for _, ref := range pr.GetRefProps() {
		if ref == nil {
			continue
		}
		dev.Assert(ref != nil, "reference is nil")
		if _, ok := references[ref.PropName]; !ok {
			references[ref.PropName] = make([]Object, len(ref.Properties))
		}
		for _, p := range ref.Properties {
			references[ref.PropName] = append(
				references[ref.PropName],
				unmarshalObject(p, p.Metadata),
			)
		}
	}
	return Object{
		Properties: properties,
		References: references,
		Metadata: ObjectMetadata{
			UUID:               unmarshalUUID(mr.GetIdAsBytes()),
			CreationTimeUnix:   nilPresent(mr.GetCreationTimeUnix(), mr.GetCreationTimeUnixPresent()),
			LastUpdateTimeUnix: nilPresent(mr.GetLastUpdateTimeUnix(), mr.GetLastUpdateTimeUnixPresent()),
			Distance:           nilPresent(mr.GetDistance(), mr.GetDistancePresent()),
			Certainty:          nilPresent(mr.GetCertainty(), mr.GetCertaintyPresent()),
			Score:              nilPresent(mr.GetScore(), mr.GetScorePresent()),
			ExplainScore:       nilPresent(mr.GetExplainScore(), mr.GetExplainScorePresent()),
			NamedVectors:       unmarshalVectors(mr.GetVectors()),
			UnnamedVector: Vector{
				Name:   DefaultVectorName,
				Single: unmarshalSingle(mr.GetVectorBytes()),
			},
		},
	}
}

// ptr is a helper for exporting pointers to constants.
func ptr[T any](v T) *T { return &v }

// nilZero returns the dereferenced value if the pointer is not nil
// and the zero value of type T otherwise.
func nilZero[T any](v *T) T {
	if v == nil {
		return *new(T)
	}
	return *v
}

func nilPresent[T any](v T, present bool) *T {
	if !present {
		return nil
	}
	return &v
}

func NewSet[Slice ~[]E, E comparable](values Slice) Set[E] {
	set := make(Set[E], len(values))
	for _, v := range values {
		set.Add(v)
	}
	return set
}

// Set is a lightweight set implementation based on on map[string]struct{}.
// It uses struct{} as a value type, because empty structs do not use memory.
type Set[E comparable] map[E]struct{}

func (s Set[E]) Add(v E) {
	s[v] = struct{}{}
}

func (s Set[E]) Contains(v E) bool {
	_, ok := s[v]
	return ok
}
