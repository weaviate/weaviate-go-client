package graphql

import (
	"context"
	"fmt"
	"testing"

	"github.com/semi-technologies/weaviate-go-client/v3/test/testsuit"
	"github.com/semi-technologies/weaviate-go-client/v3/weaviate/graphql"
	"github.com/semi-technologies/weaviate-go-client/v3/weaviate/testenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGraphQL_integration(t *testing.T) {

	t.Run("up", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			fmt.Printf(err.Error())
			t.Fail()
		}
	})

	t.Run("Get", func(t *testing.T) {
		client := testsuit.CreateTestClient()
		testsuit.CreateTestSchemaAndData(t, client)

		resultSet, gqlErr := client.GraphQL().Get().Objects().WithClassName("Pizza").WithFields("name").Do(context.Background())
		assert.Nil(t, gqlErr)

		get := resultSet.Data["Get"].(map[string]interface{})
		pizza := get["Pizza"].([]interface{})
		assert.Equal(t, 4, len(pizza))

		withNearObject := client.GraphQL().NearObjectArgBuilder().
			WithID("5b6a08ba-1d46-43aa-89cc-8b070790c6f2")
		resultSet, gqlErr = client.GraphQL().Get().Objects().
			WithClassName("Pizza").
			WithFields("name").
			WithNearObject(withNearObject).
			Do(context.Background())
		assert.Nil(t, gqlErr)
		require.Nil(t, resultSet.Errors)

		get = resultSet.Data["Get"].(map[string]interface{})
		pizza = get["Pizza"].([]interface{})
		assert.Equal(t, 4, len(pizza))

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("Explore", func(t *testing.T) {
		client := testsuit.CreateTestClient()
		testsuit.CreateTestSchemaAndData(t, client)

		fields := []graphql.ExploreFields{graphql.Certainty, graphql.Beacon, graphql.ClassName}
		concepts := []string{"pineapple slices", "ham"}
		moveTo := &graphql.MoveParameters{
			Concepts: []string{"Pizza"},
			Force:    0.3,
		}
		moveAwayFrom := &graphql.MoveParameters{
			Concepts: []string{"toast", "bread"},
			Force:    0.4,
		}

		withNearText := client.GraphQL().NearTextArgBuilder().
			WithConcepts(concepts).
			WithCertainty(0.71).
			WithMoveTo(moveTo).
			WithMoveAwayFrom(moveAwayFrom)

		resultSet, gqlErr := client.GraphQL().Explore().
			WithFields(fields).
			WithNearText(withNearText).
			Do(context.Background())

		assert.Nil(t, gqlErr)
		assert.NotNil(t, resultSet)

		withNearObject := client.GraphQL().NearObjectArgBuilder().
			WithID("5b6a08ba-1d46-43aa-89cc-8b070790c6f2")

		resultSet, gqlErr = client.GraphQL().Explore().
			WithFields(fields).
			WithNearObject(withNearObject).
			Do(context.Background())

		assert.Nil(t, gqlErr)
		assert.NotNil(t, resultSet)

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("Aggregate", func(t *testing.T) {
		client := testsuit.CreateTestClient()
		testsuit.CreateTestSchemaAndData(t, client)

		fields := "meta {count}"
		resultSet, gqlErr := client.GraphQL().Aggregate().Objects().WithFields(fields).WithClassName("Pizza").Do(context.Background())

		assert.Nil(t, gqlErr)
		assert.NotNil(t, resultSet)

		resultSet, gqlErr = client.GraphQL().
			Aggregate().
			Objects().
			WithFields(fields).
			WithWhere(`{path:["id"] operator:Equal valueString:"5b6a08ba-1d46-43aa-89cc-8b070790c6f2"}`).
			WithClassName("Pizza").
			Do(context.Background())

		assert.Nil(t, gqlErr)
		assert.NotNil(t, resultSet)

		resultSet, gqlErr = client.GraphQL().
			Aggregate().
			Objects().
			WithFields(fields).
			WithGroupBy("name").
			WithClassName("Pizza").
			Do(context.Background())

		assert.Nil(t, gqlErr)
		assert.NotNil(t, resultSet)

		resultSet, gqlErr = client.GraphQL().
			Aggregate().
			Objects().
			WithFields(fields).
			WithWhere(`{path:["id"] operator:Equal valueString:"5b6a08ba-1d46-43aa-89cc-8b070790c6f2"}`).
			WithGroupBy("name").
			WithClassName("Pizza").
			Do(context.Background())

		assert.Nil(t, gqlErr)
		assert.NotNil(t, resultSet)

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("Get with group filter", func(t *testing.T) {
		client := testsuit.CreateTestClient()
		testsuit.CreateTestSchemaAndData(t, client)

		resultSet, gqlErr := client.GraphQL().
			Get().
			Objects().
			WithClassName("Pizza").
			WithFields("name").
			WithGroup("{type:merge force:1.0}").
			WithLimit(7).
			Do(context.Background())
		assert.Nil(t, gqlErr)

		get := resultSet.Data["Get"].(map[string]interface{})
		require.Equal(t, 1, len(get))

		pizza := get["Pizza"].([]interface{})
		assert.Equal(t, 1, len(pizza))

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			fmt.Printf(err.Error())
			t.Fail()
		}
	})
}
