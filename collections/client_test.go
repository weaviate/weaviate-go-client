package collections_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/collections"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/types"
)

func TestNewClient(t *testing.T) {
	require.Panics(t, func() {
		collections.NewClient(nil)
	}, "nil transport")
}

func TestClient_Use(t *testing.T) {
	c := collections.NewClient(testkit.NopTransport)
	require.NotNil(t, c, "nil client")

	checkNamespaces := func(t *testing.T, h *collections.Handle) {
		t.Helper()
		assert.NotNil(t, h.Aggregate, "nil aggregate namespace")
		assert.NotNil(t, h.Data, "nil data namespace")
		assert.NotNil(t, h.Query, "nil query namespace")
	}

	t.Run("default handle", func(t *testing.T) {
		handle := c.Use("Songs")

		assert.Equal(t, "Songs", handle.CollectionName(), "collection name")
		assert.EqualValues(t, "", handle.ConsistencyLevel(), "consistency level")
		assert.Equal(t, "", handle.Tenant(), "tenant")
		checkNamespaces(t, handle)
	})

	t.Run("use options", func(t *testing.T) {
		handle := c.Use("Songs",
			collections.WithConsistencyLevel(types.ConsistencyLevelQuorum),
			collections.WithTenant("john_doe"))

		assert.Equal(t, "Songs", handle.CollectionName(), "collection name")
		assert.EqualValues(t, types.ConsistencyLevelQuorum, handle.ConsistencyLevel(), "consistency level")
		assert.Equal(t, "john_doe", handle.Tenant(), "tenant")
		checkNamespaces(t, handle)
	})

	t.Run("derive new handle", func(t *testing.T) {
		handle := c.Use("Songs")

		derived := handle.WithOptions(
			collections.WithConsistencyLevel(types.ConsistencyLevelQuorum),
			collections.WithTenant("john_doe"))

		assert.Equal(t, handle.CollectionName(), derived.CollectionName(), "collection name")
		assert.EqualValues(t, types.ConsistencyLevelQuorum, derived.ConsistencyLevel(), "consistency level")
		assert.Equal(t, "john_doe", derived.Tenant(), "tenant")
		checkNamespaces(t, derived)

		assert.EqualValues(t, "", handle.ConsistencyLevel(), "consistency level was modified")
		assert.Equal(t, "", handle.Tenant(), "tenant was modified")
	})
}

