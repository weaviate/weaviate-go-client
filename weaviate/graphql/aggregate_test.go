package graphql

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/filters"
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

	t.Run("tenant", func(t *testing.T) {
		builder := AggregateBuilder{connection: &MockRunREST{}}
		meta := Field{
			Name:   "meta",
			Fields: []Field{{Name: "count"}},
		}

		query := builder.WithClassName("Pizza").
			WithTenant("TenantNo1").
			WithFields(meta).
			build()

		expected := `{Aggregate{Pizza(tenant: "TenantNo1"){meta{count}}}}`
		assert.Equal(t, expected, query)
	})

	t.Run("tenant and where", func(t *testing.T) {
		builder := AggregateBuilder{connection: &MockRunREST{}}
		meta := Field{
			Name:   "meta",
			Fields: []Field{{Name: "count"}},
		}
		where := filters.Where().
			WithPath([]string{"name"}).
			WithOperator(filters.Equal).
			WithValueText("Hawaii")

		query := builder.WithClassName("Pizza").
			WithTenant("TenantNo1").
			WithWhere(where).
			WithFields(meta).
			build()

		expected := `{Aggregate{Pizza(tenant: "TenantNo1", where:{operator: Equal path: ["name"] valueText: "Hawaii"}){meta{count}}}}`
		assert.Equal(t, expected, query)
	})
}

func TestAggregate_NearMedia(t *testing.T) {
	fieldMeta := Field{
		Name:   "meta",
		Fields: []Field{{Name: "count"}},
	}

	t.Run("NearImage", func(t *testing.T) {
		nearImage := (&NearImageArgumentBuilder{}).
			WithImage("iVBORw0KGgoAAAANS").
			WithCertainty(0.5)

		query := (&AggregateBuilder{}).
			WithClassName("PizzaImage").
			WithFields(fieldMeta).
			WithNearImage(nearImage).
			build()

		expected := `{Aggregate{PizzaImage(nearImage:{image: "iVBORw0KGgoAAAANS" certainty: 0.5}){meta{count}}}}`
		assert.Equal(t, expected, query)
	})

	t.Run("NearAudio", func(t *testing.T) {
		nearAudio := (&NearAudioArgumentBuilder{}).
			WithAudio("iVBORw0KGgoAAAANS").
			WithCertainty(0.5)

		query := (&AggregateBuilder{}).
			WithClassName("PizzaAudio").
			WithFields(fieldMeta).
			WithNearAudio(nearAudio).
			build()

		expected := `{Aggregate{PizzaAudio(nearAudio:{audio: "iVBORw0KGgoAAAANS" certainty: 0.5}){meta{count}}}}`
		assert.Equal(t, expected, query)
	})

	t.Run("NearVideo", func(t *testing.T) {
		nearVideo := (&NearVideoArgumentBuilder{}).
			WithVideo("iVBORw0KGgoAAAANS").
			WithCertainty(0.5)

		query := (&AggregateBuilder{}).
			WithClassName("PizzaVideo").
			WithFields(fieldMeta).
			WithNearVideo(nearVideo).
			build()

		expected := `{Aggregate{PizzaVideo(nearVideo:{video: "iVBORw0KGgoAAAANS" certainty: 0.5}){meta{count}}}}`
		assert.Equal(t, expected, query)
	})

	t.Run("NearDepth", func(t *testing.T) {
		nearDepth := (&NearDepthArgumentBuilder{}).
			WithDepth("iVBORw0KGgoAAAANS").
			WithCertainty(0.5)

		query := (&AggregateBuilder{}).
			WithClassName("PizzaDepth").
			WithFields(fieldMeta).
			WithNearDepth(nearDepth).
			build()

		expected := `{Aggregate{PizzaDepth(nearDepth:{depth: "iVBORw0KGgoAAAANS" certainty: 0.5}){meta{count}}}}`
		assert.Equal(t, expected, query)
	})

	t.Run("NearThermal", func(t *testing.T) {
		nearThermal := (&NearThermalArgumentBuilder{}).
			WithThermal("iVBORw0KGgoAAAANS").
			WithCertainty(0.5)

		query := (&AggregateBuilder{}).
			WithClassName("PizzaThermal").
			WithFields(fieldMeta).
			WithNearThermal(nearThermal).
			build()

		expected := `{Aggregate{PizzaThermal(nearThermal:{thermal: "iVBORw0KGgoAAAANS" certainty: 0.5}){meta{count}}}}`
		assert.Equal(t, expected, query)
	})

	t.Run("NearImu", func(t *testing.T) {
		nearImu := (&NearImuArgumentBuilder{}).
			WithImu("iVBORw0KGgoAAAANS").
			WithCertainty(0.5)

		query := (&AggregateBuilder{}).
			WithClassName("PizzaImu").
			WithFields(fieldMeta).
			WithNearImu(nearImu).
			build()

		expected := `{Aggregate{PizzaImu(nearIMU:{imu: "iVBORw0KGgoAAAANS" certainty: 0.5}){meta{count}}}}`
		assert.Equal(t, expected, query)
	})

	t.Run("Hybrid", func(t *testing.T) {
		hybrid := (&HybridArgumentBuilder{}).
			WithQuery("query").WithAlpha(0.5).WithFusionType(Ranked)

		query := (&AggregateBuilder{}).
			WithClassName("PizzaDepth").
			WithFields(fieldMeta).
			WithHybrid(hybrid).
			build()

		expected := `{Aggregate{PizzaDepth(hybrid:{query: "query", alpha: 0.5, fusionType: rankedFusion}){meta{count}}}}`
		assert.Equal(t, expected, query)

		hybrid = (&HybridArgumentBuilder{}).
			WithQuery("new query").WithFusionType(RelativeScore)

		query = (&AggregateBuilder{}).
			WithClassName("PizzaDepth").
			WithFields(fieldMeta).
			WithHybrid(hybrid).
			build()

		expected = `{Aggregate{PizzaDepth(hybrid:{query: "new query", fusionType: relativeScoreFusion}){meta{count}}}}`
		assert.Equal(t, expected, query)
	})
}
