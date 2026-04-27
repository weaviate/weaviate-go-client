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
				References: api.ObjectReferences{
					"band": {
						{UUID: testkit.UUID, Collection: "Drummers"},
						{UUID: testkit.UUID, Collection: "Basists"},
					},
					"label": {
						{UUID: testkit.UUID},
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
					Sharding: &api.ShardingConfig{
						DesiredCount:        3,
						DesiredVirtualCount: 150,
						VirtualPerPhysical:  50,
					},
					Replication: &api.ReplicationConfig{
						AsyncEnabled:     false,
						Factor:           6,
						DeletionStrategy: api.TimeBasedResolution,
						AsyncReplication: &api.AsyncReplicationConfig{
							DiffBatchSize:                   1,
							DiffPerNodeTimeout:              2 * time.Second,
							ReplicationConcurrency:          3,
							ReplicationFrequency:            4 * time.Millisecond,
							ReplicationFrequencyPropagating: 5 * time.Millisecond,
							PrePropagationTimeout:           6 * time.Second,
							PropagationConcurrency:          7,
							PropagationBatchSize:            8,
							PropagationLimit:                9,
							PropagationTimeout:              10 * time.Second,
							PropagationDelay:                11 * time.Millisecond,
							HashTreeHeight:                  12,
							NodePingFrequency:               13 * time.Millisecond,
							LoggingFrequency:                14 * time.Second,
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
				ShardingConfig: map[string]any{
					"desiredCount":        3,
					"desiredVirtualCount": 150,
					"virtualPerPhysical":  50,
				},
				ReplicationConfig: rest.ReplicationConfig{
					AsyncEnabled:     false,
					Factor:           6,
					DeletionStrategy: rest.TimeBasedResolution,
					AsyncConfig: rest.ReplicationAsyncConfig{
						DiffBatchSize:               1,
						DiffPerNodeTimeout:          2,
						MaxWorkers:                  3,
						Frequency:                   4,
						FrequencyWhilePropagating:   5,
						PrePropagationTimeout:       6,
						PropagationConcurrency:      7,
						PropagationBatchSize:        8,
						PropagationLimit:            9,
						PropagationTimeout:          10,
						PropagationDelay:            11,
						HashtreeHeight:              12,
						AliveNodesCheckingFrequency: 13,
						LoggingFrequency:            14,
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
	for _, tt := range []struct {
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
				ShardingConfig: map[string]any{
					"desiredCount":        3,
					"desiredVirtualCount": 150,
					"virtualPerPhysical":  50,
				},
				ReplicationConfig: rest.ReplicationConfig{
					AsyncEnabled:     false,
					Factor:           6,
					DeletionStrategy: rest.TimeBasedResolution,
					AsyncConfig: rest.ReplicationAsyncConfig{
						DiffBatchSize:               1,
						DiffPerNodeTimeout:          2,
						MaxWorkers:                  3,
						Frequency:                   4,
						FrequencyWhilePropagating:   5,
						PrePropagationTimeout:       6,
						PropagationConcurrency:      7,
						PropagationBatchSize:        8,
						PropagationLimit:            9,
						PropagationTimeout:          10,
						PropagationDelay:            11,
						HashtreeHeight:              12,
						AliveNodesCheckingFrequency: 13,
						LoggingFrequency:            14,
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
				Sharding: &api.ShardingConfig{
					DesiredCount:        3,
					DesiredVirtualCount: 150,
					VirtualPerPhysical:  50,
				},
				Replication: &api.ReplicationConfig{
					AsyncEnabled:     false,
					Factor:           6,
					DeletionStrategy: api.TimeBasedResolution,
					AsyncReplication: &api.AsyncReplicationConfig{
						DiffBatchSize:                   1,
						DiffPerNodeTimeout:              2 * time.Second,
						ReplicationConcurrency:          3,
						ReplicationFrequency:            4 * time.Millisecond,
						ReplicationFrequencyPropagating: 5 * time.Millisecond,
						PrePropagationTimeout:           6 * time.Second,
						PropagationConcurrency:          7,
						PropagationBatchSize:            8,
						PropagationLimit:                9,
						PropagationTimeout:              10 * time.Second,
						PropagationDelay:                11 * time.Millisecond,
						HashTreeHeight:                  12,
						NodePingFrequency:               13 * time.Millisecond,
						LoggingFrequency:                14 * time.Second,
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
	} {
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
