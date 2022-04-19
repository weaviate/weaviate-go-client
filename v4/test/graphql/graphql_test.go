package graphql

import (
	"context"
	"testing"

	"github.com/semi-technologies/weaviate-go-client/v4/test/testsuit"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/graphql"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/testenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGraphQL_integration(t *testing.T) {

	t.Run("up", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to setup weaviate: %s", err)
		}
	})

	t.Run("Get", func(t *testing.T) {
		client := testsuit.CreateTestClient()
		testsuit.CreateTestSchemaAndData(t, client)

		fields := []graphql.Field{
			{
				Name: "name",
			},
		}

		resultSet, gqlErr := client.GraphQL().Get().Objects().WithClassName("Pizza").WithFields(fields).Do(context.Background())
		assert.Nil(t, gqlErr)

		get := resultSet.Data["Get"].(map[string]interface{})
		pizza := get["Pizza"].([]interface{})
		assert.Equal(t, 4, len(pizza))

		withNearObject := client.GraphQL().NearObjectArgBuilder().
			WithID("5b6a08ba-1d46-43aa-89cc-8b070790c6f2")
		resultSet, gqlErr = client.GraphQL().Get().Objects().
			WithClassName("Pizza").
			WithFields(fields).
			WithNearObject(withNearObject).
			Do(context.Background())
		assert.Nil(t, gqlErr)
		require.Nil(t, resultSet.Errors)

		get = resultSet.Data["Get"].(map[string]interface{})
		pizza = get["Pizza"].([]interface{})
		assert.Equal(t, 4, len(pizza))

		where := client.GraphQL().WhereArgBuilder().
			WithPath([]string{"name"}).
			WithOperator(graphql.Equal).
			WithValueString("Frutti di Mare")

		fields = []graphql.Field{{Name: "name"}}

		resultSet, gqlErr = client.GraphQL().Get().Objects().WithClassName("Pizza").WithWhere(where).WithFields(fields).Do(context.Background())
		assert.Nil(t, gqlErr)

		get = resultSet.Data["Get"].(map[string]interface{})
		pizza = get["Pizza"].([]interface{})
		assert.Equal(t, 1, len(pizza))

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

		fields := []graphql.Field{
			{
				Name: "meta",
				Fields: []graphql.Field{
					{
						Name: "count",
					},
				},
			},
		}

		t.Run("no filters", func(t *testing.T) {
			resultSet, gqlErr := client.GraphQL().
				Aggregate().
				WithFields(fields).
				WithClassName("Pizza").
				Do(context.Background())

			assert.Nil(t, gqlErr)
			assert.NotNil(t, resultSet)
		})

		t.Run("with where filter", func(t *testing.T) {
			where := client.GraphQL().WhereArgBuilder().
				WithPath([]string{"id"}).
				WithOperator(graphql.Equal).
				WithValueString("5b6a08ba-1d46-43aa-89cc-8b070790c6f2")

			resultSet, gqlErr := client.GraphQL().
				Aggregate().
				WithFields(fields).
				WithWhere(where).
				WithClassName("Pizza").
				Do(context.Background())

			assert.Nil(t, gqlErr)
			assert.NotNil(t, resultSet)
		})

		t.Run("with groupby", func(t *testing.T) {
			resultSet, gqlErr := client.GraphQL().
				Aggregate().
				WithFields(fields).
				WithGroupBy("name").
				WithClassName("Pizza").
				Do(context.Background())

			assert.Nil(t, gqlErr)
			assert.NotNil(t, resultSet)
		})

		t.Run("with where filter and groupby", func(t *testing.T) {
			where := client.GraphQL().WhereArgBuilder().
				WithPath([]string{"id"}).
				WithOperator(graphql.Equal).
				WithValueString("5b6a08ba-1d46-43aa-89cc-8b070790c6f2")

			resultSet, gqlErr := client.GraphQL().
				Aggregate().
				WithFields(fields).
				WithWhere(where).
				WithGroupBy("name").
				WithClassName("Pizza").
				Do(context.Background())

			assert.Nil(t, gqlErr)
			assert.NotNil(t, resultSet)
		})

		t.Run("with nearVector filter", func(t *testing.T) {
			pizza := GetOnePizza(t, client)

			nearVector := &graphql.NearVectorArgumentBuilder{}
			nearVector.WithCertainty(0.85).
				WithVector(pizza.Additional.Vector)

			resultSet, gqlErr := client.GraphQL().
				Aggregate().
				WithFields(fields).
				WithNearVector(nearVector).
				WithClassName("Pizza").
				Do(context.Background())

			assert.Nil(t, gqlErr)
			assert.NotNil(t, resultSet)
		})

		t.Run("with nearObject filter", func(t *testing.T) {
			pizza := GetOnePizza(t, client)

			nearObject := &graphql.NearObjectArgumentBuilder{}
			nearObject.WithCertainty(0.85).
				WithID(pizza.Additional.ID)

			resultSet, gqlErr := client.GraphQL().
				Aggregate().
				WithFields(fields).
				WithNearObject(nearObject).
				WithClassName("Pizza").
				Do(context.Background())

			assert.Nil(t, gqlErr)
			assert.NotNil(t, resultSet)
		})

		t.Run("with nearText filter", func(t *testing.T) {
			concepts := []string{"pineapple slices", "ham"}
			moveTo := &graphql.MoveParameters{
				Concepts: []string{"Pizza"},
				Force:    0.3,
			}
			moveAwayFrom := &graphql.MoveParameters{
				Concepts: []string{"toast", "bread"},
				Force:    0.4,
			}

			nearText := &graphql.NearTextArgumentBuilder{}
			nearText.WithCertainty(0.85).
				WithConcepts(concepts).
				WithMoveTo(moveTo).
				WithMoveAwayFrom(moveAwayFrom)

			resultSet, gqlErr := client.GraphQL().
				Aggregate().
				WithFields(fields).
				WithNearText(nearText).
				WithClassName("Pizza").
				Do(context.Background())

			assert.Nil(t, gqlErr)
			assert.NotNil(t, resultSet)
		})

		t.Run("with nearVector, where, groupby", func(t *testing.T) {
			where := client.GraphQL().WhereArgBuilder().
				WithPath([]string{"id"}).
				WithOperator(graphql.Equal).
				WithValueString("5b6a08ba-1d46-43aa-89cc-8b070790c6f2")

			pizza := GetOnePizza(t, client)
			nearVector := &graphql.NearVectorArgumentBuilder{}
			nearVector.WithCertainty(0.85).
				WithVector(pizza.Additional.Vector)

			resultSet, gqlErr := client.GraphQL().
				Aggregate().
				WithFields(fields).
				WithWhere(where).
				WithGroupBy("name").
				WithNearVector(nearVector).
				WithClassName("Pizza").
				Do(context.Background())

			assert.Nil(t, gqlErr)
			assert.NotNil(t, resultSet)
		})

		t.Run("with nearObject, where, groupby", func(t *testing.T) {
			where := client.GraphQL().WhereArgBuilder().
				WithPath([]string{"id"}).
				WithOperator(graphql.Equal).
				WithValueString("5b6a08ba-1d46-43aa-89cc-8b070790c6f2")

			pizza := GetOnePizza(t, client)
			nearObject := &graphql.NearObjectArgumentBuilder{}
			nearObject.WithCertainty(0.85).
				WithID(pizza.Additional.ID)

			resultSet, gqlErr := client.GraphQL().
				Aggregate().
				WithFields(fields).
				WithWhere(where).
				WithGroupBy("name").
				WithNearObject(nearObject).
				WithClassName("Pizza").
				Do(context.Background())

			assert.Nil(t, gqlErr)
			assert.NotNil(t, resultSet)
		})

		t.Run("with nearText, where, groupby", func(t *testing.T) {
			where := client.GraphQL().WhereArgBuilder().
				WithPath([]string{"id"}).
				WithOperator(graphql.Equal).
				WithValueString("5b6a08ba-1d46-43aa-89cc-8b070790c6f2")

			concepts := []string{"pineapple slices", "ham"}
			moveTo := &graphql.MoveParameters{
				Concepts: []string{"Pizza"},
				Force:    0.3,
			}
			moveAwayFrom := &graphql.MoveParameters{
				Concepts: []string{"toast", "bread"},
				Force:    0.4,
			}

			nearText := &graphql.NearTextArgumentBuilder{}
			nearText.WithCertainty(0.85).
				WithConcepts(concepts).
				WithMoveTo(moveTo).
				WithMoveAwayFrom(moveAwayFrom)

			resultSet, gqlErr := client.GraphQL().
				Aggregate().
				WithFields(fields).
				WithWhere(where).
				WithGroupBy("name").
				WithNearText(nearText).
				WithClassName("Pizza").
				Do(context.Background())

			assert.Nil(t, gqlErr)
			assert.NotNil(t, resultSet)
		})

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("Get with group filter", func(t *testing.T) {
		client := testsuit.CreateTestClient()
		testsuit.CreateTestSchemaAndData(t, client)

		fields := []graphql.Field{{Name: "name"}}
		group := client.GraphQL().GroupArgBuilder().WithType(graphql.Merge).WithForce(1.0)

		resultSet, gqlErr := client.GraphQL().
			Get().
			Objects().
			WithClassName("Pizza").
			WithFields(fields).
			WithGroup(group).
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
			t.Fatalf("failed to tear down weaviate: %s", err)
		}
	})
}
