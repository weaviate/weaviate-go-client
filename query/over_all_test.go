package query_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/query"
	"github.com/weaviate/weaviate-go-client/v6/query/filter"
	"github.com/weaviate/weaviate-go-client/v6/types"
)

func TestOverAll(t *testing.T) {
	rd := api.RequestDefaults{
		CollectionName:   "Songs",
		Tenant:           "john_doe",
		ConsistencyLevel: api.ConsistencyLevelQuorum,
	}

	for _, tt := range testkit.WithOnly(t, []struct {
		testkit.Only

		name  string
		nv    query.OverAll
		stubs []testkit.Stub[api.SearchRequest, api.SearchResponse]
		want  *query.Result // Expected return value.
		err   testkit.Error
	}{
		{
			name: "request ok",
			nv: query.OverAll{
				Limit:     1,
				Offset:    2,
				AutoLimit: 3,
				After:     testkit.UUID,
				Filter: filter.Not{
					P: &filter.Cond{
						Target:   "album",
						Operator: filter.Like,
						Value:    ".*Blood",
					},
				},
				ReturnMetadata: query.ReturnMetadata{
					CreatedAt:    true,
					LastUpdateAt: true,
				},
				ReturnProperties: []string{"title", "duration_sec", "release_date"},
			},
			stubs: []testkit.Stub[api.SearchRequest, api.SearchResponse]{
				{
					Request: &api.SearchRequest{
						RequestDefaults: rd,
						Limit:           1,
						Offset:          2,
						AutoLimit:       3,
						After:           testkit.UUID,
						Filter: api.FilterExpr{
							Operator: api.FilterOperatorNot,
							Exprs: []api.FilterExpr{{
								Operator: api.FilterOperatorLike,
								Target:   []string{"album"},
								Value:    ".*Blood",
							}},
						},
						ReturnMetadata: api.ReturnMetadata{
							CreatedAt:    true,
							LastUpdateAt: true,
						},
						ReturnProperties: []api.ReturnProperty{
							{Name: "title"},
							{Name: "duration_sec"},
							{Name: "release_date"},
						},
					},
					Response: api.SearchResponse{
						Took: 92 * time.Second,
						Results: []api.Object{
							{
								Collection: "Songs",
								Metadata: api.ObjectMetadata{
									UUID:          testkit.UUID,
									CreatedAt:     &testkit.Now,
									LastUpdatedAt: &testkit.Now,
								},
								Properties: map[string]any{
									"title":        "High Speed Dirt",
									"duration_sec": 252,
									"release_date": testkit.Now,
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
								"title":        "High Speed Dirt",
								"duration_sec": 252,
								"release_date": testkit.Now,
							},
							CreatedAt:     &testkit.Now,
							LastUpdatedAt: &testkit.Now,
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

			got, err := c.OverAll(t.Context(), tt.nv)
			tt.err.Require(t, err, "near vector query")
			require.EqualExportedValues(t, tt.want, got, "query result")
		})
	}
}
