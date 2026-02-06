package api

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	proto "github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/proto/v1"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
)

type SearchRequest struct {
	RequestDefaults

	Limit            int32
	Offset           int32
	AutoLimit        int32
	After            uuid.UUID
	ReturnProperties []ReturnProperty
	ReturnReferences []ReturnReference
	ReturnVectors    []string
	ReturnMetadata   ReturnMetadata
	GroupBy          *GroupBy

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
		ReturnMetadata   ReturnMetadata
		ReturnVectors    []string
		ReturnProperties []ReturnProperty
		ReturnReferences []ReturnReference
	}
	GroupBy struct {
		Property       string // Property to group by.
		ObjectLimit    int32  // Maximum number of objects per group.
		NumberOfGroups int32  // Maximum number of groups to return.
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
	dev.AssertNotNil(r, "r")

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

	if r.GroupBy != nil {
		req.GroupBy = &proto.GroupBy{
			Path:            []string{r.GroupBy.Property},
			ObjectsPerGroup: r.GroupBy.ObjectLimit,
			NumberOfGroups:  r.GroupBy.NumberOfGroups,
		}
	}

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
	dev.AssertNotNil(req, "req")

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
	dev.AssertNotNil(req, "req")

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
	dev.AssertNotNil(req, "req")

	if len(vectors) == 0 && vectors != nil {
		// ReturnVectors were explicitly set to an empty slice, include the "only" vector.
		req.Vector = true
		req.Vectors = nil
	} else {
		req.Vectors = vectors
	}
}

