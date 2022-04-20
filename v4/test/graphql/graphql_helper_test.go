package graphql

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/semi-technologies/weaviate-go-client/v4/weaviate"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/graphql"
	"github.com/stretchr/testify/require"
)

type GetPizzaResponse struct {
	Get struct {
		Pizzas []Pizza `json:"Pizza"`
	}
}

type AggregatePizzaResponse struct {
	Aggregate struct {
		Pizza []struct {
			Meta struct {
				Count int `json:"count"`
			} `json:"meta"`
		}
	}
}

type Pizza struct {
	Additional struct {
		ID     string    `json:"id"`
		Vector []float32 `json:"vector"`
	} `json:"_additional"`
}

func GetOnePizza(t *testing.T, c *weaviate.Client) *Pizza {
	fields := []graphql.Field{
		{
			Name: "_additional { id vector }",
		},
	}

	resp, err := c.GraphQL().
		Get().
		WithClassName("Pizza").
		WithFields(fields).
		Do(context.Background())
	if err != nil {
		t.Fatalf("failed to get an object: %s", err)
	}

	b, err := json.Marshal(resp.Data)
	require.Nil(t, err)

	var pizza GetPizzaResponse
	err = json.Unmarshal(b, &pizza)

	require.Nil(t, err)
	require.NotEmpty(t, pizza.Get.Pizzas)

	return &pizza.Get.Pizzas[0]
}
