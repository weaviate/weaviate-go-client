package ssb

import (
	"context"
	"errors"
	"io"
	"iter"
	"maps"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
)

type ClientConfig struct {
	api.RequestDefaults
	Context   context.Context
	Transport Transport
	QueueSize int // Queue buffer size.
	BatchSize int // Initial batch capacity.
	RetryFunc RetryFunc
}

func NewClient(conf ClientConfig) (*Client, error) {
	ctx, cancel := context.WithCancelCause(conf.Context)
	c := &Client{
		ctx:       ctx,
		cancel:    cancel,
		defaults:  conf.RequestDefaults,
		transport: conf.Transport,
		queue:     make(chan *Task, conf.QueueSize),
		retry:     make(chan []*Task, conf.BatchSize),
		wip: &cache{
			m: make(map[string]*Task, conf.BatchSize),
		},
		batch: newBatch(conf.Transport.NewRequest, conf.BatchSize),
		state: newState(notStarted),
	}
	if err := c.init(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) Add(d Data) (*Task, error) {
	t := &Task{
		data: d,
		done: make(chan struct{}),
	}
	select {
	case c.queue <- t:
		return t, nil
	case <-c.ctx.Done():
		// TODO(dyma): should Add have its own context? I feel like yes.
		return nil, context.Canceled
	}
}

type Task struct {
	data Data

	val     atomic.Value
	retries atomic.Uint32
	err     atomic.Value
	done    chan struct{}
}

func (t *Task) ID() string            { return t.data.ID() }
func (t *Task) Done() <-chan struct{} { return t.done }
func (t *Task) Err() error            { return t.err.Load().(error) }
func (t *Task) TimesRetried() int     { return int(t.retries.Load()) }

type Transport interface {
	// NewStream opens a new streaming client and starts the batch stream.
	NewStream(context.Context, api.RequestDefaults) (Stream, error)

	// NewRequest returns a new batch request
	// sized appropriately to the transport's capacity.
	NewRequest() BatchRequest

	// Prepare marshals the data to a value accepted by [BatchRequest.Add].
	Prepare(Data) (any, error)
}

type Stream interface {
	// Send marshaled batch request. The batch is guaranteed to be
	// of the same type as returned by [Transport.NewRequest].
	Send(BatchRequest) error

	// Recv receives the next server-side event.
	// An [io.EOF] indicates the stream has terminated successfully.
	Recv() (Event, error)

	// Close the batch stream.
	Close() error
}

// BatchRequest is a sized container, which accumulates batch items
// until their size reaches the parent stream's capacity.
type BatchRequest interface {
	Add(v any) (added, full bool)
}

// RetryFunc determines if a task should be retried after the following error.
type RetryFunc func(id string, retries int, err error) bool

// check is safe to call on a nil RetryFunc.
func (rf RetryFunc) check(t *Task, err error) bool {
	return rf != nil && rf(t.ID(), t.TimesRetried(), err)
}

type Data struct {
	Object    *api.BatchObject
	Reference *api.Reference
}

func (d *Data) ID() (id string) {
	switch {
	case d.Object != nil:
		id = d.Object.UUID.String()
	case d.Reference != nil:
		id = d.Reference.String()
	default:
		dev.Unreachable()
	}
	return
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

type actionFlags uint8

const (
	canPrepare actionFlags = 1 << iota
	canSend
)

const (
	inFlight actionFlags = 0
	reconnecting

	notStarted = canPrepare
	shuttingDown

	active = canPrepare | canSend
)

func newState(actions actionFlags) *state {
	return &state{
		flags:   actions,
		changed: make(chan struct{}),
	}
}

type state struct {
	mu      sync.RWMutex
	flags   actionFlags
	changed chan struct{}
}

// set state flags to a new value.
func (s *state) set(af actionFlags) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.flags = af
}

// set state flags and notify the awaiting goroutine.
func (s *state) setNotify(af actionFlags) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.flags = af
	close(s.changed)
	s.changed = make(chan struct{})
}

