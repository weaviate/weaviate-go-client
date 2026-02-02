package api

import (
	"encoding/json"
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

func TestVectors_UnmarshalJSON(t *testing.T) {
	t.Run("valid vectors", func(t *testing.T) {
		want := Vectors{
			"single": {Name: "single", Single: []float32{1, 2, 3}},
			"multi":  {Name: "multi", Multi: [][]float32{{1, 2, 3}, {1, 2, 3}}},
		}
		data := map[string]any{
			"single": want["single"].Single,
			"multi":  want["multi"].Multi,
		}
		b, err := json.Marshal(data)
		require.NoError(t, err, "marshal input data")

		var got Vectors
		err = json.Unmarshal(b, &got)
		require.NoError(t, err, "unmarshal vectors")
		require.Equal(t, want, got, "bad vectors")
	})

	t.Run("invalid vector", func(t *testing.T) {
		data := []byte(`{"letters": ["a", "b", "c"]}`)

		var got Vectors
		err := json.Unmarshal(data, &got)
		require.Error(t, err)
	})

	t.Run("empty vectors map", func(t *testing.T) {
		var got Vectors
		err := json.Unmarshal([]byte(`{}`), &got)
		require.NoError(t, err, "unmarshal vectors")
		require.Nil(t, got, "vectors map was initialized")
	})
}

func TestVectorBytes(t *testing.T) {
	t.Run("trivial single", func(t *testing.T) {
		v := []float32{1, 2, 3}
		want := []byte{
			0x0, 0x0, 0x80, 0x3f,
			0x0, 0x0, 0x0, 0x40,
			0x0, 0x0, 0x40, 0x40,
		}
		got := marshalSingle(v)
		require.Equal(t, want, got, "bad vector bytes")
	})

	t.Run("trivial multi", func(t *testing.T) {
		v := [][]float32{{1, 2, 3}, {1, 2, 3}}
		want := []byte{
			0x3, 0x0, // inner array size, uint16(3)
			0x0, 0x0, 0x80, 0x3f, // first vector
			0x0, 0x0, 0x0, 0x40,
			0x0, 0x0, 0x40, 0x40,
			0x0, 0x0, 0x80, 0x3f, // second vector
			0x0, 0x0, 0x0, 0x40,
			0x0, 0x0, 0x40, 0x40,
		}
		got := marshalMulti(v)
		require.Equal(t, want, got, "bad vector bytes")
	})
}
