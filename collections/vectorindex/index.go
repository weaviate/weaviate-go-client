package vectorindex

import (
	"github.com/weaviate/weaviate-go-client/v6/internal"
)

// Registry stores all vector index types defined by this package.
var Registry internal.Modules[Type]

func init() {
	Registry.Register(*new(HFresh))
	Registry.Register(*new(Flat))
}

type Type string

var (
	_ internal.Module[Type] = (*HFresh)(nil)
	_ internal.Module[Type] = (*Flat)(nil)
)

type HFresh struct {
	Distance         Distance `json:"distance,omitempty"`
	MaxPostingSizeKB int      `json:"maxPostingSizeKB,omitempty"`
	ReplicaCount     int      `json:"replicas,omitempty"`
	SearchProbe      int      `json:"searchProbe,omitempty"`
}

func (HFresh) Name() Type { return "hfresh" }

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
