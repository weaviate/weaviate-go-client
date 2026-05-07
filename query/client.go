package query

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
	"github.com/weaviate/weaviate-go-client/v6/query/filter"
	"github.com/weaviate/weaviate-go-client/v6/types"
)

func NewClient(t internal.Transport, rd api.RequestDefaults) *Client {
	dev.AssertNotNil(t, "transport")

	return &Client{
		transport:  t,
		defaults:   rd,
		NearVector: nearVectorFunc(t, rd),
		NearText:   nearTextFunc(t, rd),
		Hybrid:     hybridFunc(t, rd),
	}
}

type Client struct {
	transport internal.Transport
	defaults  api.RequestDefaults

	NearVector NearVectorFunc
	NearText   NearTextFunc
	Hybrid     HybridFunc
}

type (
	ReturnMetadata api.ReturnMetadata
	GroupBy        struct {
		Property       string // Property to group by.
		ObjectLimit    int32  // Maximum number of objects per group.
		NumberOfGroups int32  // Maximum number of groups to return.
	}
)

type NestedProperty struct {
	Name                   string
	ReturnProperties       []string
	ReturnNestedProperties []NestedProperty
}

type Reference struct {
	PropertyName     string // Name of the reference property. Required.
	TargetCollection string // Target collection. Required for multi-target reference properties.

	ReturnMetadata         ReturnMetadata   // Select object metadata to return for each reference.
	ReturnVectors          []string         // List vectors to return for each reference in the result set.
	ReturnReferences       []Reference      // Select reference properties to return.
	ReturnNestedProperties []NestedProperty // Return object properties and a subset of their nested properties.

	// Select a subset of properties to return. By default, all properties are returned.
	// To not return any properties, initialize this value to an empty slice explicitly.
	ReturnProperties []string
}

// VectorTarget can be used as an input for a vector similarity search.
type VectorTarget interface {
	// Vectors returns vectors included in the search target.
	Vectors() []api.TargetVector
}

type Result struct {
	Took    time.Duration
	Objects []Object[map[string]any]
}

type Object[T any] struct {
	types.Object[T]
	Metadata Metadata
}

type Metadata struct {
	Distance     *float32 // Distance from the search vector. Nil if not requested via [ReturnMetadata].
	Certainty    *float32 // Normalized distance metric. Nil if not requested via [ReturnMetadata].
	Score        *float32
	ExplainScore *string
}

type GroupByResult struct {
	Took    time.Duration
	Objects []GroupObject[map[string]any]
	Groups  map[string]Group[map[string]any]
}

type Group[T any] struct {
	Name                     string
	MinDistance, MaxDistance float32
	Size                     int64
	Objects                  []GroupObject[T]
}

type GroupObject[T any] struct {
	Object[T]
	BelongsToGroup string
}

type request struct {
	api.RequestDefaults
	Limit                  int32
	Offset                 int32
	AutoLimit              int32
	After                  uuid.UUID
	Filter                 filter.Expr
	ReturnMetadata         ReturnMetadata
	ReturnVectors          []string
	ReturnReferences       []Reference
	ReturnNestedProperties []NestedProperty
	ReturnProperties       []string
	GroupBy                *GroupBy
}

