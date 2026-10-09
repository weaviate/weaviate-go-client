package api

import (
	"github.com/weaviate/weaviate-go-client/v6/internal/api/transport"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"

	proto "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
	protoutil "google.golang.org/protobuf/proto"
)

// StartStreamRequest initiates an SSB stream with the given consistency level.
type StartStreamRequest struct {
	ConsistencyLevel ConsistencyLevel // Consistency level for incoming writes.
}

var _ transport.MessageMarshaler[proto.BatchStreamRequest] = (*StartStreamRequest)(nil)

func (r *StartStreamRequest) MarshalMessage() (*proto.BatchStreamRequest, error) {
	return &proto.BatchStreamRequest{
		Message: &proto.BatchStreamRequest_Start_{
			Start: &proto.BatchStreamRequest_Start{
				ConsistencyLevel: r.ConsistencyLevel.proto(),
			},
		},
	}, nil
}

// StopStreamRequest terminates the SSB stream.
var StopStreamRequest = &stopStreamRequest{}

type stopStreamRequest struct{}

var _ transport.MessageMarshaler[proto.BatchStreamRequest] = (*stopStreamRequest)(nil)

func (*stopStreamRequest) MarshalMessage() (*proto.BatchStreamRequest, error) {
	return &proto.BatchStreamRequest{
		Message: &proto.BatchStreamRequest_Stop_{
			Stop: new(proto.BatchStreamRequest_Stop),
		},
	}, nil
}

// BatchRequest accumulates objects and references into batched requests.
type BatchRequest struct {
	MaxSize int // Maximum size request size in bytes.

	objects    []*proto.BatchObject
	references []*proto.BatchReference
}

var _ transport.MessageMarshaler[proto.BatchStreamRequest] = (*BatchRequest)(nil)

// MarshalMessage creates a new Data request.
func (req *BatchRequest) MarshalMessage() (*proto.BatchStreamRequest, error) {
	return &proto.BatchStreamRequest{
		Message: &proto.BatchStreamRequest_Data_{
			Data: &proto.BatchStreamRequest_Data{
				Objects: &proto.BatchStreamRequest_Data_Objects{
					Values: req.objects,
				},
				References: &proto.BatchStreamRequest_Data_References{
					Values: req.references,
				},
			},
		},
	}, nil
}

func (req *BatchRequest) Add(v any) (added, full bool) {
	var pop func()

	switch v := v.(type) {
	case *proto.BatchObject:
		req.objects = append(req.objects, v)
		pop = req.popObject
	case *proto.BatchReference:
		req.references = append(req.references, v)
		pop = req.popReference
	default:
		dev.Unreachable()
	}

	reqSize := req.size()
	full = reqSize >= req.MaxSize
	added = reqSize <= req.MaxSize
	if !added {
		pop()
	}
	return
}

// size returns the estimated size of the marshaled message, in bytes.
func (req *BatchRequest) size() int {
	m, err := req.MarshalMessage()
	dev.Assert(err == nil, "marshal message failed")
	return protoutil.Size(m)
}

// Pop removes the last added [BatchObject] from the batch.
// Safe to call on an empty batch.
func (req *BatchRequest) popObject() {
	req.objects = req.objects[:len(req.objects)-1]
}

// popReference removes the last added [Reference] from the batch.
// Safe to call on an empty batch.
func (req *BatchRequest) popReference() {
	req.references = req.references[:len(req.references)-1]
}