func marshalNearVector(req *NearVector) (*proto.NearVector, error) {
	dev.AssertNotNil(req, "req")

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
		dev.AssertNotNil(v, "v")
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
		dev.AssertNotNil(v, "v")

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
	dev.AssertNotNil(v, "v")

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

type SearchResponse struct {
	Took           time.Duration
	Results        []Object
	GroupByResults []Group
}

var _ MessageUnmarshaler[proto.SearchReply] = (*SearchResponse)(nil)

type (
	Object struct {
		Collection string
		Metadata   ObjectMetadata
		Properties map[string]any
		References map[string][]Object
	}
	ObjectMetadata struct {
		UUID          uuid.UUID
		CreatedAt     *time.Time
		LastUpdatedAt *time.Time
		Distance      *float32
		Certainty     *float32
		Score         *float32
		ExplainScore  *string
		Vectors       Vectors
	}
	Group struct {
		Name                     string
		MinDistance, MaxDistance float32
		Size                     int64
		Objects                  []GroupObject
	}
	GroupObject struct {
		Object
		BelongsToGroup string
	}
)

// UnmarshalMessage reads proto.SearchReply into this SearchResponse.
func (r *SearchResponse) UnmarshalMessage(reply *proto.SearchReply) error {
	dev.AssertNotNil(reply, "reply")

	objects := make([]Object, len(reply.Results))
	if err := unmarshalObjects(reply.Results, func(i int, o Object) {
		objects[i] = o
	}); err != nil {
		return err
	}

	groups := make([]Group, len(reply.GroupByResults))
	for gi, group := range reply.GroupByResults {
		if group == nil {
			continue
		}
		dev.AssertNotNil(group, "group")

		objects := make([]GroupObject, len(group.Objects))
		if err := unmarshalObjects(group.Objects, func(i int, o Object) {
			objects[i] = GroupObject{
				BelongsToGroup: group.Name,
				Object:         o,
			}
		}); err != nil {
			return nil
		}

		groups[gi] = Group{
			Name:        group.Name,
			MinDistance: group.MinDistance,
			MaxDistance: group.MaxDistance,
			Size:        group.NumberOfObjects,
			Objects:     objects,
		}
	}

	*r = SearchResponse{
		Took:           time.Duration(reply.Took) * time.Second,
		Results:        objects,
		GroupByResults: groups,
	}
	return nil
}

// unmarshalObjects calls [unmarshalObject] for each search result
// and passes the object and its index to the consumer f.
// Any non-nil error stops the iteration and is returned immediately.
func unmarshalObjects(objects []*proto.SearchResult, f func(int, Object)) error {
	dev.AssertNotNil(f, "f")

	for i, object := range objects {
		if object == nil {
			continue
		}
		dev.AssertNotNil(object, "object")

		o, err := unmarshalObject(object.Properties, object.Metadata)
		if err != nil {
			return err
		}
		dev.AssertNotNil(o, "o")

		f(i, *o)
	}
	return nil
}

func unmarshalObject(pr *proto.PropertiesResult, mr *proto.MetadataResult) (*Object, error) {
	properties, err := unmarshalProperties(pr.GetNonRefProps())
	if err != nil {
		return nil, err
	}
	dev.AssertNotNil(properties, "properties")

	references := make(map[string][]Object, len(pr.GetRefProps()))
	for _, ref := range pr.GetRefProps() {
		if ref == nil {
			continue
		}
		dev.AssertNotNil(ref, "ref")

		if _, ok := references[ref.PropName]; !ok {
			references[ref.PropName] = make([]Object, 0, len(ref.Properties))
		}
		for _, p := range ref.Properties {
			o, err := unmarshalObject(p, p.Metadata)
			if err != nil {
				return nil, err
			}
			dev.AssertNotNil(o, "o")
			references[ref.PropName] = append(references[ref.PropName], *o)
		}
	}

	var metadata ObjectMetadata
	if mr != nil {
		var id uuid.UUID
		if bytes := mr.GetIdAsBytes(); bytes != nil {
			fromBytes, err := uuid.FromBytes(bytes)
			if err != nil {
				return nil, err
			}
			id = fromBytes
		}

		vectors, err := unmarshalVectors(mr)
		if err != nil {
			return nil, err
		}
		dev.AssertNotNil(vectors, "vectors")

		metadata = ObjectMetadata{
			UUID:          id,
			CreatedAt:     nilPresent(timeFromUnix(mr.CreationTimeUnix), mr.CreationTimeUnixPresent),
			LastUpdatedAt: nilPresent(timeFromUnix(mr.LastUpdateTimeUnix), mr.LastUpdateTimeUnixPresent),
			Distance:      nilPresent(mr.Distance, mr.DistancePresent),
			Certainty:     nilPresent(mr.Certainty, mr.CertaintyPresent),
			Score:         nilPresent(mr.Score, mr.ScorePresent),
			ExplainScore:  nilPresent(mr.ExplainScore, mr.ExplainScorePresent),
			Vectors:       vectors,
		}
	}

	return &Object{
		Collection: pr.GetTargetCollection(),
		Metadata:   metadata,
		Properties: properties,
		References: references,
	}, nil
}

// unmarshalProperties unmarshals map[string]proto.Value into map[string]any.
// ps can be nil, in which case an empty map is returned.
func unmarshalProperties(ps *proto.Properties) (map[string]any, error) {
	out := make(map[string]any, len(ps.GetFields()))
	for name, f := range ps.GetFields() {
		var v any
		switch f.GetKind().(type) {
		case *proto.Value_NullValue:
			v = nil
		case *proto.Value_TextValue:
			v = f.GetTextValue()
		case *proto.Value_IntValue:
			v = f.GetIntValue()
		case *proto.Value_NumberValue:
			v = f.GetNumberValue()
		case *proto.Value_BoolValue:
			v = f.GetBoolValue()
		case *proto.Value_BlobValue:
			v = f.GetBlobValue()
		case *proto.Value_DateValue:
			t, err := timeFromString(f.GetDateValue())
			if err != nil {
				return nil, err
			}
			v = t
		case *proto.Value_UuidValue:
			id, err := uuid.Parse(f.GetUuidValue())
			if err != nil {
				return nil, err
			}
			v = id
		case *proto.Value_ObjectValue:
			properties, err := unmarshalProperties(f.GetObjectValue())
			if err != nil {
				return nil, err
			}
			dev.AssertNotNil(properties, "properties")
			v = properties
		default:
			// TODO(dyma): support array types
		}
		out[name] = v
	}
	return out, nil
}

func unmarshalVectors(mr *proto.MetadataResult) (Vectors, error) {
	out := make(Vectors, len(mr.GetVectors()))
	for _, vector := range mr.GetVectors() {
		v := Vector{Name: vector.Name}
		bytes := vector.GetVectorBytes()
		switch vector.Type {
		case proto.Vectors_VECTOR_TYPE_SINGLE_FP32:
			v.Single = unmarshalSingle(bytes)
		case proto.Vectors_VECTOR_TYPE_MULTI_FP32:
			v.Multi = unmarshalMulti(bytes)
		default:
			return nil, fmt.Errorf("unknown type for vector %q", vector.Name)
		}
		out[v.Name] = v
	}

	if v := mr.GetVectorBytes(); len(v) > 0 {
		out[DefaultVectorName] = Vector{
			Name:   DefaultVectorName,
			Single: unmarshalSingle(v),
		}
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

// nilPresent returns a pointer to v if present == true and nil otherwise.
func nilPresent[T any](v T, present bool) *T {
	if !present {
		return nil
	}
	return &v
}