func query(ctx context.Context, t internal.Transport, r request, f func(*api.SearchRequest)) (*Result, error) {
	req := &api.SearchRequest{
		RequestDefaults:  r.RequestDefaults,
		Limit:            r.Limit,
		AutoLimit:        r.AutoLimit,
		Offset:           r.Offset,
		After:            r.After,
		Filter:           marshalFilter(r.Filter),
		ReturnVectors:    r.ReturnVectors,
		ReturnMetadata:   api.ReturnMetadata(r.ReturnMetadata),
		ReturnProperties: marshalReturnProperties(r.ReturnProperties, r.ReturnNestedProperties),
		ReturnReferences: marshalReturnReferences(r.ReturnReferences),
	}

	f(req)

	if r.GroupBy != nil {
		req.GroupBy = &api.GroupBy{
			Property:       r.GroupBy.Property,
			Limit:          r.GroupBy.ObjectLimit,
			NumberOfGroups: r.GroupBy.NumberOfGroups,
		}
	}

	var resp api.SearchResponse
	if err := t.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("near vector: %w", err)
	}

	// query was called from the GroupBy() method. This means we should put
	// GroupByResult in the context, as the first return value will be discarded.
	if r.GroupBy != nil {
		groups := internal.MakeMap[string, Group[map[string]any]](len(resp.GroupByResults))
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
		internal.SetContextValue(ctx, groupByResultKey, &GroupByResult{
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

func marshalReturnProperties(properties []string, nested []NestedProperty) []api.ReturnProperty {
	if len(properties)+len(nested) == 0 {
		return nil
	}

	out := make([]api.ReturnProperty, len(properties)+len(nested))
	for i, p := range properties {
		out[i] = api.ReturnProperty{Name: p}
	}

	for i, np := range nested {
		out[i+len(properties)] = api.ReturnProperty{
			Name:             np.Name,
			NestedProperties: marshalReturnProperties(np.ReturnProperties, np.ReturnNestedProperties),
		}
	}
	return out
}

func marshalReturnReferences(in []Reference) []api.ReturnReference {
	if len(in) == 0 {
		return nil
	}
	out := make([]api.ReturnReference, len(in))
	for i, ref := range in {
		out[i] = api.ReturnReference{
			PropertyName:     ref.PropertyName,
			TargetCollection: ref.TargetCollection,
			ReturnVectors:    ref.ReturnVectors,
			ReturnMetadata:   api.ReturnMetadata(ref.ReturnMetadata),
			ReturnProperties: marshalReturnProperties(ref.ReturnProperties, ref.ReturnNestedProperties),
		}
	}
	return out
}

func marshalFilter(expr filter.Expr) api.Filter {
	if expr == nil {
		return api.Filter{}
	}
	f := api.Filter{
		Operator: expr.Operator(),
		Target:   expr.Target(),
		Value:    expr.Value(),
	}

	for _, e := range expr.Exprs() {
		f.Exprs = append(f.Exprs, marshalFilter(e))
	}

	return f
}

func marshalSearchTarget(target VectorTarget) api.SearchTarget {
	dev.AssertNotNil(target, "target")

	out := api.SearchTarget{Vectors: target.Vectors()}
	if cm, ok := target.(interface{ CombinationMethod() api.CombinationMethod }); ok {
		out.CombinationMethod = cm.CombinationMethod()
	}
	return out
}

func unmarshalObject(o *api.Object) Object[map[string]any] {
	vectors := internal.MakeMap[string, types.Vector](len(o.Metadata.Vectors))
	for k, v := range o.Metadata.Vectors {
		vectors[k] = types.Vector(v)
	}

	references := internal.MakeMap[string, []types.Object[map[string]any]](len(o.References))
	for k, refs := range o.References {
		objects := make([]types.Object[map[string]any], len(refs))
		for i, ref := range refs {
			objects[i] = unmarshalObject(&ref).Object
		}
		references[k] = objects
	}

	return Object[map[string]any]{
		Object: types.Object[map[string]any]{
			UUID:          o.Metadata.UUID,
			Collection:    o.Collection,
			Properties:    o.Properties,
			Vectors:       vectors,
			References:    references,
			CreatedAt:     o.Metadata.CreatedAt,
			LastUpdatedAt: o.Metadata.LastUpdatedAt,
		},
		Metadata: Metadata{
			Distance:     o.Metadata.Distance,
			Certainty:    o.Metadata.Certainty,
			Score:        o.Metadata.Score,
			ExplainScore: o.Metadata.ExplainScore,
		},
	}
}

// groupByResultKey is used to pass grouped query results to the GroupBy caller.
var groupByResultKey = internal.ContextKey{}

func queryGroupBy[In any](ctx context.Context, f func(context.Context, In) (*Result, error), in In) (*GroupByResult, error) {
	ctx = internal.ContextWithPlaceholder[GroupByResult](ctx, groupByResultKey)
	if _, err := f(ctx, in); err != nil {
		return nil, err
	}
	return internal.ValueFromContext[GroupByResult](ctx, groupByResultKey), nil
}
