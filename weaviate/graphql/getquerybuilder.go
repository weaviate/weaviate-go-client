package graphql

import (
	"context"
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/weaviate/semantics"
	"github.com/semi-technologies/weaviate/entities/models"
	"strings"
)

// GetBuilder for GraphQL
type GetBuilder struct {
	connection rest
	semanticKind semantics.Kind
	className string
	withFields string

	includesFilterClause bool // true if brackets behind class is needed
	withWhereFilter      string
	includesLimit        bool
	limit                int
	withExploreFilter string
	withGroupFilter string
}

// WithClassName that should be queried
func (gb *GetBuilder) WithClassName(name string) *GetBuilder {
	gb.className = name
	return gb
}

// WithFields included in the result set
func (gb *GetBuilder) WithFields(fields string) *GetBuilder {
	gb.withFields = fields
	return gb
}

// WithWhere filter
func (gb *GetBuilder) WithWhere(filter string) *GetBuilder {
	gb.includesFilterClause = true
	gb.withWhereFilter = filter
	return gb
}

// WithLimit of objects in the result set
func (gb *GetBuilder) WithLimit(limit int) *GetBuilder {
	gb.includesFilterClause = true
	gb.includesLimit = true
	gb.limit = limit
	return gb
}

// WithExplore clause to find close objects
func (gb *GetBuilder) WithExplore(explore string) *GetBuilder {
	gb.includesFilterClause = true
	gb.withExploreFilter = explore
	return gb
}

// WithGroup statement
func (gb *GetBuilder) WithGroup(group string) *GetBuilder {
	gb.includesFilterClause = true
	gb.withGroupFilter = group
	return gb
}

// Do execute the GraphQL query
func (gb *GetBuilder) Do(ctx context.Context) (*models.GraphQLResponse, error) {
	return runGraphQLQuery(ctx, gb.connection, gb.build())
}

// build the GraphQL query string (not needed when Do is executed)
func (gb *GetBuilder) build() string {
	semanticKind := strings.Title(string(gb.semanticKind))

	filterClause := ""
	if gb.includesFilterClause {
		filterClause = gb.createFilterClause()
	}

	query := fmt.Sprintf("{Get {%v {%v %v {%v}}}}", semanticKind, gb.className, filterClause, gb.withFields)

	return query
}

func (gb *GetBuilder) createFilterClause() string {
	clause := "("
	if len(gb.withWhereFilter) > 0 {
		clause += fmt.Sprintf("where: %v", gb.withWhereFilter)
	}
	if len(gb.withExploreFilter) > 0 {
		if string(clause[len(clause)-1]) == "(" {
			clause += fmt.Sprintf("explore: %v", gb.withExploreFilter)
		} else {
			clause += fmt.Sprintf(", explore: %v", gb.withExploreFilter)
		}
	}
	if len(gb.withGroupFilter) > 0 {
		if string(clause[len(clause)-1]) == "(" {
			clause += fmt.Sprintf("group: %v", gb.withGroupFilter)
		} else {
			clause += fmt.Sprintf(", group: %v", gb.withGroupFilter)
		}
	}
	if gb.includesLimit {
		if string(clause[len(clause)-1]) == "(" {
			clause += fmt.Sprintf("limit: %v", gb.limit)
		} else {
			clause += fmt.Sprintf(", limit: %v", gb.limit)
		}
	}
	clause += ")"
	return clause
}

