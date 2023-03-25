package misc

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/weaviate/weaviate-go-client/v4/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/testenv"
)

func TestMisc_version_check(t *testing.T) {
	cfg := &weaviate.Config{
		Host:    "localhost:8080",
		Scheme:  "http",
		Headers: map[string]string{},
	}

	client := weaviate.New(*cfg)

	t.Run("Weaviate is not live, perform live check", func(t *testing.T) {
		isReady, err := client.Misc().ReadyChecker().Do(context.Background())
		assert.NotNil(t, err)
		assert.False(t, isReady)
	})

	t.Run("Start Weaviate", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			fmt.Printf(err.Error())
			t.Fail()
		}
	})
	client.WaitForWeavaite(time.Second * 60)
	t.Run("Weaviate is live, perform ready check", func(t *testing.T) {
		isReady, err := client.Misc().ReadyChecker().Do(context.Background())
		assert.Nil(t, err)
		assert.True(t, isReady)
	})

	t.Run("Create sample schema food, try to perform queries using /v1/objects?class={className}", func(t *testing.T) {
		testsuit.CreateWeaviateTestSchemaFood(t, client)

		_, errCreate := client.Data().Creator().WithClassName("Pizza").WithProperties(map[string]string{
			"name":        "Pepperoni",
			"description": "meat",
		}).Do(context.Background())
		assert.Nil(t, errCreate)

		_, errCreate = client.Data().Creator().WithClassName("Soup").WithProperties(map[string]string{
			"name":        "Chicken",
			"description": "meat",
		}).Do(context.Background())
		assert.Nil(t, errCreate)

		pizzas, pizzasErr := client.Data().ObjectsGetter().WithClassName("Pizza").Do(context.Background())
		assert.Nil(t, pizzasErr)
		assert.Equal(t, 1, len(pizzas))

		soups, soupsErr := client.Data().ObjectsGetter().WithClassName("Soup").Do(context.Background())
		assert.Nil(t, soupsErr)
		assert.Equal(t, 1, len(soups))

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

func TestMisc_empty_version_check(t *testing.T) {
	// needs to be the same client for the whole test
	cfg := &weaviate.Config{
		Host:    "localhost:8080",
		Scheme:  "http",
		Headers: map[string]string{},
	}

	client := weaviate.New(*cfg)

	t.Run("Weaviate is not live, perform live check", func(t *testing.T) {
		isReady, err := client.Misc().ReadyChecker().Do(context.Background())
		assert.NotNil(t, err)
		assert.False(t, isReady)
	})

	t.Run("Start Weaviate", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			fmt.Printf(err.Error())
			t.Fail()
		}
	})
	client.WaitForWeavaite(60 * time.Second)
	t.Run("Create sample schema food, try to perform queries using /v1/objects?class={className}", func(t *testing.T) {
		testsuit.CreateWeaviateTestSchemaFood(t, client)

		_, errCreate := client.Data().Creator().WithClassName("Pizza").WithProperties(map[string]string{
			"name":        "Pepperoni",
			"description": "meat",
		}).Do(context.Background())
		assert.Nil(t, errCreate)

		_, errCreate = client.Data().Creator().WithClassName("Soup").WithProperties(map[string]string{
			"name":        "Chicken",
			"description": "meat",
		}).Do(context.Background())
		assert.Nil(t, errCreate)

		pizzas, pizzasErr := client.Data().ObjectsGetter().WithClassName("Pizza").Do(context.Background())
		assert.Nil(t, pizzasErr)
		assert.Equal(t, 1, len(pizzas))

		soups, soupsErr := client.Data().ObjectsGetter().WithClassName("Soup").Do(context.Background())
		assert.Nil(t, soupsErr)
		assert.Equal(t, 1, len(soups))

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
