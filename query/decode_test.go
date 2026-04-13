package query_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/query"
	"github.com/weaviate/weaviate-go-client/v6/types"
)

func TestDecode(t *testing.T) {
	type Song struct {
		Title    string `json:"title"`
		Duration int    `json:"duration_sec"`
		Artist   string `json:"artist"`
	}

	r := query.Result{
		Objects: []query.Object[map[string]any]{
			{
				Object: types.Object[map[string]any]{
					UUID: testkit.UUID,
					Properties: map[string]any{
						"title":        "Golden Silver Surfer",
						"artist":       "Telebrains",
						"duration_sec": 321,
					},
					CreatedAt:     &testkit.Now,
					LastUpdatedAt: &testkit.Now,
				},
				Metadata: query.Metadata{Distance: testkit.Ptr[float32](.22)},
			},
			{
				Object: types.Object[map[string]any]{
					Properties: map[string]any{
						"title":        "Justifier",
						"artist":       "Telebrains",
						"duration_sec": 158,
					},
				},
			},
		},
	}

	t.Run("nil dest", func(t *testing.T) {
		var dest []query.Object[Song]

		err := query.Decode(&r, &dest)
		require.NoError(t, err, "decode error")
		require.Len(t, dest, len(r.Objects))

		for i := range len(r.Objects) {
			mapSurfer, structSurfer := r.Objects[i], dest[i]
			assert.Equal(t, mapSurfer.Properties["title"], structSurfer.Properties.Title, "title")
			assert.Equal(t, mapSurfer.Properties["artist"], structSurfer.Properties.Artist, "artist")
			assert.Equal(t, mapSurfer.Properties["duration_sec"], structSurfer.Properties.Duration, "duration_sec")

			assert.Equal(t, mapSurfer.UUID, structSurfer.UUID, "uuid")
			assert.Equal(t, mapSurfer.CreatedAt, structSurfer.CreatedAt, "created at")
			assert.Equal(t, mapSurfer.LastUpdatedAt, structSurfer.LastUpdatedAt, "last updated at")
		}
	})

	t.Run("reuse dest", func(t *testing.T) {
		dest := []query.Object[Song]{
			// The "old" objects have some other UUIDs,
			// we want to make sure that Decode will overwrite them.
			{Object: types.Object[Song]{UUID: uuid.New()}},
			{Object: types.Object[Song]{UUID: uuid.New()}},
		}

		err := query.Decode(&r, &dest)
		require.NoError(t, err, "decode error")
		require.Len(t, dest, len(r.Objects))

		for i := range len(r.Objects) {
			mapSurfer, structSurfer := r.Objects[i], dest[i]
			assert.Equal(t, mapSurfer.Properties["title"], structSurfer.Properties.Title, "title")
			assert.Equal(t, mapSurfer.Properties["artist"], structSurfer.Properties.Artist, "artist")
			assert.Equal(t, mapSurfer.Properties["duration_sec"], structSurfer.Properties.Duration, "duration_sec")

			assert.Equal(t, mapSurfer.UUID, structSurfer.UUID, "uuid")
		}
	})
}

func TestDecodeGrouped(t *testing.T) {
	type Song struct {
		Title    string `json:"title"`
		Duration int    `json:"duration_sec"`
		Artist   string `json:"artist"`
	}

	// DecodeGrouped only reads from Groups, using Objects only to
	// grow the dest slice appropriately, so it's enough to just allocate it.
	// This test is, in a way, aware of the implementation, to the extent
	// that it simplifies the setup.
	r := query.GroupByResult{
		Groups: map[string]query.Group[map[string]any]{
			"My Thoughts Changed Directions": {
				Name: "My Thoughts Changed Directions",
				Size: 2,
				Objects: []query.GroupObject[map[string]any]{
					{
						BelongsToGroup: "My Thoughts Changed Directions",
						Object: query.Object[map[string]any]{
							Object: types.Object[map[string]any]{
								UUID: testkit.UUID,
								Properties: map[string]any{
									"title":        "Golden Silver Surfer",
									"artist":       "Telebrains",
									"duration_sec": 321,
								},
								CreatedAt:     &testkit.Now,
								LastUpdatedAt: &testkit.Now,
							},
							Metadata: query.Metadata{Distance: testkit.Ptr[float32](.22)},
						},
					},
					{
						BelongsToGroup: "My Thoughts Changed Directions",
						Object: query.Object[map[string]any]{
							Object: types.Object[map[string]any]{
								Properties: map[string]any{
									"title":        "Justifier",
									"artist":       "Telebrains",
									"duration_sec": 158,
								},
							},
						},
					},
				},
			},
			"Tomorrow": {Name: "Tomorrow"}, // empty group to force an edge case
		},
		Objects: make([]query.GroupObject[map[string]any], 2),
	}

	t.Run("nil dest", func(t *testing.T) {
		var dest []query.GroupObject[Song]

		groups, err := query.DecodeGrouped(&r, &dest)
		require.NoError(t, err, "decode error")

		assert.Len(t, groups, len(r.Groups), "no. groups")
		assert.Len(t, dest, len(r.Objects), "no. objects")
		assert.Contains(t, groups, "My Thoughts Changed Directions")
		assert.Contains(t, groups, "Tomorrow")

		got := groups["My Thoughts Changed Directions"]
		want := r.Groups["My Thoughts Changed Directions"]
		require.Len(t, got.Objects, len(want.Objects), "objects in the group")

		for i := range len(want.Objects) {
			mapSurfer, structSurfer := want.Objects[i], got.Objects[i]
			assert.Equal(t, mapSurfer.Properties["title"], structSurfer.Properties.Title, "title")
			assert.Equal(t, mapSurfer.Properties["artist"], structSurfer.Properties.Artist, "artist")
			assert.Equal(t, mapSurfer.Properties["duration_sec"], structSurfer.Properties.Duration, "duration_sec")

			assert.Equal(t, mapSurfer.UUID, structSurfer.UUID, "uuid")
		}
	})
}
