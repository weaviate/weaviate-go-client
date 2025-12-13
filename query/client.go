package query

import (
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

// commonOptions are parameters applicable to all search types.
// Concrete request structs should embed this struct.
type commonOptions struct {
	Limit            *int
	Offset           *int
	AutoLimit        *int
	After            *string
	ReturnProperties []api.ReturnProperty
	ReturnReferences []api.ReturnReference // TODO: add functional option
	ReturnVectors    []string
	ReturnMetadata   []api.Metadata
	GroupBy          *GroupBy
}

// WithLimit sets the `limit` parameter.
type WithLimit int

// Compile-time assertion that WithLimit implements NearVectorOption.
var _ NearVectorOption = (*WithLimit)(nil)

// WithOffset sets the `offset` parameter.
type WithOffset int

// Compile-time assertion that WithOffset implements NearVectorOption.
var _ NearVectorOption = (*WithOffset)(nil)

// WithAutoLimit sets the `autocut` parameter.
type WithAutoLimit int

// Compile-time assertion that WithAutoLimit implements NearVectorOption.
var _ NearVectorOption = (*WithAutoLimit)(nil)

// WithAfter sets the `after` parameter.
type WithAfter string

// Compile-time assertion that WithAfter implements NearVectorOption.
var _ NearVectorOption = (*WithAfter)(nil)

type ReturnVectorsOption []string

// Compile-time assertion that ReturnVectorsOption implements NearVectorOption.
var _ NearVectorOption = (*ReturnVectorsOption)(nil)

func WithReturnVector(vectors ...string) ReturnVectorsOption {
	if len(vectors) == 0 {
		return api.ReturnOnlyVector
	}
	return ReturnVectorsOption(vectors)
}

// returnPropertiesOption selects properties to include in the response.
type returnPropertiesOption []api.ReturnProperty

// Compile-time assertion that returnPropertiesOption implements NearVectorOption.
var _ NearVectorOption = (*returnPropertiesOption)(nil)

// WithReturnProperties selects properties to include in the response.
// By default, all properties are returned.
func WithReturnProperties(properties ...string) returnPropertiesOption {
	out := make(returnPropertiesOption, len(properties))
	for _, p := range properties {
		out = append(out, api.ReturnProperty{Name: p})
	}
	return out
}

// WithReturnNestedProperties selects properties to include in the response.
// By default, all properties are returned.
func WithReturnNestedProperties(propertyName string, nestedProperties ...string) returnPropertiesOption {
	return returnPropertiesOption{{Name: propertyName, NestedProperties: nestedProperties}}
}

type returnMetadataOption []api.Metadata

// Compile-time assertion that returnMetadataOption implements NearVectorOption.
var _ NearVectorOption = (*returnMetadataOption)(nil)

type Metadata api.Metadata

const (
	MetadataCreationTimeUnix   Metadata = Metadata(api.MetadataCreationTimeUnix)
	MetadataLastUpdateTimeUnix Metadata = Metadata(api.MetadataLastUpdateTimeUnix)
	MetadataDistance           Metadata = Metadata(api.MetadataDistance)
	MetadataCertainty          Metadata = Metadata(api.MetadataCertainty)
	MetadataScore              Metadata = Metadata(api.MetadataScore)
	MetadataExplainScore       Metadata = Metadata(api.MetadataExplainScore)
)

func WithReturnMetadataOption(metadata ...Metadata) returnMetadataOption {
	out := make(returnMetadataOption, len(metadata))
	for _, m := range metadata {
		out = append(out, api.Metadata(m))
	}
	return out
}

// TODO: define GroupBy parameters
type GroupBy struct {
	Property       string // Property to group by.
	ObjectLimit    int    // Maximum number of objects per group.
	NumberOfGroups int    // Maximum number of groups to return.
}

// groupByOption is used internally to support grouped queries.
type groupByOption GroupBy

var _ NearVectorOption = (*groupByOption)(nil)

func withGroupBy(property string) groupByOption {
	return groupByOption(GroupBy{Property: property})
}

type Result struct {
	Objects []types.Object[types.Map]
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
	types.Object[P]
	Metadata       QueryMetadata
	BelongsToGroup string
}

type GroupByResult struct {
	Objects []GroupByObject[types.Map]
	Groups  map[string]Group[types.Map]
}
