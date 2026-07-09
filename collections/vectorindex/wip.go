package vectorindex

import (
	"github.com/weaviate/weaviate-go-client/v6/internal"
)

var (
	_ internal.Module[Type] = (*HNSW)(nil)
	_ internal.Module[Type] = (*Flat)(nil)
	_ internal.Module[Type] = (*Dynamic)(nil)
)

type HNSW struct {
	Distance               Distance    `json:"distance"`
	Ef                     int         `json:"ef"`
	EfConstruction         int         `json:"efConstruction"`
	MaxConnections         int         `json:"maxConnections"`
	VectorMaxCacheObjects  int         `json:"vectorCacheMaxObjects"`
	CleanupIntervalSeconds any         `json:"cleanupIntervalSeconds"` // TODO(dyma): time.Duration?
	FilterStrategy         any         `json:"filterStrategy"`
	MultiVector            MultiVector `json:"multivector"`

	DynamicEfMin      int  `json:"dynamicEfMin"`
	DynamicEfMax      int  `json:"dynamicEfMax"`
	DynamicEfFactor   int  `json:"dynamicEfFactor"`
	FlatSearchCutoff  int  `json:"flatSearchCutoff"`
	SkipVectorization bool `json:"skip"`
}

func (HNSW) Name() Type { return "hnsw" }

type Flat struct {
	VectorMaxCacheObjects int `json:"vectorCacheMaxObjects"`
}

func (Flat) Name() Type { return "flat" }

type Dynamic struct {
	HNSW      HNSW  `json:"hnsw"`
	Flat      Flat  `json:"flat"`
	Threshold int64 `json:"threshold"`
}

func (Dynamic) Name() Type { return "dynamic" }

type MultiVector struct {
	Enabled bool `json:"enabled"`
}
