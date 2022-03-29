package schema

import (
	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/semi-technologies/weaviate/usecases/config"
)

var defaultInvertedIndexConfig = &models.InvertedIndexConfig{
	CleanupIntervalSeconds: 60,
	Bm25: &models.BM25Config{
		K1: config.DefaultBM25k1,
		B:  config.DefaultBM25b,
	},
	Stopwords: &models.StopwordConfig{
		Preset: "en",
	},
}

var defaultModuleConfig = map[string]interface{}{
	"text2vec-contextionary": map[string]interface{}{
		"vectorizeClassName": true,
	},
}

var defaultShardingConfig = map[string]interface{}{
	"actualCount":         float64(1),
	"actualVirtualCount":  float64(128),
	"desiredCount":        float64(1),
	"desiredVirtualCount": float64(128),
	"function":            "murmur3",
	"key":                 "_id",
	"strategy":            "hash",
	"virtualPerPhysical":  float64(128),
}

var defaultVectorIndexConfig = map[string]interface{}{
	"cleanupIntervalSeconds": float64(300),
	"efConstruction":         float64(128),
	"maxConnections":         float64(64),
	"vectorCacheMaxObjects":  float64(500000),
	"ef":                     float64(-1),
	"skip":                   false,
	"dynamicEfFactor":        float64(8),
	"dynamicEfMax":           float64(500),
	"dynamicEfMin":           float64(100),
	"flatSearchCutoff":       float64(40000),
}
