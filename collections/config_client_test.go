package collections_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/collections"
	"github.com/weaviate/weaviate-go-client/v6/collections/compression"
	"github.com/weaviate/weaviate-go-client/v6/collections/vectorindex"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/query/filter"
)

func TestNewConfigClient(t *testing.T) {
	require.Panics(t, func() {
		collections.NewConfigClient(nil, api.RequestDefaults{
			CollectionName: "Songs",
		})
	}, "nil transport")
}

func TestConfigClient_update(t *testing.T) {
	rd := api.RequestDefaults{
		CollectionName: "Songs",
	}

	for _, tt := range testkit.WithOnly(t, []struct {
		testkit.Only

		name          string
		current, want api.Collection
		update        func(ctx context.Context, c *collections.ConfigClient) error
	}{
		{
			name: "update property description",
			current: api.Collection{
				Name: rd.CollectionName,
				Properties: []api.Property{
					{Name: "album", DataType: api.DataTypeText},
					{Name: "artist", DataType: api.DataTypeText},
				},
			},
			update: func(ctx context.Context, c *collections.ConfigClient) error {
				return c.SetPropertyDescription(
					ctx, "album",
					"Songs released together",
				)
			},
			want: api.Collection{
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
		},
		{
			name: "update vector index",
			current: api.Collection{
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
			},
			update: func(ctx context.Context, c *collections.ConfigClient) error {
				return c.UpdateVectorConfig(ctx, "lyrics_vec", func(vc *collections.VectorConfig) {
					require.NotNil(t, vc, "vector config")

					vc.Compression = compression.BQ{
						RescoreLimit: 92,
					}

					if hfresh, ok := vc.Index.(vectorindex.HFresh); ok {
						hfresh.SearchProbe = 666
						vc.Index = hfresh
					}
				})
			},
			want: api.Collection{
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
		},
		{
			name: "update replication config",
			current: api.Collection{
				Name: rd.CollectionName,
				Replication: &api.ReplicationConfig{
					DeletionStrategy: api.NoAutomatedResolution,
					AsyncReplication: &api.AsyncReplicationConfig{
						PropagationTimeout: time.Hour,
						DiffBatchSize:      92,
					},
				},
			},
			update: func(ctx context.Context, c *collections.ConfigClient) error {
				return c.UpdateReplicationConfig(t.Context(), func(rc *collections.ReplicationConfig) {
					require.NotNil(t, rc, "replication config")

					rc.DeletionStrategy = collections.DeleteOnConflict
					if assert.NotNil(t, rc.AsyncReplication, "async replication config") {
						rc.AsyncReplication.PropagationTimeout = 19 * time.Second
					}
				})
			},
			want: api.Collection{
				Name: rd.CollectionName,
				Replication: &api.ReplicationConfig{
					DeletionStrategy: api.DeleteOnConflict,
					AsyncReplication: &api.AsyncReplicationConfig{
						PropagationTimeout: 19 * time.Second,
						DiffBatchSize:      92,
					},
				},
			},
		},
		{
			name: "update inverted index config",
			current: api.Collection{
				Name: rd.CollectionName,
				InvertedIndex: &api.InvertedIndexConfig{
					IndexTimestamps:     true,
					IndexPropertyLength: false,
				},
			},
			update: func(ctx context.Context, c *collections.ConfigClient) error {
				return c.UpdateInvertedIndexConfig(t.Context(), func(iic *collections.InvertedIndexConfig) {
					require.NotNil(t, iic, "inverted index config")
					iic.IndexPropertyLength = true
				})
			},
			want: api.Collection{
				Name: rd.CollectionName,
				InvertedIndex: &api.InvertedIndexConfig{
					IndexTimestamps:     true,
					IndexPropertyLength: true,
				},
			},
		},
		{
			name: "update object ttl config",
			current: api.Collection{
				Name: rd.CollectionName,
				ObjectTTL: &api.ObjectTTLConfig{
					Enabled:      false,
					DefaultTTL:   5 * time.Second,
					PropertyName: filter.CreatedAt,
				},
			},
			update: func(ctx context.Context, c *collections.ConfigClient) error {
				return c.UpdateObjectTTLConfig(t.Context(), func(ttl *collections.ObjectTTLConfig) {
					require.NotNil(t, ttl, "object TTL config")
					ttl.Enabled = true
					ttl.PropertyName = filter.LastUpdatedAt
				})
			},
			want: api.Collection{
				Name: rd.CollectionName,
				ObjectTTL: &api.ObjectTTLConfig{
					Enabled:      true,
					DefaultTTL:   5 * time.Second,
					PropertyName: filter.LastUpdatedAt,
				},
			},
		},
		{
			name: "update multi-tenancy config",
			current: api.Collection{
				Name: rd.CollectionName,
				MultiTenancy: &api.MultiTenancyConfig{
					Enabled:              true,
					AutoTenantCreation:   false,
					AutoTenantActivation: false,
				},
			},
			update: func(ctx context.Context, c *collections.ConfigClient) error {
				return c.UpdateMultiTenancyConfig(t.Context(), func(mt *collections.MultiTenancyConfig) {
					require.NotNil(t, mt, "multi-tenancy config")
					mt.AutoTenantCreation = true
					mt.AutoTenantActivation = true
				})
			},
			want: api.Collection{
				Name: rd.CollectionName,
				MultiTenancy: &api.MultiTenancyConfig{
					Enabled:              true,
					AutoTenantCreation:   true,
					AutoTenantActivation: true,
				},
			},
		},
	}) {
		t.Run(tt.name, func(t *testing.T) {
			// Expect the client to fetch the current collection config first,
			// and then send the expected update in the next request.
			transport := testkit.NewTransport(t, []testkit.Stub[any, any]{
				{
					Request:  new(api.GetCollectionRequest("Songs")),
					Response: tt.current,
				},
				{
					Request: testkit.Ptr[any](&api.UpdateCollectionConfigRequest{
						Collection: tt.want,
					}),
					Response: tt.current,
				},
			})
			c := collections.NewConfigClient(transport, rd)
			require.NotNil(t, c, "nil client")

			err := tt.update(t.Context(), c)
			require.NoError(t, err, "update collection")
		})
	}
}
