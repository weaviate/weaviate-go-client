package schema

import (
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/usecases/config"
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

var defaultPQConfig = map[string]interface{}{
	"bitCompression": false,
	"centroids":      float64(256),
	"enabled":        false,
	"encoder": map[string]interface{}{
		"distribution": "log-normal",
		"type":         "kmeans",
	},
	"segments":      float64(0),
	"trainingLimit": float64(100_000),
}

var defaultBQConfig = map[string]interface{}{
	"enabled": false,
}

var defaultSQConfig = map[string]interface{}{
	"enabled":       false,
	"rescoreLimit":  float64(20),
	"trainingLimit": float64(100_000),
}

var defaultRQConfig = map[string]interface{}{
	"enabled":      false,
	"rescoreLimit": float64(20),
	"bits":         float64(8),
}

var defaultVectorIndexConfig = map[string]interface{}{
	"cleanupIntervalSeconds": float64(300),
	"efConstruction":         float64(128),
	"maxConnections":         float64(32),
	"vectorCacheMaxObjects":  float64(500000),
	"ef":                     float64(-1),
	"skip":                   false,
	"dynamicEfFactor":        float64(8),
	"dynamicEfMax":           float64(500),
	"dynamicEfMin":           float64(100),
	"flatSearchCutoff":       float64(40000),
	"distance":               "cosine",
	"pq":                     defaultPQConfig,
	"bq":                     defaultBQConfig,
	"sq":                     defaultSQConfig,
	"rq":                     defaultRQConfig,
	"filterStrategy":         "sweeping",
	"multivector":            defaultMultivectorConfig,
}

var defaultMultivectorConfig = map[string]interface{}{
	"enabled":     false,
	"aggregation": "maxSim",
	"muvera": map[string]interface{}{
		"enabled":      false,
		"dprojections": float64(16),
		"ksim":         float64(4),
		"repetitions":  float64(10),
	},
}

var defaultReplicationConfig = &models.ReplicationConfig{
	Factor:           1,
	DeletionStrategy: models.ReplicationConfigDeletionStrategyDeleteOnConflict,
}

var defaultMultiTenancyConfig = &models.MultiTenancyConfig{
	Enabled: false,
}
