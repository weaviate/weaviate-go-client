package replication_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/replication"
)

func TestNewClient(t *testing.T) {
	require.Panics(t, func() {
		replication.NewClient(nil)
	}, "nil transport")
}

// Move and Copy and both thin wrappers around the CreateReplication
// request, so its enought that we test one of them.
func TestClient_Move(t *testing.T) {
	for _, tt := range []struct {
		name    string
		options replication.MoveOptions
		stubs   []testkit.Stub[api.CreateReplicationRequest, any]
		want    *replication.Operation
		err     testkit.Error // Expected error.
	}{
		{
			name: "successfully",
			options: replication.MoveOptions{
				Collection: "Songs",
				Shard:      "abc",
				Source:     "node-0",
				Target:     "node-1",
			},
			stubs: []testkit.Stub[api.CreateReplicationRequest, any]{
				{
					Request: &api.CreateReplicationRequest{
						Type:       api.ReplicationMove,
						Collection: "Songs",
						Shard:      "abc",
						Source:     "node-0",
						Target:     "node-1",
					},
					Response: api.Replication{
						Type:       api.ReplicationMove,
						ID:         testkit.UUID,
						Collection: "Songs",
						Shard:      "abc",
						Source:     "node-0",
						Target:     "node-1",
						StartedAt:  testkit.Now,
					},
				},
			},
			want: &replication.Operation{
				Type:       replication.Move,
				ID:         testkit.UUID,
				Collection: "Songs",
				Shard:      "abc",
				Source:     "node-0",
				Target:     "node-1",
				StartedAt:  testkit.Now,
			},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[api.CreateReplicationRequest, any]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := replication.NewClient(transport)
			require.NotNil(t, c, "nil client")

			got, err := c.Move(t.Context(), tt.options)
			tt.err.Require(t, err, "move error")
			require.EqualExportedValues(t, tt.want, got, "returned operation")
		})
	}
}

func TestClient_Get(t *testing.T) {
	for _, tt := range []struct {
		name    string
		options replication.GetOptions
		stubs   []testkit.Stub[api.GetReplicationRequest, any]
		want    *replication.Operation
		err     testkit.Error // Expected error.
	}{
		{
			name: "successfully",
			options: replication.GetOptions{
				ID:             testkit.UUID,
				IncludeHistory: false,
			},
			stubs: []testkit.Stub[api.GetReplicationRequest, any]{
				{
					Request: &api.GetReplicationRequest{
						ID:             testkit.UUID,
						IncludeHistory: false,
					},
					Response: api.Replication{
						Type:       api.ReplicationCopy,
						ID:         testkit.UUID,
						Collection: "Songs",
						Shard:      "abc",
						Source:     "node-0",
						Target:     "node-1",
						StartedAt:  testkit.Now,
					},
				},
			},
			want: &replication.Operation{
				Type:       replication.Copy,
				ID:         testkit.UUID,
				Collection: "Songs",
				Shard:      "abc",
				Source:     "node-0",
				Target:     "node-1",
				StartedAt:  testkit.Now,
			},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[api.GetReplicationRequest, any]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := replication.NewClient(transport)
			require.NotNil(t, c, "nil client")

			got, err := c.Get(t.Context(), tt.options)
			tt.err.Require(t, err, "get error")
			require.EqualExportedValues(t, tt.want, got, "returned operation")
		})
	}
}

func TestClient_List(t *testing.T) {
	for _, tt := range []struct {
		name    string
		options replication.ListOptions
		stubs   []testkit.Stub[api.ListReplicationsRequest, any]
		want    []replication.Operation
		err     testkit.Error // Expected error.
	}{
		{
			name: "successfully",
			options: replication.ListOptions{
				Collection: "Songs",
				Shard:      "abc",
				Target:     "node-1",
			},
			stubs: []testkit.Stub[api.ListReplicationsRequest, any]{
				{
					Request: &api.ListReplicationsRequest{
						Collection: "Songs",
						Shard:      "abc",
						Target:     "node-1",
					},
					Response: []api.Replication{
						{
							Type:       api.ReplicationCopy,
							ID:         testkit.UUID,
							Collection: "Songs",
							Shard:      "abc",
							Source:     "node-0",
							Target:     "node-1",
							StartedAt:  testkit.Now,
							Current: api.ReplicationStage{
								State: api.ReplicationStateFinalizing,
							},
						},
					},
				},
			},
			want: []replication.Operation{
				{
					Type:       replication.Copy,
					ID:         testkit.UUID,
					Collection: "Songs",
					Shard:      "abc",
					Source:     "node-0",
					Target:     "node-1",
					StartedAt:  testkit.Now,
					Current: replication.Stage{
						State: replication.StateFinalizing,
					},
				},
			},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[api.ListReplicationsRequest, any]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := replication.NewClient(transport)
			require.NotNil(t, c, "nil client")

			got, err := c.List(t.Context(), tt.options)
			tt.err.Require(t, err, "list error")
			require.EqualExportedValues(t, tt.want, got, "returned operations")
		})
	}
}

func TestClient_Cancel(t *testing.T) {
	for _, tt := range []struct {
		name  string
		id    uuid.UUID
		stubs []testkit.Stub[any, any]
		err   testkit.Error // Expected error.
	}{
		{
			name: "successfully",
			id:   testkit.UUID,
			stubs: []testkit.Stub[any, any]{
				{Request: testkit.Ptr(api.CancelReplicationRequest(testkit.UUID))},
			},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[any, any]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := replication.NewClient(transport)
			require.NotNil(t, c, "nil client")

			err := c.Cancel(t.Context(), tt.id)
			tt.err.Require(t, err, "cancel error")
		})
	}
}

func TestClient_Delete(t *testing.T) {
	for _, tt := range []struct {
		name  string
		id    uuid.UUID
		stubs []testkit.Stub[any, any]
		err   testkit.Error // Expected error.
	}{
		{
			name: "successfully",
			id:   testkit.UUID,
			stubs: []testkit.Stub[any, any]{
				{Request: testkit.Ptr(api.DeleteReplicationRequest(testkit.UUID))},
			},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[any, any]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := replication.NewClient(transport)
			require.NotNil(t, c, "nil client")

			err := c.Delete(t.Context(), tt.id)
			tt.err.Require(t, err, "delete error")
		})
	}
}

func TestClient_DeleteAll(t *testing.T) {
	for _, tt := range []struct {
		name  string
		id    uuid.UUID
		stubs []testkit.Stub[any, any]
		err   testkit.Error // Expected error.
	}{
		{
			name: "successfully",
			id:   testkit.UUID,
			stubs: []testkit.Stub[any, any]{
				{Request: testkit.Ptr[any](api.DeleteAllReplicationsRequest)},
			},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[any, any]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := replication.NewClient(transport)
			require.NotNil(t, c, "nil client")

			err := c.DeleteAll(t.Context())
			tt.err.Require(t, err, "delete all error")
		})
	}
}
