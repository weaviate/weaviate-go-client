package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
	"google.golang.org/api/iterator"
)

// NewObjectIterator creates a new [ObjectIterator] for objects in the collection
// which client c targets. The context ctx will be used for all requests throughtout the
// iterators lifecycle, and canceling it will prevent the iterator from fetching any more objects.
func NewObjectIterator(ctx context.Context, c *Client) *ObjectIterator {
	dev.AssertNotNil(ctx, "ctx")
	dev.AssertNotNil(c, "client")

	it := &ObjectIterator{
		ctx: ctx,
		c:   c,
	}

	it.pageInfo, it.nextFunc = iterator.NewPageInfo(
		it.fetch,
		func() int { return len(it.items) },
		func() any { b := it.items; it.items = nil; return b })

	dev.AssertNotNil(it.pageInfo, "page info")
	dev.AssertNotNil(it.nextFunc, "next func")

	return it
}

// ObjectIterator implements a [Google-style iterator] for objects in the target collection.
// This iterator does not support concurrent modification or iteration.
//
// [Google-style iterator]: https://github.com/googleapis/google-cloud-go/wiki/Iterator-Guidelines
type ObjectIterator struct {
	ctx      context.Context
	c        *Client
	pageInfo *iterator.PageInfo
	nextFunc func() error
	items    []*Object[map[string]any]
}

var _ iterator.Pageable = (*ObjectIterator)(nil)

// PageInfo returns iterator's [iterator.PageInfo], which supports pagination.
//
// You can control the page size by setting [iterator.PageInfo.MaxSize] on the
// returned page info and update the cursor by setting [iterator.PageInfo.Token].
// These MUST be update synchronously before using the iterator again.
func (it *ObjectIterator) PageInfo() *iterator.PageInfo { return it.pageInfo }

// Next fetches the next object in the collection.
func (it *ObjectIterator) Next() (*Object[map[string]any], error) {
	if err := it.nextFunc(); err != nil {
		return nil, err
	}
	item := it.items[0]
	it.items = it.items[1:]
	return item, nil
}

// fetch fetches the next batch of pageSize objects and populates the internal buffer it.items.
func (it *ObjectIterator) fetch(pageSize int, pageToken string) (string, error) {
	req := OverAll{Limit: pageSize}
	if pageToken != "" {
		after, err := uuid.Parse(pageToken)
		if err != nil {
			return "", err
		}
		req.After = after
	}

	r, err := it.c.OverAll(it.ctx, req)
	if err != nil {
		return "", err
	}

	// iterator.Pager will mark the iterator as exhausted and return
	// [iterator.Done] on every subsequent call to Next or NextPage.
	//
	// TODO(dyma): should we also return iterator.Done in this case?
	// I find it odd that the user should check for iterator.Done when
	// using the Iterator and check for nextToken == "" when using the
	// same inside a Pager. I would prefer if we could return this error
	// every time the iterator's exhausted its resource.
	// Feels like we should also check if len(r.Objects) < pageSize,
	// which suggests there weren't any more objects on the server;
	// otherwise we're likely to need another call to get here if
	// the last request "undef-fetches".
	if len(r.Objects) == 0 {
		return "", nil
	}

	var i int
	for i = range r.Objects {
		it.items = append(it.items, &r.Objects[i])
	}
	return r.Objects[i].UUID.String(), nil
}
