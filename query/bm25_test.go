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

func TestKeywordSimilarity(t *testing.T) {
	t.Run("all tokens match", func(t *testing.T) {
		assert.True(t, query.AllTokensMatch.AllTokensMatch(), "all tokens match")
		assert.Nil(t, query.AllTokensMatch.MinimumTokensMatch(), "minimum tokens match")
	})
	t.Run("all tokens match cross property", func(t *testing.T) {
		assert.True(t, query.AllTokensMatchCross.AllTokensMatch(), "all tokens match")
		assert.True(t, query.AllTokensMatchCross.CrossProperty(), "cross property")
		assert.Nil(t, query.AllTokensMatchCross.MinimumTokensMatch(), "minimum tokens match")
	})
	t.Run("minimum tokens match", func(t *testing.T) {
		similarity := query.MinimumTokensMatch(5)
		assert.False(t, similarity.AllTokensMatch(), "all tokens match")
		assert.NotNil(t, similarity.MinimumTokensMatch(), "minimum tokens match")
	})
}

func TestBM25(t *testing.T) {
	rd := api.RequestDefaults{
		CollectionName:   "Songs",
		Tenant:           "john_doe",
		ConsistencyLevel: api.ConsistencyLevelQuorum,
	}

	for _, tt := range testkit.WithOnly(t, []struct {
		testkit.Only

		name  string
		bm25  query.BM25
		stubs []testkit.Stub[api.SearchRequest, api.SearchResponse]
		want  *query.Result // Expected return value.
		err   testkit.Error
	}{
		{
			name: "ok",
			bm25: query.BM25{
				Query:             "yellow submarine",
				QueryProperties:   []string{"title", "lyrics"},
				KeywordSimilarity: query.MinimumTokensMatch(3),
			},
			stubs: []testkit.Stub[api.SearchRequest, api.SearchResponse]{
				{
					Request: &api.SearchRequest{
						RequestDefaults: rd,
						BM25: &api.BM25{
							Query:           "yellow submarine",
							QueryProperties: []string{"title", "lyrics"},
							KeywordSimilarity: api.KeywordSimilarity{
								MinimumTokensMatch: testkit.Ptr[int32](3),
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

			got, err := c.BM25(t.Context(), tt.bm25)
			tt.err.Require(t, err, "bm25 query")
			require.EqualExportedValues(t, tt.want, got, "query result")
		})
	}
}