func TestClient_Create(t *testing.T) {
	for _, tt := range []struct {
		name       string
		collection collections.Collection // Collection to be created.
		stubs      []testkit.Stub[api.CreateCollectionRequest, api.Collection]
		err        testkit.Error
	}{
		{
			name: "full config",
			collection: collections.Collection{
				Name:        "Songs",
				Description: "My favorite songs",
				Properties: []collections.Property{
					{Name: "title", DataType: collections.DataTypeText},
					{Name: "genres", DataType: collections.DataTypeTextArray},
					{Name: "single", DataType: collections.DataTypeBool},
					{Name: "year", DataType: collections.DataTypeInt},
					{
						Name:              "lyrics",
						DataType:          collections.DataTypeInt,
						Tokenization:      collections.TokenizationTrigram,
						IndexFilterable:   true,
						IndexRangeFilters: true,
						IndexSearchable:   true,
					},
					{
						Name: "metadata", DataType: collections.DataTypeObject,
						NestedProperties: []collections.Property{
							{Name: "duration", DataType: collections.DataTypeNumber},
							{Name: "uploadedTime", DataType: collections.DataTypeDate},
						},
						Tokenization:      collections.TokenizationWhitespace,
						IndexFilterable:   true,
						IndexRangeFilters: true,
						IndexSearchable:   true,
					},
				},
				References: []collections.Reference{
					{
						Name:        "artist",
						Collections: []string{"Singers", "Bands"},
					},
				},
				Sharding: &collections.ShardingConfig{
					DesiredCount:        3,
					DesiredVirtualCount: 150,
					VirtualPerPhysical:  50,
				},
				Replication: &collections.ReplicationConfig{
					AsyncEnabled:     false,
					Factor:           6,
					DeletionStrategy: collections.TimeBasedResolution,
					AsyncReplication: &collections.AsyncReplicationConfig{
						DiffBatchSize:                   1,
						DiffPerNodeTimeout:              2 * time.Second,
						ReplicationConcurrency:          3,
						ReplicationFrequency:            4 * time.Millisecond,
						ReplicationFrequencyPropagating: 5 * time.Millisecond,
						PrePropagationTimeout:           6 * time.Second,
						PropagationConcurrency:          7,
						PropagationBatchSize:            8,
						PropagationLimit:                9,
						PropagationTimeout:              10 * time.Second,
						PropagationDelay:                11 * time.Millisecond,
						HashTreeHeight:                  12,
						NodePingFrequency:               13 * time.Millisecond,
						LoggingFrequency:                14 * time.Second,
					},
				},
				InvertedIndex: &collections.InvertedIndexConfig{
					IndexNullState:         true,
					IndexPropertyLength:    true,
					IndexTimestamps:        true,
					UsingBlockMaxWAND:      true,
					CleanupIntervalSeconds: 92,
					BM25: &collections.BM25Config{
						B:  25,
						K1: 1,
					},
					Stopwords: &collections.StopwordConfig{
						Preset:    "standard-please-stop",
						Additions: []string{"end"},
						Removals:  []string{"terminate"},
					},
				},
				MultiTenancy: &collections.MultiTenancyConfig{
					Enabled:              true,
					AutoTenantActivation: true,
					AutoTenantCreation:   false,
				},
			},
			stubs: []testkit.Stub[api.CreateCollectionRequest, api.Collection]{
				{
					Request: &api.CreateCollectionRequest{
						Collection: api.Collection{
							Name:        "Songs",
							Description: "My favorite songs",
							Properties: []api.Property{
								{Name: "title", DataType: api.DataTypeText},
								{Name: "genres", DataType: api.DataTypeTextArray},
								{Name: "single", DataType: api.DataTypeBool},
								{Name: "year", DataType: api.DataTypeInt},
								{
									Name:              "lyrics",
									DataType:          api.DataTypeInt,
									Tokenization:      api.TokenizationTrigram,
									IndexFilterable:   true,
									IndexRangeFilters: true,
									IndexSearchable:   true,
								},
								{
									Name: "metadata", DataType: api.DataTypeObject,
									NestedProperties: []api.Property{
										{Name: "duration", DataType: api.DataTypeNumber},
										{Name: "uploadedTime", DataType: api.DataTypeDate},
									},
									Tokenization:      api.TokenizationWhitespace,
									IndexFilterable:   true,
									IndexRangeFilters: true,
									IndexSearchable:   true,
								},
							},
							References: []api.ReferenceProperty{
								{
									Name:        "artist",
									Collections: []string{"Singers", "Bands"},
								},
							},
							Sharding: &api.ShardingConfig{
								DesiredCount:        3,
								DesiredVirtualCount: 150,
								VirtualPerPhysical:  50,
							},
							Replication: &api.ReplicationConfig{
								AsyncEnabled:     false,
								Factor:           6,
								DeletionStrategy: api.TimeBasedResolution,
								AsyncReplication: &api.AsyncReplicationConfig{
									DiffBatchSize:                   1,
									DiffPerNodeTimeout:              2 * time.Second,
									ReplicationConcurrency:          3,
									ReplicationFrequency:            4 * time.Millisecond,
									ReplicationFrequencyPropagating: 5 * time.Millisecond,
									PrePropagationTimeout:           6 * time.Second,
									PropagationConcurrency:          7,
									PropagationBatchSize:            8,
									PropagationLimit:                9,
									PropagationTimeout:              10 * time.Second,
									PropagationDelay:                11 * time.Millisecond,
									HashTreeHeight:                  12,
									NodePingFrequency:               13 * time.Millisecond,
									LoggingFrequency:                14 * time.Second,
								},
							},
							InvertedIndex: &api.InvertedIndexConfig{
								IndexNullState:         true,
								IndexPropertyLength:    true,
								IndexTimestamps:        true,
								UsingBlockMaxWAND:      true,
								CleanupIntervalSeconds: 92,
								BM25: &api.BM25Config{
									B:  25,
									K1: 1,
								},
								Stopwords: &api.StopwordConfig{
									Preset:    "standard-please-stop",
									Additions: []string{"end"},
									Removals:  []string{"terminate"},
								},
							},
							MultiTenancy: &api.MultiTenancyConfig{
								Enabled:              true,
								AutoTenantActivation: true,
								AutoTenantCreation:   false,
							},
						},
					},
				},
			},
		},
		{
			name: "partial config",
			collection: collections.Collection{
				Name:        "Songs",
				Description: "My favorite songs",
			},
			stubs: []testkit.Stub[api.CreateCollectionRequest, api.Collection]{
				{
					Request: &api.CreateCollectionRequest{
						Collection: api.Collection{
							Name:        "Songs",
							Description: "My favorite songs",
						},
					},
				},
			},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[api.CreateCollectionRequest, api.Collection]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// Return the exact collection that was passed
			// in the request to simplify the test cases.
			for i, stub := range tt.stubs {
				if stub.Request != nil {
					tt.stubs[i].Response = stub.Request.Collection
				}
			}

			transport := testkit.NewTransport(t, tt.stubs)
			c := collections.NewClient(transport)
			require.NotNil(t, c, "nil client")

			handle, err := c.Create(t.Context(), tt.collection)
			tt.err.Require(t, err, "create error")

			if tt.err == nil {
				require.Equal(t, tt.collection.Name, handle.CollectionName(), "collection handle name")
			} else {
				require.Nil(t, handle, "handle on error")
			}
		})
	}
}

