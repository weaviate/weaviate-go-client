package modules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/modules"
	"github.com/weaviate/weaviate-go-client/v6/modules/model2vec"
	"github.com/weaviate/weaviate-go-client/v6/modules/selfprovided"
)

// TestModules ensures that all modules are registerred with [modules.Registry]
// and that they produce correct configurations when serialized.
func TestModules(t *testing.T) {
	for _, tt := range []struct {
		name   string         // Module name.
		module modules.Module // Module configuration.
		conf   map[string]any // Expected configuration.
	}{
		{
			name:   "none",
			module: selfprovided.Vectorizer,
			conf:   make(map[string]any),
		},
		{
			name: "text2vec-model2vec",
			module: model2vec.Text2Vec{
				URL:        "example.com",
				Properties: []string{"title", "lyrics"},
			},
			conf: map[string]any{
				"inferenceURL": "example.com",
				"properties":   []string{"title", "lyrics"},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.name, tt.module.Name(), "module name")

			conf, err := modules.Registry.Encode(tt.module)
			require.NoError(t, err, "encode")
			assert.Equal(t, tt.conf, conf, "encoded configuration")

			module, err := modules.Registry.Decode(tt.name, conf)
			require.NoError(t, err, "decode")
			assert.EqualExportedValues(t, tt.module, module, "decoded module")
		})
	}
}
