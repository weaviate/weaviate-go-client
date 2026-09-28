package vectorindex

import (
	"time"

	"github.com/weaviate/weaviate-go-client/v6/collections/compression"
	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
)

// Registry stores all vector index types defined by this package.
var Registry internal.Modules[Type]

func init() {
	Registry.Register(*new(HFresh))
	Registry.Register(*new(HNSW))
	Registry.Register(*new(Flat))
	Registry.Register(*new(Dynamic))
}

type Type string

var (
	_ internal.Module[Type] = (*HFresh)(nil)
	_ internal.Module[Type] = (*HNSW)(nil)
	_ internal.Module[Type] = (*Flat)(nil)
	_ internal.Module[Type] = (*Dynamic)(nil)
)

type HFresh struct {
	Distance         Distance `json:"distance,omitempty"`
	MaxPostingSizeKB int      `json:"maxPostingSizeKB,omitempty"`
	ReplicaCount     int      `json:"replicas,omitempty"`
	SearchProbe      int      `json:"searchProbe,omitempty"`
}

func (HFresh) Name() Type { return "hfresh" }

type HNSW struct {
	Distance              Distance
	FilterStrategy        FilterStrategy
	Ef                    int
	EfConstruction        int
	MaxConnections        int
	VectorCacheMaxObjects int64
	CleanupInterval       time.Duration

	// TODO(dyma): support multi-vector
	// MultiVector            MultiVector

	DynamicEfMin      int
	DynamicEfMax      int
	DynamicEfFactor   int
	FlatSearchCutoff  int
	SkipVectorization bool
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

type hnswJSON struct {
	Distance              Distance       `json:"distance,omitempty"`
	FilterStrategy        FilterStrategy `json:"filterStrategy,omitempty"`
	Ef                    int            `json:"ef,omitempty"`
	EfConstruction        int            `json:"efConstruction,omitempty"`
	MaxConnections        int            `json:"maxConnections,omitempty"`
	VectorCacheMaxObjects int64          `json:"vectorCacheMaxObjects,omitempty"`
	CleanupInterval       int            `json:"cleanupIntervalSeconds,omitempty"`

	// TODO(dyma): support multi-vector
	// MultiVector            MultiVector    `json:"multivector,omitmepty"`

	DynamicEfMin      int  `json:"dynamicEfMin,omitempty"`
	DynamicEfMax      int  `json:"dynamicEfMax,omitempty"`
	DynamicEfFactor   int  `json:"dynamicEfFactor,omitempty"`
	FlatSearchCutoff  int  `json:"flatSearchCutoff,omitempty"`
	SkipVectorization bool `json:"skip,omitempty"`
}

func (hnsw HNSW) EncodeMap() (map[string]any, error) {
	dest := make(map[string]any)
	if err := internal.Encode(hnswJSON{
		Distance:              hnsw.Distance,
		FilterStrategy:        hnsw.FilterStrategy,
		Ef:                    hnsw.Ef,
		EfConstruction:        hnsw.EfConstruction,
		MaxConnections:        hnsw.MaxConnections,
		VectorCacheMaxObjects: hnsw.VectorCacheMaxObjects,
		CleanupInterval:       int(hnsw.CleanupInterval.Seconds()),

		// TODO(dyma): support multi-vector
		// MultiVector            MultiVector    `json:"multivector,omitmepty"`

		DynamicEfMin:      hnsw.DynamicEfMin,
		DynamicEfMax:      hnsw.DynamicEfMax,
		DynamicEfFactor:   hnsw.DynamicEfFactor,
		FlatSearchCutoff:  hnsw.FlatSearchCutoff,
		SkipVectorization: hnsw.SkipVectorization,
	}, dest); err != nil {
		return nil, err
	}
	return dest, nil
}

func (hnsw *HNSW) DecodeMap(m map[string]any) error {
	var dest hnswJSON
	if err := internal.Decode(m, &dest); err != nil {
		return nil
	}
	*hnsw = HNSW{
		Distance:              dest.Distance,
		FilterStrategy:        dest.FilterStrategy,
		Ef:                    dest.Ef,
		EfConstruction:        dest.EfConstruction,
		MaxConnections:        dest.MaxConnections,
		VectorCacheMaxObjects: dest.VectorCacheMaxObjects,
		CleanupInterval:       time.Duration(dest.CleanupInterval) * time.Second,

		// TODO(dyma): support multi-vector
		// MultiVector            MultiVector    `json:"multivector,omitmepty"`

		DynamicEfMin:      dest.DynamicEfMin,
		DynamicEfMax:      dest.DynamicEfMax,
		DynamicEfFactor:   dest.DynamicEfFactor,
		FlatSearchCutoff:  dest.FlatSearchCutoff,
		SkipVectorization: dest.SkipVectorization,
	}
	return nil
}

type Flat struct {
	VectorCacheMaxObjects int64 `json:"vectorCacheMaxObjects,omitempty"`
}

func (Flat) Name() Type { return "flat" }

type Dynamic struct {
	Distance  Distance
	Threshold int64

	HNSW            HNSW
	HNSWCompression internal.Module[compression.Type]

	Flat            Flat
	FlatCompression internal.Module[compression.Type]

	// TODO(dyma): find a good way to disable compression on either.
	// Possibly the simplest way is compression.None with a custom codec.
}

func (Dynamic) Name() Type { return "dynamic" }

var (
	_ internal.MapEncoder = (*Dynamic)(nil)
	_ internal.MapDecoder = (*Dynamic)(nil)
)

func (d Dynamic) EncodeMap() (map[string]any, error) {
	dev.AssertNotNil(d, "dynamic")

	hnsw := make(map[string]any)
	if err := internal.Encode(d.HNSW, hnsw); err != nil {
		return nil, err
	}
	if c := d.HNSWCompression; c != nil {
		conf, err := compression.Registry.Encode(c)
		if err != nil {
			return nil, err
		}
		hnsw[string(c.Name())] = conf
	}

	flat := make(map[string]any)
	if err := internal.Encode(d.Flat, flat); err != nil {
		return nil, err
	}
	if c := d.FlatCompression; c != nil {
		conf, err := compression.Registry.Encode(c)
		if err != nil {
			return nil, err
		}
		hnsw[string(c.Name())] = conf
	}

	return map[string]any{
		"distance":  d.Distance,
		"threshold": d.Threshold,
		"hnsw":      hnsw,
		"flat":      flat,
	}, nil
}

func (d *Dynamic) DecodeMap(m map[string]any) error {
	dev.AssertNotNil(d, "dynamic")
	if m == nil {
		return nil
	}

	var decoded Dynamic
	if distance, ok := m["distance"]; ok {
		switch distance := distance.(type) {
		case string:
			decoded.Distance = Distance(distance)
		case Distance:
			decoded.Distance = distance
		}
	}

	if threshold, ok := m["threshold"].(int64); ok {
		decoded.Threshold = threshold
	}

	if m, ok := m["hnsw"].(map[string]any); ok {
		hnsw, err := Registry.Decode("hnsw", m)
		if err != nil {
			return err
		}
		dev.AssertType[HNSW](hnsw, "hnsw module")
		decoded.HNSW = hnsw.(HNSW)

		if key, ok := compression.Registry.Find(m); ok {
			if m, ok := m[key].(map[string]any); ok {
				c, err := compression.Registry.Decode(key, m)
				if err != nil {
					return err
				}
				decoded.HNSWCompression = c
			}
		}
	}

	if m, ok := m["flat"].(map[string]any); ok {
		flat, err := Registry.Decode("flat", m)
		if err != nil {
			return err
		}
		dev.AssertType[Flat](flat, "flat module")
		decoded.Flat = flat.(Flat)

		// TODO(dyma): maybe we can apply this approach to how compression config
		// is decoded for "normal" modules, where api/ layer extracts it and
		// the public layer decodes it.
		// The key idea is that presenting Compression as a separate field and
		// not as part of the Index configuration is a DX requirement (!!).
		// Whereas api/ is meant to only handle Weaviate API's quirks.
		//
		// collection/compression should expose something like:
		// - compression.Encode(m Module, conf map[string]any)
		// - compression.Decode(conf map[string]any) Module
		// to handle this. The api/ layer should not need to care.
		if key, ok := compression.Registry.Find(m); ok {
			if m, ok := m[key].(map[string]any); ok {
				c, err := compression.Registry.Decode(key, m)
				if err != nil {
					return err
				}
				decoded.FlatCompression = c
			}
		}
	}

	*d = decoded
	return nil
}
