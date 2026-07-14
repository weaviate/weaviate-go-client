package api_test

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/rest"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/internal/transports"
)

// TestRESTRequests verifies the parameters of REST requests provided by the 'api' package.
// Because of its exhaustive nature, this test doubles as a documentation
// of the REST requests supported by the client.
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
func TestRESTRequests(t *testing.T) {
	for _, tt := range testkit.WithOnly(t, []struct {
		testkit.Only

		name string
		req  any // Request object.

		wantMethod string     // Expected HTTP Method.
		wantPath   string     // Expected endpoint + path parameters.
		wantQuery  url.Values // Expected query parameters.
		wantBody   any        // Expected request body. JSON strings are compared.
	}{
		{
			name:       "check is live",
			req:        api.IsLiveRequest,
			wantMethod: http.MethodGet,
			wantPath:   "/.well-known/live",
		},
		{
			name:       "check is ready",
			req:        api.IsReadyRequest,
			wantMethod: http.MethodGet,
			wantPath:   "/.well-known/ready",
		},
		{
			name:       "get instance metadata",
			req:        api.GetInstanceMetadataRequest,
			wantMethod: http.MethodGet,
			wantPath:   "/meta",
		},
		{
			name: "replace object (no consistency_level)",
			req: &api.ReplaceObjectRequest{
				RequestDefaults: api.RequestDefaults{
					CollectionName: "Songs",
					Tenant:         "john_doe",
				},
				UUID: &testkit.UUID,
				Properties: map[string]any{
					"title":  "High Speed Dirt",
					"genres": []string{"thrash metal", "blues"},
					"single": false,
					"year":   1992,
				},
				References: api.References{
					"band": {
						{Target: api.ObjectPath{UUID: testkit.UUID, Collection: "Drummers"}},
						{Target: api.ObjectPath{UUID: testkit.UUID, Collection: "Basists"}},
					},
					"label": {
						{Target: api.ObjectPath{UUID: testkit.UUID}},
					},
				},
				Vectors: []api.Vector{
					{Name: "lyrics", Single: []float32{1, 2, 3}},
				},
			},
			wantMethod: http.MethodPut,
			wantPath:   "/objects/Songs/" + testkit.UUID.String(),
			wantBody: &rest.Object{
				Tenant: "john_doe",
				Properties: map[string]any{
					"title":  "High Speed Dirt",
					"genres": []string{"thrash metal", "blues"},
					"single": false,
					"year":   1992,
					"band": []string{
						"weaviate://localhost/Drummers/" + testkit.UUID.String(),
						"weaviate://localhost/Basists/" + testkit.UUID.String(),
					},
					"label": []string{
						"weaviate://localhost/" + testkit.UUID.String(),
					},
				},
				Vectors: map[string]any{
					"lyrics": []float32{1, 2, 3},
				},
			},
		},
		{
			name: "replace object (consistency_level=ONE)",
			req: &api.ReplaceObjectRequest{
				RequestDefaults: api.RequestDefaults{
					CollectionName:   "Songs",
					ConsistencyLevel: api.ConsistencyLevelOne,
				},
				UUID: &testkit.UUID,
			},
			wantMethod: http.MethodPut,
			wantPath:   "/objects/Songs/" + testkit.UUID.String(),
			wantQuery:  url.Values{"consistency_level": {string(api.ConsistencyLevelOne)}},
			wantBody:   &rest.Object{},
		},
		{
			name: "delete object (no consistency_level)",
			req: &api.DeleteObjectRequest{
				RequestDefaults: api.RequestDefaults{
					CollectionName: "Songs",
					Tenant:         "john_doe",
				},
				UUID: testkit.UUID,
			},
			wantMethod: http.MethodDelete,
			wantPath:   "/objects/Songs/" + testkit.UUID.String(),
			wantQuery:  url.Values{"tenant": {"john_doe"}},
		},
		{
			name: "delete object (no tenant)",
			req: &api.DeleteObjectRequest{
				RequestDefaults: api.RequestDefaults{
					CollectionName:   "Songs",
					ConsistencyLevel: api.ConsistencyLevelOne,
				},
				UUID: testkit.UUID,
			},
			wantMethod: http.MethodDelete,
			wantPath:   "/objects/Songs/" + testkit.UUID.String(),
			wantQuery:  url.Values{"consistency_level": {string(api.ConsistencyLevelOne)}},
		},
		{
			name: "delete object (no tenant, no consistency_level)",
			req: &api.DeleteObjectRequest{
				RequestDefaults: api.RequestDefaults{CollectionName: "Songs"},
				UUID:            testkit.UUID,
			},
			wantMethod: http.MethodDelete,
			wantPath:   "/objects/Songs/" + testkit.UUID.String(),
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
					Vectors: map[string]api.VectorConfig{
						"lyrics_vec": {
							Index: &api.Module{
								Name: "hfresh",
								Conf: map[string]any{
									"distance": "cosine",
								},
							},
							Compression: &api.Module{
								Name: "rq",
								Conf: map[string]any{
									"enabled":       true,
									"cache":         true,
									"bits":          1,
									"rescore_limit": 16,
								},
							},
							Vectorizer: &api.Module{
								Name: testkit.ModuleName,
								Conf: map[string]any{
									"url": "example.com",
								},
							},
						},
						"title_vec": {},
						"audio_vec": {SkipDefaultCompression: true},
					},
					Sharding: &api.ShardingConfig{
						DesiredCount:        3,
						DesiredVirtualCount: 150,
						VirtualPerPhysical:  50,
					},
					Replication: &api.ReplicationConfig{
						Factor:           6,
						DeletionStrategy: api.TimeBasedResolution,
						AsyncReplication: &api.AsyncReplicationConfig{
							DiffBatchSize:                   1,
							DiffPerNodeTimeout:              2 * time.Second,
							ReplicationFrequency:            3 * time.Millisecond,
							ReplicationFrequencyPropagating: 4 * time.Millisecond,
							PrePropagationTimeout:           5 * time.Second,
							PropagationConcurrency:          6,
							PropagationBatchSize:            7,
							PropagationLimit:                8,
							PropagationTimeout:              9 * time.Second,
							PropagationDelay:                10 * time.Millisecond,
							HashTreeHeight:                  11,
							LoggingFrequency:                12 * time.Second,
						},
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
				VectorConfig: map[string]rest.VectorConfig{
					"lyrics_vec": {
						VectorIndexType: "hfresh",
						VectorIndexConfig: map[string]any{
							"distance": "cosine",
							"rq": map[string]any{
								"enabled":       true,
								"cache":         true,
								"bits":          1,
								"rescore_limit": 16,
							},
						},
						Vectorizer: map[string]any{
							testkit.ModuleName: map[string]any{
								"url": "example.com",
							},
						},
					},
					"title_vec": {},
					"audio_vec": {
						VectorIndexConfig: map[string]any{
							"skipDefaultQuantization": true,
						},
						Vectorizer: map[string]any{
							"none": map[string]any{},
						},
					},
				},
				ShardingConfig: map[string]any{
					"desiredCount":        3,
					"desiredVirtualCount": 150,
					"virtualPerPhysical":  50,
				},
				ReplicationConfig: rest.ReplicationConfig{
					Factor:           6,
					DeletionStrategy: rest.TimeBasedResolution,
					AsyncConfig: rest.ReplicationAsyncConfig{
						DiffBatchSize:             1,
						DiffPerNodeTimeout:        2,
						Frequency:                 3,
						FrequencyWhilePropagating: 4,
						PrePropagationTimeout:     5,
						PropagationConcurrency:    6,
						PropagationBatchSize:      7,
						PropagationLimit:          8,
						PropagationTimeout:        9,
						PropagationDelay:          10,
						HashtreeHeight:            11,
						LoggingFrequency:          12,
					},
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
				PrefixIncremental:  "incr-bak-",
				MaxCPUPercentage:   92,
				ChunkSizeMiB:       20,
				CompressionLevel:   api.BackupCompressionLevelDefault,
			},
			wantMethod: http.MethodPost,
			wantPath:   "/backups/filesystem",
			wantBody: &rest.BackupCreateRequest{
				Id:                      "bak-1",
				Include:                 []string{"Songs"},
				Exclude:                 []string{"Pizza"},
				IncrementalBaseBackupId: "incr-bak-",
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
				Operation: api.BackupOperationCreate,
			},
			wantMethod: http.MethodGet,
			wantPath:   "/backups/filesystem/bak-1",
		},
		{
			name: "get backup restore status",
			req: &api.BackupStatusRequest{
				Backend:   "filesystem",
				ID:        "bak-1",
				Operation: api.BackupOperationRestore,
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
			name: "cancel backup create",
			req: &api.CancelBackupRequest{
				Backend:   "filesystem",
				ID:        "bak-1",
				Operation: api.BackupOperationCreate,
			},
			wantMethod: http.MethodDelete,
			wantPath:   "/backups/filesystem/bak-1",
		},
		{
			name: "cancel backup restore",
			req: &api.CancelBackupRequest{
				Backend:   "filesystem",
				ID:        "bak-1",
				Operation: api.BackupOperationRestore,
			},
			wantMethod: http.MethodDelete,
			wantPath:   "/backups/filesystem/bak-1/restore",
		},
		{
			name: "create role",
			req: &api.CreateRoleRequest{
				Role: api.Role{
					ID: "rock-n-role",
					Permissions: api.Permissions{
						Aliases: []api.AliasPermission{
							{
								Collection: "Songs", Alias: "Records",
								Create: true, Read: true,
							},
							{
								Collection: "Artists", Alias: "Musicians",
								Update: true, Delete: true,
							},
						},
						Backups: []api.BackupsPermission{
							{Collection: "Songs", Manage: true},
						},
						Cluster: []api.ClusterPermission{{Read: true}},
						Collections: []api.CollectionPermission{
							{Collection: "Songs", Create: true, Read: true},
							{Collection: "Artists", Update: true, Delete: true},
						},
						Data: []api.DataPermission{
							{
								Collection: "Songs", Tenant: "john_doe",
								Create: true, Read: true,
							},
							{
								Collection: "Artists",
								Update:     true, Delete: true,
							},
						},
						Groups: []api.GroupPermission{
							{
								GroupID: "external", Type: "oidc",
								Read: true, AssignAndRevoke: true,
							},
						},
						Namespaces: []api.NamespacePermission{
							{Namespace: "sandbox", Manage: true},
						},
						Nodes: []api.NodesPermission{
							{Collection: "Songs", Verbosity: api.NodeVerbosityMinimal, Read: true},
							{Collection: "Artists", Verbosity: api.NodeVerbosityVerbose, Read: true},
						},
						MCP: []api.MCPPermission{
							{Create: true, Read: true, Update: true},
						},
						Replication: []api.ReplicationPermission{
							{
								Collection: "Songs", Shard: "abc",
								Create: true, Read: true,
							},
							{
								Collection: "Songs", Shard: "xyz",
								Update: true, Delete: true,
							},
						},
						Roles: []api.RolePermission{
							{
								RoleID: "rock-n-role", Scope: api.RoleScopeAll,
								Create: true, Read: true,
							},
							{
								RoleID: "rock-n-role", Scope: api.RoleScopeMatch,
								Update: true, Delete: true,
							},
						},
						Tenants: []api.TenantPermission{
							{
								Collection: "Songs", Tenant: "john_doe",
								Create: true, Read: true,
							},
							{
								Collection: "Artists", Tenant: "*",
								Update: true, Delete: true,
							},
						},
						Users: []api.UserPermission{
							{UserID: "john_doe", Create: true, Read: true},
							{UserID: "jane_doe", Update: true, Delete: true},
							{UserID: "jim_beam", AssignAndRevoke: true},
						},
					},
				},
			},
			wantMethod: http.MethodPost,
			wantPath:   "/authz/roles",
			wantBody: map[string]any{
				"name": "rock-n-role",
				"permissions": []map[string]any{
					{
						"action": "create_aliases",
						"aliases": map[string]any{
							"collection": "Songs",
							"alias":      "Records",
						},
					},
					{
						"action": "read_aliases",
						"aliases": map[string]any{
							"collection": "Songs",
							"alias":      "Records",
						},
					},
					{
						"action": "update_aliases",
						"aliases": map[string]any{
							"collection": "Artists",
							"alias":      "Musicians",
						},
					},
					{
						"action": "delete_aliases",
						"aliases": map[string]any{
							"collection": "Artists",
							"alias":      "Musicians",
						},
					},
					{
						"action": "manage_backups",
						"backups": map[string]any{
							"collection": "Songs",
						},
					},
					{
						"action": "read_cluster",
					},
					{
						"action": "create_collections",
						"collections": map[string]any{
							"collection": "Songs",
						},
					},
					{
						"action": "read_collections",
						"collections": map[string]any{
							"collection": "Songs",
						},
					},
					{
						"action": "update_collections",
						"collections": map[string]any{
							"collection": "Artists",
						},
					},
					{
						"action": "delete_collections",
						"collections": map[string]any{
							"collection": "Artists",
						},
					},
					{
						"action": "create_data",
						"data": map[string]any{
							"collection": "Songs",
							"tenant":     "john_doe",
						},
					},
					{
						"action": "read_data",
						"data": map[string]any{
							"collection": "Songs",
							"tenant":     "john_doe",
						},
					},
					{
						"action": "update_data",
						"data": map[string]any{
							"collection": "Artists",
						},
					},
					{
						"action": "delete_data",
						"data": map[string]any{
							"collection": "Artists",
						},
					},
					{
						"action": "read_groups",
						"groups": map[string]any{
							"group":     "external",
							"groupType": "oidc",
						},
					},
					{
						"action": "assign_and_revoke_groups",
						"groups": map[string]any{
							"group":     "external",
							"groupType": "oidc",
						},
					},
					{
						"action": "manage_namespaces",
						"namespaces": map[string]any{
							"namespace": "sandbox",
						},
					},
					{
						"action": "read_nodes",
						"nodes": map[string]any{
							"collection": "Songs",
							"verbosity":  api.NodeVerbosityMinimal,
						},
					},
					{
						"action": "read_nodes",
						"nodes": map[string]any{
							"collection": "Artists",
							"verbosity":  api.NodeVerbosityVerbose,
						},
					},
					{
						"action": "create_mcp",
					},
					{
						"action": "read_mcp",
					},
					{
						"action": "update_mcp",
					},
					{
						"action": "create_replicate",
						"replicate": map[string]any{
							"collection": "Songs",
							"shard":      "abc",
						},
					},
					{
						"action": "read_replicate",
						"replicate": map[string]any{
							"collection": "Songs",
							"shard":      "abc",
						},
					},
					{
						"action": "update_replicate",
						"replicate": map[string]any{
							"collection": "Songs",
							"shard":      "xyz",
						},
					},
					{
						"action": "delete_replicate",
						"replicate": map[string]any{
							"collection": "Songs",
							"shard":      "xyz",
						},
					},
					{
						"action": "create_roles",
						"roles": map[string]any{
							"role":  "rock-n-role",
							"scope": api.RoleScopeAll,
						},
					},
					{
						"action": "read_roles",
						"roles": map[string]any{
							"role":  "rock-n-role",
							"scope": api.RoleScopeAll,
						},
					},
					{
						"action": "update_roles",
						"roles": map[string]any{
							"role":  "rock-n-role",
							"scope": api.RoleScopeMatch,
						},
					},
					{
						"action": "delete_roles",
						"roles": map[string]any{
							"role":  "rock-n-role",
							"scope": api.RoleScopeMatch,
						},
					},
					{
						"action": "create_tenants",
						"tenants": map[string]any{
							"collection": "Songs",
							"tenant":     "john_doe",
						},
					},
					{
						"action": "read_tenants",
						"tenants": map[string]any{
							"collection": "Songs",
							"tenant":     "john_doe",
						},
					},
					{
						"action": "update_tenants",
						"tenants": map[string]any{
							"collection": "Artists",
							"tenant":     "*",
						},
					},
					{
						"action": "delete_tenants",
						"tenants": map[string]any{
							"collection": "Artists",
							"tenant":     "*",
						},
					},
					{
						"action": "create_users",
						"users": map[string]any{
							"user": "john_doe",
						},
					},
					{
						"action": "read_users",
						"users": map[string]any{
							"user": "john_doe",
						},
					},
					{
						"action": "update_users",
						"users": map[string]any{
							"user": "jane_doe",
						},
					},
					{
						"action": "delete_users",
						"users": map[string]any{
							"user": "jane_doe",
						},
					},
					{
						"action": "assign_and_revoke_users",
						"users": map[string]any{
							"user": "jim_beam",
						},
					},
				},
			},
		},
		{
			name:       "get role",
			req:        api.GetRoleRequest("rock-n-role"),
			wantMethod: http.MethodGet,
			wantPath:   "/authz/roles/rock-n-role",
		},
		{
			name:       "delete role",
			req:        api.DeleteRoleRequest("rock-n-role"),
			wantMethod: http.MethodDelete,
			wantPath:   "/authz/roles/rock-n-role",
		},
		{
			name:       "list roles",
			req:        api.ListRolesRequest,
			wantMethod: http.MethodGet,
			wantPath:   "/authz/roles",
		},
		{
			name:       "get assigned users",
			req:        api.GetAssignedUsersRequest("rock-n-role"),
			wantMethod: http.MethodGet,
			wantPath:   "/authz/roles/rock-n-role/users",
		},
		{
			name:       "get user assignments",
			req:        api.GetUserAssignmentsRequest("rock-n-role"),
			wantMethod: http.MethodGet,
			wantPath:   "/authz/roles/rock-n-role/user-assignments",
		},
		{
			name:       "get group assignments",
			req:        api.GetGroupAssignmentsRequest("rock-n-role"),
			wantMethod: http.MethodGet,
			wantPath:   "/authz/roles/rock-n-role/group-assignments",
		},
		{
			name: "add permissions to role",
			req: &api.ManagePermissionsRequest{
				RoleID: "rock-n-role",
				Verb:   api.PermissionVerbAdd,
				Permissions: api.Permissions{
					Cluster: []api.ClusterPermission{
						{Read: true},
					},
				},
			},
			wantMethod: http.MethodPost,
			wantPath:   "/authz/roles/rock-n-role/add-permissions",
			wantBody: map[string]any{
				"permissions": []map[string]any{
					{"action": "read_cluster"},
				},
			},
		},
		{
			name: "remove permissions from role",
			req: &api.ManagePermissionsRequest{
				RoleID: "rock-n-role",
				Verb:   api.PermissionVerbRemove,
				Permissions: api.Permissions{
					Cluster: []api.ClusterPermission{
						{Read: true},
					},
				},
			},
			wantMethod: http.MethodPost,
			wantPath:   "/authz/roles/rock-n-role/remove-permissions",
			wantBody: map[string]any{
				"permissions": []map[string]any{
					{"action": "read_cluster"},
				},
			},
		},
		{
			name: "check role has permission",
			req: &api.HasPermissionRequest{
				RoleID:  "rock-n-role",
				Cluster: api.ClusterPermission{Read: true},
			},
			wantMethod: http.MethodPost,
			wantPath:   "/authz/roles/rock-n-role/has-permission",
			wantBody: map[string]any{
				"action": "read_cluster",
			},
		},
		{
			name:       "get own user info",
			req:        api.GetOwnUserInfoRequest,
			wantMethod: http.MethodGet,
			wantPath:   "/users/own-info",
		},
		{
			name: "get roles assigned to the user (include permissions)",
			req: &api.GetAssignedRolesRequest{
				Entity:             api.UserID("john-malkovich"),
				Kind:               api.RBACKindDB,
				IncludePermissions: true,
			},
			wantMethod: http.MethodGet,
			wantPath:   "/authz/users/john-malkovich/roles/db",
			wantQuery:  url.Values{"includePermissions": {"true"}},
		},
		{
			name: "get roles assigned to the group (include permissions)",
			req: &api.GetAssignedRolesRequest{
				Entity:             api.GroupID("external"),
				Kind:               api.RBACKindOIDC,
				IncludePermissions: true,
			},
			wantMethod: http.MethodGet,
			wantPath:   "/authz/groups/external/roles/oidc",
			wantQuery:  url.Values{"includePermissions": {"true"}},
		},
		{
			name: "assign roles to user",
			req: &api.ManageRolesRequest{
				Entity: api.UserID("john-malkovich"),
				Kind:   api.RBACKindDB,
				Verb:   api.RoleVerbAssign,
				Roles:  []string{"rock-n-role", "sushi-role"},
			},
			wantMethod: http.MethodPost,
			wantPath:   "/authz/users/john-malkovich/assign",
			wantBody: rest.AssignRoleToUserJSONBody{
				UserType: rest.UserTypeInputDb,
				Roles:    []string{"rock-n-role", "sushi-role"},
			},
		},
		{
			name: "assign roles to group",
			req: &api.ManageRolesRequest{
				Entity: api.GroupID("external"),
				Kind:   api.RBACKindOIDC,
				Verb:   api.RoleVerbAssign,
				Roles:  []string{"rock-n-role", "sushi-role"},
			},
			wantMethod: http.MethodPost,
			wantPath:   "/authz/groups/external/assign",
			wantBody: rest.AssignRoleToGroupJSONBody{
				GroupType: rest.GroupTypeOidc,
				Roles:     []string{"rock-n-role", "sushi-role"},
			},
		},
		{
			name: "revoke roles from user",
			req: &api.ManageRolesRequest{
				Entity: api.UserID("john-malkovich"),
				Kind:   api.RBACKindDB,
				Verb:   api.RoleVerbRevoke,
				Roles:  []string{"rock-n-role", "sushi-role"},
			},
			wantMethod: http.MethodPost,
			wantPath:   "/authz/users/john-malkovich/revoke",
			wantBody: rest.RevokeRoleFromUserJSONBody{
				UserType: rest.UserTypeInputDb,
				Roles:    []string{"rock-n-role", "sushi-role"},
			},
		},
		{
			name: "revoke roles from group",
			req: &api.ManageRolesRequest{
				Entity: api.GroupID("external"),
				Kind:   api.RBACKindOIDC,
				Verb:   api.RoleVerbRevoke,
				Roles:  []string{"rock-n-role", "sushi-role"},
			},
			wantMethod: http.MethodPost,
			wantPath:   "/authz/groups/external/revoke",
			wantBody: rest.RevokeRoleFromGroupJSONBody{
				GroupType: rest.GroupTypeOidc,
				Roles:     []string{"rock-n-role", "sushi-role"},
			},
		},
		{
			name:       "list groups",
			req:        api.ListGroupsRequest,
			wantMethod: http.MethodGet,
			wantPath:   "/authz/groups/oidc",
		},
		{
			name:       "create db user",
			req:        api.CreateUserRequest("john-malkovich"),
			wantMethod: http.MethodPost,
			wantPath:   "/users/db/john-malkovich",
		},
		{
			name:       "delete db user",
			req:        api.DeleteUserRequest("john-malkovich"),
			wantMethod: http.MethodDelete,
			wantPath:   "/users/db/john-malkovich",
		},
		{
			name:       "activate db user",
			req:        api.ActivateUserRequest("john-malkovich"),
			wantMethod: http.MethodPost,
			wantPath:   "/users/db/john-malkovich/activate",
		},
		{
			name: "deactivate db user (revoke api key)",
			req: &api.DeactivateUserRequest{
				ID:        "john-malkovich",
				RevokeKey: true,
			},
			wantMethod: http.MethodPost,
			wantPath:   "/users/db/john-malkovich/deactivate",
			wantBody: rest.DeactivateUserJSONRequestBody{
				RevokeKey: true,
			},
		},
		{
			name:       "rotate db user api key",
			req:        api.RotateAPIKeyRequest("john-malkovich"),
			wantMethod: http.MethodPost,
			wantPath:   "/users/db/john-malkovich/rotate-key",
		},
		{
			name: "get db user info",
			req: &api.GetUserInfoRequest{
				ID:                "john-malkovich",
				IncludeLastUsedAt: true,
			},
			wantMethod: http.MethodGet,
			wantPath:   "/users/db/john-malkovich",
			wantQuery:  url.Values{"includeLastUsedTime": {"true"}},
		},
		{
			name: "list db users",
			req: &api.ListUsersRequest{
				IncludeLastUsedAt: true,
			},
			wantMethod: http.MethodGet,
			wantPath:   "/users/db",
			wantQuery:  url.Values{"includeLastUsedTime": {"true"}},
		},
		{
			name:       "list aliases",
			req:        api.ListAliasesRequest,
			wantMethod: http.MethodGet,
			wantPath:   "/aliases",
		},
		{
			name: "create alias",
			req: &api.CreateAliasRequest{
				Alias: api.Alias{
					Collection: "GeorgeBarnes",
					Alias:      "MachineGunKelly",
				},
			},
			wantMethod: http.MethodPost,
			wantPath:   "/aliases",
			wantBody: rest.AliasesCreateJSONRequestBody{
				Class: "GeorgeBarnes",
				Alias: "MachineGunKelly",
			},
		},
		{
			name:       "get an alias",
			req:        api.GetAliasRequest("MachineGunKelly"),
			wantMethod: http.MethodGet,
			wantPath:   "/aliases/MachineGunKelly",
		},
		{
			name: "update an alias",
			req: &api.UpdateAliasRequest{
				Alias: api.Alias{
					Alias:      "MachineGunKelly",
					Collection: "ColsonBaker",
				},
			},
			wantMethod: http.MethodPut,
			wantPath:   "/aliases/MachineGunKelly",
			wantBody: rest.AliasesUpdateJSONRequestBody{
				Class: "ColsonBaker",
			},
		},
		{
			name:       "delete an alias",
			req:        api.DeleteAliasRequest("MachineGunKelly"),
			wantMethod: http.MethodDelete,
			wantPath:   "/aliases/MachineGunKelly",
		},
		{
			name: "create tenants",
			req: &api.CreateTenantsRequest{
				Collection: "Songs",
				Tenants: []api.Tenant{
					{Name: "john_doe", Status: api.TenantStatusActive},
					{Name: "jane_doe", Status: api.TenantStatusFrozen},
				},
			},
			wantMethod: http.MethodPost,
			wantPath:   "/schema/Songs/tenants",
			wantBody: rest.TenantsCreateJSONRequestBody{
				{Name: "john_doe", ActivityStatus: rest.ACTIVE},
				{Name: "jane_doe", ActivityStatus: rest.FROZEN},
			},
		},
		{
			name: "update tenants",
			req: &api.UpdateTenantsRequest{
				Collection: "Songs",
				Tenants: []api.Tenant{
					{Name: "john_doe", Status: api.TenantStatusActive},
					{Name: "jane_doe", Status: api.TenantStatusFrozen},
				},
			},
			wantMethod: http.MethodPut,
			wantPath:   "/schema/Songs/tenants",
			wantBody: rest.TenantsUpdateJSONRequestBody{
				{Name: "john_doe", ActivityStatus: rest.ACTIVE},
				{Name: "jane_doe", ActivityStatus: rest.FROZEN},
			},
		},
		{
			name: "delete tenants",
			req: &api.DeleteTenantsRequest{
				Collection: "Songs",
				Tenants:    []string{"john_doe", "jane_doe"},
			},
			wantMethod: http.MethodDelete,
			wantPath:   "/schema/Songs/tenants",
			wantBody: rest.TenantsDeleteJSONRequestBody{
				"john_doe", "jane_doe",
			},
		},
		{
			name:       "get all shards",
			req:        &api.GetShardsRequest{Collection: "Songs"},
			wantMethod: http.MethodGet,
			wantPath:   "/replication/sharding-state",
			wantQuery:  url.Values{"collection": {"Songs"}},
		},
		{
			name: "get one shard",
			req: &api.GetShardsRequest{
				Collection: "Songs",
				Shard:      "xyz",
			},
			wantMethod: http.MethodGet,
			wantPath:   "/replication/sharding-state",
			wantQuery: url.Values{
				"collection": {"Songs"},
				"shard":      {"xyz"},
			},
		},
		{
			name: "get all nodes in the cluster",
			req: &api.GetNodesRequest{
				Verbosity: api.NodeVerbosityVerbose,
			},
			wantMethod: http.MethodGet,
			wantPath:   "/nodes",
			wantQuery:  url.Values{"output": {"verbose"}},
		},
		{
			name: "get all nodes in a collection",
			req: &api.GetNodesRequest{
				Collection: "Songs",
			},
			wantMethod: http.MethodGet,
			wantPath:   "/nodes/Songs",
		},
		{
			name: "get all nodes for shard",
			req: &api.GetNodesRequest{
				Shard: "xyz",
			},
			wantMethod: http.MethodGet,
			wantPath:   "/nodes",
			wantQuery:  url.Values{"shardName": {"xyz"}},
		},
		{
			name: "create replication",
			req: &api.CreateReplicationRequest{
				Type:       api.ReplicationCopy,
				Collection: "Songs",
				Shard:      "abc",
				Source:     "node-0",
				Target:     "node-2",
			},
			wantMethod: http.MethodPost,
			wantPath:   "/replication/replicate",
			wantBody: rest.ReplicateJSONRequestBody{
				Type:       rest.ReplicationReplicateReplicaRequestTypeCOPY,
				Collection: "Songs",
				Shard:      "abc",
				SourceNode: "node-0",
				TargetNode: "node-2",
			},
		},
		{
			name:       "get replication",
			req:        &api.GetReplicationRequest{ID: testkit.UUID},
			wantMethod: http.MethodGet,
			wantPath:   "/replication/replicate/" + testkit.UUID.String(),
			wantQuery:  url.Values{"includeHistory": {"false"}},
		},
		{
			name: "list replications",
			req: &api.ListReplicationsRequest{
				Collection:     "Songs",
				Target:         "node-1",
				IncludeHistory: true,
			},
			wantMethod: http.MethodGet,
			wantPath:   "/replication/replicate/list",
			wantQuery: url.Values{
				"includeHistory": {"true"},
				"collection":     {"Songs"},
				"targetNode":     {"node-1"},
				// "shard" is not included because it's empty
			},
		},
		{
			name:       "cancel replication",
			req:        api.CancelReplicationRequest(testkit.UUID),
			wantMethod: http.MethodPost,
			wantPath:   "/replication/replicate/" + testkit.UUID.String() + "/cancel",
		},
		{
			name:       "delete a replication",
			req:        api.DeleteReplicationRequest(testkit.UUID),
			wantMethod: http.MethodDelete,
			wantPath:   "/replication/replicate/" + testkit.UUID.String(),
		},
		{
			name:       "delete all replications",
			req:        api.DeleteAllReplicationsRequest,
			wantMethod: http.MethodDelete,
			wantPath:   "/replication/replicate",
		},
	}) {
		t.Run(tt.name, func(t *testing.T) {
			require.Implements(t, (*transports.Endpoint)(nil), tt.req)
			endpoint := (tt.req).(transports.Endpoint)

			assert.Equal(t, tt.wantMethod, endpoint.Method(), "bad method")
			assert.Equal(t, tt.wantPath, endpoint.Path(), "bad path")
			assert.Equal(t, tt.wantQuery, endpoint.Query(), "bad query")

			gotBody := endpoint.Body()
			gotJSON, err := json.Marshal(gotBody)
			require.NoError(t, err, "marshal request body")

			wantJSON, err := json.Marshal(tt.wantBody)
			require.NoError(t, err, "marshal wantBody")

			assert.JSONEq(t, string(wantJSON), string(gotJSON), "bad body")
		})
	}
}

// TestRESTResponses verifies that response objects in the 'api' package
// unmarshal response JSONs correctly.
func TestRESTResponses(t *testing.T) {
	for _, tt := range testkit.WithOnly(t, []struct {
		testkit.Only

		name string
		body any // Response body.
		dest any // Set dest to a pointer to the response struct.
		want any // Expected value after deserialization.
	}{
		{
			name: "collection config",
			body: &rest.Class{
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
				VectorConfig: map[string]rest.VectorConfig{
					"lyrics_vec": {
						VectorIndexType: "hfresh",
						VectorIndexConfig: map[string]any{
							"distance": "cosine",
							"rq": map[string]any{
								"enabled":       true,
								"cache":         true,
								"bits":          1,
								"rescore_limit": 16,
							},
						},
						Vectorizer: map[string]any{
							testkit.ModuleName: map[string]any{
								"url": "example.com",
							},
						},
					},
					"title_vec": {},
					"audio_vec": {
						VectorIndexType: "hfresh",
						VectorIndexConfig: map[string]any{
							"distance":                "manhattan",
							"skipDefaultQuantization": true,
						},
					},
				},
				ShardingConfig: map[string]any{
					"desiredCount":        3,
					"desiredVirtualCount": 150,
					"virtualPerPhysical":  50,
				},
				ReplicationConfig: rest.ReplicationConfig{
					Factor:           6,
					DeletionStrategy: rest.TimeBasedResolution,
					AsyncConfig: rest.ReplicationAsyncConfig{
						DiffBatchSize:             1,
						DiffPerNodeTimeout:        2,
						Frequency:                 3,
						FrequencyWhilePropagating: 4,
						PrePropagationTimeout:     5,
						PropagationConcurrency:    6,
						PropagationBatchSize:      7,
						PropagationLimit:          8,
						PropagationTimeout:        9,
						PropagationDelay:          10,
						HashtreeHeight:            11,
						LoggingFrequency:          12,
					},
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
			dest: new(api.Collection),
			want: &api.Collection{
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
				Vectors: map[string]api.VectorConfig{
					"lyrics_vec": {
						Index: &api.Module{
							Name: "hfresh",
							Conf: map[string]any{
								"distance": "cosine",
							},
						},
						Compression: &api.Module{
							Name: "rq",
							Conf: map[string]any{
								"enabled":       true,
								"cache":         true,
								"bits":          float64(1),
								"rescore_limit": float64(16),
							},
						},
						Vectorizer: &api.Module{
							Name: testkit.ModuleName,
							Conf: map[string]any{
								"url": "example.com",
							},
						},
					},
					"title_vec": {},
					"audio_vec": {
						Index: &api.Module{
							Name: "hfresh",
							Conf: map[string]any{
								"distance": "manhattan",
							},
						},
						SkipDefaultCompression: true,
					},
				},
				Sharding: &api.ShardingConfig{
					DesiredCount:        3,
					DesiredVirtualCount: 150,
					VirtualPerPhysical:  50,
				},
				Replication: &api.ReplicationConfig{
					Factor:           6,
					DeletionStrategy: api.TimeBasedResolution,
					AsyncReplication: &api.AsyncReplicationConfig{
						DiffBatchSize:                   1,
						DiffPerNodeTimeout:              2 * time.Second,
						ReplicationFrequency:            3 * time.Millisecond,
						ReplicationFrequencyPropagating: 4 * time.Millisecond,
						PrePropagationTimeout:           5 * time.Second,
						PropagationConcurrency:          6,
						PropagationBatchSize:            7,
						PropagationLimit:                8,
						PropagationTimeout:              9 * time.Second,
						PropagationDelay:                10 * time.Millisecond,
						HashTreeHeight:                  11,
						LoggingFrequency:                12 * time.Second,
					},
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
		{
			name: "backup create response",
			body: &rest.BackupCreateResponse{
				Backend: "filesystem",
				Id:      "bak-1",
				Bucket:  "my-backups",
				Path:    "/path/to/backup",
				Classes: []string{"Songs"},
				Error:   "whaam!",
				Status:  rest.BackupCreateResponseStatusFAILED,
			},
			dest: new(api.BackupInfo),
			want: &api.BackupInfo{
				Backend:             "filesystem",
				ID:                  "bak-1",
				Bucket:              "my-backups",
				Path:                "/path/to/backup",
				IncludesCollections: []string{"Songs"},
				Error:               "whaam!",
				Status:              api.BackupStatusFailed,
			},
		},
		{
			name: "backup restore response",
			body: &rest.BackupRestoreResponse{
				Backend: "filesystem",
				Id:      "bak-1",
				Path:    "/path/to/backup",
				Classes: []string{"Songs"},
				Error:   "whaam!",
				Status:  rest.BackupRestoreResponseStatusFAILED,
			},
			dest: new(api.BackupInfo),
			want: &api.BackupInfo{
				Backend:             "filesystem",
				ID:                  "bak-1",
				Path:                "/path/to/backup",
				IncludesCollections: []string{"Songs"},
				Error:               "whaam!",
				Status:              api.BackupStatusFailed,
			},
		},
		{
			name: "backup create status response",
			body: &rest.BackupCreateStatusResponse{
				Backend: "filesystem",
				Id:      "bak-1",
				Path:    "/path/to/backup",
				Status:  rest.BackupCreateStatusResponseStatusSUCCESS,
				Size:    92,
			},
			dest: new(api.BackupInfo),
			want: &api.BackupInfo{
				Backend:     "filesystem",
				ID:          "bak-1",
				Path:        "/path/to/backup",
				Status:      api.BackupStatusSuccess,
				SizeGiB:     testkit.Ptr[float32](92),
				StartedAt:   testkit.Ptr(time.Time{}),
				CompletedAt: testkit.Ptr(time.Time{}),
			},
		},
		{
			name: "backup list response",
			body: rest.BackupListResponse{
				{
					Id:      "bak-1",
					Classes: []string{"Artists"},
					Status:  rest.BackupListResponseStatusTRANSFERRING,
					Size:    92,
				},
				{
					Id:      "bak-2",
					Classes: []string{"Songs"},
					Status:  rest.BackupListResponseStatusTRANSFERRED,
					Size:    80085,
				},
			},
			dest: new([]api.BackupInfo),
			want: &[]api.BackupInfo{
				{
					ID:                  "bak-1",
					IncludesCollections: []string{"Artists"},
					Status:              api.BackupStatusTransferring,
					SizeGiB:             testkit.Ptr[float32](92),
					StartedAt:           testkit.Ptr(time.Time{}),
					CompletedAt:         testkit.Ptr(time.Time{}),
				},
				{
					ID:                  "bak-2",
					IncludesCollections: []string{"Songs"},
					Status:              api.BackupStatusTransferred,
					SizeGiB:             testkit.Ptr[float32](80085),
					StartedAt:           testkit.Ptr(time.Time{}),
					CompletedAt:         testkit.Ptr(time.Time{}),
				},
			},
		},
		{
			name: "instance metadata",
			body: &rest.Meta{
				Hostname: "example.com",
				Version:  "v1.37.0",
				Modules: map[string]any{
					"text2vec-weaviate": true,
					"backup-s3":         true,
				},
				GrpcMaxMessageSize: 4096,
			},
			dest: new(api.GetInstanceMetadataResponse),
			want: &api.GetInstanceMetadataResponse{
				Hostname: "example.com",
				Version:  "v1.37.0",
				Modules: map[string]any{
					"text2vec-weaviate": true,
					"backup-s3":         true,
				},
				GRPCMaxMessageSize: 4096,
			},
		},
		{
			name: "get role response",
			body: map[string]any{
				"name": "rock-n-role",
				"permissions": []map[string]any{
					{
						"action": "create_aliases",
						"aliases": map[string]any{
							"collection": "Songs",
							"alias":      "Records",
						},
					},
					{
						"action": "read_aliases",
						"aliases": map[string]any{
							"collection": "Songs",
							"alias":      "Records",
						},
					},
					{
						"action": "update_aliases",
						"aliases": map[string]any{
							"collection": "Artists",
							"alias":      "Musicians",
						},
					},
					{
						"action": "delete_aliases",
						"aliases": map[string]any{
							"collection": "Artists",
							"alias":      "Musicians",
						},
					},
					{
						"action": "manage_backups",
						"backups": map[string]any{
							"collection": "Songs",
						},
					},
					{
						"action": "read_cluster",
					},
					{
						"action": "create_collections",
						"collections": map[string]any{
							"collection": "Songs",
						},
					},
					{
						"action": "read_collections",
						"collections": map[string]any{
							"collection": "Songs",
						},
					},
					{
						"action": "update_collections",
						"collections": map[string]any{
							"collection": "Artists",
						},
					},
					{
						"action": "delete_collections",
						"collections": map[string]any{
							"collection": "Artists",
						},
					},
					{
						"action": "create_data",
						"data": map[string]any{
							"collection": "Songs",
							"tenant":     "john_doe",
						},
					},
					{
						"action": "read_data",
						"data": map[string]any{
							"collection": "Songs",
							"tenant":     "john_doe",
						},
					},
					{
						"action": "update_data",
						"data": map[string]any{
							"collection": "Artists",
						},
					},
					{
						"action": "delete_data",
						"data": map[string]any{
							"collection": "Artists",
						},
					},
					{
						"action": "read_groups",
						"groups": map[string]any{
							"group":     "external",
							"groupType": "oidc",
						},
					},
					{
						"action": "assign_and_revoke_groups",
						"groups": map[string]any{
							"group":     "external",
							"groupType": "oidc",
						},
					},
					{
						"action": "manage_namespaces",
						"namespaces": map[string]any{
							"namespace": "sandbox",
						},
					},
					{
						"action": "read_nodes",
						"nodes": map[string]any{
							"collection": "Songs",
							"verbosity":  api.NodeVerbosityMinimal,
						},
					},
					{
						"action": "read_nodes",
						"nodes": map[string]any{
							"collection": "Artists",
							"verbosity":  api.NodeVerbosityVerbose,
						},
					},
					{
						"action": "create_mcp",
					},
					{
						"action": "read_mcp",
					},
					{
						"action": "update_mcp",
					},
					{
						"action": "create_replicate",
						"replicate": map[string]any{
							"collection": "Songs",
							"shard":      "abc",
						},
					},
					{
						"action": "read_replicate",
						"replicate": map[string]any{
							"collection": "Songs",
							"shard":      "abc",
						},
					},
					{
						"action": "update_replicate",
						"replicate": map[string]any{
							"collection": "Songs",
							"shard":      "xyz",
						},
					},
					{
						"action": "delete_replicate",
						"replicate": map[string]any{
							"collection": "Songs",
							"shard":      "xyz",
						},
					},
					{
						"action": "create_roles",
						"roles": map[string]any{
							"role":  "rock-n-role",
							"scope": api.RoleScopeAll,
						},
					},
					{
						"action": "read_roles",
						"roles": map[string]any{
							"role":  "rock-n-role",
							"scope": api.RoleScopeAll,
						},
					},
					{
						"action": "update_roles",
						"roles": map[string]any{
							"role":  "rock-n-role",
							"scope": api.RoleScopeMatch,
						},
					},
					{
						"action": "delete_roles",
						"roles": map[string]any{
							"role":  "rock-n-role",
							"scope": api.RoleScopeMatch,
						},
					},
					{
						"action": "create_tenants",
						"tenants": map[string]any{
							"collection": "Songs",
							"tenant":     "john_doe",
						},
					},
					{
						"action": "read_tenants",
						"tenants": map[string]any{
							"collection": "Songs",
							"tenant":     "john_doe",
						},
					},
					{
						"action": "update_tenants",
						"tenants": map[string]any{
							"collection": "Artists",
							"tenant":     "*",
						},
					},
					{
						"action": "delete_tenants",
						"tenants": map[string]any{
							"collection": "Artists",
							"tenant":     "*",
						},
					},
					{
						"action": "create_users",
						"users": map[string]any{
							"user": "john_doe",
						},
					},
					{
						"action": "read_users",
						"users": map[string]any{
							"user": "john_doe",
						},
					},
					{
						"action": "update_users",
						"users": map[string]any{
							"user": "jane_doe",
						},
					},
					{
						"action": "delete_users",
						"users": map[string]any{
							"user": "jane_doe",
						},
					},
					{
						"action": "assign_and_revoke_users",
						"users": map[string]any{
							"user": "jim_beam",
						},
					},
				},
			},
			dest: new(api.Role),
			want: &api.Role{
				ID: "rock-n-role",
				Permissions: api.Permissions{
					Aliases: []api.AliasPermission{
						{
							Collection: "Songs", Alias: "Records",
							Create: true, Read: true,
						},
						{
							Collection: "Artists", Alias: "Musicians",
							Update: true, Delete: true,
						},
					},
					Backups: []api.BackupsPermission{
						{Collection: "Songs", Manage: true},
					},
					Cluster: []api.ClusterPermission{{Read: true}},
					Collections: []api.CollectionPermission{
						{Collection: "Songs", Create: true, Read: true},
						{Collection: "Artists", Update: true, Delete: true},
					},
					Data: []api.DataPermission{
						{
							Collection: "Songs", Tenant: "john_doe",
							Create: true, Read: true,
						},
						{
							Collection: "Artists",
							Update:     true, Delete: true,
						},
					},
					Groups: []api.GroupPermission{
						{
							GroupID: "external", Type: "oidc",
							Read: true, AssignAndRevoke: true,
						},
					},
					Namespaces: []api.NamespacePermission{
						{Namespace: "sandbox", Manage: true},
					},
					Nodes: []api.NodesPermission{
						{Collection: "Songs", Verbosity: api.NodeVerbosityMinimal, Read: true},
						{Collection: "Artists", Verbosity: api.NodeVerbosityVerbose, Read: true},
					},
					MCP: []api.MCPPermission{
						{Create: true, Read: true, Update: true},
					},
					Replication: []api.ReplicationPermission{
						{
							Collection: "Songs", Shard: "abc",
							Create: true, Read: true,
						},
						{
							Collection: "Songs", Shard: "xyz",
							Update: true, Delete: true,
						},
					},
					Roles: []api.RolePermission{
						{
							RoleID: "rock-n-role", Scope: api.RoleScopeAll,
							Create: true, Read: true,
						},
						{
							RoleID: "rock-n-role", Scope: api.RoleScopeMatch,
							Update: true, Delete: true,
						},
					},
					Tenants: []api.TenantPermission{
						{
							Collection: "Songs", Tenant: "john_doe",
							Create: true, Read: true,
						},
						{
							Collection: "Artists", Tenant: "*",
							Update: true, Delete: true,
						},
					},
					Users: []api.UserPermission{
						{UserID: "john_doe", Create: true, Read: true},
						{UserID: "jane_doe", Update: true, Delete: true},
						{UserID: "jim_beam", AssignAndRevoke: true},
					},
				},
			},
		},
		{
			name: "role has permission",
			body: testkit.Ptr(true),
			dest: new(api.HasPermissionResponse),
			want: testkit.Ptr[api.HasPermissionResponse](true),
		},
		{
			name: "user assignments",
			body: []map[string]any{{
				"userId":   "john_doe",
				"userType": rest.DBUserInfoDbUserTypeDbEnvUser,
			}},
			dest: new([]api.UserInfo),
			want: &[]api.UserInfo{{
				ID:   "john_doe",
				Type: api.UserTypeDBEnv,
			}},
		},
		{
			name: "group assignment",
			body: map[string]any{
				"groupId":   "external",
				"groupType": "oidc",
			},
			dest: new(api.GroupInfo),
			want: &api.GroupInfo{
				ID:   "external",
				Type: "oidc",
			},
		},
		{
			name: "own user info",
			body: &rest.UserOwnInfo{
				Username: "john-malkovich",
				Roles: []rest.Role{
					{Name: "rock-n-role"},
					{Name: "sushi-role"},
				},
				Groups: []string{"external"},
			},
			dest: new(api.GetOwnUserInfoResponse),
			want: &api.GetOwnUserInfoResponse{
				ID: "john-malkovich",
				Roles: []api.Role{
					{ID: "rock-n-role"},
					{ID: "sushi-role"},
				},
				Groups: []string{"external"},
			},
		},
		{
			name: "create db user",
			body: &rest.UserApiKey{
				Apikey: "abracadabra",
			},
			dest: new(api.CreateUserResponse),
			want: &api.CreateUserResponse{
				APIKey: "abracadabra",
			},
		},
		{
			name: "rotate db user api key",
			body: &rest.UserApiKey{
				Apikey: "abracadabra",
			},
			dest: new(api.RotateAPIKeyResponse),
			want: &api.RotateAPIKeyResponse{
				APIKey: "abracadabra",
			},
		},
		{
			name: "get user info",
			body: &rest.DBUserInfo{
				UserId:     "john-malkovich",
				DbUserType: rest.DBUserInfoDbUserTypeDbUser,
				Active:     true,
				Roles:      []string{"external"},
			},
			dest: new(api.UserInfo),
			want: &api.UserInfo{
				ID:     "john-malkovich",
				Type:   "db_user",
				Active: true,
				Roles:  []string{"external"},
			},
		},
		{
			name: "list db users",
			body: []rest.DBUserInfo{{
				UserId:     "john-malkovich",
				DbUserType: rest.DBUserInfoDbUserTypeDbUser,
				Active:     true,
				Roles:      []string{"external"},
			}},
			dest: new([]api.UserInfo),
			want: &[]api.UserInfo{{
				ID:     "john-malkovich",
				Type:   "db_user",
				Active: true,
				Roles:  []string{"external"},
			}},
		},
		{
			name: "list groups",
			body: []string{"internal", "external"},
			dest: new(api.ListGroupsResponse),
			want: &api.ListGroupsResponse{"internal", "external"},
		},
		{
			name: "get alias",
			body: &rest.Alias{
				Class: "GeorgeBarnes",
				Alias: "MachineGunKelly",
			},
			dest: new(api.Alias),
			want: &api.Alias{
				Collection: "GeorgeBarnes",
				Alias:      "MachineGunKelly",
			},
		},
		{
			name: "list aliases",
			body: &rest.AliasResponse{
				Aliases: []rest.Alias{
					{Class: "GeorgeBarnes", Alias: "MachineGunKelly"},
					{Class: "LouisAttanasio", Alias: "LouieHaHa"},
				},
			},
			dest: new(api.ListAliasesResponse),
			want: &api.ListAliasesResponse{
				{Collection: "GeorgeBarnes", Alias: "MachineGunKelly"},
				{Collection: "LouisAttanasio", Alias: "LouieHaHa"},
			},
		},
		{
			name: "shard replicas",
			body: &rest.ReplicationShardingStateResponse{
				ShardingState: rest.ReplicationShardingState{
					Shards: []rest.ReplicationShardReplicas{
						{Shard: "abc", Replicas: []string{"abc-1", "abc-2"}},
						{Shard: "xyz", Replicas: []string{"xyz-1", "xyz-2"}},
					},
				},
			},
			dest: new(api.GetShardsResponse),
			want: &api.GetShardsResponse{
				{Shard: "abc", Replicas: []string{"abc-1", "abc-2"}},
				{Shard: "xyz", Replicas: []string{"xyz-1", "xyz-2"}},
			},
		},
		{
			name: "nodes",
			body: &rest.NodesStatusResponse{
				Nodes: []rest.NodeStatus{
					{
						Name:            "node-1",
						Status:          rest.NodeStatusStatusHEALTHY,
						GitHash:         "5cc3aa3",
						Version:         "1.37.3",
						OperationalMode: rest.ScaleOut,
						Stats: rest.NodeStats{
							ObjectCount: 4096,
							ShardCount:  8,
						},
						Shards: []rest.NodeShardStatus{
							{
								Name:                 "abc",
								Class:                "Song",
								ObjectCount:          256,
								Compressed:           true,
								Loaded:               true,
								VectorIndexingStatus: "ok",
								VectorQueueLength:    16,
								AsyncReplicationStatus: []rest.AsyncReplicationStatus{
									{
										TargetNode:              "node-2",
										ObjectsPropagated:       1024,
										StartDiffTimeUnixMillis: 3600,
									},
								},
							},
						},
						BatchStats: rest.BatchStats{
							QueueLength:   512,
							RatePerSecond: 64,
						},
					},
				},
			},
			dest: new(api.GetNodesResponse),
			want: &api.GetNodesResponse{
				{
					Name:    "node-1",
					Status:  api.NodeStatusHealthy,
					GitHash: "5cc3aa3",
					Version: "1.37.3",
					Mode:    api.NodeModeScaleOut,
					Stats: api.NodeStats{
						ObjectCount: 4096,
						ShardCount:  8,
					},
					Shards: []api.Shard{
						{
							Name:                 "abc",
							Collection:           "Song",
							ObjectCount:          256,
							Compressed:           true,
							Loaded:               true,
							VectorIndexingStatus: "ok",
							VectorQueueLength:    16,
							OngoingReplications: []api.ReplicationStatus{
								{
									TargetNode:         "node-2",
									ObjectsPropagated:  1024,
									LastIterationStart: time.UnixMilli(3600),
								},
							},
						},
					},
					BatchStats: api.BatchStats{
						QueueLength:   512,
						RatePerSecond: 64,
					},
				},
			},
		},
		{
			name: "replication info",
			body: &rest.ReplicationReplicateDetailsReplicaResponse{
				Type:               rest.ReplicationReplicateDetailsReplicaResponseTypeMOVE,
				Id:                 &testkit.UUID,
				Collection:         "Songs",
				Shard:              "abc",
				SourceNode:         "node-0",
				TargetNode:         "node-1",
				Uncancelable:       true,
				ScheduledForCancel: true,
				ScheduledForDelete: true,
				WhenStartedUnixMs:  testkit.Now.UnixMilli(),
				Status: rest.ReplicationReplicateDetailsReplicaStatus{
					State:             rest.FINALIZING,
					WhenStartedUnixMs: testkit.Now.UnixMilli(),
				},
				StatusHistory: []rest.ReplicationReplicateDetailsReplicaStatus{
					{
						State: rest.INTEGRATING,
						Errors: []rest.ReplicationReplicateDetailsReplicaStatusError{
							{Message: "Whaam!", WhenErroredUnixMs: testkit.Now.UnixMilli()},
						},
						WhenStartedUnixMs: testkit.Now.UnixMilli(),
					},
				},
			},
			dest: new(api.Replication),
			want: &api.Replication{
				Type:            api.ReplicationMove,
				ID:              testkit.UUID,
				Collection:      "Songs",
				Shard:           "abc",
				Source:          "node-0",
				Target:          "node-1",
				CanCancel:       false,
				CancelScheduled: true,
				DeleteScheduled: true,
				StartedAt:       testkit.Now,
				Current: api.ReplicationStage{
					State:     api.ReplicationStateFinalizing,
					StartedAt: testkit.Now,
				},
				History: []api.ReplicationStage{
					{
						State:     api.ReplicationStateIntegrating,
						StartedAt: testkit.Now,
						Errors: []api.ReplicationError{
							{Message: "Whaam!", Time: testkit.Now},
						},
					},
				},
			},
		},
	}) {
		t.Run(tt.name, func(t *testing.T) {
			require.NotNil(t, tt.body, "incomplete test case: body is nil")
			testkit.RequirePointer(t, tt.body, "body")
			testkit.RequirePointer(t, tt.dest, "dest")

			body, err := json.Marshal(tt.body)
			require.NoError(t, err, "marshal expected body")

			err = json.Unmarshal(body, tt.dest)
			assert.NoError(t, err, "unmarshal response body")
			assert.Equal(t, tt.want, tt.dest, "bad unmarshaled value")
		})
	}
}
