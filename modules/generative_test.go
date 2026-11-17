package modules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/generate"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/modules"
	"github.com/weaviate/weaviate-go-client/v6/modules/anthropic"
	"github.com/weaviate/weaviate-go-client/v6/modules/aws"
	"github.com/weaviate/weaviate-go-client/v6/modules/cohere"
	"github.com/weaviate/weaviate-go-client/v6/modules/contextualai"
	"github.com/weaviate/weaviate-go-client/v6/modules/databricks"
	"github.com/weaviate/weaviate-go-client/v6/modules/deepseek"
	"github.com/weaviate/weaviate-go-client/v6/modules/digitalocean"
	"github.com/weaviate/weaviate-go-client/v6/modules/friendliai"
	"github.com/weaviate/weaviate-go-client/v6/modules/google"
	"github.com/weaviate/weaviate-go-client/v6/modules/meta"
	"github.com/weaviate/weaviate-go-client/v6/modules/mistral"
	"github.com/weaviate/weaviate-go-client/v6/modules/nvidia"
	"github.com/weaviate/weaviate-go-client/v6/modules/ollama"
	"github.com/weaviate/weaviate-go-client/v6/modules/openai"
	"github.com/weaviate/weaviate-go-client/v6/modules/xai"
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
			name:   "generative-openai without provider options",
			module: openai.Generative{Model: "gpt-5-mini"},
			want: &proto.GenerativeProvider{
				Kind: &proto.GenerativeProvider_Openai{
					Openai: &proto.GenerativeOpenAI{
						Model:   new("gpt-5-mini"),
						IsAzure: new(false),
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
		{
			name: "generative-anthropic",
			module: anthropic.Generative{
				BaseURL:       "anthropic.com",
				Model:         "claude-haiku-5-5",
				MaxTokens:     new(int64(12)),
				Temperature:   new(0.7),
				TopK:          new(int64(5)),
				TopP:          new(0.9),
				StopSequences: []string{"<enough>"},

				Provider: &anthropic.Provider{
					Images:          []string{"pony.png"},
					ImageProperties: []string{"profile_pic"},
					ReturnMetadata:  true,
				},
			},
			want: &proto.GenerativeProvider{
				ReturnMetadata: true,
				Kind: &proto.GenerativeProvider_Anthropic{
					Anthropic: &proto.GenerativeAnthropic{
						BaseUrl:       new("anthropic.com"),
						Model:         new("claude-haiku-5-5"),
						MaxTokens:     new(int64(12)),
						Temperature:   new(0.7),
						TopK:          new(int64(5)),
						TopP:          new(0.9),
						StopSequences: &proto.TextArray{Values: []string{"<enough>"}},

						Images:          &proto.TextArray{Values: []string{"pony.png"}},
						ImageProperties: &proto.TextArray{Values: []string{"profile_pic"}},
					},
				},
			},
		},
		{
			name:   "generative-anthropic without provider options",
			module: anthropic.Generative{Model: "claude-haiku-5-5"},
			want: &proto.GenerativeProvider{
				Kind: &proto.GenerativeProvider_Anthropic{
					Anthropic: &proto.GenerativeAnthropic{
						Model: new("claude-haiku-5-5"),
					},
				},
			},
		},
		{
			name: "generative-aws",
			module: aws.Generative{
				Service:           aws.SageMaker,
				Region:            "us-east-1",
				BedrockModel:      "anthropic.claude-haiku-5-5",
				SageMakerEndpoint: "my-endpoint",
				SageMakerModel:    "my-target-model",
				SageMakerVariant:  "my-variant",
				MaxTokens:         new(int64(12)),
				StopSequences:     []string{"<enough>"},
				Temperature:       new(0.7),

				Provider: &aws.Provider{
					Images:          []string{"pony.png"},
					ImageProperties: []string{"profile_pic"},
					ReturnMetadata:  true,
				},
			},
			want: &proto.GenerativeProvider{
				ReturnMetadata: true,
				Kind: &proto.GenerativeProvider_Aws{
					Aws: &proto.GenerativeAWS{
						Service:       new("sagemaker"),
						Region:        new("us-east-1"),
						Model:         new("anthropic.claude-haiku-5-5"),
						Endpoint:      new("my-endpoint"),
						TargetModel:   new("my-target-model"),
						TargetVariant: new("my-variant"),
						MaxTokens:     new(int64(12)),
						StopSequences: &proto.TextArray{Values: []string{"<enough>"}},
						Temperature:   new(0.7),

						Images:          &proto.TextArray{Values: []string{"pony.png"}},
						ImageProperties: &proto.TextArray{Values: []string{"profile_pic"}},
					},
				},
			},
		},
		{
			name: "generative-cohere",
			module: cohere.Generative{
				BaseURL:       "cohere.com",
				Model:         "command-a-03-2025",
				MaxTokens:     new(int64(12)),
				StopSequences: []string{"<enough>"},
				Temperature:   new(0.7),
				K:             new(int64(5)),

				Provider: &cohere.Provider{
					P:                new(0.9),
					FrequencyPenalty: new(0.1),
					PresencePenalty:  new(0.2),
					Images:           []string{"pony.png"},
					ImageProperties:  []string{"profile_pic"},
					ReturnMetadata:   true,
				},
			},
			want: &proto.GenerativeProvider{
				ReturnMetadata: true,
				Kind: &proto.GenerativeProvider_Cohere{
					Cohere: &proto.GenerativeCohere{
						BaseUrl:       new("cohere.com"),
						Model:         new("command-a-03-2025"),
						MaxTokens:     new(int64(12)),
						StopSequences: &proto.TextArray{Values: []string{"<enough>"}},
						Temperature:   new(0.7),
						K:             new(int64(5)),

						P:                new(0.9),
						FrequencyPenalty: new(0.1),
						PresencePenalty:  new(0.2),
						Images:           &proto.TextArray{Values: []string{"pony.png"}},
						ImageProperties:  &proto.TextArray{Values: []string{"profile_pic"}},
					},
				},
			},
		},
		{
			name: "generative-contextualai",
			module: contextualai.Generative{
				Model:           "v2",
				MaxNewTokens:    new(int64(12)),
				Temperature:     new(0.7),
				TopP:            new(0.9),
				SystemPrompt:    "Be brief.",
				AvoidCommentary: new(true),
				Knowledge:       []string{"The sky is blue."},

				Provider: &contextualai.Provider{ReturnMetadata: true},
			},
			want: &proto.GenerativeProvider{
				ReturnMetadata: true,
				Kind: &proto.GenerativeProvider_Contextualai{
					Contextualai: &proto.GenerativeContextualAI{
						Model:           new("v2"),
						MaxNewTokens:    new(int64(12)),
						Temperature:     new(0.7),
						TopP:            new(0.9),
						SystemPrompt:    new("Be brief."),
						AvoidCommentary: new(true),
						Knowledge:       &proto.TextArray{Values: []string{"The sky is blue."}},
					},
				},
			},
		},
		{
			name: "generative-databricks",
			module: databricks.Generative{
				Endpoint:    "databricks.com/serving-endpoints/my-model/invocations",
				MaxTokens:   new(int64(12)),
				Temperature: new(0.7),
				TopP:        new(0.9),

				Provider: &databricks.Provider{
					Model:            "my-model",
					N:                new(int64(2)),
					StopSequences:    []string{"<enough>"},
					FrequencyPenalty: new(0.1),
					PresencePenalty:  new(0.2),
					LogProbs:         new(true),
					TopLogProbs:      new(int64(3)),
					ReturnMetadata:   true,
				},
			},
			want: &proto.GenerativeProvider{
				ReturnMetadata: true,
				Kind: &proto.GenerativeProvider_Databricks{
					Databricks: &proto.GenerativeDatabricks{
						Endpoint:    new("databricks.com/serving-endpoints/my-model/invocations"),
						MaxTokens:   new(int64(12)),
						Temperature: new(0.7),
						TopP:        new(0.9),

						Model:            new("my-model"),
						N:                new(int64(2)),
						Stop:             &proto.TextArray{Values: []string{"<enough>"}},
						FrequencyPenalty: new(0.1),
						PresencePenalty:  new(0.2),
						LogProbs:         new(true),
						TopLogProbs:      new(int64(3)),
					},
				},
			},
		},
		{
			name: "generative-deepseek",
			module: deepseek.Generative{
				BaseURL:          "deepseek.com",
				Model:            "deepseek-chat",
				MaxTokens:        new(int64(12)),
				StopSequences:    []string{"<enough>"},
				Temperature:      new(0.7),
				TopP:             new(0.9),
				FrequencyPenalty: new(0.1),
				PresencePenalty:  new(0.2),

				Provider: &deepseek.Provider{ReturnMetadata: true},
			},
			want: &proto.GenerativeProvider{
				ReturnMetadata: true,
				Kind: &proto.GenerativeProvider_Deepseek{
					Deepseek: &proto.GenerativeDeepseek{
						BaseUrl:          new("deepseek.com"),
						Model:            new("deepseek-chat"),
						MaxTokens:        new(int64(12)),
						Stop:             &proto.TextArray{Values: []string{"<enough>"}},
						Temperature:      new(0.7),
						TopP:             new(0.9),
						FrequencyPenalty: new(0.1),
						PresencePenalty:  new(0.2),
					},
				},
			},
		},
		{
			name: "generative-digitalocean",
			module: digitalocean.Generative{
				BaseURL:          "digitalocean.com",
				Model:            "llama3.3-70b-instruct",
				MaxTokens:        new(int64(12)),
				StopSequences:    []string{"<enough>"},
				Temperature:      new(0.7),
				TopP:             new(0.9),
				FrequencyPenalty: new(0.1),
				PresencePenalty:  new(0.2),

				Provider: &digitalocean.Provider{ReturnMetadata: true},
			},
			want: &proto.GenerativeProvider{
				ReturnMetadata: true,
				Kind: &proto.GenerativeProvider_Digitalocean{
					Digitalocean: &proto.GenerativeDigitalOcean{
						BaseUrl:          new("digitalocean.com"),
						Model:            new("llama3.3-70b-instruct"),
						MaxTokens:        new(int64(12)),
						Stop:             &proto.TextArray{Values: []string{"<enough>"}},
						Temperature:      new(0.7),
						TopP:             new(0.9),
						FrequencyPenalty: new(0.1),
						PresencePenalty:  new(0.2),
					},
				},
			},
		},
		{
			name: "generative-friendliai",
			module: friendliai.Generative{
				BaseURL:     "friendli.ai",
				Model:       "meta-llama-3.1-70b-instruct",
				MaxTokens:   new(int64(12)),
				Temperature: new(0.7),

				Provider: &friendliai.Provider{
					N:              new(int64(2)),
					TopP:           new(0.9),
					ReturnMetadata: true,
				},
			},
			want: &proto.GenerativeProvider{
				ReturnMetadata: true,
				Kind: &proto.GenerativeProvider_Friendliai{
					Friendliai: &proto.GenerativeFriendliAI{
						BaseUrl:     new("friendli.ai"),
						Model:       new("meta-llama-3.1-70b-instruct"),
						MaxTokens:   new(int64(12)),
						Temperature: new(0.7),
						N:           new(int64(2)),
						TopP:        new(0.9),
					},
				},
			},
		},
		{
			name: "generative-google",
			module: google.Generative{
				APIEndpoint:      "us-east1-aiplatform.googleapis.com",
				Model:            "gemini-2.5-flash",
				MaxTokens:        new(int64(12)),
				Temperature:      new(0.7),
				TopK:             new(int64(5)),
				TopP:             new(0.9),
				VertexProjectID:  "my-project",
				VertexEndpointID: "my-endpoint",
				VertexRegion:     "us-east1",
				VertexLocation:   "global",

				Provider: &google.Provider{
					StopSequences:    []string{"<enough>"},
					FrequencyPenalty: new(0.1),
					PresencePenalty:  new(0.2),
					Images:           []string{"pony.png"},
					ImageProperties:  []string{"profile_pic"},
					ReturnMetadata:   true,
				},
			},
			want: &proto.GenerativeProvider{
				ReturnMetadata: true,
				Kind: &proto.GenerativeProvider_Google{
					Google: &proto.GenerativeGoogle{
						ApiEndpoint: new("us-east1-aiplatform.googleapis.com"),
						Model:       new("gemini-2.5-flash"),
						MaxTokens:   new(int64(12)),
						Temperature: new(0.7),
						TopK:        new(int64(5)),
						TopP:        new(0.9),
						ProjectId:   new("my-project"),
						EndpointId:  new("my-endpoint"),
						Region:      new("us-east1"),
						Location:    new("global"),

						StopSequences:    &proto.TextArray{Values: []string{"<enough>"}},
						FrequencyPenalty: new(0.1),
						PresencePenalty:  new(0.2),
						Images:           &proto.TextArray{Values: []string{"pony.png"}},
						ImageProperties:  &proto.TextArray{Values: []string{"profile_pic"}},
					},
				},
			},
		},
		{
			name: "generative-meta",
			module: meta.Generative{
				BaseURL:          "meta.ai",
				Model:            "muse-spark-1.2",
				MaxTokens:        new(int64(12)),
				Temperature:      new(0.7),
				TopP:             new(0.9),
				FrequencyPenalty: new(0.1),
				PresencePenalty:  new(0.2),
				ReasoningEffort:  meta.XHighEffort,

				Provider: &meta.Provider{
					Images:          []string{"pony.png"},
					ImageProperties: []string{"profile_pic"},
					ReturnMetadata:  true,
				},
			},
			want: &proto.GenerativeProvider{
				ReturnMetadata: true,
				Kind: &proto.GenerativeProvider_Meta{
					Meta: &proto.GenerativeMeta{
						BaseUrl:          new("meta.ai"),
						Model:            new("muse-spark-1.2"),
						MaxTokens:        new(int64(12)),
						Temperature:      new(0.7),
						TopP:             new(0.9),
						FrequencyPenalty: new(0.1),
						PresencePenalty:  new(0.2),
						ReasoningEffort:  proto.GenerativeMeta_REASONING_EFFORT_XHIGH.Enum(),

						Images:          &proto.TextArray{Values: []string{"pony.png"}},
						ImageProperties: &proto.TextArray{Values: []string{"profile_pic"}},
					},
				},
			},
		},
		{
			name: "generative-mistral",
			module: mistral.Generative{
				BaseURL:     "mistral.ai",
				Model:       "mistral-large-latest",
				MaxTokens:   new(int64(12)),
				Temperature: new(0.7),

				Provider: &mistral.Provider{
					TopP:           new(0.9),
					ReturnMetadata: true,
				},
			},
			want: &proto.GenerativeProvider{
				ReturnMetadata: true,
				Kind: &proto.GenerativeProvider_Mistral{
					Mistral: &proto.GenerativeMistral{
						BaseUrl:     new("mistral.ai"),
						Model:       new("mistral-large-latest"),
						MaxTokens:   new(int64(12)),
						Temperature: new(0.7),
						TopP:        new(0.9),
					},
				},
			},
		},
		{
			name: "generative-nvidia",
			module: nvidia.Generative{
				BaseURL:     "nvidia.com",
				Model:       "nvidia/llama-3.1-nemotron-70b-instruct",
				MaxTokens:   new(int64(12)),
				Temperature: new(0.7),
				TopP:        new(0.9),

				Provider: &nvidia.Provider{ReturnMetadata: true},
			},
			want: &proto.GenerativeProvider{
				ReturnMetadata: true,
				Kind: &proto.GenerativeProvider_Nvidia{
					Nvidia: &proto.GenerativeNvidia{
						BaseUrl:     new("nvidia.com"),
						Model:       new("nvidia/llama-3.1-nemotron-70b-instruct"),
						MaxTokens:   new(int64(12)),
						Temperature: new(0.7),
						TopP:        new(0.9),
					},
				},
			},
		},
		{
			name: "generative-ollama",
			module: ollama.Generative{
				APIEndpoint: "http://host.docker.internal:11434",
				Model:       "llama3.2",

				Provider: &ollama.Provider{
					Temperature:     new(0.7),
					Images:          []string{"pony.png"},
					ImageProperties: []string{"profile_pic"},
					ReturnMetadata:  true,
				},
			},
			want: &proto.GenerativeProvider{
				ReturnMetadata: true,
				Kind: &proto.GenerativeProvider_Ollama{
					Ollama: &proto.GenerativeOllama{
						ApiEndpoint: new("http://host.docker.internal:11434"),
						Model:       new("llama3.2"),
						Temperature: new(0.7),

						Images:          &proto.TextArray{Values: []string{"pony.png"}},
						ImageProperties: &proto.TextArray{Values: []string{"profile_pic"}},
					},
				},
			},
		},
		{
			name: "generative-xai",
			module: xai.Generative{
				BaseURL:     "x.ai",
				Model:       "grok-4",
				MaxTokens:   new(int64(12)),
				Temperature: new(0.7),
				TopP:        new(0.9),

				Provider: &xai.Provider{
					Images:          []string{"pony.png"},
					ImageProperties: []string{"profile_pic"},
					ReturnMetadata:  true,
				},
			},
			want: &proto.GenerativeProvider{
				ReturnMetadata: true,
				Kind: &proto.GenerativeProvider_Xai{
					Xai: &proto.GenerativeXAI{
						BaseUrl:     new("x.ai"),
						Model:       new("grok-4"),
						MaxTokens:   new(int64(12)),
						Temperature: new(0.7),
						TopP:        new(0.9),

						Images:          &proto.TextArray{Values: []string{"pony.png"}},
						ImageProperties: &proto.TextArray{Values: []string{"profile_pic"}},
					},
				},
			},
		},
		{
			name:   "generative-meta without provider options",
			module: meta.Generative{Model: "muse-spark-1.2"},
			want: &proto.GenerativeProvider{
				Kind: &proto.GenerativeProvider_Meta{
					Meta: &proto.GenerativeMeta{Model: new("muse-spark-1.2")},
				},
			},
		},
	}) {
		t.Run(tt.name, func(t *testing.T) {
			if assert.Implements(t, (*generate.Provider)(nil), tt.module) {
				p := tt.module.(generate.Provider)
				got := p.GenerativeProvider()
				require.EqualExportedValues(t, tt.want, got)
			}
		})
	}
}
