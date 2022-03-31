package graphql

import (
	"context"
	"fmt"
	"strings"

	"github.com/semi-technologies/weaviate-go-client/v3/weaviate/connection"
	"github.com/semi-technologies/weaviate/entities/models"
)

// Aggregate allows the building of an aggregation query
type Aggregate struct {
	connection *connection.Connection
}

// Objects aggregate objects
func (a *Aggregate) Objects() *AggregateBuilder {
	return &AggregateBuilder{
		connection: a.connection,
	}
}

// AggregateBuilder for the aggregate GraphQL query string
type AggregateBuilder struct {
	connection                rest
	fields                    string
	className                 string
	includesFilterClause      bool // true if brackets behind class is needed
	groupByClausePropertyName string
	withWhereFilter           *WhereArgumentBuilder
}

// WithFields that should be included in the aggregation query e.g. `meta{count}`
func (ab *AggregateBuilder) WithFields(fields string) *AggregateBuilder {
	ab.fields = fields
	return ab
}

// WithClassName that should be aggregated
func (ab *AggregateBuilder) WithClassName(name string) *AggregateBuilder {
	ab.className = name
	return ab
}

// WithWhere adds the where filter.
func (ab *AggregateBuilder) WithWhere(where *WhereArgumentBuilder) *AggregateBuilder {
	ab.includesFilterClause = true
	ab.withWhereFilter = where
	return ab
}

// WithGroupBy adds the group by property clause as the filter.
//  The group by value/path clause still needs to be set in the WithFields field.
func (ab *AggregateBuilder) WithGroupBy(propertyName string) *AggregateBuilder {
	ab.includesFilterClause = true
	ab.groupByClausePropertyName = propertyName
	return ab
}

// Do execute the aggregation query
func (ab *AggregateBuilder) Do(ctx context.Context) (*models.GraphQLResponse, error) {
	return runGraphQLQuery(ctx, ab.connection, ab.build())
}

func (ab *AggregateBuilder) createFilterClause() string {
	filters := []string{}
	if len(ab.groupByClausePropertyName) > 0 {
		filters = append(filters, fmt.Sprintf(`groupBy: "%v"`, ab.groupByClausePropertyName))
	}
	if ab.withWhereFilter != nil {
		filters = append(filters, ab.withWhereFilter.build())
	}
	return fmt.Sprintf("(%s)", strings.Join(filters, ", "))
}

// build the query string
func (ab *AggregateBuilder) build() string {
	filterClause := ""
	if ab.includesFilterClause {
		filterClause = ab.createFilterClause()
	}
	return fmt.Sprintf(`{Aggregate{%v%v{%v}}}`, ab.className, filterClause, ab.fields)
}
