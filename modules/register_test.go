package modules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/modules"
	"github.com/weaviate/weaviate-go-client/v6/modules/aws"
	"github.com/weaviate/weaviate-go-client/v6/modules/bind"
	"github.com/weaviate/weaviate-go-client/v6/modules/clip"
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
	"github.com/weaviate/weaviate-go-client/v6/modules/twelvelabs"
	"github.com/weaviate/weaviate-go-client/v6/modules/voyageai"
	"github.com/weaviate/weaviate-go-client/v6/modules/weaviate"
)

// TestModules ensures that all modules are registered with [modules.Registry]
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
			name: "multi2vec-google",
			module: google.Multi2Vec{
				APIEndpoint: "us-east1-aiplatform.googleapis.com",
				Model:       "multimodalembedding@001",
				Dimensions:  512,
				TextFields:  []string{"title"},
				ImageFields: []string{"cover"},
				VideoFields: []string{"clip"},
				AudioFields: []string{"track"},
				Weights: google.Weights{
					TextFields:  []float32{0.4},
					ImageFields: []float32{0.3},
					VideoFields: []float32{0.2},
					AudioFields: []float32{0.1},
				},
				ProjectID:            "my-project",
				Location:             "us-east1",
				VideoIntervalSeconds: 10,
			},
			conf: map[string]any{
				"apiEndpoint": "us-east1-aiplatform.googleapis.com",
				"model":       "multimodalembedding@001",
				"dimensions":  512,
				"textFields":  []string{"title"},
				"imageFields": []string{"cover"},
				"videoFields": []string{"clip"},
				"audioFields": []string{"track"},
				"weights": map[string]any{
					"textFields":  []float32{0.4},
					"imageFields": []float32{0.3},
					"videoFields": []float32{0.2},
					"audioFields": []float32{0.1},
				},
				"projectId":            "my-project",
				"location":             "us-east1",
				"videoIntervalSeconds": 10,
			},
		},
		{
			name:   "multi2vec-google",
			module: google.Multi2Vec{},
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
			name: "multi2vec-cohere",
			module: cohere.Multi2Vec{
				Model:       "embed-multilingual-v3.0",
				Truncate:    cohere.TruncateEnd,
				BaseURL:     "https://proxy.example.com",
				Dimensions:  1024,
				TextFields:  []string{"title", "lyrics"},
				ImageFields: []string{"cover"},
				Weights: cohere.Weights{
					TextFields:  []float32{0.2, 0.3},
					ImageFields: []float32{0.5},
				},
			},
			conf: map[string]any{
				"model":       "embed-multilingual-v3.0",
				"truncate":    cohere.TruncateEnd,
				"baseURL":     "https://proxy.example.com",
				"dimensions":  1024,
				"textFields":  []string{"title", "lyrics"},
				"imageFields": []string{"cover"},
				"weights": map[string]any{
					"textFields":  []float32{0.2, 0.3},
					"imageFields": []float32{0.5},
				},
			},
		},
		{
			name:   "multi2vec-cohere",
			module: cohere.Multi2Vec{},
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
			name: "multi2vec-jinaai",
			module: jinaai.Multi2Vec{
				BaseURL:     "https://proxy.example.com",
				Model:       "jina-clip-v2",
				Dimensions:  512,
				TextFields:  []string{"title", "lyrics"},
				ImageFields: []string{"cover"},
				Weights: jinaai.Weights{
					TextFields:  []float32{0.2, 0.3},
					ImageFields: []float32{0.5},
				},
			},
			conf: map[string]any{
				"baseURL":     "https://proxy.example.com",
				"model":       "jina-clip-v2",
				"dimensions":  512,
				"textFields":  []string{"title", "lyrics"},
				"imageFields": []string{"cover"},
				"weights": map[string]any{
					"textFields":  []float32{0.2, 0.3},
					"imageFields": []float32{0.5},
				},
			},
		},
		{
			name:   "multi2vec-jinaai",
			module: jinaai.Multi2Vec{},
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
			name: "multi2vec-voyageai",
			module: voyageai.Multi2Vec{
				BaseURL:     "https://proxy.example.com",
				Model:       "voyage-multimodal-3",
				TextFields:  []string{"title"},
				ImageFields: []string{"cover"},
				VideoFields: []string{"clip"},
				Weights: voyageai.Weights{
					TextFields:  []float32{0.5},
					ImageFields: []float32{0.3},
					VideoFields: []float32{0.2},
				},
				Truncate: testkit.Ptr(false),
			},
			conf: map[string]any{
				"baseURL":     "https://proxy.example.com",
				"model":       "voyage-multimodal-3",
				"textFields":  []string{"title"},
				"imageFields": []string{"cover"},
				"videoFields": []string{"clip"},
				"weights": map[string]any{
					"textFields":  []float32{0.5},
					"imageFields": []float32{0.3},
					"videoFields": []float32{0.2},
				},
				"truncate": testkit.Ptr(false),
			},
		},
		{
			name:   "multi2vec-voyageai",
			module: voyageai.Multi2Vec{},
			conf:   map[string]any{},
		},
		{
			name: "multi2vec-twelvelabs",
			module: twelvelabs.Multi2Vec{
				BaseURL:     "https://proxy.example.com",
				Model:       "Marengo-retrieval-2.7",
				TextFields:  []string{"title", "lyrics"},
				ImageFields: []string{"cover"},
				Weights: twelvelabs.Weights{
					TextFields:  []float32{0.2, 0.3},
					ImageFields: []float32{0.5},
				},
			},
			conf: map[string]any{
				"baseURL":     "https://proxy.example.com",
				"model":       "Marengo-retrieval-2.7",
				"textFields":  []string{"title", "lyrics"},
				"imageFields": []string{"cover"},
				"weights": map[string]any{
					"textFields":  []float32{0.2, 0.3},
					"imageFields": []float32{0.5},
				},
			},
		},
		{
			name:   "multi2vec-twelvelabs",
			module: twelvelabs.Multi2Vec{},
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
			name: "multi2vec-nvidia",
			module: nvidia.Multi2Vec{
				BaseURL:     "https://proxy.example.com",
				Model:       "nvidia/nvclip",
				TextFields:  []string{"title", "lyrics"},
				ImageFields: []string{"cover"},
				Weights: nvidia.Weights{
					TextFields:  []float32{0.2, 0.3},
					ImageFields: []float32{0.5},
				},
			},
			conf: map[string]any{
				"baseURL":     "https://proxy.example.com",
				"model":       "nvidia/nvclip",
				"textFields":  []string{"title", "lyrics"},
				"imageFields": []string{"cover"},
				"weights": map[string]any{
					"textFields":  []float32{0.2, 0.3},
					"imageFields": []float32{0.5},
				},
			},
		},
		{
			name:   "multi2vec-nvidia",
			module: nvidia.Multi2Vec{},
			conf:   map[string]any{},
		},
		{
			name: "text2vec-aws",
			module: aws.Text2Vec{
				Service:          aws.SageMaker,
				Region:           "us-east-1",
				Model:            "amazon.titan-embed-text-v2:0",
				Endpoint:         "my-endpoint",
				SageMakerModel:   "my-target-model",
				SageMakerVariant: "my-variant",
				Dimensions:       512,
				Properties:       []string{"title", "lyrics"},
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
			name: "multi2vec-aws",
			module: aws.Multi2Vec{
				Model:       "amazon.titan-embed-image-v1",
				Dimensions:  384,
				TextFields:  []string{"title", "lyrics"},
				ImageFields: []string{"cover"},
				Weights: aws.Weights{
					TextFields:  []float32{0.2, 0.3},
					ImageFields: []float32{0.5},
				},
				Region: "us-east-1",
			},
			conf: map[string]any{
				"model":       "amazon.titan-embed-image-v1",
				"dimensions":  384,
				"textFields":  []string{"title", "lyrics"},
				"imageFields": []string{"cover"},
				"weights": map[string]any{
					"textFields":  []float32{0.2, 0.3},
					"imageFields": []float32{0.5},
				},
				"region": "us-east-1",
			},
		},
		{
			name:   "multi2vec-aws",
			module: aws.Multi2Vec{},
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
			},
			conf: map[string]any{
				"baseURL":          "openai.com",
				"model":            "o3-mini",
				"apiVersion":       "v1",
				"temperature":      new(36.6),
				"topP":             new(11.0),
				"maxTokens":        new(int64(12)),
				"frequencyPenalty": new(13.0),
				"presencePenalty":  new(14.0),
				"reasoningEffort":  "low",
				"verbosity":        "high",
				"resourceName":     "iron-ore",
				"deploymentId":     "azure-123",
			},
		},
		{
			name:   "generative-openai",
			module: openai.Generative{},
			conf:   map[string]any{},
		},
		{
			name: "multi2vec-clip",
			module: clip.Multi2Vec{
				InferenceURL: "http://multi2vec-clip:8080",
				TextFields:   []string{"title", "lyrics"},
				ImageFields:  []string{"cover"},
				Weights: clip.Weights{
					TextFields:  []float32{0.2, 0.3},
					ImageFields: []float32{0.5},
				},
			},
			conf: map[string]any{
				"inferenceUrl": "http://multi2vec-clip:8080",
				"textFields":   []string{"title", "lyrics"},
				"imageFields":  []string{"cover"},
				"weights": map[string]any{
					"textFields":  []float32{0.2, 0.3},
					"imageFields": []float32{0.5},
				},
			},
		},
		{
			name:   "multi2vec-clip",
			module: clip.Multi2Vec{},
			conf:   map[string]any{},
		},
		{
			name: "multi2vec-bind",
			module: bind.Multi2Vec{
				TextFields:    []string{"title"},
				ImageFields:   []string{"cover"},
				AudioFields:   []string{"track"},
				VideoFields:   []string{"clip"},
				IMUFields:     []string{"motion"},
				ThermalFields: []string{"heat"},
				DepthFields:   []string{"depth"},
				Weights: bind.Weights{
					TextFields:    []float32{0.1},
					ImageFields:   []float32{0.1},
					AudioFields:   []float32{0.1},
					VideoFields:   []float32{0.1},
					IMUFields:     []float32{0.2},
					ThermalFields: []float32{0.2},
					DepthFields:   []float32{0.2},
				},
			},
			conf: map[string]any{
				"textFields":    []string{"title"},
				"imageFields":   []string{"cover"},
				"audioFields":   []string{"track"},
				"videoFields":   []string{"clip"},
				"imuFields":     []string{"motion"},
				"thermalFields": []string{"heat"},
				"depthFields":   []string{"depth"},
				"weights": map[string]any{
					"textFields":    []float32{0.1},
					"imageFields":   []float32{0.1},
					"audioFields":   []float32{0.1},
					"videoFields":   []float32{0.1},
					"imuFields":     []float32{0.2},
					"thermalFields": []float32{0.2},
					"depthFields":   []float32{0.2},
				},
			},
		},
		{
			name:   "multi2vec-bind",
			module: bind.Multi2Vec{},
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
