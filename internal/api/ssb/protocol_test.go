package ssb_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"slices"
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
	BatchSize   int
	RetryLimit  int
	ReconnCount int
	ReconnLimit int
	TooLarge    sync.Map
	CanOOM      bool
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
// Assertion guideline: During the simulation, use `require` only for testing logic,
// as it may interrupt the test and cause a panic. Use `require` for final checks.
func TestClient(t *testing.T) {
	sim := Simulation{
		T:           t,
		prng:        testkit.NewPRNG(t),
		TaskCount:   124512,
		BatchSize:   53,
		RetryLimit:  3,
		ReconnLimit: 5,
		CanOOM:      true,
	}

	conn := make(chan *Stream)
	srv := Server{
		Simulation: &sim,
		conn:       conn,
		work:       make(chan *Batch, 8),
		seen:       make(map[string]taskStat, sim.TaskCount),
		done:       make(chan struct{}),
	}
	go srv.run()

	c := ssb.NewClient(ssb.ClientConfig{
		Context: t.Context(),
		Transport: &Transport{
			Simulation: &sim,
			conn:       conn,
		},
		QueueSize: 32,
		BatchSize: sim.BatchSize,
		Reconnect: ssb.ReconnectPolicy{
			Limit:     sim.ReconnLimit,
			DelayFunc: func(int) time.Duration { return 0 },
		},
		CanRetry: func(id string, retries int, _ error) bool {
			return retries < sim.RetryLimit
		},
	})

	// Add all data to the client.
	tasks := make([]*ssb.Task, 0, sim.TaskCount)
	for i := 0; i < sim.TaskCount; i++ {
		id := uuid.New()
		if sim.prng.Chance(1, 50) {
			sim.TooLarge.Store(id, true)
		}
		task, err := c.Add(
			t.Context(),
			ssb.Data{Object: &api.BatchObject{UUID: id}},
		)

		// If the test server can OOM, then we won't try and guess
		// whether the error was expected, as OOM happens randomly.
		if !sim.CanOOM {
			assert.NoError(t, err, "add error")
		}
		if err != nil {
			if assert.NotNil(t, task, "nil task") {
				tasks = append(tasks, task)
			}
		}
	}

	log.Println("added all data -> close client")

	// Close the client.
	err := c.Close()
	if sim.CanOOM && err != nil {
		assert.ErrorContains(t, err, "OOM", "close error, CanOOM=%t", sim.CanOOM)
	} else {
		assert.NoError(t, err, "close error")
	}

	// Shutdown the server.
	close(conn)
	<-srv.done

	for _, task := range tasks {
		<-task.Done()
		id, err := task.ID(), task.Err()
		if _, tooLarge := sim.TooLarge.Load(id); tooLarge {
			require.ErrorIs(t, err, ssb.ErrTooLarge, "%s is too large", id)
			continue
		}
		stat, ok := srv.seen[task.ID()]
		if ok {
			require.Equal(t, stat.Err, err, "task %s error", id)
			require.Equal(t, stat.Retries, task.TimesRetried(), "task %s error", id)
		}
	}
	t.Log("FINISH TEST")
}

type Batch struct {
	*Simulation
	values []string
}
type Event ssb.Event

type Server struct {
	*Simulation
	conn <-chan *Stream
	work chan *Batch
	seen map[string]taskStat
	done chan struct{}
}

type taskStat struct {
	Retries int
	Err     error
}

