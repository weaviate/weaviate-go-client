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
	prng     *testkit.PRNG
	tooLarge sync.Map

	TaskCount   int
	BatchSize   int
	MessageCap  int
	RetryLimit  int
	ReconnCount int
	ReconnLimit int
	CanOOM      bool
}

func (sim *Simulation) newClient(conn chan *Stream) *ssb.Client {
	return ssb.NewClient(ssb.ClientConfig{
		Context: sim.Context(),
		Transport: &Transport{
			Simulation: sim,
			conn:       conn,
		},
		QueueSize: sim.prng.RangeInclusive(sim.BatchSize/2, sim.BatchSize*2),
		BatchSize: sim.BatchSize,
		Reconnect: ssb.ReconnectPolicy{
			Limit:     sim.ReconnLimit,
			DelayFunc: func(int) time.Duration { return 0 },
		},
		CanRetry: func(id string, retries int, _ error) bool {
			return retries < sim.RetryLimit
		},
	})
}

func (sim *Simulation) newServer(conn chan *Stream) *Server {
	return &Server{
		Simulation: sim,
		conn:       conn,
		work:       make(chan *Batch, 8),
		seen:       make(map[string]int, sim.TaskCount),
		done:       make(chan struct{}),
	}
}

func (sim *Simulation) Backoff() bool      { return sim.prng.Chance(1, 3) }
func (sim *Simulation) BadNetwork() bool   { return sim.prng.Chance(1, 30) }
func (sim *Simulation) ShuttingDown() bool { return sim.prng.Chance(1, 50) }
func (sim *Simulation) OOM() bool          { return sim.CanOOM && sim.prng.Chance(1, 50) }

// CheckSize stores the ID of this task if it is "too large".
func (sim *Simulation) CheckSize(id string) {
	if sim.prng.Chance(1, 50) {
		sim.tooLarge.Store(id, true)
	}
}

// TooLarge returns true if the task is "too large".
func (sim *Simulation) TooLarge(id string) bool {
	_, ok := sim.tooLarge.Load(id)
	return ok
}

func TestClient(t *testing.T) {
	prng := testkit.NewPRNG(t)

	const (
		N            = 10
		retryLimit   = 3
		maxTaskCount = 1000
	)

	var (
		added int // Added to batch stream.
		seen  int // Arrived to the server.
		ok    int // Succeeded.
		fail  int // Failed.
	)

	for range N {
		t.Run("ssb client fuzz", func(t *testing.T) {
			sim := Simulation{
				T:    t,
				prng: prng,

				TaskCount:   prng.IntInclusive(maxTaskCount),
				RetryLimit:  retryLimit,
				BatchSize:   prng.RangeInclusive(16, 64),
				MessageCap:  prng.RangeInclusive(32, 128),
				ReconnLimit: prng.RangeInclusive(3, 5),
				CanOOM:      prng.Bool(),
			}

			conn := make(chan *Stream)
			c := sim.newClient(conn)
			srv := sim.newServer(conn)

			go srv.run()

			// Add all data to the client.
			tasks := make([]*ssb.Task, 0, sim.TaskCount)
			for i := 0; i < sim.TaskCount; i++ {
				id := uuid.New()
				sim.CheckSize(id.String())

				task, err := c.Add(
					t.Context(),
					ssb.Data{Object: &api.BatchObject{UUID: id}},
				)

				// If the test server can OOM, then we won't try and guess
				// whether the error was expected, as OOM happens randomly.
				if !sim.CanOOM && err != context.Canceled {
					assert.NoError(t, err, "add error")
				}

				if err == nil {
					if assert.NotNil(t, task, "nil task") {
						tasks = append(tasks, task)
					}
				}
			}

			// Close the client.
			c.Close() // nolint:errcheck

			// Shutdown the server.
			close(conn)
			<-srv.done

			// Wait for all tasks to complete.
			for _, task := range tasks {
				<-task.Done()
				require.LessOrEqual(t, task.TimesRetried(), retryLimit, "task %s retries", task.ID())

				switch task.Err() {
				case nil:
					ok++
				default:
					fail++
				}
			}

			added += sim.TaskCount
			seen += len(srv.seen)
		})
	}

	require.GreaterOrEqual(t, seen, int(float64(added)*.9), "over 90% of all data arrive at the server")
	require.GreaterOrEqual(t, ok, int(float64(seen)*.75), "over 75% of submitted tasks succeed")
	require.LessOrEqual(t, fail, int(float64(seen)*.25), "under 25% of submitted tasks fail")
}

