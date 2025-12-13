package api

import (
	"github.com/weaviate/weaviate-go-client/v6/internal/gen/proto/v1"
)

type Request any

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
		ReturnMetadata   Set[Metadata]
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
		// TODO: ReturnMetadata
	}
	NearVector struct {
		Target    NearVectorTarget
		Certainty *float64
		Distance  *float64
	}
	NearVectorTarget interface {
		CombinationMethod() CombinationMethod
		Targets() []TargetVector
	}
	TargetVector interface {
		Vector() *Vector
		Weight() float64
	}
	SearchResponse proto.SearchResult
)

// ReturnOnlyVector is a sentinel value the caller can pass in SearchRequest.ReturnVectors
// to request the single vector in collection to be returned.
// This is different from passing api.DefaultVectorName, as the "only" vector might have a different name.
var ReturnOnlyVector []string

type Metadata string

const (
	MetadataCreationTimeUnix   Metadata = "CreationTimeUnix"
	MetadataLastUpdateTimeUnix Metadata = "LastUpdateTimeUnix"
	MetadataDistance           Metadata = "Distance"
	MetadataCertainty          Metadata = "Certainty"
	MetadataScore              Metadata = "Score"
	MetadataExplainScore       Metadata = "ExplainScore"
)

type ConsistencyLevel string

const (
	consistencyLevelUnspecified ConsistencyLevel = ""
	ConsistencyLevelOne         ConsistencyLevel = "ONE"
	ConsistencyLevelQuorum      ConsistencyLevel = "QUORUM"
	ConsistencyLevelAll         ConsistencyLevel = "ALL"
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

// ptr is a helper for exporting pointers to constants.
func ptr[T any](v T) *T { return &v }

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