// await blocks until flag is set or ctx expires
// and returns [Context.Err] in the latter case.
func (s *state) await(ctx context.Context, af actionFlags) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for s.flags&af != af {
		changed := s.changed
		s.mu.RUnlock()

		select {
		case <-changed:
			s.mu.RLock()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

var (
	// ErrBatchFull is a sentinel error returned if [Batch] cannot fit
	// the latest item in the remaining space. The caller should send
	// the current batch as-is and retry the same item with a new batch.
	ErrBatchFull = errors.New("batch is full")

	// ErrTooLarge is a sentinel error returned if the latest item exceeds
	// the maximum request size supported by the [Batch].
	// This error is not retried, and surfaced to the user instead.
	ErrTooLarge = errors.New("batch item exceeds maximum request size")
)

func (t *Task) incrTimesRetried()  { t.retries.Add(1) }
func (t *Task) setValue(v any) any { return t.val.CompareAndSwap(nil, v) }
func (t *Task) value() any         { return t.val.Load() }
func (t *Task) setErr(err error)   { t.err.Store(err) }
func (t *Task) complete(err error) {
	t.setErr(err)
	close(t.done)
}

type Client struct {
	// Parent context. It is derived from [ClientConfig.Context]
	// and is used internally to halt the client in the event of an error.
	ctx    context.Context
	cancel context.CancelCauseFunc // Cancels client context.

	defaults  api.RequestDefaults // Defaults for all outgoing requests.
	transport Transport           // Transport provides Stream and BatchRequest.
	queue     chan *Task          // Task queue.
	retry     chan []*Task        // Retry queue.
	batch     *batch              // Batch container.
	state     *state              // State controls
	wip       *cache              // Tasks taken off the queue but not yet completed.
	canRetry  RetryFunc           // Retry decides if a task will be retried.

	// FIXME(dyma): this needs to be protected by mutex
	cancelSend context.CancelCauseFunc // Cancels the context of the "send" goroutine.
}

func (c *Client) init() error {
	s, err := c.transport.NewStream(c.ctx, c.defaults)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancelCause(c.ctx)
	c.cancelSend = cancel
	go c.send(ctx, s)
	go c.recv(s)

	return nil
}

func (c *Client) send(ctx context.Context, s Stream) {
	var fatal error

	maybeSend := func() {
		req := c.batch.prepare()
		if req != nil {
			if fatal = c.state.await(ctx, canSend); fatal != nil {
				return // essentially ctx.Done()
			}
			if fatal = s.Send(req); fatal != nil {
				return // TODO(dyma): handle connection error
			}
			c.state.set(inFlight)
		}
		return
	}

	if fatal = c.state.await(ctx, canPrepare); fatal != nil {
		goto Exit
	}

	for {
		select {
		case t, ok := <-c.queue:
			if !ok {
				goto Drain
			}

			v, err := c.transport.Prepare(t.data)
			if err != nil {
				t.complete(err)
			}

			t.setValue(v)
			c.wip.put(t)
			c.batch.add(v)

		case tasks, _ := <-c.retry:
			for _, t := range tasks {
				c.batch.add(t.value())
			}

		case <-ctx.Done():
			goto Exit
		}

		if maybeSend(); fatal != nil {
			goto Exit
		}
	}

Drain:
	c.batch.disableGrowth()
	c.batch.resize(c.wip.size())
	for {
		if maybeSend(); fatal != nil {
			goto Exit
		}

		select {
		case tasks, _ := <-c.retry:
			// Every time we receive Results, the cache is updated.
			// Resizing the batch to the current cache size guarantees
			// that it will eventually fill up.
			// FIXME(dyma): ensure Backoff messages do not resize the cache back up.
			c.batch.resize(c.wip.size())
			for _, t := range tasks {
				c.batch.add(t.value())
			}
		case <-ctx.Done():
			goto Exit
		}
	}

Exit:
	if errors.Is(fatal, context.Cause(ctx)) {
		// The context got canceled, we just need to exit.
		return
	}
	// TODO(dyma): Send failed. Close the stream?
	print("bye")
}

func (c *Client) recv(s Stream) {
	for {
		event, err := s.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
		}

		switch {
		case event.Started:
			c.state.setNotify(active)

		case event.Acks != nil:
			// NOTE(dyma): the protocol guarantees that Acks message
			// includes all data from the previous batch. Do we need
			// to verify that that is the case?
			c.state.setNotify(active)

		case event.Results != nil:
			c.wip.walk(slices.Values(event.Results.OK), func(t *Task) bool {
				t.complete(nil)
				return true
			})

			// NOTE(dyma): we _could_ re-use the array between send/recv
			// but that requires another synchronization channel. Later.
			failed := event.Results.Failed
			retry := make([]*Task, 0, len(failed))
			c.wip.walk(maps.Keys(failed), func(t *Task) (remove bool) {
				err := errors.New(failed[t.ID()])
				if c.canRetry.check(t, err) {
					t.setErr(err)
					retry = append(retry, t)
				} else {
					t.complete(err)
				}
				return
			})

			// Adding tasks to c.retry may block, so we kick off a goroutine.
			go func() {
				select {
				case c.retry <- retry:
					for _, t := range retry {
						t.incrTimesRetried()
					}
				case <-c.ctx.Done():
					// TODO(dyma): probably we should fail the entries?
					return
				}
			}()

		case event.Backoff != nil:
			c.batch.resize(*event.Backoff)

		// TODO(dyma): handle
		case event.ShuttingDown:
		case event.OOM != nil:
		}
	}
}

func (c *Client) Close() error {
	if c.cancelSend != nil {
		c.cancelSend(nil)
	}
	return c.ctx.Err()
}

func newBatch(newRequest func() BatchRequest, size int) *batch {
	return &batch{
		newRequest: newRequest,
		req:        newRequest(),
		cap:        size,
		buf:        make([]any, 0, size),
	}
}

// batchFlags store internal batch state.
type batchFlags uint8

const (
	full   batchFlags = 1 << iota // Batch is full.
	noGrow                        // Batch capacity may not increase.
)

type batch struct {
	mu sync.Mutex

	flags batchFlags
	buf   []any
	req   BatchRequest
	cap   int
	len   int

	newRequest func() BatchRequest
}

func (b *batch) add(v any) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.buf = append(b.buf, v)
	if b.flags&full != full {
		b.addLocked(v)
	}
}

