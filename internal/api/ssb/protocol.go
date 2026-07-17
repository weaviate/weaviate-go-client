package ssb

import (
	"context"
	"errors"
	"time"

	"github.com/weaviate/weaviate-go-client/v6/internal/api"
)

type OpenFunc func(context.Context, api.RequestDefaults) (Stream, error)

type Stream interface {
	// NewBatch returns a new batch sized appropriately
	// to the transport's capacity.
	NewBatch() Batch

	// Recv receives the next server-side event.
	// An [io.EOF] indicates the stream has terminated successfully.
	Recv() (Event, error)

	// Close the batch stream.
	Close() error
}

type Data struct {
	Object    *api.BatchObject
	Reference *api.Reference
}

type Event struct {
	Started      bool
	ShuttingDown bool
	Backoff      *int
	Acks         []string
	OOM          *OOM
	Results      *Results
}

type OOM struct {
	ReconnectAfter time.Duration
}

type Results struct {
	OK     []string
	Failed map[string]string
}

// Batch is a sized container, which accumulates batch items
// until their size reaches the parent transport's capacity.
type Batch interface {
	Add(Data) error

	// Send marshals a request with all items currently contained in the batch
	// and sends it using the parent stream object. See [Stream.NewBatch].
	Send() error
}

var (
	// ErrBatchFull is a sentinel error returned if [Batch] cannot fit
	// the latest item in the remaining space. The caller should send
	// the current batch as-is and retry the same item with a new batch.
	ErrBatchFull = errors.New("batch is full")

	// ErrTooLarge is a sentinel error returned if the latest item exceeds
	// the maximum request size supported by the [Batch]. This error should
	// not be retried, and surfaced to the user instead.
	ErrTooLarge = errors.New("batch item exceeds maximum request size")
)

func NewClient(ctx context.Context, openFunc OpenFunc) (*Client, error) {
	return &Client{
		ctx:  ctx,
		open: openFunc,
	}, nil
}

type Client struct {
	ctx  context.Context
	open OpenFunc
	b    Batch
}
