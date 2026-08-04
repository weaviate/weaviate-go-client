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
	QueueSize int             // Queue buffer size.
	BatchSize int             // Initial batch capacity.
	Reconnect ReconnectPolicy // Policy for dealing with dropped connections.
	CanRetry  CanRetryFunc    // Determines if a task will be retried.
}

type ReconnectPolicy struct {
	DelayFunc func(n int) time.Duration // Returns delay before the next re-connect.
	Limit     int                       // Maximum number of re-connect attempts.
}

func ExponentialDelay(n int) time.Duration { return 1 << n * time.Second }

// Maximum number of Result events the client will accept
// before it starts putting backpressure on the Recv.
const retryBuffer = 1 << 4

func NewClient(conf ClientConfig) *Client {
	ctx, cancel := context.WithCancelCause(conf.Context)
	c := &Client{
		ctx:         ctx,
		finish:      cancel,
		transport:   conf.Transport,
		queue:       make(chan *Task, conf.QueueSize),
		retry:       make(chan []*Task, retryBuffer),
		batch:       newBatch(conf.Transport.NewRequest, conf.BatchSize),
		wip:         newCache(conf.BatchSize),
		state:       newState(canPrepare),
		canRetry:    conf.CanRetry,
		delayFunc:   conf.Reconnect.DelayFunc,
		reconnLimit: conf.Reconnect.Limit,
	}
	c.init()
	return c
}

