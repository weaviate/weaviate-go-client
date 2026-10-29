package modules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/generate"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/modules"
	"github.com/weaviate/weaviate-go-client/v6/modules/digitalocean"
	proto "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

func TestGenerativeDigitalOcean(t *testing.T) {
	for _, tt := range testkit.WithOnly(t, []struct {
		testkit.Only
		name   string
		module generate.Provider
		want   *proto.GenerativeProvider
	}{
		{
			name: "configured",
			module: digitalocean.Generative{
				BaseURL:          "https://proxy.example.com",
				Model:            "llama-4-maverick",
				Temperature:      new(0.7),
				TopP:             new(0.9),
				MaxTokens:        new(int64(128)),
				FrequencyPenalty: new(0.1),
				PresencePenalty:  new(0.2),
				Stop:             []string{"<end>"},
				ReturnMetadata:   true,
			},
			want: &proto.GenerativeProvider{
				ReturnMetadata: true,
				Kind: &proto.GenerativeProvider_Digitalocean{
					Digitalocean: &proto.GenerativeDigitalOcean{
						BaseUrl:          new("https://proxy.example.com"),
						Model:            new("llama-4-maverick"),
						Temperature:      new(0.7),
						TopP:             new(0.9),
						MaxTokens:        new(int64(128)),
						FrequencyPenalty: new(0.1),
						PresencePenalty:  new(0.2),
						Stop:             &proto.TextArray{Values: []string{"<end>"}},
					},
				},
			},
		},
		{
			name:   "defaults",
			module: digitalocean.Generative{},
			want: &proto.GenerativeProvider{
				Kind: &proto.GenerativeProvider_Digitalocean{
					Digitalocean: &proto.GenerativeDigitalOcean{},
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

func TestGenerativeDigitalOceanConfig(t *testing.T) {
	module := digitalocean.Generative{
		BaseURL:          "https://proxy.example.com",
		Model:            "llama-4-maverick",
		Temperature:      new(0.7),
		TopP:             new(0.9),
		MaxTokens:        new(int64(128)),
		FrequencyPenalty: new(0.1),
		PresencePenalty:  new(0.2),
		Stop:             []string{"<end>"},
		ReturnMetadata:   true,
	}

	conf, err := modules.Registry.Encode(module)
	require.NoError(t, err, "encode")
	assert.Equal(t, map[string]any{
		"baseURL":          "https://proxy.example.com",
		"model":            "llama-4-maverick",
		"temperature":      new(0.7),
		"topP":             new(0.9),
		"maxTokens":        new(int64(128)),
		"frequencyPenalty": new(0.1),
		"presencePenalty":  new(0.2),
		"stop":             []string{"<end>"},
	}, conf, "encoded configuration")

	decoded, err := modules.Registry.Decode(module.Name(), conf)
	require.NoError(t, err, "decode")

	// ReturnMetadata is a query-time option, so it is not part of the collection config.
	want := module
	want.ReturnMetadata = false
	assert.EqualExportedValues(t, want, decoded, "decoded module")
}
