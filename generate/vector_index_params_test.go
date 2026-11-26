package generate_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/generate"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/query"
)

// TestVectorIndexParams checks that every generative search type forwards
// the per-request vector index params (near text is covered in its own test).
func TestVectorIndexParams(t *testing.T) {
	rd := api.RequestDefaults{CollectionName: "Songs"}
	params := query.VectorIndexParams{Ef: 256, FilterStrategy: query.FilterStrategyACORN}
	perTarget := map[string]query.VectorIndexParams{"title_vec": {Ef: 64}}
	target := api.SearchTarget{Vectors: []api.TargetVector{{Vector: api.Vector{Name: "title_vec"}}}}
	single := &generate.SinglePrompt{Prompt: "What time is it?"}

	for _, tt := range testkit.WithOnly(t, []struct {
		testkit.Only

		name string
		run  func(context.Context, *generate.Client) error
		want api.SearchRequest
	}{
		{
			name: "near vector",
			run: func(ctx context.Context, c *generate.Client) error {
				_, err := c.NearVector(ctx, generate.NearVector{Single: single, Query: &query.NearVector{
					Target: query.VectorName("title_vec"), VectorIndexParams: params, TargetVectorIndexParams: perTarget,
				}})
				return err
			},
			want: api.SearchRequest{NearVector: &api.NearVector{Target: target}},
		},
		{
			name: "near object",
			run: func(ctx context.Context, c *generate.Client) error {
				_, err := c.NearObject(ctx, generate.NearObject{Single: single, Query: &query.NearObject{
					UUID: testkit.UUID, Target: query.VectorName("title_vec"), VectorIndexParams: params, TargetVectorIndexParams: perTarget,
				}})
				return err
			},
			want: api.SearchRequest{NearObject: &api.NearObject{UUID: testkit.UUID, Target: target}},
		},
		{
			name: "near media",
			run: func(ctx context.Context, c *generate.Client) error {
				_, err := c.NearMedia(ctx, generate.NearMedia{Single: single, Query: &query.NearMedia{
					Media: query.Image("hounds-mid-race=="), Target: query.VectorName("title_vec"),
					VectorIndexParams: params, TargetVectorIndexParams: perTarget,
				}})
				return err
			},
			want: api.SearchRequest{NearMedia: &api.NearMedia{Kind: api.MediaImage, Media: "hounds-mid-race==", Target: target}},
		},
		{
			name: "hybrid",
			run: func(ctx context.Context, c *generate.Client) error {
				_, err := c.Hybrid(ctx, generate.Hybrid{Single: single, Query: &query.Hybrid{
					Query: "yellow submarine", VectorIndexParams: params, TargetVectorIndexParams: perTarget,
				}})
				return err
			},
			want: api.SearchRequest{Hybrid: &api.Hybrid{Query: "yellow submarine"}},
		},
	}) {
		t.Run(tt.name, func(t *testing.T) {
			want := tt.want
			want.RequestDefaults = rd
			want.VectorIndexParams = api.VectorIndexParams{Ef: 256, FilterStrategy: api.FilterStrategyACORN}
			want.TargetVectorIndexParams = map[string]api.VectorIndexParams{"title_vec": {Ef: 64}}
			want.Generate = api.GenerateRequest{Single: &api.SingleTask{Prompt: "What time is it?"}}
			transport := testkit.NewTransport(t, []testkit.Stub[api.SearchRequest, api.SearchResponse]{
				{Request: &want, Response: api.SearchResponse{}},
			})

			require.NoError(t, tt.run(t.Context(), generate.NewClient(transport, rd)))
		})
	}
}
