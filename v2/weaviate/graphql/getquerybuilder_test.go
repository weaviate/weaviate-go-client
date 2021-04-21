package graphql

import (
	"context"
	"testing"

	"github.com/semi-technologies/weaviate-go-client/v2/weaviate/connection"
	"github.com/stretchr/testify/assert"
)

// MockRunREST for testing
type MockRunREST struct {
	ArgPath            string
	ArgRestMethod      string
	ArgRequestBody     interface{}
	ReturnResponseData *connection.ResponseData
	ReturnError        error
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
		}

		query := builder.WithClassName("Pizza").WithFields("name").build()

		expected := "{Get {Pizza  {name}}}"
		assert.Equal(t, expected, query)
	})

	t.Run("Multiple fields", func(t *testing.T) {
		conMock := &MockRunREST{}

		builder := GetBuilder{
			connection: conMock,
		}

		query := builder.WithClassName("Pizza").WithFields("name description").build()

		expected := "{Get {Pizza  {name description}}}"
		assert.Equal(t, expected, query)
	})

	t.Run("Where filter", func(t *testing.T) {
		conMock := &MockRunREST{}

		builder := GetBuilder{
			connection: conMock,
		}

		query := builder.WithClassName("Pizza").
			WithFields("name").
			WithWhere(`{path: ["name"] operator: Equal valueString: "Hawaii" }`).
			build()

		expected := `{Get {Pizza (where: {path: ["name"] operator: Equal valueString: "Hawaii" }) {name}}}`
		assert.Equal(t, expected, query)

		query = builder.WithClassName("Pizza").
			WithFields("name").
			WithWhere(`{operator: Or operands: [{path: ["name"] operator: Equal valueString: "Hawaii"}, {path: ["name"] operator: Equal valueString: "Doener"}]}`).
			build()

		expected = `{Get {Pizza (where: {operator: Or operands: [{path: ["name"] operator: Equal valueString: "Hawaii"}, {path: ["name"] operator: Equal valueString: "Doener"}]}) {name}}}`
		assert.Equal(t, expected, query)
	})

	t.Run("Limit Get", func(t *testing.T) {
		conMock := &MockRunREST{}

		builder := GetBuilder{
			connection:           conMock,
			includesFilterClause: false,
		}

		query := builder.WithClassName("Pizza").WithFields("name").WithLimit(2).build()

		expected := "{Get {Pizza (limit: 2) {name}}}"
		assert.Equal(t, expected, query)
	})

	t.Run("NearText filter", func(t *testing.T) {
		conMock := &MockRunREST{}

		builder := GetBuilder{
			connection:           conMock,
			includesFilterClause: false,
		}

		nearText := &NearTextArgumentBuilder{}
		nearText = nearText.WithConcepts([]string{"good"})
		query := builder.WithClassName("Pizza").WithFields("name").WithNearText(nearText).build()

		expected := `{Get {Pizza (nearText:{concepts: ["good"]}) {name}}}`
		assert.Equal(t, expected, query)
	})

	t.Run("NearVector filter", func(t *testing.T) {
		conMock := &MockRunREST{}

		builder := GetBuilder{
			connection:           conMock,
			includesFilterClause: false,
		}

		query := builder.WithClassName("Pizza").WithFields("name").WithNearVector("{vector: [0, 1, 0.8]}").build()

		expected := `{Get {Pizza (nearVector: {vector: [0, 1, 0.8]}) {name}}}`
		assert.Equal(t, expected, query)
	})

	t.Run("Group filter", func(t *testing.T) {
		conMock := &MockRunREST{}

		builder := GetBuilder{
			connection:           conMock,
			includesFilterClause: false,
		}

		query := builder.WithClassName("Pizza").WithFields("name").WithGroup(`{type: closest force: 0.4}`).build()

		expected := `{Get {Pizza (group: {type: closest force: 0.4}) {name}}}`
		assert.Equal(t, expected, query)
	})

	t.Run("Multiple filter", func(t *testing.T) {
		conMock := &MockRunREST{}

		builder := GetBuilder{
			connection:           conMock,
			includesFilterClause: false,
		}

		nearText := &NearTextArgumentBuilder{}
		nearText = nearText.WithConcepts([]string{"good"})

		query := builder.WithClassName("Pizza").
			WithFields("name").
			WithNearText(nearText).
			WithLimit(2).
			WithWhere(`{path: ["name"] operator: Equal valueString: "Hawaii"}`).
			build()

		expected := `{Get {Pizza (where: {path: ["name"] operator: Equal valueString: "Hawaii"}, nearText:{concepts: ["good"]}, limit: 2) {name}}}`
		assert.Equal(t, expected, query)
	})

	t.Run("Multiple filters", func(t *testing.T) {
		conMock := &MockRunREST{}

		builder := GetBuilder{
			connection:           conMock,
			includesFilterClause: false,
		}

		nearText := &NearTextArgumentBuilder{}
		nearText = nearText.WithConcepts([]string{"good"})

		query := builder.WithClassName("Pizza").
			WithFields("name").
			WithNearText(nearText).
			WithNearVector("{vector: [0, 1, 0.8]}").
			WithLimit(2).
			WithWhere(`{path: ["name"] operator: Equal valueString: "Hawaii"}`).
			build()

		expected := `{Get {Pizza (where: {path: ["name"] operator: Equal valueString: "Hawaii"}, nearText:{concepts: ["good"]}, nearVector: {vector: [0, 1, 0.8]}, limit: 2) {name}}}`
		assert.Equal(t, expected, query)
	})

	t.Run("Missuse", func(t *testing.T) {
		conMock := &MockRunREST{}
		builder := GetBuilder{
			connection:           conMock,
			includesFilterClause: false,
		}
		query := builder.build()
		assert.NotEmpty(t, query, "Check that there is no panic if query is not validly build")
	})

	t.Run("NearText filter with concepts", func(t *testing.T) {
		conMock := &MockRunREST{}

		builder := GetBuilder{
			connection:           conMock,
			includesFilterClause: false,
		}

		nearText := &NearTextArgumentBuilder{}
		nearText = nearText.WithConcepts([]string{"good"})

		query := builder.WithClassName("Pizza").
			WithFields("name").
			WithNearText(nearText).
			build()

		expected := `{Get {Pizza (nearText:{concepts: ["good"]}) {name}}}`
		assert.Equal(t, expected, query)
	})

	t.Run("NearObject filter with all fields", func(t *testing.T) {
		conMock := &MockRunREST{}

		builder := GetBuilder{
			connection:           conMock,
			includesFilterClause: false,
		}

		nearObject := &NearObjectArgumentBuilder{}
		nearObject = nearObject.WithBeacon("weawiate/some-uuid")
		query := builder.WithClassName("Pizza").
			WithFields("name").
			WithNearObject(nearObject).
			build()

		expected := `{Get {Pizza (nearObject:{beacon: "weawiate/some-uuid"}) {name}}}`
		assert.Equal(t, expected, query)

		nearObject = &NearObjectArgumentBuilder{}
		nearObject = nearObject.WithBeacon("weawiate/some-uuid").WithID("some-uuid")
		query = builder.WithClassName("Pizza").
			WithFields("name").
			WithNearObject(nearObject).
			build()

		expected = `{Get {Pizza (nearObject:{id: "some-uuid" beacon: "weawiate/some-uuid"}) {name}}}`
		assert.Equal(t, expected, query)

		nearObject = &NearObjectArgumentBuilder{}
		nearObject = nearObject.WithBeacon("weawiate/some-uuid").WithID("some-uuid").WithCertainty(0.8)
		query = builder.WithClassName("Pizza").
			WithFields("name").
			WithNearObject(nearObject).
			build()

		expected = `{Get {Pizza (nearObject:{id: "some-uuid" beacon: "weawiate/some-uuid" certainty: 0.8}) {name}}}`
		assert.Equal(t, expected, query)

		nearObject = &NearObjectArgumentBuilder{}
		nearObject = nearObject.WithBeacon("weawiate/some-uuid").WithID("some-uuid").WithCertainty(0.8)
		nearText := &NearTextArgumentBuilder{}
		nearText = nearText.WithConcepts([]string{"good"})
		query = builder.WithClassName("Pizza").
			WithFields("name").
			WithNearObject(nearObject).
			WithNearText(nearText).
			build()

		expected = `{Get {Pizza (nearText:{concepts: ["good"]}, nearObject:{id: "some-uuid" beacon: "weawiate/some-uuid" certainty: 0.8}) {name}}}`
		assert.Equal(t, expected, query)
	})
}
