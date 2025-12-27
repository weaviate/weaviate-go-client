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
	Properties map[string]any
	References any // TODO(dyma): define
	Vectors    []Vector
}

var (
	_ transport.Endpoint = (*InsertObjectRequest)(nil)
	_ json.Marshaler     = (*InsertObjectRequest)(nil)
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

type InsertObjectResponse struct {
	UUID               uuid.UUID
	Properties         map[string]any
	References         any // TODO(dyma): define
	Vectors            Vectors
	CreationTimeUnix   int64
	LastUpdateTimeUnix int64
}

var _ json.Unmarshaler = (*InsertObjectResponse)(nil)

// ReplaceObjectRequest mirrors InsertObjectRequest but uses PUT method instead of POST.
// Also the name of the collection is sent as a path parameter.
type ReplaceObjectRequest InsertObjectRequest

var _ transport.Endpoint = (*ReplaceObjectRequest)(nil)

func (*ReplaceObjectRequest) Method() string { return http.MethodPut }
func (*ReplaceObjectRequest) Path() string   { return "/objects" }
func (r *ReplaceObjectRequest) Query() url.Values {
	if r.ConsistencyLevel != consistencyLevelUndefined {
		return url.Values{"consistency_level": {string(r.ConsistencyLevel)}}
	}
	return nil
}
func (r *ReplaceObjectRequest) Body() any { return (*InsertObjectRequest)(r) }

type ReplaceObjectResponse InsertObjectResponse

type DeleteObjectRequest struct {
	transport.BaseEndpoint
	RequestDefaults
	UUID uuid.UUID
}

var _ transport.Endpoint = (*DeleteObjectRequest)(nil)

func (*DeleteObjectRequest) Method() string { return http.MethodDelete }
func (r *DeleteObjectRequest) Path() string {
	return "/objects/" + r.CollectionName + "/" + r.UUID.String()
}

func (r *DeleteObjectRequest) Query() url.Values {
	if r.Tenant == "" && r.ConsistencyLevel == consistencyLevelUndefined {
		return nil
	}

	q := make(url.Values)
	if r.Tenant != "" {
		q.Add("tenant", r.Tenant)
	}
	if r.ConsistencyLevel != consistencyLevelUndefined {
		q.Add("consistency_level", string(r.ConsistencyLevel))
	}
	return q
}

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
		UUID:       res.Id,
		Properties: res.Properties,
		References: nil,
	}
	return nil
}
