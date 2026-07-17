package stream

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	proto "github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/proto/v1"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/ssb"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	protoutil "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestBatch(t *testing.T) {
	rd := api.RequestDefaults{
		CollectionName: "Songs",
		Tenant:         "john_doe",
	}
	properties := map[string]any{
		"title":        "Man Made of Meat",
		"artist":       "Viagra Boys",
		"album":        "viagr aboys",
		"duration_sec": 189,
	}

	protoProperties, err := structpb.NewStruct(properties)
	require.NoError(t, err, "marshal properties")

	object := &api.BatchObject{
		UUID:       testkit.UUID,
		Properties: properties,
	}

	// protoObject is the marshaled version of object.
	protoObject := &proto.BatchObject{
		Collection: rd.CollectionName,
		Tenant:     rd.Tenant,
		Uuid:       testkit.UUID.String(),
		Properties: &proto.BatchObject_Properties{
			NonRefProperties: protoProperties,
		},
	}

	// expectObjects calculates batch size in bytes that can fit n protoObject copies.
	expectObjects := func(n int) int {
		objects := make([]*proto.BatchObject, n)
		for i := range n {
			objects[i] = protoObject
		}
		return protoutil.Size(&proto.BatchStreamRequest{
			Message: &proto.BatchStreamRequest_Data_{
				Data: &proto.BatchStreamRequest_Data{
					Objects: &proto.BatchStreamRequest_Data_Objects{
						Values: objects,
					},
					References: &proto.BatchStreamRequest_Data_References{
						Values: nil,
					},
				},
			},
		})
	}

	t.Run("add", func(t *testing.T) {
		t.Run("all objects fit the batch", func(t *testing.T) {
			bs := batchStream{maxSize: expectObjects(3)}
			b := bs.NewBatch()

			for i := range 3 {
				err := b.Add(ssb.Data{Object: object})
				if assert.LessOrEqual(t, batchSize(t, b.(*batch)), bs.maxSize) {
					assert.NoErrorf(t, err, "add item #%d to batch", i+1)
				}
			}

			err := b.Add(ssb.Data{Object: object})
			assert.ErrorIs(t, err, ssb.ErrBatchFull, "add 4th item")
		})

		t.Run("object is too large", func(t *testing.T) {
			bs := batchStream{RequestDefaults: rd, maxSize: expectObjects(1) / 2}
			b := bs.NewBatch()

			err := b.Add(ssb.Data{Object: object})
			assert.ErrorIs(t, err, ssb.ErrTooLarge)
		})
	})

	t.Run("send", func(t *testing.T) {
		var c client
		bs := batchStream{
			RequestDefaults: rd,
			c:               &c,
			maxSize:         expectObjects(5),
		}
		b := bs.NewBatch()

		require.NoError(t, b.Add(ssb.Data{Object: object}), "add item to batch")

		err := b.Send()
		assert.NoError(t, err, "send batch")
		assert.NotNil(t, c.sent, "sent request")

		objects := c.sent.GetData().GetObjects().GetValues()
		assert.Len(t, objects, 1)
		assert.EqualExportedValues(t, protoObject, objects[0])
	})
}

func batchSize(t *testing.T, b *batch) int {
	size, err := b.size()
	require.NoError(t, err, "get batch size")
	return size
}

type client struct {
	// Embedded Client helps client satisfy the interface
	// without implementing methods not used in the test.
	Client

	sent *proto.BatchStreamRequest // Last request that was sent via this client.
	recv *proto.BatchStreamReply   // Return value for [Client.Recv].
}

func (c *client) Send(req *proto.BatchStreamRequest) error {
	c.sent = req
	return nil
}

func (c *client) Recv() (*proto.BatchStreamReply, error) {
	return c.recv, nil
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
			want: ssb.Event{
				Started: true,
			},
		},
		{
			name: "shutting down",
			recv: &proto.BatchStreamReply{
				Message: &proto.BatchStreamReply_ShuttingDown_{
					ShuttingDown: &proto.BatchStreamReply_ShuttingDown{},
				},
			},
			want: ssb.Event{
				ShuttingDown: true,
			},
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
			want: ssb.Event{
				Backoff: testkit.Ptr(92),
			},
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
			want: ssb.Event{
				OOM: &ssb.OOM{
					ReconnectAfter: 92 * time.Second,
				},
			},
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
			want: ssb.Event{
				Acks: []string{"uuid-0", "uuid-1", "ref-0", "ref-1"},
			},
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
								Error: "Whaam!",
							},
						},
					},
				},
			},
			want: ssb.Event{
				Results: &ssb.Results{
					OK: []string{"uuid-0", "ref-0"},
					Failed: map[string]string{
						"ref-1": "Whaam!",
					},
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			bs := batchStream{c: &client{recv: tt.recv}}

			got, err := bs.Recv()
			require.NoError(t, err, "recv error")

			assert.EqualExportedValues(t, tt.want, got, "received event")
		})
	}
}
