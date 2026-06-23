package backup_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/backup"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
)

func TestNewClient(t *testing.T) {
	require.Panics(t, func() {
		backup.NewClient(nil)
	}, "nil transport")
}

func TestClient_Create(t *testing.T) {
	for _, tt := range []struct {
		name   string
		create backup.Create
		stubs  []testkit.Stub[api.CreateBackupRequest, api.BackupInfo]
		want   *backup.Info  // Expected return value.
		err    testkit.Error // Expected error.
	}{
		{
			name: "successfully",
			create: backup.Create{
				Backend:            "filesystem",
				ID:                 "bak-1",
				Path:               "/path/to/backup",
				Endpoint:           "s3.amazonaws.com",
				Bucket:             "my-backups",
				IncludeCollections: []string{"Songs"},
				ExcludeCollections: []string{"Pizza"},
				PrefixIncremental:  "incr-bak-",
				MaxCPUPercentage:   92,
				ChunkSizeMiB:       20,
				CompressionLevel:   backup.CompressionLevelDefault,
			},
			stubs: []testkit.Stub[api.CreateBackupRequest, api.BackupInfo]{
				{
					Request: &api.CreateBackupRequest{
						Backend:            "filesystem",
						ID:                 "bak-1",
						BackupPath:         "/path/to/backup",
						Endpoint:           "s3.amazonaws.com",
						Bucket:             "my-backups",
						IncludeCollections: []string{"Songs"},
						ExcludeCollections: []string{"Pizza"},
						PrefixIncremental:  "incr-bak-",
						MaxCPUPercentage:   92,
						ChunkSizeMiB:       20,
						CompressionLevel:   api.BackupCompressionLevelDefault,
					},
					Response: api.BackupInfo{
						Backend:             "filesystem",
						ID:                  "bak-1",
						Path:                "/path/to/backup",
						Status:              api.BackupStatusStarted,
						IncludesCollections: []string{"Songs"},
					},
				},
			},
			want: &backup.Info{
				Backend:             "filesystem",
				ID:                  "bak-1",
				Path:                "/path/to/backup",
				Status:              backup.StatusStarted,
				IncludesCollections: []string{"Songs"},
			},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[api.CreateBackupRequest, api.BackupInfo]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := backup.NewClient(transport)
			require.NotNil(t, c, "nil client")

			got, err := c.Create(t.Context(), tt.create)
			tt.err.Require(t, err, "create error")
			require.EqualExportedValues(t, tt.want, got, "returned info")
		})
	}
}

func TestClient_Restore(t *testing.T) {
	for _, tt := range []struct {
		name    string
		restore backup.Restore
		stubs   []testkit.Stub[api.RestoreBackupRequest, api.BackupInfo]
		want    *backup.Info  // Expected return value.
		err     testkit.Error // Expected error.
	}{
		{
			name: "successfully",
			restore: backup.Restore{
				Backend:            "filesystem",
				ID:                 "bak-1",
				Path:               "/path/to/backup",
				Endpoint:           "s3.amazonaws.com",
				Bucket:             "my-backups",
				IncludeCollections: []string{"Songs"},
				ExcludeCollections: []string{"Pizza"},
				MaxCPUPercentage:   92,
				OverwriteAlias:     true,
				RestoreUsers:       backup.RBACRestoreAll,
				RestoreRoles:       backup.RBACRestoreNone,
				NodeMapping:        map[string]string{"node-1": "node-a"},
			},
			stubs: []testkit.Stub[api.RestoreBackupRequest, api.BackupInfo]{
				{
					Request: &api.RestoreBackupRequest{
						Backend:            "filesystem",
						ID:                 "bak-1",
						BackupPath:         "/path/to/backup",
						Endpoint:           "s3.amazonaws.com",
						Bucket:             "my-backups",
						IncludeCollections: []string{"Songs"},
						ExcludeCollections: []string{"Pizza"},
						MaxCPUPercentage:   92,
						OverwriteAlias:     true,
						RestoreUsers:       api.RBACRestoreAll,
						RestoreRoles:       api.RBACRestoreNone,
						NodeMapping:        map[string]string{"node-1": "node-a"},
					},
					Response: api.BackupInfo{
						Backend:             "filesystem",
						ID:                  "bak-1",
						Path:                "/path/to/backup",
						Status:              api.BackupStatusSuccess,
						IncludesCollections: []string{"Songs"},
						SizeGiB:             testkit.Ptr[float32](.6),
					},
				},
			},
			want: &backup.Info{
				Backend:             "filesystem",
				ID:                  "bak-1",
				Path:                "/path/to/backup",
				Status:              backup.StatusSuccess,
				IncludesCollections: []string{"Songs"},
				SizeGiB:             testkit.Ptr[float32](.6),
			},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[api.RestoreBackupRequest, api.BackupInfo]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := backup.NewClient(transport)
			require.NotNil(t, c, "nil client")

			got, err := c.Restore(t.Context(), tt.restore)
			tt.err.Require(t, err, "restore error")
			require.EqualExportedValues(t, tt.want, got, "returned info")
		})
	}
}

