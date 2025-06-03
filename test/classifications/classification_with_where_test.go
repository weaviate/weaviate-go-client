package classifications

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/classifications"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/testenv"
	"github.com/weaviate/weaviate/usecases/classification"
)

func TestClassifications_With_Where_Filters_integration(t *testing.T) {
	t.Run("up", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			fmt.Print(err.Error())
			t.Fail()
		}
	})

	t.Run("POST /classifications", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		createClassificationClasses(t, client)

		sourceWhere := filters.Where().
			WithOperator(filters.Like).
			WithPath([]string{"id"}).
			WithValueString("*")
		classifyProperties := []string{"tagged"}
		basedOnProperties := []string{"description"}
		classification, err := client.Classifications().Scheduler().
			WithType(classifications.Contextual).
			WithClassName("Pizza").
			WithSourceWhereFilter(sourceWhere).
			WithClassifyProperties(classifyProperties).
			WithBasedOnProperties(basedOnProperties).
			Do(context.Background())
		assert.Nil(t, err)
		assert.Contains(t, classification.BasedOnProperties, "description")
		assert.Contains(t, classification.ClassifyProperties, "tagged")
		classification, err = client.Classifications().Scheduler().
			WithType(classifications.Contextual).
			WithClassName("Pizza").
			WithSourceWhereFilter(sourceWhere).
			WithClassifyProperties(classifyProperties).
			WithBasedOnProperties(basedOnProperties).
			WithWaitForCompletion().
			Do(context.Background())
		assert.Nil(t, err)
		assert.Contains(t, classification.BasedOnProperties, "description")
		assert.Contains(t, classification.ClassifyProperties, "tagged")

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("GET /classifications/{id}", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		createClassificationClasses(t, client)

		sourceWhere := filters.Where().
			WithOperator(filters.Like).
			WithPath([]string{"id"}).
			WithValueString("*")
		var k int32 = 3
		c, err := client.Classifications().Scheduler().
			WithType(classifications.KNN).
			WithSettings(&classification.ParamsKNN{K: &k}).
			WithSourceWhereFilter(sourceWhere).
			WithClassName("Pizza").
			WithClassifyProperties([]string{"tagged"}).
			WithBasedOnProperties([]string{"description"}).
			Do(context.Background())
		assert.Nil(t, err)

		getC, getErr := client.Classifications().Getter().WithID(string(c.ID)).Do(context.Background())
		assert.Nil(t, getErr)
		assert.Equal(t, c.ID, getC.ID)
		knn, ok := getC.Settings.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, float64(3), knn["k"])

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			fmt.Print(err.Error())
			t.Fail()
		}
	})
}
