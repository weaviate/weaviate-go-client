package query

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

func NewObjectIterator(ctx context.Context, c *Client) *ObjectIterator {
	it := &ObjectIterator{
		ctx: ctx,
		c:   c,
	}

	it.pageInfo, it.nextFunc = iterator.NewPageInfo(
		it.fetch,
		func() int { return len(it.items) },
		func() any { b := it.items; it.items = nil; return b })

	return it
}

type ObjectIterator struct {
	ctx      context.Context
	c        *Client
	pageInfo *iterator.PageInfo
	nextFunc func() error
	items    []*Object[map[string]any]
}

var _ iterator.Pageable = (*ObjectIterator)(nil)

func (it *ObjectIterator) PageInfo() *iterator.PageInfo { return it.pageInfo }

func (it *ObjectIterator) Next() (*Object[map[string]any], error) {
	if err := it.nextFunc(); err != nil {
		return nil, err
	}
	item := it.items[0]
	it.items = it.items[1:]
	return item, nil
}

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

	var i int
	for i = range r.Objects {
		it.items = append(it.items, &r.Objects[i])
	}
	return r.Objects[i].UUID.String(), nil
}
