package schema

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate/entities/models"
)

// ShardsUpdater builder object to update all shards of a class
type ShardsUpdater struct {
	connection *connection.Connection
	className  string
	// shardNames []string
	status models.ShardStatus
}

// WithClassName specifies the class to which the shards belong
func (s *ShardsUpdater) WithClassName(className string) *ShardsUpdater {
	s.className = className
	return s
}

// WithStatus specifies the status with which the shards will be updated
func (s *ShardsUpdater) WithStatus(targetStatus string) *ShardsUpdater {
	s.status = models.ShardStatus{Status: targetStatus}
	return s
}

// Do update the status of the shards of the class specified in ShardsUpdater
func (s *ShardsUpdater) Do(ctx context.Context) (UpdateShardsResponse, error) {
	shards, err := getShards(ctx, s.connection, s.className)
	if err != nil {
		return nil, err
	}

	var payload UpdateShardsResponse

	for _, shard := range shards {
		resp, err := updateShard(ctx, s.connection, s.className, shard.Name, s.status)
		payload = append(payload, &UpdateShardResponse{Name: shard.Name, Status: resp.Status})
		if err != nil {
			return payload, err
		}
	}

	return payload, nil
}

type UpdateShardsResponse []*UpdateShardResponse

type UpdateShardResponse struct {
	Name   string
	Status string
}
