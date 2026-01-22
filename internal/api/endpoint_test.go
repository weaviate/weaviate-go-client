package api_test

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
	"github.com/google/uuid"
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
				UUID: &uuid.Nil,
				Properties: map[string]any{
					"title":  "High Speed Dirt",
					"genres": []string{"thrash metal", "blues"},
					"single": false,
					"year":   1992,
				},
				References: api.ObjectReferences{
					"band": {
						{UUID: uuid.Nil, Collection: "Drummers"},
						{UUID: uuid.Nil, Collection: "Basists"},
					},
					"label": {
						{UUID: uuid.Nil},
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
				Id:     &uuid.Nil,
				Properties: map[string]any{
					"title":  "High Speed Dirt",
					"genres": []string{"thrash metal", "blues"},
					"single": false,
					"year":   1992,
					"band": []string{
						"weaviate://localhost/Drummers/" + uuid.Nil.String(),
						"weaviate://localhost/Basists/" + uuid.Nil.String(),
					},
					"label": []string{
						"weaviate://localhost/" + uuid.Nil.String(),
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
				UUID: &uuid.Nil,
				Properties: map[string]any{
					"title":  "High Speed Dirt",
					"genres": []string{"thrash metal", "blues"},
					"single": false,
					"year":   1992,
				},
				References: api.ObjectReferences{
					"band": {
						{UUID: uuid.Nil, Collection: "Drummers"},
						{UUID: uuid.Nil, Collection: "Basists"},
					},
					"label": {
						{UUID: uuid.Nil},
					},
				},
				Vectors: []api.Vector{
					{Name: "lyrics", Single: []float32{1, 2, 3}},
				},
			},
			wantMethod: http.MethodPut,
			wantPath:   "/objects/Songs/" + uuid.Nil.String(),
			wantBody: &rest.Object{
				Tenant: "john_doe",
				Properties: map[string]any{
					"title":  "High Speed Dirt",
					"genres": []string{"thrash metal", "blues"},
					"single": false,
					"year":   1992,
					"band": []string{
						"weaviate://localhost/Drummers/" + uuid.Nil.String(),
						"weaviate://localhost/Basists/" + uuid.Nil.String(),
					},
					"label": []string{
						"weaviate://localhost/" + uuid.Nil.String(),
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
				UUID: &uuid.Nil,
			},
			wantMethod: http.MethodPut,
			wantPath:   "/objects/Songs/" + uuid.Nil.String(),
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
				UUID: uuid.Nil,
			},
			wantMethod: http.MethodDelete,
			wantPath:   "/objects/Songs/" + uuid.Nil.String(),
			wantQuery:  url.Values{"tenant": {"john_doe"}},
		},
		{
			name: "delete object (no tenant)",
			req: &api.DeleteObjectRequest{
				RequestDefaults: api.RequestDefaults{
					CollectionName:   "Songs",
					ConsistencyLevel: api.ConsistencyLevelOne,
				},
				UUID: uuid.Nil,
			},
			wantMethod: http.MethodDelete,
			wantPath:   "/objects/Songs/" + uuid.Nil.String(),
			wantQuery:  url.Values{"consistency_level": {string(api.ConsistencyLevelOne)}},
		},
		{
			name: "delete object (no tenant, no consistency_level)",
			req: &api.DeleteObjectRequest{
				RequestDefaults: api.RequestDefaults{CollectionName: "Songs"},
				UUID:            uuid.Nil,
			},
			wantMethod: http.MethodDelete,
			wantPath:   "/objects/Songs/" + uuid.Nil.String(),
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
				Id:                 &uuid.Nil,
				CreationTimeUnix:   testkit.Now.UnixMilli(),
				LastUpdateTimeUnix: testkit.Now.UnixMilli(),
				Properties: map[string]any{
					"title":  "High Speed Dirt",
					"genres": []string{"thrash metal", "blues"},
					"single": false,
					"year":   1992,
					"band": []string{
						"weaviate://localhost/Drummers/" + uuid.Nil.String(),
						"weaviate://localhost/Basists/" + uuid.Nil.String(),
					},
					"label": []string{
						"weaviate://localhost/" + uuid.Nil.String(),
					},
				},
				Vectors: map[string]any{
					"lyrics": []float32{1, 2, 3},
				},
			},
			dest: new(api.InsertObjectResponse),
			want: &api.InsertObjectResponse{
				UUID:          uuid.Nil,
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
						{UUID: uuid.Nil, Collection: "Drummers"},
						{UUID: uuid.Nil, Collection: "Basists"},
					},
					"label": {
						{UUID: uuid.Nil},
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
				Id:                 &uuid.Nil,
				CreationTimeUnix:   testkit.Now.UnixMilli(),
				LastUpdateTimeUnix: testkit.Now.UnixMilli(),
				Properties: map[string]any{
					"title":  "High Speed Dirt",
					"genres": []string{"thrash metal", "blues"},
					"single": false,
					"year":   1992,
					"band": []string{
						"weaviate://localhost/Drummers/" + uuid.Nil.String(),
						"weaviate://localhost/Basists/" + uuid.Nil.String(),
					},
					"label": []string{
						"weaviate://localhost/" + uuid.Nil.String(),
					},
				},
				Vectors: map[string]any{
					"lyrics": []float32{1, 2, 3},
				},
			},
			dest: new(api.ReplaceObjectResponse),
			want: &api.ReplaceObjectResponse{
				UUID:          uuid.Nil,
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
						{UUID: uuid.Nil, Collection: "Drummers"},
						{UUID: uuid.Nil, Collection: "Basists"},
					},
					"label": {
						{UUID: uuid.Nil},
					},
				},
				Vectors: map[string]api.Vector{
					"lyrics": {Name: "lyrics", Single: []float32{1, 2, 3}},
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			require.NotNil(t, tt.body, "incomplete test case: body is nil")
			testkit.IsPointer(t, tt.body, "body")
			testkit.IsPointer(t, tt.dest, "dest")

			body, err := json.Marshal(tt.body)
			require.NoError(t, err, "marshal expected body")

			err = json.Unmarshal(body, tt.dest)
			assert.NoError(t, err, "unmarshal response body")
			assert.Equal(t, tt.want, tt.dest, "bad unmarshaled value")
		})
	}
}
