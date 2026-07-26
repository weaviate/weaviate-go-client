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
	Context   context.Context
	Transport Transport
	QueueSize int       // Queue buffer size.
	BatchSize int       // Initial batch capacity.
	RetryFunc RetryFunc // Controls if a task will be retried.
	Reconnect ReconnectPolicy
}

type ReconnectPolicy struct {
	DelayFunc func(n int) time.Duration // Returns delay before the next re-connect.
	Limit     int                       // Maximum number of re-connect attempts.
}

func ExponentialDelay(n int) time.Duration { return 1 << n * time.Second }

func NewClient(conf ClientConfig) *Client {
	ctx, cancel := context.WithCancelCause(conf.Context)
	c := &Client{
		ctx:         ctx,
		finish:      cancel,
		transport:   conf.Transport,
		queue:       make(chan *Task, conf.QueueSize),
		retry:       make(chan []*Task),
		batch:       newBatch(conf.Transport.NewRequest, conf.BatchSize),
		wip:         newCache(conf.BatchSize),
		state:       newState(canPrepare),
		canRetry:    conf.RetryFunc,
		delayFunc:   conf.Reconnect.DelayFunc,
		reconnLimit: conf.Reconnect.Limit,
	}
	c.init()
	return c
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
	NewStream(context.Context) (Stream, error)

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
	ExitAfter time.Duration
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

// clear all flags set previously.
func (s *state) clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.flags = 0
}

// set state flags and notify the awaiting goroutine.
func (s *state) set(af actionFlags) {
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

// ErrTooLarge is a sentinel error returned if the latest item exceeds
// the maximum request size supported by the [Batch].
// This error is not retried, and surfaced to the user instead.
var ErrTooLarge = errors.New("batch item exceeds maximum request size")

func (t *Task) setValue(v any) { t.val.CompareAndSwap(nil, v) }
func (t *Task) value() any     { return t.val.Load() }

// retry sets the error and increments retry count.
func (t *Task) retry(err error) {
	t.err.Store(err)
	t.retries.Add(1)
}

// complete sets the error and closes the done channel.
func (t *Task) complete(err error) {
	t.err.Store(err)
	close(t.done)
}

type Client struct {
	// Parent context. It is derived from [ClientConfig.Context]
	// and is used internally to halt the client in the event of an error.
	ctx    context.Context
	finish context.CancelCauseFunc // Cancels client context.

	transport Transport    // Transport provides Stream and BatchRequest.
	state     *state       // Stream state.
	queue     chan *Task   // Task queue.
	retry     chan []*Task // Retry queue.
	batch     *batch       // Batch container.
	wip       *cache       // Tasks taken off the queue but not yet completed.
	canRetry  RetryFunc    // Retry decides if a task will be retried.

	// Tally of failed reconnect attempts. It is only accessed
	// by the 'recv' goroutine, so it does not need a guard.
	reconn      int
	reconnLimit int                     // Maximum number of failed reconnects.
	delayFunc   func(int) time.Duration // Delay before the next reconnect.
}

func (c *Client) init() {
	type span struct {
		ctx    context.Context
		stream Stream
	}

	tick := make(chan span)
	go func() {
		for s := range tick {
			c.send(s.ctx, s.stream)
		}
	}()

	go func() {
		var err error
		defer c.finish(err)
		defer close(tick)

		for ; c.reconn < c.reconnLimit; c.reconn++ {
			var s Stream
			if s, err = c.transport.NewStream(c.ctx); err == nil {
				ctx, cancel := context.WithCancel(c.ctx)
				tick <- span{ctx: ctx, stream: s}
				if err = c.recv(s, cancel); err == io.EOF {
					return
				}
				<-ctx.Done()
			}

			c.state.clear()

			select {
			case <-time.After(c.delayFunc(c.reconn)):
				c.batch.clear()
				c.wip.all(func(t *Task) { c.batch.add(t.value()) })
			case <-c.ctx.Done():
				return
			}
		}
	}()
}

func (c *Client) send(ctx context.Context, s Stream) {
	defer s.Close() //nolint:errcheck

	// maybeSend returns an error if Stream.Send fails
	// or the context expires while waiting for canSend.
	maybeSend := func() (err error) {
		req := c.batch.prepare()
		if req != nil {
			if err = c.state.await(ctx, canSend); err != nil {
				return
			}
			if err = s.Send(req); err != nil {
				return
			}
			c.state.clear()
		}
		return
	}

	if err := c.state.await(ctx, canPrepare); err != nil {
		return
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

		case tasks := <-c.retry:
			for _, t := range tasks {
				c.batch.add(t.value())
			}

		case <-ctx.Done():
			return
		}

		if err := maybeSend(); err != nil {
			return
		}
	}

Drain:
	c.batch.disableGrowth()
	for {
		wipCount := c.wip.size()
		if wipCount == 0 {
			return
		}
		c.batch.resize(wipCount)

		if err := maybeSend(); err != nil {
			return
		}

		select {
		case tasks := <-c.retry:
			// Every time we receive Results, wip is updated.
			// Resizing the batch to the current wip size guarantees
			// that the former will eventually fill up.
			for _, t := range tasks {
				c.batch.add(t.value())
			}
		case <-ctx.Done():
			return
		}
	}
}