func (srv *Server) run() {
	srv.Helper()

	defer close(srv.work)
	defer close(srv.done)

	for stream := range srv.conn {
		if err := stream.srvSend(Event{Started: true}); err != nil {
			continue
		}

		var oom bool
	Conn:
		for {
			if srv.prng.Chance(1, 30) || oom {
				srv.Logf("[server]: Shutting down (OOM=%t)", oom)
				stream.srvSend(Event{ShuttingDown: true})
				break Conn
			}

			if srv.prng.Chance(1, 3) {
				srv.Logf("[server]: Backoff")
				backoff := srv.prng.RangeInclusive(srv.BatchSize/2, srv.BatchSize*2)
				if err := stream.srvSend(Event{Backoff: &backoff}); err != nil {
					break Conn
				}
			}

			select {
			case batch, ok := <-stream.srvRecv():
				if !ok {
					for len(srv.work) > 0 {
						results := srv.processBatch(<-srv.work)
						if err := stream.srvSend(Event{Results: results}); err != nil {
							break Conn
						}
					}
					break Conn
				}
				assert.NotNil(srv.T, batch)
				assert.NotEmpty(srv.T, batch.values)
				srv.Logf("[server]: Received batch (%d) %p", len(batch.values), batch)

				if srv.CanOOM && srv.prng.Chance(1, 20) {
					oom = true
					var exitAfter time.Duration
					// if srv.prng.Bool() {
					exitAfter = 5 * time.Second
					// }
					if err := stream.srvSend(Event{OOM: &ssb.OOM{ExitAfter: exitAfter}}); err != nil {
						break Conn
					}
					continue
				}
				srv.work <- batch
				srv.Log("[server]: Ack batch")
				if err := stream.srvSend(Event{Ack: true}); err != nil {
					srv.Log("[server]: Failed to ack batch")
					break Conn
				}
				srv.Log("[server]: Ack-ed batch")

			case batch := <-srv.work:
				// if srv.prng.Chance(1, 10) {
				// 	srv.Log("[server]: Busy, return batch to queue")
				// 	srv.work <- batch
				// 	continue
				// }
				results := srv.processBatch(batch)
				if err := stream.srvSend(Event{Results: results}); err != nil {
					break Conn
				}
			}
		}

		srv.Logf("[server]: Discard any remaining work (%d)", len(srv.work))
		for len(srv.work) > 0 {
			<-srv.work
		}
		srv.Logf("[server]: Close stream")
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
		if srv.prng.Chance(1, srv.RetryLimit-t.Retries+4) {
			results.Failed[id] = testkit.ErrWhaam
			t.Err = testkit.ErrWhaam
			failed++
		} else {
			results.OK = append(results.OK, id)
			t.Err = nil
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

	ctx, cancel := context.WithCancelCause(ctx)
	s := &Stream{
		Simulation: t.Simulation,
		ctx:        ctx,
		cancel:     cancel,
		batchc:     make(chan *Batch),
		eventc:     make(chan Event),
	}

	t.conn <- s
	return s, nil
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
	assert.NotNil(s.T, req, "nil request")
	assert.IsType(s.T, (*Batch)(nil), req, "bad request")

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
	if s.prng.Chance(1, 25) {
		s.cancel(errors.New("bad network"))
	}

	select {
	case ev := <-s.eventc:
		return ssb.Event(ev), nil
	case <-s.ctx.Done():
		return ssb.Event{}, context.Cause(s.ctx)
	}
}

func (s *Stream) Close() error {
	s.Log("[stream.Close]: Client closed it's half of the stream")
	close(s.batchc)
	return nil
}

func (s *Stream) srvSend(ev Event) error {
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
	return &Batch{
		Simulation: t.Simulation,
		values:     make([]string, 0, 92), // FIXME: random cap, not 92
	}
}

func (t *Transport) Prepare(data ssb.Data) (any, error) {
	t.Helper()
	assert.NotNil(t.T, data.Object, "nil batch object")

	id := data.Object.UUID.String()
	if _, ok := t.TooLarge.Load(id); ok {
		return nil, ssb.ErrTooLarge
	}
	return id, nil
}

func (b *Batch) Add(v any) (added, full bool) {
	b.Helper()
	defer require.LessOrEqual(b.T, len(b.values), cap(b.values))

	assert.IsType(b.T, *new(string), v, "bad value in Add")
	if len(b.values) == cap(b.values) {
		return false, true
	}
	if slices.Contains(b.values, v.(string)) {
		panic(v)
	}
	b.values = append(b.values, v.(string))
	return true, len(b.values) == cap(b.values)
}

func (b *Batch) String() string { return fmt.Sprint(b.values) }
