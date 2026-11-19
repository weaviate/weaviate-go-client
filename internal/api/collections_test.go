package api_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/rest"
)

// The server reports every compression algorithm in the vector index
// configuration and marks the ones that are not in use with "enabled": false.
// Only the enabled one is the vector's compression.
func TestCollection_UnmarshalJSON_compressionEnabled(t *testing.T) {
	serverConfig := func(enabled string) map[string]any {
		conf := map[string]any{"distance": "cosine"}
		for _, name := range []string{"bq", "pq", "rq", "sq"} {
			conf[name] = map[string]any{
				"enabled":      name == enabled,
				"rescoreLimit": 7,
			}
		}
		return conf
	}

	for _, tt := range []struct {
		name    string
		enabled string
		want    *api.Module
	}{
		{
			name:    "one algorithm enabled",
			enabled: "rq",
			want: &api.Module{
				Name: "rq",
				Conf: map[string]any{"enabled": true, "rescoreLimit": float64(7)},
			},
		},
		{
			name:    "no algorithm enabled",
			enabled: "",
			want:    nil,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// The registry is iterated in a random order;
			// repeat to make a wrong pick visible.
			for range 32 {
				body, err := json.Marshal(&rest.Class{
					Class: "Songs",
					// The server always sends BM25 settings.
					InvertedIndexConfig: rest.InvertedIndexConfig{Bm25: &rest.BM25Config{}},
					VectorConfig: map[string]rest.VectorConfig{
						"title_vec": {
							VectorIndexType:   "hnsw",
							VectorIndexConfig: serverConfig(tt.enabled),
						},
					},
				})
				require.NoError(t, err, "marshal response")

				var got api.Collection
				require.NoError(t, json.Unmarshal(body, &got), "unmarshal collection")

				vc := got.Vectors["title_vec"]
				require.Equal(t, tt.want, vc.Compression, "compression")
				require.Equal(t, &api.Module{
					Name: "hnsw",
					Conf: map[string]any{"distance": "cosine"},
				}, vc.Index, "index config without compression entries")
			}
		})
	}
}

// A vectorizer may be stored without settings, e.g. {"none": null} on a
// collection created over REST. It is still the vector's vectorizer.
func TestCollection_UnmarshalJSON_vectorizerWithoutSettings(t *testing.T) {
	body, err := json.Marshal(&rest.Class{
		Class:               "Songs",
		InvertedIndexConfig: rest.InvertedIndexConfig{Bm25: &rest.BM25Config{}},
		VectorConfig: map[string]rest.VectorConfig{
			"title_vec": {Vectorizer: map[string]any{"none": nil}},
		},
	})
	require.NoError(t, err, "marshal response")

	var got api.Collection
	require.NoError(t, json.Unmarshal(body, &got), "unmarshal collection")
	require.Equal(t, &api.Module{Name: "none"}, got.Vectors["title_vec"].Vectorizer, "vectorizer")
}