// Add creates a new task and puts it on the work queue.
// If [Client.Context] expires, Add returns [context.Canceled],
// otherwise the error is nil. Calling Add after closing the batch panics.
func (c *Client) Add(ctx context.Context, d Data) (*Task, error) {
	t := &Task{
		data: d,
		done: make(chan struct{}),
	}
	select {
	case c.queue <- t:
		return t, nil
	case <-ctx.Done():
	case <-c.ctx.Done():
	}
	return nil, context.Canceled
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
func (t *Task) Err() error {
	if err := t.err.Load(); err != nil {
		return err.(error)
	}
	return nil
}
func (t *Task) TimesRetried() int { return int(t.retries.Load()) }

type Transport interface {
	// NewStream opens a new streaming client and starts the batch stream.
	NewStream(context.Context) (Stream, error)

	// NewRequest returns a new batch request
	// sized appropriately to the transport's capacity.
	NewRequest() BatchRequest

	// Prepare marshals the data to a value accepted by [BatchRequest.Add].
	Prepare(Data) (any, error)
}

// ErrTooLarge is an error [Transport] must return if the latest item
// exceeds the maximum request size supported by the transport.
// This error is not retried, and surfaced to the user instead.
var ErrTooLarge = errors.New("batch item exceeds maximum request size")

type Stream interface {
	// Send marshaled batch request. The batch must be
	// the value returned by [Transport.NewRequest].
	Send(any) error

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

// CanRetryFunc determines if a task should be retried after the following error.
type CanRetryFunc func(id string, retries int, err error) bool

// check is safe to call on a nil RetryFunc.
func (can CanRetryFunc) check(t *Task, err error) bool {
	return can != nil && can(t.ID(), t.TimesRetried(), err)
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

// Event is a union of expected server-side messages.
type Event struct {
	Started      bool     // Server is ready to accept data.
	ShuttingDown bool     // Server is shutting down.
	Backoff      *int     // Limit the number of objects in the future batch.
	Ack          bool     // The previous batch is ack'ed.
	OOM          *OOM     // Server is OOM and will be shutting down.
	Results      *Results // Results for a previously-acked batch.
}

type OOM struct {
	ExitAfter time.Duration // Grace period for the server to send a ShuttingDown event.
}

type Results struct {
	OK     []string         // IDs of tasks that completed successfully.
	Failed map[string]error // Batch insertion errors keyed by task ID.
}

type permissionFlags uint8

const (
	canPrepare permissionFlags = 1 << iota // Client can prepare the next batch.
	canSend                                // Client can send the next batch.
)

func newState(permissions permissionFlags) state {
	return state{
		permissions: permissions,
		changed:     make(chan struct{}),
	}
}

// state stores actions the client is currently allowed to take.
type state struct {
	mu          sync.Mutex
	permissions permissionFlags
	changed     chan struct{}
}

// clear all flags set previously.
func (s *state) clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.permissions = 0
}

// set state flags and notify the awaiting goroutine.
func (s *state) set(permissions permissionFlags) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.permissions = permissions
	close(s.changed)
	s.changed = make(chan struct{})
}

// await blocks until permission is set and atomically consumes it.
// If ctx expires await exits early with [Context.Err].
func (s *state) await(ctx context.Context, permissions permissionFlags) error {
	// await works like a context-aware [sync.Cond.Wait].
	// We check if [s.permissions] contain the permissions we need
	// while holding the lock. If the check fails, we release the lock
	// and wait until new permissions are set.
	s.mu.Lock()
	for s.permissions&permissions != permissions {
		// s.changed is closed and re-created on set, so we must copy
		// the current channel to avoid accessing s.changed concurrently.
		changed := s.changed
		s.mu.Unlock()

		select {
		case <-changed:
			s.mu.Lock()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	s.permissions &^= permissions
	s.mu.Unlock()
	return nil
}

func (t *Task) setValue(v any) { t.val.CompareAndSwap(nil, v) }
func (t *Task) value() any     { return t.val.Load() }

// retry sets the error and increments retry count.
func (t *Task) retry(err error) {
	t.err.Store(err)
	t.retries.Add(1)
}

// complete sets the error and closes the done channel.
func (t *Task) complete(err error) {
	if err != nil {
		t.err.Store(err)
	}
	close(t.done)
}

type Client struct {
	// Parent context. It is derived from [ClientConfig.Context]
	// and is used internally to halt the client in the event of an error.
	ctx    context.Context
	finish context.CancelCauseFunc // Cancels client context.

	transport Transport    // Transport provides Stream and BatchRequest.
	state     state        // Stream state.
	queue     chan *Task   // Task queue.
	retry     chan []*Task // Retry queue.
	batch     batch        // Batch container.
	wip       cache        // Tasks taken off the queue but not yet completed.
	canRetry  CanRetryFunc // Retry decides if a task will be retried.

	// Tally of failed reconnect attempts. It is only accessed
	// by the 'recv' goroutine, so it does not need a guard.
	reconnCount int                     // Count of consecutive reconnects.
	reconnLimit int                     // Maximum number of failed reconnects.
	delayFunc   func(int) time.Duration // Delay before the next reconnect.
}

// init kicks off the 'send' and 'recv' routines.
func (c *Client) init() {
	type span struct {
		ctx    context.Context
		stream Stream
	}

	tick := make(chan span)
	tock := make(chan struct{})
	go func() {
		defer close(tock)
		for s := range tick {
			c.send(s.ctx, s.stream)
			tock <- struct{}{}
		}
	}()

	go func() {
		var err error
		defer func() {
			c.finish(err)
			close(tick)
			close(c.retry)
		}()

		// Start a stream, unblock the 'send' goroutine, and continue reconnecting
		// until batch completes, the server is deemed unresponsive, or context expires.
		for err = io.EOF; c.reconnCount < c.reconnLimit; c.reconnCount++ {
			var s Stream
			if s, err = c.transport.NewStream(c.ctx); err == nil {
				ctx, cancel := context.WithCancel(c.ctx)
				tick <- span{ctx: ctx, stream: s}
				err = c.recv(s, cancel)
				<-tock       // Wait until current 'send' span completes.
				<-ctx.Done() // Make sure recv canceled the context on exit
				if err == io.EOF {
					// io.EOF means the stream ended successfully,
					// and everything else means "try to reconnect".
					// IF the server OOMs and stops responding, the
					// canceled c.ctx will prevent us from reconnecting.
					return
				}
			}

			c.state.clear()
			// If connection drops, we assume that any in-progress tasks
			// have failed on the server and we have to redo them all.
			c.batch.clear()
			c.wip.all(func(t *Task) { c.batch.add(t) })
			c.retry = make(chan []*Task, retryBuffer)

			select {
			case <-time.After(c.delayFunc(c.reconnCount)):
			case <-c.ctx.Done():
				return
			}
		}
	}()
}

// send consumes queue and retry channels, prepares and sends the batch.
func (c *Client) send(ctx context.Context, s Stream) {
	defer s.Close() //nolint:errcheck

	// The caller should exit if maybeSend returns a non-nil error.
	maybeSend := func() (err error) {
		req := c.batch.prepare()
		if req != nil {
			if err = c.state.await(ctx, canSend); err != nil {
				return
			}
			if err = s.Send(req); err != nil {
				return
			}
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
				continue
			}

			t.setValue(v)
			c.wip.put(t)
			c.batch.add(t)

		case tasks := <-c.retry:
			c.batch.add(tasks...)

		case <-ctx.Done():
			return
		}

		if err := maybeSend(); err != nil {
			return
		}
	}

Drain:
	// Every time we receive Results, the wip shrinks and the batch's capacity
	// is reduced to the wip's size. This guarantees the batch will eventually
	// fill up. Disable growth to prevent Backoff from increasing the capacity.
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
			c.batch.add(tasks...)
		case <-ctx.Done():
			return
		}
	}
}

