package modules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/modules"
	"github.com/weaviate/weaviate-go-client/v6/modules/aws"
	"github.com/weaviate/weaviate-go-client/v6/modules/cohere"
	"github.com/weaviate/weaviate-go-client/v6/modules/databricks"
	"github.com/weaviate/weaviate-go-client/v6/modules/digitalocean"
	"github.com/weaviate/weaviate-go-client/v6/modules/google"
	"github.com/weaviate/weaviate-go-client/v6/modules/huggingface"
	"github.com/weaviate/weaviate-go-client/v6/modules/jinaai"
	"github.com/weaviate/weaviate-go-client/v6/modules/mistral"
	"github.com/weaviate/weaviate-go-client/v6/modules/model2vec"
	"github.com/weaviate/weaviate-go-client/v6/modules/morph"
	"github.com/weaviate/weaviate-go-client/v6/modules/nvidia"
	"github.com/weaviate/weaviate-go-client/v6/modules/ollama"
	"github.com/weaviate/weaviate-go-client/v6/modules/openai"
	"github.com/weaviate/weaviate-go-client/v6/modules/selfprovided"
	"github.com/weaviate/weaviate-go-client/v6/modules/transformers"
	"github.com/weaviate/weaviate-go-client/v6/modules/voyageai"
	"github.com/weaviate/weaviate-go-client/v6/modules/weaviate"
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
				"inferenceUrl": "example.com",
				"properties":   []string{"title", "lyrics"},
			},
		},
		{
			name:   "text2vec-model2vec",
			module: model2vec.Text2Vec{},
			conf:   map[string]any{},
		},
		{
			name: "text2vec-weaviate",
			module: weaviate.Text2Vec{
				URL:        "example.com",
				Properties: []string{"title", "lyrics"},
				Model:      weaviate.SnowflakeArcticEmbedMv1_5,
				Dimensions: 92,
			},
			conf: map[string]any{
				"baseURL":    "example.com",
				"properties": []string{"title", "lyrics"},
				"model":      "Snowflake/snowflake-arctic-embed-m-v1.5",
				"dimensions": 92,
			},
		},
		{
			name:   "text2vec-weaviate",
			module: weaviate.Text2Vec{},
			conf:   map[string]any{},
		},
		{
			name: "text2vec-openai",
			module: openai.Text2Vec{
				Model:        "text-embedding-3-large",
				Dimensions:   1024,
				ModelType:    openai.TextModel,
				ModelVersion: "3",
				BaseURL:      "https://proxy.example.com",
				Endpoint:     "/v2/embeddings",
				ResourceName: "my-resource",
				DeploymentID: "my-deployment",
				Properties:   []string{"title", "lyrics"},
			},
			conf: map[string]any{
				"model":        "text-embedding-3-large",
				"dimensions":   1024,
				"type":         openai.TextModel,
				"modelVersion": "3",
				"baseURL":      "https://proxy.example.com",
				"endpoint":     "/v2/embeddings",
				"resourceName": "my-resource",
				"deploymentId": "my-deployment",
				"properties":   []string{"title", "lyrics"},
			},
		},
		{
			name:   "text2vec-openai",
			module: openai.Text2Vec{Dimensions: 256},
			conf:   map[string]any{"dimensions": 256},
		},
		{
			name:   "text2vec-openai",
			module: openai.Text2Vec{},
			conf:   map[string]any{},
		},
		{
			name: "text2vec-google",
			module: google.Text2Vec{
				APIEndpoint:   "us-east1-aiplatform.googleapis.com",
				ProjectID:     "my-project",
				Model:         "gemini-embedding-001",
				Location:      "us-east1",
				Dimensions:    1536,
				TaskType:      google.SemanticSimilarity,
				TitleProperty: "title",
				Properties:    []string{"title", "lyrics"},
			},
			conf: map[string]any{
				"apiEndpoint":   "us-east1-aiplatform.googleapis.com",
				"projectId":     "my-project",
				"model":         "gemini-embedding-001",
				"location":      "us-east1",
				"dimensions":    1536,
				"taskType":      google.SemanticSimilarity,
				"titleProperty": "title",
				"properties":    []string{"title", "lyrics"},
			},
		},
		{
			name:   "text2vec-google",
			module: google.Text2Vec{},
			conf:   map[string]any{},
		},
		{
			name: "text2vec-huggingface",
			module: huggingface.Text2Vec{
				Model:        "sentence-transformers/all-MiniLM-L6-v2",
				EndpointURL:  "https://my-endpoint.huggingface.cloud",
				WaitForModel: testkit.Ptr(true),
				UseGPU:       testkit.Ptr(true),
				UseCache:     testkit.Ptr(false),
				Properties:   []string{"title", "lyrics"},
			},
			conf: map[string]any{
				"model":       "sentence-transformers/all-MiniLM-L6-v2",
				"endpointURL": "https://my-endpoint.huggingface.cloud",
				"options": map[string]any{
					"waitForModel": testkit.Ptr(true),
					"useGPU":       testkit.Ptr(true),
					"useCache":     testkit.Ptr(false),
				},
				"properties": []string{"title", "lyrics"},
			},
		},
		{
			name:   "text2vec-huggingface",
			module: huggingface.Text2Vec{UseCache: testkit.Ptr(false)},
			conf: map[string]any{
				"options": map[string]any{"useCache": testkit.Ptr(false)},
			},
		},
		{
			name:   "text2vec-huggingface",
			module: huggingface.Text2Vec{},
			conf:   map[string]any{},
		},
		{
			name: "text2vec-cohere",
			module: cohere.Text2Vec{
				Model:      "embed-english-v3.0",
				Truncate:   cohere.TruncateStart,
				BaseURL:    "https://proxy.example.com",
				Dimensions: 512,
				Properties: []string{"title", "lyrics"},
			},
			conf: map[string]any{
				"model":      "embed-english-v3.0",
				"truncate":   cohere.TruncateStart,
				"baseURL":    "https://proxy.example.com",
				"dimensions": 512,
				"properties": []string{"title", "lyrics"},
			},
		},
		{
			name:   "text2vec-cohere",
			module: cohere.Text2Vec{},
			conf:   map[string]any{},
		},
		{
			name: "text2vec-jinaai",
			module: jinaai.Text2Vec{
				Model:      "jina-embeddings-v3",
				BaseURL:    "https://proxy.example.com",
				Dimensions: 256,
				Properties: []string{"title", "lyrics"},
			},
			conf: map[string]any{
				"model":      "jina-embeddings-v3",
				"baseURL":    "https://proxy.example.com",
				"dimensions": 256,
				"properties": []string{"title", "lyrics"},
			},
		},
		{
			name:   "text2vec-jinaai",
			module: jinaai.Text2Vec{},
			conf:   map[string]any{},
		},
		{
			name: "text2vec-voyageai",
			module: voyageai.Text2Vec{
				Model:      "voyage-3-large",
				Truncate:   testkit.Ptr(false),
				BaseURL:    "https://proxy.example.com",
				Dimensions: 1024,
				Properties: []string{"title", "lyrics"},
			},
			conf: map[string]any{
				"model":      "voyage-3-large",
				"truncate":   testkit.Ptr(false),
				"baseURL":    "https://proxy.example.com",
				"dimensions": 1024,
				"properties": []string{"title", "lyrics"},
			},
		},
		{
			name:   "text2vec-voyageai",
			module: voyageai.Text2Vec{},
			conf:   map[string]any{},
		},
		{
			name: "text2vec-mistral",
			module: mistral.Text2Vec{
				Model:      "mistral-embed",
				BaseURL:    "https://proxy.example.com/v1/embeddings",
				Properties: []string{"title", "lyrics"},
			},
			conf: map[string]any{
				"model":      "mistral-embed",
				"baseURL":    "https://proxy.example.com/v1/embeddings",
				"properties": []string{"title", "lyrics"},
			},
		},
		{
			name:   "text2vec-mistral",
			module: mistral.Text2Vec{},
			conf:   map[string]any{},
		},
		{
			name: "text2vec-nvidia",
			module: nvidia.Text2Vec{
				Model:      "nvidia/nv-embedqa-e5-v5",
				Truncate:   nvidia.TruncateEnd,
				BaseURL:    "https://proxy.example.com",
				Properties: []string{"title", "lyrics"},
			},
			conf: map[string]any{
				"model":      "nvidia/nv-embedqa-e5-v5",
				"truncate":   nvidia.TruncateEnd,
				"baseURL":    "https://proxy.example.com",
				"properties": []string{"title", "lyrics"},
			},
		},
		{
			name:   "text2vec-nvidia",
			module: nvidia.Text2Vec{},
			conf:   map[string]any{},
		},
		{
			name: "text2vec-aws",
			module: aws.Text2Vec{
				Service:       aws.SageMaker,
				Region:        "us-east-1",
				Model:         "amazon.titan-embed-text-v2:0",
				Endpoint:      "my-endpoint",
				TargetModel:   "my-target-model",
				TargetVariant: "my-variant",
				Dimensions:    512,
				Properties:    []string{"title", "lyrics"},
			},
			conf: map[string]any{
				"service":       aws.SageMaker,
				"region":        "us-east-1",
				"model":         "amazon.titan-embed-text-v2:0",
				"endpoint":      "my-endpoint",
				"targetModel":   "my-target-model",
				"targetVariant": "my-variant",
				"dimensions":    512,
				"properties":    []string{"title", "lyrics"},
			},
		},
		{
			name:   "text2vec-aws",
			module: aws.Text2Vec{},
			conf:   map[string]any{},
		},
		{
			name: "text2vec-databricks",
			module: databricks.Text2Vec{
				Endpoint:    "https://my-workspace.cloud.databricks.com/serving-endpoints/my-model/invocations",
				Instruction: "Represent this document for retrieval",
				Properties:  []string{"title", "lyrics"},
			},
			conf: map[string]any{
				"endpoint":    "https://my-workspace.cloud.databricks.com/serving-endpoints/my-model/invocations",
				"instruction": "Represent this document for retrieval",
				"properties":  []string{"title", "lyrics"},
			},
		},
		{
			name:   "text2vec-databricks",
			module: databricks.Text2Vec{},
			conf:   map[string]any{},
		},
		{
			name: "text2vec-digitalocean",
			module: digitalocean.Text2Vec{
				Model:      "qwen3-embedding-0.6b",
				BaseURL:    "https://proxy.example.com",
				Properties: []string{"title", "lyrics"},
			},
			conf: map[string]any{
				"model":      "qwen3-embedding-0.6b",
				"baseURL":    "https://proxy.example.com",
				"properties": []string{"title", "lyrics"},
			},
		},
		{
			name:   "text2vec-digitalocean",
			module: digitalocean.Text2Vec{},
			conf:   map[string]any{},
		},
		{
			name: "text2vec-morph",
			module: morph.Text2Vec{
				Model:      "morph-embedding-v3",
				BaseURL:    "https://proxy.example.com",
				Endpoint:   "/v2/embeddings",
				Properties: []string{"title", "lyrics"},
			},
			conf: map[string]any{
				"model":      "morph-embedding-v3",
				"baseURL":    "https://proxy.example.com",
				"endpoint":   "/v2/embeddings",
				"properties": []string{"title", "lyrics"},
			},
		},
		{
			name:   "text2vec-morph",
			module: morph.Text2Vec{},
			conf:   map[string]any{},
		},
		{
			name: "text2vec-ollama",
			module: ollama.Text2Vec{
				APIEndpoint: "http://host.docker.internal:11434",
				Model:       "mxbai-embed-large",
				Properties:  []string{"title", "lyrics"},
			},
			conf: map[string]any{
				"apiEndpoint": "http://host.docker.internal:11434",
				"model":       "mxbai-embed-large",
				"properties":  []string{"title", "lyrics"},
			},
		},
		{
			name:   "text2vec-ollama",
			module: ollama.Text2Vec{},
			conf:   map[string]any{},
		},
		{
			name: "text2vec-transformers",
			module: transformers.Text2Vec{
				PoolingStrategy:     transformers.CLS,
				InferenceURL:        "http://t2v-transformers:8080",
				PassageInferenceURL: "http://t2v-passage:8080",
				QueryInferenceURL:   "http://t2v-query:8080",
				Dimensions:          384,
				Properties:          []string{"title", "lyrics"},
			},
			conf: map[string]any{
				"poolingStrategy":     transformers.CLS,
				"inferenceUrl":        "http://t2v-transformers:8080",
				"passageInferenceUrl": "http://t2v-passage:8080",
				"queryInferenceUrl":   "http://t2v-query:8080",
				"dimensions":          384,
				"properties":          []string{"title", "lyrics"},
			},
		},
		{
			name:   "text2vec-transformers",
			module: transformers.Text2Vec{},
			conf:   map[string]any{},
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
