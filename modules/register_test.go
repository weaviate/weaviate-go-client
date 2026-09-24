package modules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/modules"
	"github.com/weaviate/weaviate-go-client/v6/modules/google"
	"github.com/weaviate/weaviate-go-client/v6/modules/huggingface"
	"github.com/weaviate/weaviate-go-client/v6/modules/model2vec"
	"github.com/weaviate/weaviate-go-client/v6/modules/openai"
	"github.com/weaviate/weaviate-go-client/v6/modules/selfprovided"
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
				ModelType:    openai.ModelTypeText,
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
				"type":         openai.ModelTypeText,
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
				TaskType:      google.TaskTypeSemanticSimilarity,
				TitleProperty: "title",
				Properties:    []string{"title", "lyrics"},
			},
			conf: map[string]any{
				"apiEndpoint":   "us-east1-aiplatform.googleapis.com",
				"projectId":     "my-project",
				"model":         "gemini-embedding-001",
				"location":      "us-east1",
				"dimensions":    1536,
				"taskType":      google.TaskTypeSemanticSimilarity,
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
				PassageModel: "sentence-transformers/facebook-dpr-ctx_encoder-single-nq-base",
				QueryModel:   "sentence-transformers/facebook-dpr-question_encoder-single-nq-base",
				EndpointURL:  "https://my-endpoint.huggingface.cloud",
				Options: &huggingface.Options{
					WaitForModel: true,
					UseGPU:       true,
					UseCache:     testkit.Ptr(false),
				},
				Properties: []string{"title", "lyrics"},
			},
			conf: map[string]any{
				"passageModel": "sentence-transformers/facebook-dpr-ctx_encoder-single-nq-base",
				"queryModel":   "sentence-transformers/facebook-dpr-question_encoder-single-nq-base",
				"endpointURL":  "https://my-endpoint.huggingface.cloud",
				"options": map[string]any{
					"waitForModel": true,
					"useGPU":       true,
					"useCache":     testkit.Ptr(false),
				},
				"properties": []string{"title", "lyrics"},
			},
		},
		{
			name:   "text2vec-huggingface",
			module: huggingface.Text2Vec{},
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
