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
	"github.com/weaviate/weaviate-go-client/v6/internal/transports"
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
			name: "insert object (consistency_level=ONE)",
			req: &api.InsertObjectRequest{
				RequestDefaults: api.RequestDefaults{
					CollectionName:   "Songs",
					Tenant:           "john_doe",
					ConsistencyLevel: api.ConsistencyLevelOne,
				},
				UUID: &uuid.Nil,
				Properties: map[string]any{
					"title":  "High Speed Dirt",
					"genres": []string{"thrash metal", "blues"},
					"single": false,
					"year":   1992,
				},
				Vectors: []api.Vector{
					{Name: "lyrics", Single: []float32{1, 2, 3}},
				},
			},
			wantMethod: http.MethodPost,
			wantPath:   "/objects",
			wantQuery:  url.Values{"consistency_level": {string(api.ConsistencyLevelOne)}},
			wantBody: &rest.Object{
				Class:  "Songs",
				Tenant: "john_doe",
				Id:     &uuid.Nil,
				Properties: map[string]any{
					"title":  "High Speed Dirt",
					"genres": []string{"thrash metal", "blues"},
					"single": false,
					"year":   1992,
				},
				Vectors: map[string]any{
					"lyrics": []float32{1, 2, 3},
				},
			},
		},
		{
			name: "insert object (no consistency_level)",
			req: &api.InsertObjectRequest{
				RequestDefaults: api.RequestDefaults{CollectionName: "Songs"},
			},
			wantMethod: http.MethodPost,
			wantPath:   "/objects",
			wantBody:   &rest.Object{Class: "Songs"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			require.Implements(t, (*transports.Endpoint)(nil), tt.req)
			endpoint := (tt.req).(transports.Endpoint)

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
