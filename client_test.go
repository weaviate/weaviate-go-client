package weaviate_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/weaviate/weaviate-go-client/v5"
	"github.com/weaviate/weaviate-go-client/v5/query"
	"github.com/weaviate/weaviate-go-client/v5/types"
)

func TestClient(*testing.T) {
	c, err := weaviate.NewClient()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	// Create a collection and get a handle
	h := c.Collections.Create(ctx, "Songs")

	var single []float32
	var multi [][]float32

	h.Query.NearVector(ctx,
		types.Vector{Single: single},
		query.WithLimit(5),
		query.WithDistance(.34),
	)

	h.Query.NearVector(ctx, query.Average(
		types.Vector{Name: "title", Single: single},
		types.Vector{Name: "lyrics", Multi: multi},
	))

	res, _ := h.Query.NearVector(ctx, query.ManualWeights(
		query.Target(types.Vector{Name: "title", Single: single}, 0.4),
		query.Target(types.Vector{Name: "lyrics", Multi: multi}, 0.6),
	))

	fmt.Printf("got %d objects\n", len(res.Objects))
	for i, obj := range res.Objects {
		fmt.Printf("#%d: id=%s\n\tproperties=%v", obj.UUID, obj.Properties)
	}

	h.Query.NearVector.GroupBy(ctx,
		types.Vector{Single: single},
		"group by album",
		query.WithAutoLimit(2),
	)
}