func TestClient_GetConfig(t *testing.T) {
	for _, tt := range []struct {
		name       string
		collection string
		stubs      []testkit.Stub[any, api.Collection]
		want       *collections.Collection
		err        testkit.Error
	}{
		{
			name:       "ok",
			collection: "Songs",
			stubs: []testkit.Stub[any, api.Collection]{
				{
					Request: testkit.Ptr(api.GetCollectionRequest("Songs")),
					Response: api.Collection{
						Name:        "Songs",
						Description: "My favorite songs",
						Properties: []api.Property{
							{Name: "title", DataType: api.DataTypeText},
							{Name: "genres", DataType: api.DataTypeTextArray},
							{Name: "single", DataType: api.DataTypeBool},
							{Name: "year", DataType: api.DataTypeInt},
							{
								Name:              "lyrics",
								DataType:          api.DataTypeInt,
								Tokenization:      api.TokenizationTrigram,
								IndexFilterable:   true,
								IndexRangeFilters: true,
								IndexSearchable:   true,
							},
							{
								Name: "metadata", DataType: api.DataTypeObject,
								NestedProperties: []api.Property{
									{Name: "duration", DataType: api.DataTypeNumber},
									{Name: "uploadedTime", DataType: api.DataTypeDate},
								},
								Tokenization:      api.TokenizationWhitespace,
								IndexFilterable:   true,
								IndexRangeFilters: true,
								IndexSearchable:   true,
							},
						},
						References: []api.ReferenceProperty{
							{
								Name:        "artist",
								Collections: []string{"Singers", "Bands"},
							},
						},
						Sharding: &api.ShardingConfig{
							DesiredCount:        3,
							DesiredVirtualCount: 150,
							VirtualPerPhysical:  50,
						},
						Replication: &api.ReplicationConfig{
							AsyncEnabled:     false,
							Factor:           6,
							DeletionStrategy: api.TimeBasedResolution,
							AsyncReplication: &api.AsyncReplicationConfig{
								DiffBatchSize:                   1,
								DiffPerNodeTimeout:              2 * time.Second,
								ReplicationConcurrency:          3,
								ReplicationFrequency:            4 * time.Millisecond,
								ReplicationFrequencyPropagating: 5 * time.Millisecond,
								PrePropagationTimeout:           6 * time.Second,
								PropagationConcurrency:          7,
								PropagationBatchSize:            8,
								PropagationLimit:                9,
								PropagationTimeout:              10 * time.Second,
								PropagationDelay:                11 * time.Millisecond,
								HashTreeHeight:                  12,
								NodePingFrequency:               13 * time.Millisecond,
								LoggingFrequency:                14 * time.Second,
							},
						},
						InvertedIndex: &api.InvertedIndexConfig{
							IndexNullState:         true,
							IndexPropertyLength:    true,
							IndexTimestamps:        true,
							UsingBlockMaxWAND:      true,
							CleanupIntervalSeconds: 92,
							BM25: &api.BM25Config{
								B:  25,
								K1: 1,
							},
							Stopwords: &api.StopwordConfig{
								Preset:    "standard-please-stop",
								Additions: []string{"end"},
								Removals:  []string{"terminate"},
							},
						},
						MultiTenancy: &api.MultiTenancyConfig{
							Enabled:              true,
							AutoTenantActivation: true,
							AutoTenantCreation:   false,
						},
					},
				},
			},
			want: &collections.Collection{
				Name:        "Songs",
				Description: "My favorite songs",
				Properties: []collections.Property{
					{Name: "title", DataType: collections.DataTypeText},
					{Name: "genres", DataType: collections.DataTypeTextArray},
					{Name: "single", DataType: collections.DataTypeBool},
					{Name: "year", DataType: collections.DataTypeInt},
					{
						Name:              "lyrics",
						DataType:          collections.DataTypeInt,
						Tokenization:      collections.TokenizationTrigram,
						IndexFilterable:   true,
						IndexRangeFilters: true,
						IndexSearchable:   true,
					},
					{
						Name: "metadata", DataType: collections.DataTypeObject,
						NestedProperties: []collections.Property{
							{Name: "duration", DataType: collections.DataTypeNumber},
							{Name: "uploadedTime", DataType: collections.DataTypeDate},
						},
						Tokenization:      collections.TokenizationWhitespace,
						IndexFilterable:   true,
						IndexRangeFilters: true,
						IndexSearchable:   true,
					},
				},
				References: []collections.Reference{
					{
						Name:        "artist",
						Collections: []string{"Singers", "Bands"},
					},
				},
				Sharding: &collections.ShardingConfig{
					DesiredCount:        3,
					DesiredVirtualCount: 150,
					VirtualPerPhysical:  50,
				},
				Replication: &collections.ReplicationConfig{
					AsyncEnabled:     false,
					Factor:           6,
					DeletionStrategy: collections.TimeBasedResolution,
					AsyncReplication: &collections.AsyncReplicationConfig{
						DiffBatchSize:                   1,
						DiffPerNodeTimeout:              2 * time.Second,
						ReplicationConcurrency:          3,
						ReplicationFrequency:            4 * time.Millisecond,
						ReplicationFrequencyPropagating: 5 * time.Millisecond,
						PrePropagationTimeout:           6 * time.Second,
						PropagationConcurrency:          7,
						PropagationBatchSize:            8,
						PropagationLimit:                9,
						PropagationTimeout:              10 * time.Second,
						PropagationDelay:                11 * time.Millisecond,
						HashTreeHeight:                  12,
						NodePingFrequency:               13 * time.Millisecond,
						LoggingFrequency:                14 * time.Second,
					},
				},
				InvertedIndex: &collections.InvertedIndexConfig{
					IndexNullState:         true,
					IndexPropertyLength:    true,
					IndexTimestamps:        true,
					UsingBlockMaxWAND:      true,
					CleanupIntervalSeconds: 92,
					BM25: &collections.BM25Config{
						B:  25,
						K1: 1,
					},
					Stopwords: &collections.StopwordConfig{
						Preset:    "standard-please-stop",
						Additions: []string{"end"},
						Removals:  []string{"terminate"},
					},
				},
				MultiTenancy: &collections.MultiTenancyConfig{
					Enabled:              true,
					AutoTenantActivation: true,
					AutoTenantCreation:   false,
				},
			},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[any, api.Collection]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := collections.NewClient(transport)
			require.NotNil(t, c, "nil client")

			got, err := c.GetConfig(t.Context(), tt.collection)
			tt.err.Require(t, err, "list error")

			require.Equal(t, tt.want, got, "collection config")
		})
	}
}

