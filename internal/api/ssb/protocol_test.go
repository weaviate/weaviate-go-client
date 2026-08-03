package ssb_test

import (
	"context"
	"errors"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/ssb"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
)

type Simulation struct {
	*testing.T
	prng *testkit.PRNG

	TaskCount   int
	RetryLimit  int
	ReconnCount int
	ReconnLimit int
	TooLarge    sync.Map
}

// What I want to test:
//
//  1. Client never crashes.
//  2. If context is not canceled, all data added to the batch makes it to the server.
//     This includes reconnects, which can happen up to N times in a row.
//  3. OOM terminates the stream at any given time (+ timeout).
//  4. Tricky?: No new data is sent before an ACK.
//  5. Canceling the context fails all in-progress tasks.
//
// Of those, 1 and 2 seem to be something that's easy to model in a test.
// I think what I will try and do is first write tests where I come up with
// the scenario and then progressively add randomness into it.
//
// If I can build a server model with a good size-to-weight ratio, then I can
// just let it rip FRNG style and check that it doesn't produce invalid outcomes.
//
// Luckily (!!) the client can only do Add, so that is a trivial part really.
// I can fuzz the smaller components (like batch or cache) and PBT the client.
// protocol_fuzz_test.go (package ssb) + protocol_test.go (package ssb_test)
//
// Assertion guideline: `assert` for Client-stuff, `require` for test logic.
func TestClient(t *testing.T) {
	sim := Simulation{
		T:           t,
		prng:        testkit.NewPRNG(t),
		TaskCount:   32,
		RetryLimit:  3,
		ReconnLimit: 5,
	}

	conn := make(chan *Stream)
	t.Cleanup(func() { close(conn) })

	srv := Server{
		Simulation: &sim,
		conn:       conn,
		work:       make(chan *Batch, 8),
		seen:       make(map[uuid.UUID]taskStat, sim.TaskCount),
	}
	go srv.run()

	c := ssb.NewClient(ssb.ClientConfig{
		Context: t.Context(),
		Transport: &Transport{
			Simulation: &sim,
			conn:       conn,
		},
		QueueSize: 32,
		BatchSize: 16,
		Reconnect: ssb.ReconnectPolicy{
			Limit:     sim.ReconnLimit,
			DelayFunc: func(int) time.Duration { return 0 },
		},
		CanRetry: func(_ string, retries int, _ error) bool {
			return retries < sim.RetryLimit
		},
	})

	tasks := make([]*ssb.Task, sim.TaskCount)
	for i := 0; i < sim.TaskCount; i++ {
		id := uuid.New()
		if sim.prng.Chance(1, 50) {
			sim.TooLarge.Store(id, true)
		}
		task, err := c.Add(
			t.Context(),
			ssb.Data{Object: &api.BatchObject{UUID: id}},
		)
		assert.NoError(t, err, "add error")
		require.NotNil(t, task, "nil task")
		tasks[i] = task
	}

	t.Log("added all data -> close client")
	assert.NoError(t, c.Close(), "close client")
}

type Batch struct {
	*Simulation
	values []uuid.UUID
}
type Event ssb.Event

type Server struct {
	*Simulation
	conn <-chan *Stream
	work chan *Batch
	seen map[uuid.UUID]taskStat
}

type taskStat struct {
	Retries int
	Err     int
}

func (srv *Server) run() {
	srv.Helper()
	defer close(srv.work)

	// srv.Log("[server]: Run")

	for stream := range srv.conn {
		if err := stream.srvSend(Event{Started: true}); err != nil {
			continue
		}
		// srv.Log("[server]: Stream started")

		// var oom bool
	Conn:
		for {
			// if srv.prng.Chance(1, 30) || oom {
			// 	srv.Logf("[server]: Shutting down (OOM=%t)", oom)
			// 	stream.srvSend() <- Event{ShuttingDown: true}
			// 	break Conn
			// }
			select {
			case batch, ok := <-stream.srvRecv():
				if !ok {
					// srv.Logf("client closed the stream, process remaining %d batches", len(srv.work))
					for len(srv.work) > 0 {
						results := srv.processBatch(<-srv.work)
						if err := stream.srvSend(Event{Results: results}); err != nil {
							break Conn
						}
					}
					break Conn
				}
				require.NotNil(srv.T, batch)
				srv.Logf("[server]: Received batch (%d)", len(batch.values))
				require.NotEmpty(srv.T, batch.values)

				// if srv.prng.Chance(1, 20) {
				// 	srv.Logf("[server]: Out Of Memory")
				// 	oom = true
				// 	if srv.prng.Bool() {
				// 		stream.srvSend() <- Event{OOM: new(ssb.OOM)}
				// 	} else {
				// 		stream.srvSend() <- Event{OOM: &ssb.OOM{ExitAfter: 5 * time.Second}}
				// 	}
				// 	continue
				// }

				srv.work <- batch
				srv.Log("[server]: Ack batch")
				if err := stream.srvSend(Event{Ack: true}); err != nil {
					srv.Log("[server]: Failed to ack batch")
					break Conn
				}
				srv.Log("[server]: Ack-ed batch")

			case batch := <-srv.work:
				if srv.prng.Chance(1, 10) {
					srv.Log("[server]: Busy, return batch to queue")
					srv.work <- batch
					continue
				}
				results := srv.processBatch(batch)
				if err := stream.srvSend(Event{Results: results}); err != nil {
					break Conn
				}
			}
		}

		srv.Log("[server]: Close stream + discard any remaining work")
		for len(srv.work) > 0 {
			<-srv.work
		}
		stream.srvClose()
	}
}

