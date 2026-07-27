package stream

import (
	"context"
	"time"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	proto "github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/proto/v1"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/ssb"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/transport"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
	protoutil "google.golang.org/protobuf/proto"
)

// Transport supports asynchronous bidirectional streaming.
type Transport interface {
	// NewStream opens a new [transport.BatchStream].
	NewStream(context.Context) (transport.BatchStream, error)
	// The maximum request size in bytes that can be sent via the stream.
	MaxSize() int
}

func NewAdapter(t Transport, rd api.RequestDefaults) *TransportAdapter {
	return &TransportAdapter{
		transport: t,
		defaults:  rd,
	}
}

type TransportAdapter struct {
	transport Transport
	defaults  api.RequestDefaults
}

func (tad *TransportAdapter) NewStream(ctx context.Context) (ssb.Stream, error) {
	stream, err := tad.transport.NewStream(ctx)
	if err != nil {
		return nil, err
	}
	sad := &streamAdapter{stream: stream}
	if err := sad.send(&api.StartStreamRequest{
		ConsistencyLevel: tad.defaults.ConsistencyLevel,
	}); err != nil {
		return nil, err
	}
	return sad, nil
}

func (tad *TransportAdapter) NewRequest() ssb.BatchRequest {
	maxSize := tad.transport.MaxSize()
	dev.Assert(maxSize > 0, "max grpc request <= 0")
	return &api.BatchRequest{
		MaxSize: tad.transport.MaxSize(),
	}
}

func (tad *TransportAdapter) Prepare(d ssb.Data) (v any, err error) {
	switch {
	case d.Object != nil:
		v, err = api.MarshalBatchObject(d.Object, tad.defaults)
	case d.Reference != nil:
		v = api.MarshalBatchReference(d.Reference, tad.defaults)
	default:
		dev.Unreachable()
	}

	if err == nil &&
		protoutil.Size(v.(protoutil.Message)) > tad.transport.MaxSize() {
		err = ssb.ErrTooLarge
	}
	return
}

// streamAdapter implements [ssb.Stream] on top of [transport.BatchStream].
type streamAdapter struct{ stream transport.BatchStream }

func (sad *streamAdapter) Send(req ssb.BatchRequest) error {
	dev.AssertType[*api.BatchRequest](req, "batch request")
	return sad.send(req.(*api.BatchRequest))
}

func (sad *streamAdapter) Recv() (ssb.Event, error) {
	var event ssb.Event

	reply, err := sad.stream.Recv()
	if err != nil {
		return event, err
	}

	switch msg := reply.GetMessage().(type) {
	case *proto.BatchStreamReply_Started_:
		event.Started = msg.Started != nil
	case *proto.BatchStreamReply_ShuttingDown_:
		event.ShuttingDown = msg.ShuttingDown != nil
	case *proto.BatchStreamReply_Backoff_:
		if msg.Backoff != nil {
			batchSize := int(msg.Backoff.BatchSize)
			event.Backoff = &batchSize
		}
	case *proto.BatchStreamReply_OutOfMemory_:
		if msg.OutOfMemory != nil {
			event.OOM = &ssb.OOM{
				ExitAfter: time.Duration(msg.OutOfMemory.WaitTime) * time.Second,
			}
		}
	case *proto.BatchStreamReply_Acks_:
		if msg.Acks != nil {
			event.Acks = append(event.Acks, msg.Acks.Uuids...)
			event.Acks = append(event.Acks, msg.Acks.Beacons...)
		}

	case *proto.BatchStreamReply_Results_:
		event.Results = &ssb.Results{
			OK:     make([]string, len(msg.Results.Successes)),
			Failed: internal.MakeMap[string, string](len(msg.Results.Errors)),
		}
		for i, ok := range msg.Results.Successes {
			switch d := ok.Detail.(type) {
			case *proto.BatchStreamReply_Results_Success_Uuid:
				event.Results.OK[i] = d.Uuid
			case *proto.BatchStreamReply_Results_Success_Beacon:
				event.Results.OK[i] = d.Beacon
			}
		}

		for _, e := range msg.Results.Errors {
			switch d := e.Detail.(type) {
			case *proto.BatchStreamReply_Results_Error_Uuid:
				event.Results.Failed[d.Uuid] = e.Error
			case *proto.BatchStreamReply_Results_Error_Beacon:
				event.Results.Failed[d.Beacon] = e.Error
			}
		}
	}

	return event, nil
}

func (sad *streamAdapter) Close() error {
	if err := sad.send(api.StopStreamRequest); err != nil {
		return err
	}
	return sad.stream.CloseSend()
}

func (sad *streamAdapter) send(mm transport.MessageMarshaler[proto.BatchStreamRequest]) error {
	req, err := mm.MarshalMessage()
	if err != nil {
		return err
	}
	return sad.stream.Send(req)
}
