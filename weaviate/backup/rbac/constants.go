package rbac

// RBACScope defines the scope for Role-Based Access Control in backup operations
type (
	RBACScope string
	UserScope string
)

const (
	// RBACNone excludes all RBAC data from backup/restore operations (default)
	RBACNone RBACScope = "noRestore"

	// RBACAll includes all RBAC settings (roles, users) in backup/restore operations
	RBACAll RBACScope = "all"

	// RBACNone excludes all RBAC data from backup/restore operations (default)
	UserNone UserScope = "noRestore"

	// RBACAll includes all RBAC settings (roles, users) in backup/restore operations
	UserAll UserScope = "all"
)