func (c *Client) recv(s Stream, cancelSend context.CancelFunc) error {
	defer cancelSend()

	var oomTimer *time.Timer // Waiting for OOM to resolve.
	var shutdown bool        // Server is shutting down.

	for {
		event, err := s.Recv()
		if err != nil {
			if err == io.EOF && shutdown {
				// The client should reconnect.
				err = errors.New("server shutdown")
			}
			return err
		}

		switch {
		case event.Started:
			c.reconn = 0
			c.state.set(canPrepare | canSend)

		case event.Acks != nil:
			// NOTE(dyma): the protocol guarantees that Acks message
			// includes all data from the previous batch. Do we need
			// to verify that that is the case?
			c.state.set(canPrepare | canSend)

		case event.Results != nil:
			c.wip.walk(slices.Values(event.Results.OK), func(t *Task) bool {
				t.complete(nil)
				return true
			})

			failed := event.Results.Failed
			if len(failed) == 0 {
				continue
			}

			// Adding tasks to c.retry may block, so we kick off a goroutine.
			// Its lifetime is bounded by [Client.ctx].
			go func() {
				// NOTE(dyma): we could share []*Task slice between send and recv.
				// This requires another synchronization chan, maybe not worth it.
				retry := make([]*Task, 0, len(failed))
				c.wip.walk(maps.Keys(failed), func(t *Task) (remove bool) {
					err := errors.New(failed[t.ID()])
					if c.canRetry.check(t, err) {
						t.retry(err)
						retry = append(retry, t)
					} else {
						t.complete(err)
						remove = true
					}
					return
				})

				select {
				case c.retry <- retry:
				case <-c.ctx.Done():
					return
				}
			}()

		case event.Backoff != nil:
			c.batch.resize(*event.Backoff)

		case event.OOM != nil:
			oomTimer = time.AfterFunc(event.OOM.ExitAfter, func() {
				c.finish(errors.New("server OOM"))
			})
			c.state.clear()

		case event.ShuttingDown:
			cancelSend()
			if oomTimer != nil {
				// If the timer has already fired, the context would be canceled
				// and s.Recv() must return an error. The only reason that we got
				// here is that the message arrived _just_ before oomTimer fired.
				// This means the server is responsive and we can try to reconnect.
				oomTimer.Stop()
				oomTimer = nil
			}
			shutdown = true
			c.state.clear()
		}
	}
}

func (c *Client) Close() error {
	defer close(c.retry)

	close(c.queue)
	<-c.ctx.Done()
	err := context.Cause(c.ctx)
	if err == io.EOF {
		return nil
	}

	c.wip.all(func(t *Task) { t.complete(err) })
	for t := range c.queue {
		t.complete(err)
	}
	return err
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

// addLocked adds v to request and updates [batch.len].
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

// clear empties the batch, preserving the capacity and noGrow flag, if set.
func (b *batch) clear() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.req = b.newRequest()
	b.buf = b.buf[:0]
	b.len = 0
	b.flags ^= full
}

func newCache(size int) *cache {
	return &cache{m: make(map[string]*Task, size)}
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

// all calls f for all tasks in the cache.
func (c *cache) all(f func(*Task)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, t := range c.m {
		f(t)
	}
}
