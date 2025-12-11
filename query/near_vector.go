package query

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/types"
)

// For demo purposes only
var (
	mockObject_1 = types.Object[types.Map]{
		UUID:       "uuid-1",
		Properties: map[string]any{"number": "one"},
		Vectors: map[string]types.Vector{
			"1d": {Name: "1d", Single: []float32{1, 2, 3}},
			"2d": {Name: "2d", Multi: [][]float32{{1, 2, 3}, {1, 2, 3}}},
		},
	}
	mockObject_2 = types.Object[types.Map]{
		UUID:       "uuid-2",
		Properties: map[string]any{"number": "two"},
		Vectors: map[string]types.Vector{
			"1d": {Name: "1d", Single: []float32{1, 2, 3}},
			"2d": {Name: "2d", Multi: [][]float32{{1, 2, 3}, {1, 2, 3}}},
		},
	}
	mockObject_3 = types.Object[types.Map]{
		UUID:       "uuid-3",
		Properties: map[string]any{"number": "three"},
		Vectors: map[string]types.Vector{
			"1d": {Name: "1d", Single: []float32{1, 2, 3}},
			"2d": {Name: "2d", Multi: [][]float32{{1, 2, 3}, {1, 2, 3}}},
		},
	}
)

// NearVectorFunc runs plain near vector search.
type NearVectorFunc func(context.Context, NearVectorTarget, ...NearVectorOption) (*Result, error)

// GroupBy runs near vector search with a group by clause.
func (nv NearVectorFunc) GroupBy(ctx context.Context, target NearVectorTarget, groupBy string, options ...NearVectorOption) (*GroupByResult, error) {
	ctxcpy := internal.ContextWithGroupByResult(ctx)
	_, err := nv(ctxcpy, target, NearVectorOptions(options).Add(withGroupBy(groupBy)))
	if err != nil {
		return nil, err
	}
	_ = internal.GroupByResultFromContext(ctxcpy)
	return &GroupByResult{
		Objects: []GroupByObject[types.Map]{
			{Object: mockObject_1, BelongsToGroup: "a"},
			{Object: mockObject_2, BelongsToGroup: "b"},
			{Object: mockObject_3, BelongsToGroup: "b"},
		},
		Groups: map[string]Group[types.Map]{
			"a": {
				Name: "a", Size: 1, Objects: []GroupByObject[types.Map]{
					{Object: mockObject_1, BelongsToGroup: "a"},
				},
			},
			"b": {
				Name: "b", Size: 2, Objects: []GroupByObject[types.Map]{
					{Object: mockObject_2, BelongsToGroup: "b"},
					{Object: mockObject_3, BelongsToGroup: "b"},
				},
			},
		},
	}, nil
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
	return &Result{
		Objects: []types.Object[types.Map]{mockObject_1, mockObject_2, mockObject_3},
	}, nil
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
