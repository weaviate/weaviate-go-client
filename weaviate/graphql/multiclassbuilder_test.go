package graphql

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/filters"
)

func TestMultiClassQueryBuilder(t *testing.T) {
	t.Run("Simple Get", func(t *testing.T) {
		conMock := &MockRunREST{}

		builder := MultiClassBuilder{
			connection:    conMock,
			classBuilders: make(map[string]*GetBuilder),
		}

		name1 := Field{Name: "name1"}
		name2 := Field{Name: "name2"}

		query := builder.
			AddQueryClass(NewQueryClassBuilder("Pizza").WithFields(name1)).
			AddQueryClass(NewQueryClassBuilder("Risotto").WithFields(name2)).
			build()

		expected := "{Get {Pizza  {name1} Risotto  {name2}}}"
		assert.Equal(t, expected, query)
	})

	t.Run("Multiple fields", func(t *testing.T) {
		conMock := &MockRunREST{}

		builder := MultiClassBuilder{
			connection:    conMock,
			classBuilders: make(map[string]*GetBuilder),
		}

		fields1 := []Field{
			{Name: "name"},
			{Name: "description"},
		}

		fields2 := []Field{
			{Name: "email"},
			{Name: "password"},
		}

		query := builder.
			AddQueryClass(NewQueryClassBuilder("Pizza").WithFields(fields1...)).
			AddQueryClass(NewQueryClassBuilder("Risotto").WithFields(fields2...)).
			build()

		expected := "{Get {Pizza  {name description} Risotto  {email password}}}"
		assert.Equal(t, expected, query)
	})

	t.Run("Where filter", func(t *testing.T) {
		conMock := &MockRunREST{}

		builder := MultiClassBuilder{
			connection:    conMock,
			classBuilders: make(map[string]*GetBuilder),
		}

		name := Field{Name: "name"}
		city := Field{Name: "city"}

		where1 := filters.Where().
			WithPath([]string{"name"}).
			WithOperator(filters.Equal).
			WithValueString("Hawaii")

		where2 := filters.Where().
			WithPath([]string{"city"}).
			WithOperator(filters.Equal).
			WithValueString("New York")

		query := builder.
			AddQueryClass(NewQueryClassBuilder("Pizza").WithFields(name).WithWhere(where1)).
			AddQueryClass(NewQueryClassBuilder("Risotto").WithFields(city).WithWhere(where2)).
			build()

		expected := `{Get {Pizza (where:{operator: Equal path: ["name"] valueString: "Hawaii"}) {name} Risotto (where:{operator: Equal path: ["city"] valueString: "New York"}) {city}}}`
		assert.Equal(t, expected, query)
	})

	t.Run("Limit And Offset Get", func(t *testing.T) {
		conMock := &MockRunREST{}

		builder := MultiClassBuilder{
			connection:    conMock,
			classBuilders: make(map[string]*GetBuilder),
		}

		name := Field{Name: "name"}
		city := Field{Name: "city"}

		query := builder.
			AddQueryClass(NewQueryClassBuilder("Pizza").WithFields(name).WithOffset(0).WithLimit(2)).
			AddQueryClass(NewQueryClassBuilder("Risotto").WithFields(city).WithOffset(5).WithLimit(8)).
			build()

		expected := "{Get {Pizza (limit: 2, offset: 0) {name} Risotto (limit: 8, offset: 5) {city}}}"
		assert.Equal(t, expected, query)
	})

	t.Run("NearVector filter", func(t *testing.T) {
		conMock := &MockRunREST{}

		builder := MultiClassBuilder{
			connection:    conMock,
			classBuilders: make(map[string]*GetBuilder),
		}

		name := Field{Name: "name"}
		city := Field{Name: "city"}

		nearVector1 := &NearVectorArgumentBuilder{}
		nearVector1.WithVector([]float32{0, 1, 0.8})

		nearVector2 := &NearVectorArgumentBuilder{}
		nearVector2.WithVector([]float32{0.3, 0.5, 0.6})

		query := builder.
			AddQueryClass(NewQueryClassBuilder("Pizza").WithFields(name).WithNearVector(nearVector1)).
			AddQueryClass(NewQueryClassBuilder("Risotto").WithFields(city).WithNearVector(nearVector2)).
			build()

		expected := `{Get {Pizza (nearVector:{vector: [0,1,0.8]}) {name} Risotto (nearVector:{vector: [0.3,0.5,0.6]}) {city}}}`
		assert.Equal(t, expected, query)
	})

	t.Run("Group filter", func(t *testing.T) {
		conMock := &MockRunREST{}

		builder := MultiClassBuilder{
			connection:    conMock,
			classBuilders: make(map[string]*GetBuilder),
		}

		name := Field{Name: "name"}
		city := Field{Name: "city"}

		group1 := &GroupArgumentBuilder{}
		group1 = group1.WithType(Closest).WithForce(0.4)

		group2 := &GroupArgumentBuilder{}
		group2 = group2.WithType(Merge).WithForce(0.8)

		query := builder.
			AddQueryClass(NewQueryClassBuilder("Pizza").WithFields(name).WithGroup(group1)).
			AddQueryClass(NewQueryClassBuilder("Risotto").WithFields(city).WithGroup(group2)).
			build()

		expected := `{Get {Pizza (group:{type: closest force: 0.4}) {name} Risotto (group:{type: merge force: 0.8}) {city}}}`
		assert.Equal(t, expected, query)
	})

	t.Run("Multiple filter", func(t *testing.T) {
		conMock := &MockRunREST{}

		builder := MultiClassBuilder{
			connection:    conMock,
			classBuilders: make(map[string]*GetBuilder),
		}

		name := Field{Name: "name"}
		city := Field{Name: "city"}

		nearText := &NearTextArgumentBuilder{}
		nearText = nearText.WithConcepts([]string{"good"})

		where := filters.Where().
			WithPath([]string{"city"}).
			WithOperator(filters.Equal).
			WithValueString("New York")

		query := builder.
			AddQueryClass(NewQueryClassBuilder("Pizza").WithFields(name).WithNearText(nearText)).
			AddQueryClass(NewQueryClassBuilder("Risotto").WithFields(city).WithWhere(where)).
			build()

		expected := `{Get {Pizza (nearText:{concepts: ["good"]}) {name} Risotto (where:{operator: Equal path: ["city"] valueString: "New York"}) {city}}}`
		assert.Equal(t, expected, query)
	})

	t.Run("Multiple filter", func(t *testing.T) {
		conMock := &MockRunREST{}

		builder := MultiClassBuilder{
			connection:    conMock,
			classBuilders: make(map[string]*GetBuilder),
		}

		name := Field{Name: "name"}
		city := Field{Name: "city"}

		nearText1 := &NearTextArgumentBuilder{}
		nearText1 = nearText1.WithConcepts([]string{"good"})

		where1 := filters.Where().
			WithPath([]string{"name"}).
			WithOperator(filters.Equal).
			WithValueString("Hawaii")

		nearText2 := &NearTextArgumentBuilder{}
		nearText2 = nearText2.WithConcepts([]string{"best"})

		where2 := filters.Where().
			WithPath([]string{"city"}).
			WithOperator(filters.Equal).
			WithValueString("New York")

		query := builder.
			AddQueryClass(NewQueryClassBuilder("Pizza").WithFields(name).
				WithNearText(nearText1).
				WithLimit(2).
				WithWhere(where1)).
			AddQueryClass(NewQueryClassBuilder("Risotto").WithFields(city).
				WithNearText(nearText2).
				WithLimit(5).
				WithWhere(where2)).
			build()

		expected := `{Get {Pizza (where:{operator: Equal path: ["name"] valueString: "Hawaii"}, nearText:{concepts: ["good"]}, limit: 2) {name} Risotto (where:{operator: Equal path: ["city"] valueString: "New York"}, nearText:{concepts: ["best"]}, limit: 5) {city}}}`
		assert.Equal(t, expected, query)
	})

	t.Run("NearObject filter with all fields", func(t *testing.T) {
		t.Run("with certainty", func(t *testing.T) {
			conMock := &MockRunREST{}

			builder := MultiClassBuilder{
				connection:    conMock,
				classBuilders: make(map[string]*GetBuilder),
			}

			name := Field{Name: "name"}
			city := Field{Name: "city"}

			nearObject1 := &NearObjectArgumentBuilder{}
			nearObject1 = nearObject1.WithBeacon("uuid1").WithID("uuid1").WithCertainty(0.8)
			nearText1 := &NearTextArgumentBuilder{}
			nearText1 = nearText1.WithConcepts([]string{"good"})

			nearObject2 := &NearObjectArgumentBuilder{}
			nearObject2 = nearObject2.WithBeacon("uuid2").WithID("uuid2").WithCertainty(0.5)
			nearText2 := &NearTextArgumentBuilder{}
			nearText2 = nearText2.WithConcepts([]string{"best"})

			query := builder.
				AddQueryClass(NewQueryClassBuilder("Pizza").
					WithFields(name).
					WithNearObject(nearObject1).
					WithNearText(nearText1)).
				AddQueryClass(NewQueryClassBuilder("Risotto").
					WithFields(city).
					WithNearObject(nearObject2).
					WithNearText(nearText2)).
				build()

			expected := `{Get {Pizza (nearText:{concepts: ["good"]}, nearObject:{id: "uuid1" beacon: "uuid1" certainty: 0.8}) {name} Risotto (nearText:{concepts: ["best"]}, nearObject:{id: "uuid2" beacon: "uuid2" certainty: 0.5}) {city}}}`
			assert.Equal(t, expected, query)
		})

		t.Run("with distance", func(t *testing.T) {
			conMock := &MockRunREST{}

			builder := MultiClassBuilder{
				connection:    conMock,
				classBuilders: make(map[string]*GetBuilder),
			}

			name := Field{Name: "name"}
			city := Field{Name: "city"}

			nearObject1 := &NearObjectArgumentBuilder{}
			nearObject1 = nearObject1.WithBeacon("uuid1").WithID("uuid1").WithDistance(0.2)
			nearText1 := &NearTextArgumentBuilder{}
			nearText1 = nearText1.WithConcepts([]string{"good"})

			nearObject2 := &NearObjectArgumentBuilder{}
			nearObject2 = nearObject2.WithBeacon("uuid2").WithID("uuid2").WithDistance(0.5)
			nearText2 := &NearTextArgumentBuilder{}
			nearText2 = nearText2.WithConcepts([]string{"best"})

			query := builder.
				AddQueryClass(NewQueryClassBuilder("Pizza").
					WithFields(name).
					WithNearObject(nearObject1).
					WithNearText(nearText1)).
				AddQueryClass(NewQueryClassBuilder("Risotto").
					WithFields(city).
					WithNearObject(nearObject2).
					WithNearText(nearText2)).
				build()

			expected := `{Get {Pizza (nearText:{concepts: ["good"]}, nearObject:{id: "uuid1" beacon: "uuid1" distance: 0.2}) {name} Risotto (nearText:{concepts: ["best"]}, nearObject:{id: "uuid2" beacon: "uuid2" distance: 0.5}) {city}}}`
			assert.Equal(t, expected, query)
		})
	})

	t.Run("Ask filter with all fields", func(t *testing.T) {
		t.Run("with certainty", func(t *testing.T) {
			conMock := &MockRunREST{}

			builder := MultiClassBuilder{
				connection:    conMock,
				classBuilders: make(map[string]*GetBuilder),
			}

			name := Field{Name: "name"}
			city := Field{Name: "city"}

			ask1 := &AskArgumentBuilder{}
			ask1 = ask1.WithQuestion("What is Weaviate?").
				WithProperties([]string{"prop1", "prop2"}).
				WithCertainty(0.8).
				WithRerank(true)

			ask2 := &AskArgumentBuilder{}
			ask2 = ask2.WithQuestion("How to use Weaviate?").
				WithProperties([]string{"prop3", "prop4"}).
				WithCertainty(0.5).
				WithRerank(false)

			query := builder.
				AddQueryClass(NewQueryClassBuilder("Pizza").
					WithFields(name).
					WithAsk(ask1)).
				AddQueryClass(NewQueryClassBuilder("Risotto").
					WithFields(city).
					WithAsk(ask2)).
				build()

			expected := `{Get {Pizza (ask:{question: "What is Weaviate?" properties: ["prop1","prop2"] certainty: 0.8 rerank: true}) {name} Risotto (ask:{question: "How to use Weaviate?" properties: ["prop3","prop4"] certainty: 0.5 rerank: false}) {city}}}`
			assert.Equal(t, expected, query)
		})

		t.Run("with distance", func(t *testing.T) {
			conMock := &MockRunREST{}

			builder := MultiClassBuilder{
				connection:    conMock,
				classBuilders: make(map[string]*GetBuilder),
			}

			name := Field{Name: "name"}
			city := Field{Name: "city"}

			ask1 := &AskArgumentBuilder{}
			ask1 = ask1.WithQuestion("What is Weaviate?").
				WithProperties([]string{"prop1", "prop2"}).
				WithDistance(0.2).
				WithRerank(true)

			ask2 := &AskArgumentBuilder{}
			ask2 = ask2.WithQuestion("How to use Weaviate?").
				WithProperties([]string{"prop3", "prop4"}).
				WithDistance(0.5).
				WithRerank(false)

			query := builder.
				AddQueryClass(NewQueryClassBuilder("Pizza").
					WithFields(name).
					WithAsk(ask1)).
				AddQueryClass(NewQueryClassBuilder("Risotto").
					WithFields(city).
					WithAsk(ask2)).
				build()

			expected := `{Get {Pizza (ask:{question: "What is Weaviate?" properties: ["prop1","prop2"] distance: 0.2 rerank: true}) {name} Risotto (ask:{question: "How to use Weaviate?" properties: ["prop3","prop4"] distance: 0.5 rerank: false}) {city}}}`
			assert.Equal(t, expected, query)
		})
	})
}

