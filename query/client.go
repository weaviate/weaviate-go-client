package query

import (
	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/request"
	"github.com/weaviate/weaviate-go-client/v6/types"
)

func NewClient(t internal.Transport, rd request.Defaults) *Client {
	return &Client{
		transport:  t,
		defaults:   rd,
		NearVector: nearVectorFunc(t),
	}
}

type Client struct {
	transport internal.Transport
	defaults  request.Defaults

	NearVector NearVectorFunc
}

type commonOptions struct {
	Limit            *int
	Offset           *int
	AutoLimit        *int
	After            *string
	ReturnProperties []string
	IncludeVectors   []string
	GroupBy          *GroupBy
}

// LimitOption sets the `limit` parameter.
type LimitOption int

var _ NearVectorOption = (*LimitOption)(nil)

func WithLimit(l int) LimitOption {
	return LimitOption(l)
}

// OffsetOption sets the `limit` parameter.
type OffsetOption int

var _ NearVectorOption = (*OffsetOption)(nil)

func WithOffset(l int) OffsetOption {
	return OffsetOption(l)
}

// AutoLimitOption sets the `limit` parameter.
type AutoLimitOption int

var _ NearVectorOption = (*AutoLimitOption)(nil)

func WithAutoLimit(l int) AutoLimitOption {
	return AutoLimitOption(l)
}

// TODO: define GroupBy parameters
type GroupBy struct {
	Property string
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
	// Should these be pointers? *float32
	Distance     float32
	Certainty    float32
	Score        float32
	ExplainScore string
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
