package api

import (
	"encoding"
	"encoding/json"
	"maps"
	"net/http"
	"net/url"
	"time"

	proto "github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/proto/v1"

	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/rest"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/transport"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
	"github.com/weaviate/weaviate-go-client/v6/internal/transports"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

// InsertObjectsRequest inserts a batch of objects into a collection.
type InsertObjectsRequest struct {
	RequestDefaults
	Objects []BatchObject
}

func (*InsertObjectsRequest) Method() transport.MethodFunc[proto.BatchObjectsRequest, proto.BatchObjectsReply] {
	return proto.WeaviateClient.BatchObjects
}

func (r *InsertObjectsRequest) Body() transport.MessageMarshaler[proto.BatchObjectsRequest] {
	return r
}

type BatchObject struct {
	// Batch API does not allow inserting objects without UUIDs,
	// the way POST /objects does. Normally, api package's policy
	// would be to deal with such quirks internally and not expose
	// that to the caller.
	// However, in order to map an error from the batch response
	// to the right UUID on return, the caller MUST know the UUID
	// prior to sending the request.
	//
	// The zero value of BatchObject is useful. If is perfectly OK
	// to insert an object with [uuid.Nil] ID, and new(BatchObject)
	// will produce exactly that.
	UUID       uuid.UUID
	Properties map[string]any
	References ObjectReferences
	Vectors    []Vector
}

var (
	_ transport.Message[proto.BatchObjectsRequest, proto.BatchObjectsReply] = (*InsertObjectsRequest)(nil)
	_ transport.MessageMarshaler[proto.BatchObjectsRequest]                 = (*InsertObjectsRequest)(nil)
)

// MarshalMessage implements [transport.MessageMarshaler].
func (r *InsertObjectsRequest) MarshalMessage() (*proto.BatchObjectsRequest, error) {
	dev.AssertNotNil(r, "r")
	batch := make([]*proto.BatchObject, len(r.Objects))
	for i := range r.Objects {
		bo, err := marshalBatchObject(&r.Objects[i], r.RequestDefaults)
		if err != nil {
			return nil, err
		}
		batch[i] = bo
	}

	return &proto.BatchObjectsRequest{
		ConsistencyLevel: r.ConsistencyLevel.proto(),
		Objects:          batch,
	}, nil
}

type InsertObjectsResponse struct {
	Took      time.Duration
	Positions []int32  // Positional indices of the failed objects. Aligned with Errors.
	Errors    []string // Error messages for failed objects. Aligned with Indices.
}

var _ transport.MessageUnmarshaler[proto.BatchObjectsReply] = (*InsertObjectsResponse)(nil)

// UnmarshalMessage implements [transport.MessageUnmarshaler].
func (r *InsertObjectsResponse) UnmarshalMessage(reply *proto.BatchObjectsReply) error {
	*r = InsertObjectsResponse{
		Took: time.Duration(reply.Took) * time.Second,
	}

	for _, e := range reply.GetErrors() {
		r.Positions = append(r.Positions, e.Index)
		r.Errors = append(r.Errors, e.Error)
	}
	return nil
}

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
func (r *ReplaceObjectRequest) MarshalJSON() ([]byte, error) {
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

// ReplaceObjectRequest mirrors InsertObjectRequest but uses PUT method.
// The collection name is sent as a path parameter.
type ReplaceObjectRequest struct {
	RequestDefaults
	UUID       *uuid.UUID
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
	dev.AssertNotNil(r.UUID, "object uuid")

	// InsertObjectRequest already implements json.Marshaler.
	// For replace, CollectionName and UUID should not part of the payload.
	return &ReplaceObjectRequest{
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

func marshalBatchObject(bo *BatchObject, rd RequestDefaults) (*proto.BatchObject, error) {
	var vectors []*proto.Vectors
	for i := range bo.Vectors {
		v, err := marshalVector(&bo.Vectors[i])
		if err != nil {
			return nil, err
		}
		vectors = append(vectors, v)
	}

	var properties *proto.BatchObject_Properties
	if len(bo.Properties) > 0 || len(bo.References) > 0 {
		properties = new(proto.BatchObject_Properties)
		if err := marshalProperties(bo.Properties, properties); err != nil {
			return nil, err
		}
		if err := marshalReferences(bo.References, properties); err != nil {
			return nil, err
		}
	}

	return &proto.BatchObject{
		Uuid:       bo.UUID.String(),
		Collection: rd.CollectionName,
		Tenant:     rd.Tenant,
		Vectors:    vectors,
		Properties: properties,
	}, nil
}

func marshalProperties(properties map[string]any, dest *proto.BatchObject_Properties) error {
	if len(properties) == 0 {
		return nil
	}
	nonRef, err := structpb.NewStruct(properties)
	if err != nil {
		return err
	}

	// TODO(dyma): move object / array properties out of nonRef
	dest.NonRefProperties = nonRef
	return nil
}

func marshalReferences(references ObjectReferences, dest *proto.BatchObject_Properties) error {
	if len(references) == 0 {
		return nil
	}
	var single []*proto.BatchObject_SingleTargetRefProps
	var multi []*proto.BatchObject_MultiTargetRefProps
	for name, refs := range references {
		uuids := make(map[string][]string, 0)
		for _, ref := range refs {
			uuids[ref.Collection] = append(uuids[ref.Collection], ref.UUID.String())
		}

		for collection := range uuids {
			if collection == "" {
				single = append(single, &proto.BatchObject_SingleTargetRefProps{
					PropName: name,
					Uuids:    uuids[collection],
				})
			} else {
				multi = append(multi, &proto.BatchObject_MultiTargetRefProps{
					PropName:         name,
					Uuids:            uuids[collection],
					TargetCollection: collection,
				})
			}
		}
	}
	dest.SingleTargetRefProps = single
	dest.MultiTargetRefProps = multi
	return nil
}
