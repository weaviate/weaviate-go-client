package api

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/rest"
	"github.com/weaviate/weaviate-go-client/v6/internal/transports"
)

type CreateReplicationRequest struct {
	transports.BaseEndpoint

	Type       ReplicationType
	Collection string // Collection to be replicated.
	Shard      string // Shard to be replicated.
	Source     string // The source node.
	Target     string // The target node.
}

var _ transports.Endpoint = (*CreateReplicationRequest)(nil)

func (CreateReplicationRequest) Method() string  { return http.MethodPost }
func (r *CreateReplicationRequest) Path() string { return "/replication/replicate" }
func (r *CreateReplicationRequest) Body() any {
	return &rest.ReplicateJSONRequestBody{
		Type:       rest.ReplicationReplicateReplicaRequestType(r.Type),
		Collection: r.Collection,
		Shard:      r.Shard,
		SourceNode: r.Source,
		TargetNode: r.Target,
	}
}

type ReplicationType rest.ReplicationReplicateReplicaRequestType

const (
	ReplicationCopy = ReplicationType(rest.ReplicationReplicateReplicaRequestTypeCOPY)
	ReplicationMove = ReplicationType(rest.ReplicationReplicateReplicaRequestTypeMOVE)
)

// GetReplicationRequest retrieves a replication operation by its UUID.
type GetReplicationRequest struct {
	transports.BaseEndpoint

	UUID           uuid.UUID // Replication ID.
	IncludeHistory bool      // Include history of status changes in the reply.
}

func (GetReplicationRequest) Method() string  { return http.MethodGet }
func (r *GetReplicationRequest) Path() string { return "/replication/replicate/" + r.UUID.String() }
func (r *GetReplicationRequest) Query() url.Values {
	return url.Values{
		"includeHistory": {fmt.Sprintf("%t", r.IncludeHistory)},
	}
}

type ListReplicationsRequest struct {
	transports.BaseEndpoint

	Collection     string // Collection to be replicated.
	Shard          string // Shard to be replicated.
	Target         string // The target node.
	IncludeHistory bool   // Include history of status changes in the reply.
}

func (ListReplicationsRequest) Method() string  { return http.MethodGet }
func (r *ListReplicationsRequest) Path() string { return "/replication/replicate/list" }
func (r *ListReplicationsRequest) Query() url.Values {
	v := url.Values{
		"includeHistory": {fmt.Sprintf("%t", r.IncludeHistory)},
	}
	if r.Collection == "" && r.Shard == "" && r.Target == "" {
		return nil
	}
	if r.Collection != "" {
		v.Set("collection", r.Collection)
	}
	if r.Shard != "" {
		v.Set("shard", r.Shard)
	}
	if r.Target != "" {
		v.Set("targetNode", r.Target)
	}
	return v
}

var (
	CancelReplicationRequest     = transports.IdentityEndpoint[uuid.UUID](http.MethodPost, "/replication/replicate/%s/cancel")
	DeleteReplicationRequest     = transports.IdentityEndpoint[uuid.UUID](http.MethodDelete, "/replication/replicate/%s")
	DeleteAllReplicationsRequest = transports.StaticEndpoint(http.MethodDelete, "/replication/replicate")
)
