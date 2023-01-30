package graphql

import (
	"context"
	"fmt"
	"sort"
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
	// sorting className to have consistent order in query
	s := make([]string, 0, len(mb.classBuilders))
	for key := range mb.classBuilders {
		s = append(s, key)
	}
	sort.Strings(s)
	for _, className := range s {
		filterClause := ""
		if mb.classBuilders[className].includesFilterClause {
			filterClause = mb.classBuilders[className].createFilterClause()
		}
		fieldsClause := mb.classBuilders[className].createFieldsClause()
		query += fmt.Sprintf("%v %v {%v}", className, filterClause, fieldsClause) + " "
	}
	query = strings.TrimSpace(query)
	query = fmt.Sprintf("{Get {%v}}", query)
	return query
}
