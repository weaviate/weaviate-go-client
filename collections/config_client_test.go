package collections_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/collections"
	"github.com/weaviate/weaviate-go-client/v6/collections/compression"
	"github.com/weaviate/weaviate-go-client/v6/collections/vectorindex"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
)

func TestNewConfigClient(t *testing.T) {
	require.Panics(t, func() {
		collections.NewConfigClient(nil, api.RequestDefaults{
			CollectionName: "Songs",
		})
	}, "nil transport")
}

func TestConfigClient_update(t *testing.T) {
	returnCollection := func(c api.Collection) testkit.Stub[any, any] {
		return testkit.Stub[any, any]{
			Request:  new(api.GetCollectionRequest("Songs")),
			Response: c,
		}
	}

	rd := api.RequestDefaults{
		CollectionName: "Songs",
	}

	for _, tt := range testkit.WithOnly(t, []struct {
		testkit.Only

		name   string
		update func(ctx context.Context, c *collections.ConfigClient) error
		stubs  []testkit.Stub[any, any]
	}{
		{
			name: "update property description",
			update: func(ctx context.Context, c *collections.ConfigClient) error {
				return c.SetPropertyDescription(
					ctx, "album",
					"Songs released together",
				)
			},
			stubs: []testkit.Stub[any, any]{
				returnCollection(api.Collection{
					Name: rd.CollectionName,
					Properties: []api.Property{
						{Name: "album", DataType: api.DataTypeText},
						{Name: "artist", DataType: api.DataTypeText},
					},
				}),
				{
					Request: testkit.Ptr[any](&api.UpdateCollectionConfigRequest{
						Collection: api.Collection{
							Name: rd.CollectionName,
							Properties: []api.Property{
								{
									Name:        "album",
									DataType:    api.DataTypeText,
									Description: "Songs released together",
								},
								{Name: "artist", DataType: api.DataTypeText},
							},
						},
					}),
				},
			},
		},
		{
			name: "update vector index",
			update: func(ctx context.Context, c *collections.ConfigClient) error {
				return c.UpdateVectorConfig(ctx, "lyrics_vec", func(vc *collections.VectorConfig) {
					vc.Compression = compression.BQ{
						RescoreLimit: 92,
					}

					if hfresh, ok := vc.Index.(vectorindex.HFresh); ok {
						hfresh.SearchProbe = 666
						vc.Index = hfresh
					}
				})
			},
			stubs: []testkit.Stub[any, any]{
				returnCollection(api.Collection{
					Name: rd.CollectionName,
					Vectors: map[string]api.VectorConfig{
						"title_vec": {},
						"lyrics_vec": {
							Index: &api.Module{
								Name: "hfresh",
								Conf: map[string]any{
									"searchProbe": 123,
								},
							},
						},
					},
				}),
				{
					Request: testkit.Ptr[any](&api.UpdateCollectionConfigRequest{
						Collection: api.Collection{
							Name: rd.CollectionName,
							Vectors: map[string]api.VectorConfig{
								"title_vec": {},
								"lyrics_vec": {
									Compression: &api.Module{
										Name: "bq",
										Conf: map[string]any{
											"rescoreLimit": 92,
										},
									},
									Index: &api.Module{
										Name: "hfresh",
										Conf: map[string]any{
											"searchProbe": 666,
										},
									},
								},
							},
						},
					}),
				},
			},
		},
	}) {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := collections.NewConfigClient(transport, rd)
			require.NotNil(t, c, "nil client")

			err := tt.update(t.Context(), c)
			require.NoError(t, err, "update collection")
		})
	}
}
