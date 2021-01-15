package graphql

import (
	"context"
	"fmt"

	"github.com/semi-technologies/weaviate/entities/models"
)

// Explore query builder
type Explore struct {
	connection   rest
	fields       []ExploreFields
	withNearText *NearTextArgumentBuilder
}

// WithNearText adds nearText to clause
func (e *Explore) WithNearText(nearText *NearTextArgumentBuilder) *Explore {
	e.withNearText = nearText
	return e
}

// WithFields that should be included in the result set
func (e *Explore) WithFields(fields []ExploreFields) *Explore {
	e.fields = fields
	return e
}

func (e *Explore) createFilterClause() string {
	var clause string
	if e.withNearText != nil {
		clause = e.withNearText.build()
	}
	return clause
}

// build query
func (e *Explore) build() string {
	fields := ""
	for _, field := range e.fields {
		fields += fmt.Sprintf("%v ", field)
	}

	filterClause := e.createFilterClause()

	query := fmt.Sprintf("{Explore(%v){%v}}", filterClause, fields)

	return query
}

// Do execute explore search
func (e *Explore) Do(ctx context.Context) (*models.GraphQLResponse, error) {
	return runGraphQLQuery(ctx, e.connection, e.build())
}
