# RBAC Integration Tests

This directory contains comprehensive integration tests for RBAC (Role-Based Access Control) functionality in the Weaviate Go client, specifically focusing on backup and restore operations with RBAC data.

## Test Files

1. **`rbac_backup_integration_test.go`** - Main integration tests that mirror the CLI operations from the shell script:
   - Basic RBAC backup and restore workflow
   - Different RBAC backup configurations
   - Error handling and edge cases
   - Backup creation and validation patterns

2. **`rbac_advanced_test.go`** - Advanced RBAC scenarios:
   - Complex role hierarchies
   - Incremental backup strategies
   - Comprehensive user and role management

## Test Coverage

These tests cover the same functionality as described in the original shell script:

### Core Operations Tested
- ✅ User creation and deletion (`CLI user create testuser testkey` / `CLI user delete testuser`)
- ✅ Class creation and deletion (`CLI class create Article` / `CLI class delete Article`)
- ✅ Multiple backup creation (`CLI backup create rbactest1..10`)
- ✅ Backup status verification (`CLI backup status rbactest1`)
- ✅ Backup restoration (`CLI backup restore rbactest1 --roles all --users all`)
- ✅ Restore status verification (`CLI backup restore-status rbactest1`)
- ✅ Verification of restored data (classes and users exist after restore)

### RBAC Configuration Options Tested
- `WithRBACAll()` - Include all RBAC data (equivalent to `--roles all --users all`)
- `WithRBACNone()` - Exclude all RBAC data
- `WithRBACSpecificRoles(...)` - Include only specific roles
- `WithRBACRolesOnly()` - Include role definitions only
- `WithRBACUsersOnly()` - Include user assignments only
- `WithRBAC(config)` - Custom RBAC configuration
- `WithRBACFlags(flags)` - Use bitwise flags (e.g., `rbac.Roles | rbac.Users`)
- `WithRBACFlags(...).WithSpecificRoles(...)` - Combine flags with specific roles

### Error Scenarios Tested
- Non-existent roles in backup configuration
- Invalid backend configurations
- Conflicting RBAC configurations between backup and restore
- Timeout and polling scenarios

## Running the Tests

### Prerequisites
- Docker and Docker Compose installed
- Weaviate server with RBAC enabled
- Go 1.19+ installed

### Running Individual Tests

```bash
# Run the main integration test
go test -v ./test/rbac_integration -run TestRBACBackupIntegration

# Run different configuration tests
go test -v ./test/rbac_integration -run TestRBACBackupDifferentConfigurations

# Run error handling tests
go test -v ./test/rbac_integration -run TestRBACBackupErrorHandling

# Run advanced scenarios
go test -v ./test/rbac_integration -run TestRBACBackupWithUserRoleManagement

# Run all RBAC integration tests
go test -v ./test/rbac_integration
```

### Running with Custom Environment

```bash
# Set up test environment variables if needed
export WEAVIATE_VERSION=1.25.0
export INTEGRATION_TESTS_AUTH=auth_enabled

# Run tests
go test -v ./test/rbac_integration
```

## Test Environment Setup

The tests use the existing test infrastructure with RBAC-enabled containers:

```yaml
# From docker-compose-rbac.yml
services:
  weaviate-rbac:
    environment:
      AUTHENTICATION_APIKEY_ENABLED: 'true'
      AUTHORIZATION_RBAC_ENABLED: 'true'
      AUTHENTICATION_APIKEY_ALLOWED_KEYS: 'my-secret-key'
      AUTHENTICATION_APIKEY_USERS: 'adam-the-admin'
      AUTHORIZATION_ADMIN_USERS: 'adam-the-admin'
      AUTHENTICATION_DB_USERS_ENABLED: "true"
      BACKUP_FILESYSTEM_PATH: "/tmp/backups"
      ENABLE_MODULES: "backup-filesystem"
```

## Key Test Patterns

### 1. Basic Integration Test Pattern
```go
// Mirrors the CLI script operations:
// 1. Create user and class
// 2. Create multiple backups with RBAC data
// 3. Delete user and class
// 4. Restore backup
// 5. Verify user and class exist after restore
```

### 2. Configuration Testing Pattern
```go
// Tests different RBAC backup configurations:
tests := []struct {
    name         string
    setupBackup  func(*backup.BackupCreator) *backup.BackupCreator
    setupRestore func(*backup.BackupRestorer) *backup.BackupRestorer
    description  string
}{
    // Convenience methods
    creator.WithRBACAll()
    creator.WithRBACNone()
    creator.WithRBACSpecificRoles("role1", "role2")
    
    // Bitwise flags (new API)
    creator.WithRBACFlags(rbac.None)
    creator.WithRBACFlags(rbac.Roles | rbac.Users)
    creator.WithRBACFlags(rbac.Roles | rbac.Users).WithSpecificRoles("admin")
    
    // Custom configuration
    creator.WithRBAC(&rbac.RBACConfig{...})
}
```

### 3. Error Handling Pattern
```go
// Tests error conditions and edge cases:
// - Non-existent roles
// - Invalid backends
// - Conflicting configurations
```

