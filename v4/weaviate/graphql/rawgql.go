package graphql

import (
	"context"

	"github.com/weaviate/weaviate/entities/models"
)

// Raw for accepting a prebuilt query from the user
type Raw struct {
	connection rest
	query      string
}

// Do execute the GraphQL query
func (gql *Raw) Do(ctx context.Context) (*models.GraphQLResponse, error) {
	return runGraphQLQuery(ctx, gql.connection, gql.build())
}

// WithQuery the query string
func (b *Raw) WithQuery(query string) *Raw {
	b.query = query
	return b
}

// return the query string
func (gql *Raw) build() string {
	return gql.query
}
