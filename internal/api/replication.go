package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/rest"
	"github.com/weaviate/weaviate-go-client/v6/internal/transports"
)

type Replication struct {
	ID              uuid.UUID
	Type            ReplicationType
	Collection      string             // The collection being replicated.
	Shard           string             // The shard being replicated.
	Source          string             // The source node.
	Target          string             // The target node.
	CanCancel       bool               // Whether this replication can be canceled.
	CancelScheduled bool               // Whether this replication is scheduled for cancelation.
	DeleteScheduled bool               // Whether this replication is scheduled for deletion.
	StartedAt       time.Time          // Time at which the replication was initiated.
	Current         ReplicationStage   // Current stage of the replication.
	History         []ReplicationStage // History of previous replication stages.
}

var _ json.Unmarshaler = (*Replication)(nil)

type ReplicationStage struct {
	State     ReplicationState
	Errors    []ReplicationError
	StartedAt time.Time
}

type ReplicationError struct {
	Message string    // Error message.
	Time    time.Time // Time when the error has occurred.
}

// CreateReplicationRequest starts a new replication.
// Use with [Replication] as response destination.
type CreateReplicationRequest struct {
	transports.BaseEndpoint

	Type       ReplicationType // Operation type.
	Collection string          // Collection to be replicated.
	Shard      string          // Shard to be replicated.
	Source     string          // The source node.
	Target     string          // The target node.
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

type ReplicationState rest.ReplicationReplicateDetailsReplicaStatusState

const (
	ReplicationStateCanceled    = ReplicationState(rest.CANCELLED)
	ReplicationStateDehydrating = ReplicationState(rest.DEHYDRATING)
	ReplicationStateFinalizing  = ReplicationState(rest.FINALIZING)
	ReplicationStateHydrating   = ReplicationState(rest.HYDRATING)
	ReplicationStateIntegrating = ReplicationState(rest.INTEGRATING)
	ReplicationStateReady       = ReplicationState(rest.READY)
	ReplicationStateRegistered  = ReplicationState(rest.REGISTERED)
)

type ReplicationType rest.ReplicationReplicateReplicaRequestType

const (
	ReplicationCopy = ReplicationType(rest.ReplicationReplicateReplicaRequestTypeCOPY)
	ReplicationMove = ReplicationType(rest.ReplicationReplicateReplicaRequestTypeMOVE)
)

// GetReplicationRequest retrieves a replication operation by its UUID.
// Use with [Replication] as response destination.
type GetReplicationRequest struct {
	transports.BaseEndpoint

	UUID           uuid.UUID // Replication ID.
	IncludeHistory bool      // Include history of status changes in the reply.
}

var _ transports.Endpoint = (*GetReplicationRequest)(nil)

func (GetReplicationRequest) Method() string  { return http.MethodGet }
func (r *GetReplicationRequest) Path() string { return "/replication/replicate/" + r.UUID.String() }
func (r *GetReplicationRequest) Query() url.Values {
	return url.Values{
		"includeHistory": {fmt.Sprintf("%t", r.IncludeHistory)},
	}
}

// ListReplicationsRequest retrieves all replications, optionally filtered.
// Use with [ListReplicationsResponse].
type ListReplicationsRequest struct {
	transports.BaseEndpoint

	Collection     string // Collection to be replicated.
	Shard          string // Shard to be replicated.
	Target         string // The target node.
	IncludeHistory bool   // Include history of status changes in the reply.
}

var _ transports.Endpoint = (*ListReplicationsRequest)(nil)

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

type ListReplicationsResponse []Replication

var (
	// Cancel a replication by its UUID.
	CancelReplicationRequest = transports.IdentityEndpoint[uuid.UUID](http.MethodPost, "/replication/replicate/%s/cancel")
	// Delete a replication by its UUID.
	DeleteReplicationRequest = transports.IdentityEndpoint[uuid.UUID](http.MethodDelete, "/replication/replicate/%s")
	// Delete all replication operations in a cluster.
	DeleteAllReplicationsRequest = transports.StaticEndpoint(http.MethodDelete, "/replication/replicate")
)

func (r *Replication) UnmarshalJSON(data []byte) error {
	var resp rest.ReplicationReplicateDetailsReplicaResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return err
	}

	var id uuid.UUID
	if resp.Id != nil {
		id = *resp.Id
	}

	stage := func(s rest.ReplicationReplicateDetailsReplicaStatus) ReplicationStage {
		errs := slices.Grow([]ReplicationError(nil), len(s.Errors))
		for _, e := range s.Errors {
			errs = append(errs, ReplicationError{
				Message: e.Message,
				Time:    time.UnixMilli(e.WhenErroredUnixMs),
			})
		}

		return ReplicationStage{
			State:     ReplicationState(s.State),
			StartedAt: time.UnixMilli(s.WhenStartedUnixMs),
			Errors:    errs,
		}
	}

	history := slices.Grow([]ReplicationStage(nil), len(resp.StatusHistory))
	for _, s := range resp.StatusHistory {
		history = append(history, stage(s))
	}

	*r = Replication{
		Type:            ReplicationType(resp.Type),
		ID:              id,
		Collection:      resp.Collection,
		Shard:           resp.Shard,
		Source:          resp.SourceNode,
		Target:          resp.TargetNode,
		CanCancel:       !resp.Uncancelable,
		CancelScheduled: resp.ScheduledForCancel,
		DeleteScheduled: resp.ScheduledForDelete,
		StartedAt:       time.UnixMilli(resp.WhenStartedUnixMs),
		Current:         stage(resp.Status),
		History:         history,
	}
	return nil
}
