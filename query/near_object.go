package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/query/boost"
	"github.com/weaviate/weaviate-go-client/v6/query/filter"
)

type NearObject struct {
	Limit                  int              // Limit the number of results returned for the query.
	Offset                 int              // Skip the first N objects in the collection.
	AutoLimit              int              // Return objects in the first N similarity clusters.
	After                  uuid.UUID        // Skip all objects before the one with this ID.
	Filter                 filter.Expr      // Filter results based on their properties.
	Boost                  boost.Expr       // Rerank search results using a decay function.
	ReturnMetadata         ReturnMetadata   // Select query and object metadata to return for each object.
	ReturnVectors          []string         // List vectors to return for each object in the result set.
	ReturnReferences       []Reference      // Select reference properties to return.
	ReturnNestedProperties []NestedProperty // Return object properties and a subset of their nested properties.

	// Select a subset of properties to return. By default, all properties are returned.
	// To not return any properties, initialize this value to an empty slice explicitly.
	ReturnProperties []string

	// Vectors of this object will be used as the search target.
	UUID uuid.UUID

	// ExcludeSelf removes the target object from the result set.
	ExcludeSelf bool

	// Target vector or a combination of multiple vector targets.
	// By default, the resulting vectors are compared against the "default"
	// vector, or the _only_ vector, if the collection only has a single vector index.
	// See [MultiVectorTarget] for examples of providing multiple targets.
	Target VectorTarget

	// Similarity specifies a cutoff point for query results.
	// Prefer expressing similarity in terms of vector distance, as that is a more conventional metric.
	Similarity VectorSimilarity

	// groupBy can only be set by [NearObjectFunc.GroupBy], as it changes the shape of the response.
	groupBy *GroupBy
}

// NearObjectFunc runs plain near object search.
type NearObjectFunc func(context.Context, NearObject) (*Result, error)

// nearObjectFunc makes internal.Transport available to [query] via a closure.
func nearObjectFunc(t internal.Transport, rd api.RequestDefaults) NearObjectFunc {
	return func(ctx context.Context, no NearObject) (*Result, error) {
		// TODO(dyma): make sure this is also covered in aggregation
		if no.ExcludeSelf {
			exclude := filter.Not{
				P: &filter.Cond{
					Target:   filter.UUID,
					Operator: filter.Equal,
					Value:    no.UUID,
				},
			}
			if no.Filter == nil {
				no.Filter = exclude
			} else {
				no.Filter = filter.And{no.Filter, exclude}
			}
		}

		return query(ctx, t, request{
			RequestDefaults:        rd,
			Limit:                  int32(no.Limit),
			AutoLimit:              int32(no.AutoLimit),
			Offset:                 int32(no.Offset),
			After:                  no.After,
			Filter:                 no.Filter,
			Boost:                  no.Boost,
			ReturnVectors:          no.ReturnVectors,
			ReturnMetadata:         no.ReturnMetadata,
			ReturnProperties:       no.ReturnProperties,
			ReturnNestedProperties: no.ReturnNestedProperties,
			ReturnReferences:       no.ReturnReferences,
			GroupBy:                no.groupBy,
		}, func(req *api.SearchRequest) { req.NearObject = nearObject(&no) })
	}
}

// nearObject convers [NearObject] to [api.NearObject].
func nearObject(no *NearObject) *api.NearObject {
	if no == nil {
		return nil
	}

	out := &api.NearObject{
		UUID: no.UUID,
		Similarity: api.VectorSimilarity{
			Distance:  no.Similarity.Distance(),
			Certainty: no.Similarity.Certainty(),
		},
	}

	if no.Target != nil {
		out.Target = marshalSearchTarget(no.Target)
	}

	return out
}

// GroupBy runs near object search with a GroupBy clause.
func (nof NearObjectFunc) GroupBy(ctx context.Context, no NearObject, groupBy GroupBy) (*GroupByResult, error) {
	no.groupBy = &groupBy
	return queryGroupBy(ctx, nof, no)
}

func (no NearObject) Search() *api.NearObject { return nearObject(&no) }
