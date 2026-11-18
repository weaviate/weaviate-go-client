package cluster_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/cluster"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
)

func TestNewClient(t *testing.T) {
	require.Panics(t, func() {
		cluster.NewClient(nil)
	}, "nil transport")
}

func TestClient_ListShardReplicas(t *testing.T) {
	for _, tt := range []struct {
		name    string
		options cluster.ListShardReplicasOptions
		stubs   []testkit.Stub[api.GetShardsRequest, api.GetShardsResponse]
		want    []cluster.ShardReplica
		err     testkit.Error // Expected error.
	}{
		{
			name: "successfully",
			options: cluster.ListShardReplicasOptions{
				Collection: "Songs",
				Shard:      "xyz",
			},
			stubs: []testkit.Stub[api.GetShardsRequest, api.GetShardsResponse]{
				{
					Request: &api.GetShardsRequest{
						Collection: "Songs",
						Shard:      "xyz",
					},
					Response: api.GetShardsResponse{
						{Shard: "abc", Replicas: []string{"abc-1", "abc-2"}},
						{Shard: "xyz", Replicas: []string{"xyz-1", "xyz-2"}},
					},
				},
			},
			want: []cluster.ShardReplica{
				{Shard: "abc", Replicas: []string{"abc-1", "abc-2"}},
				{Shard: "xyz", Replicas: []string{"xyz-1", "xyz-2"}},
			},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[api.GetShardsRequest, api.GetShardsResponse]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := cluster.NewClient(transport)
			require.NotNil(t, c, "nil client")

			got, err := c.ListShardReplicas(t.Context(), tt.options)
			tt.err.Require(t, err, "list-shards error")
			require.EqualExportedValues(t, tt.want, got, "returned shards")
		})
	}
}

func TestClient_ListNodes(t *testing.T) {
	for _, tt := range []struct {
		name    string
		options cluster.ListNodesOptions
		stubs   []testkit.Stub[api.GetNodesRequest, api.GetNodesResponse]
		want    []cluster.Node
		err     testkit.Error // Expected error.
	}{
		{
			name: "successfully",
			options: cluster.ListNodesOptions{
				Collection: "Songs",
				Shard:      "xyz",
				Verbose:    true,
			},
			stubs: []testkit.Stub[api.GetNodesRequest, api.GetNodesResponse]{
				{
					Request: &api.GetNodesRequest{
						Collection: "Songs",
						Shard:      "xyz",
						Verbosity:  api.NodeVerbosityVerbose,
					},
					Response: api.GetNodesResponse{
						{
							Name:    "node-1",
							Status:  api.NodeStatusHealthy,
							GitHash: "5cc3aa3",
							Version: "1.37.3",
							Mode:    api.NodeModeScaleOut,
						},
					},
				},
			},
			want: []cluster.Node{
				{
					Name:    "node-1",
					Status:  cluster.NodeStatusHealthy,
					GitHash: "5cc3aa3",
					Version: "1.37.3",
					Mode:    cluster.NodeModeScaleOut,
				},
			},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[api.GetNodesRequest, api.GetNodesResponse]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := cluster.NewClient(transport)
			require.NotNil(t, c, "nil client")

			got, err := c.ListNodes(t.Context(), tt.options)
			tt.err.Require(t, err, "list-shards error")
			require.EqualExportedValues(t, tt.want, got, "returned nodes")
		})
	}
}