func TestClient_List(t *testing.T) {
	for _, tt := range []struct {
		name  string
		stubs []testkit.Stub[any, api.ListCollectionsResponse]
		want  []collections.Collection
		err   testkit.Error
	}{
		{
			name: "empty response",
			stubs: []testkit.Stub[any, api.ListCollectionsResponse]{
				{
					Request: testkit.Ptr[any](api.ListCollectionsRequest),
				},
			},
		},
		{
			name: "several collections",
			stubs: []testkit.Stub[any, api.ListCollectionsResponse]{
				{
					Request: testkit.Ptr[any](api.ListCollectionsRequest),
					Response: api.ListCollectionsResponse{
						{Name: "Songs"},
						{Name: "Artists"},
						{Name: "Albums"},
					},
				},
			},
			want: []collections.Collection{
				{Name: "Songs"},
				{Name: "Artists"},
				{Name: "Albums"},
			},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[any, api.ListCollectionsResponse]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := collections.NewClient(transport)
			require.NotNil(t, c, "nil client")

			got, err := c.List(t.Context())
			tt.err.Require(t, err, "list error")

			require.Equal(t, tt.want, got, "collection config")
		})
	}
}

func TestClient_Delete(t *testing.T) {
	transport := testkit.NewTransport(t, []testkit.Stub[any, any]{
		{Request: testkit.Ptr(api.DeleteCollectionRequest("Songs"))},
	})
	c := collections.NewClient(transport)
	require.NotNil(t, c, "nil client")

	err := c.Delete(t.Context(), "Songs")
	require.NoError(t, err, "delete error")
}

