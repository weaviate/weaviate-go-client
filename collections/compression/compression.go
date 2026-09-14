package compression

import "github.com/weaviate/weaviate-go-client/v6/internal"

// Registry stores all compression algorithms defined by this package.
var Registry internal.Modules[Type]

func init() {
	Registry.Register(*new(BQ))
	Registry.Register(*new(PQ))
	Registry.Register(*new(RQ))
}

type Type string

var (
	_ internal.Module[Type] = (*BQ)(nil)
	_ internal.Module[Type] = (*PQ)(nil)
	_ internal.Module[Type] = (*RQ)(nil)
)

// Rotational quantization.
type RQ struct {
	Bits         int  `json:"bits,omitempty"`
	RescoreLimit int  `json:"rescore_limit,omitempty"`
	Cache        bool `json:"cache,omitempty"`
}

func (RQ) Name() Type { return "rq" }

// Binary quantization.
type BQ struct {
	RescoreLimit int  `json:"rescore_limit,omitempty"`
	Cache        bool `json:"cache,omitempty"`
}

func (BQ) Name() Type { return "bq" }

// Product quantization.
type PQ struct {
	Centroids           int                   `json:"centroids,omitempty"`
	Segments            int                   `json:"segments,omitempty"`
	TrainingLimit       int                   `json:"training_limit,omitempty"`
	Encoder             PQEncoder             `json:"encoder_type,omitempty"`
	EncoderDistribution PQEncoderDistribution `json:"encoder_distribution,omitempty"`
	BitCompression      bool                  `json:"bit_compression,omitempty"`
}

func (PQ) Name() Type { return "pq" }

type PQEncoder string

const (
	PQEncoderKmeans PQEncoder = "kmeans"
	PQEncoderTile   PQEncoder = "tile"
)

type PQEncoderDistribution string

const (
	PQDistributionNormal    PQEncoderDistribution = "normal"
	PQDistributionLogNormal PQEncoderDistribution = "log-normal"
)
