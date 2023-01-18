package graphql

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
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
	Description string `json:"description"`
	Additional  struct {
		ID                 string    `json:"id"`
		Vector             []float32 `json:"vector"`
		CreationTimeUnix   string    `json:"creationTimeUnix"`
		LastUpdateTimeUnix string    `json:"lastUpdateTimeUnix"`
	} `json:"_additional"`
}

func GetOnePizza(t *testing.T, c *weaviate.Client) *Pizza {
	_additional := graphql.Field{
		Name: "_additional",
		Fields: []graphql.Field{
			{Name: "id"},
			{Name: "vector"},
			{Name: "creationTimeUnix"},
			{Name: "lastUpdateTimeUnix"},
		},
	}

	resp, err := c.GraphQL().
		Get().
		WithClassName("Pizza").
		WithFields(_additional).
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
