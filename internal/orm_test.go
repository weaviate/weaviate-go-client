package internal_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal"
)

func TestDecodeEncode(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		type Song struct {
			Title    string `json:"title"`
			Duration int    `json:"duration_sec"`
			Artist   string `json:"artist"`
		}

		song := map[string]any{
			"title":        "Golden Silver Surfer",
			"artist":       "Telebrains",
			"duration_sec": 321,
		}

		var s Song
		err := internal.Decode(song, &s, nil)
		require.NoError(t, err, "decode error")

		require.Equal(t, Song{
			Title:    "Golden Silver Surfer",
			Artist:   "Telebrains",
			Duration: 321,
		}, s, "bad decode result")

		m := make(map[string]any)
		err = internal.Encode(&s, m, nil)
		require.NoError(t, err, "encode err")

		require.Equal(t, song, m, "bad encode result")
	})

	t.Run("decode hook", func(t *testing.T) {
		type Song struct {
			Title  string
			hooked bool // Can only be set/read through custom hooks
		}

		song := map[string]any{"title": "Poison"}

		var s Song
		err := internal.Decode(song, &s, func(from map[string]any) (any, error) {
			require.IsType(t, (map[string]any)(nil), from, "decode source data")
			return Song{
				Title:  from["title"].(string),
				hooked: true,
			}, nil
		})
		require.NoError(t, err, "decode error")

		require.Equal(t, Song{
			Title:  "Poison",
			hooked: true,
		}, s, "bad decode result")

		m := make(map[string]any)
		err = internal.Encode(&s, m, func(from any) (map[string]any, error) {
			require.IsType(t, (*Song)(nil), from, "encode source data")
			return map[string]any{
				"title":  from.(*Song).Title,
				"hooked": "yes",
			}, nil
		})
		require.NoError(t, err, "encode err")

		require.Equal(t, map[string]any{
			"title":  "Poison",
			"hooked": "yes",
		}, m, "bad encode result")
	})
}
