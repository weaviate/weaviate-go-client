package query

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
)

type NearVectorTarget api.NearVectorTarget

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
	return &GroupByResult{}, nil
}

type nearVectorRequest struct {
	commonOptions
	Target              NearVectorTarget
	Distance, Certainty *float64
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

func nearVector(ctx context.Context, t internal.Transport, rd api.RequestDefaults, target NearVectorTarget, options ...NearVectorOption) (*Result, error) {
	var nv nearVectorRequest
	for _, opt := range options {
		opt.apply(&nv)
	}

	req := &api.SearchRequest{
		RequestDefaults:  rd,
		Limit:            nv.Limit,
		AutoLimit:        nv.AutoLimit,
		Offset:           nv.Offset,
		After:            nv.After,
		ReturnProperties: nv.ReturnProperties,
		ReturnReferences: nv.ReturnReferences,
		ReturnVectors:    nv.ReturnVectors,
		ReturnMetadata:   api.NewSet(nv.ReturnMetadata),
		NearVector: &api.NearVector{
			Target:    nv.Target,
			Distance:  nv.Distance,
			Certainty: nv.Certainty,
		},
	}

	var resp api.SearchResponse
	if err := t.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("near vector: %w")
	}
	return &Result{}, nil
}

// nearVectorFunc makes internal.Transport available to nearVector via a closure.
func nearVectorFunc(t internal.Transport, rd api.RequestDefaults) NearVectorFunc {
	return func(ctx context.Context, target NearVectorTarget, options ...NearVectorOption) (*Result, error) {
		return nearVector(ctx, t, rd, target, options...)
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

func (opt ReturnVectorOption) apply(r *nearVectorRequest) {
	r.ReturnVectors = opt
}

func (opt returnMetadataOption) apply(r *nearVectorRequest) {
	r.ReturnMetadata = opt
}

func (opt WithCertainty) apply(r *nearVectorRequest) {
	r.Certainty = (*float64)(&opt)
}

func (opt WithDistance) apply(r *nearVectorRequest) {
	r.Distance = (*float64)(&opt)
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
