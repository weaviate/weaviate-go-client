package query

import (
	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
	"github.com/weaviate/weaviate-go-client/v6/types"
)

func NewClient(t internal.Transport, rd api.RequestDefaults) *Client {
	dev.AssertNotNil(t, "t")

	return &Client{
		transport:  t,
		defaults:   rd,
		NearVector: nearVectorFunc(t, rd),
	}
}

type Client struct {
	transport internal.Transport
	defaults  api.RequestDefaults

	NearVector NearVectorFunc
}

type ReturnMetadata api.ReturnMetadata

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

type Result struct {
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

func marshalReturnReferences(references []Reference) []api.ReturnReference {
	out := make([]api.ReturnReference, len(references))
	for i, ref := range references {
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

func unmarshalObject(o *api.Object) Object[map[string]any] {
	vectors := make(types.Vectors, len(o.Metadata.Vectors))
	for k, v := range o.Metadata.Vectors {
		vectors[k] = types.Vector(v)
	}

	references := make(map[string][]types.Object[map[string]any], len(o.References))
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
