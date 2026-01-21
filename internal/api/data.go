package api

import (
	"bytes"
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/rest"
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
type InsertObjectResponse struct {
	UUID          uuid.UUID
	Properties    map[string]any
	References    ObjectReferences
	Vectors       Vectors
	CreatedAt     time.Time
	LastUpdatedAt time.Time
}

var _ json.Unmarshaler = (*InsertObjectResponse)(nil)

type (
	ObjectReferences map[string][]ObjectReference
	ObjectReference  struct {
		Collection string    // Collection the referenced object belongs to.
		UUID       uuid.UUID // UUID of the referenced object.
	}
)

var (
	_ encoding.TextMarshaler   = (*ObjectReference)(nil)
	_ encoding.TextUnmarshaler = (*ObjectReference)(nil)
)

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

func (o *ObjectReference) UnmarshalText(text []byte) error {
	text, ok := bytes.CutPrefix(text, beaconPrefix)
	if !ok {
		return errors.New("not a beacon")
	}
	parts := bytes.Split(text, beaconSep)

	if len(parts) == 2 {
		o.Collection, parts = string(parts[0]), parts[1:]
	}

	if len(parts) != 1 {
		return fmt.Errorf("beacon %s is malformed", string(text))
	}

	id, err := uuid.ParseBytes(parts[0])
	if err != nil {
		return err
	}
	o.UUID = id
	return nil
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
	var o struct {
		rest.Object
		Properties map[string]json.RawMessage `json:"properties,omitempty"`
		Vectors    Vectors                    `json:"vectors,omitempty"`
	}
	if err := json.Unmarshal(data, &o); err != nil {
		return err
	}

	// Expect most of the properties will be data, not references.
	properties := make(map[string]any, len(o.Properties))
	references := make(ObjectReferences)
	for k, data := range o.Properties {
		var refs []ObjectReference
		if err := json.Unmarshal(data, &refs); err == nil {
			references[k] = refs
			continue
		}

		var v any
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		properties[k] = v
	}

	*i = InsertObjectResponse{
		UUID:          *o.Id,
		CreatedAt:     time.UnixMilli(o.CreationTimeUnix),
		LastUpdatedAt: time.UnixMilli(o.LastUpdateTimeUnix),
		Vectors:       o.Vectors,
		Properties:    properties,
		References:    references,
	}
	return nil
}

// ReplaceObjectRequest mirrors InsertObjectRequest but uses PUT method.
// The collection name is sent as a path parameter.
type ReplaceObjectRequest struct {
	RequestDefaults
	UUID       uuid.UUID
	Properties map[string]any
	References ObjectReferences
	Vectors    []Vector
}

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

type ReplaceObjectResponse InsertObjectResponse

var _ json.Unmarshaler = (*ReplaceObjectResponse)(nil)

func (r *ReplaceObjectResponse) UnmarshalJSON(data []byte) error {
	// InsertObjectResponse implements json.Unmarshaler,
	// and response structs are identical.
	var ior InsertObjectResponse
	if err := json.Unmarshal(data, &ior); err != nil {
		return err
	}
	*r = ReplaceObjectResponse(ior)
	return nil
}
