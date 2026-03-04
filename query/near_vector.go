package query

import (
	"context"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
)

type NearVector struct {
	Limit                  int32            // Limit the number of results returned for the query.
	Offset                 int32            // Skip the first N objects in the collection.
	AutoLimit              int32            // Return objects in the first N similarity clusters.
	After                  uuid.UUID        // Skip all objects before the one with this ID.
	ReturnMetadata         ReturnMetadata   // Select query and object metadata to return for each object.
	ReturnVectors          []string         // List vectors to return for each object in the result set.
	ReturnReferences       []Reference      // Select reference properties to return.
	ReturnNestedProperties []NestedProperty // Return object properties and a subset of their nested properties.

	// Select a subset of properties to return. By default, all properties are returned.
	// To not return any properties, initialize this value to an empty slice explicitly.
	ReturnProperties []string

	// Target vector or a combination of multiple vector targets. Required parameter.
	// See [MultiVectorTarget] for examples of providing multiple targets.
	Target VectorTarget

	// Similarity specifies a cutoff point for query results.
	// Use Distance() to set it as maximum distance between vectors.
	// Use Certainty() to set it to a normalized value between 0 and 1.
	// Prefer expressing Similarity in terms of vector distance, as that is a more conventional metric.
	Similarity *Similarity

	// groupBy can only be set by [NearVectorFunc.GroupBy], as it changes the shape of the response.
	groupBy *GroupBy
}

// Distance sets a similarity cutoff in terms of maximum vector distance.
func Distance(d float64) *Similarity { return &Similarity{distance: &d} }

// Certainty sets a similarity cutoff in terms of certainty.
func Certainty(c float64) *Similarity { return &Similarity{certainty: &c} }

// NearVectorFunc runs plain near vector search.
type NearVectorFunc func(context.Context, NearVector) (*Result, error)

// nearVectorFunc makes internal.Transport available to nearVector via a closure.
func nearVectorFunc(t internal.Transport, rd api.RequestDefaults) NearVectorFunc {
	return func(ctx context.Context, nv NearVector) (*Result, error) {
		return nearVector(ctx, t, rd, nv)
	}
}

func nearVector(ctx context.Context, t internal.Transport, rd api.RequestDefaults, nv NearVector) (*Result, error) {
	req := &api.SearchRequest{
		RequestDefaults:  rd,
		Limit:            nv.Limit,
		AutoLimit:        nv.AutoLimit,
		Offset:           nv.Offset,
		After:            nv.After,
		ReturnVectors:    nv.ReturnVectors,
		ReturnMetadata:   api.ReturnMetadata(nv.ReturnMetadata),
		ReturnProperties: marshalReturnProperties(nv.ReturnProperties, nv.ReturnNestedProperties),
		ReturnReferences: marshalReturnReferences(nv.ReturnReferences),
	}

	if nv.Target != nil {
		req.NearVector = &api.NearVector{
			Target:    marshalSearchTarget(nv.Target),
			Distance:  nv.Similarity.Distance(),
			Certainty: nv.Similarity.Certainty(),
		}
	}

	if nv.groupBy != nil {
		req.GroupBy = (*api.GroupBy)(nv.groupBy)
	}

	var resp api.SearchResponse
	if err := t.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("near vector: %w", err)
	}

	// nearVector was called from the NearVectorFunc.GroupBy() method.
	// This means we should put GroupByResult in the context, as the first
	// return value will be discarded.
	if nv.groupBy != nil {
		groups := make(map[string]Group[map[string]any], len(resp.GroupByResults))
		objects := make([]GroupObject[map[string]any], 0)
		for _, group := range resp.GroupByResults {

			objects = slices.Grow(objects, len(group.Objects))
			for _, obj := range group.Objects {
				objects = append(objects, GroupObject[map[string]any]{
					BelongsToGroup: group.Name,
					Object:         unmarshalObject(&obj.Object),
				})
			}

			// Create a view into the objects slice rather than allocating a separate one.
			from, to := len(objects)-len(group.Objects), len(objects)
			groups[group.Name] = Group[map[string]any]{
				Name:        group.Name,
				MinDistance: group.MinDistance,
				MaxDistance: group.MaxDistance,
				Size:        group.Size,
				Objects:     objects[from:to],
			}
		}
		setGroupByResult(ctx, &GroupByResult{
			Groups:  groups,
			Objects: objects,
		})
		return nil, nil
	}

	objects := make([]Object[map[string]any], len(resp.Results))
	for i, obj := range resp.Results {
		objects[i] = unmarshalObject(&obj)
	}
	return &Result{Objects: objects}, nil
}

// GroupBy runs near vector search with a GroupBy clause.
func (nvf NearVectorFunc) GroupBy(ctx context.Context, nv NearVector, groupBy GroupBy) (*GroupByResult, error) {
	nv.groupBy = &groupBy
	ctx = contextWithGroupByResult(ctx) // safe to reassign since we hold the copy of the original context.
	if _, err := nvf(ctx, nv); err != nil {
		return nil, err
	}
	return getGroupByResult(ctx), nil
}

// Similarity is a cutoff point for query results.
type Similarity struct{ distance, certainty *float64 }

func (s *Similarity) Distance() *float64 {
	if s == nil {
		return nil
	}
	return s.distance
}

func (s *Similarity) Certainty() *float64 {
	if s == nil {
		return nil
	}
	return s.certainty
}
