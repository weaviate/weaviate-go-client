package api_test

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/gen/rest"
	"github.com/weaviate/weaviate-go-client/v6/internal/transport"
)

func TestRESTEndpoints(t *testing.T) {
	for _, tt := range []struct {
		name string
		req  any // Request object.

		wantMethod string     // Expected HTTP Method.
		wantPath   string     // Expected endpoint + path parameters.
		wantQuery  url.Values // Expected query parameters.
		wantBody   any        // Expected request body.
	}{
		{
			name:       "delete alias",
			req:        api.DeleteAliasRequest("abc"),
			wantMethod: http.MethodDelete,
			wantPath:   "/aliases/abc",
		},
		{
			name: "create backup request",
			req: &api.CreateBackupRequest{
				Backend:            "filesystem",
				ID:                 "bak-1",
				BackupPath:         "/path/to/backup",
				Endpoint:           "s3.amazonaws.com",
				Bucket:             "my-backups",
				IncludeCollections: []string{"Songs"},
				ExcludeCollections: []string{"Pizza"},
				MaxCPUPercentage:   92,
				ChunkSizeMiB:       20,
				CompressionLevel:   api.BackupCompressionLevelDefault,
			},
			wantMethod: http.MethodPost,
			wantPath:   "/backups/filesystem",
			wantBody: &rest.BackupCreateRequest{
				Id:      "bak-1",
				Include: []string{"Songs"},
				Exclude: []string{"Pizza"},
				Config: rest.BackupConfig{
					Path:             "/path/to/backup",
					Bucket:           "my-backups",
					Endpoint:         "s3.amazonaws.com",
					CPUPercentage:    92,
					ChunkSize:        20,
					CompressionLevel: rest.DefaultCompression,
				},
			},
		},
		{
			name: "restore backup request",
			req: &api.RestoreBackupRequest{
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
			wantMethod: http.MethodPost,
			wantPath:   "/backups/filesystem/bak-1/restore",
			wantBody: &rest.BackupRestoreRequest{
				Include:        []string{"Songs"},
				Exclude:        []string{"Pizza"},
				OverwriteAlias: true,
				NodeMapping:    map[string]string{"node-1": "node-a"},
				Config: rest.RestoreConfig{
					Path:          "/path/to/backup",
					Bucket:        "my-backups",
					Endpoint:      "s3.amazonaws.com",
					CPUPercentage: 92,
					UsersOptions:  rest.All,
					RolesOptions:  rest.RestoreConfigRolesOptionsNoRestore,
				},
			},
		},
		{
			name: "get backup create status",
			req: &api.BackupStatusRequest{
				Backend:   "filesystem",
				ID:        "bak-1",
				Operation: api.CreateBackup,
			},
			wantMethod: http.MethodGet,
			wantPath:   "/backups/filesystem/bak-1",
		},
		{
			name: "get backup restore status",
			req: &api.BackupStatusRequest{
				Backend:   "filesystem",
				ID:        "bak-1",
				Operation: api.RestoreBackup,
			},
			wantMethod: http.MethodGet,
			wantPath:   "/backups/filesystem/bak-1/restore",
		},
		{
			name: "list backups",
			req: &api.ListBackupsRequest{
				Backend: "filesystem",
			},
			wantMethod: http.MethodGet,
			wantPath:   "/backups/filesystem",
		},
		{
			name: "list backups order by starting time",
			req: &api.ListBackupsRequest{
				Backend:         "filesystem",
				StartingTimeAsc: true,
			},
			wantMethod: http.MethodGet,
			wantPath:   "/backups/filesystem",
			wantQuery:  url.Values{"order": {"asc"}},
		},
		{
			name: "cancel backup",
			req: &api.CancelBackupRequest{
				Backend: "filesystem",
				ID:      "bak-1",
			},
			wantMethod: http.MethodDelete,
			wantPath:   "/backups/filesystem/bak-1",
		},
	} {
		t.Run(fmt.Sprintf("%s (%T)", tt.name, tt.req), func(t *testing.T) {
			require.Implements(t, (*transport.Endpoint)(nil), tt.req)
			endpoint := (tt.req).(transport.Endpoint)

			assert.Equal(t, tt.wantMethod, endpoint.Method(), "bad method")
			assert.Equal(t, tt.wantPath, endpoint.Path(), "bad path")
			assert.Equal(t, tt.wantQuery, endpoint.Query(), "bad query")
			assert.Equal(t, tt.wantBody, endpoint.Body(), "bad body")
		})
	}
}
