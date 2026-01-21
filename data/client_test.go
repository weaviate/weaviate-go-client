package data_test

import (
	"testing"

	"github.com/go-openapi/testify/v2/require"
	"github.com/google/uuid"
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
		name         string
		object       *data.Object                  // Object to be inserted.
		want         *types.Object[map[string]any] // Expected return value.
		stubs        []testkit.Stub[api.InsertObjectRequest, api.InsertObjectResponse]
		requireError require.ErrorAssertionFunc // Set to require.Error to expect an error.
	}{
		{
			name: "nil object",
			stubs: []testkit.Stub[api.InsertObjectRequest, api.InsertObjectResponse]{{
				Request: &api.InsertObjectRequest{RequestDefaults: rd},
				Response: api.InsertObjectResponse{
					UUID:          uuid.Nil,
					CreatedAt:     testkit.Now,
					LastUpdatedAt: testkit.Now,
				},
			}},
			want: &types.Object[map[string]any]{
				UUID:          uuid.Nil,
				CreatedAt:     testkit.Now,
				LastUpdatedAt: testkit.Now,
				References:    (data.References)(nil), // References must be a typed null.
			},
		},
		{
			name: "with data",
			object: &data.Object{
				Vectors: []types.Vector{
					{Name: "single", Single: []float32{1, 2, 3}},
				},
				Properties: map[string]any{"foo": "bar"},
				References: data.References{
					"ref": []data.Reference{
						{Collection: "Foo", UUID: uuid.Nil},
						{Collection: "Bar", UUID: uuid.Nil},
					},
				},
			},
			stubs: []testkit.Stub[api.InsertObjectRequest, api.InsertObjectResponse]{{
				Request: &api.InsertObjectRequest{
					RequestDefaults: rd,
					Vectors: []api.Vector{
						{Name: "single", Single: []float32{1, 2, 3}},
					},
					Properties: map[string]any{"foo": "bar"},
					References: api.ObjectReferences{
						"ref": []api.ObjectReference{
							{Collection: "Foo", UUID: uuid.Nil},
							{Collection: "Bar", UUID: uuid.Nil},
						},
					},
				},
				Response: api.InsertObjectResponse{
					UUID:          uuid.Nil,
					CreatedAt:     testkit.Now,
					LastUpdatedAt: testkit.Now,
					Vectors: map[string]api.Vector{
						"single": {Name: "single", Single: []float32{1, 2, 3}},
					},
					Properties: map[string]any{"foo": "bar"},
					References: api.ObjectReferences{
						"ref": []api.ObjectReference{
							{Collection: "Foo", UUID: uuid.Nil},
							{Collection: "Bar", UUID: uuid.Nil},
						},
					},
				},
			}},
			want: &types.Object[map[string]any]{
				UUID:          uuid.Nil,
				CreatedAt:     testkit.Now,
				LastUpdatedAt: testkit.Now,
				Vectors: map[string]types.Vector{
					"single": {Name: "single", Single: []float32{1, 2, 3}},
				},
				Properties: map[string]any{"foo": "bar"},
				References: data.References{
					"ref": []data.Reference{
						{Collection: "Foo", UUID: uuid.Nil},
						{Collection: "Bar", UUID: uuid.Nil},
					},
				},
			},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[api.InsertObjectRequest, api.InsertObjectResponse]{
				{Err: testkit.ErrWhaam},
			},
			requireError: require.Error,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := data.NewClient(transport, rd)
			require.NotNil(t, c, "nil client")

			requireError := require.NoError
			if tt.requireError != nil {
				requireError = tt.requireError
			}

			got, err := c.Insert(t.Context(), tt.object)
			requireError(t, err, "insert error")
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
		name         string
		object       data.Object                   // Object to be replaced.
		want         *types.Object[map[string]any] // Expected return value.
		stub         []testkit.Stub[api.ReplaceObjectRequest, api.ReplaceObjectResponse]
		requireError require.ErrorAssertionFunc // Set to require.Error to expect an error.
	}{
		{
			name: "with data",
			object: data.Object{
				UUID: &uuid.Nil,
				Vectors: []types.Vector{
					{Name: "single", Single: []float32{1, 2, 3}},
				},
				Properties: map[string]any{"foo": "bar"},
				References: data.References{
					"ref": []data.Reference{
						{Collection: "Foo", UUID: uuid.Nil},
						{Collection: "Bar", UUID: uuid.Nil},
					},
				},
			},
			stub: []testkit.Stub[api.ReplaceObjectRequest, api.ReplaceObjectResponse]{{
				Request: &api.ReplaceObjectRequest{
					RequestDefaults: rd,
					UUID:            &uuid.Nil,
					Vectors: []api.Vector{
						{Name: "single", Single: []float32{1, 2, 3}},
					},
					Properties: map[string]any{"foo": "bar"},
					References: api.ObjectReferences{
						"ref": []api.ObjectReference{
							{Collection: "Foo", UUID: uuid.Nil},
							{Collection: "Bar", UUID: uuid.Nil},
						},
					},
				},
				Response: api.ReplaceObjectResponse{
					UUID:          uuid.Nil,
					CreatedAt:     testkit.Now,
					LastUpdatedAt: testkit.Now,
					Properties:    map[string]any{"foo": "bar"},
					References: api.ObjectReferences{
						"ref": []api.ObjectReference{
							{Collection: "Foo", UUID: uuid.Nil},
							{Collection: "Bar", UUID: uuid.Nil},
						},
					},
				},
			}},
			want: &types.Object[map[string]any]{
				UUID:          uuid.Nil,
				CreatedAt:     testkit.Now,
				LastUpdatedAt: testkit.Now,
				Properties:    map[string]any{"foo": "bar"},
				References: data.References{
					"ref": []data.Reference{
						{Collection: "Foo", UUID: uuid.Nil},
						{Collection: "Bar", UUID: uuid.Nil},
					},
				},
			},
		},
		{
			name:         "error on nil uuid",
			requireError: require.Error,
		},
		{
			name:   "with error",
			object: data.Object{UUID: &uuid.Nil},
			stub: []testkit.Stub[api.ReplaceObjectRequest, api.ReplaceObjectResponse]{
				{Err: testkit.ErrWhaam},
			},
			requireError: require.Error,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stub)
			c := data.NewClient(transport, rd)
			require.NotNil(t, c, "nil client")

			requireError := require.NoError
			if tt.requireError != nil {
				requireError = tt.requireError
			}

			got, err := c.Replace(t.Context(), tt.object)
			requireError(t, err, "replace error")
			require.Equal(t, tt.want, got, "returned object")
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
		name         string
		uuid         uuid.UUID // ID of the object to be deleted.
		stub         []testkit.Stub[api.DeleteObjectRequest, any]
		requireError require.ErrorAssertionFunc // Set to require.Error to expect an error.
	}{
		{
			name: "ok",
			uuid: uuid.Nil,
			stub: []testkit.Stub[api.DeleteObjectRequest, any]{{
				Request: &api.DeleteObjectRequest{
					RequestDefaults: rd,
					UUID:            uuid.Nil,
				},
			}},
		},
		{
			name: "with error",
			uuid: uuid.Nil,
			stub: []testkit.Stub[api.DeleteObjectRequest, any]{
				{Err: testkit.ErrWhaam},
			},
			requireError: require.Error,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stub)
			c := data.NewClient(transport, rd)
			require.NotNil(t, c, "nil client")

			requireError := require.NoError
			if tt.requireError != nil {
				requireError = tt.requireError
			}

			err := c.Delete(t.Context(), tt.uuid)
			requireError(t, err, "delete error")
		})
	}
}
