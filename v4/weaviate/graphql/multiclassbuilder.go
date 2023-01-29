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
	// sorting classBuilder based on className to have consistent order in query
	s := make(classBuildersSorter, 0, len(mb.classBuilders))
	for _, k := range mb.classBuilders {
		s = append(s, k)
	}
	sort.Sort(classBuildersSorter(s))
	for _, class := range s {
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

// Implementing Sort Interface for *GetBuilder Struct
type classBuildersSorter []*GetBuilder

func (mb classBuildersSorter) Len() int {
	return len(mb)
}

func (mb classBuildersSorter) Less(i, j int) bool {
	return mb[i].className < mb[j].className
}

func (mb classBuildersSorter) Swap(i, j int) {
	mb[i], mb[j] = mb[j], mb[i]
}
