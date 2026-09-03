package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/query/boost"
	"github.com/weaviate/weaviate-go-client/v6/query/filter"
)

type NearMedia struct {
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

	// Media is represented as a base64 string.
	Media Media

	// Target vector or a combination of multiple vector targets.
	// By default, the resulting vectors are compared against the "default"
	// vector, or the _only_ vector, if the collection only has a single vector index.
	// See [MultiVectorTarget] for examples of providing multiple targets.
	Target VectorTarget

	// Similarity specifies a cutoff point for query results.
	// Prefer expressing similarity in terms of vector distance, as that is a more conventional metric.
	Similarity VectorSimilarity

	// groupBy can only be set by [NearMediaFunc.GroupBy], as it changes the shape of the response.
	groupBy *GroupBy
}

type Media struct {
	data string
	kind api.MediaKind
}

func (m Media) Kind() api.MediaKind { return m.kind }
func (m Media) Data() string        { return m.data }

// Image returns near image query target. The value must be a base64-encoded string.
func Image(s string) Media { return Media{kind: api.MediaImage, data: s} }

// Audio returns near audio query target. The value must be a base64-encoded string.
func Audio(s string) Media { return Media{kind: api.MediaAudio, data: s} }

// Video returns near video query target. The value must be a base64-encoded string.
func Video(s string) Media { return Media{kind: api.MediaVideo, data: s} }

// Depth returns near depth query target. The value must be a base64-encoded string.
func Depth(s string) Media { return Media{kind: api.MediaDepth, data: s} }

// Thermal returns near thermal query target. The value must be a base64-encoded string.
func Thermal(s string) Media { return Media{kind: api.MediaThermal, data: s} }

// IMU returns near IMU query target. The value must be a base64-encoded string.
func IMU(s string) Media { return Media{kind: api.MediaIMU, data: s} }

// NearMediaFunc runs plain near vector search.
type NearMediaFunc func(context.Context, NearMedia) (*Result, error)

// nearMediaFunc makes internal.Transport available to [query] via a closure.
func nearMediaFunc(t internal.Transport, rd api.RequestDefaults) NearMediaFunc {
	return func(ctx context.Context, nm NearMedia) (*Result, error) {
		return query(ctx, t, request{
			RequestDefaults:        rd,
			Limit:                  int32(nm.Limit),
			AutoLimit:              int32(nm.AutoLimit),
			Offset:                 int32(nm.Offset),
			After:                  nm.After,
			Filter:                 nm.Filter,
			Boost:                  nm.Boost,
			ReturnVectors:          nm.ReturnVectors,
			ReturnMetadata:         nm.ReturnMetadata,
			ReturnProperties:       nm.ReturnProperties,
			ReturnNestedProperties: nm.ReturnNestedProperties,
			ReturnReferences:       nm.ReturnReferences,
			GroupBy:                nm.groupBy,
		}, func(req *api.SearchRequest) { req.NearMedia = nearMedia(&nm) })
	}
}

// nearMedia convers [NearMedia] to [api.NearMedia].
func nearMedia(nm *NearMedia) *api.NearMedia {
	if nm == nil {
		return nil
	}

	out := &api.NearMedia{
		Kind:  nm.Media.Kind(),
		Media: nm.Media.Data(),
		Similarity: api.VectorSimilarity{
			Distance:  nm.Similarity.Distance(),
			Certainty: nm.Similarity.Certainty(),
		},
	}

	if nm.Target != nil {
		out.Target = marshalSearchTarget(nm.Target)
	}

	return out
}

// GroupBy runs near vector search with a GroupBy clause.
func (nmf NearMediaFunc) GroupBy(ctx context.Context, nv NearMedia, groupBy GroupBy) (*GroupByResult, error) {
	nv.groupBy = &groupBy
	return queryGroupBy(ctx, nmf, nv)
}

func (nm NearMedia) Search() *api.NearMedia { return nearMedia(&nm) }
