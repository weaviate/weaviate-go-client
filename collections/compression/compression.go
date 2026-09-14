package compression

import "github.com/weaviate/weaviate-go-client/v6/internal"

// Registry stores all compression algorithms defined by this package.
var Registry internal.Modules[Type]

func init() {
	Registry.Register(*new(BQ))
	Registry.Register(*new(RQ))
}

type Type string

var (
	_ internal.Module[Type] = (*BQ)(nil)
	_ internal.Module[Type] = (*RQ)(nil)
)

type RQ struct {
	Bits         int  `json:"bits,omitempty"`
	RescoreLimit int  `json:"rescore_limit,omitempty"`
	Cache        bool `json:"cache,omitempty"`
}

func (RQ) Name() Type { return "rq" }

type BQ struct {
	RescoreLimit int  `json:"rescore_limit,omitempty"`
	Cache        bool `json:"cache,omitempty"`
}

func (BQ) Name() Type { return "bq" }
