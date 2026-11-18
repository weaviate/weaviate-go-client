package replication

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
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
type Operation struct {
	ID              uuid.UUID
	Type            Type
	Collection      string    // The collection being replicated.
	Shard           string    // The shard being replicated.
	Source          string    // The source node.
	Target          string    // The target node.
	CanCancel       bool      // Whether this replication can be canceled.
	CancelScheduled bool      // Whether this replication is scheduled for cancelation.
	DeleteScheduled bool      // Whether this replication is scheduled for deletion.
	StartedAt       time.Time // Time at which the replication was initiated.
	Current         Stage     // Current stage of the replication.
	History         []Stage   // History of previous replication stages.
}

type Stage struct {
	State     State
	Errors    []Error
	StartedAt time.Time
}

type Error api.ReplicationError

type State api.ReplicationState

const (
	StateCanceled    = State(api.ReplicationStateCanceled)
	StateDehydrating = State(api.ReplicationStateDehydrating)
	StateFinalizing  = State(api.ReplicationStateFinalizing)
	StateHydrating   = State(api.ReplicationStateHydrating)
	StateIntegrating = State(api.ReplicationStateIntegrating)
	StateReady       = State(api.ReplicationStateReady)
	StateRegistered  = State(api.ReplicationStateRegistered)
)

type Type api.ReplicationType

const (
	Copy = Type(api.ReplicationCopy)
	Move = Type(api.ReplicationMove)
)

type (
	MoveOptions CreateOptions
	CopyOptions CreateOptions
)

// Move starts a MOVE operation for the shard data.
func (c *Client) Move(ctx context.Context, options MoveOptions) (*Operation, error) {
	return c.create(ctx, api.ReplicationMove, CreateOptions(options))
}

// Copy starts a COPY operation for the shard data.
func (c *Client) Copy(ctx context.Context, options CopyOptions) (*Operation, error) {
	return c.create(ctx, api.ReplicationCopy, CreateOptions(options))
}

type CreateOptions struct {
	Collection string // Required: Collection to be replicated.
	Shard      string // Required: Shard to be replicated.
	Source     string // Required: The source node.
	Target     string // Required: The target node.
}

// create starts a new replication operation.
func (c *Client) create(ctx context.Context, rt api.ReplicationType, options CreateOptions) (*Operation, error) {
	req := &api.CreateReplicationRequest{
		Type:       rt,
		Collection: options.Collection,
		Shard:      options.Shard,
		Source:     options.Source,
		Target:     options.Target,
	}

	var resp api.Replication
	if err := c.transport.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("create replication: %w", err)
	}
	op := newOperation(&resp)
	return &op, nil
}

type GetOptions struct {
	ID             uuid.UUID // Required: Operation ID.
	IncludeHistory bool      // Include history of status changes in the reply.
}

// Get information about the replication by its ID.
func (c *Client) Get(ctx context.Context, options GetOptions) (*Operation, error) {
	req := &api.GetReplicationRequest{
		ID:             options.ID,
		IncludeHistory: options.IncludeHistory,
	}

	var resp api.Replication
	if err := c.transport.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("get replication: %w", err)
	}

	op := newOperation(&resp)
	return &op, nil
}

type ListOptions struct {
	Collection     string // Collection to be replicated.
	Shard          string // Shard to be replicated.
	Target         string // The target node.
	IncludeHistory bool   // Include history of status changes in the reply.
}

// List replication operations which match the filters.
func (c *Client) List(ctx context.Context, options ListOptions) ([]Operation, error) {
	req := &api.ListReplicationsRequest{
		Collection:     options.Collection,
		Shard:          options.Shard,
		Target:         options.Target,
		IncludeHistory: options.IncludeHistory,
	}

	var resp api.ListReplicationsResponse
	if err := c.transport.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("list replications: %w", err)
	}

	out := slices.Grow([]Operation(nil), len(resp))
	for i := range resp {
		out = append(out, newOperation(&resp[i]))
	}
	return out, nil
}

func newOperation(r *api.Replication) Operation {
	stage := func(s api.ReplicationStage) Stage {
		errs := slices.Grow([]Error(nil), len(s.Errors))
		for i := range s.Errors {
			errs = append(errs, Error(s.Errors[i]))
		}

		return Stage{
			State:     State(s.State),
			StartedAt: s.StartedAt,
			Errors:    errs,
		}
	}

	history := slices.Grow([]Stage(nil), len(r.History))
	for i := range history {
		history = append(history, stage(r.History[i]))
	}

	return Operation{
		Type:            Type(r.Type),
		ID:              r.ID,
		Collection:      r.Collection,
		Shard:           r.Shard,
		Source:          r.Source,
		Target:          r.Target,
		CanCancel:       r.CanCancel,
		CancelScheduled: r.CancelScheduled,
		DeleteScheduled: r.DeleteScheduled,
		StartedAt:       r.StartedAt,
		Current:         stage(r.Current),
		History:         history,
	}
}

// Cancel a replication operation by its ID.
func (c *Client) Cancel(ctx context.Context, id uuid.UUID) error {
	req := api.CancelReplicationRequest(id)

	if err := c.transport.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("cancel replication: %w", err)
	}
	return nil
}

// Delete deletes information about a replication operation by its ID.
func (c *Client) Delete(ctx context.Context, id uuid.UUID) error {
	req := api.DeleteReplicationRequest(id)

	if err := c.transport.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("delete replication: %w", err)
	}
	return nil
}

// DeleteAll deletes information about all replication operations.
func (c *Client) DeleteAll(ctx context.Context) error {
	if err := c.transport.Do(ctx, api.DeleteAllReplicationsRequest, nil); err != nil {
		return fmt.Errorf("delete all replications: %w", err)
	}
	return nil
}
