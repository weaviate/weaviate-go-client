package graphql

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v4/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/testenv"
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
		client := testsuit.CreateTestClient(8080, nil)
		testsuit.CreateTestSchemaAndData(t, client)
		defer testsuit.CleanUpWeaviate(t, client)

		name := graphql.Field{Name: "name"}

		// what what
		t.Run("get raw", func(t *testing.T) {
			resultSet, gqlErr := client.GraphQL().Raw().WithQuery("{Get {Pizza {name}}}").Do(context.Background())
			assert.Nil(t, gqlErr)

			get := resultSet.Data["Get"].(map[string]interface{})
			pizza := get["Pizza"].([]interface{})
			assert.Equal(t, 4, len(pizza))
		})

		t.Run("get all", func(t *testing.T) {
			resultSet, gqlErr := client.GraphQL().Get().WithClassName("Pizza").WithFields(name).Do(context.Background())
			assert.Nil(t, gqlErr)

			get := resultSet.Data["Get"].(map[string]interface{})
			pizza := get["Pizza"].([]interface{})
			assert.Equal(t, 4, len(pizza))
		})

		t.Run("by near object", func(t *testing.T) {
			withNearObject := client.GraphQL().NearObjectArgBuilder().
				WithID("5b6a08ba-1d46-43aa-89cc-8b070790c6f2")
			resultSet, gqlErr := client.GraphQL().Get().
				WithClassName("Pizza").
				WithFields(name).
				WithNearObject(withNearObject).
				Do(context.Background())
			assert.Nil(t, gqlErr)
			require.Nil(t, resultSet.Errors)

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
						{ID: "5b6a08ba-1d46-43aa-89cc-8b070790c6f1"},
						{Beacon: "weaviate://localhost/5b6a08ba-1d46-43aa-89cc-8b070790c6f1"},
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

				assert.Nil(t, gqlErr)
				assert.NotNil(t, resultSet)
				assert.Nil(t, resultSet.Errors)
			})

			t.Run("with distance", func(t *testing.T) {
				concepts := []string{"pineapple slices", "ham"}
				moveTo := &graphql.MoveParameters{
					Force: 0.3,
					Objects: []graphql.MoverObject{
						{ID: "5b6a08ba-1d46-43aa-89cc-8b070790c6f1"},
						{Beacon: "weaviate://localhost/5b6a08ba-1d46-43aa-89cc-8b070790c6f1"},
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

				assert.Nil(t, gqlErr)
				assert.NotNil(t, resultSet)
				assert.Nil(t, resultSet.Errors)
			})
		})

		t.Run("with where filter (string)", func(t *testing.T) {
			where := filters.Where().
				WithPath([]string{"name"}).
				WithOperator(filters.Equal).
				WithValueString("Frutti di Mare")

			name = graphql.Field{Name: "name"}

			resultSet, gqlErr := client.GraphQL().Get().WithClassName("Pizza").WithWhere(where).WithFields(name).Do(context.Background())
			assert.Nil(t, gqlErr)

			get := resultSet.Data["Get"].(map[string]interface{})
			pizza := get["Pizza"].([]interface{})
			assert.Equal(t, 1, len(pizza))
		})

		t.Run("with where filter (date)", func(t *testing.T) {
			// two pizzas have best_before dates in the test set, one of them expires
			// May 3rd, the other May 5th, so the filter below should match exactly
			// one.
			targetDate, err := time.Parse(time.RFC3339, "2022-05-04T12:00:00+02:00")
			require.Nil(t, err)

			where := filters.Where().
				WithPath([]string{"best_before"}).
				WithOperator(filters.LessThan).
				WithValueDate(targetDate)

			name = graphql.Field{Name: "name"}

			resultSet, gqlErr := client.GraphQL().Get().WithClassName("Pizza").WithWhere(where).WithFields(name).Do(context.Background())
			assert.Nil(t, gqlErr)

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
			assert.Nil(t, gqlErr)
			assertSortResult(resultSet, "Pizza", []string{"Doener", "Frutti di Mare", "Hawaii", "Quattro Formaggi"})

			byNameDesc := graphql.Sort{Path: []string{"name"}, Order: graphql.Desc}
			resultSet, gqlErr = client.GraphQL().Get().WithClassName("Pizza").
				WithSort(byNameDesc).WithFields(name).Do(context.Background())
			assert.Nil(t, gqlErr)
			assertSortResult(resultSet, "Pizza", []string{"Quattro Formaggi", "Hawaii", "Frutti di Mare", "Doener"})

			byPriceAsc := graphql.Sort{Path: []string{"price"}, Order: graphql.Asc}
			resultSet, gqlErr = client.GraphQL().Get().WithClassName("Soup").
				WithSort(byPriceAsc).WithFields(name).Do(context.Background())
			assert.Nil(t, gqlErr)
			assertSortResult(resultSet, "Soup", []string{"ChickenSoup", "Beautiful"})

			resultSet, gqlErr = client.GraphQL().Get().WithClassName("Pizza").
				WithSort(byPriceAsc, byNameDesc).WithFields(name).Do(context.Background())
			assert.Nil(t, gqlErr)
			assertSortResult(resultSet, "Pizza", []string{"Quattro Formaggi", "Frutti di Mare", "Hawaii", "Doener"})
		})
	})

	t.Run("Explore", func(t *testing.T) {
		t.Run("with certainty", func(t *testing.T) {
			client := testsuit.CreateTestClient(8080, nil)
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

			assert.Nil(t, gqlErr)
			assert.NotNil(t, resultSet)

			withNearObject := client.GraphQL().NearObjectArgBuilder().
				WithID("5b6a08ba-1d46-43aa-89cc-8b070790c6f2")

			resultSet, gqlErr = client.GraphQL().Explore().
				WithFields(graphql.Certainty, graphql.Beacon, graphql.ClassName).
				WithNearObject(withNearObject).
				Do(context.Background())

			assert.Nil(t, gqlErr)
			assert.NotNil(t, resultSet)
		})

		t.Run("with distance", func(t *testing.T) {
			client := testsuit.CreateTestClient(8080, nil)
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

			assert.Nil(t, gqlErr)
			assert.NotNil(t, resultSet)

			withNearObject := client.GraphQL().NearObjectArgBuilder().
				WithID("5b6a08ba-1d46-43aa-89cc-8b070790c6f2")

			resultSet, gqlErr = client.GraphQL().Explore().
				WithFields(graphql.Certainty, graphql.Beacon, graphql.ClassName).
				WithNearObject(withNearObject).
				Do(context.Background())

			assert.Nil(t, gqlErr)
			assert.NotNil(t, resultSet)
		})
	})

	t.Run("Aggregate", func(t *testing.T) {
		client := testsuit.CreateTestClient(8080, nil)
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

			assert.Nil(t, gqlErr)
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

			assert.Nil(t, gqlErr)
			assert.NotNil(t, resultSet)
		})

		t.Run("with groupby", func(t *testing.T) {
			resultSet, gqlErr := client.GraphQL().
				Aggregate().
				WithFields(meta).
				WithGroupBy("name").
				WithClassName("Pizza").
				Do(context.Background())

			assert.Nil(t, gqlErr)
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

			assert.Nil(t, gqlErr)
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

				assert.Nil(t, gqlErr)
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

				assert.Nil(t, gqlErr)
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

				assert.Nil(t, gqlErr)
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

				assert.Nil(t, gqlErr)
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

				assert.Nil(t, gqlErr)
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

				assert.Nil(t, gqlErr)
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

				assert.Nil(t, gqlErr)
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

				assert.Nil(t, gqlErr)
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

				assert.Nil(t, gqlErr)
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

				assert.Nil(t, gqlErr)
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

				assert.Nil(t, gqlErr)
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

				assert.Nil(t, gqlErr)
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

				assert.Nil(t, gqlErr)
				assert.NotNil(t, resultSet)
				assert.NotNil(t, resultSet.Data)

				b, err := json.Marshal(resultSet.Data)
				require.Nil(t, err)

				var resp AggregatePizzaResponse
				err = json.Unmarshal(b, &resp)
				require.Nil(t, err)

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

				assert.Nil(t, gqlErr)
				assert.NotNil(t, resultSet)
				assert.NotNil(t, resultSet.Data)

				b, err := json.Marshal(resultSet.Data)
				require.Nil(t, err)

				var resp AggregatePizzaResponse
				err = json.Unmarshal(b, &resp)
				require.Nil(t, err)

				assert.NotEmpty(t, resp.Aggregate.Pizza)
				assert.Equal(t, objectLimit, resp.Aggregate.Pizza[0].Meta.Count)
			})
		})
	})

	t.Run("Get with group filter", func(t *testing.T) {
		client := testsuit.CreateTestClient(8080, nil)
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
		assert.Nil(t, gqlErr)

		get := resultSet.Data["Get"].(map[string]interface{})
		require.Equal(t, 1, len(get))

		pizza := get["Pizza"].([]interface{})
		assert.Equal(t, 1, len(pizza))
	})

	t.Run("Get with creationTimeUnix filters", func(t *testing.T) {
		client := testsuit.CreateTestClient(8080, nil)
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

		require.Nil(t, err)
		require.Nil(t, result.Errors)
		require.NotNil(t, result)
		require.NotNil(t, result.Data)

		b, err := json.Marshal(result.Data)
		require.Nil(t, err)

		var resp GetPizzaResponse
		err = json.Unmarshal(b, &resp)
		require.Nil(t, err)
		require.NotEmpty(t, resp.Get.Pizzas)

		assert.Equal(t, expectedCreateTime, resp.Get.Pizzas[0].Additional.CreationTimeUnix)
	})

	t.Run("Get with lastUpdateTimeUnix filters", func(t *testing.T) {
		client := testsuit.CreateTestClient(8080, nil)
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

		require.Nil(t, err)
		require.Nil(t, result.Errors)
		require.NotNil(t, result)
		require.NotNil(t, result.Data)

		b, err := json.Marshal(result.Data)
		require.Nil(t, err)

		var resp GetPizzaResponse
		err = json.Unmarshal(b, &resp)
		require.Nil(t, err)
		require.NotEmpty(t, resp.Get.Pizzas)

		assert.Equal(t, expectedUpdateTime, resp.Get.Pizzas[0].Additional.LastUpdateTimeUnix)
	})

	t.Run("Get bm25 filter", func(t *testing.T) {
		client := testsuit.CreateTestClient(8080, nil)
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

		require.Nil(t, err)
		require.Nil(t, result.Errors)
		require.NotNil(t, result)
		require.NotNil(t, result.Data)

		b, err := json.Marshal(result.Data)
		require.Nil(t, err)

		var resp GetPizzaResponse
		err = json.Unmarshal(b, &resp)
		require.Nil(t, err)
		require.NotEmpty(t, resp.Get.Pizzas)

		assert.Equal(t, 1, len(resp.Get.Pizzas))
		assert.Equal(t, "A innovation, some say revolution, in the pizza industry.", resp.Get.Pizzas[0].Description)
	})

	t.Run("Get hybrid filter", func(t *testing.T) {
		client := testsuit.CreateTestClient(8080, nil)
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

		require.Nil(t, err)
		require.Nil(t, result.Errors)
		require.NotNil(t, result)
		require.NotNil(t, result.Data)

		b, err := json.Marshal(result.Data)
		require.Nil(t, err)

		var resp GetPizzaResponse
		err = json.Unmarshal(b, &resp)
		require.Nil(t, err)
		require.NotEmpty(t, resp.Get.Pizzas)

		assert.Equal(t, 4, len(resp.Get.Pizzas))
	})

	t.Run("MultiClass Get", func(t *testing.T) {
		client := testsuit.CreateTestClient(8080, nil)
		testsuit.CreateTestSchemaAndData(t, client)
		defer testsuit.CleanUpWeaviate(t, client)

		name := graphql.Field{Name: "name"}
		description := graphql.Field{Name: "description"}

		// what what
		t.Run("multiclass get raw", func(t *testing.T) {
			resultSet, gqlErr := client.GraphQL().Raw().WithQuery("{Get {Pizza {name} Risotto{description}}}").Do(context.Background())
			assert.Nil(t, gqlErr)

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
			assert.Nil(t, gqlErr)

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
			assert.Nil(t, gqlErr)
			require.Nil(t, resultSet.Errors)

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
						{ID: "5b6a08ba-1d46-43aa-89cc-8b070790c6f1"},
						{Beacon: "weaviate://localhost/5b6a08ba-1d46-43aa-89cc-8b070790c6f1"},
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

				assert.Nil(t, gqlErr)
				assert.NotNil(t, resultSet)
				assert.Nil(t, resultSet.Errors)
			})

			t.Run("with distance", func(t *testing.T) {
				concepts := []string{"pineapple slices", "ham"}
				moveTo := &graphql.MoveParameters{
					Force: 0.3,
					Objects: []graphql.MoverObject{
						{ID: "5b6a08ba-1d46-43aa-89cc-8b070790c6f1"},
						{Beacon: "weaviate://localhost/5b6a08ba-1d46-43aa-89cc-8b070790c6f1"},
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

				assert.Nil(t, gqlErr)
				assert.NotNil(t, resultSet)
				assert.Nil(t, resultSet.Errors)
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
			assert.Nil(t, gqlErr)

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
			require.Nil(t, err)

			// two risotto have best_before dates in the test set, one of them expires
			// May 3rd, the other May 5th, so the filter below should match both
			risottoTargetDate, err := time.Parse(time.RFC3339, "2022-05-02T12:00:00+02:00")
			require.Nil(t, err)

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
			assert.Nil(t, gqlErr)

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
			assert.Nil(t, gqlErr)
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
			assert.Nil(t, gqlErr)
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
			assert.Nil(t, gqlErr)
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

			assert.Nil(t, gqlErr)
			assertSortResult(resultSet, "Pizza", []string{"Quattro Formaggi", "Frutti di Mare", "Hawaii", "Doener"})
			assertSortResult(resultSet, "Risotto", []string{"Risi e bisi", "Risotto alla pilota", "Risotto al nero di seppia"})
		})
	})

	t.Run("MultiClassGet with group filter", func(t *testing.T) {
		client := testsuit.CreateTestClient(8080, nil)
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
		assert.Nil(t, gqlErr)

		get := resultSet.Data["Get"].(map[string]interface{})
		require.Equal(t, 2, len(get))

		pizza := get["Pizza"].([]interface{})
		assert.Equal(t, 1, len(pizza))
		risotto := get["Pizza"].([]interface{})
		assert.Equal(t, 1, len(risotto))
	})

	t.Run("MultiClassGet with creationTimeUnix filters", func(t *testing.T) {
		client := testsuit.CreateTestClient(8080, nil)
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

		require.Nil(t, err)
		require.Nil(t, result.Errors)
		require.NotNil(t, result)
		require.NotNil(t, result.Data)

		b, err := json.Marshal(result.Data)
		require.Nil(t, err)

		var pizzaResp GetPizzaResponse
		err = json.Unmarshal(b, &pizzaResp)
		require.Nil(t, err)
		require.NotEmpty(t, pizzaResp.Get.Pizzas)
		assert.Equal(t, pizzaExpectedCreateTime, pizzaResp.Get.Pizzas[0].Additional.CreationTimeUnix)

		var risottoResp GetRisottoResponse
		err = json.Unmarshal(b, &risottoResp)
		require.Nil(t, err)
		require.NotEmpty(t, risottoResp.Get.Risotto)
		assert.Equal(t, risottoExpectedCreateTime, risottoResp.Get.Risotto[0].Additional.CreationTimeUnix)
	})

	t.Run("MultiClassGet with lastUpdateTimeUnix filters", func(t *testing.T) {
		client := testsuit.CreateTestClient(8080, nil)
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

		require.Nil(t, err)
		require.Nil(t, result.Errors)
		require.NotNil(t, result)
		require.NotNil(t, result.Data)

		b, err := json.Marshal(result.Data)
		require.Nil(t, err)

		var pizzaResp GetPizzaResponse
		err = json.Unmarshal(b, &pizzaResp)
		require.Nil(t, err)
		require.NotEmpty(t, pizzaResp.Get.Pizzas)
		assert.Equal(t, pizzaExpectedUpdateTime, pizzaResp.Get.Pizzas[0].Additional.LastUpdateTimeUnix)

		var risottoResp GetRisottoResponse
		err = json.Unmarshal(b, &risottoResp)
		require.Nil(t, err)
		require.NotEmpty(t, risottoResp.Get.Risotto)
		assert.Equal(t, risottoExpectedUpdateTime, risottoResp.Get.Risotto[0].Additional.LastUpdateTimeUnix)
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to tear down weaviate: %s", err)
		}
	})
}
