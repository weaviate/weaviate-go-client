package graphql

import (
	"context"
	"fmt"
	"strings"

	"github.com/semi-technologies/weaviate-go-client/v2/weaviate/connection"
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
	fields                    []Field
	className                 string
	withGroupByClause         bool
	groupByClausePropertyName string
}

// WithFields that should be included in the aggregation query e.g. `meta{count}`
func (ab *AggregateBuilder) WithFields(fields []Field) *AggregateBuilder {
	ab.fields = fields
	return ab
}

// WithClassName that should be aggregated
func (ab *AggregateBuilder) WithClassName(name string) *AggregateBuilder {
	ab.className = name
	return ab
}

// WithGroupBy adds the group by property clause as the filter.
//  The group by value/path clause still needs to be set in the WithFields field.
func (ab *AggregateBuilder) WithGroupBy(propertyName string) *AggregateBuilder {
	ab.withGroupByClause = true
	ab.groupByClausePropertyName = propertyName
	return ab
}

// Do execute the aggregation query
func (ab *AggregateBuilder) Do(ctx context.Context) (*models.GraphQLResponse, error) {
	return runGraphQLQuery(ctx, ab.connection, ab.build())
}

func (gb *AggregateBuilder) createFieldsClause() string {
	if len(gb.fields) > 0 {
		fields := make([]string, len(gb.fields))
		for i := range gb.fields {
			fields[i] = gb.fields[i].build()
		}
		return strings.Join(fields, " ")
	}
	return ""
}

// build the query string
func (ab *AggregateBuilder) build() string {
	filter := ""
	if ab.withGroupByClause {
		filter = fmt.Sprintf(`(groupBy: "%v")`, ab.groupByClausePropertyName)
	}
	fields := ab.createFieldsClause()
	return fmt.Sprintf(`{Aggregate{%v%v{%v}}}`, ab.className, filter, fields)
}
