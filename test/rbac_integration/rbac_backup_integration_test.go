package rbac_integration

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/backup"
	backupRbac "github.com/weaviate/weaviate-go-client/v5/weaviate/backup/rbac"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/rbac"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/testenv"
	"github.com/weaviate/weaviate/entities/models"
)

// TestRBACBackupIntegration tests the complete RBAC backup and restore functionality
func TestRBACBackupIntegration(t *testing.T) {
	ctx := context.Background()
	container, stop := testenv.SetupLocalContainer(t, ctx, test.RBAC, true)
	t.Cleanup(stop)

	client := testsuit.CreateTestClientForContainer(t, container)
	testsuit.CleanUpWeaviate(t, client)

	backend := backup.BACKEND_FILESYSTEM

	const (
		testUserID    = "testuser"
		testUserKey   = "testkey"
		testClassName = "Article"
		testClassDesc = "News articles"
		backupIDBase  = "rbactest"
	)

	// Test mirrors the CLI script operations
	t.Run("RBAC Backup and Restore Integration Test", func(t *testing.T) {
		// Step 1: Create test user (equivalent to CLI user create)
		t.Run("create test user", func(t *testing.T) {
			apiKey, err := client.Users().DB().Creator().WithUserID(testUserID).Do(ctx)
			require.NoError(t, err, "failed to create test user")
			require.NotEmpty(t, apiKey, "API key should not be empty")
			t.Logf("Created user '%s' with API key", testUserID)

			// Cleanup user at the end of the test
			t.Cleanup(func() {
				deleted, err := client.Users().DB().Deleter().WithUserID(testUserID).Do(ctx)
				if err != nil {
					t.Logf("Error cleaning up user %s: %v", testUserID, err)
				} else if deleted {
					t.Logf("Successfully cleaned up user %s", testUserID)
				}
			})
		})

		// Step 2: Create test class (equivalent to CLI class create)
		t.Run("create test class", func(t *testing.T) {
			class := &models.Class{
				Class:       testClassName,
				Description: testClassDesc,
				Properties: []*models.Property{
					{
						Name:     "title",
						DataType: []string{"text"},
					},
					{
						Name:     "content",
						DataType: []string{"text"},
					},
				},
			}

			err := client.Schema().ClassCreator().WithClass(class).Do(ctx)
			require.NoError(t, err, "failed to create test class")
			t.Logf("Created class '%s'", testClassName)

			// Verify class exists in schema
			t.Run("verify class exists in schema", func(t *testing.T) {
				schema, err := client.Schema().Getter().Do(ctx)
				require.NoError(t, err, "failed to get schema")

				found := false
				for _, class := range schema.Classes {
					if class.Class == testClassName {
						found = true
						break
					}
				}
				require.True(t, found, "Class '%s' should exist in schema", testClassName)
				t.Logf("Class '%s' verified in schema", testClassName)
			})

			// Cleanup class at the end of the test
			t.Cleanup(func() {
				err := client.Schema().ClassDeleter().WithClassName(testClassName).Do(ctx)
				if err != nil {
					t.Logf("Error cleaning up class %s: %v", testClassName, err)
				} else {
					t.Logf("Successfully cleaned up class %s", testClassName)
				}
			})
		})

		// Step 3: Create multiple backups with RBAC data (equivalent to CLI backup create loop)
		var backupIDs []string
		for i := 1; i <= 10; i++ {
			backupID := fmt.Sprintf("%s%d", backupIDBase, i)
			backupIDs = append(backupIDs, backupID)

			t.Run(fmt.Sprintf("create backup %s", backupID), func(t *testing.T) {
				// Create backup with all RBAC data
				createResponse, err := client.Backup().Creator().
					WithBackend(backend).
					WithBackupID(backupID).
					WithRBACAll(). // Include all RBAC data
					WithWaitForCompletion(true).
					Do(ctx)

				require.NoError(t, err, "failed to create backup %s", backupID)
				require.NotNil(t, createResponse, "create response should not be nil")
				assert.Equal(t, backupID, createResponse.ID)
				assert.Equal(t, models.BackupCreateResponseStatusSUCCESS, *createResponse.Status)
				t.Logf("Successfully created backup %s", backupID)

				// Verify backup status
				t.Run("verify backup status", func(t *testing.T) {
					statusResponse, err := client.Backup().CreateStatusGetter().
						WithBackend(backend).
						WithBackupID(backupID).
						Do(ctx)

					require.NoError(t, err, "failed to get backup status")
					require.NotNil(t, statusResponse, "status response should not be nil")
					assert.Equal(t, models.BackupCreateStatusResponseStatusSUCCESS, *statusResponse.Status)
					t.Logf("Backup %s status: %s", backupID, *statusResponse.Status)
				})
			})
		}

		// Step 4: Delete the test user (equivalent to CLI user delete)
		t.Run("delete test user", func(t *testing.T) {
			deleted, err := client.Users().DB().Deleter().WithUserID(testUserID).Do(ctx)
			require.NoError(t, err, "failed to delete test user")
			require.True(t, deleted, "user should be successfully deleted")
			t.Logf("Successfully deleted user '%s'", testUserID)

			// Verify user deletion
			t.Run("verify user deletion", func(t *testing.T) {
				users, err := client.Users().DB().Lister().Do(ctx)
				require.NoError(t, err, "failed to list users")

				found := false
				for _, user := range users {
					if user.UserID == testUserID {
						found = true
						break
					}
				}
				assert.False(t, found, "User '%s' should not exist after deletion", testUserID)
				t.Logf("User '%s' deletion verified", testUserID)
			})
		})

		// Step 5: Delete the test class (equivalent to CLI class delete)
		t.Run("delete test class", func(t *testing.T) {
			err := client.Schema().ClassDeleter().WithClassName(testClassName).Do(ctx)
			require.NoError(t, err, "failed to delete test class")
			t.Logf("Successfully deleted class '%s'", testClassName)

			// Verify class deletion
			t.Run("verify class deletion", func(t *testing.T) {
				schema, err := client.Schema().Getter().Do(ctx)
				require.NoError(t, err, "failed to get schema")

				found := false
				for _, class := range schema.Classes {
					if class.Class == testClassName {
						found = true
						break
					}
				}
				assert.False(t, found, "Class '%s' should not exist after deletion", testClassName)
				t.Logf("Class '%s' deletion verified", testClassName)
			})
		})

		// Step 6: Restore the first backup with all RBAC data (equivalent to CLI backup restore)
		restoreBackupID := backupIDs[0]
		t.Run(fmt.Sprintf("restore backup %s", restoreBackupID), func(t *testing.T) {
			// Restore backup with all RBAC data
			restoreResponse, err := client.Backup().Restorer().
				WithBackend(backend).
				WithBackupID(restoreBackupID).
				WithRBACAll(). // Restore all RBAC data
				WithWaitForCompletion(true).
				Do(ctx)

			require.NoError(t, err, "failed to restore backup %s", restoreBackupID)
			require.NotNil(t, restoreResponse, "restore response should not be nil")
			assert.Equal(t, restoreBackupID, restoreResponse.ID)
			assert.Equal(t, models.BackupRestoreResponseStatusSUCCESS, *restoreResponse.Status)
			t.Logf("Successfully restored backup %s", restoreBackupID)

			// Verify restore status
			t.Run("verify restore status", func(t *testing.T) {
				statusResponse, err := client.Backup().RestoreStatusGetter().
					WithBackend(backend).
					WithBackupID(restoreBackupID).
					Do(ctx)

				require.NoError(t, err, "failed to get restore status")
				require.NotNil(t, statusResponse, "status response should not be nil")
				assert.Equal(t, models.BackupRestoreStatusResponseStatusSUCCESS, *statusResponse.Status)
				t.Logf("Restore %s status: %s", restoreBackupID, *statusResponse.Status)
			})
		})

		// Step 7: Verify that class exists after restore
		t.Run("verify class exists after restore", func(t *testing.T) {
			schema, err := client.Schema().Getter().Do(ctx)
			require.NoError(t, err, "failed to get schema after restore")

			found := false
			for _, class := range schema.Classes {
				if class.Class == testClassName {
					found = true
					assert.Equal(t, testClassDesc, class.Description)
					break
				}
			}
			require.True(t, found, "Class '%s' should exist after restore", testClassName)
			t.Logf("Class '%s' exists after restore", testClassName)
		})

		// Step 8: Verify that user exists after restore
		t.Run("verify user exists after restore", func(t *testing.T) {
			users, err := client.Users().DB().Lister().Do(ctx)
			require.NoError(t, err, "failed to list users after restore")

			found := false
			for _, user := range users {
				if user.UserID == testUserID {
					found = true
					break
				}
			}
			require.True(t, found, "User '%s' should exist after restore", testUserID)
			t.Logf("User '%s' exists after restore", testUserID)
		})
	})
}

