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

	fmt.Printf("NearVector: got %d objects\n", len(res.Objects))
	for i, obj := range res.Objects {
		fmt.Printf("#%d: id=%s\n\tproperties=%v\n", i, obj.UUID, obj.Properties)
	}

	grouped, _ := h.Query.NearVector.GroupBy(ctx,
		types.Vector{Single: single},
		"group by album",
		query.WithAutoLimit(2),
	)

	fmt.Printf("NearVector.GroupBy: got %d objects in %d groups\n", len(grouped.Objects), len(grouped.Groups))
	for key, group := range grouped.Groups {
		fmt.Printf("Group %q has %d objects \n", key, group.Size)
	}
	// grouped.Objects[0].Properties
}
