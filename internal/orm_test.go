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
	})

	t.Run("decode/encode hooks direct", func(t *testing.T) {
		song := map[string]any{"title": "Poison"}

		var s SpecialSong
		err := internal.Decode(song, &s)
		require.NoError(t, err, "decode error")

		require.Equal(t, SpecialSong{
			Title:  "Poison",
			hooked: true,
		}, s, "bad decode result")

		m := make(map[string]any)
		err = internal.Encode(&s, m)
		require.NoError(t, err, "encode err")

		require.Equal(t, map[string]any{
			"title":  "Poison",
			"hooked": "yes",
		}, m, "bad encode result")
	})

	t.Run("decode hook indirect", func(t *testing.T) {
		song := map[string]any{"title": "Poison"}

		var s any = SpecialSong{}
		err := internal.Decode(song, &s)
		require.NoError(t, err, "decode error")

		require.Equal(t, SpecialSong{
			Title:  "Poison",
			hooked: true,
		}, s, "bad decode result")
	})
}

type SpecialSong struct {
	Title  string
	hooked bool // Can only be set/read through custom hooks
}

var (
	_ internal.MapDecoder = (*SpecialSong)(nil)
	_ internal.MapEncoder = (*SpecialSong)(nil)
)

func (ss *SpecialSong) DecodeMap(m map[string]any) error {
	*ss = SpecialSong{
		Title:  m["title"].(string),
		hooked: true,
	}
	return nil
}

func (ss *SpecialSong) EncodeMap() (map[string]any, error) {
	return map[string]any{
		"title":  ss.Title,
		"hooked": "yes",
	}, nil
}
