package graphql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiClassQueryBuilder(t *testing.T) {

	t.Run("Simple Get", func(t *testing.T) {
		conMock := &MockRunREST{}

		builder := MultiClassBuilder{
			connection:    conMock,
			classBuilders: make(map[string]*builderBase),
		}

		name1 := Field{Name: "name1"}
		name2 := Field{Name: "name2"}

		query := builder.
			AddQueryClass(NewQueryClassBuilder("Pizza").WithFields(name1)).
			AddQueryClass(NewQueryClassBuilder("Pasta").WithFields(name2)).
			build()

		expected := "{Get {Pizza  {name1} Pasta  {name2}}}"
		assert.Equal(t, expected, query)
	})

}
