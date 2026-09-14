package compression_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/collections/compression"
	"github.com/weaviate/weaviate-go-client/v6/internal"
)

// TestCompression ensures that all modules are registerred with [compression.Registry]
// and that they produce correct configurations when serialized.
func TestCompression(t *testing.T) {
	for _, tt := range []struct {
		name   compression.Type                  // Module name.
		module internal.Module[compression.Type] // Compression configuration.
		conf   map[string]any                    // Expected configuration.
	}{
		{
			name: "rq",
			module: compression.RQ{
				Bits:         1,
				RescoreLimit: 2,
				Cache:        true,
			},
			conf: map[string]any{
				"bits":          1,
				"rescore_limit": 2,
				"cache":         true,
			},
		},
		{
			name:   "rq default",
			module: compression.RQ{},
			conf:   make(map[string]any),
		},
		{
			name: "bq",
			module: compression.BQ{
				RescoreLimit: 1,
				Cache:        true,
			},
			conf: map[string]any{
				"rescore_limit": 1,
				"cache":         true,
			},
		},
		{
			name:   "bq default",
			module: compression.BQ{},
			conf:   make(map[string]any),
		},
		{
			name: "pq",
			module: compression.PQ{
				Centroids:           1,
				Segments:            2,
				TrainingLimit:       3,
				Encoder:             compression.PQEncoderKmeans,
				EncoderDistribution: compression.PQDistributionLogNormal,
				BitCompression:      true,
			},
			conf: map[string]any{
				"centroids":            1,
				"segments":             2,
				"training_limit":       3,
				"encoder_type":         compression.PQEncoderKmeans,
				"encoder_distribution": compression.PQDistributionLogNormal,
				"bit_compression":      true,
			},
		},
		{
			name:   "pq default",
			module: compression.PQ{},
			conf:   make(map[string]any),
		},
	} {
		t.Run(string(tt.name), func(t *testing.T) {
			name := strings.Split(string(tt.name), " ")[0]
			assert.EqualValues(t, name, tt.module.Name(), "module name")

			conf, err := compression.Registry.Encode(tt.module)
			require.NoError(t, err, "encode")
			assert.Equal(t, tt.conf, conf, "encoded configuration")

			module, err := compression.Registry.Decode(name, conf)
			require.NoError(t, err, "decode")
			assert.EqualExportedValues(t, tt.module, module, "decoded module")
		})
	}
}
