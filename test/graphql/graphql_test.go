package graphql

import (
	"context"
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/test/testsuit"
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

		//resultSet, gqlErr := client.GraphQL.Explore().WithFields([]paragons.ExploreFields{paragons.Certainty, paragons.Beacon}).WithConcepts([]string{"Apple"}).WithLimit(3).WithCertainty(0.71).WithMoveTo([]string{"pizza"}, 0.2).WithMoveAwayFrom([]string{"Fish"}, 0.1).Do()
		//assert.Nil(t, gqlErr)
		//
		//assert.NotNil(t, resultSet)

		t.Fail()


		testsuit.CleanUpWeaviate(t, client)



		/*
		{
			Explore(
				concepts: "apple",
			limit: 3,
			certainty: 0.71,
			moveTo: {
		concepts: "pizza"
		force: 0.2
		}
		moveAwayFrom: {
		concepts: "fish",
			force: 0.1
		}
			) {
			certainty
			beacon
			className
		}
		}
		*/
	})

	t.Run("", func(t *testing.T) {
		t.Fail()
	})

	t.Run("", func(t *testing.T) {
		t.Fail()
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			fmt.Printf(err.Error())
			t.Fail()
		}
	})
}