func TestClient_List(t *testing.T) {
	for _, tt := range []struct {
		name  string
		list  backup.List
		stubs []testkit.Stub[api.ListBackupsRequest, []api.BackupInfo]
		want  []backup.Info // Expected return value.
		err   testkit.Error // Expected error.
	}{
		{
			name: "successfully",
			list: backup.List{
				Backend:         "filesystem",
				StartingTimeAsc: true,
			},
			stubs: []testkit.Stub[api.ListBackupsRequest, []api.BackupInfo]{
				{
					Request: &api.ListBackupsRequest{
						Backend:         "filesystem",
						StartingTimeAsc: true,
					},
					Response: []api.BackupInfo{
						{ID: "bak-1"},
						{ID: "bak-2"},
						{ID: "bak-3"},
					},
				},
			},
			want: []backup.Info{
				{ID: "bak-1"},
				{ID: "bak-2"},
				{ID: "bak-3"},
			},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[api.ListBackupsRequest, []api.BackupInfo]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := backup.NewClient(transport)
			require.NotNil(t, c, "nil client")

			got, err := c.List(t.Context(), tt.list)
			tt.err.Require(t, err, "list error")
			require.EqualExportedValues(t, tt.want, got, "returned info")
		})
	}
}

// TestClient_IncrementalBaseBackupID verifies that the IncrementalBaseBackupID
// returned by the server is surfaced on backup.Info for both the create-status
// getter and the list endpoint, and that it stays empty for non-incremental
// backups. This mirrors the behavior added on main (commit 480aa7e), adapted to
// v6's mock-based client tests.
func TestClient_IncrementalBaseBackupID(t *testing.T) {
	const (
		backend      = "filesystem"
		baseBackupID = "incr-base-1"
		incrBackupID = "incr-child-1"
	)

	t.Run("create status surfaces IncrementalBaseBackupID", func(t *testing.T) {
		stubs := []testkit.Stub[api.BackupStatusRequest, api.BackupInfo]{
			{
				Request: &api.BackupStatusRequest{
					Backend:   backend,
					ID:        incrBackupID,
					Operation: api.BackupOperationCreate,
				},
				Response: api.BackupInfo{
					Backend:                 backend,
					ID:                      incrBackupID,
					Status:                  api.BackupStatusSuccess,
					IncrementalBaseBackupID: baseBackupID,
				},
			},
		}
		transport := testkit.NewTransport(t, stubs)
		c := backup.NewClient(transport)
		require.NotNil(t, c, "nil client")

		got, err := c.GetCreateStatus(t.Context(), backup.GetStatus{
			Backend: backend,
			ID:      incrBackupID,
		})
		require.NoError(t, err, "get create status")
		require.NotNil(t, got, "nil info")
		require.Equal(t, baseBackupID, got.IncrementalBaseBackupID,
			"IncrementalBaseBackupID should match the base backup ID")
	})

	t.Run("create status leaves IncrementalBaseBackupID empty for base backup", func(t *testing.T) {
		stubs := []testkit.Stub[api.BackupStatusRequest, api.BackupInfo]{
			{
				Request: &api.BackupStatusRequest{
					Backend:   backend,
					ID:        baseBackupID,
					Operation: api.BackupOperationCreate,
				},
				Response: api.BackupInfo{
					Backend: backend,
					ID:      baseBackupID,
					Status:  api.BackupStatusSuccess,
				},
			},
		}
		transport := testkit.NewTransport(t, stubs)
		c := backup.NewClient(transport)
		require.NotNil(t, c, "nil client")

		got, err := c.GetCreateStatus(t.Context(), backup.GetStatus{
			Backend: backend,
			ID:      baseBackupID,
		})
		require.NoError(t, err, "get create status")
		require.NotNil(t, got, "nil info")
		require.Empty(t, got.IncrementalBaseBackupID,
			"base backup should not have IncrementalBaseBackupID")
	})

	t.Run("list surfaces IncrementalBaseBackupID", func(t *testing.T) {
		stubs := []testkit.Stub[api.ListBackupsRequest, []api.BackupInfo]{
			{
				Request: &api.ListBackupsRequest{Backend: backend},
				Response: []api.BackupInfo{
					{ID: baseBackupID, Status: api.BackupStatusSuccess},
					{
						ID:                      incrBackupID,
						Status:                  api.BackupStatusSuccess,
						IncrementalBaseBackupID: baseBackupID,
					},
				},
			},
		}
		transport := testkit.NewTransport(t, stubs)
		c := backup.NewClient(transport)
		require.NotNil(t, c, "nil client")

		got, err := c.List(t.Context(), backup.List{Backend: backend})
		require.NoError(t, err, "list backups")
		require.Len(t, got, 2, "returned backups")

		byID := make(map[string]backup.Info, len(got))
		for _, b := range got {
			byID[b.ID] = b
		}

		require.Equal(t, baseBackupID, byID[incrBackupID].IncrementalBaseBackupID,
			"incremental backup should reference the base backup ID")
		require.Empty(t, byID[baseBackupID].IncrementalBaseBackupID,
			"base backup should not have IncrementalBaseBackupID")
	})
}

