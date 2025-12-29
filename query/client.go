package query

import "github.com/weaviate/weaviate-go-client/v6/types"

type Client struct {
	gRPC any // gRPCClient

	collectionName string
	NearVector     NearVectorFunc
}

func NewClient(gRPC any, collectionName string) *Client {
	return &Client{
		gRPC:       gRPC,
		NearVector: nearVectorFunc(gRPC),
	}
}

type commonOptions struct {
	Limit            *int
	Offset           *int
	AutoLimit        *int
	After            *string
	ReturnProperties []string
	ReturnMetadata   []Metadata
	IncludeVectors   []string
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

// returnPropertiesOption selects properties to include in the response.
type returnPropertiesOption []string

// Compile-time assertion that returnPropertiesOption implements NearVectorOption.
var _ NearVectorOption = (*returnPropertiesOption)(nil)

// WithReturnProperties selects properties to include in the response.
// By default, all properties are returned.
func WithReturnProperties(properties ...string) returnPropertiesOption {
	return returnPropertiesOption(properties)
}

type returnMetadataOption []Metadata

// Compile-time assertion that returnMetadataOption implements NearVectorOption.
var _ NearVectorOption = (*returnMetadataOption)(nil)

type Metadata string

const (
	MetadataCreationTimeUnix   Metadata = "CreationTimeUnix"
	MetadataLastUpdateTimeUnix Metadata = "LastUpdateTimeUnix"
	MetadataDistance           Metadata = "Distance"
	MetadataCertainty          Metadata = "Certainty"
	MetadataScore              Metadata = "Score"
	MetadataExplainScore       Metadata = "ExplainScore"
)

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

type GroupByResult[P types.Properties] struct {
	Objects []GroupByObject[P]
	Groups  map[string]Group[P]
}
