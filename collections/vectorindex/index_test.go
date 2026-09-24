package vectorindex_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/collections/vectorindex"
	"github.com/weaviate/weaviate-go-client/v6/internal"
)

// TestIndex ensures that all modules are registerred with [vectorindex.Registry]
// and that they produce correct configurations when serialized.
func TestIndex(t *testing.T) {
	for _, tt := range []struct {
		name   vectorindex.Type                  // Module name.
		module internal.Module[vectorindex.Type] // Index configuration.
		conf   map[string]any                    // Expected configuration.
	}{
		{
			name: "hfresh",
			module: vectorindex.HFresh{
				Distance:         vectorindex.DistanceCosine,
				MaxPostingSizeKB: 1,
				ReplicaCount:     2,
				SearchProbe:      3,
			},
			conf: map[string]any{
				"distance":         vectorindex.DistanceCosine,
				"maxPostingSizeKB": 1,
				"replicas":         2,
				"searchProbe":      3,
			},
		},
		{
			name:   "hfresh default",
			module: vectorindex.HFresh{},
			conf:   make(map[string]any),
		},
		{
			name: "flat",
			module: vectorindex.Flat{
				VectorCacheMaxObjects: 80085,
			},
			conf: map[string]any{
				"vectorCacheMaxObjects": int64(80085),
			},
		},
		{
			name:   "flat default",
			module: vectorindex.Flat{},
			conf:   make(map[string]any),
		},
	} {
		t.Run(string(tt.name), func(t *testing.T) {
			name := strings.Split(string(tt.name), " ")[0]
			assert.EqualValues(t, name, tt.module.Name(), "module name")

			conf, err := vectorindex.Registry.Encode(tt.module)
			require.NoError(t, err, "encode")
			assert.Equal(t, tt.conf, conf, "encoded configuration")

			module, err := vectorindex.Registry.Decode(name, conf)
			require.NoError(t, err, "decode")
			assert.EqualExportedValues(t, tt.module, module, "decoded module")
		})
	}
}