func TestClient_Cancel(t *testing.T) {
	for _, tt := range []struct {
		name   string
		cancel func(context.Context, *backup.Client) error // Call appropriate cancel function.
		stubs  []testkit.Stub[api.CancelBackupRequest, any]
		err    testkit.Error // Expected error.
	}{
		{
			name: "cancel create successfully",
			cancel: func(ctx context.Context, c *backup.Client) error {
				return c.CancelCreate(ctx, backup.Cancel{
					Backend: "filesystem",
					ID:      "bak-1",
				})
			},
			stubs: []testkit.Stub[api.CancelBackupRequest, any]{
				{
					Request: &api.CancelBackupRequest{
						Backend:   "filesystem",
						ID:        "bak-1",
						Operation: api.BackupOperationCreate,
					},
				},
			},
		},
		{
			name: "cancel create with error",
			cancel: func(ctx context.Context, c *backup.Client) error {
				return c.CancelCreate(ctx, backup.Cancel{})
			},
			stubs: []testkit.Stub[api.CancelBackupRequest, any]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
		{
			name: "cancel restore successfully",
			cancel: func(ctx context.Context, c *backup.Client) error {
				return c.CancelRestore(ctx, backup.Cancel{
					Backend: "filesystem",
					ID:      "bak-1",
				})
			},
			stubs: []testkit.Stub[api.CancelBackupRequest, any]{
				{
					Request: &api.CancelBackupRequest{
						Backend:   "filesystem",
						ID:        "bak-1",
						Operation: api.BackupOperationRestore,
					},
				},
			},
		},
		{
			name: "cancel restore with error",
			cancel: func(ctx context.Context, c *backup.Client) error {
				return c.CancelRestore(ctx, backup.Cancel{})
			},
			stubs: []testkit.Stub[api.CancelBackupRequest, any]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := backup.NewClient(transport)
			require.NotNil(t, c, "nil client")

			err := tt.cancel(t.Context(), c)
			tt.err.Require(t, err, "cancel error")
		})
	}
}
