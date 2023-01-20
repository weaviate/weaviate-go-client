package graphql

import (
	"context"
	"fmt"
	"strings"

	"github.com/weaviate/weaviate/entities/models"
)

type MultiClassBuilder struct {
	connection    rest
	classBuilders map[string]*GetBuilder
}

// ClassName that should be queried
func NewQueryClassBuilder(className string) *GetBuilder {
	return &GetBuilder{className: className}
}

func (mb *MultiClassBuilder) AddQueryClass(class *GetBuilder) *MultiClassBuilder {
	mb.classBuilders[class.className] = class
	return mb
}

// Do execute the GraphQL query
func (mb *MultiClassBuilder) Do(ctx context.Context) (*models.GraphQLResponse, error) {
	return runGraphQLQuery(ctx, mb.connection, mb.build())
}

// build the GraphQL query string (not needed when Do is executed)
func (mb *MultiClassBuilder) build() string {
	var query string
	for _, class := range mb.classBuilders {
		filterClause := ""
		if class.includesFilterClause {
			filterClause = class.createFilterClause()
		}
		fieldsClause := class.createFieldsClause()
		query += fmt.Sprintf("%v %v {%v}", class.className, filterClause, fieldsClause) + " "
	}
	query = strings.TrimSpace(query)
	query = fmt.Sprintf("{Get {%v}}", query)
	return query
}
