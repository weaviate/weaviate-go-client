package vectorindex

import (
	"github.com/weaviate/weaviate-go-client/v6/internal"
)

// Registry stores all vector index types defined by this package.
var Registry internal.Modules[Type]

func init() {
	Registry.Register(*new(HFresh))
}

type Type string

var _ internal.Module[Type] = (*HFresh)(nil)

type HFresh struct {
	Distance         Distance `json:"distance"`
	MaxPostingSizeKB int      `json:"maxPostingSizeKB"`
	ReplicaCount     int      `json:"replicas"`
	SearchProbe      int      `json:"searchProbe"`
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