## Comparison with Shell Script

| Shell Script Operation | Go Test Equivalent |
|----------------------|-------------------|
| `$CLI $AUTH user create testuser testkey` | `client.Users().DB().Creator().WithUserID(testUserID).Do(ctx)` |
| `$CLI $AUTH class create Article` | `client.Schema().ClassCreator().WithClass(class).Do(ctx)` |
| `$CLI $AUTH backup create rbactest1 --backend s3` | `client.Backup().Creator().WithBackend(backend).WithBackupID(backupID).WithRBACAll().Do(ctx)` |
| `$CLI $AUTH backup status rbactest1` | `client.Backup().CreateStatusGetter().WithBackend(backend).WithBackupID(backupID).Do(ctx)` |
| `$CLI $AUTH backup restore rbactest1 --roles all --users all` | `client.Backup().Restorer().WithBackend(backend).WithBackupID(backupID).WithRBACAll().Do(ctx)` |
| `$CLI $AUTH user delete testuser` | `client.Users().DB().Deleter().WithUserID(testUserID).Do(ctx)` |
| `$CLI $AUTH class delete Article` | `client.Schema().ClassDeleter().WithClassName(testClassName).Do(ctx)` |

## RBAC API Examples

The tests demonstrate multiple ways to configure RBAC backup/restore:

```go
// Convenience methods (most common)
client.Backup().Creator().WithRBACAll()                    // All RBAC data
client.Backup().Creator().WithRBACNone()                   // No RBAC data  
client.Backup().Creator().WithRBACRolesOnly()              // Roles only
client.Backup().Creator().WithRBACUsersOnly()              // Users only
client.Backup().Creator().WithRBACSpecificRoles("admin")   // Specific roles

// Bitwise flags (new pattern)
client.Backup().Creator().WithRBACFlags(rbac.None)                        // No RBAC
client.Backup().Creator().WithRBACFlags(rbac.Roles | rbac.Users)          // Roles + Users
client.Backup().Creator().WithRBACFlags(rbac.Roles).WithSpecificRoles("admin")  // Specific roles

// Custom configuration (advanced)
client.Backup().Creator().WithRBAC(&rbac.RBACConfig{
    Scope:                   rbac.RBACAll,
    RoleSelection:          rbac.RoleSelectionSpecific,
    SpecificRoles:          []string{"admin", "user"},
    IncludeUserAssignments: true,
    IncludeGroupAssignments: false,
})
```

## Expected Output

When running the tests, you should see output similar to:

```
=== RUN   TestRBACBackupIntegration
=== RUN   TestRBACBackupIntegration/RBAC_Backup_and_Restore_Integration_Test
=== RUN   TestRBACBackupIntegration/RBAC_Backup_and_Restore_Integration_Test/create_test_user
    rbac_backup_integration_test.go:48: Created user 'testuser' with API key
=== RUN   TestRBACBackupIntegration/RBAC_Backup_and_Restore_Integration_Test/create_test_class
    rbac_backup_integration_test.go:75: Created class 'Article'
=== RUN   TestRBACBackupIntegration/RBAC_Backup_and_Restore_Integration_Test/create_backup_rbactest1
    rbac_backup_integration_test.go:119: Successfully created backup rbactest1
=== RUN   TestRBACBackupIntegration/RBAC_Backup_and_Restore_Integration_Test/delete_test_user
    rbac_backup_integration_test.go:149: Successfully deleted user 'testuser'
=== RUN   TestRBACBackupIntegration/RBAC_Backup_and_Restore_Integration_Test/restore_backup_rbactest1
    rbac_backup_integration_test.go:187: Successfully restored backup rbactest1
=== RUN   TestRBACBackupIntegration/RBAC_Backup_and_Restore_Integration_Test/verify_class_exists_after_restore
    rbac_backup_integration_test.go:215: Class 'Article' exists after restore
=== RUN   TestRBACBackupIntegration/RBAC_Backup_and_Restore_Integration_Test/verify_user_exists_after_restore
    rbac_backup_integration_test.go:231: User 'testuser' exists after restore
--- PASS: TestRBACBackupIntegration (XX.XXs)
```

## Notes

- Tests use the `testenv.SetupLocalContainer` with `test.RBAC` flag to ensure RBAC is enabled
- All tests include proper cleanup using `t.Cleanup()` functions
- Tests wait for operations to complete using `WithWaitForCompletion(true)` or manual polling
- Error cases are tested to ensure robust behavior
- Tests mirror the CLI script behavior but use the Go client API instead

## Troubleshooting

If tests fail:

1. **Check Docker**: Ensure Docker is running and can pull Weaviate images
2. **Check Ports**: Ensure test ports (8089, etc.) are available
3. **Check Logs**: Look at container logs for RBAC configuration issues
4. **Check Permissions**: Ensure the test user has proper admin permissions
5. **Check Filesystem**: Ensure backup filesystem path is writable in the container

For more debugging, you can enable verbose logging:

```bash
go test -v -args -test.v ./test/rbac_integration
```
