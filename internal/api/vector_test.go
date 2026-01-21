package api_test

import (
	"encoding/json"
	"testing"

	"github.com/go-openapi/testify/v2/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
)

func TestVectors_UnmarshalJSON(t *testing.T) {
	t.Run("valid vectors", func(t *testing.T) {
		want := api.Vectors{
			"single": {Name: "single", Single: []float32{1, 2, 3}},
			"multi":  {Name: "multi", Multi: [][]float32{{1, 2, 3}, {1, 2, 3}}},
		}
		data := map[string]any{
			"single": want["single"].Single,
			"multi":  want["multi"].Multi,
		}
		b, err := json.Marshal(data)
		require.NoError(t, err, "marshal input data")

		var got api.Vectors
		err = json.Unmarshal(b, &got)
		require.NoError(t, err, "unmarshal vectors")
		require.Equal(t, want, got, "bad vectors")
	})

	t.Run("invalid vector", func(t *testing.T) {
		data := []byte(`{"letters": ["a", "b", "c"]}`)

		var got api.Vectors
		err := json.Unmarshal(data, &got)
		require.Error(t, err)
	})
}
