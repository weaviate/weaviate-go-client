package batch

import (
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/testenv"
	"testing"
)

func TestBatch_integration(t *testing.T) {

	t.Run("up", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			fmt.Printf(err.Error())
			t.Fail()
		}
	})

	t.Run("POST /batching/{type}", func(t *testing.T) {




	})

	t.Run("POST /batching/references", func(t *testing.T) {

	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			fmt.Printf(err.Error())
			t.Fail()
		}
	})
}