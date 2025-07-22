package graphql

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/weaviate/weaviate-go-client/v5/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/data/replication"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/graphql"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/testenv"
	"github.com/weaviate/weaviate/entities/models"
)

func TestGraphQL_integration(t *testing.T) {
	t.Run("up", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to setup weaviate: %s", err)
		}
	})

	t.Run("Get", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		testsuit.CreateTestSchemaAndData(t, client)
		defer testsuit.CleanUpWeaviate(t, client)

		name := graphql.Field{Name: "name"}

		t.Run("get raw", func(t *testing.T) {
			resultSet, gqlErr := client.GraphQL().Raw().WithQuery("{Get {Pizza {name}}}").Do(context.Background())
			assert.NoError(t, gqlErr)

			get := resultSet.Data["Get"].(map[string]interface{})
			pizza := get["Pizza"].([]interface{})
			assert.Equal(t, 4, len(pizza))
		})

		t.Run("get all", func(t *testing.T) {
			resultSet, gqlErr := client.GraphQL().Get().WithClassName("Pizza").WithFields(name).Do(context.Background())
			assert.NoError(t, gqlErr)

			get := resultSet.Data["Get"].(map[string]interface{})
			pizza := get["Pizza"].([]interface{})
			assert.Equal(t, 4, len(pizza))
		})

		t.Run("get all with cursor", func(t *testing.T) {
			afterID := "00000000-0000-0000-0000-000000000000"
			resultSet, gqlErr := client.GraphQL().Get().WithClassName("Pizza").
				WithFields(name).WithAfter(afterID).WithLimit(10).Do(context.Background())
			assert.NoError(t, gqlErr)

			get := resultSet.Data["Get"].(map[string]interface{})
			pizza := get["Pizza"].([]interface{})
			assert.Equal(t, 3, len(pizza)) // 3 instead of 4
		})

		t.Run("by near object", func(t *testing.T) {
			withNearObject := client.GraphQL().NearObjectArgBuilder().
				WithID("5b6a08ba-1d46-43aa-89cc-8b070790c6f2")
			resultSet, gqlErr := client.GraphQL().Get().
				WithClassName("Pizza").
				WithFields(name).
				WithNearObject(withNearObject).
				Do(context.Background())
			assert.NoError(t, gqlErr)
			require.Empty(t, prettyErrors(resultSet.Errors))

			get := resultSet.Data["Get"].(map[string]interface{})
			pizza := get["Pizza"].([]interface{})
			assert.Equal(t, 4, len(pizza))
		})

		t.Run("by near text and movers", func(t *testing.T) {
			t.Run("with certainty", func(t *testing.T) {
				concepts := []string{"pineapple slices", "ham"}
				moveTo := &graphql.MoveParameters{
					Force: 0.3,
					Objects: []graphql.MoverObject{
						{ID: "00000000-0000-0000-0000-000000000000"},
						{Beacon: "weaviate://localhost/00000000-0000-0000-0000-000000000000"},
					},
				}
				moveAwayFrom := &graphql.MoveParameters{
					Force: 0.3,
					Objects: []graphql.MoverObject{
						{},
						{ID: "5b6a08ba-1d46-43aa-89cc-8b070790c6f2"},
						{
							ID:     "5b6a08ba-1d46-43aa-89cc-8b070790c6f2",
							Beacon: "weaviate://localhost/5b6a08ba-1d46-43aa-89cc-8b070790c6f2",
						},
					},
				}

				withNearText := client.GraphQL().NearTextArgBuilder().
					WithConcepts(concepts).
					WithCertainty(0.71).
					WithMoveTo(moveTo).
					WithMoveAwayFrom(moveAwayFrom)

				resultSet, gqlErr := client.GraphQL().Get().WithClassName("Pizza").
					WithFields(name).
					WithNearText(withNearText).
					Do(context.Background())

				assert.NoError(t, gqlErr)
				assert.NotNil(t, resultSet)
				require.Empty(t, prettyErrors(resultSet.Errors))
			})

			t.Run("with distance", func(t *testing.T) {
				concepts := []string{"pineapple slices", "ham"}
				moveTo := &graphql.MoveParameters{
					Force: 0.3,
					Objects: []graphql.MoverObject{
						{ID: "00000000-0000-0000-0000-000000000000"},
						{Beacon: "weaviate://localhost/00000000-0000-0000-0000-000000000000"},
					},
				}
				moveAwayFrom := &graphql.MoveParameters{
					Force: 0.3,
					Objects: []graphql.MoverObject{
						{},
						{ID: "5b6a08ba-1d46-43aa-89cc-8b070790c6f2"},
						{
							ID:     "5b6a08ba-1d46-43aa-89cc-8b070790c6f2",
							Beacon: "weaviate://localhost/5b6a08ba-1d46-43aa-89cc-8b070790c6f2",
						},
					},
				}

				withNearText := client.GraphQL().NearTextArgBuilder().
					WithConcepts(concepts).
					WithDistance(0.29).
					WithMoveTo(moveTo).
					WithMoveAwayFrom(moveAwayFrom)

				resultSet, gqlErr := client.GraphQL().Get().WithClassName("Pizza").
					WithFields(name).
					WithNearText(withNearText).
					Do(context.Background())

				assert.NoError(t, gqlErr)
				assert.NotNil(t, resultSet)
				require.Empty(t, prettyErrors(resultSet.Errors))
			})
		})

		t.Run("with where filter (string)", func(t *testing.T) {
			where := filters.Where().
				WithPath([]string{"name"}).
				WithOperator(filters.Equal).
				WithValueString("Frutti di Mare")

			name = graphql.Field{Name: "name"}

			resultSet, gqlErr := client.GraphQL().Get().WithClassName("Pizza").WithWhere(where).WithFields(name).Do(context.Background())
			assert.NoError(t, gqlErr)

			get := resultSet.Data["Get"].(map[string]interface{})
			pizza := get["Pizza"].([]interface{})
			assert.Equal(t, 1, len(pizza))
		})

		t.Run("with where filter (date)", func(t *testing.T) {
			// two pizzas have best_before dates in the test set, one of them expires
			// May 3rd, the other May 5th, so the filter below should match exactly
			// one.
			targetDate, err := time.Parse(time.RFC3339, "2022-05-04T12:00:00+02:00")
			require.NoError(t, err)

			where := filters.Where().
				WithPath([]string{"best_before"}).
				WithOperator(filters.LessThan).
				WithValueDate(targetDate)

			name = graphql.Field{Name: "name"}

			resultSet, gqlErr := client.GraphQL().Get().WithClassName("Pizza").WithWhere(where).WithFields(name).Do(context.Background())
			assert.NoError(t, gqlErr)

			get := resultSet.Data["Get"].(map[string]interface{})
			pizza := get["Pizza"].([]interface{})
			assert.Equal(t, 1, len(pizza))
		})

		t.Run("with sorting", func(t *testing.T) {
			var pizza []interface{}

			assertSortResult := func(resultSet *models.GraphQLResponse, className string, expectedPizzas []string) {
				get := resultSet.Data["Get"].(map[string]interface{})
				pizza = get[className].([]interface{})
				require.Equal(t, len(expectedPizzas), len(pizza))
				result := make([]string, len(pizza))
				for i := range pizza {
					p := pizza[i].(map[string]interface{})
					result[i] = p["name"].(string)
				}
				assert.Equal(t, expectedPizzas, result)
			}

			byNameAsc := graphql.Sort{Path: []string{"name"}, Order: graphql.Asc}
			resultSet, gqlErr := client.GraphQL().Get().WithClassName("Pizza").
				WithSort(byNameAsc).WithFields(name).Do(context.Background())
			assert.NoError(t, gqlErr)
			assertSortResult(resultSet, "Pizza", []string{"Doener", "Frutti di Mare", "Hawaii", "Quattro Formaggi"})

			byNameDesc := graphql.Sort{Path: []string{"name"}, Order: graphql.Desc}
			resultSet, gqlErr = client.GraphQL().Get().WithClassName("Pizza").
				WithSort(byNameDesc).WithFields(name).Do(context.Background())
			assert.NoError(t, gqlErr)
			assertSortResult(resultSet, "Pizza", []string{"Quattro Formaggi", "Hawaii", "Frutti di Mare", "Doener"})

			byPriceAsc := graphql.Sort{Path: []string{"price"}, Order: graphql.Asc}
			resultSet, gqlErr = client.GraphQL().Get().WithClassName("Soup").
				WithSort(byPriceAsc).WithFields(name).Do(context.Background())
			assert.NoError(t, gqlErr)
			assertSortResult(resultSet, "Soup", []string{"ChickenSoup", "Beautiful"})

			resultSet, gqlErr = client.GraphQL().Get().WithClassName("Pizza").
				WithSort(byPriceAsc, byNameDesc).WithFields(name).Do(context.Background())
			assert.NoError(t, gqlErr)
			assertSortResult(resultSet, "Pizza", []string{"Quattro Formaggi", "Frutti di Mare", "Hawaii", "Doener"})
		})

		t.Run("generative OpenAI", func(t *testing.T) {
			t.Skip("skipping all generative OpenAI tests due to OpenAI API being unstable")
			if os.Getenv("OPENAI_APIKEY") == "" {
				t.Skip("No open-ai api key added")
			}

			t.Run("with generative search single result", func(t *testing.T) {
				gs := graphql.NewGenerativeSearch().SingleResult("Describe this pizza : {name}")

				resultSet, gqlErr := client.GraphQL().Get().
					WithClassName("Pizza").
					WithFields(name).
					WithGenerativeSearch(gs).
					Do(context.Background())
				assert.NoError(t, gqlErr)

				get := resultSet.Data["Get"].(map[string]interface{})
				pizzas := get["Pizza"].([]interface{})
				assert.Equal(t, 4, len(pizzas))
				for _, pizza := range pizzas {
					_additional := pizza.(map[string]interface{})["_additional"]
					assert.NotNil(t, _additional)

					generate := _additional.(map[string]interface{})["generate"].(map[string]interface{})
					generateErr := generate["error"]
					assert.Nil(t, generateErr)
					singleResult := generate["singleResult"].(string)
					assert.NotEmpty(t, singleResult)
				}
			})

			t.Run("with generative search grouped result", func(t *testing.T) {
				gs := graphql.NewGenerativeSearch().GroupedResult("Describe these pizzas")

				resultSet, gqlErr := client.GraphQL().Get().
					WithClassName("Pizza").
					WithFields(name).
					WithGenerativeSearch(gs).
					Do(context.Background())
				assert.NoError(t, gqlErr)

				get := resultSet.Data["Get"].(map[string]interface{})
				pizza := get["Pizza"].([]interface{})
				assert.Equal(t, 4, len(pizza))

				_additional := pizza[0].(map[string]interface{})["_additional"]
				assert.NotNil(t, _additional)

				generate := _additional.(map[string]interface{})["generate"].(map[string]interface{})
				generateErr := generate["error"]
				assert.Nil(t, generateErr)
				groupedResult := generate["groupedResult"].(string)
				assert.NotEmpty(t, groupedResult)
			})

			t.Run("with generative search single result and grouped result", func(t *testing.T) {
				gs := graphql.NewGenerativeSearch().
					SingleResult("Describe this pizza : {name}").
					GroupedResult("Describe these pizzas")

				resultSet, gqlErr := client.GraphQL().Get().
					WithClassName("Pizza").
					WithFields(name).
					WithGenerativeSearch(gs).
					Do(context.Background())
				assert.NoError(t, gqlErr)

				get := resultSet.Data["Get"].(map[string]interface{})
				pizzas := get["Pizza"].([]interface{})
				assert.Equal(t, 4, len(pizzas))

				for _, pizza := range pizzas {
					_additional := pizza.(map[string]interface{})["_additional"]
					assert.NotNil(t, _additional)

					generate := _additional.(map[string]interface{})["generate"].(map[string]interface{})
					generateErr := generate["error"]
					assert.Nil(t, generateErr)
					singleResult := generate["singleResult"].(string)
					assert.NotEmpty(t, singleResult)
				}

				_additional := pizzas[0].(map[string]interface{})["_additional"]
				assert.NotNil(t, _additional)

				generate := _additional.(map[string]interface{})["generate"].(map[string]interface{})
				generateErr := generate["error"]
				assert.Nil(t, generateErr)
				groupedResult := generate["groupedResult"].(string)
				assert.NotEmpty(t, groupedResult)
			})

			t.Run("with generative search grouped result with properties", func(t *testing.T) {
				gs := graphql.NewGenerativeSearch().GroupedResult("Describe these pizzas", "title", "description")

				resultSet, gqlErr := client.GraphQL().Get().
					WithClassName("Pizza").
					WithFields(name).
					WithGenerativeSearch(gs).
					Do(context.Background())
				assert.NoError(t, gqlErr)

				get := resultSet.Data["Get"].(map[string]interface{})
				pizza := get["Pizza"].([]interface{})
				assert.Equal(t, 4, len(pizza))

				_additional := pizza[0].(map[string]interface{})["_additional"]
				assert.NotNil(t, _additional)

				generate := _additional.(map[string]interface{})["generate"].(map[string]interface{})
				generateErr := generate["error"]
				assert.Nil(t, generateErr)
				groupedResult := generate["groupedResult"].(string)
				assert.NotEmpty(t, groupedResult)
			})
		})
	})

	t.Run("Explore", func(t *testing.T) {
		t.Run("with certainty", func(t *testing.T) {
			client := testsuit.CreateTestClient(false)
			testsuit.CreateTestSchemaAndData(t, client)
			defer testsuit.CleanUpWeaviate(t, client)

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
				WithFields(graphql.Certainty, graphql.Beacon, graphql.ClassName).
				WithNearText(withNearText).
				Do(context.Background())

			assert.NoError(t, gqlErr)
			assert.NotNil(t, resultSet)

			withNearObject := client.GraphQL().NearObjectArgBuilder().
				WithID("5b6a08ba-1d46-43aa-89cc-8b070790c6f2")

			resultSet, gqlErr = client.GraphQL().Explore().
				WithFields(graphql.Certainty, graphql.Beacon, graphql.ClassName).
				WithNearObject(withNearObject).
				Do(context.Background())

			assert.NoError(t, gqlErr)
			assert.NotNil(t, resultSet)
		})

		t.Run("with distance", func(t *testing.T) {
			client := testsuit.CreateTestClient(false)
			testsuit.CreateTestSchemaAndData(t, client)
			defer testsuit.CleanUpWeaviate(t, client)

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
				WithDistance(0.29).
				WithMoveTo(moveTo).
				WithMoveAwayFrom(moveAwayFrom)

			resultSet, gqlErr := client.GraphQL().Explore().
				WithFields(graphql.Certainty, graphql.Beacon, graphql.ClassName).
				WithNearText(withNearText).
				Do(context.Background())

			assert.NoError(t, gqlErr)
			assert.NotNil(t, resultSet)

			withNearObject := client.GraphQL().NearObjectArgBuilder().
				WithID("5b6a08ba-1d46-43aa-89cc-8b070790c6f2")

			resultSet, gqlErr = client.GraphQL().Explore().
				WithFields(graphql.Certainty, graphql.Beacon, graphql.ClassName).
				WithNearObject(withNearObject).
				Do(context.Background())

			assert.NoError(t, gqlErr)
			assert.NotNil(t, resultSet)
		})
	})

	t.Run("Aggregate", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		testsuit.CreateTestSchemaAndData(t, client)
		defer testsuit.CleanUpWeaviate(t, client)

		meta := graphql.Field{
			Name:   "meta",
			Fields: []graphql.Field{{Name: "count"}},
		}

		t.Run("no filters", func(t *testing.T) {
			resultSet, gqlErr := client.GraphQL().
				Aggregate().
				WithFields(meta).
				WithClassName("Pizza").
				Do(context.Background())

			assert.NoError(t, gqlErr)
			assert.NotNil(t, resultSet)
		})

		t.Run("with where filter", func(t *testing.T) {
			where := filters.Where().
				WithPath([]string{"id"}).
				WithOperator(filters.Equal).
				WithValueString("5b6a08ba-1d46-43aa-89cc-8b070790c6f2")

			resultSet, gqlErr := client.GraphQL().
				Aggregate().
				WithFields(meta).
				WithWhere(where).
				WithClassName("Pizza").
				Do(context.Background())

			assert.NoError(t, gqlErr)
			assert.NotNil(t, resultSet)
		})

		t.Run("with groupby", func(t *testing.T) {
			resultSet, gqlErr := client.GraphQL().
				Aggregate().
				WithFields(meta).
				WithGroupBy("name").
				WithClassName("Pizza").
				Do(context.Background())

			assert.NoError(t, gqlErr)
			assert.NotNil(t, resultSet)
		})

		t.Run("with where filter and groupby", func(t *testing.T) {
			where := filters.Where().
				WithPath([]string{"id"}).
				WithOperator(filters.Equal).
				WithValueString("5b6a08ba-1d46-43aa-89cc-8b070790c6f2")

			resultSet, gqlErr := client.GraphQL().
				Aggregate().
				WithFields(meta).
				WithWhere(where).
				WithGroupBy("name").
				WithClassName("Pizza").
				Do(context.Background())

			assert.NoError(t, gqlErr)
			assert.NotNil(t, resultSet)
		})

		t.Run("with nearVector filter", func(t *testing.T) {
			t.Run("with certainty", func(t *testing.T) {
				pizza := GetOnePizza(t, client)

				nearVector := &graphql.NearVectorArgumentBuilder{}
				nearVector.WithCertainty(0.85).
					WithVector(pizza.Additional.Vector)

				resultSet, gqlErr := client.GraphQL().
					Aggregate().
					WithFields(meta).
					WithNearVector(nearVector).
					WithClassName("Pizza").
					Do(context.Background())

				assert.NoError(t, gqlErr)
				assert.NotNil(t, resultSet)
			})

			t.Run("with distance", func(t *testing.T) {
				pizza := GetOnePizza(t, client)

				nearVector := &graphql.NearVectorArgumentBuilder{}
				nearVector.WithDistance(0.15).
					WithVector(pizza.Additional.Vector)

				resultSet, gqlErr := client.GraphQL().
					Aggregate().
					WithFields(meta).
					WithNearVector(nearVector).
					WithClassName("Pizza").
					Do(context.Background())

				assert.NoError(t, gqlErr)
				assert.NotNil(t, resultSet)
			})
		})

		t.Run("with nearObject filter", func(t *testing.T) {
			t.Run("with certainty", func(t *testing.T) {
				pizza := GetOnePizza(t, client)

				nearObject := &graphql.NearObjectArgumentBuilder{}
				nearObject.WithCertainty(0.85).
					WithID(pizza.Additional.ID)

				resultSet, gqlErr := client.GraphQL().
					Aggregate().
					WithFields(meta).
					WithNearObject(nearObject).
					WithClassName("Pizza").
					Do(context.Background())

				assert.NoError(t, gqlErr)
				assert.NotNil(t, resultSet)
			})

			t.Run("with distance", func(t *testing.T) {
				pizza := GetOnePizza(t, client)

				nearObject := &graphql.NearObjectArgumentBuilder{}
				nearObject.WithDistance(0.15).
					WithID(pizza.Additional.ID)

				resultSet, gqlErr := client.GraphQL().
					Aggregate().
					WithFields(meta).
					WithNearObject(nearObject).
					WithClassName("Pizza").
					Do(context.Background())

				assert.NoError(t, gqlErr)
				assert.NotNil(t, resultSet)
			})
		})

		t.Run("with nearText filter", func(t *testing.T) {
			t.Run("with certainty", func(t *testing.T) {
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
					WithFields(meta).
					WithNearText(nearText).
					WithClassName("Pizza").
					Do(context.Background())

				assert.NoError(t, gqlErr)
				assert.NotNil(t, resultSet)
			})

			t.Run("with distance", func(t *testing.T) {
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
				nearText.WithDistance(0.15).
					WithConcepts(concepts).
					WithMoveTo(moveTo).
					WithMoveAwayFrom(moveAwayFrom)

				resultSet, gqlErr := client.GraphQL().
					Aggregate().
					WithFields(meta).
					WithNearText(nearText).
					WithClassName("Pizza").
					Do(context.Background())

				assert.NoError(t, gqlErr)
				assert.NotNil(t, resultSet)
			})
		})

		t.Run("with nearVector, where, groupby", func(t *testing.T) {
			t.Run("with certainty", func(t *testing.T) {
				where := filters.Where().
					WithPath([]string{"id"}).
					WithOperator(filters.Equal).
					WithValueString("5b6a08ba-1d46-43aa-89cc-8b070790c6f2")

				pizza := GetOnePizza(t, client)
				nearVector := &graphql.NearVectorArgumentBuilder{}
				nearVector.WithCertainty(0.85).
					WithVector(pizza.Additional.Vector)

				resultSet, gqlErr := client.GraphQL().
					Aggregate().
					WithFields(meta).
					WithWhere(where).
					WithGroupBy("name").
					WithNearVector(nearVector).
					WithClassName("Pizza").
					Do(context.Background())

				assert.NoError(t, gqlErr)
				assert.NotNil(t, resultSet)
			})

			t.Run("with distance", func(t *testing.T) {
				where := filters.Where().
					WithPath([]string{"id"}).
					WithOperator(filters.Equal).
					WithValueString("5b6a08ba-1d46-43aa-89cc-8b070790c6f2")

				pizza := GetOnePizza(t, client)
				nearVector := &graphql.NearVectorArgumentBuilder{}
				nearVector.WithDistance(0.15).
					WithVector(pizza.Additional.Vector)

				resultSet, gqlErr := client.GraphQL().
					Aggregate().
					WithFields(meta).
					WithWhere(where).
					WithGroupBy("name").
					WithNearVector(nearVector).
					WithClassName("Pizza").
					Do(context.Background())

				assert.NoError(t, gqlErr)
				assert.NotNil(t, resultSet)
			})
		})

		t.Run("with nearObject, where, groupby", func(t *testing.T) {
			t.Run("with certainty", func(t *testing.T) {
				where := filters.Where().
					WithPath([]string{"id"}).
					WithOperator(filters.Equal).
					WithValueString("5b6a08ba-1d46-43aa-89cc-8b070790c6f2")

				pizza := GetOnePizza(t, client)
				nearObject := &graphql.NearObjectArgumentBuilder{}
				nearObject.WithCertainty(0.85).
					WithID(pizza.Additional.ID)

				resultSet, gqlErr := client.GraphQL().
					Aggregate().
					WithFields(meta).
					WithWhere(where).
					WithGroupBy("name").
					WithNearObject(nearObject).
					WithClassName("Pizza").
					Do(context.Background())

				assert.NoError(t, gqlErr)
				assert.NotNil(t, resultSet)
			})

			t.Run("with distance", func(t *testing.T) {
				where := filters.Where().
					WithPath([]string{"id"}).
					WithOperator(filters.Equal).
					WithValueString("5b6a08ba-1d46-43aa-89cc-8b070790c6f2")

				pizza := GetOnePizza(t, client)
				nearObject := &graphql.NearObjectArgumentBuilder{}
				nearObject.WithDistance(0.15).
					WithID(pizza.Additional.ID)

				resultSet, gqlErr := client.GraphQL().
					Aggregate().
					WithFields(meta).
					WithWhere(where).
					WithGroupBy("name").
					WithNearObject(nearObject).
					WithClassName("Pizza").
					Do(context.Background())

				assert.NoError(t, gqlErr)
				assert.NotNil(t, resultSet)
			})
		})

		t.Run("with nearText, where, groupby", func(t *testing.T) {
			t.Run("with certainty", func(t *testing.T) {
				where := filters.Where().
					WithPath([]string{"id"}).
					WithOperator(filters.Equal).
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
					WithFields(meta).
					WithWhere(where).
					WithGroupBy("name").
					WithNearText(nearText).
					WithClassName("Pizza").
					Do(context.Background())

				assert.NoError(t, gqlErr)
				assert.NotNil(t, resultSet)
			})

			t.Run("with distance", func(t *testing.T) {
				where := filters.Where().
					WithPath([]string{"id"}).
					WithOperator(filters.Equal).
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
				nearText.WithDistance(0.15).
					WithConcepts(concepts).
					WithMoveTo(moveTo).
					WithMoveAwayFrom(moveAwayFrom)

				resultSet, gqlErr := client.GraphQL().
					Aggregate().
					WithFields(meta).
					WithWhere(where).
					WithGroupBy("name").
					WithNearText(nearText).
					WithClassName("Pizza").
					Do(context.Background())

				assert.NoError(t, gqlErr)
				assert.NotNil(t, resultSet)
			})
		})

		t.Run("with objectLimit", func(t *testing.T) {
			t.Run("with certainty", func(t *testing.T) {
				pizza := GetOnePizza(t, client)

				nearObject := &graphql.NearObjectArgumentBuilder{}
				nearObject.WithCertainty(0.5).
					WithID(pizza.Additional.ID)

				objectLimit := 1

				resultSet, gqlErr := client.GraphQL().
					Aggregate().
					WithFields(meta).
					WithNearObject(nearObject).
					WithObjectLimit(objectLimit).
					WithClassName("Pizza").
					Do(context.Background())

				assert.NoError(t, gqlErr)
				assert.NotNil(t, resultSet)
				assert.NotNil(t, resultSet.Data)

				b, err := json.Marshal(resultSet.Data)
				require.NoError(t, err)

				var resp AggregatePizzaResponse
				err = json.Unmarshal(b, &resp)
				require.NoError(t, err)

				assert.NotEmpty(t, resp.Aggregate.Pizza)
				assert.Equal(t, objectLimit, resp.Aggregate.Pizza[0].Meta.Count)
			})

			t.Run("with distance", func(t *testing.T) {
				pizza := GetOnePizza(t, client)

				nearObject := &graphql.NearObjectArgumentBuilder{}
				nearObject.WithDistance(0.5).
					WithID(pizza.Additional.ID)

				objectLimit := 1

				resultSet, gqlErr := client.GraphQL().
					Aggregate().
					WithFields(meta).
					WithNearObject(nearObject).
					WithObjectLimit(objectLimit).
					WithClassName("Pizza").
					Do(context.Background())

				assert.NoError(t, gqlErr)
				assert.NotNil(t, resultSet)
				assert.NotNil(t, resultSet.Data)

				b, err := json.Marshal(resultSet.Data)
				require.NoError(t, err)

				var resp AggregatePizzaResponse
				err = json.Unmarshal(b, &resp)
				require.NoError(t, err)

				assert.NotEmpty(t, resp.Aggregate.Pizza)
				assert.Equal(t, objectLimit, resp.Aggregate.Pizza[0].Meta.Count)
			})
		})

		t.Run("with hybrid", func(t *testing.T) {
			t.Run("with where also", func(t *testing.T) {
				where := filters.Where().
					WithPath([]string{"id"}).
					WithOperator(filters.Equal).
					WithValueString("5b6a08ba-1d46-43aa-89cc-8b070790c6f2")

				hybrid := client.GraphQL().HybridArgumentBuilder().WithQuery("toast")

				resultSet, gqlErr := client.GraphQL().
					Aggregate().
					WithFields(meta).
					WithWhere(where).
					WithHybrid(hybrid).
					WithObjectLimit(10).
					WithClassName("Pizza").
					Do(context.Background())

				assert.NoError(t, gqlErr)
				require.Empty(t, prettyErrors(resultSet.Errors))
				assert.NotNil(t, resultSet)
			})

			t.Run("with just hybrid", func(t *testing.T) {
				hybrid := client.GraphQL().HybridArgumentBuilder().
					WithQuery("toast").WithAlpha(0.6)

				resultSet, gqlErr := client.GraphQL().
					Aggregate().
					WithFields(meta).
					WithHybrid(hybrid).
					WithObjectLimit(10).
					WithClassName("Pizza").
					Do(context.Background())

				assert.NoError(t, gqlErr)
				require.Empty(t, prettyErrors(resultSet.Errors))
				assert.NotNil(t, resultSet)
			})

			t.Run("with hybrid and nearText", func(t *testing.T) {
				hybrid := client.GraphQL().HybridArgumentBuilder().WithQuery("toast")

				nearText := client.GraphQL().NearTextArgBuilder().
					WithConcepts([]string{"toast"}).WithCertainty(0.7)

				resultSet, gqlErr := client.GraphQL().
					Aggregate().
					WithFields(meta).
					WithHybrid(hybrid).
					WithNearText(nearText).
					WithClassName("Pizza").
					Do(context.Background())

				assert.NoError(t, gqlErr)
				require.Empty(t, prettyErrors(resultSet.Errors))
				assert.NotNil(t, resultSet)
			})
		})
	})

	t.Run("Get with group filter", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		testsuit.CreateTestSchemaAndData(t, client)
		defer testsuit.CleanUpWeaviate(t, client)

		name := graphql.Field{Name: "name"}
		group := client.GraphQL().GroupArgBuilder().WithType(graphql.Merge).WithForce(1.0)

		resultSet, gqlErr := client.GraphQL().
			Get().
			WithClassName("Pizza").
			WithFields(name).
			WithGroup(group).
			WithLimit(7).
			Do(context.Background())
		assert.NoError(t, gqlErr)

		get := resultSet.Data["Get"].(map[string]interface{})
		require.Equal(t, 1, len(get))

		pizza := get["Pizza"].([]interface{})
		assert.Equal(t, 1, len(pizza))
	})

	t.Run("Get with creationTimeUnix filters", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		testsuit.CreateTestSchemaAndData(t, client)
		defer testsuit.CleanUpWeaviate(t, client)

		pizza := GetOnePizza(t, client)
		expectedCreateTime := pizza.Additional.CreationTimeUnix

		additional := graphql.Field{
			Name: "_additional", Fields: []graphql.Field{
				{Name: "creationTimeUnix"},
			},
		}

		whereCreateTime := filters.Where().
			WithPath([]string{"_creationTimeUnix"}).
			WithOperator(filters.Equal).
			WithValueString(expectedCreateTime)

		result, err := client.GraphQL().Get().
			WithClassName("Pizza").
			WithFields(additional).
			WithWhere(whereCreateTime).
			Do(context.Background())

		require.NoError(t, err)
		require.Empty(t, prettyErrors(result.Errors))
		require.NotNil(t, result)
		require.NotNil(t, result.Data)

		b, err := json.Marshal(result.Data)
		require.NoError(t, err)

		var resp GetPizzaResponse
		err = json.Unmarshal(b, &resp)
		require.NoError(t, err)
		require.NotEmpty(t, resp.Get.Pizzas)

		assert.Equal(t, expectedCreateTime, resp.Get.Pizzas[0].Additional.CreationTimeUnix)
	})

	t.Run("Get with lastUpdateTimeUnix filters", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		testsuit.CreateTestSchemaAndData(t, client)
		defer testsuit.CleanUpWeaviate(t, client)

		pizza := GetOnePizza(t, client)
		expectedUpdateTime := pizza.Additional.LastUpdateTimeUnix

		additional := graphql.Field{
			Name: "_additional", Fields: []graphql.Field{
				{Name: "lastUpdateTimeUnix"},
			},
		}

		whereUpdateTime := filters.Where().
			WithPath([]string{"_lastUpdateTimeUnix"}).
			WithOperator(filters.Equal).
			WithValueString(expectedUpdateTime)

		result, err := client.GraphQL().Get().
			WithClassName("Pizza").
			WithFields(additional).
			WithWhere(whereUpdateTime).
			Do(context.Background())

		require.NoError(t, err)
		require.Empty(t, prettyErrors(result.Errors))
		require.NotNil(t, result)
		require.NotNil(t, result.Data)

		b, err := json.Marshal(result.Data)
		require.NoError(t, err)

		var resp GetPizzaResponse
		err = json.Unmarshal(b, &resp)
		require.NoError(t, err)
		require.NotEmpty(t, resp.Get.Pizzas)

		assert.Equal(t, expectedUpdateTime, resp.Get.Pizzas[0].Additional.LastUpdateTimeUnix)
	})

	t.Run("Get bm25 filter", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		testsuit.CreateTestSchemaAndData(t, client)
		defer testsuit.CleanUpWeaviate(t, client)

		description := graphql.Field{
			Name: "description",
		}

		bm25 := client.GraphQL().Bm25ArgBuilder().WithQuery("innovation").WithProperties("description")

		result, err := client.GraphQL().Get().
			WithClassName("Pizza").
			WithFields(description).
			WithBM25(bm25).
			Do(context.Background())

		require.NoError(t, err)
		require.Empty(t, prettyErrors(result.Errors))
		require.NotNil(t, result)
		require.NotNil(t, result.Data)

		b, err := json.Marshal(result.Data)
		require.NoError(t, err)

		var resp GetPizzaResponse
		err = json.Unmarshal(b, &resp)
		require.NoError(t, err)
		require.NotEmpty(t, resp.Get.Pizzas)

		assert.Equal(t, 1, len(resp.Get.Pizzas))
		assert.Equal(t, "A innovation, some say revolution, in the pizza industry.", resp.Get.Pizzas[0].Description)

		t.Run("search operator", func(t *testing.T) {
			for _, tt := range []struct {
				name     string
				operator func(graphql.BM25SearchOperatorBuilder) graphql.BM25SearchOperatorBuilder
			}{
				{name: graphql.BM25SearchOperatorAnd, operator: func(builder graphql.BM25SearchOperatorBuilder) graphql.BM25SearchOperatorBuilder {
					return *builder.WithOperator(graphql.BM25SearchOperatorAnd)
				}},
				{name: graphql.BM25SearchOperatorOr, operator: func(builder graphql.BM25SearchOperatorBuilder) graphql.BM25SearchOperatorBuilder {
					return *builder.WithOperator(graphql.BM25SearchOperatorOr).
						WithMinimumMatch(1)
				}},
			} {
				t.Run(tt.name, func(t *testing.T) {
					var builder graphql.BM25SearchOperatorBuilder

					bm25 := client.GraphQL().
						Bm25ArgBuilder().
						WithQuery("innovation").
						WithProperties("description").
						WithSearchOperator(tt.operator(builder))

					result, err := client.GraphQL().Get().
						WithClassName("Pizza").
						WithFields(description).
						WithBM25(bm25).
						Do(context.Background())

					require.NoError(t, err)
					require.Empty(t, prettyErrors(result.Errors))
					require.NotNil(t, result)
					require.NotNil(t, result.Data)

					b, err := json.Marshal(result.Data)
					require.NoError(t, err)

					var resp GetPizzaResponse
					err = json.Unmarshal(b, &resp)
					require.NoError(t, err)
					require.NotEmpty(t, resp.Get.Pizzas)

					assert.Equal(t, 1, len(resp.Get.Pizzas))
					assert.Equal(t, "A innovation, some say revolution, in the pizza industry.", resp.Get.Pizzas[0].Description)
				})
			}
		})
	})

	t.Run("Get hybrid filter", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		testsuit.CreateTestSchemaAndData(t, client)
		defer testsuit.CleanUpWeaviate(t, client)

		description := graphql.Field{
			Name: "description",
		}

		hybrid := client.GraphQL().HybridArgumentBuilder().WithQuery("some say revolution").WithAlpha(0.8)

		result, err := client.GraphQL().Get().
			WithClassName("Pizza").
			WithFields(description).
			WithHybrid(hybrid).
			Do(context.Background())

		require.NoError(t, err)
		require.Empty(t, prettyErrors(result.Errors))
		require.NotNil(t, result)
		require.NotNil(t, result.Data)

		b, err := json.Marshal(result.Data)
		require.NoError(t, err)

		var resp GetPizzaResponse
		err = json.Unmarshal(b, &resp)
		require.NoError(t, err)
		require.NotEmpty(t, resp.Get.Pizzas)

		assert.Equal(t, 4, len(resp.Get.Pizzas))

		t.Run("search operator", func(t *testing.T) {
			for _, tt := range []struct {
				name     string
				operator func(graphql.BM25SearchOperatorBuilder) graphql.BM25SearchOperatorBuilder
			}{
				{name: graphql.BM25SearchOperatorAnd, operator: func(builder graphql.BM25SearchOperatorBuilder) graphql.BM25SearchOperatorBuilder {
					return *builder.WithOperator(graphql.BM25SearchOperatorAnd)
				}},
				{name: graphql.BM25SearchOperatorOr, operator: func(builder graphql.BM25SearchOperatorBuilder) graphql.BM25SearchOperatorBuilder {
					return *builder.WithOperator(graphql.BM25SearchOperatorOr).
						WithMinimumMatch(1)
				}},
			} {
				t.Run(tt.name, func(t *testing.T) {
					var builder graphql.BM25SearchOperatorBuilder

					hybrid := client.GraphQL().
						HybridArgumentBuilder().
						WithQuery("some say revolution").
						WithAlpha(0.8).
						WithBM25SearchOperator(tt.operator(builder))

					result, err := client.GraphQL().Get().
						WithClassName("Pizza").
						WithFields(description).
						WithHybrid(hybrid).
						Do(context.Background())

					require.NoError(t, err)
					require.Empty(t, prettyErrors(result.Errors))
					require.NotNil(t, result)
					require.NotNil(t, result.Data)

					b, err := json.Marshal(result.Data)
					require.NoError(t, err)

					var resp GetPizzaResponse
					err = json.Unmarshal(b, &resp)
					require.NoError(t, err)
					require.NotEmpty(t, resp.Get.Pizzas)

					assert.Equal(t, 4, len(resp.Get.Pizzas))
				})
			}
		})
	})

	tests := []struct {
		name        string
		properties  []string
		num_results int
	}{
		{name: "Get hybrid Properties nil", properties: nil, num_results: 1},
		{name: "Get hybrid Properties empty", properties: []string{}, num_results: 1},
		{name: "Get hybrid Properties name", properties: []string{"name"}, num_results: 0},
		{name: "Get hybrid Properties name,description", properties: []string{"name", "description"}, num_results: 1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := testsuit.CreateTestClient(false)
			testsuit.CreateTestSchemaAndData(t, client)
			defer testsuit.CleanUpWeaviate(t, client)

			hybrid := client.GraphQL().HybridArgumentBuilder().WithQuery("mussels").WithAlpha(0.0).WithProperties(tc.properties)

			fields := []graphql.Field{
				{Name: "name"},
				{Name: "description"},
				{Name: "best_before"},
			}

			result, err := client.GraphQL().Get().
				WithClassName("Pizza").
				WithFields(fields...).
				WithHybrid(hybrid).
				Do(context.Background())

			require.NoError(t, err)
			require.Empty(t, prettyErrors(result.Errors))
			require.NotNil(t, result)
			require.NotNil(t, result.Data)

			b, err := json.Marshal(result.Data)
			require.NoError(t, err)

			var resp GetPizzaResponse
			err = json.Unmarshal(b, &resp)
			require.NoError(t, err)

			assert.Equal(t, tc.num_results, len(resp.Get.Pizzas))
		})
	}

	t.Run("MultiClass Get", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		testsuit.CreateTestSchemaAndData(t, client)
		defer testsuit.CleanUpWeaviate(t, client)

		name := graphql.Field{Name: "name"}
		description := graphql.Field{Name: "description"}

		// what what
		t.Run("multiclass get raw", func(t *testing.T) {
			resultSet, gqlErr := client.GraphQL().Raw().WithQuery("{Get {Pizza {name} Risotto{description}}}").Do(context.Background())
			assert.NoError(t, gqlErr)

			get := resultSet.Data["Get"].(map[string]interface{})
			pizza := get["Pizza"].([]interface{})
			assert.Equal(t, 4, len(pizza))
			risotto := get["Risotto"].([]interface{})
			assert.Equal(t, 3, len(risotto))
		})

		t.Run("multiclass get all", func(t *testing.T) {
			resultSet, gqlErr := client.GraphQL().MultiClassGet().
				AddQueryClass(graphql.NewQueryClassBuilder("Pizza").
					WithFields(name)).
				AddQueryClass(graphql.NewQueryClassBuilder("Risotto").
					WithFields(description)).
				Do(context.Background())
			assert.NoError(t, gqlErr)

			get := resultSet.Data["Get"].(map[string]interface{})
			pizza := get["Pizza"].([]interface{})
			assert.Equal(t, 4, len(pizza))
			risotto := get["Risotto"].([]interface{})
			assert.Equal(t, 3, len(risotto))
		})

		t.Run("multiclass by near object", func(t *testing.T) {
			pizzaWithNearObject := client.GraphQL().NearObjectArgBuilder().
				WithID("5b6a08ba-1d46-43aa-89cc-8b070790c6f2")
			risottoWithNearObject := client.GraphQL().NearObjectArgBuilder().
				WithID("696bf381-7f98-40a4-bcad-841780e00e0e")

			resultSet, gqlErr := client.GraphQL().MultiClassGet().
				AddQueryClass(graphql.NewQueryClassBuilder("Pizza").
					WithFields(name).
					WithNearObject(pizzaWithNearObject)).
				AddQueryClass(graphql.NewQueryClassBuilder("Risotto").
					WithFields(description).
					WithNearObject(risottoWithNearObject)).
				Do(context.Background())
			assert.NoError(t, gqlErr)
			require.Empty(t, prettyErrors(resultSet.Errors))

			get := resultSet.Data["Get"].(map[string]interface{})
			pizza := get["Pizza"].([]interface{})
			assert.Equal(t, 4, len(pizza))
			risotto := get["Risotto"].([]interface{})
			assert.Equal(t, 3, len(risotto))
		})

		t.Run("multiclass by near text and movers", func(t *testing.T) {
			t.Run("with certainty", func(t *testing.T) {
				concepts := []string{"pineapple slices", "ham"}
				moveTo := &graphql.MoveParameters{
					Force: 0.3,
					Objects: []graphql.MoverObject{
						{ID: "00000000-0000-0000-0000-000000000000"},
						{Beacon: "weaviate://localhost/00000000-0000-0000-0000-000000000000"},
					},
				}
				moveAwayFrom := &graphql.MoveParameters{
					Force: 0.3,
					Objects: []graphql.MoverObject{
						{},
						{ID: "5b6a08ba-1d46-43aa-89cc-8b070790c6f2"},
						{
							ID:     "5b6a08ba-1d46-43aa-89cc-8b070790c6f2",
							Beacon: "weaviate://localhost/5b6a08ba-1d46-43aa-89cc-8b070790c6f2",
						},
					},
				}

				withNearText := client.GraphQL().NearTextArgBuilder().
					WithConcepts(concepts).
					WithCertainty(0.71).
					WithMoveTo(moveTo).
					WithMoveAwayFrom(moveAwayFrom)

				resultSet, gqlErr := client.GraphQL().MultiClassGet().
					AddQueryClass(graphql.NewQueryClassBuilder("Pizza").
						WithFields(name).
						WithNearText(withNearText)).
					Do(context.Background())

				assert.NoError(t, gqlErr)
				assert.NotNil(t, resultSet)
				require.Empty(t, prettyErrors(resultSet.Errors))
			})

			t.Run("with distance", func(t *testing.T) {
				concepts := []string{"pineapple slices", "ham"}
				moveTo := &graphql.MoveParameters{
					Force: 0.3,
					Objects: []graphql.MoverObject{
						{ID: "00000000-0000-0000-0000-000000000000"},
						{Beacon: "weaviate://localhost/00000000-0000-0000-0000-000000000000"},
					},
				}
				moveAwayFrom := &graphql.MoveParameters{
					Force: 0.3,
					Objects: []graphql.MoverObject{
						{},
						{ID: "5b6a08ba-1d46-43aa-89cc-8b070790c6f2"},
						{
							ID:     "5b6a08ba-1d46-43aa-89cc-8b070790c6f2",
							Beacon: "weaviate://localhost/5b6a08ba-1d46-43aa-89cc-8b070790c6f2",
						},
					},
				}

				withNearText := client.GraphQL().NearTextArgBuilder().
					WithConcepts(concepts).
					WithDistance(0.29).
					WithMoveTo(moveTo).
					WithMoveAwayFrom(moveAwayFrom)

				resultSet, gqlErr := client.GraphQL().MultiClassGet().
					AddQueryClass(graphql.NewQueryClassBuilder("Pizza").
						WithFields(name).
						WithNearText(withNearText)).
					Do(context.Background())

				assert.NoError(t, gqlErr)
				assert.NotNil(t, resultSet)
				require.Empty(t, prettyErrors(resultSet.Errors))
			})
		})

		t.Run("multiclass with where filter (string)", func(t *testing.T) {
			wherePizza := filters.Where().
				WithPath([]string{"name"}).
				WithOperator(filters.Equal).
				WithValueString("Frutti di Mare")

			whereRisotto := filters.Where().
				WithPath([]string{"name"}).
				WithOperator(filters.Equal).
				WithValueString("Risotto alla pilota")

			resultSet, gqlErr := client.GraphQL().MultiClassGet().
				AddQueryClass(graphql.NewQueryClassBuilder("Pizza").
					WithWhere(wherePizza).
					WithFields(name)).
				AddQueryClass(graphql.NewQueryClassBuilder("Risotto").
					WithWhere(whereRisotto).
					WithFields(description)).
				Do(context.Background())
			assert.NoError(t, gqlErr)

			get := resultSet.Data["Get"].(map[string]interface{})
			pizza := get["Pizza"].([]interface{})
			assert.Equal(t, 1, len(pizza))
			risotto := get["Risotto"].([]interface{})
			assert.Equal(t, 1, len(risotto))
		})

		t.Run("multiclass with where filter (date)", func(t *testing.T) {
			// two pizzas have best_before dates in the test set, one of them expires
			// May 3rd, the other May 5th, so the filter below should match exactly
			// one.
			pizzaTargetDate, err := time.Parse(time.RFC3339, "2022-05-04T12:00:00+02:00")
			require.NoError(t, err)

			// two risotto have best_before dates in the test set, one of them expires
			// May 3rd, the other May 5th, so the filter below should match both
			risottoTargetDate, err := time.Parse(time.RFC3339, "2022-05-02T12:00:00+02:00")
			require.NoError(t, err)

			pizzaWhere := filters.Where().
				WithPath([]string{"best_before"}).
				WithOperator(filters.LessThan).
				WithValueDate(pizzaTargetDate)

			risottoWhere := filters.Where().
				WithPath([]string{"best_before"}).
				WithOperator(filters.GreaterThan).
				WithValueDate(risottoTargetDate)

			resultSet, gqlErr := client.GraphQL().MultiClassGet().
				AddQueryClass(graphql.NewQueryClassBuilder("Pizza").
					WithWhere(pizzaWhere).
					WithFields(name)).
				AddQueryClass(graphql.NewQueryClassBuilder("Risotto").
					WithWhere(risottoWhere).
					WithFields(name)).
				Do(context.Background())
			assert.NoError(t, gqlErr)

			get := resultSet.Data["Get"].(map[string]interface{})
			pizza := get["Pizza"].([]interface{})
			assert.Equal(t, 1, len(pizza))
			risotto := get["Risotto"].([]interface{})
			assert.Equal(t, 2, len(risotto))
		})

		t.Run("multiclass with sorting", func(t *testing.T) {
			var resultData []interface{}

			assertSortResult := func(resultSet *models.GraphQLResponse, className string, expected []string) {
				get := resultSet.Data["Get"].(map[string]interface{})
				resultData = get[className].([]interface{})
				require.Equal(t, len(expected), len(resultData))
				result := make([]string, len(resultData))
				for i := range resultData {
					p := resultData[i].(map[string]interface{})
					result[i] = p["name"].(string)
				}
				assert.Equal(t, expected, result)
			}

			byNameAsc := graphql.Sort{Path: []string{"name"}, Order: graphql.Asc}
			resultSet, gqlErr := client.GraphQL().MultiClassGet().
				AddQueryClass(graphql.NewQueryClassBuilder("Pizza").
					WithSort(byNameAsc).
					WithFields(name)).
				AddQueryClass(graphql.NewQueryClassBuilder("Risotto").
					WithSort(byNameAsc).
					WithFields(name)).
				Do(context.Background())
			assert.NoError(t, gqlErr)
			assertSortResult(resultSet, "Pizza", []string{"Doener", "Frutti di Mare", "Hawaii", "Quattro Formaggi"})
			assertSortResult(resultSet, "Risotto", []string{"Risi e bisi", "Risotto al nero di seppia", "Risotto alla pilota"})

			byNameDesc := graphql.Sort{Path: []string{"name"}, Order: graphql.Desc}
			resultSet, gqlErr = client.GraphQL().MultiClassGet().
				AddQueryClass(graphql.NewQueryClassBuilder("Pizza").
					WithSort(byNameDesc).
					WithFields(name)).
				AddQueryClass(graphql.NewQueryClassBuilder("Risotto").
					WithSort(byNameDesc).
					WithFields(name)).
				Do(context.Background())
			assert.NoError(t, gqlErr)
			assertSortResult(resultSet, "Pizza", []string{"Quattro Formaggi", "Hawaii", "Frutti di Mare", "Doener"})
			assertSortResult(resultSet, "Risotto", []string{"Risotto alla pilota", "Risotto al nero di seppia", "Risi e bisi"})

			byPriceAsc := graphql.Sort{Path: []string{"price"}, Order: graphql.Asc}
			resultSet, gqlErr = client.GraphQL().MultiClassGet().
				AddQueryClass(graphql.NewQueryClassBuilder("Soup").
					WithSort(byPriceAsc).
					WithFields(name)).
				AddQueryClass(graphql.NewQueryClassBuilder("Risotto").
					WithSort(byPriceAsc).
					WithFields(name)).
				Do(context.Background())
			assert.NoError(t, gqlErr)
			assertSortResult(resultSet, "Soup", []string{"ChickenSoup", "Beautiful"})
			assertSortResult(resultSet, "Risotto", []string{"Risi e bisi", "Risotto alla pilota", "Risotto al nero di seppia"})

			resultSet, gqlErr = client.GraphQL().MultiClassGet().
				AddQueryClass(graphql.NewQueryClassBuilder("Pizza").
					WithSort(byPriceAsc, byNameDesc).
					WithFields(name)).
				AddQueryClass(graphql.NewQueryClassBuilder("Risotto").
					WithSort(byPriceAsc, byNameDesc).
					WithFields(name)).
				Do(context.Background())

			assert.NoError(t, gqlErr)
			assertSortResult(resultSet, "Pizza", []string{"Quattro Formaggi", "Frutti di Mare", "Hawaii", "Doener"})
			assertSortResult(resultSet, "Risotto", []string{"Risi e bisi", "Risotto alla pilota", "Risotto al nero di seppia"})
		})
	})

	t.Run("MultiClassGet with group filter", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		testsuit.CreateTestSchemaAndData(t, client)
		defer testsuit.CleanUpWeaviate(t, client)

		name := graphql.Field{Name: "name"}
		group := client.GraphQL().GroupArgBuilder().WithType(graphql.Merge).WithForce(1.0)

		resultSet, gqlErr := client.GraphQL().MultiClassGet().
			AddQueryClass(graphql.NewQueryClassBuilder("Pizza").
				WithFields(name).
				WithGroup(group).
				WithLimit(7)).
			AddQueryClass(graphql.NewQueryClassBuilder("Risotto").
				WithFields(name).
				WithGroup(group).
				WithLimit(7)).
			Do(context.Background())
		assert.NoError(t, gqlErr)

		get := resultSet.Data["Get"].(map[string]interface{})
		require.Equal(t, 2, len(get))

		pizza := get["Pizza"].([]interface{})
		assert.Equal(t, 1, len(pizza))
		risotto := get["Pizza"].([]interface{})
		assert.Equal(t, 1, len(risotto))
	})

	t.Run("MultiClassGet with creationTimeUnix filters", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		testsuit.CreateTestSchemaAndData(t, client)
		defer testsuit.CleanUpWeaviate(t, client)

		pizza := GetOnePizza(t, client)
		pizzaExpectedCreateTime := pizza.Additional.CreationTimeUnix

		risotto := GetOneRisotto(t, client)
		risottoExpectedCreateTime := risotto.Additional.CreationTimeUnix

		additional := graphql.Field{
			Name: "_additional", Fields: []graphql.Field{
				{Name: "creationTimeUnix"},
			},
		}

		pizzaWhereCreateTime := filters.Where().
			WithPath([]string{"_creationTimeUnix"}).
			WithOperator(filters.Equal).
			WithValueString(pizzaExpectedCreateTime)

		risottoWhereCreateTime := filters.Where().
			WithPath([]string{"_creationTimeUnix"}).
			WithOperator(filters.Equal).
			WithValueString(risottoExpectedCreateTime)

		result, err := client.GraphQL().MultiClassGet().
			AddQueryClass(graphql.NewQueryClassBuilder("Pizza").
				WithFields(additional).
				WithWhere(pizzaWhereCreateTime)).
			AddQueryClass(graphql.NewQueryClassBuilder("Risotto").
				WithFields(additional).
				WithWhere(risottoWhereCreateTime)).
			Do(context.Background())

		require.NoError(t, err)
		require.Empty(t, prettyErrors(result.Errors))
		require.NotNil(t, result)
		require.NotNil(t, result.Data)

		b, err := json.Marshal(result.Data)
		require.NoError(t, err)

		var pizzaResp GetPizzaResponse
		err = json.Unmarshal(b, &pizzaResp)
		require.NoError(t, err)
		require.NotEmpty(t, pizzaResp.Get.Pizzas)
		assert.Equal(t, pizzaExpectedCreateTime, pizzaResp.Get.Pizzas[0].Additional.CreationTimeUnix)

		var risottoResp GetRisottoResponse
		err = json.Unmarshal(b, &risottoResp)
		require.NoError(t, err)
		require.NotEmpty(t, risottoResp.Get.Risotto)
		assert.Equal(t, risottoExpectedCreateTime, risottoResp.Get.Risotto[0].Additional.CreationTimeUnix)
	})

	t.Run("MultiClassGet with lastUpdateTimeUnix filters", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		testsuit.CreateTestSchemaAndData(t, client)
		defer testsuit.CleanUpWeaviate(t, client)

		pizza := GetOnePizza(t, client)
		pizzaExpectedUpdateTime := pizza.Additional.LastUpdateTimeUnix

		risotto := GetOneRisotto(t, client)
		risottoExpectedUpdateTime := risotto.Additional.LastUpdateTimeUnix

		additional := graphql.Field{
			Name: "_additional", Fields: []graphql.Field{
				{Name: "lastUpdateTimeUnix"},
			},
		}

		pizzaWhereUpdateTime := filters.Where().
			WithPath([]string{"_lastUpdateTimeUnix"}).
			WithOperator(filters.Equal).
			WithValueString(pizzaExpectedUpdateTime)

		risottoWhereUpdateTime := filters.Where().
			WithPath([]string{"_lastUpdateTimeUnix"}).
			WithOperator(filters.Equal).
			WithValueString(risottoExpectedUpdateTime)

		result, err := client.GraphQL().MultiClassGet().
			AddQueryClass(graphql.NewQueryClassBuilder("Pizza").
				WithFields(additional).
				WithWhere(pizzaWhereUpdateTime)).
			AddQueryClass(graphql.NewQueryClassBuilder("Risotto").
				WithFields(additional).
				WithWhere(risottoWhereUpdateTime)).
			Do(context.Background())

		require.NoError(t, err)
		require.Empty(t, prettyErrors(result.Errors))
		require.NotNil(t, result)
		require.NotNil(t, result.Data)

		b, err := json.Marshal(result.Data)
		require.NoError(t, err)

		var pizzaResp GetPizzaResponse
		err = json.Unmarshal(b, &pizzaResp)
		require.NoError(t, err)
		require.NotEmpty(t, pizzaResp.Get.Pizzas)
		assert.Equal(t, pizzaExpectedUpdateTime, pizzaResp.Get.Pizzas[0].Additional.LastUpdateTimeUnix)

		var risottoResp GetRisottoResponse
		err = json.Unmarshal(b, &risottoResp)
		require.NoError(t, err)
		require.NotEmpty(t, risottoResp.Get.Risotto)
		assert.Equal(t, risottoExpectedUpdateTime, risottoResp.Get.Risotto[0].Additional.LastUpdateTimeUnix)
	})

	t.Run("group by with nearObject", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		testsuit.CreateTestDocumentAndPassageSchemaAndData(t, client)
		defer testsuit.CleanUpWeaviate(t, client)

		additional := graphql.Field{
			Name: "_additional", Fields: []graphql.Field{
				{Name: "group", Fields: []graphql.Field{
					{Name: "id"},
					{Name: "groupedBy", Fields: []graphql.Field{
						{Name: "value"},
						{Name: "path"},
					}},
					{Name: "count"},
					{Name: "maxDistance"},
					{Name: "minDistance"},
					{Name: "hits", Fields: []graphql.Field{
						{Name: "ofDocument{... on Document{_additional{id}}}"},
						{Name: "_additional", Fields: []graphql.Field{
							{Name: "id"},
							{Name: "distance"},
						}},
					}},
				}},
			},
		}

		groupBy := client.GraphQL().GroupByArgBuilder().
			WithPath([]string{"ofDocument"}).WithGroups(3).WithObjectsPerGroup(10)

		nearObject := client.GraphQL().NearObjectArgBuilder().
			WithID("00000000-0000-0000-0000-000000000001")

		result, err := client.GraphQL().Get().
			WithClassName("Passage").
			WithNearObject(nearObject).
			WithGroupBy(groupBy).
			WithFields(additional).
			Do(context.Background())

		require.NoError(t, err)
		require.Empty(t, prettyErrors(result.Errors))
		require.NotNil(t, result)
		require.NotNil(t, result.Data)

		getGroup := func(value interface{}) map[string]interface{} {
			group := value.(map[string]interface{})["_additional"].(map[string]interface{})["group"].(map[string]interface{})
			return group
		}
		groups := []map[string]interface{}{}
		passages := result.Data["Get"].(map[string]interface{})["Passage"].([]interface{})
		for _, passage := range passages {
			groups = append(groups, getGroup(passage))
		}
		getGroupHits := func(group map[string]interface{}) (string, []string) {
			result := []string{}
			hits := group["hits"].([]interface{})
			for _, hit := range hits {
				additional := hit.(map[string]interface{})["_additional"].(map[string]interface{})
				result = append(result, additional["id"].(string))
			}
			groupedBy := group["groupedBy"].(map[string]interface{})
			groupedByValue := groupedBy["value"].(string)
			return groupedByValue, result
		}

		require.Len(t, groups, 3)
		expectedGroups := map[string][]string{}
		group1 := "weaviate://localhost/Document/00000000-0000-0000-0000-00000000000a"
		expectedGroups[group1] = []string{
			"00000000-0000-0000-0000-000000000001",
			"00000000-0000-0000-0000-000000000009",
			"00000000-0000-0000-0000-000000000007",
			"00000000-0000-0000-0000-000000000008",
			"00000000-0000-0000-0000-000000000006",
			"00000000-0000-0000-0000-000000000010",
			"00000000-0000-0000-0000-000000000005",
			"00000000-0000-0000-0000-000000000004",
			"00000000-0000-0000-0000-000000000003",
			"00000000-0000-0000-0000-000000000002",
		}
		group2 := "weaviate://localhost/Document/00000000-0000-0000-0000-00000000000b"
		expectedGroups[group2] = []string{
			"00000000-0000-0000-0000-000000000011",
			"00000000-0000-0000-0000-000000000013",
			"00000000-0000-0000-0000-000000000012",
			"00000000-0000-0000-0000-000000000014",
		}
		group3 := ""
		expectedGroups[group3] = []string{
			"00000000-0000-0000-0000-000000000016",
			"00000000-0000-0000-0000-000000000017",
			"00000000-0000-0000-0000-000000000015",
			"00000000-0000-0000-0000-000000000020",
			"00000000-0000-0000-0000-000000000019",
			"00000000-0000-0000-0000-000000000018",
		}

		groupsOrder := []string{group1, group2, group3}
		for i, current := range groups {
			groupedBy, ids := getGroupHits(current)
			assert.Equal(t, groupsOrder[i], groupedBy)
			assert.ElementsMatch(t, expectedGroups[groupedBy], ids)
		}
	})

	t.Run("group by with hybrid", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		testsuit.CreateTestDocumentAndPassageSchemaAndData(t, client)
		defer testsuit.CleanUpWeaviate(t, client)

		additional := graphql.Field{
			Name: "_additional", Fields: []graphql.Field{
				{Name: "group", Fields: []graphql.Field{
					{Name: "id"},
					{Name: "groupedBy", Fields: []graphql.Field{
						{Name: "value"},
						{Name: "path"},
					}},
					{Name: "count"},
					{Name: "maxDistance"},
					{Name: "minDistance"},
					{Name: "hits", Fields: []graphql.Field{
						{Name: "_additional", Fields: []graphql.Field{
							{Name: "id"},
							{Name: "distance"},
						}},
					}},
				}},
			},
		}

		groupBy := client.GraphQL().GroupByArgBuilder().
			WithPath([]string{"content"}).WithGroups(3).WithObjectsPerGroup(10)

		searches := client.GraphQL().HybridSearchesArgumentBuilder().
			WithNearText(client.GraphQL().NearTextArgBuilder().WithConcepts([]string{"Passage content 2"}))

		hybrid := client.GraphQL().HybridArgumentBuilder().
			WithQuery("Passage content 2").
			WithSearches(searches).
			WithAlpha(0.9)

		result, err := client.GraphQL().Get().
			WithClassName("Passage").
			WithHybrid(hybrid).
			WithGroupBy(groupBy).
			WithFields(additional).
			Do(context.Background())

		require.NoError(t, err)
		require.Empty(t, prettyErrors(result.Errors))
		require.NotNil(t, result)
		require.NotNil(t, result.Data)

		getGroup := func(value interface{}) map[string]interface{} {
			group := value.(map[string]interface{})["_additional"].(map[string]interface{})["group"].(map[string]interface{})
			return group
		}
		groups := []map[string]interface{}{}
		passages := result.Data["Get"].(map[string]interface{})["Passage"].([]interface{})
		for _, passage := range passages {
			groups = append(groups, getGroup(passage))
		}
		getGroupHits := func(group map[string]interface{}) (string, []string) {
			result := []string{}
			hits := group["hits"].([]interface{})
			for _, hit := range hits {
				additional := hit.(map[string]interface{})["_additional"].(map[string]interface{})
				result = append(result, additional["id"].(string))
			}
			groupedBy := group["groupedBy"].(map[string]interface{})
			groupedByValue := groupedBy["value"].(string)
			return groupedByValue, result
		}

		require.Len(t, groups, 3)
		group1 := "Passage content 2"
		group2 := "Passage content 1"
		group3 := "Passage content 3"

		groupsOrder := []string{group1, group2, group3}
		for i, current := range groups {
			groupedBy, _ := getGroupHits(current)
			assert.Equal(t, groupsOrder[i], groupedBy)
		}
	})

	t.Run("query with consistency level", func(t *testing.T) {
		ctx := context.Background()
		container, stop := testenv.SetupLocalContainer(t, ctx, test.Cluster, true)
		t.Cleanup(stop)
		client := weaviate.New(weaviate.Config{
			Host:   container.HTTPAddress(),
			Scheme: "http",
		})
		fields := []graphql.Field{
			{
				Name: "name",
			},
			{
				Name: "_additional",
				Fields: []graphql.Field{
					{
						Name: "isConsistent",
					},
				},
			},
		}

		testsuit.CreateTestSchemaAndData(t, client, testsuit.WithReplication)
		defer testsuit.CleanUpWeaviate(t, client)

		t.Run("All", func(t *testing.T) {
			resp, err := client.GraphQL().Get().
				WithClassName("Pizza").
				WithFields(fields...).
				WithConsistencyLevel(replication.ConsistencyLevel.ALL).
				Do(ctx)
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Empty(t, prettyErrors(resp.Errors))
		})

		t.Run("Quorum", func(t *testing.T) {
			resp, err := client.GraphQL().Get().
				WithClassName("Pizza").
				WithFields(fields...).
				WithConsistencyLevel(replication.ConsistencyLevel.QUORUM).
				Do(ctx)
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Empty(t, prettyErrors(resp.Errors))
		})

		t.Run("One", func(t *testing.T) {
			resp, err := client.GraphQL().Get().
				WithClassName("Pizza").
				WithFields(fields...).
				WithConsistencyLevel(replication.ConsistencyLevel.ONE).
				Do(ctx)
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Empty(t, prettyErrors(resp.Errors))
		})
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to tear down weaviate: %s", err)
		}
	})
}

