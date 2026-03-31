package backup_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/backup"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
)

func TestAwaitStatus(t *testing.T) {
	t.Run("nil backup", func(t *testing.T) {
		got, err := backup.AwaitStatus(t.Context(), nil, backup.StatusSuccess)
		assert.Nil(t, got, "must return nil *backup.Info")
		assert.Error(t, err)
	})

	t.Run("invalid backup (must have *backup.Client reference)", func(t *testing.T) {
		bak := &backup.Info{ID: "backup-1", Status: backup.StatusTransferring}

		got, err := backup.AwaitStatus(t.Context(), bak, backup.StatusSuccess)

		assert.Nil(t, got, "must return nil *backup.Info")
		assert.Error(t, err)
	})

	// The cases below describe valid backups. Each case defines the resposes
	// that AwaitStatus is expected to consume. Cases with nil/empty stubs
	// expect AwaitStatus to return without making any more requests.
	for _, tt := range []struct {
		name        string
		initStatus  backup.Status                       // Initial backup status, Started by default.
		stubs       []testkit.Stub[any, api.BackupInfo] // Responses for AwaitStatus.
		awaitStatus backup.Status                       // Passed to AwaitStatus.
		wantStatus  backup.Status                       // Latest observed status.
		ctx         context.Context                     // Using t.Context() if nil.
		err         testkit.Error                       // Set to testkit.ExpectError or testkit.ErrorIs
	}{
		{
			name:        "backup in desired state",
			initStatus:  backup.StatusTransferring,
			awaitStatus: backup.StatusTransferring,
			wantStatus:  backup.StatusTransferring,
		},
		{
			name:        "backup is already completed (status fallthrough)",
			initStatus:  backup.StatusSuccess,
			awaitStatus: backup.StatusTransferring,
			wantStatus:  backup.StatusSuccess,
			err:         testkit.ExpectError,
		},
		{
			name: "successful await",
			stubs: []testkit.Stub[any, api.BackupInfo]{
				{Response: api.BackupInfo{Status: api.BackupStatusTransferring}},
				{Response: api.BackupInfo{Status: api.BackupStatusTransferring}},
				{Response: api.BackupInfo{Status: api.BackupStatusTransferred}},
			},
			awaitStatus: backup.StatusTransferred,
			wantStatus:  backup.StatusTransferred,
		},
		{
			name: "backup is canceled abruptly (status fallthrough)",
			stubs: []testkit.Stub[any, api.BackupInfo]{
				{Response: api.BackupInfo{Status: api.BackupStatusTransferring}},
				{Response: api.BackupInfo{Status: api.BackupStatusTransferring}},
				{Response: api.BackupInfo{Status: api.BackupStatusCanceled}},
			},
			awaitStatus: backup.StatusTransferred,
			wantStatus:  backup.StatusCanceled,
			err:         testkit.ExpectError,
		},
		{
			name: "error while awaiting",
			stubs: []testkit.Stub[any, api.BackupInfo]{
				{Response: api.BackupInfo{Status: api.BackupStatusTransferring}},
				{Err: errors.New("whaam!")},
			},
			awaitStatus: backup.StatusSuccess,
			wantStatus:  backup.StatusTransferring,
			err:         testkit.ExpectError,
		},
		{
			name:        "context is canceled",
			ctx:         ctxCanceled(),
			awaitStatus: backup.StatusSuccess,
			wantStatus:  backup.StatusStarted,
			err:         testkit.ErrorIs(context.Canceled),
		},
		{
			name: "context deadline exceeded",
			stubs: []testkit.Stub[any, api.BackupInfo]{
				{Response: api.BackupInfo{Status: api.BackupStatusStarted}},
			},
			ctx:         ctxDoneOnSecondCheck(),
			awaitStatus: backup.StatusSuccess,
			wantStatus:  backup.StatusStarted,
			err:         testkit.ErrorIs(context.DeadlineExceeded),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := t.Context()
			if tt.ctx != nil {
				ctx = tt.ctx
			}

			initStatus := api.BackupStatusStarted
			if tt.initStatus != "" {
				initStatus = api.BackupStatus(tt.initStatus)
			}

			// The first response is always consumed by the test itself to GetCreateStatus.
			transport := testkit.NewTransport(t, append([]testkit.Stub[any, api.BackupInfo]{
				{Response: api.BackupInfo{Status: initStatus}},
			}, tt.stubs...))
			c := backup.NewClient(transport)
			require.NotNil(t, c, "nil client")

			// GetCreateStatus is part of test setup, always called with t.Context()
			bak, err := c.GetCreateStatus(t.Context(), backup.GetStatus{})
			require.NoError(t, err)
			require.NotNil(t, bak, "nil backup from get-status")

			// Act
			got, err := backup.AwaitStatus(
				ctx, bak, tt.awaitStatus,
				backup.WithPollingInterval(0),
			)

			// Assert
			tt.err.Require(t, err, "await status error")
			assert.NotNil(t, got, "must return latest backup status")
			assert.Equal(t, tt.wantStatus, got.Status, "latest status")
		})
	}
}

// ctxCanceled returns a canceled context.
func ctxCanceled() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return ctx
}

// ctxDoneOnSecondCheck returns a context that is Done
// after its Done method has been called twice.
//
// While this breaches the boundaries black-box testing, this
// lets us reach a case in AwaitStatus where the deadline
// expires while the goroutine is asleep without relying
// on the real clock.
func ctxDoneOnSecondCheck() context.Context {
	return testkit.NewTickingContext(2)
}

func TestInfo_IsCompleted(t *testing.T) {
	for _, tt := range []struct {
		bak  backup.Info
		want bool
	}{
		{bak: backup.Info{Status: backup.StatusStarted}, want: false},
		{bak: backup.Info{Status: backup.StatusTransferring}, want: false},
		{bak: backup.Info{Status: backup.StatusTransferred}, want: false},
		{bak: backup.Info{Status: backup.StatusSuccess}, want: true},
		{bak: backup.Info{Status: backup.StatusFailed}, want: true},
		{bak: backup.Info{Status: backup.StatusCanceled}, want: true},
	} {
		t.Run(fmt.Sprintf("status=%s", tt.bak.Status), func(t *testing.T) {
			require.Equal(t, tt.want, tt.bak.IsCompleted())
		})
	}

	t.Run("listed backups", func(t *testing.T) {
		transport := testkit.NewTransport(t, []testkit.Stub[any, []api.BackupInfo]{
			{Response: []api.BackupInfo{
				{ID: "1"}, {ID: "2"}, {ID: "3"},
			}},
		})

		c := backup.NewClient(transport)
		require.NotNil(t, c, "nil backup client")

		all, err := c.List(t.Context(), backup.List{})
		assert.NoError(t, err)
		for _, bak := range all {
			assert.True(t, bak.IsCompleted(), "bak-%s: List must return completed backups", bak.ID)
		}
	})
}
