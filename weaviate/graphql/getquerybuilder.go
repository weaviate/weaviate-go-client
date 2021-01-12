package graphql

import (
	"context"
	"fmt"

	"github.com/semi-technologies/weaviate/entities/models"
)

// GetBuilder for GraphQL
type GetBuilder struct {
	connection rest
	className  string
	withFields string

	includesFilterClause bool // true if brackets behind class is needed
	withWhereFilter      string
	includesLimit        bool
	limit                int
	withNearTextFilter   string
	withNearVectorFilter string
	withGroupFilter      string
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

// WithNearText clause to find close objects
func (gb *GetBuilder) WithNearText(nearText string) *GetBuilder {
	gb.includesFilterClause = true
	gb.withNearTextFilter = nearText
	return gb
}

// WithNearVector clause to find close objects
func (gb *GetBuilder) WithNearVector(nearVector string) *GetBuilder {
	gb.includesFilterClause = true
	gb.withNearVectorFilter = nearVector
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
	filterClause := ""
	if gb.includesFilterClause {
		filterClause = gb.createFilterClause()
	}

	query := fmt.Sprintf("{Get {%v %v {%v}}}", gb.className, filterClause, gb.withFields)

	return query
}

func (gb *GetBuilder) createFilterClause() string {
	clause := "("
	if len(gb.withWhereFilter) > 0 {
		clause += fmt.Sprintf("where: %v", gb.withWhereFilter)
	}
	if len(gb.withNearTextFilter) > 0 {
		if string(clause[len(clause)-1]) == "(" {
			clause += fmt.Sprintf("nearText: %v", gb.withNearTextFilter)
		} else {
			clause += fmt.Sprintf(", nearText: %v", gb.withNearTextFilter)
		}
	}
	if len(gb.withNearVectorFilter) > 0 {
		if string(clause[len(clause)-1]) == "(" {
			clause += fmt.Sprintf("nearVector: %v", gb.withNearVectorFilter)
		} else {
			clause += fmt.Sprintf(", nearVector: %v", gb.withNearVectorFilter)
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
