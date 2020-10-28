package graphql

import (
	"context"
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/test/testsuit"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/paragons"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/testenv"
	"github.com/stretchr/testify/assert"
	"testing"
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

		resultSet, gqlErr := client.GraphQL.Get().Things().WithClassName("Pizza").WithFields("name").Do(context.Background())
		assert.Nil(t, gqlErr)

		get := resultSet.Data["Get"].(map[string]interface{})
		things := get["Things"].(map[string]interface{})
		pizza := things["Pizza"].([]interface{})
		assert.Equal(t, 4, len(pizza))


		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("Explore", func(t *testing.T) {
		client := testsuit.CreateTestClient()
		testsuit.CreateTestSchemaAndData(t, client)

		fields := []paragons.ExploreFields{paragons.Certainty, paragons.Beacon, paragons.ClassName}
		concepts := []string{"pineapple slices", "ham"}
		moveTo := &paragons.MoveParameters{
			Concepts: []string{"Pizza"},
			Force:    0.3,
		}
		moveAwayFrom := &paragons.MoveParameters{
			Concepts: []string{"toast", "bread"},
			Force:    0.4,
		}

		resultSet, gqlErr := client.GraphQL.Explore().WithFields(fields).WithConcepts(concepts).WithLimit(3).WithCertainty(0.71).WithMoveTo(moveTo).WithMoveAwayFrom(moveAwayFrom).Do(context.Background())
		assert.Nil(t, gqlErr)
		assert.NotNil(t, resultSet)
		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("Aggregate", func(t *testing.T) {
		client := testsuit.CreateTestClient()
		testsuit.CreateTestSchemaAndData(t, client)

		fields := "meta {count}"
		resultSet, gqlErr := client.GraphQL.Aggregate().Things().WithFields(fields).WithClassName("Pizza").Do(context.Background())

		assert.Nil(t, gqlErr)
		assert.NotNil(t, resultSet)
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

