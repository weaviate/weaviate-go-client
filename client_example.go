package weaviate

import (
	"github.com/weaviate/weaviate-go-client/v5/query"
	"github.com/weaviate/weaviate-go-client/v5/util"
)

func ExampleClient() {
	c, err := NewClient()
	if err != nil {
		panic(err)
	}

	c.Query.NearVector(
		util.Vector{Single: []float32{}},
		query.WithLimit(5),
		query.WithDistance(.34),
	)

	var single []float32
	var multi [][]float32

	c.Query.NearVector(query.Average(
		util.Vector{Name: "title", Single: single},
		util.Vector{Name: "lyrics", Multi: multi},
	))

	c.Query.NearVector(query.ManualWeights(
		query.Target(util.Vector{Name: "title", Single: single}, 0.4),
		query.Target(util.Vector{Name: "lyrics", Multi: multi}, 0.6),
	))

	c.Query.NearVector.GroupBy(
		util.Vector{Single: single},
		"group by album",
	)
}
