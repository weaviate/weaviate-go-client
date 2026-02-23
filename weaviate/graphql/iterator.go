package graphql

import (
	"context"

	"github.com/google/uuid"
)

// Iterator pages through all objects in a collection using per-shard cursor
// state. Obtain one via Search.Iterator(). The iterator is not safe for
// concurrent use.
type Iterator struct {
	search       *Search
	after        string
	shardCursors map[string]string
	done         bool
}

// Iterator returns a stateful iterator starting from the Search's current
// After position (if any). All other Search parameters (collection, filters,
// limit, properties, etc.) are preserved on each page fetch.
func (s *Search) Iterator() *Iterator {
	return &Iterator{
		search: s,
		after:  s.after,
	}
}

// HasNext reports whether there are more results to fetch.
func (it *Iterator) HasNext() bool {
	return !it.done
}

// Next fetches the next batch of results. Returns (nil, nil) when the
// iterator is exhausted.
func (it *Iterator) Next(ctx context.Context) ([]SearchResult, error) {
	if it.done {
		return nil, nil
	}

	// Build request from the template search, then apply iterator cursor state.
	req := it.search.togrpc()
	if it.after != "" {
		req.After = &it.after
	} else {
		req.After = nil
	}
	req.ShardCursors = it.shardCursors

	reply, err := it.search.grpcClient.Search(ctx, req)
	if err != nil {
		return nil, err
	}

	results := toResults(reply.Results)

	// Advance cursor state for the next call.
	it.shardCursors = reply.ShardCursors
	if len(results) > 0 {
		it.after = results[len(results)-1].ID
	}

	// Declare exhaustion when there are no results or all shards report nil UUID.
	if len(results) == 0 || isExhausted(reply.ShardCursors) {
		it.done = true
	}

	return results, nil
}

// isExhausted returns true when every shard cursor is the nil UUID, meaning
// no shard has more objects to scan.
func isExhausted(cursors map[string]string) bool {
	if len(cursors) == 0 {
		// No shard cursor info returned — cannot declare done based on cursors alone.
		return false
	}
	nilUUID := uuid.Nil.String()
	for _, v := range cursors {
		if v != nilUUID {
			return false
		}
	}
	return true
}
