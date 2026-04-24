package api

import (
	"encoding"
	"encoding/json"
	"maps"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/rest"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
	"github.com/weaviate/weaviate-go-client/v6/internal/transports"
)

// InsertObjectRequest inserts a new object into a collection.
type InsertObjectRequest struct {
	RequestDefaults
	UUID       *uuid.UUID
	Properties map[string]any
	References ObjectReferences
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
type InsertObjectResponse uuid.UUID

var _ json.Unmarshaler = (*InsertObjectResponse)(nil)

func (r *InsertObjectResponse) UUID() *uuid.UUID { return (*uuid.UUID)(r) }

type (
	ObjectReferences map[string][]ObjectReference
	ObjectReference  struct {
		Collection string    // Collection the referenced object belongs to.
		UUID       uuid.UUID // UUID of the referenced object.
	}
)

var _ encoding.TextMarshaler = (*ObjectReference)(nil)

var (
	beaconPrefix = []byte("weaviate://localhost/")
	beaconSep    = []byte("/")
)

// MarshalText formats the object reference as a beacon.
// json.Marshal will call this method and encode the result as a JSON string.
func (o *ObjectReference) MarshalText() ([]byte, error) {
	id, err := o.UUID.MarshalText()
	if err != nil {
		return nil, err
	}
	b := append([]byte(nil), beaconPrefix...)
	if o.Collection != "" {
		b = append(b, o.Collection...)
		b = append(b, beaconSep...)
	}
	return append(b, id...), nil
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

	properties := make(map[string]any, len(r.Properties)+len(r.References))
	maps.Copy(properties, r.Properties)

	for name, ref := range r.References {
		properties[name] = ref
	}

	req := &rest.Object{
		Class:      r.CollectionName,
		Tenant:     r.Tenant,
		Id:         r.UUID,
		Properties: properties,
		Vectors:    vectors,
	}

	return json.Marshal(req)
}

// UnmarshalJSON implements json.Unmarshaler via [rest.Object].
func (i *InsertObjectResponse) UnmarshalJSON(data []byte) error {
	var o rest.Object
	if err := json.Unmarshal(data, &o); err != nil {
		return err
	}
	*i = InsertObjectResponse(*o.Id)
	return nil
}

// ReplaceObjectRequest mirrors InsertObjectRequest but uses PUT method.
// The collection name is sent as a path parameter.
type ReplaceObjectRequest InsertObjectRequest

var _ transports.Endpoint = (*ReplaceObjectRequest)(nil)

func (*ReplaceObjectRequest) Method() string { return http.MethodPut }
func (r *ReplaceObjectRequest) Path() string {
	return "/objects/" + r.CollectionName + "/" + r.UUID.String()
}

func (r *ReplaceObjectRequest) Query() url.Values {
	if r.ConsistencyLevel != consistencyLevelUndefined {
		return url.Values{"consistency_level": {string(r.ConsistencyLevel)}}
	}
	return nil
}

func (r *ReplaceObjectRequest) Body() any {
	dev.AssertNotNil(r.UUID, "object uuid")

	// InsertObjectRequest already implements json.Marshaler.
	// For replace, CollectionName and UUID should not part of the payload.
	return &InsertObjectRequest{
		RequestDefaults: RequestDefaults{
			Tenant:           r.Tenant,
			ConsistencyLevel: r.ConsistencyLevel,
		},
		Properties: r.Properties,
		References: r.References,
		Vectors:    r.Vectors,
	}
}

// DeleteObjectRequest deletes an object by its UUID.
type DeleteObjectRequest struct {
	transports.BaseEndpoint

	RequestDefaults
	UUID uuid.UUID
}

var _ transports.Endpoint = (*DeleteObjectRequest)(nil)

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