func (srv *Server) processBatch(b *Batch) *ssb.Results {
	results := &ssb.Results{
		OK:     make([]string, 0),
		Failed: make(map[string]error),
	}
	var OK, failed int
	for _, id := range b.values {
		t, ok := srv.seen[id]
		if !ok {
			t = taskStat{Retries: -1}
		}
		t.Retries++

		// On it's last retry the task fails with a 1/4 chance.
		if id := id.String(); srv.prng.Chance(1, srv.RetryLimit-t.Retries+4) {
			results.Failed[id] = testkit.ErrWhaam
			failed++
		} else {
			results.OK = append(results.OK, id)
			OK++
		}
		srv.seen[id] = t
	}
	srv.Logf("[server]: Send results: ok=%d, failed=%d", OK, failed)
	return results
}

type Transport struct {
	*Simulation
	conn chan<- *Stream
}

var _ ssb.Transport = (*Transport)(nil)

func (t *Transport) NewStream(ctx context.Context) (ssb.Stream, error) {
	t.Helper()
	assert.NotNil(t.T, ctx, "nil stream context")

	t.Log("[transport]: New stream")

	ctx, cancel := context.WithCancelCause(ctx)
	s := &Stream{
		Simulation: t.Simulation,
		ctx:        ctx,
		cancel:     cancel,
		batchc:     make(chan *Batch),
		eventc:     make(chan Event),
	}

	t.conn <- s
	return s, nil // TODO(dyma): sometimes error for reconnects?
}

type Stream struct {
	*Simulation
	ctx    context.Context
	cancel context.CancelCauseFunc
	batchc chan *Batch
	eventc chan Event
}

func (s *Stream) Send(req any) error {
	s.Helper()
	require.NotNil(s.T, req, "nil request")
	require.IsType(s.T, (*Batch)(nil), req, "bad request")

	if s.prng.Chance(1, 25) {
		s.cancel(errors.New("bad network"))
	}

	select {
	case s.batchc <- req.(*Batch):
	case <-s.ctx.Done():
		s.Log("[stream.Send]: context canceled")
		return io.EOF
	}
	return nil
}

func (s *Stream) Recv() (ssb.Event, error) {
	select {
	case ev := <-s.eventc:
		return ssb.Event(ev), nil
	case <-s.ctx.Done():
		// s.Logf("[stream.Recv]:  context canceled: %v", context.Cause(s.ctx))
		return ssb.Event{}, context.Cause(s.ctx)
	}
}

func (s *Stream) Close() error {
	s.Log("[stream.Close]: Client closed it's half of the stream")
	close(s.batchc)
	return nil
}

func (s *Stream) srvSend(ev Event) error {
	// s.Log("[server::srvSend] start")
	// defer s.Log("[server::srvSend] done")
	select {
	case s.eventc <- ev:
		return nil
	case <-s.ctx.Done():
		return s.ctx.Err()
	}
}
func (s *Stream) srvRecv() <-chan *Batch { return s.batchc }
func (s *Stream) srvClose()              { s.cancel(io.EOF) }

func (t *Transport) NewRequest() ssb.BatchRequest {
	b := &Batch{
		Simulation: t.Simulation,
		values:     make([]uuid.UUID, 0, 92), // FIXME: random cap, not 92
	}
	return b
}

func (t *Transport) Prepare(data ssb.Data) (any, error) {
	t.Helper()
	require.NotNil(t.T, data.Object, "nil batch object")

	id := data.Object.UUID
	if _, ok := t.TooLarge.Load(id); ok {
		t.Log("TOO LARGE ", id)
		return nil, ssb.ErrTooLarge
	}
	return id, nil
}

func (b *Batch) Add(v any) (added, full bool) {
	b.Helper()
	defer require.LessOrEqual(b.T, len(b.values), cap(b.values))

	require.IsType(b.T, *new(uuid.UUID), v, "bad value in Add")
	if len(b.values) == cap(b.values) {
		return false, true
	}
	b.values = append(b.values, v.(uuid.UUID))
	return true, len(b.values) == cap(b.values)
}
