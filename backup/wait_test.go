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

	t.Run("backup in desired state", func(t *testing.T) {
		transport := testkit.NewResponder(t, []testkit.Response[api.BackupInfo]{
			{Value: api.BackupInfo{Status: api.BackupStatusTransferred}},
		})
		c := backup.NewClient(transport)

		bak, err := c.GetCreateStatus(t.Context(), backup.GetStatus{})
		require.NoError(t, err)

		got, err := backup.AwaitStatus(t.Context(), bak, backup.StatusTransferred)

		assert.Equal(t, *bak, *got)
		assert.NoError(t, err, "await error")
	})

	t.Run("backup is already completed (status fallthrough)", func(t *testing.T) {
		transport := testkit.NewResponder(t, []testkit.Response[api.BackupInfo]{
			{Value: api.BackupInfo{Status: api.BackupStatusSuccess}},
		})
		c := backup.NewClient(transport)

		bak, err := c.GetCreateStatus(t.Context(), backup.GetStatus{})
		require.NoError(t, err)
		require.NotNil(t, bak, "nil backup")

		got, err := backup.AwaitStatus(t.Context(), bak, backup.StatusCanceled)

		assert.NotNil(t, got, "must return completed backup")
		assert.Equal(t, *bak, *got)
		assert.Error(t, err, "cannot await a completed backup")
	})

	t.Run("successful await", func(t *testing.T) {
		transport := testkit.NewResponder(t, []testkit.Response[api.BackupInfo]{
			{Value: api.BackupInfo{Status: api.BackupStatusStarted}},
			{Value: api.BackupInfo{Status: api.BackupStatusTransferring}},
			{Value: api.BackupInfo{Status: api.BackupStatusTransferring}},
			{Value: api.BackupInfo{Status: api.BackupStatusTransferred}},
		})
		c := backup.NewClient(transport)

		bak, err := c.GetCreateStatus(t.Context(), backup.GetStatus{})
		require.NoError(t, err)

		got, err := backup.AwaitStatus(t.Context(),
			bak, backup.StatusTransferred,
			backup.WithPollingInterval(0),
		)

		require.NoError(t, err, "await error")
		assert.NotNil(t, got, "must return latest backup status")
		assert.Equal(t, backup.StatusTransferred, got.Status, "latest status")
	})

	t.Run("backup is canceled abruptly (status fallthrough)", func(t *testing.T) {
		transport := testkit.NewResponder(t, []testkit.Response[api.BackupInfo]{
			{Value: api.BackupInfo{Status: api.BackupStatusStarted}},
			{Value: api.BackupInfo{Status: api.BackupStatusTransferring}},
			{Value: api.BackupInfo{Status: api.BackupStatusTransferring}},
			{Value: api.BackupInfo{Status: api.BackupStatusCanceled}},
		})
		c := backup.NewClient(transport)

		bak, err := c.GetCreateStatus(t.Context(), backup.GetStatus{})
		require.NoError(t, err)

		got, err := backup.AwaitStatus(t.Context(),
			bak, backup.StatusTransferred,
			backup.WithPollingInterval(0),
		)

		assert.Error(t, err, "cannot await a completed backup")
		assert.NotNil(t, got, "must return latest backup status")
		assert.Equal(t, backup.StatusCanceled, got.Status, "latest status")
	})

	t.Run("error while awaiting", func(t *testing.T) {
		transport := testkit.NewResponder(t, []testkit.Response[api.BackupInfo]{
			{Value: api.BackupInfo{Status: api.BackupStatusStarted}},
			{Value: api.BackupInfo{Status: api.BackupStatusTransferring}},
			{Err: errors.New("whaam!")},
		})
		c := backup.NewClient(transport)

		bak, err := c.GetCreateStatus(t.Context(), backup.GetStatus{})
		require.NoError(t, err)

		got, err := backup.AwaitStatus(t.Context(),
			bak, backup.StatusSuccess,
			backup.WithPollingInterval(0),
		)

		assert.Error(t, err, "must propagate get-status error")
		assert.NotNil(t, got, "must return latest backup status")
		assert.Equal(t, backup.StatusTransferring, got.Status, "latest status")
	})

	t.Run("context is canceled", func(t *testing.T) {
		transport := testkit.NewResponder(t, []testkit.Response[api.BackupInfo]{
			{Value: api.BackupInfo{Status: api.BackupStatusStarted}}, // consumed before await
			{Value: api.BackupInfo{Status: api.BackupStatusStarted}}, // first status check
		})
		c := backup.NewClient(transport)

		bak, err := c.GetCreateStatus(t.Context(), backup.GetStatus{})
		require.NoError(t, err)

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
	t.Run("context times out", func(t *testing.T) {
		transport := testkit.NewResponder(t, []testkit.Response[api.BackupInfo]{
			{Value: api.BackupInfo{Status: api.BackupStatusStarted}}, // consumed before await
			{Value: api.BackupInfo{Status: api.BackupStatusStarted}}, // first status check
		})
		c := backup.NewClient(transport)

		bak, err := c.GetCreateStatus(t.Context(), backup.GetStatus{})
		require.NoError(t, err)

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
