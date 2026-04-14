package data_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/data"
)

func TestEncode(t *testing.T) {
	type Song struct {
		Title    string `json:"title"`
		Duration int    `json:"duration_sec"`
		Artist   string `json:"artist"`
	}

	t.Run("ok", func(t *testing.T) {
		song, err := data.Encode(&Song{
			Title:    "This Is My Bassdrum",
			Artist:   "Telebrains",
			Duration: 202,
		})
		require.NoError(t, err, "encode error")

		assert.Equal(t, map[string]any{
			"title":        "This Is My Bassdrum",
			"artist":       "Telebrains",
			"duration_sec": 202,
		}, song)
	})

	t.Run("invalid input type", func(t *testing.T) {
		var f func()
		song, err := data.Encode(&f)
		require.Error(t, err, "encode error")
		require.Nil(t, song, "result map")
	})
}

func TestMustEncode(t *testing.T) {
	var f func()
	require.Panics(t, func() { data.MustEncode(&f) })
}