func (c *Client) recv(s Stream, cancelSend context.CancelFunc) error {
	defer cancelSend()

	var oomTimer *time.Timer // Waiting for OOM to resolve.
	var shutdown bool        // Server is shutting down.

	defer func() {
		// Stop the timer unconditionally, even if the stream hangs up
		// while waiting for ShuttingDown after seeing an OOM.
		if oomTimer != nil {
			oomTimer.Stop()
		}
	}()

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
			c.reconnCount = 0
			c.state.set(canPrepare | canSend)

		case event.Ack:
			c.state.set(canPrepare | canSend)

		case event.Results != nil:
			c.wip.walk(slices.Values(event.Results.OK), func(t *Task) bool {
				t.complete(nil)
				return true
			})

			// NOTE(dyma): we could share []*Task slice between send and recv.
			// This requires another synchronization chan, maybe not worth it.
			failed := event.Results.Failed
			retry := make([]*Task, 0, len(failed))
			if len(failed) > 0 {
				c.wip.walk(maps.Keys(failed), func(t *Task) (remove bool) {
					err := failed[t.ID()]
					if c.canRetry.check(t, err) {
						t.retry(err)
						retry = append(retry, t)
					} else {
						t.complete(err)
						remove = true
					}
					return
				})
			}

			select {
			case c.retry <- retry:
			case <-c.ctx.Done():
				continue
			}

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
				// The only reason that we got here is because the message
				// had arrived before oomTimer fired. Server is responsive
				// and we can try to reconnect.
				oomTimer.Stop()
				oomTimer = nil
			}
			shutdown = true
			c.state.clear()
		}
	}
}

// Close closes the work queue, drains the batch, and blocks until either
// all accepted tasks have been processed. If the context expires, Close
// fails all pending tasks and returns the cause of the context's expiry.
func (c *Client) Close() error {
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

func newBatch(newRequest func() BatchRequest, size int) batch {
	return batch{
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

// add all [Task.value] to the batch.
func (b *batch) add(tasks ...*Task) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, t := range tasks {
		v := t.value()
		b.buf = append(b.buf, v)
		if b.flags&full != full {
			b.addLocked(v)
		}
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
	b.flags &^= full
	for _, v := range b.buf {
		b.addLocked(v)
		if b.flags&full == full {
			break
		}
	}
}

// clear empties the batch, preserving the capacity and noGrow flag, if set.
func (b *batch) clear() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.req = b.newRequest()
	b.len = 0
	b.buf = b.buf[:0]
	b.flags &^= full
}

func newCache(size int) cache {
	return cache{m: make(map[string]*Task, size)}
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
