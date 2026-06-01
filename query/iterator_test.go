package query_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/query"
	"google.golang.org/api/iterator"
)

func TestObjectIterator(t *testing.T) {
	rd := api.RequestDefaults{
		CollectionName:   "Songs",
		ConsistencyLevel: api.ConsistencyLevelOne,
		Tenant:           "john_doe",
	}

	t.Run("iterate", func(t *testing.T) {
		t.Skip()
		after1 := uuid.New() // UUID of the last object in 1st batch
		after2 := uuid.New() // UUID of the last object in 2nd batch

		transport := testkit.NewTransport(t, []testkit.Stub[api.SearchRequest, api.SearchResponse]{
			{
				Request: &api.SearchRequest{RequestDefaults: rd},
				Response: api.SearchResponse{Results: []api.Object{
					{Metadata: api.ObjectMetadata{UUID: testkit.UUID}},
					{Metadata: api.ObjectMetadata{UUID: testkit.UUID}},
					{Metadata: api.ObjectMetadata{UUID: after1}},
				}},
			},
			{
				Request: &api.SearchRequest{RequestDefaults: rd, After: after1},
				Response: api.SearchResponse{Results: []api.Object{
					{Metadata: api.ObjectMetadata{UUID: testkit.UUID}},
					{Metadata: api.ObjectMetadata{UUID: after2}},
				}},
			},
			{
				Request: &api.SearchRequest{RequestDefaults: rd, After: after2},
				// The collection only has 5 items, no objects in the response.
			},
		})
		c := query.NewClient(transport, rd)
		require.NotNil(t, c, "nil client")

		it := query.NewObjectIterator(t.Context(), c)
		require.NotNil(t, it, "nil iterator")

		// The collection has 5 items, which should all be fetched.
		for i := range 5 {
			o, err := it.Next()
			assert.NoErrorf(t, err, "next #%d", i)
			assert.NotNil(t, o, "object #%d", i)
		}

		// Each subsequent call to Next should return iterator.Done.
		for range 2 {
			o, err := it.Next()
			assert.ErrorIs(t, err, iterator.Done)
			assert.Nil(t, o, "object after iteration is done")
		}
	})

	t.Run("paginate", func(t *testing.T) {
		after1 := uuid.New() // UUID of the last object in 1st batch
		after2 := uuid.New() // UUID of the last object in 2nd batch

		transport := testkit.NewTransport(t, []testkit.Stub[api.SearchRequest, api.SearchResponse]{
			{
				Request: &api.SearchRequest{RequestDefaults: rd, Limit: 3},
				Response: api.SearchResponse{Results: []api.Object{
					{Metadata: api.ObjectMetadata{UUID: testkit.UUID}},
					{Metadata: api.ObjectMetadata{UUID: testkit.UUID}},
					{Metadata: api.ObjectMetadata{UUID: after1}},
				}},
			},
			{
				Request: &api.SearchRequest{RequestDefaults: rd, After: after1, Limit: 3},
				Response: api.SearchResponse{Results: []api.Object{
					{Metadata: api.ObjectMetadata{UUID: testkit.UUID}},
					{Metadata: api.ObjectMetadata{UUID: after2}},
				}},
			},
			{
				// Page will try to fill the page up to the desired size 3 by fetching 1 more item.
				// The collection only has 5 items, no objects in the response.
				Request: &api.SearchRequest{RequestDefaults: rd, After: after2, Limit: 1},
			},
		})
		c := query.NewClient(transport, rd)
		require.NotNil(t, c, "nil client")

		it := query.NewObjectIterator(t.Context(), c)
		require.NotNil(t, it, "nil iterator")

		p := iterator.NewPager(it, 3, "")
		require.NotNil(t, p, "nil pager")

		var total int

		// The first batch contains pageSize items (3), but the second one is undersized (2).
		// Pager will try to fetch the remaining 1 item before NextPage returns. As there are
		// only 5 objects in the collection, the iteration finishes there with an empty nextPageToken.
		wantToken := map[int]string{1: after1.String(), 2: ""}
		for i := 1; ; i++ {
			var objects []*query.Object[map[string]any]

			nextPageToken, err := p.NextPage(&objects)
			assert.NoError(t, err, "next page")
			assert.Equal(t, wantToken[i], nextPageToken, "bad next page token for batch #%d", i)

			total += len(objects)

			if nextPageToken == "" {
				break
			}
		}

		assert.Equal(t, 5, total, "number of fetched objects")
	})
}
