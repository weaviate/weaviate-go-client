package graphql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAggregateBuilder(t *testing.T) {
	t.Run("Simple Explore", func(t *testing.T) {
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
		query := builder.WithClassName("Pizza").WithFields(fields).build()

		expected := `{Aggregate{Pizza{meta{count}}}}`
		assert.Equal(t, expected, query)
	})

	t.Run("Group by", func(t *testing.T) {
		conMock := &MockRunREST{}
		builder := AggregateBuilder{
			connection: conMock,
		}

		fields := []Field{
			{
				Name: "groupedBy",
				Fields: []Field{
					{
						Name: "value",
					},
				},
			},
			{
				Name: "name",
				Fields: []Field{
					{
						Name: "count",
					},
				},
			},
		}

		query := builder.WithClassName("Pizza").WithFields(fields).WithGroupBy("name").build()

		expected := `{Aggregate{Pizza(groupBy: "name"){groupedBy{value} name{count}}}}`
		assert.Equal(t, expected, query)
	})

	t.Run("Where", func(t *testing.T) {
		conMock := &MockRunREST{}
		builder := AggregateBuilder{
			connection: conMock,
		}

		fields := []Field{
			{
				Name:   "meta",
				Fields: []Field{{Name: "count"}},
			},
		}

		where := newWhereArgBuilder().
			WithPath([]string{"id"}).WithOperator(Equal).WithValueString("uuid")

		query := builder.WithClassName("Pizza").WithWhere(where).WithFields(fields).build()

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

		where := newWhereArgBuilder().
			WithPath([]string{"id"}).WithOperator(Equal).WithValueString("uuid")

		query := builder.WithClassName("Pizza").
			WithGroupBy("name").
			WithWhere(where).
			WithFields(fields).
			build()

		expected := `{Aggregate{Pizza(groupBy: "name", where:{operator: Equal path: ["id"] valueString: "uuid"}){meta{count}}}}`
		assert.Equal(t, expected, query)
	})

	t.Run("nearVector", func(t *testing.T) {
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

		withNearVector := &NearVectorArgumentBuilder{}
		withNearVector.WithVector([]float32{1, 2, 3}).WithCertainty(0.9)

		query := builder.WithClassName("Pizza").
			WithNearVector(withNearVector).
			WithFields(fields).
			build()

		expected := `{Aggregate{Pizza(nearVector:{certainty: 0.9 vector: [1,2,3]}){meta{count}}}}`
		assert.Equal(t, expected, query)
	})

	t.Run("nearObject", func(t *testing.T) {
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

		withNearObject := &NearObjectArgumentBuilder{}
		withNearObject.WithID("123").WithBeacon("weaviate://test").WithCertainty(0.7878)

		query := builder.WithClassName("Pizza").
			WithNearObject(withNearObject).
			WithFields(fields).
			build()

		expected := `{Aggregate{Pizza(nearObject:{id: "123" beacon: "weaviate://test" certainty: 0.7878}){meta{count}}}}`
		assert.Equal(t, expected, query)
	})

	t.Run("nearText", func(t *testing.T) {
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

		withNearText := &NearTextArgumentBuilder{}
		withNearText.WithConcepts([]string{"pepperoni"}).WithCertainty(0.987)

		query := builder.WithClassName("Pizza").
			WithNearText(withNearText).
			WithFields(fields).
			build()

		expected := `{Aggregate{Pizza(nearText:{concepts: ["pepperoni"] certainty: 0.987}){meta{count}}}}`
		assert.Equal(t, expected, query)
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
