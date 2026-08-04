package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	proto "github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/proto/v1"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	protoutil "google.golang.org/protobuf/proto"
)

func TestBatchRequest(t *testing.T) {
	rd := api.RequestDefaults{
		CollectionName: "Songs",
		Tenant:         "john_doe",
	}

	protoObject, err := api.MarshalBatchObject(&api.BatchObject{
		UUID: testkit.UUID,
		Properties: map[string]any{
			"title":        "Man Made of Meat",
			"artist":       "Viagra Boys",
			"album":        "viagr aboys",
			"duration_sec": 189,
		},
	}, rd)
	require.NoError(t, err, "marshal object")

	protoReference := api.MarshalBatchReference(&api.Reference{
		Origin: api.ObjectPath{Collection: "Songs", Property: "performedBy"},
		Target: api.ObjectPath{UUID: testkit.UUID},
	}, rd)
	require.NoError(t, err, "marshal reference")

	type want struct{ objects, references int }

	// reserveSize calculates batch size in bytes that can fit want data.
	reserveSize := func(want want) int {
		objects := make([]*proto.BatchObject, want.objects)
		for i := range objects {
			objects[i] = protoObject
		}
		references := make([]*proto.BatchReference, want.references)
		for i := range references {
			references[i] = protoReference
		}
		return protoutil.Size(&proto.BatchStreamRequest{
			Message: &proto.BatchStreamRequest_Data_{
				Data: &proto.BatchStreamRequest_Data{
					Objects: &proto.BatchStreamRequest_Data_Objects{
						Values: objects,
					},
					References: &proto.BatchStreamRequest_Data_References{
						Values: references,
					},
				},
			},
		})
	}

	br := api.BatchRequest{
		MaxSize: reserveSize(want{
			objects:    2,
			references: 1,
		}),
	}

	checkContains := func(want want) {
		req, err := br.MarshalMessage()
		require.NoError(t, err, "marshal batch request")

		data := req.GetData()
		require.NotNil(t, data, "request data")

		assert.Len(t, data.GetObjects().GetValues(), want.objects, "objects in batch")
		assert.Len(t, data.GetReferences().GetValues(), want.references, "references in batch")
	}

	var added, full bool

	added, full = br.Add(protoObject)
	assert.True(t, added, "added object #1")
	assert.False(t, full, "full")
	checkContains(want{objects: 1})

	added, full = br.Add(protoObject)
	assert.True(t, added, "added object #2")
	assert.False(t, full, "full")
	checkContains(want{objects: 2})

	added, full = br.Add(protoReference)
	assert.True(t, added, "added reference #1")
	assert.True(t, full, "full")
	checkContains(want{objects: 2, references: 1})

	added, full = br.Add(protoReference)
	assert.False(t, added, "added object #3")
	assert.True(t, full, "full")
	checkContains(want{objects: 2, references: 1})

	added, full = br.Add(protoReference)
	assert.False(t, added, "added reference #2")
	assert.True(t, full, "full")
	checkContains(want{objects: 2, references: 1})
}
