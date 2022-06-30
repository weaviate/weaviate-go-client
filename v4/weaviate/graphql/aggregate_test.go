package graphql

import (
	"testing"

	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/filters"
	"github.com/stretchr/testify/assert"
)

func TestAggregateBuilder(t *testing.T) {
	t.Run("Simple Aggregate", func(t *testing.T) {
		conMock := &MockRunREST{}

		builder := AggregateBuilder{
			connection: conMock,
		}

		meta := Field{
			Name:   "meta",
			Fields: []Field{{Name: "count"}},
		}

		query := builder.WithClassName("Pizza").WithFields(meta).build()

		expected := `{Aggregate{Pizza{meta{count}}}}`
		assert.Equal(t, expected, query)
	})

	t.Run("Group by", func(t *testing.T) {
		conMock := &MockRunREST{}
		builder := AggregateBuilder{
			connection: conMock,
		}

		groupedBy := Field{
			Name:   "groupedBy",
			Fields: []Field{{Name: "value"}},
		}
		name := Field{
			Name:   "name",
			Fields: []Field{{Name: "count"}},
		}

		query := builder.WithClassName("Pizza").WithFields(groupedBy, name).WithGroupBy("name").build()

		expected := `{Aggregate{Pizza(groupBy: "name"){groupedBy{value} name{count}}}}`
		assert.Equal(t, expected, query)
	})

	t.Run("Where", func(t *testing.T) {
		conMock := &MockRunREST{}
		builder := AggregateBuilder{
			connection: conMock,
		}

		meta := Field{
			Name:   "meta",
			Fields: []Field{{Name: "count"}},
		}

		where := filters.Where().
			WithPath([]string{"id"}).
			WithOperator(filters.Equal).
			WithValueString("uuid")

		query := builder.WithClassName("Pizza").
			WithWhere(where).
			WithFields(meta).build()

		expected := `{Aggregate{Pizza(where:{operator: Equal path: ["id"] valueString: "uuid"}){meta{count}}}}`
		assert.Equal(t, expected, query)
	})

	t.Run("Where and Group by", func(t *testing.T) {
		conMock := &MockRunREST{}
		builder := AggregateBuilder{
			connection: conMock,
		}

		fields := []Field{
			{
				Name: "meta",
				Fields: []Field{
					{
						Name: "count",
					},
				},
			},
		}

		where := filters.Where().
			WithPath([]string{"id"}).
			WithOperator(filters.Equal).
			WithValueString("uuid")

		query := builder.WithClassName("Pizza").
			WithGroupBy("name").
			WithWhere(where).
			WithFields(fields...).
			build()

		expected := `{Aggregate{Pizza(groupBy: "name", where:{operator: Equal path: ["id"] valueString: "uuid"}){meta{count}}}}`
		assert.Equal(t, expected, query)
	})

	t.Run("nearVector", func(t *testing.T) {
		t.Run("with certainty", func(t *testing.T) {
			conMock := &MockRunREST{}
			builder := AggregateBuilder{
				connection: conMock,
			}

			meta := Field{
				Name:   "meta",
				Fields: []Field{{Name: "count"}},
			}

			withNearVector := &NearVectorArgumentBuilder{}
			withNearVector.WithVector([]float32{1, 2, 3}).WithCertainty(0.9)

			query := builder.WithClassName("Pizza").
				WithNearVector(withNearVector).
				WithFields(meta).
				build()

			expected := `{Aggregate{Pizza(nearVector:{certainty: 0.9 vector: [1,2,3]}){meta{count}}}}`
			assert.Equal(t, expected, query)
		})

		t.Run("with distance", func(t *testing.T) {
			conMock := &MockRunREST{}
			builder := AggregateBuilder{
				connection: conMock,
			}

			meta := Field{
				Name:   "meta",
				Fields: []Field{{Name: "count"}},
			}

			withNearVector := &NearVectorArgumentBuilder{}
			withNearVector.WithVector([]float32{1, 2, 3}).WithDistance(0.1)

			query := builder.WithClassName("Pizza").
				WithNearVector(withNearVector).
				WithFields(meta).
				build()

			expected := `{Aggregate{Pizza(nearVector:{distance: 0.1 vector: [1,2,3]}){meta{count}}}}`
			assert.Equal(t, expected, query)
		})
	})

	t.Run("nearObject", func(t *testing.T) {
		t.Run("with certainty", func(t *testing.T) {
			conMock := &MockRunREST{}
			builder := AggregateBuilder{
				connection: conMock,
			}

			meta := Field{
				Name:   "meta",
				Fields: []Field{{Name: "count"}},
			}

			withNearObject := &NearObjectArgumentBuilder{}
			withNearObject.WithID("123").WithBeacon("weaviate://test").WithCertainty(0.7878)

			query := builder.WithClassName("Pizza").
				WithNearObject(withNearObject).
				WithFields(meta).
				build()

			expected := `{Aggregate{Pizza(nearObject:{id: "123" beacon: "weaviate://test" certainty: 0.7878}){meta{count}}}}`
			assert.Equal(t, expected, query)
		})

		t.Run("with distance", func(t *testing.T) {
			conMock := &MockRunREST{}
			builder := AggregateBuilder{
				connection: conMock,
			}

			meta := Field{
				Name:   "meta",
				Fields: []Field{{Name: "count"}},
			}

			withNearObject := &NearObjectArgumentBuilder{}
			withNearObject.WithID("123").WithBeacon("weaviate://test").WithDistance(0.7878)

			query := builder.WithClassName("Pizza").
				WithNearObject(withNearObject).
				WithFields(meta).
				build()

			expected := `{Aggregate{Pizza(nearObject:{id: "123" beacon: "weaviate://test" distance: 0.7878}){meta{count}}}}`
			assert.Equal(t, expected, query)
		})
	})

	t.Run("nearText", func(t *testing.T) {
		t.Run("with certainty", func(t *testing.T) {
			conMock := &MockRunREST{}
			builder := AggregateBuilder{
				connection: conMock,
			}

			meta := Field{
				Name:   "meta",
				Fields: []Field{{Name: "count"}},
			}

			withNearText := &NearTextArgumentBuilder{}
			withNearText.WithConcepts([]string{"pepperoni"}).WithCertainty(0.987)

			query := builder.WithClassName("Pizza").
				WithNearText(withNearText).
				WithFields(meta).
				build()

			expected := `{Aggregate{Pizza(nearText:{concepts: ["pepperoni"] certainty: 0.987}){meta{count}}}}`
			assert.Equal(t, expected, query)
		})

		t.Run("with distance", func(t *testing.T) {
			conMock := &MockRunREST{}
			builder := AggregateBuilder{
				connection: conMock,
			}

			meta := Field{
				Name:   "meta",
				Fields: []Field{{Name: "count"}},
			}

			withNearText := &NearTextArgumentBuilder{}
			withNearText.WithConcepts([]string{"pepperoni"}).WithDistance(0.987)

			query := builder.WithClassName("Pizza").
				WithNearText(withNearText).
				WithFields(meta).
				build()

			expected := `{Aggregate{Pizza(nearText:{concepts: ["pepperoni"] distance: 0.987}){meta{count}}}}`
			assert.Equal(t, expected, query)
		})
	})

	t.Run("objectLimit", func(t *testing.T) {
		t.Run("with certainty", func(t *testing.T) {
			conMock := &MockRunREST{}
			builder := AggregateBuilder{
				connection: conMock,
			}

			meta := Field{
				Name:   "meta",
				Fields: []Field{{Name: "count"}},
			}

			withNearText := &NearTextArgumentBuilder{}
			withNearText.WithConcepts([]string{"pepperoni"}).WithCertainty(0.987)

			query := builder.WithClassName("Pizza").
				WithFields(meta).
				WithNearText(withNearText).
				WithObjectLimit(3).
				build()

			expected := `{Aggregate{Pizza(nearText:{concepts: ["pepperoni"] certainty: 0.987}, objectLimit: 3){meta{count}}}}`
			assert.Equal(t, expected, query)
		})

		t.Run("with distance", func(t *testing.T) {
			conMock := &MockRunREST{}
			builder := AggregateBuilder{
				connection: conMock,
			}

			meta := Field{
				Name:   "meta",
				Fields: []Field{{Name: "count"}},
			}

			withNearText := &NearTextArgumentBuilder{}
			withNearText.WithConcepts([]string{"pepperoni"}).WithDistance(0.987)

			query := builder.WithClassName("Pizza").
				WithFields(meta).
				WithNearText(withNearText).
				WithObjectLimit(3).
				build()

			expected := `{Aggregate{Pizza(nearText:{concepts: ["pepperoni"] distance: 0.987}, objectLimit: 3){meta{count}}}}`
			assert.Equal(t, expected, query)
		})
	})

	t.Run("objectLimit and limit", func(t *testing.T) {
		t.Run("with certainty", func(t *testing.T) {
			conMock := &MockRunREST{}
			builder := AggregateBuilder{
				connection: conMock,
			}

			meta := Field{
				Name:   "meta",
				Fields: []Field{{Name: "count"}},
			}

			withNearText := &NearTextArgumentBuilder{}
			withNearText.WithConcepts([]string{"pepperoni"}).WithCertainty(0.987)

			query := builder.WithClassName("Pizza").
				WithFields(meta).
				WithNearText(withNearText).
				WithObjectLimit(3).
				WithLimit(10).
				build()

			expected := `{Aggregate{Pizza(nearText:{concepts: ["pepperoni"] certainty: 0.987}, objectLimit: 3, limit: 10){meta{count}}}}`
			assert.Equal(t, expected, query)
		})

		t.Run("with distance", func(t *testing.T) {
			conMock := &MockRunREST{}
			builder := AggregateBuilder{
				connection: conMock,
			}

			meta := Field{
				Name:   "meta",
				Fields: []Field{{Name: "count"}},
			}

			withNearText := &NearTextArgumentBuilder{}
			withNearText.WithConcepts([]string{"pepperoni"}).WithDistance(0.987)

			query := builder.WithClassName("Pizza").
				WithFields(meta).
				WithNearText(withNearText).
				WithObjectLimit(3).
				WithLimit(10).
				build()

			expected := `{Aggregate{Pizza(nearText:{concepts: ["pepperoni"] distance: 0.987}, objectLimit: 3, limit: 10){meta{count}}}}`
			assert.Equal(t, expected, query)
		})
	})

	t.Run("limit", func(t *testing.T) {
		t.Run("with certainty", func(t *testing.T) {
			conMock := &MockRunREST{}
			builder := AggregateBuilder{
				connection: conMock,
			}

			meta := Field{
				Name:   "meta",
				Fields: []Field{{Name: "count"}},
			}

			withNearText := &NearTextArgumentBuilder{}
			withNearText.WithConcepts([]string{"pepperoni"}).WithCertainty(0.987)

			query := builder.WithClassName("Pizza").
				WithFields(meta).
				WithNearText(withNearText).
				WithLimit(10).
				build()

			expected := `{Aggregate{Pizza(nearText:{concepts: ["pepperoni"] certainty: 0.987}, limit: 10){meta{count}}}}`
			assert.Equal(t, expected, query)
		})

		t.Run("with distance", func(t *testing.T) {
			conMock := &MockRunREST{}
			builder := AggregateBuilder{
				connection: conMock,
			}

			meta := Field{
				Name:   "meta",
				Fields: []Field{{Name: "count"}},
			}

			withNearText := &NearTextArgumentBuilder{}
			withNearText.WithConcepts([]string{"pepperoni"}).WithDistance(0.987)

			query := builder.WithClassName("Pizza").
				WithFields(meta).
				WithNearText(withNearText).
				WithLimit(10).
				build()

			expected := `{Aggregate{Pizza(nearText:{concepts: ["pepperoni"] distance: 0.987}, limit: 10){meta{count}}}}`
			assert.Equal(t, expected, query)
		})
	})

	t.Run("nearImage", func(t *testing.T) {
		t.Run("with certainty", func(t *testing.T) {
			conMock := &MockRunREST{}
			builder := AggregateBuilder{
				connection: conMock,
			}

			meta := Field{
				Name:   "meta",
				Fields: []Field{{Name: "count"}},
			}

			withNearImage := &NearImageArgumentBuilder{}
			withNearImage.WithImage("iVBORw0KGgoAAAANS").WithCertainty(0.9)

			query := builder.WithClassName("Pizza").
				WithNearImage(withNearImage).
				WithFields(meta).
				build()

			expected := `{Aggregate{Pizza(nearImage:{image: "iVBORw0KGgoAAAANS" certainty: 0.9}){meta{count}}}}`
			assert.Equal(t, expected, query)
		})

		t.Run("with distance", func(t *testing.T) {
			conMock := &MockRunREST{}
			builder := AggregateBuilder{
				connection: conMock,
			}

			meta := Field{
				Name:   "meta",
				Fields: []Field{{Name: "count"}},
			}

			withNearImage := &NearImageArgumentBuilder{}
			withNearImage.WithImage("iVBORw0KGgoAAAANS").WithDistance(0.9)

			query := builder.WithClassName("Pizza").
				WithNearImage(withNearImage).
				WithFields(meta).
				build()

			expected := `{Aggregate{Pizza(nearImage:{image: "iVBORw0KGgoAAAANS" distance: 0.9}){meta{count}}}}`
			assert.Equal(t, expected, query)
		})
	})

	t.Run("ask", func(t *testing.T) {
		t.Run("with certainty", func(t *testing.T) {
			conMock := &MockRunREST{}
			builder := AggregateBuilder{
				connection: conMock,
			}

			meta := Field{
				Name:   "meta",
				Fields: []Field{{Name: "count"}},
			}

			withAsk := &AskArgumentBuilder{}
			withAsk.WithQuestion("question?").WithAutocorrect(true).WithCertainty(0.5).WithProperties([]string{"property"})

			query := builder.WithClassName("Pizza").
				WithAsk(withAsk).
				WithFields(meta).
				build()

			expected := `{Aggregate{Pizza(ask:{question: "question?" properties: ["property"] certainty: 0.5 autocorrect: true}){meta{count}}}}`
			assert.Equal(t, expected, query)
		})

		t.Run("with distance", func(t *testing.T) {
			conMock := &MockRunREST{}
			builder := AggregateBuilder{
				connection: conMock,
			}

			meta := Field{
				Name:   "meta",
				Fields: []Field{{Name: "count"}},
			}

			withAsk := &AskArgumentBuilder{}
			withAsk.WithQuestion("question?").WithAutocorrect(true).WithDistance(0.5).WithProperties([]string{"property"})

			query := builder.WithClassName("Pizza").
				WithAsk(withAsk).
				WithFields(meta).
				build()

			expected := `{Aggregate{Pizza(ask:{question: "question?" properties: ["property"] distance: 0.5 autocorrect: true}){meta{count}}}}`
			assert.Equal(t, expected, query)
		})
	})

	t.Run("Missuse", func(t *testing.T) {
		conMock := &MockRunREST{}

		t.Run("empty query builder", func(t *testing.T) {
			builder := AggregateBuilder{connection: conMock}
			query := builder.build()
			assert.NotEmpty(t, query, "Check that there is no panic if query is not validly build")
		})
	})
}
