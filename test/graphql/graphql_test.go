package graphql

import (
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/testenv"
	"testing"
)

func TestClassifications_integration(t *testing.T) {

	t.Run("up", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			fmt.Printf(err.Error())
			t.Fail()
		}
	})

	t.Run("", func(t *testing.T) {

		t.Fail()
	})

	t.Run("", func(t *testing.T) {
		t.Fail()
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
