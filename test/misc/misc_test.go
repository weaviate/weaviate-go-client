package misc

import (
	"context"
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/test/testsuit"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/testenv"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMisc_integration(t *testing.T) {

	t.Run("up", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			fmt.Printf(err.Error())
			t.Fail()
		}
	})

	t.Run("GET /.well-known/ready", func(t *testing.T) {

		client := testsuit.CreateTestClient()
		isReady, err := client.Misc.ReadyChecker().Do(context.Background())

		assert.Nil(t, err)
		assert.True(t, isReady)
	})

	t.Run("GET /.well-known/live", func(t *testing.T) {

		client := testsuit.CreateTestClient()
		isLive, err := client.Misc.LiveChecker().Do(context.Background())

		assert.Nil(t, err)
		assert.True(t, isLive)
	})

	t.Run("GET /.well-known/openid-configuration", func(t *testing.T) {

		client := testsuit.CreateTestClient()
		openIDconfig, err := client.Misc.OpenIDConfigurationGetter().Do(context.Background())

		assert.Nil(t, err)
		assert.Nil(t, openIDconfig)
	})

	t.Run("GET /meta", func(t *testing.T) {
		client := testsuit.CreateTestClient()
		meta, err := client.Misc.MetaGetter().Do(context.Background())
		assert.Nil(t, err)
		assert.NotEmpty(t, meta.Version)
		assert.NotEmpty(t, meta.ContextionaryVersion)
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			fmt.Printf(err.Error())
			t.Fail()
		}
	})

}

func TestMisc_connection_error(t *testing.T) {
	t.Run("ready", func(t *testing.T) {
		cfg := weaviateclient.Config{
			Host:   "localhorst",
			Scheme: "http",
		}

		client := weaviateclient.New(cfg)
		isReady, err := client.Misc.ReadyChecker().Do(context.Background())

		assert.NotNil(t, err)
		assert.False(t, isReady)
	})

	t.Run("live", func(t *testing.T) {
		cfg := weaviateclient.Config{
			Host:   "localhorst",
			Scheme: "http",
		}

		client := weaviateclient.New(cfg)
		isReady, err := client.Misc.LiveChecker().Do(context.Background())

		assert.NotNil(t, err)
		assert.False(t, isReady)
	})
}
