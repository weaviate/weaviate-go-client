package stream

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	proto "github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/proto/v1"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/ssb"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/transport"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
)

type mockTransport struct {
	maxSize int
	stream  *mockStream
}

func (mt *mockTransport) MaxSize() int { return mt.maxSize }
func (mt *mockTransport) NewStream(context.Context) (transport.BatchStream, error) {
	return mt.stream, nil
}

func TestTransportAdapter_NewStream(t *testing.T) {
	rd := api.RequestDefaults{ConsistencyLevel: api.ConsistencyLevelQuorum}
	ms := new(mockStream)
	tad := NewAdapter(&mockTransport{stream: ms}, rd)
	assert.NotNil(t, tad, "nil transport adapter")

	s, err := tad.NewStream(t.Context())
	assert.NoError(t, err, "new stream")
	assert.NotNil(t, s, "nil stream")

	assert.NotNil(t, ms.sent, "sent one message down the stream")
	if assert.IsType(t, ms.sent.Message, (*proto.BatchStreamRequest_Start_)(nil)) {
		assert.Equal(t,
			proto.ConsistencyLevel_CONSISTENCY_LEVEL_QUORUM,
			ms.sent.GetStart().GetConsistencyLevel(),
			"stream started with a wrong consistency level",
		)
	}
}

func TestTransportAdapter_NewRequest(t *testing.T) {
	tad := NewAdapter(&mockTransport{maxSize: 1 << 10}, api.RequestDefaults{})
	assert.NotNil(t, tad, "nil transport adapter")

	req := tad.NewRequest()
	assert.NotNil(t, req, "nil batch request")

	if assert.IsType(t, req, (*api.BatchRequest)(nil)) {
		assert.Equal(t, 1<<10, req.(*api.BatchRequest).MaxSize, "max batch size")
	}
}

type mockStream struct {
	// Embedded BatchStream helps client satisfy the interface
	// without implementing methods not used in the test.
	transport.BatchStream

	sent *proto.BatchStreamRequest // Last request that was sent via this client.
	recv *proto.BatchStreamReply   // Return value for [Client.Recv].
}

func (ms *mockStream) Send(req *proto.BatchStreamRequest) error {
	ms.sent = req
	return nil
}

func (ms *mockStream) Recv() (*proto.BatchStreamReply, error) {
	return ms.recv, nil
}

func TestStream_Send(t *testing.T) {
	req := api.BatchRequest{
		MaxSize: 1 << 10,
	}

	objects := make([]*proto.BatchObject, 5)
	for i := range objects {
		objects[i] = &proto.BatchObject{Uuid: uuid.New().String()}
		added, _ := req.Add(objects[i])
		require.Truef(t, added, "added object #%d", i)
	}

	var ms mockStream
	sad := streamAdapter{stream: &ms}
	err := sad.Send(&req)
	assert.NoError(t, err, "send error")

	if assert.IsType(t, ms.sent.Message, (*proto.BatchStreamRequest_Data_)(nil)) {
		assert.EqualExportedValues(t, objects, ms.sent.GetData().GetObjects().GetValues())
	}
}

func TestStream_Recv(t *testing.T) {
	for _, tt := range []struct {
		name string
		recv *proto.BatchStreamReply
		want ssb.Event
	}{
		{
			name: "started",
			recv: &proto.BatchStreamReply{
				Message: &proto.BatchStreamReply_Started_{
					Started: &proto.BatchStreamReply_Started{},
				},
			},
			want: ssb.Event{Started: true},
		},
		{
			name: "shutting down",
			recv: &proto.BatchStreamReply{
				Message: &proto.BatchStreamReply_ShuttingDown_{
					ShuttingDown: &proto.BatchStreamReply_ShuttingDown{},
				},
			},
			want: ssb.Event{ShuttingDown: true},
		},
		{
			name: "backoff",
			recv: &proto.BatchStreamReply{
				Message: &proto.BatchStreamReply_Backoff_{
					Backoff: &proto.BatchStreamReply_Backoff{
						BatchSize: 92,
					},
				},
			},
			want: ssb.Event{Backoff: testkit.Ptr(92)},
		},
		{
			name: "oom",
			recv: &proto.BatchStreamReply{
				Message: &proto.BatchStreamReply_OutOfMemory_{
					OutOfMemory: &proto.BatchStreamReply_OutOfMemory{
						WaitTime: 92,
					},
				},
			},
			want: ssb.Event{OOM: &ssb.OOM{ExitAfter: 92 * time.Second}},
		},
		{
			name: "acks",
			recv: &proto.BatchStreamReply{
				Message: &proto.BatchStreamReply_Acks_{
					Acks: &proto.BatchStreamReply_Acks{
						Uuids:   []string{"uuid-0", "uuid-1"},
						Beacons: []string{"ref-0", "ref-1"},
					},
				},
			},
			want: ssb.Event{Ack: true},
		},
		{
			name: "results",
			recv: &proto.BatchStreamReply{
				Message: &proto.BatchStreamReply_Results_{
					Results: &proto.BatchStreamReply_Results{
						Successes: []*proto.BatchStreamReply_Results_Success{
							{Detail: &proto.BatchStreamReply_Results_Success_Uuid{
								Uuid: "uuid-0",
							}},
							{Detail: &proto.BatchStreamReply_Results_Success_Beacon{
								Beacon: "ref-0",
							}},
						},
						Errors: []*proto.BatchStreamReply_Results_Error{
							{
								Detail: &proto.BatchStreamReply_Results_Error_Beacon{
									Beacon: "ref-1",
								},
								Error: testkit.ErrWhaam.Error(),
							},
						},
					},
				},
			},
			want: ssb.Event{
				Results: &ssb.Results{
					OK: []string{"uuid-0", "ref-0"},
					Failed: map[string]error{
						"ref-1": testkit.ErrWhaam,
					},
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			sad := streamAdapter{stream: &mockStream{recv: tt.recv}}

			got, err := sad.Recv()
			require.NoError(t, err, "recv error")

			assert.EqualExportedValues(t, tt.want, got, "received event")
		})
	}
}
