package query

import (
	"context"

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

	Selection Selection

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

type (
	// Move adjusts the bias of the query results.
	Move api.Move
	MMR  api.SelectionMMR
)

// SelectionMMR creates [Selection] using [MMR] algorithm.
func SelectionMMR(mmr MMR) Selection { return Selection{mmr: &mmr} }

type Selection struct {
	mmr *MMR
}

func (s *Selection) MMR() *MMR { return s.mmr }

// NearTextFunc runs plain near text search.
type NearTextFunc func(context.Context, NearText) (*Result, error)

// nearTextFunc makes internal.Transport available to nearText via a closure.
func nearTextFunc(t internal.Transport, rd api.RequestDefaults) NearTextFunc {
	return func(ctx context.Context, nt NearText) (*Result, error) {
		return query(ctx, t, request{
			RequestDefaults:        rd,
			Limit:                  nt.Limit,
			AutoLimit:              nt.AutoLimit,
			Offset:                 nt.Offset,
			After:                  nt.After,
			ReturnVectors:          nt.ReturnVectors,
			ReturnMetadata:         nt.ReturnMetadata,
			ReturnProperties:       nt.ReturnProperties,
			ReturnNestedProperties: nt.ReturnNestedProperties,
			ReturnReferences:       nt.ReturnReferences,
			GroupBy:                nt.groupBy,
		}, func(req *api.SearchRequest) {
			req.NearText = &api.NearText{
				Concepts:  nt.Concepts,
				Distance:  nt.Similarity.Distance(),
				Certainty: nt.Similarity.Certainty(),
				MoveTo:    (*api.Move)(nt.MoveTo),
				MoveAway:  (*api.Move)(nt.MoveAway),
				Selection: api.Selection{
					MMR: (*api.SelectionMMR)(nt.Selection.MMR()),
				},
			}
			if nt.Target != nil {
				req.NearText.Target = marshalSearchTarget(nt.Target)
			}
		})
	}
}

// GroupBy runs near vector search with a GroupBy clause.
func (ntf NearTextFunc) GroupBy(ctx context.Context, nv NearText, groupBy GroupBy) (*GroupByResult, error) {
	nv.groupBy = &groupBy
	ctx = contextWithGroupByResult(ctx) // safe to reassign since we hold the copy of the original context.
	if _, err := ntf(ctx, nv); err != nil {
		return nil, err
	}
	return getGroupByResult(ctx), nil
}