// addLocked adds v to request and updates [batch.len] and [batch.full].
func (b *batch) addLocked(v any) {
	added, reqFull := b.req.Add(v)
	if added {
		b.len++
	}
	if reqFull || b.len == b.cap {
		b.flags |= full
	}
}

// prepare returns the request object and refills the batch from [batch.buf].
// If the batch is full, prepare is a no-op and returns nil.
func (b *batch) prepare() BatchRequest {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.flags&full != full {
		return nil
	}

	req := b.req          // Copy request data.
	b.buf = b.buf[b.len:] // Discard data that's already in the request.
	b.refillLocked()
	return req
}

// resize sets a new capacity and resizes the request accordingly.
// If [noGrow] flag is set, resize will not increase the capacity.
func (b *batch) resize(size int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if size > b.cap && b.flags&noGrow == noGrow {
		return
	}

	b.cap = size
	if b.len == b.cap {
		b.flags |= full
	} else if b.len > b.cap {
		b.refillLocked()
	}
}

// disableGrowth sets [noGrow] flag.
func (b *batch) disableGrowth() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.flags |= noGrow
}

// refillLocked creates a new request and re-populates it
// from [batch.buf] up to the current capacity.
func (b *batch) refillLocked() {
	b.req = b.newRequest()
	b.len = 0
	b.flags ^= full
	for _, v := range b.buf {
		b.addLocked(v)
		if b.flags&full == full {
			break
		}
	}
	b.buf = b.buf[b.len:]
}

// cache is a synchronized map of in-progress tasks.
type cache struct {
	mu sync.Mutex
	m  map[string]*Task
}

// put a task in the cache.
func (c *cache) put(t *Task) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[t.ID()] = t
}

// walk calls f for every cache entry in keys.
// If f returns true, the entry is removed.
func (c *cache) walk(keys iter.Seq[string], f func(*Task) bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k := range keys {
		t, ok := c.m[k]
		if ok && f(t) {
			delete(c.m, k)
		}
	}
}

// size returns the number of entries in the cache.
func (c *cache) size() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.m)
}
