package modules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/generate"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/modules"
	"github.com/weaviate/weaviate-go-client/v6/modules/meta"
	proto "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

func TestGenerativeMeta(t *testing.T) {
	for _, tt := range testkit.WithOnly(t, []struct {
		testkit.Only
		name   string
		module generate.Provider
		want   *proto.GenerativeProvider
	}{
		{
			name: "configured",
			module: meta.Generative{
				BaseURL:          "https://proxy.example.com",
				Model:            "muse-spark-1.2",
				Temperature:      new(0.7),
				TopP:             new(0.9),
				MaxTokens:        new(int64(128)),
				FrequencyPenalty: new(0.1),
				PresencePenalty:  new(0.2),
				ReasoningEffort:  meta.HighEffort,
				Images:           []string{"pony.png"},
				ImageProperties:  []string{"cover"},
				ReturnMetadata:   true,
			},
			want: &proto.GenerativeProvider{
				ReturnMetadata: true,
				Kind: &proto.GenerativeProvider_Meta{
					Meta: &proto.GenerativeMeta{
						BaseUrl:          new("https://proxy.example.com"),
						Model:            new("muse-spark-1.2"),
						Temperature:      new(0.7),
						TopP:             new(0.9),
						MaxTokens:        new(int64(128)),
						FrequencyPenalty: new(0.1),
						PresencePenalty:  new(0.2),
						ReasoningEffort:  proto.GenerativeMeta_REASONING_EFFORT_HIGH.Enum(),
						Images:           &proto.TextArray{Values: []string{"pony.png"}},
						ImageProperties:  &proto.TextArray{Values: []string{"cover"}},
					},
				},
			},
		},
		{
			name:   "defaults",
			module: meta.Generative{},
			want: &proto.GenerativeProvider{
				Kind: &proto.GenerativeProvider_Meta{
					Meta: &proto.GenerativeMeta{},
				},
			},
		},
	}) {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.module.GenerativeProvider()
			require.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestGenerativeMetaConfig(t *testing.T) {
	module := meta.Generative{
		BaseURL:          "https://proxy.example.com",
		Model:            "muse-spark-1.2",
		Temperature:      new(0.7),
		TopP:             new(0.9),
		MaxTokens:        new(int64(128)),
		FrequencyPenalty: new(0.1),
		PresencePenalty:  new(0.2),
		ReasoningEffort:  meta.XHighEffort,
		Images:           []string{"pony.png"},
		ImageProperties:  []string{"cover"},
		ReturnMetadata:   true,
	}

	conf, err := modules.Registry.Encode(module)
	require.NoError(t, err, "encode")
	assert.Equal(t, map[string]any{
		"baseURL":          "https://proxy.example.com",
		"model":            "muse-spark-1.2",
		"temperature":      new(0.7),
		"topP":             new(0.9),
		"maxTokens":        new(int64(128)),
		"frequencyPenalty": new(0.1),
		"presencePenalty":  new(0.2),
		"reasoningEffort":  "xhigh",
	}, conf, "encoded configuration")

	decoded, err := modules.Registry.Decode(module.Name(), conf)
	require.NoError(t, err, "decode")

	// Images, ImageProperties and ReturnMetadata are query-time options,
	// so they are not part of the collection config.
	want := module
	want.Images = nil
	want.ImageProperties = nil
	want.ReturnMetadata = false
	assert.EqualExportedValues(t, want, decoded, "decoded module")
}
