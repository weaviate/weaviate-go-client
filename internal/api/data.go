package api

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/rest"
	"github.com/weaviate/weaviate-go-client/v6/internal/transports"
)

// InsertObjectRequest inserts a new object into a collection.
type InsertObjectRequest struct {
	RequestDefaults
	UUID       *uuid.UUID
	Properties map[string]any
	Vectors    []Vector
}

var (
	_ transports.Endpoint = (*InsertObjectRequest)(nil)
	_ json.Marshaler      = (*InsertObjectRequest)(nil)
)

func (*InsertObjectRequest) Method() string { return http.MethodPost }
func (*InsertObjectRequest) Path() string   { return "/objects" }
func (r *InsertObjectRequest) Query() url.Values {
	if r.ConsistencyLevel != consistencyLevelUndefined {
		return url.Values{"consistency_level": {string(r.ConsistencyLevel)}}
	}
	return nil
}
func (r *InsertObjectRequest) Body() any { return r }

// InsertObjectResponses reads the data for the newly inserted object.
type InsertObjectResponse struct {
	UUID               uuid.UUID
	Properties         map[string]any
	Vectors            Vectors
	CreationTimeUnix   int64
	LastUpdateTimeUnix int64
}

var _ json.Unmarshaler = (*InsertObjectResponse)(nil)

// MarshalJSON implements json.Marshaler via [rest.Object].
func (r *InsertObjectRequest) MarshalJSON() ([]byte, error) {
	vectors := make(map[string]any, len(r.Vectors))
	for _, v := range r.Vectors {
		if v.Single != nil {
			vectors[v.Name] = v.Single
		} else if v.Multi != nil {
			vectors[v.Name] = v.Multi
		}
	}

	req := &rest.Object{
		Class:      r.CollectionName,
		Tenant:     r.Tenant,
		Id:         r.UUID,
		Properties: r.Properties,
		Vectors:    vectors,
	}
	return json.Marshal(req)
}

// UnmarshalJSON implements json.Unmarshaler via [rest.Object].
func (i *InsertObjectResponse) UnmarshalJSON(data []byte) error {
	var res rest.Object
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	*i = InsertObjectResponse{
		UUID:       *res.Id,
		Properties: res.Properties,
	}
	return nil
}
