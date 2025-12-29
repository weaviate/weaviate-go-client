package query

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/types"
)

type NearVectorTarget any

// NearVectorFunc runs plain near vector search.
type NearVectorFunc func(context.Context, NearVectorTarget, ...NearVectorOption) (*Result, error)

// GroupBy runs near vector search with a group by clause.
func (nv NearVectorFunc) GroupBy(ctx context.Context, target NearVectorTarget, groupBy string, options ...NearVectorOption) (*GroupByResult[types.Map], error) {
	return &GroupByResult[types.Map]{}, nil
}

type nearVectorRequest struct {
	commonOptions
	Target              NearVectorTarget
	Distance, Certainty *float64
	GroupBy             *GroupBy
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

// WithDistance sets the `distance` parameter.
type WithDistance float64

var _ NearVectorOption = (*WithDistance)(nil)

// WithCertainty sets the `certainty` parameter.
type WithCertainty float64

var _ NearVectorOption = (*WithCertainty)(nil)

func nearVector(ctx context.Context, t internal.Transport, target NearVectorTarget, options ...NearVectorOption) (*Result, error) {
	var nv nearVectorRequest
	for _, opt := range options {
		opt.apply(&nv)
	}
	return &Result{}, nil
}

// nearVectorFunc makes internal.Transport available to nearVector via a closure.
func nearVectorFunc(t internal.Transport) NearVectorFunc {
	return func(ctx context.Context, target NearVectorTarget, options ...NearVectorOption) (*Result, error) {
		return nearVector(ctx, t, target, options...)
	}
}

func (opt WithLimit) apply(r *nearVectorRequest) {
	r.Limit = (*int)(&opt)
}

func (opt WithOffset) apply(r *nearVectorRequest) {
	r.Offset = (*int)(&opt)
}

func (opt WithAutoLimit) apply(r *nearVectorRequest) {
	r.AutoLimit = (*int)(&opt)
}

func (opt WithAfter) apply(r *nearVectorRequest) {
	r.After = (*string)(&opt)
}

func (opt returnPropertiesOption) apply(r *nearVectorRequest) {
	r.ReturnProperties = opt
}

func (opt WithCertainty) apply(r *nearVectorRequest) {
	r.Certainty = (*float64)(&opt)
}

func (opt WithDistance) apply(r *nearVectorRequest) {
	r.Distance = (*float64)(&opt)
}

// apply implements NearVectorOption.
func (r *returnMetadataOption) apply(req *nearVectorRequest) {
	req.ReturnMetadata = append(req.ReturnMetadata, *r...)
}

// apply implements NearVectorOption.
func (gb groupByOption) apply(r *nearVectorRequest) {
	r.GroupBy = (*GroupBy)(&gb)
}

// NearVectorOption can be applied as a single option,
// in which case it will individually apply the options it comprises.
func (opts NearVectorOptions) apply(r *nearVectorRequest) {
	for _, opt := range opts {
		opt.apply(r)
	}
}
