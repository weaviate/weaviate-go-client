package api_test

import (
	"encoding/json"
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

// TestRESTEndpoints
// Because of its exhaustive nature, this test doubles as a documentation of the
// REST requests supported by the client and their declarative implementation.
//
// Important: we do not impose any restrictions on how each request implements Body().
// I.e., an endpoint may choose to return a stub from internal/api/gen/rest directly,
// or return some value that will use the stub to implement json.Marshaler.
// Functionally, we only care that the body produces a valid JSON once marshaled.
//
// While it may be tempting to copy-paste the request as the expected body
// (after all, many endpoint implementations will be returning themselves),
// it defies the purpose of this test. Instead, populate wantBody with a stub
// from internal/api/gen/rest package, as it is guaranteed to produce a valid
// JSON, giving you a more useful comparison in the tests.
func TestRESTEndpoints(t *testing.T) {
	for _, tt := range []struct {
		name string
		req  any // Request object.

		wantMethod string     // Expected HTTP Method.
		wantPath   string     // Expected endpoint + path parameters.
		wantQuery  url.Values // Expected query parameters.
		wantBody   any        // Expected request body. JSON strings are compared.
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
		{
			name: "create collection (full config)",
			req: &api.CreateCollectionRequest{
				Collection: api.Collection{
					Name:        "Songs",
					Description: "My favorite songs",
					Properties: []api.Property{
						{Name: "title", DataType: api.DataTypeText},
						{Name: "genres", DataType: api.DataTypeTextArray},
						{Name: "single", DataType: api.DataTypeBool},
						{Name: "year", DataType: api.DataTypeInt},
						{
							Name:              "lyrics",
							DataType:          api.DataTypeInt,
							Tokenization:      api.TokenizationTrigram,
							IndexFilterable:   true,
							IndexRangeFilters: true,
							IndexSearchable:   true,
						},
						{
							Name: "metadata", DataType: api.DataTypeObject,
							NestedProperties: []api.Property{
								{Name: "duration", DataType: api.DataTypeNumber},
								{Name: "uploadedTime", DataType: api.DataTypeDate},
							},
							Tokenization:      api.TokenizationWhitespace,
							IndexFilterable:   true,
							IndexRangeFilters: true,
							IndexSearchable:   true,
						},
					},
					References: []api.ReferenceProperty{
						{
							Name:        "artist",
							Collections: []string{"Singers", "Bands"},
						},
					},
					Sharding: &api.ShardingConfig{
						DesiredCount:        3,
						DesiredVirtualCount: 150,
						VirtualPerPhysical:  50,
					},
					Replication: &api.ReplicationConfig{
						AsyncEnabled:     false,
						Factor:           6,
						DeletionStrategy: api.TimeBasedResolution,
					},
					InvertedIndex: &api.InvertedIndexConfig{
						IndexNullState:         true,
						IndexPropertyLength:    true,
						IndexTimestamps:        true,
						UsingBlockMaxWAND:      true,
						CleanupIntervalSeconds: 92,
						BM25: &api.BM25Config{
							B:  25,
							K1: 1,
						},
						Stopwords: &api.StopwordConfig{
							Preset:    "standard-please-stop",
							Additions: []string{"end"},
							Removals:  []string{"terminate"},
						},
					},
					MultiTenancy: &api.MultiTenancyConfig{
						Enabled:              true,
						AutoTenantActivation: true,
						AutoTenantCreation:   false,
					},
				},
			},
			wantMethod: http.MethodPost,
			wantPath:   "/schema",
			wantBody: &rest.Class{
				Class:       "Songs",
				Description: "My favorite songs",
				Properties: []rest.Property{
					{Name: "title", DataType: []string{string(api.DataTypeText)}},
					{Name: "genres", DataType: []string{string(api.DataTypeTextArray)}},
					{Name: "single", DataType: []string{string(api.DataTypeBool)}},
					{Name: "year", DataType: []string{string(api.DataTypeInt)}},
					{
						Name:              "lyrics",
						DataType:          []string{string(api.DataTypeInt)},
						Tokenization:      rest.PropertyTokenizationTrigram,
						IndexFilterable:   true,
						IndexRangeFilters: true,
						IndexSearchable:   true,
					},
					{
						Name: "metadata", DataType: []string{string(api.DataTypeObject)},
						NestedProperties: []rest.NestedProperty{
							{Name: "duration", DataType: []string{string(api.DataTypeNumber)}},
							{Name: "uploadedTime", DataType: []string{string(api.DataTypeDate)}},
						},
						Tokenization:      rest.PropertyTokenizationWhitespace,
						IndexFilterable:   true,
						IndexRangeFilters: true,
						IndexSearchable:   true,
					},
					{
						Name:     "artist",
						DataType: []string{"Singers", "Bands"},
					},
				},
				ShardingConfig: map[string]any{
					"desiredCount":        3,
					"desiredVirturlCount": 150,
					"virtualPerPhysical":  50,
				},
				ReplicationConfig: rest.ReplicationConfig{
					AsyncEnabled:     false,
					Factor:           6,
					DeletionStrategy: rest.TimeBasedResolution,
				},
				InvertedIndexConfig: rest.InvertedIndexConfig{
					IndexNullState:         true,
					IndexPropertyLength:    true,
					IndexTimestamps:        true,
					UsingBlockMaxWAND:      true,
					CleanupIntervalSeconds: 92,
					Bm25: rest.BM25Config{
						B:  25,
						K1: 1,
					},
					Stopwords: rest.StopwordConfig{
						Preset:    "standard-please-stop",
						Additions: []string{"end"},
						Removals:  []string{"terminate"},
					},
				},
				MultiTenancyConfig: rest.MultiTenancyConfig{
					Enabled:              true,
					AutoTenantActivation: true,
					AutoTenantCreation:   false,
				},
			},
		},
		{
			name: "create collection (partial config)",
			req: &api.CreateCollectionRequest{
				Collection: api.Collection{
					Name:        "Songs",
					Description: "My favorite songs",
					Properties: []api.Property{
						{Name: "title", DataType: api.DataTypeText},
						{Name: "genres", DataType: api.DataTypeTextArray},
						{Name: "single", DataType: api.DataTypeBool},
						{Name: "year", DataType: api.DataTypeInt},
					},
				},
			},
			wantMethod: http.MethodPost,
			wantPath:   "/schema",
			wantBody: &rest.Class{
				Class:       "Songs",
				Description: "My favorite songs",
				Properties: []rest.Property{
					{Name: "title", DataType: []string{string(api.DataTypeText)}},
					{Name: "genres", DataType: []string{string(api.DataTypeTextArray)}},
					{Name: "single", DataType: []string{string(api.DataTypeBool)}},
					{Name: "year", DataType: []string{string(api.DataTypeInt)}},
				},
			},
		},
		{
			name:       "get collection config",
			req:        api.GetCollectionRequest("Songs"),
			wantMethod: http.MethodGet,
			wantPath:   "/schema/Songs",
		},
		{
			name:       "list collections",
			req:        api.ListCollectionsRequest,
			wantMethod: http.MethodGet,
			wantPath:   "/schema",
		},
		{
			name:       "delete collection",
			req:        api.DeleteCollectionRequest("Songs"),
			wantMethod: http.MethodDelete,
			wantPath:   "/schema/Songs",
		},
	} {
		t.Run(fmt.Sprintf("%s (%T)", tt.name, tt.req), func(t *testing.T) {
			require.Implements(t, (*transport.Endpoint)(nil), tt.req)
			endpoint := (tt.req).(transport.Endpoint)

			assert.Equal(t, tt.wantMethod, endpoint.Method(), "bad method")
			assert.Equal(t, tt.wantPath, endpoint.Path(), "bad path")
			assert.Equal(t, tt.wantQuery, endpoint.Query(), "bad query")

			gotBody := endpoint.Body()
			if gotBody == tt.wantBody {
				return // If two objects are equal, so will be their JSON representations.
			}

			gotJSON, err := json.Marshal(gotBody)
			require.NoError(t, err, "marshal request body")

			wantJSON, err := json.Marshal(tt.wantBody)
			require.NoError(t, err, "marshal wantBody")

			assert.JSONEq(t, string(wantJSON), string(gotJSON), "bad body")
		})
	}
}
