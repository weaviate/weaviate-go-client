package query_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/query"
	"github.com/weaviate/weaviate-go-client/v6/types"
)

func TestHybrid(t *testing.T) {
	rd := api.RequestDefaults{
		CollectionName:   "Songs",
		Tenant:           "john_doe",
		ConsistencyLevel: api.ConsistencyLevelQuorum,
	}

	for _, tt := range testkit.WithOnly(t, []struct {
		testkit.Only

		name  string
		nt    query.Hybrid
		stubs []testkit.Stub[api.SearchRequest, api.SearchResponse]
		want  *query.Result // Expected return value.
		err   testkit.Error
	}{
		{
			name: "with near text",
			nt: query.Hybrid{
				Query:             "yellow submarine",
				QueryProperties:   []string{"title", "lyrics"},
				KeywordSimilarity: query.AllTokensMatch,
				Alpha:             query.Alpha(0.44),
				Fusion:            query.HybridFusionRelativeScore,
				NearText: &query.NearText{
					Concepts: []string{"sea"},
					Target:   query.VectorName("title_vec"),
				},
			},
			stubs: []testkit.Stub[api.SearchRequest, api.SearchResponse]{
				{
					Request: &api.SearchRequest{
						RequestDefaults: rd,
						Hybrid: &api.Hybrid{
							Query:           "yellow submarine",
							QueryProperties: []string{"title", "lyrics"},
							KeywordSimilarity: &api.KeywordSimilarity{
								AllTokensMatch: true,
							},
							Alpha:  testkit.Ptr[float32](0.44),
							Fusion: api.HybridFusionRelativeScore,
							NearText: &api.NearText{
								Concepts: []string{"sea"},
								Target: api.SearchTarget{
									Vectors: []api.TargetVector{{
										Vector: api.Vector{Name: "title_vec"},
									}},
								},
							},
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
								Properties: map[string]any{
									"title":        "Yellow Submarine",
									"duration_sec": 160,
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
							Properties: map[string]any{
								"title":        "Yellow Submarine",
								"duration_sec": 160,
							},
						},
					},
				},
			},
		},
		{
			name: "with near vector",
			nt: query.Hybrid{
				Query:             "yellow submarine",
				QueryProperties:   []string{"title", "lyrics"},
				KeywordSimilarity: query.MinimumTokensMatch(2),
				Fusion:            query.HybridFusionRanked,
				NearVector: &query.NearVector{
					Target: &types.Vector{Single: singleVector},
				},
			},
			stubs: []testkit.Stub[api.SearchRequest, api.SearchResponse]{
				{
					Request: &api.SearchRequest{
						RequestDefaults: rd,
						Hybrid: &api.Hybrid{
							Query:           "yellow submarine",
							QueryProperties: []string{"title", "lyrics"},
							KeywordSimilarity: &api.KeywordSimilarity{
								MinimumTokensMatch: testkit.Ptr[int32](2),
							},
							Fusion: api.HybridFusionRanked,
							NearVector: &api.NearVector{
								Target: api.SearchTarget{
									Vectors: []api.TargetVector{{
										Vector: api.Vector{Single: singleVector},
									}},
								},
							},
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
								Properties: map[string]any{
									"title":        "Yellow Submarine",
									"duration_sec": 160,
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
							Properties: map[string]any{
								"title":        "Yellow Submarine",
								"duration_sec": 160,
							},
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

			got, err := c.Hybrid(t.Context(), tt.nt)
			tt.err.Require(t, err, "near vector query")
			require.EqualExportedValues(t, tt.want, got, "query result")
		})
	}
}
