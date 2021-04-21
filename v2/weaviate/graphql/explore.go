package graphql

import (
	"context"
	"fmt"
	"strings"

	"github.com/semi-technologies/weaviate/entities/models"
)

// Explore query builder
type Explore struct {
	connection     rest
	fields         []ExploreFields
	withNearText   *NearTextArgumentBuilder
	withNearObject *NearObjectArgumentBuilder
}

// WithNearText adds nearText to clause
func (e *Explore) WithNearText(nearText *NearTextArgumentBuilder) *Explore {
	e.withNearText = nearText
	return e
}

// WithNearObject adds nearObject to clause
func (e *Explore) WithNearObject(nearObject *NearObjectArgumentBuilder) *Explore {
	e.withNearObject = nearObject
	return e
}

// WithFields that should be included in the result set
func (e *Explore) WithFields(fields []ExploreFields) *Explore {
	e.fields = fields
	return e
}

func (e *Explore) createFilterClause() string {
	filters := []string{}
	if e.withNearText != nil {
		filters = append(filters, e.withNearText.build())
	}
	if e.withNearObject != nil {
		filters = append(filters, e.withNearObject.build())
	}
	return strings.Join(filters, ", ")
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
