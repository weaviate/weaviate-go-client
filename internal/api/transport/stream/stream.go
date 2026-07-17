package stream

import (
	"context"
	"time"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	proto "github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/proto/v1"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/ssb"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
	"google.golang.org/grpc"
	protoutil "google.golang.org/protobuf/proto"
)

type Client grpc.BidiStreamingClient[proto.BatchStreamRequest, proto.BatchStreamReply]

// NewStreamFunc opens a new streaming client. The underlying stream
// supports messages up to maxSize bytes when marshaled.
type NewStreamFunc func(context.Context) (client Client, maxSize int, err error)

func (f NewStreamFunc) NewStream(ctx context.Context, rd api.RequestDefaults) (ssb.Stream, error) {
	client, maxSize, err := f(ctx)
	if err != nil {
		return nil, err
	}
	dev.AssertNotNil(client, "client")
	return &batchStream{
		RequestDefaults: rd,
		c:               client,
		maxSize:         maxSize,
	}, nil
}

type batchStream struct {
	api.RequestDefaults
	c       Client
	maxSize int
}

// NewBatch cretes a new sized container that accumulates
// batch data until it reaches [Client]'s maximum request size.
func (bs *batchStream) NewBatch() ssb.Batch {
	return &batch{
		c:       bs.c,
		maxSize: bs.maxSize,
		BatchRequest: api.BatchRequest{
			RequestDefaults: bs.RequestDefaults,
		},
	}
}

func (bs *batchStream) Recv() (ssb.Event, error) {
	var event ssb.Event

	reply, err := bs.c.Recv()
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
				ReconnectAfter: time.Duration(msg.OutOfMemory.WaitTime) * time.Second,
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

func (bs *batchStream) Close() error {
	return nil
}

type batch struct {
	api.BatchRequest
	c       Client
	maxSize int
	curSize int
}

func (b *batch) Add(d ssb.Data) error {
	var pop func()

	switch {
	case d.Object != nil:
		if err := b.AddObject(d.Object); err != nil {
			return err
		}
		pop = b.PopObject
	}

	reqSize, err := b.size()
	if err != nil {
		return err
	}
	itemSize := reqSize - b.curSize
	if reqSize <= b.maxSize {
		b.curSize = reqSize
		return nil
	}

	// The item does not fit in the batch.
	// Defer its removal, then narrow down the reason.
	defer pop()
	if itemSize > b.maxSize {
		return ssb.ErrTooLarge
	} else if reqSize > b.maxSize {
		return ssb.ErrBatchFull
	}

	dev.Unreachable()
	return nil
}

func (b *batch) Send() error {
	req, err := b.MarshalMessage()
	if err != nil {
		return err
	}
	return b.c.Send(req)
}

// size returns the estimated size of the marshaled message, in bytes.
func (b *batch) size() (int, error) {
	m, err := b.MarshalMessage()
	if err != nil {
		return 0, err
	}
	return protoutil.Size(m), nil
}
