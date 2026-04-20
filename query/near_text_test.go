package query_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/query"
	"github.com/weaviate/weaviate-go-client/v6/types"
)

func TestNearText(t *testing.T) {
	rd := api.RequestDefaults{
		CollectionName:   "Songs",
		Tenant:           "john_doe",
		ConsistencyLevel: api.ConsistencyLevelQuorum,
	}

	for _, tt := range testkit.WithOnly(t, []struct {
		testkit.Only

		name  string
		nt    query.NearText
		stubs []testkit.Stub[api.SearchRequest, api.SearchResponse]
		want  *query.Result // Expected return value.
		err   testkit.Error
	}{
		{
			name: "request ok",
			nt: query.NearText{
				Concepts: []string{"birds", "chestnuts"},
				Target:   query.VectorName("title_vec"),
				MoveTo: &query.Move{
					Force:   .92,
					Objects: []uuid.UUID{uuid.Nil},
				},
				MoveAway: &query.Move{
					Force:    .46,
					Concepts: []string{"train"},
				},
				Selection: query.SelectionMMR(query.MMR{
					Limit:   1,
					Balance: .2,
				}),
			},
			stubs: []testkit.Stub[api.SearchRequest, api.SearchResponse]{
				{
					Request: &api.SearchRequest{
						RequestDefaults: rd,
						NearText: &api.NearText{
							Concepts: []string{"birds", "chestnuts"},
							Target: api.SearchTarget{
								Vectors: []api.TargetVector{{
									Vector: api.Vector{Name: "title_vec"},
								}},
							},
							MoveTo: &api.Move{
								Force:   .92,
								Objects: []uuid.UUID{uuid.Nil},
							},
							MoveAway: &api.Move{
								Force:    .46,
								Concepts: []string{"train"},
							},
							Selection: api.Selection{
								MMR: &api.SelectionMMR{
									Limit:   1,
									Balance: .2,
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
									"title":        "I Like Birds",
									"duration_sec": 151,
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
								"title":        "I Like Birds",
								"duration_sec": 151,
							},
							Vectors:    make(types.Vectors, 0),
							References: make(map[string][]types.Object[map[string]any], 0),
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

			got, err := c.NearText(t.Context(), tt.nt)
			tt.err.Require(t, err, "near vector query")
			require.Equal(t, tt.want, got, "query result")
		})
	}
}
