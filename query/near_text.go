package query

import (
	"context"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
)

type NearText struct {
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

	// Concepts are vectorized and used and the actual target for the similarity search.
	// Required parameter.
	Concepts []string

	// Bias the results towards or away from concepts and/or vectors.
	MoveTo, MoveAway *Move

	// Target vector or a combination of multiple vector targets.
	// By default, the resulting vectors are compared against the "default"
	// vector, or the _only_ vector, if the collection only has a single vector index.
	// See [MultiVectorTarget] for examples of providing multiple targets.
	Target VectorTarget

	// Similarity specifies a cutoff point for query results.
	// Use Distance() to set it as maximum distance between vectors.
	// Use Certainty() to set it to a normalized value between 0 and 1.
	// Prefer expressing Similarity in terms of vector distance, as that is a more conventional metric.
	Similarity *Similarity

	// groupBy can only be set by [NearTextFunc.GroupBy], as it changes the shape of the response.
	groupBy *GroupBy
}

type Move api.Move

// NearTextFunc runs plain near text search.
type NearTextFunc func(context.Context, NearText) (*Result, error)

// nearTextFunc makes internal.Transport available to nearText via a closure.
func nearTextFunc(t internal.Transport, rd api.RequestDefaults) NearTextFunc {
	return func(ctx context.Context, nv NearText) (*Result, error) {
		return nearText(ctx, t, rd, nv)
	}
}

func nearText(ctx context.Context, t internal.Transport, rd api.RequestDefaults, nt NearText) (*Result, error) {
	req := &api.SearchRequest{
		RequestDefaults:  rd,
		Limit:            nt.Limit,
		AutoLimit:        nt.AutoLimit,
		Offset:           nt.Offset,
		After:            nt.After,
		ReturnVectors:    nt.ReturnVectors,
		ReturnMetadata:   api.ReturnMetadata(nt.ReturnMetadata),
		ReturnProperties: marshalReturnProperties(nt.ReturnProperties, nt.ReturnNestedProperties),
		ReturnReferences: marshalReturnReferences(nt.ReturnReferences),
	}

	if nt.Target != nil {
		req.NearText = &api.NearText{
			Concepts:  nt.Concepts,
			Target:    marshalSearchTarget(nt.Target),
			Distance:  nt.Similarity.Distance(),
			Certainty: nt.Similarity.Certainty(),
			MoveTo:    (*api.Move)(nt.MoveTo),
			MoveAway:  (*api.Move)(nt.MoveAway),
		}
	}

	if nt.groupBy != nil {
		req.GroupBy = &api.GroupBy{
			Property:       nt.groupBy.Property,
			Limit:          nt.groupBy.ObjectLimit,
			NumberOfGroups: nt.groupBy.NumberOfGroups,
		}
	}

	var resp api.SearchResponse
	if err := t.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("near vector: %w", err)
	}

	// nearText was called from the NearTextFunc.GroupBy() method.
	// This means we should put GroupByResult in the context, as the first
	// return value will be discarded.
	if nt.groupBy != nil {
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
			Took:    resp.Took,
			Groups:  groups,
			Objects: objects,
		})
		return nil, nil
	}

	objects := make([]Object[map[string]any], len(resp.Results))
	for i, obj := range resp.Results {
		objects[i] = unmarshalObject(&obj)
	}
	return &Result{Took: resp.Took, Objects: objects}, nil
}

// GroupBy runs near vector search with a GroupBy clause.
func (nvf NearTextFunc) GroupBy(ctx context.Context, nv NearText, groupBy GroupBy) (*GroupByResult, error) {
	nv.groupBy = &groupBy
	ctx = contextWithGroupByResult(ctx) // safe to reassign since we hold the copy of the original context.
	if _, err := nvf(ctx, nv); err != nil {
		return nil, err
	}
	return getGroupByResult(ctx), nil
}
