package weaviate

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/google/uuid"

	"github.com/weaviate/weaviate-go-client/v6/types"
	"github.com/weaviate/weaviate-go-client/v6/weaviate/internal"
)

// DataClient handles CRUD operations
type DataClient struct {
	collection *CollectionClient
}

// Insert inserts a single object.
func (d *DataClient) Insert(ctx context.Context, opts ...InsertOption) (uuid.UUID, error) {
	options := &insertOptions{}
	for _, opt := range opts {
		opt(options)
	}

	// Convert data to properties map
	props, err := internal.ToPropertiesMap(options.properties)
	if err != nil {
		return uuid.Nil, errors.New("convert data to properties: " + err.Error())
	}

	// Build request body
	body := map[string]any{
		"class":      d.collection.name,
		"properties": props,
	}

	if options.id != "" {
		body["id"] = options.id
	}
	if options.vectors != nil {
		// Convert vectors for the API
		vectorMap := make(map[string]any)
		for name, vec := range options.vectors {
			var w internal.Vector = (*wrappedVector)(&vec)
			if w.IsMulti() {
				vectorMap[name] = w.ToFloat32Multi()
			} else {
				vectorMap[name] = w.ToFloat32()
			}
		}
		if len(vectorMap) == 1 {
			if v, ok := vectorMap[""]; ok {
				body["vector"] = v
				delete(vectorMap, "")
			}
		}
		if len(vectorMap) > 0 {
			body["vectors"] = vectorMap
		}
	}
	if d.collection.tenant != "" {
		body["tenant"] = d.collection.tenant
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return uuid.Nil, errors.New("marshal request body: " + err.Error())
	}

	url := "localhost:8080/v1/objects" // d.collection.client.baseURL + "/v1/objects"
	if d.collection.consistency != "" {
		url += "?consistency_level=" + string(d.collection.consistency)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return uuid.Nil, errors.New("create request: " + err.Error())
	}

	//d.collection.client.addHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	rd, _ := httputil.DumpRequest(req, true)
	fmt.Println(string(rd))

	//resp, err := d.collection.client.http.Do(req)
	// if err != nil {
	// 	return uuid.Nil, errors.New("insert object: " + err.Error())
	// }
	//defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
	// 	respBody, _ := io.ReadAll(resp.Body)
	// 	return "", fault.FromHTTPStatus(resp.StatusCode, string(respBody))
	// }

	var result struct {
		ID string `json:"id"`
	}
	// if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
	// 	return "", fault.WrapClientError("decode response", err)
	// }

	result.ID = uuid.NewString()

	return uuid.Parse(result.ID)
}

// InsertOption configures insert operations.
type InsertOption func(*insertOptions)

type insertOptions struct {
	id         string
	vectors    types.Vectors
	properties types.Properties
	references map[string][]types.ObjectReference
	validate   bool
}

// WithID specifies the object ID.
func WithID(id string) InsertOption {
	return func(o *insertOptions) {
		o.id = id
	}
}

// WithVectors specifies custom vectors.
func WithProperties(properties types.Properties) InsertOption {
	return func(o *insertOptions) {
		o.properties = properties
	}
}

type vectorInput interface {
	types.Vector | types.Vectors | []types.Vector
}

// WithVector specifies a custom vector or vectors.
func WithVector[T vectorInput](vector ...T) InsertOption {
	return func(o *insertOptions) {
		if len(vector) == 0 {
			return
		}
		if o.vectors == nil {
			o.vectors = make(types.Vectors)
		}
		// Type switch on the first element to determine the type
		switch any(vector[0]).(type) {
		case types.Vector:
			for _, v := range vector {
				vv := any(v).(types.Vector)
				o.vectors[vv.Name] = vv
			}
		case types.Vectors:
			for _, v := range vector {
				vv := any(v).(types.Vectors)
				for _, vvv := range vv {
					o.vectors[vvv.Name] = vvv
				}
			}
		case []types.Vector:
			for _, v := range vector {
				vv := any(v).([]types.Vector)
				for _, vvv := range vv {
					o.vectors[vvv.Name] = vvv
				}
			}
		default:
			panic("unsupported vector type")
		}
	}
}

// WithReferences specifies cross-references.
func WithReference(name string, ref types.ObjectReference) InsertOption {
	return func(o *insertOptions) {
		if o.references == nil {
			o.references = make(map[string][]types.ObjectReference)
		}
		o.references[name] = append(o.references[name], ref)
	}
}

// WithValidation enables schema validation.
func WithValidation() InsertOption {
	return func(o *insertOptions) {
		o.validate = true
	}
}

// // Update updates an existing object (PATCH - partial update)
// func (d *DataClient) Update(ctx context.Context, id uuid.UUID, data any, opts ...UpdateOption) error

// // Replace replaces an existing object (PUT - full replace)
// func (d *DataClient) Replace(ctx context.Context, id uuid.UUID, data any, opts ...ReplaceOption) error

// // Delete deletes an object by ID
// func (d *DataClient) Delete(ctx context.Context, id uuid.UUID) error

// // DeleteMany deletes objects matching a filter
// func (d *DataClient) DeleteMany(ctx context.Context, filter Filter, opts ...DeleteManyOption) (*DeleteManyResult, error)