func TestClient_DeleteAll(t *testing.T) {
	transport := testkit.NewTransport(t, []testkit.Stub[any, any]{
		{
			Request: testkit.Ptr[any](api.ListCollectionsRequest),
			Response: api.ListCollectionsResponse{
				{Name: "Songs"},
				{Name: "Artists"},
				{Name: "Albums"},
			},
		},
		{Request: testkit.Ptr(api.DeleteCollectionRequest("Songs"))},
		{Request: testkit.Ptr(api.DeleteCollectionRequest("Artists"))},
		{Request: testkit.Ptr(api.DeleteCollectionRequest("Albums"))},
	})
	c := collections.NewClient(transport)
	require.NotNil(t, c, "nil client")

	err := c.DeleteAll(t.Context())
	require.NoError(t, err, "delete all error")
}

func TestClient_Exists(t *testing.T) {
	for _, tt := range []struct {
		name  string
		stubs []testkit.Stub[any, api.ResourceExistsResponse]
		want  bool
		err   testkit.Error
	}{
		{
			name: "exists",
			stubs: []testkit.Stub[any, api.ResourceExistsResponse]{
				{
					Request:  testkit.Ptr(api.GetCollectionRequest("Songs")),
					Response: true,
				},
			},
			want: true,
		},
		{
			name: "not exists",
			stubs: []testkit.Stub[any, api.ResourceExistsResponse]{
				{
					Request:  testkit.Ptr(api.GetCollectionRequest("Songs")),
					Response: false,
				},
			},
			want: false,
		},
		{
			name: "with error",
			stubs: []testkit.Stub[any, api.ResourceExistsResponse]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := collections.NewClient(transport)
			require.NotNil(t, c, "nil client")

			got, err := c.Exists(t.Context(), "Songs")
			tt.err.Require(t, err, "exists error")

			require.Equal(t, tt.want, got, "exists")
		})
	}
}
