package query

import (
	"maps"

	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/types"
)

func NewClient(t api.SearchTransport, rd api.RequestDefaults) *Client {
	return &Client{
		transport:  t,
		defaults:   rd,
		NearVector: nearVectorFunc(t, rd),
	}
}

type Client struct {
	transport api.SearchTransport
	defaults  api.RequestDefaults

	NearVector NearVectorFunc
}

type NestedProperty struct {
	Name       string
	Properties []string
}

type Reference struct {
	PropertyName           string
	TargetCollection       string
	ReturnProperties       []string
	ReturnNestedProperties []NestedProperty // Return object properties and a subset of their nested properties.
	ReturnMetadata         []Metadata
}

type Metadata api.MetadataRequest

const (
	MetadataCreationTimeUnix   Metadata = Metadata(api.MetadataCreationTimeUnix)
	MetadataLastUpdateTimeUnix Metadata = Metadata(api.MetadataLastUpdateTimeUnix)
	MetadataDistance           Metadata = Metadata(api.MetadataDistance)
	MetadataCertainty          Metadata = Metadata(api.MetadataCertainty)
	MetadataScore              Metadata = Metadata(api.MetadataScore)
	MetadataExplainScore       Metadata = Metadata(api.MetadataExplainScore)
)

type GroupBy struct {
	Property       string // Property to group by.
	ObjectLimit    int    // Maximum number of objects per group.
	NumberOfGroups int    // Maximum number of groups to return.
}

type Result struct {
	Objects []Object[types.Map]
}

type Object[P types.Properties] struct {
	types.Object[P]
	Metadata QueryMetadata
}

type QueryMetadata struct {
	Distance     *float32
	Certainty    *float32
	Score        *float32
	ExplainScore *string
}

type Group[P types.Properties] struct {
	Name                     string
	MinDistance, MaxDistance float32
	Size                     int64
	Objects                  []GroupByObject[P]
}

type GroupByObject[P types.Properties] struct {
	Object[P]
	Metadata       QueryMetadata
	BelongsToGroup string
}

type GroupByResult struct {
	Objects []GroupByObject[types.Map]
	Groups  map[string]Group[types.Map]
}

type VectorTarget interface {
	Vectors() []api.TargetVector
}

func marshalSearchTarget(in VectorTarget) api.SearchTarget {
	out := api.SearchTarget{Vectors: in.Vectors()}
	if cm, ok := in.(interface{ CombinationMethod() api.CombinationMethod }); ok {
		out.CombinationMethod = cm.CombinationMethod()
	}
	return out
}

func marshalReturnMetadata(metadata []Metadata) []api.MetadataRequest {
	out := make([]api.MetadataRequest, 0, len(metadata)+1)
	for _, m := range metadata {
		out = append(out, api.MetadataRequest(m))
	}
	out = append(out, api.MetadataUUID)
	return out
}

func marshalReturnProperties(properties []string, nested []NestedProperty) []api.ReturnProperty {
	out := make([]api.ReturnProperty, 0, len(properties)+len(nested))
	for _, p := range properties {
		out = append(out, api.ReturnProperty{Name: p})
	}
	for _, np := range nested {
		out = append(out, api.ReturnProperty{
			Name:             np.Name,
			NestedProperties: np.Properties,
		})
	}
	return out
}

func marshalReturnReferences(references []Reference) []api.ReturnReference {
	out := make([]api.ReturnReference, len(references))
	for i, ref := range references {
		out[i] = api.ReturnReference{
			PropertyName:     ref.PropertyName,
			TargetCollection: ref.TargetCollection,
			ReturnMetadata:   marshalReturnMetadata(ref.ReturnMetadata),
			ReturnProperties: marshalReturnProperties(ref.ReturnProperties, ref.ReturnNestedProperties),
		}
	}
	return out
}

func unmarshalObject(in *api.Object) Object[types.Map] {
	vectors := make(types.Vectors, len(in.Metadata.NamedVectors)+1)
	maps.Copy(vectors, in.Metadata.NamedVectors)
	if in.Metadata.UnnamedVector != nil {
		vectors[api.DefaultVectorName] = *in.Metadata.UnnamedVector
	}

	// TODO(dyma): unmarshal references
	return Object[types.Map]{
		Object: types.Object[types.Map]{
			UUID:               in.Metadata.UUID,
			Vectors:            types.Vectors(vectors),
			Properties:         in.Properties,
			CreationTimeUnix:   in.Metadata.CreationTimeUnix,
			LastUpdateTimeUnix: in.Metadata.LastUpdateTimeUnix,
		},
		Metadata: QueryMetadata{
			Distance:     in.Metadata.Distance,
			Certainty:    in.Metadata.Certainty,
			Score:        in.Metadata.Score,
			ExplainScore: in.Metadata.ExplainScore,
		},
	}
}
