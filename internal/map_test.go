package internal_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weaviate/weaviate-go-client/v6/internal"
)

func TestMakeMap(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Nil(t, internal.MakeMap[string, any](0))
	})

	t.Run("not empty", func(t *testing.T) {
		assert.NotNil(t, internal.MakeMap[string, any](92))
	})
}