// TestRBACBackupDifferentConfigurations tests various RBAC backup configurations
func TestRBACBackupDifferentConfigurations(t *testing.T) {
	ctx := context.Background()
	container, stop := testenv.SetupLocalContainer(t, ctx, test.RBAC, true)
	t.Cleanup(stop)

	client := testsuit.CreateTestClientForContainer(t, container)
	testsuit.CleanUpWeaviate(t, client)

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	backend := backup.BACKEND_FILESYSTEM

	const (
		testRole = "test-role"
		testUser = "test-backup-user"
	)

	// Setup: Create a test class and role and user
	// Create test class first (required for backup to work) - at main test level
	testClass := &models.Class{
		Class:       "BackupTestClass",
		Description: "Test class for backup configuration testing",
		Properties: []*models.Property{
			{
				Name:     "title",
				DataType: []string{"text"},
			},
		},
	}

	err := client.Schema().ClassCreator().WithClass(testClass).Do(ctx)
	require.NoError(t, err, "failed to create test class")
	t.Logf("Created test class '%s'", "BackupTestClass")

	t.Cleanup(func() {
		client.Schema().ClassDeleter().WithClassName("BackupTestClass").Do(ctx)
	})

	t.Run("setup test role and user", func(t *testing.T) {
		// Create test role
		role := rbac.NewRole(testRole,
			rbac.BackupsPermission{
				Actions:    []string{models.PermissionActionManageBackups},
				Collection: "*",
			},
		)

		err = client.Roles().Creator().WithRole(role).Do(ctx)
		require.NoError(t, err, "failed to create test role")
		t.Logf("Created test role '%s'", testRole)

		// Create test user
		apiKey, err := client.Users().DB().Creator().WithUserID(testUser).Do(ctx)
		require.NoError(t, err, "failed to create test user")
		require.NotEmpty(t, apiKey, "API key should not be empty")
		t.Logf("Created test user '%s'", testUser)

		// Assign role to user
		err = client.Users().DB().RolesAssigner().WithUserID(testUser).WithRoles(testRole).Do(ctx)
		require.NoError(t, err, "failed to assign role to user")
		t.Logf("Assigned role '%s' to user '%s'", testRole, testUser)

		// Cleanup
		t.Cleanup(func() {
			client.Users().DB().Deleter().WithUserID(testUser).Do(ctx)
			client.Roles().Deleter().WithName(testRole).Do(ctx)
		})
	})

	// Test different RBAC backup configurations
	tests := []struct {
		name         string
		setupBackup  func(*backup.BackupCreator) *backup.BackupCreator
		setupRestore func(*backup.BackupRestorer) *backup.BackupRestorer
		description  string
	}{
		{
			name: "backup_with_all_rbac",
			setupBackup: func(creator *backup.BackupCreator) *backup.BackupCreator {
				return creator.WithRBACAll()
			},
			setupRestore: func(restorer *backup.BackupRestorer) *backup.BackupRestorer {
				return restorer.WithRBACAll()
			},
			description: "backup and restore with all RBAC data",
		},
		{
			name: "backup_with_no_rbac",
			setupBackup: func(creator *backup.BackupCreator) *backup.BackupCreator {
				return creator.WithRBAC(backupRbac.NewRBACConfigFromFlags(backupRbac.None))
			},
			setupRestore: func(restorer *backup.BackupRestorer) *backup.BackupRestorer {
				return restorer.WithRBACNone()
			},
			description: "backup and restore with no RBAC data",
		},
		{
			name: "backup_with_specific_roles",
			setupBackup: func(creator *backup.BackupCreator) *backup.BackupCreator {
				return creator.WithRBACSpecificRoles(testRole)
			},
			setupRestore: func(restorer *backup.BackupRestorer) *backup.BackupRestorer {
				return restorer.WithRBACSpecificRoles(testRole)
			},
			description: "backup and restore with specific roles only",
		},
		{
			name: "backup_with_roles_only",
			setupBackup: func(creator *backup.BackupCreator) *backup.BackupCreator {
				return creator.WithRBACRolesOnly()
			},
			setupRestore: func(restorer *backup.BackupRestorer) *backup.BackupRestorer {
				return restorer.WithRBACRolesOnly()
			},
			description: "backup and restore with role definitions only",
		},
		{
			name: "backup_with_users_only",
			setupBackup: func(creator *backup.BackupCreator) *backup.BackupCreator {
				return creator.WithRBACUsersOnly()
			},
			setupRestore: func(restorer *backup.BackupRestorer) *backup.BackupRestorer {
				return restorer.WithRBACUsersOnly()
			},
			description: "backup and restore with user assignments only",
		},
		{
			name: "backup_with_custom_rbac_config",
			setupBackup: func(creator *backup.BackupCreator) *backup.BackupCreator {
				rbacConfig := &backupRbac.RBACConfig{
					Scope:                    backupRbac.RBACAll,
					RoleSelection:           backupRbac.RoleSelectionSpecific,
					SpecificRoles:           []string{testRole},
					IncludeUserAssignments:  true,
					IncludeGroupAssignments: false,
				}
				return creator.WithRBAC(rbacConfig)
			},
			setupRestore: func(restorer *backup.BackupRestorer) *backup.BackupRestorer {
				rbacConfig := &backupRbac.RBACConfig{
					Scope:                    backupRbac.RBACAll,
					RoleSelection:           backupRbac.RoleSelectionSpecific,
					SpecificRoles:           []string{testRole},
					IncludeUserAssignments:  true,
					IncludeGroupAssignments: false,
				}
				return restorer.WithRBAC(rbacConfig)
			},
			description: "backup and restore with custom RBAC configuration",
		},
		{
			name: "backup_with_bitwise_flags",
			setupBackup: func(creator *backup.BackupCreator) *backup.BackupCreator {
				return creator.WithRBACFlags(backupRbac.Roles | backupRbac.Users)
			},
			setupRestore: func(restorer *backup.BackupRestorer) *backup.BackupRestorer {
				return restorer.WithRBACFlags(backupRbac.Roles | backupRbac.Users)
			},
			description: "backup and restore using bitwise flags (Roles | Users)",
		},
		{
			name: "backup_with_none_flag",
			setupBackup: func(creator *backup.BackupCreator) *backup.BackupCreator {
				return creator.WithRBACFlags(backupRbac.None)
			},
			setupRestore: func(restorer *backup.BackupRestorer) *backup.BackupRestorer {
				return restorer.WithRBACFlags(backupRbac.None)
			},
			description: "backup and restore with no RBAC data using None flag",
		},
		{
			name: "backup_with_flags_and_specific_roles",
			setupBackup: func(creator *backup.BackupCreator) *backup.BackupCreator {
				return creator.WithRBACFlags(backupRbac.Roles | backupRbac.Users).WithSpecificRoles(testRole)
			},
			setupRestore: func(restorer *backup.BackupRestorer) *backup.BackupRestorer {
				return restorer.WithRBACFlags(backupRbac.Roles | backupRbac.Users).WithSpecificRoles(testRole)
			},
			description: "backup and restore using flags with specific roles",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backupID := fmt.Sprintf("%s-%d", tt.name, random.Int63())
			t.Logf("Testing %s with backup ID: %s", tt.description, backupID)

			// Create backup with specific RBAC configuration
			t.Run("create backup", func(t *testing.T) {
				creator := client.Backup().Creator().
					WithBackend(backend).
					WithBackupID(backupID).
					WithWaitForCompletion(true)

				creator = tt.setupBackup(creator)

				createResponse, err := creator.Do(ctx)
				require.NoError(t, err, "failed to create backup with %s", tt.description)
				require.NotNil(t, createResponse, "create response should not be nil")
				assert.Equal(t, backupID, createResponse.ID)
				assert.Equal(t, models.BackupCreateResponseStatusSUCCESS, *createResponse.Status)
				t.Logf("Successfully created backup %s", backupID)
			})

			// Test restore with the same RBAC configuration
			t.Run("restore backup", func(t *testing.T) {
				restorer := client.Backup().Restorer().
					WithBackend(backend).
					WithBackupID(backupID).
					WithWaitForCompletion(true)

				restorer = tt.setupRestore(restorer)

				restoreResponse, err := restorer.Do(ctx)
				require.NoError(t, err, "failed to restore backup with %s", tt.description)
				require.NotNil(t, restoreResponse, "restore response should not be nil")
				assert.Equal(t, backupID, restoreResponse.ID)
				assert.Equal(t, models.BackupRestoreResponseStatusSUCCESS, *restoreResponse.Status)
				t.Logf("Successfully restored backup %s", backupID)
			})

			// Verify that role and user still exist (depending on the RBAC config)
			if tt.name != "backup_with_no_rbac" && tt.name != "backup_with_none_flag" {
				t.Run("verify rbac data exists after restore", func(t *testing.T) {
					// Check if role exists
					exists, err := client.Roles().Exists().WithName(testRole).Do(ctx)
					if err == nil {
						assert.True(t, exists, "Role '%s' should exist after restore", testRole)
						t.Logf("Role '%s' exists after restore", testRole)
					} else {
						t.Logf("Could not check role existence: %v", err)
					}

					// Check if user exists (depending on configuration)
					if tt.name == "backup_with_all_rbac" || tt.name == "backup_with_users_only" || tt.name == "backup_with_custom_rbac_config" || tt.name == "backup_with_bitwise_flags" || tt.name == "backup_with_flags_and_specific_roles" {
						users, err := client.Users().DB().Lister().Do(ctx)
						if err == nil {
							found := false
							for _, user := range users {
								if user.UserID == testUser {
									found = true
									break
								}
							}
							assert.True(t, found, "User '%s' should exist after restore", testUser)
							t.Logf("User '%s' exists after restore", testUser)
						} else {
							t.Logf("Could not check user existence: %v", err)
						}
					}
				})
			}
		})
	}
}

