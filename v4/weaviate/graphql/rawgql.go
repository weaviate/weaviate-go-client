package graphql

import (
	"context"

	"github.com/semi-technologies/weaviate/entities/models"
)

// RawGraphQLQueryBuilder for accepting a prebuilt query from the user
type RawGraphQLQueryBuilder struct {
	connection rest
	query      string
}

// Do execute the GraphQL query
func (gql *RawGraphQLQueryBuilder) Do(ctx context.Context) (*models.GraphQLResponse, error) {
	return runGraphQLQuery(ctx, gql.connection, gql.query)
}

// return the query string
func (gql *RawGraphQLQueryBuilder) build() string {
	return gql.query
}