func TestMultiClassBM25Builder(t *testing.T) {
	conMock := &MockRunREST{}

	builder := MultiClassBuilder{
		connection:    conMock,
		classBuilders: make(map[string]*GetBuilder),
	}

	bm25B_1 := &BM25ArgumentBuilder{}
	bm25B_1 = bm25B_1.WithQuery("good").WithProperties("name", "description")

	bm25B_2 := &BM25ArgumentBuilder{}
	bm25B_2 = bm25B_2.WithQuery("best").WithProperties("email", "password")

	query := builder.
		AddQueryClass(NewQueryClassBuilder("Pizza").WithBM25(bm25B_1)).
		AddQueryClass(NewQueryClassBuilder("Risotto").WithBM25(bm25B_2)).
		build()

	expected := `{Get {Pizza (bm25:{query: "good", properties: ["name","description"]}) {} Risotto (bm25:{query: "best", properties: ["email","password"]}) {}}}`
	assert.Equal(t, expected, query)
}

func TestMultiClassHybridBuilder(t *testing.T) {
	conMock := &MockRunREST{}

	builder := MultiClassBuilder{
		connection:    conMock,
		classBuilders: map[string]*GetBuilder{},
	}

	hybrid_1 := &HybridArgumentBuilder{}
	hybrid_1.WithQuery("query1").WithVector([]float32{1, 2, 3}).WithAlpha(0.6)

	hybrid_2 := &HybridArgumentBuilder{}
	hybrid_2.WithQuery("query2").WithVector([]float32{4, 5, 6}).WithAlpha(0.8)

	query := builder.
		AddQueryClass(NewQueryClassBuilder("Pizza").WithHybrid(hybrid_1)).
		AddQueryClass(NewQueryClassBuilder("Risotto").WithHybrid(hybrid_2)).
		build()

	expected := `{Get {Pizza (hybrid:{query: "query1", vector: [1,2,3], alpha: 0.6}) {} Risotto (hybrid:{query: "query2", vector: [4,5,6], alpha: 0.8}) {}}}`
	assert.Equal(t, expected, query)
}
