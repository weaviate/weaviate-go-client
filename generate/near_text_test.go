package generate_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/generate"
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
		nt    generate.NearText
		stubs []testkit.Stub[api.SearchRequest, api.SearchResponse]
		want  *generate.Result // Expected return value.
		err   testkit.Error
	}{
		{
			name: "with near text",
			nt: generate.NearText{
				Query: &query.NearText{
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
				Single: &generate.SingleTask{
					Prompt: "What time is it?",
				},
				Grouped: &generate.GroupedTask{
					Prompt:            "What's NOT there?",
					IncludeProperties: []string{"title", "image"},
				},
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
						Generate: api.GenerateRequest{
							Single: &api.SingleTask{
								Prompt: "What time is it?",
							},
							Grouped: &api.GroupedTask{
								Prompt:            "What's NOT there?",
								IncludeProperties: []string{"title", "image"},
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
								Generative: api.TaskOutput{
									Text: "It's hero time!",
								},
							},
						},
						Generative: api.TaskOutput{
							Text: "The spoon.",
						},
					},
				},
			},
			want: &generate.Result{
				Took: 92 * time.Second,
				Objects: []generate.Object[map[string]any]{
					{
						Object: query.Object[map[string]any]{
							Object: types.Object[map[string]any]{
								Collection: "Songs",
								UUID:       testkit.UUID,
								Properties: map[string]any{
									"title":        "Yellow Submarine",
									"duration_sec": 160,
								},
							},
						},
						Generative: generate.TaskOutput{
							Text: "It's hero time!",
						},
					},
				},
				Generative: generate.TaskOutput{
					Text: "The spoon.",
				},
			},
		},
		{
			name: "request error",
			nt:   generate.NearText{Query: &query.NearText{}},
			stubs: []testkit.Stub[api.SearchRequest, api.SearchResponse]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	}) {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)

			c := generate.NewClient(transport, rd)
			require.NotNil(t, c, "client")

			got, err := c.NearText(t.Context(), tt.nt)
			tt.err.Require(t, err, "near text query")
			require.EqualExportedValues(t, tt.want, got, "query result")
		})
	}
}
