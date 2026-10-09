package modules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/generate"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/modules"
	"github.com/weaviate/weaviate-go-client/v6/modules/openai"
	proto "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

func TestGenerativeProvider(t *testing.T) {
	for _, tt := range testkit.WithOnly(t, []struct {
		testkit.Only
		name   string
		module modules.Module
		want   *proto.GenerativeProvider
	}{
		{
			name: "generative-openai",
			module: openai.Generative{
				BaseURL:          "openai.com",
				Model:            "o3-mini",
				APIVersion:       "v1",
				Temperature:      new(36.6),
				TopP:             new(11.0),
				MaxTokens:        new(int64(12)),
				FrequencyPenalty: new(13.0),
				PresencePenalty:  new(14.0),
				ReasoningEffort:  openai.LowEffort,
				Verbosity:        openai.HighVerbosity,

				ResourceName: "iron-ore",
				DeploymentID: "azure-123",

				Provider: &openai.Provider{
					Images:          []string{"pony.png"},
					ImageProperties: []string{"profile_pic"},
					StopSequences:   []string{"<enough>"},
					N:               new(int64(92)),
					ReturnMetadata:  true,
				},
			},
			want: &proto.GenerativeProvider{
				ReturnMetadata: true,
				Kind: &proto.GenerativeProvider_Openai{
					Openai: &proto.GenerativeOpenAI{
						BaseUrl:    new("openai.com"),
						Model:      new("o3-mini"),
						ApiVersion: new("v1"),

						Temperature:      new(36.6),
						TopP:             new(11.0),
						MaxTokens:        new(int64(12)),
						FrequencyPenalty: new(13.0),
						PresencePenalty:  new(14.0),
						ReasoningEffort:  proto.GenerativeOpenAI_REASONING_EFFORT_LOW.Enum(),
						Verbosity:        proto.GenerativeOpenAI_VERBOSITY_HIGH.Enum(),

						IsAzure:      new(true),
						ResourceName: new("iron-ore"),
						DeploymentId: new("azure-123"),

						N:               new(int64(92)),
						Stop:            &proto.TextArray{Values: []string{"<enough>"}},
						Images:          &proto.TextArray{Values: []string{"pony.png"}},
						ImageProperties: &proto.TextArray{Values: []string{"profile_pic"}},
					},
				},
			},
		},
		{
			name: "generative-openai-azure",
			module: openai.Generative{
				BaseURL:      "openai.com",
				ResourceName: "iron-ore",
				DeploymentID: "azure-123",
				Provider: &openai.Provider{
					ReturnMetadata: false,
				},
			},
			want: &proto.GenerativeProvider{
				ReturnMetadata: false,
				Kind: &proto.GenerativeProvider_Openai{
					Openai: &proto.GenerativeOpenAI{
						BaseUrl:      new("openai.com"),
						IsAzure:      new(true),
						ResourceName: new("iron-ore"),
						DeploymentId: new("azure-123"),
					},
				},
			},
		},
	}) {
		t.Run(tt.name, func(t *testing.T) {
			if assert.Implements(t, (*generate.Provider)(nil), tt.module) {
				p := tt.module.(generate.Provider)
				got := p.GenerativeProvider()
				require.EqualExportedValues(t, tt.want, &got)
			}
		})
	}
}
