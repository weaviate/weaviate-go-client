package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/rest"
	"github.com/weaviate/weaviate-go-client/v6/internal/transports"
)

type ShardReplica struct {
	Shard    string   // Shard name.
	Replicas []string // Nodes which contain a replica of this shard.
}

// GetShardsRequest fetches information about collection's shards
// and their replicas.
// Use with [GetShardsResponse].
type GetShardsRequest struct {
	transports.BaseEndpoint
	Collection string // Required parameter.
	Shard      string // Optional: if ommitted, all shards are returned.
}

var _ transports.Endpoint = (*GetShardsRequest)(nil)

func (GetShardsRequest) Method() string { return http.MethodGet }
func (GetShardsRequest) Path() string   { return "/replication/sharding-state" }
func (r *GetShardsRequest) Query() url.Values {
	v := url.Values{
		"collection": {r.Collection},
	}
	if r.Shard != "" {
		v["shard"] = []string{r.Shard}
	}
	return v
}

type GetShardsResponse []ShardReplica

var _ json.Unmarshaler = (*GetShardsResponse)(nil)

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
	BatchStats rest.BatchStats
	NodeStats  rest.NodeStats
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

type NodeStatus = rest.NodeStatusStatus

const (
	NodeStatusHealthy     = NodeStatus(rest.NodeStatusStatusHEALTHY)
	NodeStatusTimeout     = NodeStatus(rest.NodeStatusStatusTIMEOUT)
	NodeStatusUnavailable = NodeStatus(rest.NodeStatusStatusUNAVAILABLE)
	NodeStatusUnhealthy   = NodeStatus(rest.NodeStatusStatusUNHEALTHY)
)

type NodeMode rest.NodeStatusOperationalMode

const (
	NodeModeReadOnly  = NodeMode(rest.ReadOnly)
	NodeModeReadWrite = NodeMode(rest.ReadWrite)
	NodeModeScaleOut  = NodeMode(rest.ScaleOut)
	NodeModeWriteOnly = NodeMode(rest.WriteOnly)
)

type GetNodesRequest struct {
	transports.BaseEndpoint
	Collection string // Optional: if ommitted, all nodes are returned.
	Shard      string // Optional: if set, only nodes containing this shard are returned.
	Verbosity  NodeVerbosity
}

var _ transports.Endpoint = (*GetNodesRequest)(nil)

func (GetNodesRequest) Method() string { return http.MethodGet }
func (r *GetNodesRequest) Path() string {
	p := "/nodes"
	if r.Collection != "" {
		p += "/" + r.Collection
	}
	return p
}

func (r *GetNodesRequest) Query() url.Values {
	if r.Shard == "" && r.Verbosity == "" {
		return nil
	}

	v := make(url.Values)
	if r.Shard != "" {
		v["shardName"] = []string{r.Shard}
	}
	if r.Verbosity != "" {
		v["output"] = []string{string(r.Verbosity)}
	}
	return v
}

type GetNodesResponse []Node

func (r *GetNodesResponse) UnmarshalJSON(data []byte) error {
	var resp rest.NodesStatusResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return err
	}

	*r = make(GetNodesResponse, len(resp.Nodes))

	for i, n := range resp.Nodes {
		shards := make([]Shard, len(n.Shards))
		for si, s := range n.Shards {
			replications := make([]ReplicationStatus, len(s.AsyncReplicationStatus))
			for ri, r := range s.AsyncReplicationStatus {
				replications[ri] = ReplicationStatus{
					TargetNode:         r.TargetNode,
					ObjectsPropagated:  r.ObjectsPropagated,
					LastIterationStart: time.UnixMilli(r.StartDiffTimeUnixMillis),
				}
			}

			shards[si] = Shard{
				Name:                 s.Name,
				Collection:           s.Class,
				ObjectCount:          s.ObjectCount,
				Compressed:           s.Compressed,
				Loaded:               s.Loaded,
				VectorIndexingStatus: s.VectorIndexingStatus,
				VectorQueueLength:    s.VectorQueueLength,
				OngoingReplications:  replications,
			}
		}

		(*r)[i] = Node{
			Name:       n.Name,
			Status:     NodeStatus(n.Status),
			GitHash:    n.GitHash,
			Version:    n.Version,
			Mode:       NodeMode(n.OperationalMode),
			Stats:      NodeStats(n.Stats),
			Shards:     shards,
			BatchStats: BatchStats(n.BatchStats),
		}
	}
	return nil
}

func (r *GetShardsResponse) UnmarshalJSON(data []byte) error {
	var resp rest.ReplicationShardingStateResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return err
	}

	*r = make(GetShardsResponse, len(resp.ShardingState.Shards))
	for i, s := range resp.ShardingState.Shards {
		(*r)[i] = ShardReplica{
			Shard:    s.Shard,
			Replicas: s.Replicas,
		}
	}
	return nil
}
