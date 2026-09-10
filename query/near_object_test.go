package query_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/query"
	"github.com/weaviate/weaviate-go-client/v6/types"
)

func TestNearObject(t *testing.T) {
	rd := api.RequestDefaults{
		CollectionName:   "Songs",
		Tenant:           "john_doe",
		ConsistencyLevel: api.ConsistencyLevelQuorum,
	}

	for _, tt := range testkit.WithOnly(t, []struct {
		testkit.Only

		name  string
		no    query.NearObject
		stubs []testkit.Stub[api.SearchRequest, api.SearchResponse]
		want  *query.Result // Expected return value.
		err   testkit.Error
	}{
		{
			name: "include all",
			no: query.NearObject{
				UUID: testkit.UUID,
			},
			stubs: []testkit.Stub[api.SearchRequest, api.SearchResponse]{
				{
					Request: &api.SearchRequest{
						RequestDefaults: rd,
						NearObject: &api.NearObject{
							UUID: testkit.UUID,
						},
					},
					Response: api.SearchResponse{
						Took: 92 * time.Second,
						Results: []api.Object{
							{
								Collection: "Songs",
								Metadata: api.ObjectMetadata{
									UUID: testkit.UUID,
								},
							},
						},
					},
				},
			},
			want: &query.Result{
				Took: 92 * time.Second,
				Objects: []query.Object[map[string]any]{
					{
						Object: types.Object[map[string]any]{
							Collection: "Songs",
							UUID:       testkit.UUID,
						},
					},
				},
			},
		},
		{
			name: "exclude self",
			no: query.NearObject{
				UUID:        testkit.UUID,
				ExcludeSelf: true,
			},
			stubs: []testkit.Stub[api.SearchRequest, api.SearchResponse]{
				{
					Request: &api.SearchRequest{
						RequestDefaults: rd,
						NearObject: &api.NearObject{
							UUID: testkit.UUID,
						},
						Filter: api.FilterExpr{
							Operator: api.FilterOperatorNot,
							Exprs: []api.FilterExpr{
								{
									Target:   []string{api.FieldUUID},
									Operator: api.FilterOperatorEqual,
									Value:    testkit.UUID,
								},
							},
						},
					},
					Response: api.SearchResponse{
						Took: 92 * time.Second,
					},
				},
			},
			want: &query.Result{Took: 92 * time.Second},
		},
		{
			name: "request error",
			stubs: []testkit.Stub[api.SearchRequest, api.SearchResponse]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	}) {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)

			c := query.NewClient(transport, rd)
			require.NotNil(t, c, "client")

			got, err := c.NearObject(t.Context(), tt.no)
			tt.err.Require(t, err, "near object query")
			require.EqualExportedValues(t, tt.want, got, "query result")
		})
	}
}
