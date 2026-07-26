package ssb_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/ssb"
)

func TestClient(t *testing.T) {
	t.Skip()

	t.Run("send a single batch", func(t *testing.T) {
		// TODO(dyma): initialize an actual client
		// ms := newMockStream(t)

		var started bool
		c := ssb.NewClient(ssb.ClientConfig{
			Context: t.Context(),
			// Start: func(context.Context, api.RequestDefaults) (ssb.Stream, error) {
			// 	started = true
			// 	return ms, nil
			// },
			BatchSize: 10,
		})
		assert.True(t, started, "new client started the batch stream")

		// TODO(dyma): find a way to check "as long as" conditions.
		// E.g.:
		// - as long as not STARTED, no message is sent
		// - as long as not RESULTS, client is not closed
		// - as long as not SHUTDOWN, new stream is not started

		var tasks []*ssb.Task
		var ids []string
		for range 10 {
			id := uuid.New()
			task, err := c.Add(ssb.Data{
				Object: &api.BatchObject{
					UUID: id,
				},
			})
			require.NoError(t, err, "add data to batch")
			require.NotNil(t, task, "task handle")
			assert.Equal(t, id.String(), task.ID())

			tasks = append(tasks, task)
			ids = append(ids, id.String())
		}

		// select {
		// case req := <-ms.recvServer():
		// 	assert.ElementsMatch(t, ids, req, "bad request")
		// case <-time.After(5 * time.Millisecond):
		// 	assert.Fail(t, "no DATA message")
		// }
		//
		// assert.NoError(t, c.Close(), "close client")
		// assert.Eventually(t, ms.closed, 5*time.Millisecond, time.Millisecond, "stream is closed")
		//
		// for _, task := range tasks {
		// 	select {
		// 	case <-task.Done():
		// 		assert.NoErrorf(t, task.Err(), "task-%s", task.ID())
		// 	case <-time.After(5 * time.Millisecond):
		// 	}
		// }
	})
}

// ----------------------------------------------------------------------------

// func newMockStream(t *testing.T) *mockStream {
// 	return &mockStream{
// 		t:       t,
// 		clientc: make(chan ssb.Event),
// 		serverc: make(chan []string, 1),
// 	}
// }
//
// var (
// 	_ ssb.Stream       = (*mockStream)(nil)
// 	_ ssb.BatchRequest = (*mockBatch)(nil)
// )
//
// type mockStream struct {
// 	t       *testing.T
// 	clientc chan ssb.Event // Channel for client-bound requests.
// 	serverc chan []string  // Channel for server-bound requests.
// }
//
// func (ms *mockStream) NewBatch() ssb.BatchRequest {
// 	return &mockBatch{stream: ms}
// }
//
// func (ms *mockStream) Recv() (ssb.Event, error) {
// 	return <-ms.clientc, nil
// }
//
// func (ms *mockStream) Close() error {
// 	close(ms.clientc)
// 	close(ms.serverc)
// 	return nil
// }
//
// func (ms *mockStream) closed() bool {
// 	var ok1, ok2 bool
// 	select {
// 	case _, ok1 = <-ms.clientc:
// 	case <-time.After(5 * time.Millisecond):
// 		require.Fail(ms.t, "clientc not closed")
// 	}
//
// 	select {
// 	case _, ok2 = <-ms.serverc:
// 	case <-time.After(5 * time.Millisecond):
// 		require.Fail(ms.t, "serverc not closed")
// 	}
//
// 	return !ok1 && !ok2
// }
//
// // recvServer reads the next server-bound message.
// func (ms *mockStream) recvServer() <-chan []string {
// 	return ms.serverc
// }
//
// // sendClient send a server-side event.
// func (ms *mockStream) sendClient(ev ssb.Event) {
// 	ms.clientc <- ev
// }
//
// type mockBatch struct {
// 	stream *mockStream
// 	ids    []string // IDs of data items in the batch.
// }
//
// func (mb *mockBatch) Add(d ssb.Data) error {
// 	mb.ids = append(mb.ids, d.ID())
// 	return nil
// }
//
// func (mb *mockBatch) Send() error {
// 	mb.stream.serverc <- mb.ids
// 	return nil
// }
