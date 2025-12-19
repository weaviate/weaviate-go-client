package query

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
	"github.com/weaviate/weaviate-go-client/v6/types"
)

type NearVector struct {
	Limit                  int               // Limit the number of results returned for the query.
	Offset                 int               // Skip the first N objects in the collection.
	AutoLimit              int               // Return objects in the first N similarity clusters.
	After                  uuid.UUID         // Skip all objects before the one with this ID.
	ReturnReferences       []*Reference      // Select reference properties to return.
	ReturnMetadata         []Metadata        // Select query and object metadata to return for each object.
	ReturnVectors          []string          // List vectors to return for each object in the result set.
	ReturnNestedProperties []*NestedProperty // Return object properties and a subset of their nested properties.

	// Select a subset of properties to return. By default, all properties are returned.
	// To not return any properties, initialize this value to an empty slice explicitly.
	ReturnProperties []string

	// Similarity specifies a cutoff point for query results.
	// Use Distance() to set it as maximum distance between vectors.
	// Use Certainty() to set it to a normalized value between 0 and 1.
	// Prefer expressing Similarity in terms of vector distance, as it is a more conventional metric.
	Similarity *Similarity

	// TODO(dyma): document
	Target NearVectorTarget

	// groupBy can only be set by NearVeectorFunc.GroupBy, as it affects the shape of the response.
	groupBy *GroupBy
}

// Similarity is a cutoff point for query results.
type Similarity struct{ distance, certainty *float64 }

// Distance sets a similarity cutoff in terms of maximum vector distance.
func Distance(d float64) *Similarity { return &Similarity{distance: &d} }

// Certainty sets a similarity cutoff in terms of certainty.
func Certainty(d float64) *Similarity { return &Similarity{certainty: &d} }

type NearVectorTarget api.NearVectorTarget

// NearVectorFunc runs plain near vector search.
type NearVectorFunc func(context.Context, NearVector) (*Result, error)

// GroupBy runs near vector search with a group by clause.
func (nvf NearVectorFunc) GroupBy(ctx context.Context, nv NearVector, groupBy GroupBy) (*GroupByResult, error) {
	nv.groupBy = &groupBy
	ctx = contextWithGroupByResult(ctx) // safe to reassign since we hold the copy of the original context.
	if _, err := nvf(ctx, nv); err != nil {
		return nil, err
	}
	return getGroupByResult(ctx), nil
}

func nearVector(ctx context.Context, t internal.Transport, rd api.RequestDefaults, nv NearVector) (*Result, error) {
	req := &api.SearchRequest{
		RequestDefaults:  rd,
		Limit:            nv.Limit,
		AutoLimit:        nv.AutoLimit,
		Offset:           nv.Offset,
		After:            nv.After,
		ReturnVectors:    nv.ReturnVectors,
		ReturnMetadata:   marshalReturnMetadata(nv.ReturnMetadata),
		ReturnProperties: marshalReturnProperties(nv.ReturnProperties, nv.ReturnNestedProperties),
		ReturnReferences: marshalReturnReferences(nv.ReturnReferences),
		NearVector: &api.NearVector{
			Target:    nv.Target,
			Distance:  nv.Similarity.distance,
			Certainty: nv.Similarity.certainty,
		},
	}

	var resp api.SearchResponse
	if err := t.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("near vector: %w", err)
	}

	// nearVector was called from the NearVectorFunc.GroupBy() method.
	// This means we should put GroupByResult in the context, as the first
	// return value will be discarded.
	if nv.groupBy != nil {
		var res GroupByResult
		groups := make(map[string]*Group[types.Map], len(resp.GroupByResults))
		objects := make([]*GroupByObject[types.Map], len(resp.GroupByResults))
		for name, group := range resp.GroupByResults {
			for _, obj := range group.Objects {
				unmarshaled := unmarshalObject(&obj.Object)
				dev.Assert(unmarshaled != nil, "nil object")
				objects = append(objects, &GroupByObject[types.Map]{
					BelongsToGroup: name,
					Object:         *unmarshaled,
				})
			}

			// Create a view into the Objects slice rather than allocating a separate one.
			from, to := len(objects)-len(group.Objects), len(objects)-1
			groups[name] = &Group[types.Map]{
				Name:        name,
				MinDistance: group.MinDistance,
				MaxDistance: group.MaxDistance,
				Size:        group.Size,
				Objects:     objects[from:to],
			}
		}
		res.Objects = objects
		setGroupByResult(ctx, &res)
		return nil, nil
	}

	var objects []*Object[types.Map]
	for _, obj := range resp.Results {
		objects = append(objects, unmarshalObject(obj))
	}
	return &Result{Objects: objects}, nil
}

// nearVectorFunc makes internal.Transport available to nearVector via a closure.
func nearVectorFunc(t internal.Transport, rd api.RequestDefaults) NearVectorFunc {
	return func(ctx context.Context, nv NearVector) (*Result, error) {
		return nearVector(ctx, t, rd, nv)
	}
}

// groupByKey is used to pass grouped query results to the GroupBy caller.
var groupByResultKey = internal.ContextKey{}

// contextWithGorupByResult creates a placeholder for *GroupByResult in the ctx.Values store.
func contextWithGroupByResult(ctx context.Context) context.Context {
	return internal.ContextWithPlaceholder[GroupByResult](ctx, groupByResultKey)
}

// getGroupByResult extracts *GroupByResult from the context.
func getGroupByResult(ctx context.Context) *GroupByResult {
	return internal.ValueFromContext[GroupByResult](ctx, groupByResultKey)
}

// setGroupByResult replaces *GroupByResult placeholder
// in the context with the value at r.
//
// We want to update the context passed to us in the request,
// rather than derive a new one. In the latter case the original
// context will stay unchanged and the caller will not see the value.
//
// Populating api.GroupByResult is NOT a part of the internal.Transport contract,
// but rather a responsibility of the layer using internal.ContextWithPlaceholder.
func setGroupByResult(ctx context.Context, r *GroupByResult) {
	internal.SetContextValue(ctx, groupByResultKey, r)
}