// TestRBACBackupErrorHandling tests error conditions and edge cases
func TestRBACBackupErrorHandling(t *testing.T) {
	ctx := context.Background()
	container, stop := testenv.SetupLocalContainer(t, ctx, test.RBAC, true)
	t.Cleanup(stop)

	client := testsuit.CreateTestClientForContainer(t, container)
	testsuit.CleanUpWeaviate(t, client)

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	backend := backup.BACKEND_FILESYSTEM

	// Create a test class for error handling tests - at main test level
	testClass := &models.Class{
		Class:       "ErrorHandlingTestClass",
		Description: "Test class for error handling tests",
		Properties: []*models.Property{
			{
				Name:     "name",
				DataType: []string{"text"},
			},
		},
	}

	err := client.Schema().ClassCreator().WithClass(testClass).Do(ctx)
	require.NoError(t, err, "failed to create test class")
	t.Cleanup(func() {
		client.Schema().ClassDeleter().WithClassName("ErrorHandlingTestClass").Do(ctx)
	})

	t.Run("test backup with non-existent role", func(t *testing.T) {
		backupID := fmt.Sprintf("non-existent-role-%d", random.Int63())
		nonExistentRole := "non-existent-role"

		// This should succeed even if the role doesn't exist
		createResponse, err := client.Backup().Creator().
			WithBackend(backend).
			WithBackupID(backupID).
			WithRBACSpecificRoles(nonExistentRole).
			WithWaitForCompletion(true).
			Do(ctx)

		// The backup should succeed but may contain no RBAC data for the non-existent role
		require.NoError(t, err, "backup should succeed even with non-existent role")
		require.NotNil(t, createResponse)
		assert.Equal(t, models.BackupCreateResponseStatusSUCCESS, *createResponse.Status)
		t.Logf("Backup with non-existent role completed successfully")
	})

	t.Run("test restore with conflicting RBAC configuration", func(t *testing.T) {
		backupID := fmt.Sprintf("conflicting-rbac-%d", random.Int63())
		testRole := "conflicting-test-role"

		// Create a test role first
		role := rbac.NewRole(testRole,
			rbac.BackupsPermission{
				Actions:    []string{models.PermissionActionManageBackups},
				Collection: "*",
			},
		)

		err := client.Roles().Creator().WithRole(role).Do(ctx)
		require.NoError(t, err, "failed to create test role")
		t.Cleanup(func() {
			client.Roles().Deleter().WithName(testRole).Do(ctx)
		})

		// Create backup with all RBAC data
		createResponse, err := client.Backup().Creator().
			WithBackend(backend).
			WithBackupID(backupID).
			WithRBACAll().
			WithWaitForCompletion(true).
			Do(ctx)

		require.NoError(t, err, "failed to create backup")
		require.NotNil(t, createResponse)
		assert.Equal(t, models.BackupCreateResponseStatusSUCCESS, *createResponse.Status)

		// Try to restore with a different RBAC configuration
		// This should still succeed
		restoreResponse, err := client.Backup().Restorer().
			WithBackend(backend).
			WithBackupID(backupID).
			WithRBACNone(). // Different from backup configuration
			WithWaitForCompletion(true).
			Do(ctx)

		require.NoError(t, err, "restore should succeed even with different RBAC config")
		require.NotNil(t, restoreResponse)
		assert.Equal(t, models.BackupRestoreResponseStatusSUCCESS, *restoreResponse.Status)
		t.Logf("Restore with conflicting RBAC configuration completed successfully")
	})

	t.Run("test backup with invalid backend", func(t *testing.T) {
		backupID := fmt.Sprintf("invalid-backend-%d", random.Int63())
		invalidBackend := "invalid-backend"

		// This should fail
		createResponse, err := client.Backup().Creator().
			WithBackend(invalidBackend).
			WithBackupID(backupID).
			WithRBACAll().
			Do(ctx)

		require.Error(t, err, "backup should fail with invalid backend")
		require.Nil(t, createResponse)
		assert.Contains(t, err.Error(), "422", "should return 422 error")
		t.Logf("Backup with invalid backend failed as expected: %v", err)
	})
}

