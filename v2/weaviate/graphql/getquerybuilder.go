package graphql

import (
	"context"
	"fmt"
	"strings"

	"github.com/semi-technologies/weaviate/entities/models"
)

// GetBuilder for GraphQL
type GetBuilder struct {
	connection rest
	className  string
	withFields string

	includesFilterClause bool // true if brackets behind class is needed
	includesLimit        bool
	limit                int
	withWhereFilter      *WhereArgumentBuilder
	withNearTextFilter   *NearTextArgumentBuilder
	withNearVectorFilter string
	withNearObjectFilter *NearObjectArgumentBuilder
	withGroupFilter      string
	withAskFilter        *AskArgumentBuilder
	withNearImageFilter  *NearImageArgumentBuilder
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
func (gb *GetBuilder) WithWhere(where *WhereArgumentBuilder) *GetBuilder {
	gb.includesFilterClause = true
	gb.withWhereFilter = where
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
func (gb *GetBuilder) WithNearText(nearText *NearTextArgumentBuilder) *GetBuilder {
	gb.includesFilterClause = true
	gb.withNearTextFilter = nearText
	return gb
}

// WithNearObject clause to find close objects
func (gb *GetBuilder) WithNearImage(nearImage *NearImageArgumentBuilder) *GetBuilder {
	gb.includesFilterClause = true
	gb.withNearImageFilter = nearImage
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

// WithAsk clause to find an aswer to the question
func (gb *GetBuilder) WithAsk(ask *AskArgumentBuilder) *GetBuilder {
	gb.includesFilterClause = true
	gb.withAskFilter = ask
	return gb
}

// WithNearObject clause to find close objects
func (gb *GetBuilder) WithNearObject(nearObject *NearObjectArgumentBuilder) *GetBuilder {
	gb.includesFilterClause = true
	gb.withNearObjectFilter = nearObject
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
	filters := []string{}
	if gb.withWhereFilter != nil {
		filters = append(filters, gb.withWhereFilter.build())
	}
	if gb.withNearTextFilter != nil {
		filters = append(filters, gb.withNearTextFilter.build())
	}
	if len(gb.withNearVectorFilter) > 0 {
		filters = append(filters, fmt.Sprintf("nearVector: %v", gb.withNearVectorFilter))
	}
	if gb.withNearObjectFilter != nil {
		filters = append(filters, gb.withNearObjectFilter.build())
	}
	if gb.withAskFilter != nil {
		filters = append(filters, gb.withAskFilter.build())
	}
	if gb.withNearImageFilter != nil {
		filters = append(filters, gb.withNearImageFilter.build())
	}
	if len(gb.withGroupFilter) > 0 {
		filters = append(filters, fmt.Sprintf("group: %v", gb.withGroupFilter))
	}
	if gb.includesLimit {
		filters = append(filters, fmt.Sprintf("limit: %v", gb.limit))
	}
	return fmt.Sprintf("(%s)", strings.Join(filters, ", "))
}
