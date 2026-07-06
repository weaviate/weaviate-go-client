package cluster

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
)

func NewClient(t internal.Transport) *Client {
	dev.AssertNotNil(t, "transport")
	return &Client{transport: t}
}

type Client struct {
	transport internal.Transport
}

type ShardReplica api.ShardReplica

type ListShardReplicasOptions struct {
	Collection string // Required parameter.
	Shard      string // Optional: if omitted, all shards are returned.
}

// ListShardReplicas returns a list of shards and all their replicas.
func (c *Client) ListShardReplicas(ctx context.Context, options ListShardReplicasOptions) ([]ShardReplica, error) {
	req := &api.GetShardsRequest{
		Collection: options.Collection,
		Shard:      options.Shard,
	}

	var resp api.GetShardsResponse
	if err := c.transport.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("list shard replicas: %w", err)
	}

	out := make([]ShardReplica, len(resp))
	for i, shard := range resp {
		out[i] = ShardReplica(shard)
	}
	return out, nil
}

type Node struct {
	Name       string
	Status     NodeStatus
	GitHash    string
	Version    string
	Mode       NodeMode
	Stats      NodeStats
	Shards     []Shard
	BatchStats BatchStats
}

type (
	BatchStats api.BatchStats
	NodeStats  api.NodeStats
)

type Shard struct {
	Name                 string
	Collection           string
	ObjectCount          int64
	Compressed           bool
	Loaded               bool
	VectorIndexingStatus string
	VectorQueueLength    int64
	OngoingReplications  []ReplicationStatus
}

type ReplicationStatus struct {
	TargetNode         string
	ObjectsPropagated  int64
	LastIterationStart time.Time
}

type NodeStatus = api.NodeStatus

const (
	NodeStatusHealthy     = NodeStatus(api.NodeStatusHealthy)
	NodeStatusTimeout     = NodeStatus(api.NodeStatusTimeout)
	NodeStatusUnavailable = NodeStatus(api.NodeStatusUnavailable)
	NodeStatusUnhealthy   = NodeStatus(api.NodeStatusUnhealthy)
)

type NodeMode api.NodeMode

const (
	NodeModeReadOnly  = NodeMode(api.NodeModeReadOnly)
	NodeModeReadWrite = NodeMode(api.NodeModeReadWrite)
	NodeModeScaleOut  = NodeMode(api.NodeModeScaleOut)
	NodeModeWriteOnly = NodeMode(api.NodeModeWriteOnly)
)

type ListNodesOptions struct {
	Collection string // Optional: if omitted, all nodes are returned.
	Shard      string // Optional: if set, only nodes containing this shard are returned.
	Verbose    bool   // Optional: set to true to fetch complete node info.
}

// ListNodes retrieves the information about nodes in this cluster.
func (c *Client) ListNodes(ctx context.Context, options ListNodesOptions) ([]Node, error) {
	req := &api.GetNodesRequest{
		Collection: options.Collection,
		Shard:      options.Shard,
		Verbosity:  api.NodeVerbosityMinimal,
	}

	if options.Verbose {
		req.Verbosity = api.NodeVerbosityVerbose
	}

	var resp api.GetNodesResponse
	if err := c.transport.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("list nodes: %w", err)
	}

	out := make([]Node, len(resp))
	for i, node := range resp {
		shards := slices.Grow([]Shard(nil), len(node.Shards))
		for si, s := range node.Shards {
			replications := slices.Grow([]ReplicationStatus(nil), len(s.OngoingReplications))
			for ri, r := range s.OngoingReplications {
				replications[ri] = ReplicationStatus(r)
			}

			shards[si] = Shard{
				Name:                 s.Name,
				Collection:           s.Collection,
				ObjectCount:          s.ObjectCount,
				Compressed:           s.Compressed,
				Loaded:               s.Loaded,
				VectorIndexingStatus: s.VectorIndexingStatus,
				VectorQueueLength:    s.VectorQueueLength,
				OngoingReplications:  replications,
			}
		}

		out[i] = Node{
			Name:       node.Name,
			Status:     NodeStatus(node.Status),
			GitHash:    node.GitHash,
			Version:    node.Version,
			Mode:       NodeMode(node.Mode),
			Stats:      NodeStats(node.Stats),
			Shards:     shards,
			BatchStats: BatchStats(node.BatchStats),
		}
	}
	return out, nil
}
