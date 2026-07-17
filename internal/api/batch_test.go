package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
)

func TestBatchRequest(t *testing.T) {
	var br api.BatchRequest

	type want struct{ objects, references int }
	checkContains := func(want want) {
		req, err := br.MarshalMessage()
		require.NoError(t, err, "marshal batch request")

		data := req.GetData()
		require.NotNil(t, data, "request data")

		assert.Len(t, data.GetObjects().GetValues(), want.objects, "objects in batch")
		assert.Len(t, data.GetReferences().GetValues(), want.references, "references in batch")
	}

	checkContains(want{objects: 0, references: 0})

	for range 3 {
		err := br.AddObject(new(api.BatchObject))
		require.NoError(t, err, "add object to batch")
	}
	checkContains(want{objects: 3, references: 0})

	br.PopObject()
	checkContains(want{objects: 2, references: 0})

	br.AddReference(new(api.Reference))
	checkContains(want{objects: 2, references: 1})

	br.PopReference()
	checkContains(want{objects: 2, references: 0})
}
