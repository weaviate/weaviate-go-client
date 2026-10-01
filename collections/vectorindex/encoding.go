package vectorindex

import "github.com/weaviate/weaviate-go-client/v6/internal"

// encoderRegistry stores all encoding types defined by this package.
var encoderRegistry internal.Modules[Encoding]

func init() {
	encoderRegistry.Register(*new(Muvera))
}

type Encoding string

type Encoder internal.Module[Encoding]

var _ Encoder = (*Muvera)(nil)

type Muvera struct {
	KSim         int `json:"ksim,omitempty"`
	DProjections int `json:"dprojections,omitempty"`
	Repetitions  int `json:"repetitions,omitempty"`
}

func (Muvera) Name() Encoding { return "muvera" }