// TestRBACBackupValidationAndWaiting tests backup creation and validation patterns
func TestRBACBackupValidationAndWaiting(t *testing.T) {
	ctx := context.Background()
	container, stop := testenv.SetupLocalContainer(t, ctx, test.RBAC, true)
	t.Cleanup(stop)

	client := testsuit.CreateTestClientForContainer(t, container)
	testsuit.CleanUpWeaviate(t, client)

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	backend := backup.BACKEND_FILESYSTEM

	// Create a test class for validation tests - at main test level
	testClass := &models.Class{
		Class:       "ValidationTestClass",
		Description: "Test class for validation tests",
		Properties: []*models.Property{
			{
				Name:     "name",
				DataType: []string{"text"},
			},
		},
	}

	err := client.Schema().ClassCreator().WithClass(testClass).Do(ctx)
	require.NoError(t, err, "failed to create test class")
	t.Cleanup(func() {
		client.Schema().ClassDeleter().WithClassName("ValidationTestClass").Do(ctx)
	})

	t.Run("test backup creation without waiting", func(t *testing.T) {
		backupID := fmt.Sprintf("no-wait-%d", random.Int63())

		// Create backup without waiting for completion
		createResponse, err := client.Backup().Creator().
			WithBackend(backend).
			WithBackupID(backupID).
			WithRBACAll().
			// Don't call WithWaitForCompletion(true)
			Do(ctx)

		require.NoError(t, err, "failed to create backup")
		require.NotNil(t, createResponse)
		assert.Equal(t, backupID, createResponse.ID)
		// Should be STARTED, not SUCCESS
		assert.Equal(t, models.BackupCreateResponseStatusSTARTED, *createResponse.Status)
		t.Logf("Backup %s started with status: %s", backupID, *createResponse.Status)

		// Wait for completion by polling status
		t.Run("wait for backup completion by polling", func(t *testing.T) {
			statusGetter := client.Backup().CreateStatusGetter().
				WithBackend(backend).
				WithBackupID(backupID)

			timeout := time.After(30 * time.Second)
			ticker := time.NewTicker(500 * time.Millisecond)
			defer ticker.Stop()

			for {
				select {
				case <-timeout:
					t.Fatal("Timeout waiting for backup completion")
				case <-ticker.C:
					statusResponse, err := statusGetter.Do(ctx)
					require.NoError(t, err, "failed to get backup status")
					require.NotNil(t, statusResponse)

					t.Logf("Backup %s status: %s", backupID, *statusResponse.Status)

					if *statusResponse.Status == models.BackupCreateStatusResponseStatusSUCCESS {
						t.Logf("Backup %s completed successfully", backupID)
						return
					}

					if *statusResponse.Status == models.BackupCreateStatusResponseStatusFAILED {
						t.Fatalf("Backup %s failed: %s", backupID, statusResponse.Error)
					}

					// Continue waiting for other statuses (STARTED, TRANSFERRING, etc.)
				}
			}
		})
	})

	t.Run("test restore without waiting", func(t *testing.T) {
		backupID := fmt.Sprintf("restore-no-wait-%d", random.Int63())

		// First create a backup
		createResponse, err := client.Backup().Creator().
			WithBackend(backend).
			WithBackupID(backupID).
			WithRBACAll().
			WithWaitForCompletion(true).
			Do(ctx)

		require.NoError(t, err, "failed to create backup for restore test")
		require.NotNil(t, createResponse)
		assert.Equal(t, models.BackupCreateResponseStatusSUCCESS, *createResponse.Status)

		// Restore without waiting
		restoreResponse, err := client.Backup().Restorer().
			WithBackend(backend).
			WithBackupID(backupID).
			WithRBACAll().
			// Don't call WithWaitForCompletion(true)
			Do(ctx)

		require.NoError(t, err, "failed to start restore")
		require.NotNil(t, restoreResponse)
		assert.Equal(t, backupID, restoreResponse.ID)
		// Should be STARTED, not SUCCESS
		assert.Equal(t, models.BackupRestoreResponseStatusSTARTED, *restoreResponse.Status)
		t.Logf("Restore %s started with status: %s", backupID, *restoreResponse.Status)

		// Wait for completion by polling status
		t.Run("wait for restore completion by polling", func(t *testing.T) {
			statusGetter := client.Backup().RestoreStatusGetter().
				WithBackend(backend).
				WithBackupID(backupID)

			timeout := time.After(30 * time.Second)
			ticker := time.NewTicker(500 * time.Millisecond)
			defer ticker.Stop()

			for {
				select {
				case <-timeout:
					t.Fatal("Timeout waiting for restore completion")
				case <-ticker.C:
					statusResponse, err := statusGetter.Do(ctx)
					require.NoError(t, err, "failed to get restore status")
					require.NotNil(t, statusResponse)

					t.Logf("Restore %s status: %s", backupID, *statusResponse.Status)

					if *statusResponse.Status == models.BackupRestoreStatusResponseStatusSUCCESS {
						t.Logf("Restore %s completed successfully", backupID)
						return
					}

					if *statusResponse.Status == models.BackupRestoreStatusResponseStatusFAILED {
						t.Fatalf("Restore %s failed: %s", backupID, statusResponse.Error)
					}

					// Continue waiting for other statuses (STARTED, TRANSFERRING, etc.)
				}
			}
		})
	})
}
