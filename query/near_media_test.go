package query_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/query"
	"github.com/weaviate/weaviate-go-client/v6/types"
)

func TestMedia(t *testing.T) {
	for _, tt := range []struct {
		name string
		make func(string) query.Media
		kind api.MediaKind
	}{
		{name: "image", make: query.Image, kind: api.MediaImage},
		{name: "audio", make: query.Audio, kind: api.MediaAudio},
		{name: "video", make: query.Video, kind: api.MediaVideo},
		{name: "depth", make: query.Depth, kind: api.MediaDepth},
		{name: "thermal", make: query.Thermal, kind: api.MediaThermal},
		{name: "imu", make: query.IMU, kind: api.MediaIMU},
	} {
		t.Run(tt.name, func(t *testing.T) {
			require.NotNil(t, tt.make, "bad test: no make function")

			media := tt.make("abc")
			assert.Equal(t, tt.kind, media.Kind(), "media kind")
			assert.Equal(t, "abc", media.Data(), "media data")
		})
	}
}

func TestNearMedia(t *testing.T) {
	rd := api.RequestDefaults{
		CollectionName:   "Songs",
		Tenant:           "john_doe",
		ConsistencyLevel: api.ConsistencyLevelQuorum,
	}

	for _, tt := range testkit.WithOnly(t, []struct {
		testkit.Only

		name  string
		nm    query.NearMedia
		stubs []testkit.Stub[api.SearchRequest, api.SearchResponse]
		want  *query.Result // Expected return value.
		err   testkit.Error
	}{
		{
			name: "default target",
			nm: query.NearMedia{
				Media: query.Image("hounds-mid-race=="),
			},
			stubs: []testkit.Stub[api.SearchRequest, api.SearchResponse]{
				{
					Request: &api.SearchRequest{
						RequestDefaults: rd,
						NearMedia: &api.NearMedia{
							Kind:  api.MediaImage,
							Media: "hounds-mid-race==",
						},
					},
					Response: api.SearchResponse{
						Took: 92 * time.Second,
						Results: []api.Object{
							{
								Collection: "Songs",
								Metadata: api.ObjectMetadata{
									UUID: testkit.UUID,
								},
							},
						},
					},
				},
			},
			want: &query.Result{
				Took: 92 * time.Second,
				Objects: []query.Object[map[string]any]{
					{
						Object: types.Object[map[string]any]{
							Collection: "Songs",
							UUID:       testkit.UUID,
						},
					},
				},
			},
		},
		{
			name: "explicit target",
			nm: query.NearMedia{
				Media:  query.Image("hounds-mid-race=="),
				Target: query.VectorName("album_cover_vec"),
			},
			stubs: []testkit.Stub[api.SearchRequest, api.SearchResponse]{
				{
					Request: &api.SearchRequest{
						RequestDefaults: rd,
						NearMedia: &api.NearMedia{
							Kind:  api.MediaImage,
							Media: "hounds-mid-race==",
							Target: api.SearchTarget{Vectors: []api.TargetVector{
								{Vector: api.Vector{Name: "album_cover_vec"}},
							}},
						},
					},
					Response: api.SearchResponse{
						Took: 92 * time.Second,
						Results: []api.Object{
							{
								Collection: "Songs",
								Metadata: api.ObjectMetadata{
									UUID: testkit.UUID,
								},
							},
						},
					},
				},
			},
			want: &query.Result{
				Took: 92 * time.Second,
				Objects: []query.Object[map[string]any]{
					{
						Object: types.Object[map[string]any]{
							Collection: "Songs",
							UUID:       testkit.UUID,
						},
					},
				},
			},
		},
		{
			name: "request error",
			stubs: []testkit.Stub[api.SearchRequest, api.SearchResponse]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	}) {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)

			c := query.NewClient(transport, rd)
			require.NotNil(t, c, "client")

			got, err := c.NearMedia(t.Context(), tt.nm)
			tt.err.Require(t, err, "near media query")
			require.EqualExportedValues(t, tt.want, got, "query result")
		})
	}
}
