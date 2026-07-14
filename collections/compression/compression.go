package compression

import "github.com/weaviate/weaviate-go-client/v6/internal"

// Registry stores all compression algorithms defined by this package.
var Registry internal.Modules[Type]

func init() {
	Registry.Register(*new(RQ))
}

type Type string

var _ internal.Module[Type] = (*RQ)(nil)

type RQ struct {
	Bits         int  `json:"bits"`
	RescoreLimit int  `json:"rescore_limit"`
	Cache        bool `json:"cache"`
}

func (RQ) Name() Type { return "rq" }
