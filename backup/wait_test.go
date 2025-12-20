package backup_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

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

	// The cases below describe valid backups. The first of the responses
	// is always consumed to fetch status outside of the await.
	// It follows that the cases with only 1 response expect AwaitStatus
	// to return without making any more requests.
	for _, tt := range []struct {
		name         string
		responses    []testkit.Response[api.BackupInfo]
		awaitStatus  backup.Status // Passed to AwaitStatus.
		expectStatus backup.Status // Latest observed status.
		errMsg       string        // Clarifying message for assert.Error.
	}{
		{
			name: "backup in desired state",
			responses: []testkit.Response[api.BackupInfo]{
				{Value: api.BackupInfo{Status: api.BackupStatusTransferred}},
			},
			awaitStatus:  backup.StatusTransferred,
			expectStatus: backup.StatusTransferred,
		},
		{
			name: "backup is already completed (status fallthrough)",
			responses: []testkit.Response[api.BackupInfo]{
				{Value: api.BackupInfo{Status: api.BackupStatusSuccess}},
			},
			awaitStatus:  backup.StatusTransferring,
			expectStatus: backup.StatusSuccess,
			errMsg:       "must not await a completed backup",
		},
		{
			name: "successful await",
			responses: []testkit.Response[api.BackupInfo]{
				{Value: api.BackupInfo{Status: api.BackupStatusStarted}},
				{Value: api.BackupInfo{Status: api.BackupStatusTransferring}},
				{Value: api.BackupInfo{Status: api.BackupStatusTransferring}},
				{Value: api.BackupInfo{Status: api.BackupStatusTransferred}},
			},
			awaitStatus:  backup.StatusTransferred,
			expectStatus: backup.StatusTransferred,
		},
		{
			name: "backup is canceled abruptly (status fallthrough)",
			responses: []testkit.Response[api.BackupInfo]{
				{Value: api.BackupInfo{Status: api.BackupStatusStarted}},
				{Value: api.BackupInfo{Status: api.BackupStatusTransferring}},
				{Value: api.BackupInfo{Status: api.BackupStatusTransferring}},
				{Value: api.BackupInfo{Status: api.BackupStatusCanceled}},
			},
			awaitStatus:  backup.StatusTransferred,
			expectStatus: backup.StatusCanceled,
			errMsg:       "must not await a completed backup",
		},
		{
			name: "error while awaiting",
			responses: []testkit.Response[api.BackupInfo]{
				{Value: api.BackupInfo{Status: api.BackupStatusStarted}},
				{Value: api.BackupInfo{Status: api.BackupStatusTransferring}},
				{Err: errors.New("whaam!")},
			},
			awaitStatus:  backup.StatusSuccess,
			expectStatus: backup.StatusTransferring,
			errMsg:       "must propagate get-status error",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewResponder(t, tt.responses)
			c := backup.NewClient(transport)
			require.NotNil(t, c, "nil client")

			bak, err := c.GetCreateStatus(t.Context(), backup.GetStatus{})
			require.NoError(t, err)
			require.NotNil(t, bak, "nil backup")

			got, err := backup.AwaitStatus(t.Context(),
				bak, tt.awaitStatus,
				backup.WithPollingInterval(0),
			)

			if tt.errMsg == "" {
				require.NoError(t, err, "await error")
			} else {
				assert.Error(t, err, tt.errMsg)
			}

			assert.NotNil(t, got, "must return latest backup status")
			assert.Equal(t, tt.expectStatus, got.Status, "latest status")
		})
	}

	t.Run("context is canceled", func(t *testing.T) {
		transport := testkit.NewResponder(t, []testkit.Response[api.BackupInfo]{
			{Value: api.BackupInfo{Status: api.BackupStatusStarted}}, // consumed before await
			{Value: api.BackupInfo{Status: api.BackupStatusStarted}}, // first status check
		})
		c := backup.NewClient(transport)
		bak, _ := c.GetCreateStatus(t.Context(), backup.GetStatus{})

		// Cancel the context right away
		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		got, err := backup.AwaitStatus(ctx,
			bak, backup.StatusSuccess,
			backup.WithPollingInterval(0),
		)

		assert.ErrorIs(t, err, context.Canceled)
		assert.NotNil(t, got, "must return latest backup status")
		assert.Equal(t, backup.StatusStarted, got.Status, "latest status")
	})

	t.Run("context timed out", func(t *testing.T) {
		transport := testkit.NewResponder(t, []testkit.Response[api.BackupInfo]{
			{Value: api.BackupInfo{Status: api.BackupStatusStarted}}, // consumed before await
			{Value: api.BackupInfo{Status: api.BackupStatusStarted}}, // first status check
		})
		c := backup.NewClient(transport)
		bak, _ := c.GetCreateStatus(t.Context(), backup.GetStatus{})

		ctx, cancel := context.WithDeadline(t.Context(), time.Now().Add(time.Nanosecond))
		defer cancel()

		got, err := backup.AwaitStatus(ctx,
			bak, backup.StatusSuccess,
			backup.WithPollingInterval(10*time.Nanosecond),
		)

		assert.ErrorIs(t, err, context.DeadlineExceeded)
		assert.NotNil(t, got, "must return latest backup status")
		assert.Equal(t, backup.StatusStarted, got.Status, "latest status")
	})
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
		transport := testkit.NewResponder(t, []testkit.Response[[]api.BackupInfo]{
			{Value: []api.BackupInfo{
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
