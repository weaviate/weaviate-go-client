package api_test

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

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
	for _, tt := range []struct {
		name string
		req  any // Request object.

		wantMethod string     // Expected HTTP Method.
		wantPath   string     // Expected endpoint + path parameters.
		wantQuery  url.Values // Expected query parameters.
		wantBody   any        // Expected request body. JSON strings are compared.
	}{
		{
			name: "insert object (no consistency_level)",
			req: &api.InsertObjectRequest{
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
			wantMethod: http.MethodPost,
			wantPath:   "/objects",
			wantBody: &rest.Object{
				Class:  "Songs",
				Tenant: "john_doe",
				Id:     &testkit.UUID,
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
			name: "insert object (consistency_level=ONE)",
			req: &api.InsertObjectRequest{
				RequestDefaults: api.RequestDefaults{
					CollectionName:   "Songs",
					ConsistencyLevel: api.ConsistencyLevelOne,
				},
			},
			wantMethod: http.MethodPost,
			wantPath:   "/objects",
			wantQuery:  url.Values{"consistency_level": {string(api.ConsistencyLevelOne)}},
			wantBody:   &rest.Object{Class: "Songs"},
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
			name: "inserted object",
			body: &rest.Object{
				Class:              "Songs",
				Tenant:             "john_doe",
				Id:                 &testkit.UUID,
				CreationTimeUnix:   testkit.Now.UnixMilli(),
				LastUpdateTimeUnix: testkit.Now.UnixMilli(),
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
			dest: new(api.InsertObjectResponse),
			want: &api.InsertObjectResponse{
				UUID:          testkit.UUID,
				CreatedAt:     testkit.Now,
				LastUpdatedAt: testkit.Now,
				Properties: map[string]any{
					"title":  "High Speed Dirt",
					"genres": []any{"thrash metal", "blues"},
					"single": false,
					"year":   float64(1992), // json.Marshal treats numbers as float64 by default
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
				Vectors: map[string]api.Vector{
					"lyrics": {Name: "lyrics", Single: []float32{1, 2, 3}},
				},
			},
		},
		{
			name: "replaced object",
			body: &rest.Object{
				Class:              "Songs",
				Tenant:             "john_doe",
				Id:                 &testkit.UUID,
				CreationTimeUnix:   testkit.Now.UnixMilli(),
				LastUpdateTimeUnix: testkit.Now.UnixMilli(),
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
			dest: new(api.ReplaceObjectResponse),
			want: &api.ReplaceObjectResponse{
				UUID:          testkit.UUID,
				CreatedAt:     testkit.Now,
				LastUpdatedAt: testkit.Now,
				Properties: map[string]any{
					"title":  "High Speed Dirt",
					"genres": []any{"thrash metal", "blues"},
					"single": false,
					"year":   float64(1992), // json.Marshal treats numbers as float64 by default
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
				Vectors: map[string]api.Vector{
					"lyrics": {Name: "lyrics", Single: []float32{1, 2, 3}},
				},
			},
		},
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
