package internal_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal"
)

func TestDecodeEncode(t *testing.T) {
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
	err := internal.Decode(song, &s)
	require.NoError(t, err, "decode error")

	require.Equal(t, Song{
		Title:    "Golden Silver Surfer",
		Artist:   "Telebrains",
		Duration: 321,
	}, s, "bad decode result")

	m := make(map[string]any)
	err = internal.Encode(&s, m)
	require.NoError(t, err, "encode err")

	require.Equal(t, song, m, "bad encode result")
}
