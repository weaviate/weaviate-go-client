package aggregate

import (
	"time"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
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

type Text struct {
	Property string

	Count               bool
	TopOccurrences      bool
	TopOccurencesCutoff int32
}

type Integer struct {
	Property string

	Count  bool
	Sum    bool
	Min    bool
	Max    bool
	Mode   bool
	Mean   bool
	Median bool
}

type Number struct {
	Property string

	Count  bool
	Sum    bool
	Min    bool
	Max    bool
	Mode   bool
	Mean   bool
	Median bool
}

type Boolean struct {
	Property string

	Count           bool
	PercentageTrue  bool
	PercentageFalse bool
	TotalTrue       bool
	TotalFalse      bool
}

type Date struct {
	Property string

	Count  bool
	Min    bool
	Max    bool
	Mode   bool
	Median bool
}

type GroupBy struct {
	Collection string
	Property   string
}

type Result struct {
	Text    map[string]TextResult
	Integer map[string]IntegerResult
	Number  map[string]NumberResult
	Boolean map[string]BooleanResult
	Date    map[string]DateResult

	TotalCount  *int
	TookSeconds float32
}

type TextResult struct {
	Count          *int64
	TopOccurrences []TopOccurrence
}

type TopOccurrence struct {
	Value       string
	OccursTimes int64
}

type IntegerResult struct {
	Count  *int64
	Sum    *int64
	Min    *int64
	Max    *int64
	Mode   *int64
	Mean   *float64
	Median *float64
}

type NumberResult struct {
	Count  *int64
	Sum    *float64
	Min    *float64
	Max    *float64
	Mode   *float64
	Mean   *float64
	Median *float64
}

type BooleanResult struct {
	Count           *int64
	PercentageTrue  *float64
	PercentageFalse *float64
	TotalTrue       *int64
	TotalFalse      *int64
}

type DateResult struct {
	Count  *int64
	Min    *time.Time
	Max    *time.Time
	Mode   *time.Time
	Median *time.Time
}
