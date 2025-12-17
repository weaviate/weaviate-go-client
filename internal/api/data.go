package api

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/gen/rest"
	"github.com/weaviate/weaviate-go-client/v6/internal/transport"
)

type InsertObjectRequest struct {
	RequestDefaults
	UUID       uuid.UUID
	Properties any
	References any // TODO(dyma): define
	Vectors    Vectors
}

var _ transport.Endpoint = (*InsertObjectRequest)(nil)

func (r *InsertObjectRequest) Method() string { return http.MethodPost }
func (r *InsertObjectRequest) Path() string   { return "/objects" }
func (r *InsertObjectRequest) Query() url.Values {
	if r.ConsistencyLevel != consistencyLevelUndefined {
		return url.Values{"consistency_level": {string(r.ConsistencyLevel)}}
	}
	return nil
}

func (r *InsertObjectRequest) Body() any {
	return &rest.Object{
		Class:      r.CollectionName,
		Tenant:     r.Tenant,
		Id:         r.UUID,
		Properties: nil, // TODO: convert to map[string]interface{}
		Vectors:    r.Vectors.toMap(),
	}
}

func (vs Vectors) toMap() map[string]any {
	out := make(map[string]any, len(vs))
	for name, v := range vs {
		if v.Single != nil {
			out[name] = v.Single
		} else if v.Multi != nil {
			out[name] = v.Multi
		}
	}
	return out
}

type InsertObjectResponse struct {
	UUID               uuid.UUID
	Properties         map[string]any
	References         any // TODO(dyma): define
	Vectors            Vectors
	CreationTimeUnix   int64
	LastUpdateTimeUnix int64
}

var _ json.Unmarshaler = (*InsertObjectResponse)(nil)

// UnmarshalJSON implements json.Unmarshaler via [rest.Object].
func (i *InsertObjectResponse) UnmarshalJSON(data []byte) error {
	var res rest.Object
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	*i = InsertObjectResponse{
		UUID:       res.Id,
		Properties: res.Properties,
		References: nil,
	}
	return nil
}

type DeleteObjectRequest struct {
	endpoint
	RequestDefaults
	UUID uuid.UUID
}

var _ transport.Endpoint = (*DeleteObjectRequest)(nil)

func (r DeleteObjectRequest) Method() string { return http.MethodDelete }
func (r DeleteObjectRequest) Path() string {
	return "/objects/" + r.CollectionName + "/" + r.UUID.String()
}
