package backup_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/backup/rbac"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/testenv"
	"github.com/weaviate/weaviate/entities/models"
)

// waitForBackupCompletion waits for a backup to complete by polling its status
func waitForBackupCompletion(t *testing.T, ctx context.Context, client *weaviate.Client, backend, backupID string) {
	t.Helper()
	
	for i := 0; i < 50; i++ { // Maximum 5 seconds (50 * 100ms)
		status, err := client.Backup().CreateStatusGetter().
			WithBackend(backend).
			WithBackupID(backupID).
			Do(ctx)
		
		if err != nil {
			// If we get a 404, the backup might not have started yet, keep waiting
			time.Sleep(100 * time.Millisecond)
			continue
		}
		
		if status != nil && status.Status != nil {
			switch *status.Status {
			case models.BackupCreateStatusResponseStatusSUCCESS:
				return // Backup completed successfully
			case models.BackupCreateStatusResponseStatusFAILED:
				t.Fatalf("Backup failed: %s", status.Error)
			case models.BackupCreateStatusResponseStatusSTARTED,
				models.BackupCreateStatusResponseStatusTRANSFERRING,
				models.BackupCreateStatusResponseStatusTRANSFERRED:
				// Still in progress, continue waiting
			}
		}
		
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatal("Backup did not complete within timeout")
}

// TestRBACBackupCreatorUsage demonstrates how to use the new RBAC features with backup creation
func TestRBACBackupCreatorUsage(t *testing.T) {
	ctx := context.Background()
	container, stop := testenv.SetupLocalContainer(t, ctx, test.RBAC, true)
	t.Cleanup(stop)

	client := testsuit.CreateTestClientForContainer(t, container)
	testsuit.CleanUpWeaviate(t, client)
	
	// Initialize random generator for unique backup IDs
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Create a test class first (required for backup to work)
	testClass := &models.Class{
		Class:       "TestBackupClass",
		Description: "Test class for backup demonstration",
		Properties: []*models.Property{
			{
				Name:     "title",
				DataType: []string{"text"},
			},
		},
	}

	err := client.Schema().ClassCreator().WithClass(testClass).Do(ctx)
	require.NoError(t, err, "failed to create test class")
	t.Cleanup(func() {
		client.Schema().ClassDeleter().WithClassName("TestBackupClass").Do(ctx)
	})

	// Example 1: Create backup with all RBAC data
	t.Run("backup with all RBAC", func(t *testing.T) {
		backupID := fmt.Sprintf("full-rbac-backup-%d", random.Int63())
		_, err := client.Backup().Creator().
			WithBackend("filesystem").
			WithBackupID(backupID).
			WithRBACAll().
			Do(ctx)
		require.NoError(t, err, "backup with all RBAC should succeed")
		
		// Wait for backup to complete
		waitForBackupCompletion(t, ctx, client, "filesystem", backupID)
	})

	// Example 2: Create backup excluding all RBAC data
	t.Run("backup excluding RBAC", func(t *testing.T) {
		backupID := fmt.Sprintf("no-rbac-backup-%d", random.Int63())
		_, err := client.Backup().Creator().
			WithBackend("filesystem").
			WithBackupID(backupID).
			WithRBACNone().
			Do(ctx)
		require.NoError(t, err, "backup excluding RBAC should succeed")
		
		// Wait for backup to complete
		waitForBackupCompletion(t, ctx, client, "filesystem", backupID)
	})

	// Example 3: Create backup with only specific roles
	t.Run("backup with specific roles", func(t *testing.T) {
		backupID := fmt.Sprintf("specific-roles-backup-%d", random.Int63())
		_, err := client.Backup().Creator().
			WithBackend("filesystem").
			WithBackupID(backupID).
			WithRBACSpecificRoles("test-manager", "test-reader", "test-writer").
			Do(ctx)
		require.NoError(t, err, "backup with specific roles should succeed")
		
		// Wait for backup to complete
		waitForBackupCompletion(t, ctx, client, "filesystem", backupID)
	})

	// Example 4: Create backup with only role definitions (no assignments)
	t.Run("backup with roles only", func(t *testing.T) {
		backupID := fmt.Sprintf("roles-only-backup-%d", random.Int63())
		_, err := client.Backup().Creator().
			WithBackend("filesystem").
			WithBackupID(backupID).
			WithRBACRolesOnly().
			Do(ctx)
		require.NoError(t, err, "backup with roles only should succeed")
		
		// Wait for backup to complete
		waitForBackupCompletion(t, ctx, client, "filesystem", backupID)
	})

	// Example 5: Create backup with custom RBAC configuration
	t.Run("backup with custom RBAC config", func(t *testing.T) {
		backupID := fmt.Sprintf("custom-rbac-backup-%d", random.Int63())
		rbacConfig := &rbac.RBACConfig{
			Scope:                   rbac.RBACAll,
			RoleSelection:          rbac.RoleSelectionSpecific,
			SpecificRoles:          []string{"test-manager", "test-reader"},
			IncludeUserAssignments: true,
			IncludeGroupAssignments: false,
		}

		_, err := client.Backup().Creator().
			WithBackend("filesystem").
			WithBackupID(backupID).
			WithRBAC(rbacConfig).
			Do(ctx)
		require.NoError(t, err, "backup with custom RBAC config should succeed")
		
		// Wait for backup to complete
		waitForBackupCompletion(t, ctx, client, "filesystem", backupID)
	})

	// Example 6: Create backup using convenience scope method
	t.Run("backup with RBAC scope", func(t *testing.T) {
		backupID := fmt.Sprintf("scope-rbac-backup-%d", random.Int63())
		_, err := client.Backup().Creator().
			WithBackend("filesystem").
			WithBackupID(backupID).
			WithRBACScope(rbac.RBACRolesOnly).
			Do(ctx)
		require.NoError(t, err, "backup with RBAC scope should succeed")
		
		// Wait for backup to complete
		waitForBackupCompletion(t, ctx, client, "filesystem", backupID)
	})

	// Example 7: Create backup using roles method
	t.Run("backup with roles method", func(t *testing.T) {
		backupID := fmt.Sprintf("roles-method-backup-%d", random.Int63())
		_, err := client.Backup().Creator().
			WithBackend("filesystem").
			WithBackupID(backupID).
			WithRoles(rbac.RoleSelectionSpecific, "test-manager", "test-writer").
			Do(ctx)
		require.NoError(t, err, "backup with roles method should succeed")
		
		// Wait for backup to complete
		waitForBackupCompletion(t, ctx, client, "filesystem", backupID)
	})
}

// TestRBACBackupRestorerUsage demonstrates how to use the new RBAC features with backup restoration
func TestRBACBackupRestorerUsage(t *testing.T) {
	ctx := context.Background()
	container, stop := testenv.SetupLocalContainer(t, ctx, test.RBAC, true)
	t.Cleanup(stop)

	client := testsuit.CreateTestClientForContainer(t, container)
	testsuit.CleanUpWeaviate(t, client)
	
	// Initialize random generator for unique backup IDs
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Create a test class first (required for backup to work)
	testClass := &models.Class{
		Class:       "TestRestoreClass",
		Description: "Test class for restore demonstration",
		Properties: []*models.Property{
			{
				Name:     "title",
				DataType: []string{"text"},
			},
		},
	}

	err := client.Schema().ClassCreator().WithClass(testClass).Do(ctx)
	require.NoError(t, err, "failed to create test class")
	t.Cleanup(func() {
		client.Schema().ClassDeleter().WithClassName("TestRestoreClass").Do(ctx)
	})

	// First create a backup to restore from
	baseBackupID := fmt.Sprintf("test-restore-backup-%d", random.Int63())
	_, err = client.Backup().Creator().
		WithBackend("filesystem").
		WithBackupID(baseBackupID).
		WithRBACAll().
		Do(ctx)
	require.NoError(t, err, "failed to create backup for restore test")
	
	// Wait for backup to complete
	waitForBackupCompletion(t, ctx, client, "filesystem", baseBackupID)

	// Example 1: Restore backup with all RBAC data
	t.Run("restore with all RBAC", func(t *testing.T) {
		_, err := client.Backup().Restorer().
			WithBackend("filesystem").
			WithBackupID(baseBackupID).
			WithRBACAll().
			Do(ctx)
		require.NoError(t, err, "restore with all RBAC should succeed")
	})

	// Example 2: Restore backup excluding all RBAC data
	t.Run("restore excluding RBAC", func(t *testing.T) {
		_, err := client.Backup().Restorer().
			WithBackend("filesystem").
			WithBackupID(baseBackupID).
			WithRBACNone().
			Do(ctx)
		require.NoError(t, err, "restore excluding RBAC should succeed")
	})

	// Example 3: Restore backup with only specific roles
	t.Run("restore with specific roles", func(t *testing.T) {
		_, err := client.Backup().Restorer().
			WithBackend("filesystem").
			WithBackupID(baseBackupID).
			WithRBACSpecificRoles("test-manager", "test-reader").
			Do(ctx)
		require.NoError(t, err, "restore with specific roles should succeed")
	})

	// Example 4: Restore backup with only user assignments
	t.Run("restore with users only", func(t *testing.T) {
		_, err := client.Backup().Restorer().
			WithBackend("filesystem").
			WithBackupID(baseBackupID).
			WithRBACUsersOnly().
			Do(ctx)
		require.NoError(t, err, "restore with users only should succeed")
	})

	// Example 5: Restore backup with custom RBAC configuration
	t.Run("restore with custom RBAC config", func(t *testing.T) {
		rbacConfig := &rbac.RBACConfig{
			Scope:                   rbac.RBACAll,
			RoleSelection:          rbac.RoleSelectionAll,
			IncludeUserAssignments: true,
			IncludeGroupAssignments: true,
		}

		_, err := client.Backup().Restorer().
			WithBackend("filesystem").
			WithBackupID(baseBackupID).
			WithRBAC(rbacConfig).
			Do(ctx)
		require.NoError(t, err, "restore with custom RBAC config should succeed")
	})
}

// TestRBACConfigHelpers demonstrates the usage of RBAC configuration helpers
func TestRBACConfigHelpers(t *testing.T) {
	t.Run("test all RBAC config", func(t *testing.T) {
		config := rbac.NewAllRBACConfig()
		if config.Scope != rbac.RBACAll {
			t.Errorf("Expected scope %v, got %v", rbac.RBACAll, config.Scope)
		}
		if !config.IncludeUserAssignments {
			t.Error("Expected IncludeUserAssignments to be true")
		}
		if !config.IncludeGroupAssignments {
			t.Error("Expected IncludeGroupAssignments to be true")
		}
	})

	t.Run("test no RBAC config", func(t *testing.T) {
		config := rbac.NewNoRBACConfig()
		if config.Scope != rbac.RBACNone {
			t.Errorf("Expected scope %v, got %v", rbac.RBACNone, config.Scope)
		}
		if config.IncludeUserAssignments {
			t.Error("Expected IncludeUserAssignments to be false")
		}
		if config.IncludeGroupAssignments {
			t.Error("Expected IncludeGroupAssignments to be false")
		}
	})

	t.Run("test specific roles config", func(t *testing.T) {
		roleNames := []string{"admin", "editor", "viewer"}
		config := rbac.NewSpecificRolesConfig(roleNames...)
		if config.RoleSelection != rbac.RoleSelectionSpecific {
			t.Errorf("Expected role selection %v, got %v", rbac.RoleSelectionSpecific, config.RoleSelection)
		}
		if len(config.SpecificRoles) != len(roleNames) {
			t.Errorf("Expected %d roles, got %d", len(roleNames), len(config.SpecificRoles))
		}
		for i, role := range roleNames {
			if config.SpecificRoles[i] != role {
				t.Errorf("Expected role %s, got %s", role, config.SpecificRoles[i])
			}
		}
	})

	t.Run("test roles only config", func(t *testing.T) {
		config := rbac.NewRolesOnlyConfig()
		if config.Scope != rbac.RBACRolesOnly {
			t.Errorf("Expected scope %v, got %v", rbac.RBACRolesOnly, config.Scope)
		}
		if config.IncludeUserAssignments {
			t.Error("Expected IncludeUserAssignments to be false")
		}
	})

	t.Run("test users only config", func(t *testing.T) {
		config := rbac.NewUsersOnlyConfig()
		if config.Scope != rbac.RBACUsersOnly {
			t.Errorf("Expected scope %v, got %v", rbac.RBACUsersOnly, config.Scope)
		}
		if !config.IncludeUserAssignments {
			t.Error("Expected IncludeUserAssignments to be true")
		}
		if config.IncludeGroupAssignments {
			t.Error("Expected IncludeGroupAssignments to be false")
		}
	})
}