type (
	Batch struct {
		*Simulation
		values []string
	}
	Event ssb.Event
)

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

func (t *Transport) NewRequest() ssb.BatchRequest {
	t.Helper()
	require.Greater(t.T, t.MessageCap, 0, "max message size")

	return &Batch{
		Simulation: t.Simulation,
		values:     make([]string, 0, t.MessageCap),
	}
}

func (t *Transport) Prepare(data ssb.Data) (any, error) {
	t.Helper()
	assert.NotNil(t.T, data.Object, "nil batch object")

	id := data.Object.UUID.String()
	if t.TooLarge(id) {
		return nil, ssb.ErrTooLarge
	}
	return id, nil
}

func (b *Batch) Add(v any) (added, full bool) {
	b.Helper()
	defer require.LessOrEqual(b.T, cap(b.values), b.MessageCap, "batch grew beyond MessageCap")

	assert.IsType(b.T, *new(string), v, "bad value in Add")
	assert.NotContains(b.T, b.values, v, "duplicate value in batch %s", v)

	if len(b.values) == cap(b.values) {
		return false, true
	}
	b.values = append(b.values, v.(string))
	return true, len(b.values) == cap(b.values)
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

	if s.BadNetwork() {
		s.cancel(errors.New("bad network"))
	}

	select {
	case s.batchc <- req.(*Batch):
	case <-s.ctx.Done():
		return io.EOF
	}
	return nil
}

func (s *Stream) Recv() (ssb.Event, error) {
	if s.BadNetwork() {
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

type Server struct {
	*Simulation
	conn <-chan *Stream
	work chan *Batch
	seen map[string]int
	done chan struct{}
}

func (srv *Server) run() {
	srv.Helper()

	defer close(srv.done)

	for stream := range srv.conn {
		if err := stream.srvSend(Event{Started: true}); err != nil {
			continue
		}

		var oom bool
		work := make(chan *Batch, 8)

	Conn:
		for {
			if srv.ShuttingDown() || oom {
				stream.srvSend(Event{ShuttingDown: true}) // nolint:errcheck
				break Conn
			}

			if srv.Backoff() {
				backoff := srv.prng.RangeInclusive(srv.BatchSize/2, srv.BatchSize*2)
				if err := stream.srvSend(Event{Backoff: &backoff}); err != nil {
					break Conn
				}
			}

			select {
			case batch, ok := <-stream.srvRecv():
				if !ok { // Client closed it's side of the stream.
					var err error
					for len(work) > 0 && err == nil {
						err = srv.processBatch(<-work, stream)
					}
					break Conn
				}

				assert.NotNil(srv.T, batch)
				assert.NotEmpty(srv.T, batch.values)

				if oom = srv.OOM(); oom {
					var exitAfter time.Duration // Let the client think server is unresponsive.
					if /*srv.prng.Bool()*/ true {
						exitAfter = 5 * time.Second // Give the server enought time to respond.
					}
					if err := stream.srvSend(Event{OOM: &ssb.OOM{ExitAfter: exitAfter}}); err != nil {
						break Conn
					}
					continue
				}

				work <- batch
				if err := stream.srvSend(Event{Ack: true}); err != nil {
					break Conn
				}

			case batch := <-work:
				if err := srv.processBatch(batch, stream); err != nil {
					break Conn
				}
			}
		}

		stream.srvClose()
	}
}

func (srv *Server) processBatch(b *Batch, s *Stream) error {
	results := &ssb.Results{
		OK:     make([]string, 0),
		Failed: make(map[string]error),
	}
	for _, id := range b.values {
		retries, ok := srv.seen[id]
		if !ok {
			retries = -1
		}
		srv.seen[id] = retries

		// On it's last retry the task fails with a 1/4 chance.
		if srv.prng.Chance(1, srv.RetryLimit-retries+4) {
			results.Failed[id] = testkit.ErrWhaam
		} else {
			results.OK = append(results.OK, id)
		}
	}
	if err := s.srvSend(Event{Results: results}); err != nil {
		return err
	}

	// Update retry stats after the client's received Results,
	// otherwise we might get out of sync.
	for _, id := range b.values {
		srv.seen[id] = srv.seen[id] + 1
	}
	return nil
}
