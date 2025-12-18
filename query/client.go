package query

import (
	"maps"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/types"
)

func NewClient(t internal.Transport, rd api.RequestDefaults) *Client {
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

func unmarshalObject(in api.Object) Object[types.Map] {
	vectors := make(types.Vectors, len(in.Metadata.NamedVectors)+1)
	maps.Copy(vectors, in.Metadata.NamedVectors)
	vectors[api.DefaultVectorName] = in.Metadata.UnnamedVector

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
