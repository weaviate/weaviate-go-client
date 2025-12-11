package query

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v5/internal"
	"github.com/weaviate/weaviate-go-client/v5/types"
)

// NearVectorFunc runs plain near vector search.
type NearVectorFunc[P types.Properties] func(context.Context, NearVectorTarget, ...NearVectorOption) (*Result[P], error)

// GroupBy runs near vector search with a group by clause.
func (nv NearVectorFunc[P]) GroupBy(ctx context.Context, target NearVectorTarget, groupBy string, options ...NearVectorOption) (*GroupByResult[P], error) {
	ctxcpy := internal.ContextWithGroupByResult(ctx)
	_, err := nv(ctxcpy, target, NearVectorOptions(options).Add(withGroupBy(groupBy)))
	if err != nil {
		return nil, err
	}
	res := internal.GroupByResultFromContext(ctxcpy)
	var p GroupByResult[P] = GroupByResult[types.Map]{Objects: []GroupByObject[types.Map]{
		{Object: types.Object[types.Map]{Properties: make(types.Map)}},
	}}
	return &p, nil
	// return &GroupByResult[types.Map]{
	// 	Objects: nil,
	// 	Groups:  nil,
	// }, nil
}

type nearVectorRequest struct {
	commonOptions
	Target              NearVectorTarget
	Distance, Certainty *float32
}

// NearVectorOption populates nearVectorRequest.
type NearVectorOption interface {
	apply(*nearVectorRequest)
}

type NearVectorOptions []NearVectorOption

func (opts NearVectorOptions) Add(options ...NearVectorOption) NearVectorOptions {
	opts = append(opts, options...)
	return opts
}

// DistanceOption sets the `distance` parameter.
type DistanceOption float32

var _ NearVectorOption = (*DistanceOption)(nil)

func WithDistance(l float32) DistanceOption {
	return DistanceOption(l)
}

// CertaintyOption sets the `certainty` parameter.
type CertaintyOption float32

var _ NearVectorOption = (*CertaintyOption)(nil)

func WithCertainty(l float32) CertaintyOption {
	return CertaintyOption(l)
}

func nearVector(context.Context, internal.Transport, NearVectorTarget, ...NearVectorOption) (*Result, error) {
	return nil, nil
}

// nearVectorFunc makes internal.Transport available to nearVector via a closure.
func nearVectorFunc(t internal.Transport) NearVectorFunc {
	return func(ctx context.Context, target NearVectorTarget, options ...NearVectorOption) (*Result, error) {
		return nearVector(ctx, t, target, options...)
	}
}

func (l LimitOption) apply(r *nearVectorRequest) {
	r.commonOptions.Limit = (*int)(&l)
}

func (l OffsetOption) apply(r *nearVectorRequest) {
	r.commonOptions.Offset = (*int)(&l)
}

func (l AutoLimitOption) apply(r *nearVectorRequest) {
	r.commonOptions.AutoLimit = (*int)(&l)
}

func (d CertaintyOption) apply(r *nearVectorRequest) {
	r.Certainty = (*float32)(&d)
}

func (d DistanceOption) apply(r *nearVectorRequest) {
	r.Distance = (*float32)(&d)
}

// apply implements NearVectorOption.
func (gb groupByOption) apply(r *nearVectorRequest) {
	r.commonOptions.GroupBy = (*GroupBy)(&gb)
}

func (opts NearVectorOptions) apply(r *nearVectorRequest) {
	for _, opt := range opts {
		opt.apply(r)
	}
}

type NearVectorTarget interface {
	ToProto()
}
