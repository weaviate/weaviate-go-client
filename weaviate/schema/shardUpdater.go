package schema

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

// ShardUpdater builder object to update the shard of a class
type ShardUpdater struct {
	connection *connection.Connection
	className  string
	shardName  string
	status     models.ShardStatus
}

// WithClassName specifies the class to which the shard belongs
func (s *ShardUpdater) WithClassName(className string) *ShardUpdater {
	s.className = className
	return s
}

// WithShardName specifies the name of the shard to update
func (s *ShardUpdater) WithShardName(shardName string) *ShardUpdater {
	s.shardName = shardName
	return s
}

// WithStatus specifies the status with which the shard will be updated
func (s *ShardUpdater) WithStatus(targetStatus string) *ShardUpdater {
	s.status = models.ShardStatus{Status: targetStatus}
	return s
}

// Do update the status of the shard specified in ShardsGetter
func (s *ShardUpdater) Do(ctx context.Context) (*models.ShardStatus, error) {
	return updateShard(ctx, s.connection, s.className, s.shardName, s.status)
}

func updateShard(ctx context.Context, conn *connection.Connection, className, shardName string,
	status models.ShardStatus,
) (*models.ShardStatus, error) {
	responseData, err := conn.RunREST(
		ctx, fmt.Sprintf("/schema/%s/shards/%s", className, shardName), http.MethodPut, status)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	if responseData.StatusCode == 200 {
		var shard models.ShardStatus
		decodeErr := responseData.DecodeBodyIntoTarget(&shard)
		return &shard, decodeErr
	}
	return nil, except.NewWeaviateClientError(responseData.StatusCode, string(responseData.Body))
}
