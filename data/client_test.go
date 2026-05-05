package data_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/data"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/types"
)

func TestNewClient(t *testing.T) {
	require.Panics(t, func() {
		data.NewClient(nil, api.RequestDefaults{CollectionName: "New"})
	}, "nil transport")
}

func TestClient_Insert(t *testing.T) {
	rd := api.RequestDefaults{
		CollectionName:   "Insert",
		ConsistencyLevel: api.ConsistencyLevelOne,
		Tenant:           "john_doe",
	}

	for _, tt := range []struct {
		name    string
		objects []*data.Object // Object to be inserted.
		stubs   []testkit.Stub[api.InsertObjectsRequest, api.InsertObjectsResponse]
		want    *data.InsertResult
		err     testkit.Error // Expected error.
	}{
		{
			name: "nil object",
			stubs: []testkit.Stub[api.InsertObjectsRequest, api.InsertObjectsResponse]{{
				Request:  &api.InsertObjectsRequest{RequestDefaults: rd},
				Response: api.InsertObjectsResponse{Took: 92 * time.Second},
			}},
			want: &data.InsertResult{Took: 92 * time.Second},
		},
		{
			name:    "nil uuid",
			objects: []*data.Object{{UUID: testkit.Ptr(uuid.Nil)}},
			stubs: []testkit.Stub[api.InsertObjectsRequest, api.InsertObjectsResponse]{{
				Request: &api.InsertObjectsRequest{
					RequestDefaults: rd,
					Objects:         []api.BatchObject{{UUID: uuid.Nil}},
				},
				Response: api.InsertObjectsResponse{Took: 92 * time.Second},
			}},
			want: &data.InsertResult{Took: 92 * time.Second},
		},
		{
			name: "with data",
			objects: []*data.Object{{
				UUID: testkit.Ptr(testkit.UUID),
				Vectors: []types.Vector{
					{Name: "single", Single: []float32{1, 2, 3}},
				},
				Properties: map[string]any{"foo": "bar"},
				References: data.References{
					"ref": []data.Reference{
						{Collection: "Foo", UUID: testkit.UUID},
						{Collection: "Bar", UUID: testkit.UUID},
					},
				},
			}},
			stubs: []testkit.Stub[api.InsertObjectsRequest, api.InsertObjectsResponse]{{
				Request: &api.InsertObjectsRequest{
					RequestDefaults: rd,
					Objects: []api.BatchObject{{
						UUID: testkit.UUID,
						Vectors: []api.Vector{
							{Name: "single", Single: []float32{1, 2, 3}},
						},
						Properties: map[string]any{"foo": "bar"},
						References: api.ObjectReferences{
							"ref": []api.ObjectReference{
								{Collection: "Foo", UUID: testkit.UUID},
								{Collection: "Bar", UUID: testkit.UUID},
							},
						},
					}},
				},
				Response: api.InsertObjectsResponse{Took: 92 * time.Second},
			}},
			want: &data.InsertResult{Took: 92 * time.Second},
		},
		{
			name:    "internal server error",
			objects: []*data.Object{{UUID: testkit.Ptr(testkit.UUID)}},
			stubs: []testkit.Stub[api.InsertObjectsRequest, api.InsertObjectsResponse]{{
				Request: &api.InsertObjectsRequest{
					RequestDefaults: rd,
					Objects:         []api.BatchObject{{UUID: testkit.UUID}},
				},
				Response: api.InsertObjectsResponse{
					Took:      92 * time.Second,
					Positions: []int32{0},
					Errors:    []string{"Whaam!"},
				},
			}},
			want: &data.InsertResult{
				Took: 92 * time.Second,
				Errors: map[uuid.UUID]string{
					testkit.UUID: "Whaam!",
				},
			},
		},
		{
			name: "with request error",
			stubs: []testkit.Stub[api.InsertObjectsRequest, api.InsertObjectsResponse]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := data.NewClient(transport, rd)
			require.NotNil(t, c, "nil client")

			got, err := c.Insert(t.Context(), tt.objects...)
			tt.err.Require(t, err, "insert error")
			require.Equal(t, tt.want, got, "returned object")
		})
	}
}

func TestClient_Replace(t *testing.T) {
	rd := api.RequestDefaults{
		CollectionName:   "Replace",
		ConsistencyLevel: api.ConsistencyLevelOne,
		Tenant:           "john_doe",
	}

	for _, tt := range []struct {
		name   string
		object data.Object // Object to be replaced.
		stub   []testkit.Stub[api.ReplaceObjectRequest, any]
		err    testkit.Error
	}{
		{
			name: "with data",
			object: data.Object{
				UUID: &testkit.UUID,
				Vectors: []types.Vector{
					{Name: "single", Single: []float32{1, 2, 3}},
				},
				Properties: map[string]any{"foo": "bar"},
				References: data.References{
					"ref": []data.Reference{
						{Collection: "Foo", UUID: testkit.UUID},
						{Collection: "Bar", UUID: testkit.UUID},
					},
				},
			},
			stub: []testkit.Stub[api.ReplaceObjectRequest, any]{{
				Request: &api.ReplaceObjectRequest{
					RequestDefaults: rd,
					UUID:            &testkit.UUID,
					Vectors: []api.Vector{
						{Name: "single", Single: []float32{1, 2, 3}},
					},
					Properties: map[string]any{"foo": "bar"},
					References: api.ObjectReferences{
						"ref": []api.ObjectReference{
							{Collection: "Foo", UUID: testkit.UUID},
							{Collection: "Bar", UUID: testkit.UUID},
						},
					},
				},
			}},
		},
		{
			name: "error on nil uuid",
			err:  testkit.ExpectError,
		},
		{
			name:   "with error",
			object: data.Object{UUID: &testkit.UUID},
			stub: []testkit.Stub[api.ReplaceObjectRequest, any]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stub)
			c := data.NewClient(transport, rd)
			require.NotNil(t, c, "nil client")

			err := c.Replace(t.Context(), tt.object)
			tt.err.Require(t, err, "replace error")
		})
	}
}

func TestClient_Delete(t *testing.T) {
	rd := api.RequestDefaults{
		CollectionName:   "Delete",
		ConsistencyLevel: api.ConsistencyLevelOne,
		Tenant:           "john_doe",
	}

	for _, tt := range []struct {
		name string
		uuid uuid.UUID // ID of the object to be deleted.
		stub []testkit.Stub[api.DeleteObjectRequest, any]
		err  testkit.Error
	}{
		{
			name: "ok",
			uuid: testkit.UUID,
			stub: []testkit.Stub[api.DeleteObjectRequest, any]{{
				Request: &api.DeleteObjectRequest{
					RequestDefaults: rd,
					UUID:            testkit.UUID,
				},
			}},
		},
		{
			name: "with error",
			uuid: testkit.UUID,
			stub: []testkit.Stub[api.DeleteObjectRequest, any]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stub)
			c := data.NewClient(transport, rd)
			require.NotNil(t, c, "nil client")

			err := c.Delete(t.Context(), tt.uuid)
			tt.err.Require(t, err, "delete error")
		})
	}
}
