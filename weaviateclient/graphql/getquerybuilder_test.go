package graphql

import (
	"context"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/paragons"
	"github.com/stretchr/testify/assert"
	"testing"
)

// MockRunREST for testing
type MockRunREST struct {
	ArgPath string
	ArgRestMethod string
	ArgRequestBody interface{}
	ReturnResponseData *connection.ResponseData
	ReturnError error
}

// RunREST store all arguments in mock and return response as defined in mock struct
func (mrr *MockRunREST) RunREST(ctx context.Context, path string, restMethod string, requestBody interface{}) (*connection.ResponseData, error) {
	mrr.ArgPath = path
	mrr.ArgRestMethod = restMethod
	mrr.ArgRequestBody = requestBody
	return mrr.ReturnResponseData, mrr.ReturnError
}


func TestQueryBuilder(t *testing.T) {

	t.Run("Simple Get", func(t *testing.T) {
		conMock := &MockRunREST{}

		builder := GetBuilder{
			connection: conMock,
			semanticKind: paragons.SemanticKindThings,
		}

		query := builder.WithClassName("Pizza").WithFields("name").build()

		expected := "{Get {Things {Pizza  {name}}}}"
		assert.Equal(t, expected, query)
	})

	t.Run("Multiple fields", func(t *testing.T) {
		conMock := &MockRunREST{}

		builder := GetBuilder{
			connection: conMock,
			semanticKind: paragons.SemanticKindThings,
		}

		query := builder.WithClassName("Pizza").WithFields("name description").build()

		expected := "{Get {Things {Pizza  {name description}}}}"
		assert.Equal(t, expected, query)
	})

	t.Run("Where filter", func(t *testing.T) {
		conMock := &MockRunREST{}

		builder := GetBuilder{
			connection: conMock,
			semanticKind: paragons.SemanticKindThings,
		}

		query := builder.WithClassName("Pizza").WithFields("name").WithWhere(`{path: ["name"] operator: Equal valueString: "Hawaii" }`).build()

		expected := `{Get {Things {Pizza (where: {path: ["name"] operator: Equal valueString: "Hawaii" }) {name}}}}`
		assert.Equal(t, expected, query)

		query = builder.WithClassName("Pizza").WithFields("name").WithWhere(`{operator: Or operands: [{path: ["name"] operator: Equal valueString: "Hawaii"}, {path: ["name"] operator: Equal valueString: "Doener"}]}`).build()

		expected = `{Get {Things {Pizza (where: {operator: Or operands: [{path: ["name"] operator: Equal valueString: "Hawaii"}, {path: ["name"] operator: Equal valueString: "Doener"}]}) {name}}}}`
		assert.Equal(t, expected, query)
	})

	t.Run("Limit Get", func(t *testing.T) {
		conMock := &MockRunREST{}

		builder := GetBuilder{
			connection:           conMock,
			semanticKind:         paragons.SemanticKindThings,
			includesFilterClause: false,
		}

		query := builder.WithClassName("Pizza").WithFields("name").WithLimit(2).build()

		expected := "{Get {Things {Pizza (limit: 2) {name}}}}"
		assert.Equal(t, expected, query)
	})

	t.Run("Explor filter", func(t *testing.T) {
		conMock := &MockRunREST{}

		builder := GetBuilder{
			connection:           conMock,
			semanticKind:         paragons.SemanticKindThings,
			includesFilterClause: false,
		}

		query := builder.WithClassName("Pizza").WithFields("name").WithExplore(`{concepts: "good"}`).build()

		expected := `{Get {Things {Pizza (explore: {concepts: "good"}) {name}}}}`
		assert.Equal(t, expected, query)
	})

	t.Run("Group filter", func(t *testing.T) {
		conMock := &MockRunREST{}

		builder := GetBuilder{
			connection:           conMock,
			semanticKind:         paragons.SemanticKindThings,
			includesFilterClause: false,
		}

		query := builder.WithClassName("Pizza").WithFields("name").WithGroup(`{type: closest force: 0.4}`).build()

		expected := `{Get {Things {Pizza (group: {type: closest force: 0.4}) {name}}}}`
		assert.Equal(t, expected, query)
	})

	t.Run("Multiple filter", func(t *testing.T) {
		conMock := &MockRunREST{}

		builder := GetBuilder{
			connection:           conMock,
			semanticKind:         paragons.SemanticKindThings,
			includesFilterClause: false,
		}

		query := builder.WithClassName("Pizza").WithFields("name").WithExplore(`{concepts: "good"}`).WithLimit(2).WithWhere(`{path: ["name"] operator: Equal valueString: "Hawaii"}`).build()

		expected := `{Get {Things {Pizza (where: {path: ["name"] operator: Equal valueString: "Hawaii"}, explore: {concepts: "good"}, limit: 2) {name}}}}`
		assert.Equal(t, expected, query)
	})

	t.Run("Missuse", func(t *testing.T) {
		conMock := &MockRunREST{}
		builder := GetBuilder{
			connection:           conMock,
			semanticKind:         paragons.SemanticKindThings,
			includesFilterClause: false,
		}
		query := builder.build()
		assert.NotEmpty(t, query, "Check that there is no panic if query is not validly build")
	})

}