package stream

import (
	"context"
	"errors"
	"time"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/ssb"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/transport"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
	proto "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
	protoutil "google.golang.org/protobuf/proto"
)

func NewAdapter(t transport.Streaming, rd api.RequestDefaults) *TransportAdapter {
	return &TransportAdapter{
		transport: t,
		defaults:  rd,
	}
}

type TransportAdapter struct {
	transport transport.Streaming
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

func (sad *streamAdapter) Send(req any) error {
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
			event.Ack = true
		}

	case *proto.BatchStreamReply_Results_:
		event.Results = &ssb.Results{
			OK:     make([]string, len(msg.Results.Successes)),
			Failed: internal.MakeMap[string, error](len(msg.Results.Errors)),
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
				event.Results.Failed[d.Uuid] = errors.New(e.Error)
			case *proto.BatchStreamReply_Results_Error_Beacon:
				event.Results.Failed[d.Beacon] = errors.New(e.Error)
			}
		}
	}

	return event, nil
}

func (sad *streamAdapter) Close() (err error) {
	defer sad.stream.CloseSend() // nolint:errcheck
	return sad.send(api.StopStreamRequest)
}

func (sad *streamAdapter) send(mm transport.MessageMarshaler[proto.BatchStreamRequest]) error {
	req, err := mm.MarshalMessage()
	if err != nil {
		return err
	}
	return sad.stream.Send(req)
}
