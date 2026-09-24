package vectorindex

import (
	"github.com/weaviate/weaviate-go-client/v6/internal"
)

// Registry stores all vector index types defined by this package.
var Registry internal.Modules[Type]

func init() {
	Registry.Register(*new(HFresh))
	Registry.Register(*new(HNSW))
	Registry.Register(*new(Flat))
}

type Type string

var (
	_ internal.Module[Type] = (*HFresh)(nil)
	_ internal.Module[Type] = (*HNSW)(nil)
	_ internal.Module[Type] = (*Flat)(nil)
)

type HFresh struct {
	Distance         Distance `json:"distance,omitempty"`
	MaxPostingSizeKB int      `json:"maxPostingSizeKB,omitempty"`
	ReplicaCount     int      `json:"replicas,omitempty"`
	SearchProbe      int      `json:"searchProbe,omitempty"`
}

func (HFresh) Name() Type { return "hfresh" }

type HNSW struct {
	Distance               Distance       `json:"distance,omitempty"`
	FilterStrategy         FilterStrategy `json:"filterStrategy,omitempty"`
	Ef                     int            `json:"ef,omitempty"`
	EfConstruction         int            `json:"efConstruction,omitempty"`
	MaxConnections         int            `json:"maxConnections,omitempty"`
	VectorCacheMaxObjects  int64          `json:"vectorCacheMaxObjects,omitempty"`
	CleanupIntervalSeconds int            `json:"cleanupIntervalSeconds,omitempty"`

	// TODO(dyma): support multi-vector
	// MultiVector            MultiVector    `json:"multivector,omitmepty"`

	DynamicEfMin      int  `json:"dynamicEfMin,omitempty"`
	DynamicEfMax      int  `json:"dynamicEfMax,omitempty"`
	DynamicEfFactor   int  `json:"dynamicEfFactor,omitempty"`
	FlatSearchCutoff  int  `json:"flatSearchCutoff,omitempty"`
	SkipVectorization bool `json:"skip,omitempty"`
}

func (HNSW) Name() Type { return "hnsw" }

// FilterStrategy is the algorithm for calculating vector distances.
type FilterStrategy string

const (
	FilterStrategySweeping = FilterStrategy("sweeping")
	FilterStrategyACORN    = FilterStrategy("acorn")
)

// Distance is the algorithm for calculating vector distances.
type Distance string

const (
	DistanceCosine    = Distance("cosine")
	DistanceDot       = Distance("dot")
	DistanceL2Squared = Distance("l2-squared")
	DistanceHamming   = Distance("hamming")
	DistanceManhattan = Distance("manhattan")
)

type Flat struct {
	VectorCacheMaxObjects int64 `json:"vectorCacheMaxObjects,omitempty"`
}

func (Flat) Name() Type { return "flat" }