func TestGraphQL_MultiTenancy(t *testing.T) {
	cleanup := func() {
		client := testsuit.CreateTestClient(false)
		err := client.Schema().AllDeleter().Do(context.Background())
		require.NoError(t, err)
	}

	t.Run("setup weaviate", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to setup weaviate: %s", err)
		}
	})

	t.Run("GraphQL Get", func(t *testing.T) {
		defer cleanup()

		tenant1 := models.Tenant{Name: "tenantNo1"}
		tenant2 := models.Tenant{Name: "tenantNo2"}
		client := testsuit.CreateTestClient(false)

		assertGetContainsIds := func(t *testing.T, response *models.GraphQLResponse,
			className string, expectedIds []string,
		) {
			require.NotNil(t, response)
			assert.Nil(t, response.Errors)
			require.NotNil(t, response.Data)

			get := response.Data["Get"].(map[string]interface{})
			objects := get[className].([]interface{})
			require.Len(t, objects, len(expectedIds))

			ids := []string{}
			for i := range objects {
				ids = append(ids, objects[i].(map[string]interface{})["_additional"].(map[string]interface{})["id"].(string))
			}
			assert.ElementsMatch(t, expectedIds, ids)
		}

		t.Run("add data", func(t *testing.T) {
			testsuit.CreateSchemaPizzaForTenants(t, client)
			testsuit.CreateTenantsPizza(t, client, tenant1, tenant2)
			testsuit.CreateDataPizzaQuattroFormaggiForTenants(t, client, tenant1.Name)
			testsuit.CreateDataPizzaFruttiDiMareForTenants(t, client, tenant1.Name)
			testsuit.CreateDataPizzaHawaiiForTenants(t, client, tenant2.Name)
			testsuit.CreateDataPizzaDoenerForTenants(t, client, tenant2.Name)
		})

		t.Run("get all data for tenant", func(t *testing.T) {
			expectedIdsByTenant := map[string][]string{
				tenant1.Name: {
					testsuit.PIZZA_QUATTRO_FORMAGGI_ID,
					testsuit.PIZZA_FRUTTI_DI_MARE_ID,
				},
				tenant2.Name: {
					testsuit.PIZZA_HAWAII_ID,
					testsuit.PIZZA_DOENER_ID,
				},
			}

			for tenant, expectedIds := range expectedIdsByTenant {
				resp, err := client.GraphQL().Get().
					WithClassName("Pizza").
					WithTenant(tenant).
					WithFields(graphql.Field{
						Name:   "_additional",
						Fields: []graphql.Field{{Name: "id"}},
					}).
					Do(context.Background())

				assert.Nil(t, err)
				assertGetContainsIds(t, resp, "Pizza", expectedIds)
			}
		})

		t.Run("get limited data for tenant", func(t *testing.T) {
			expectedIdsByTenant := map[string][]string{
				tenant1.Name: {
					testsuit.PIZZA_QUATTRO_FORMAGGI_ID,
				},
				tenant2.Name: {
					testsuit.PIZZA_HAWAII_ID,
				},
			}

			for tenant, expectedIds := range expectedIdsByTenant {
				resp, err := client.GraphQL().Get().
					WithClassName("Pizza").
					WithTenant(tenant).
					WithLimit(1).
					WithFields(graphql.Field{
						Name:   "_additional",
						Fields: []graphql.Field{{Name: "id"}},
					}).
					Do(context.Background())

				assert.Nil(t, err)
				assertGetContainsIds(t, resp, "Pizza", expectedIds)
			}
		})

		t.Run("get filtered data for tenant", func(t *testing.T) {
			expectedIdsByTenant := map[string][]string{
				tenant1.Name: {},
				tenant2.Name: {
					testsuit.PIZZA_DOENER_ID,
				},
			}
			where := filters.Where().
				WithPath([]string{"price"}).
				WithOperator(filters.GreaterThan).
				WithValueNumber(1.3)

			for tenant, expectedIds := range expectedIdsByTenant {
				resp, err := client.GraphQL().Get().
					WithClassName("Pizza").
					WithTenant(tenant).
					WithWhere(where).
					WithFields(graphql.Field{
						Name:   "_additional",
						Fields: []graphql.Field{{Name: "id"}},
					}).
					Do(context.Background())

				assert.Nil(t, err)
				assertGetContainsIds(t, resp, "Pizza", expectedIds)
			}
		})
	})

	t.Run("GraphQL Aggregate", func(t *testing.T) {
		defer cleanup()

		tenant1 := models.Tenant{Name: "tenantNo1"}
		tenant2 := models.Tenant{Name: "tenantNo2"}
		client := testsuit.CreateTestClient(false)

		assertAggregateNumFieldHasValues := func(t *testing.T, response *models.GraphQLResponse,
			className string, fieldName string, expectedAggValues map[string]*float64,
		) {
			require.NotNil(t, response)
			assert.Nil(t, response.Errors)
			require.NotNil(t, response.Data)

			agg := response.Data["Aggregate"].(map[string]interface{})
			objects := agg[className].([]interface{})
			require.Len(t, objects, 1)
			obj := objects[0].(map[string]interface{})[fieldName].(map[string]interface{})

			for name, value := range expectedAggValues {
				if value == nil {
					assert.Nil(t, obj[name])
				} else {
					assert.Equal(t, *value, obj[name])
				}
			}
		}
		ptr := func(f float64) *float64 {
			return &f
		}

		t.Run("add data", func(t *testing.T) {
			testsuit.CreateSchemaPizzaForTenants(t, client)
			testsuit.CreateTenantsPizza(t, client, tenant1, tenant2)
			testsuit.CreateDataPizzaQuattroFormaggiForTenants(t, client, tenant1.Name)
			testsuit.CreateDataPizzaFruttiDiMareForTenants(t, client, tenant1.Name)
			testsuit.CreateDataPizzaHawaiiForTenants(t, client, tenant2.Name)
			testsuit.CreateDataPizzaDoenerForTenants(t, client, tenant2.Name)
		})

		t.Run("aggregate all data for tenant", func(t *testing.T) {
			expectedAggValuesByTenant := map[string]map[string]*float64{
				tenant1.Name: {
					"count":   ptr(2),
					"maximum": ptr(1.2),
					"minimum": ptr(1.1),
					"median":  ptr(1.15),
					"mean":    ptr(1.15),
					"mode":    ptr(1.1),
					"sum":     ptr(2.3),
				},
				tenant2.Name: {
					"count":   ptr(2),
					"maximum": ptr(1.4),
					"minimum": ptr(1.3),
					"median":  ptr(1.35),
					"mean":    ptr(1.35),
					"mode":    ptr(1.3),
					"sum":     ptr(2.7),
				},
			}

			for tenant, expectedAggValues := range expectedAggValuesByTenant {
				resp, err := client.GraphQL().Aggregate().
					WithClassName("Pizza").
					WithTenant(tenant).
					WithFields(graphql.Field{
						Name: "price",
						Fields: []graphql.Field{
							{Name: "count"},
							{Name: "maximum"},
							{Name: "minimum"},
							{Name: "median"},
							{Name: "mean"},
							{Name: "mode"},
							{Name: "sum"},
						},
					}).
					Do(context.Background())

				assert.Nil(t, err)
				assertAggregateNumFieldHasValues(t, resp, "Pizza", "price", expectedAggValues)
			}
		})

		t.Run("aggregate filtered data for tenant", func(t *testing.T) {
			expectedAggValuesByTenant := map[string]map[string]*float64{
				tenant1.Name: {
					"count":   ptr(0),
					"maximum": nil,
					"minimum": nil,
					"median":  nil,
					"mean":    nil,
					"mode":    nil,
					"sum":     nil,
				},
				tenant2.Name: {
					"count":   ptr(1),
					"maximum": ptr(1.4),
					"minimum": ptr(1.4),
					"median":  ptr(1.4),
					"mean":    ptr(1.4),
					"mode":    ptr(1.4),
					"sum":     ptr(1.4),
				},
			}
			where := filters.Where().
				WithPath([]string{"price"}).
				WithOperator(filters.GreaterThan).
				WithValueNumber(1.3)

			for tenant, expectedAggValues := range expectedAggValuesByTenant {
				resp, err := client.GraphQL().Aggregate().
					WithClassName("Pizza").
					WithTenant(tenant).
					WithWhere(where).
					WithFields(graphql.Field{
						Name: "price",
						Fields: []graphql.Field{
							{Name: "count"},
							{Name: "maximum"},
							{Name: "minimum"},
							{Name: "median"},
							{Name: "mean"},
							{Name: "mode"},
							{Name: "sum"},
						},
					}).
					Do(context.Background())

				assert.Nil(t, err)
				assertAggregateNumFieldHasValues(t, resp, "Pizza", "price", expectedAggValues)
			}
		})

		t.Run("ContainsAny and ContainsAll", func(t *testing.T) {
			className := "WhereTest"
			data := testsuit.AllPropertiesData()
			id1, id2, id3 := data.ID1, data.ID2, data.ID3
			ids := []string{id1, id2, id3}

			cleanup := func() {
				err := client.Schema().AllDeleter().Do(context.Background())
				require.NoError(t, err)
			}

			t.Run("create class", func(t *testing.T) {
				cleanup()
				testsuit.AllPropertiesSchemaCreate(t, client, className, false, false)
			})

			t.Run("insert data", func(t *testing.T) {
				authors := data.Authors
				authorsArray := data.AuthorsArray
				colors := data.Colors
				colorssArray := data.ColorssArray
				numbers := data.Numbers
				numbersArray := data.NumbersArray
				ints := data.Ints
				intsArray := data.IntsArray
				uuids := data.Uuids
				uuidsArray := data.UuidsArray
				dates := data.Dates
				datesArray := data.DatesArray
				bools := data.Bools
				boolsArray := data.BoolsArray
				for i, id := range ids {
					_, err := client.Data().Creator().
						WithClassName(className).
						WithID(id).
						WithProperties(map[string]interface{}{
							"color":   colors[i],
							"colors":  colorssArray[i],
							"author":  authors[i],
							"authors": authorsArray[i],
							"number":  numbers[i],
							"numbers": numbersArray[i],
							"int":     ints[i],
							"ints":    intsArray[i],
							"uuid":    uuids[i],
							"uuids":   uuidsArray[i],
							"date":    dates[i],
							"dates":   datesArray[i],
							"bool":    bools[i],
							"bools":   boolsArray[i],
						}).
						Do(context.TODO())
					require.NoError(t, err)
				}
			})

			t.Run("where", func(t *testing.T) {
				mustGetTime := func(date string) time.Time {
					result, err := time.Parse(time.RFC3339Nano, date)
					if err != nil {
						panic(fmt.Sprintf("can't parse date: %v", date))
					}
					return result
				}
				getIds := func(data map[string]models.JSONObject) []string {
					classMap, ok := data["Get"].(map[string]interface{})
					require.True(t, ok)
					class, ok := classMap[className].([]interface{})
					require.True(t, ok)
					ids := make([]string, len(class))
					for i := range class {
						resultMap, ok := class[i].(map[string]interface{})
						require.True(t, ok)
						additional, ok := resultMap["_additional"].(map[string]interface{})
						require.True(t, ok)
						ids[i] = additional["id"].(string)
					}
					return ids
				}
				tests := []struct {
					name        string
					where       *filters.WhereBuilder
					property    string
					expectedIds []string
				}{
					// Contains operator with array types
					{
						name: "contains all authors with string array",
						where: filters.Where().
							WithPath([]string{"authors"}).
							WithOperator(filters.ContainsAll).
							WithValueString("John", "Jenny", "Joseph"),
						property:    "authors",
						expectedIds: []string{id1},
					},
					{
						name: "contains any authors with string array",
						where: filters.Where().
							WithPath([]string{"authors"}).
							WithOperator(filters.ContainsAny).
							WithValueString("John", "Jenny", "Joseph"),
						property:    "authors",
						expectedIds: []string{id1, id2, id3},
					},
					{
						name: "contains all colors with text array",
						where: filters.Where().
							WithPath([]string{"colors"}).
							WithOperator(filters.ContainsAll).
							WithValueText("red", "blue", "green"),
						property:    "colors",
						expectedIds: []string{id1},
					},
					{
						name: "contains any colors with text array",
						where: filters.Where().
							WithPath([]string{"colors"}).
							WithOperator(filters.ContainsAny).
							WithValueText("red", "blue", "green"),
						property:    "colors",
						expectedIds: []string{id1, id2, id3},
					},
					{
						name: "contains all numbers with number array",
						where: filters.Where().
							WithPath([]string{"numbers"}).
							WithOperator(filters.ContainsAll).
							WithValueNumber(1.1, 2.2, 3.3),
						property:    "numbers",
						expectedIds: []string{id1},
					},
					{
						name: "contains any numbers with number array",
						where: filters.Where().
							WithPath([]string{"numbers"}).
							WithOperator(filters.ContainsAny).
							WithValueNumber(1.1, 2.2, 3.3),
						property:    "numbers",
						expectedIds: []string{id1, id2, id3},
					},
					{
						name: "contains all ints with int array",
						where: filters.Where().
							WithPath([]string{"ints"}).
							WithOperator(filters.ContainsAll).
							WithValueInt(1, 2, 3),
						property:    "ints",
						expectedIds: []string{id1},
					},
					{
						name: "contains any ints with int array",
						where: filters.Where().
							WithPath([]string{"ints"}).
							WithOperator(filters.ContainsAny).
							WithValueInt(1, 2, 3),
						property:    "ints",
						expectedIds: []string{id1, id2, id3},
					},
					{
						name: "contains all uuids with uuid array",
						where: filters.Where().
							WithPath([]string{"uuids"}).
							WithOperator(filters.ContainsAll).
							WithValueText(id1, id2, id3),
						property:    "uuids",
						expectedIds: []string{id1},
					},
					{
						name: "contains any uuids with uuid array",
						where: filters.Where().
							WithPath([]string{"uuids"}).
							WithOperator(filters.ContainsAny).
							WithValueText(id1, id2, id3),
						property:    "uuids",
						expectedIds: []string{id1, id2, id3},
					},
					{
						name: "contains all dates with dates array",
						where: filters.Where().
							WithPath([]string{"dates"}).
							WithOperator(filters.ContainsAll).
							WithValueDate(
								mustGetTime("2009-11-01T23:00:00Z"), mustGetTime("2009-11-02T23:00:00Z"), mustGetTime("2009-11-03T23:00:00Z"),
							),
						property:    "dates",
						expectedIds: []string{id1},
					},
					{
						name: "contains any dates with dates array",
						where: filters.Where().
							WithPath([]string{"dates"}).
							WithOperator(filters.ContainsAny).
							WithValueDate(
								mustGetTime("2009-11-01T23:00:00Z"), mustGetTime("2009-11-02T23:00:00Z"), mustGetTime("2009-11-03T23:00:00Z"),
							),
						property:    "dates",
						expectedIds: []string{id1, id2, id3},
					},
					{
						name: "complex contains all ints and all numbers with AND on int array",
						where: filters.Where().
							WithOperator(filters.And).
							WithOperands([]*filters.WhereBuilder{
								filters.Where().
									WithPath([]string{"numbers"}).
									WithOperator(filters.ContainsAll).
									WithValueNumber(1.1, 2.2, 3.3),
								filters.Where().
									WithPath([]string{"ints"}).
									WithOperator(filters.ContainsAll).
									WithValueInt(1, 2, 3),
							}),
						property:    "ints",
						expectedIds: []string{id1},
					},
					{
						name: "complex contains any ints and all numbers with OR on int array",
						where: filters.Where().
							WithOperator(filters.Or).
							WithOperands([]*filters.WhereBuilder{
								filters.Where().
									WithPath([]string{"numbers"}).
									WithOperator(filters.ContainsAll).
									WithValueNumber(1.1, 2.2, 3.3),
								filters.Where().
									WithPath([]string{"ints"}).
									WithOperator(filters.ContainsAny).
									WithValueInt(1, 2, 3),
							}),
						property:    "ints",
						expectedIds: []string{id1, id2, id3},
					},
					// Contains operator with primitives
					{
						name: "contains any author with string",
						where: filters.Where().
							WithPath([]string{"author"}).
							WithOperator(filters.ContainsAny).
							WithValueString("John", "Jenny", "Joseph"),
						property:    "author",
						expectedIds: []string{id1, id2, id3},
					},
					{
						name: "contains any color with text",
						where: filters.Where().
							WithPath([]string{"color"}).
							WithOperator(filters.ContainsAny).
							WithValueText("red", "blue", "green"),
						property:    "color",
						expectedIds: []string{id1, id2, id3},
					},
					{
						name: "contains any number with number",
						where: filters.Where().
							WithPath([]string{"number"}).
							WithOperator(filters.ContainsAny).
							WithValueNumber(1.1, 2.2, 3.3),
						property:    "number",
						expectedIds: []string{id1, id2, id3},
					},
					{
						name: "contains any int with int",
						where: filters.Where().
							WithPath([]string{"int"}).
							WithOperator(filters.ContainsAny).
							WithValueInt(1, 2, 3),
						property:    "int",
						expectedIds: []string{id1, id2, id3},
					},
					{
						name: "contains any uuid with uuid",
						where: filters.Where().
							WithPath([]string{"uuid"}).
							WithOperator(filters.ContainsAny).
							WithValueText(id1, id2, id3),
						property:    "uuid",
						expectedIds: []string{id1, id2, id3},
					},
					{
						name: "contains any uuid with id",
						where: filters.Where().
							WithPath([]string{"id"}).
							WithOperator(filters.ContainsAny).
							WithValueText(id1, id2, id3),
						property:    "uuid",
						expectedIds: []string{id1, id2, id3},
					},
					{
						name: "contains any date with date",
						where: filters.Where().
							WithPath([]string{"date"}).
							WithOperator(filters.ContainsAny).
							WithValueDate(
								mustGetTime("2009-11-01T23:00:00Z"), mustGetTime("2009-11-02T23:00:00Z"), mustGetTime("2009-11-03T23:00:00Z"),
							),
						property:    "date",
						expectedIds: []string{id1, id2, id3},
					},
				}
				for _, tt := range tests {
					t.Run(tt.name, func(t *testing.T) {
						fields := []graphql.Field{
							{Name: tt.property},
							{Name: "_additional", Fields: []graphql.Field{{Name: "id"}}},
						}
						resp, err := client.GraphQL().Get().
							WithClassName(className).
							WithWhere(tt.where).
							WithFields(fields...).
							Do(context.TODO())
						require.NoError(t, err)
						require.Empty(t, prettyErrors(resp.Errors))
						resultIds := getIds(resp.Data)
						assert.Len(t, resultIds, len(tt.expectedIds))
						assert.ElementsMatch(t, resultIds, tt.expectedIds)
					})
				}
			})
		})
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to tear down weaviate: %s", err)
		}
	})
}

// prettyErrors prints error messages from the GraphQL error response
// instead of struct pointers.
// Use it to create more informative test errors:
//
//	require.Empty(t, prettyErrors(result.Errors))
type prettyErrors []*models.GraphQLError

func (errors prettyErrors) String() string {
	var msgs []string
	for i := range errors {
		msgs = append(msgs, errors[i].Message)
	}
	return fmt.Sprintf("%v", msgs)
}
