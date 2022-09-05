package graphql

import (
	"context"
	"errors"

	"github.com/semi-technologies/weaviate/entities/models"
)

// RawGQLQueryBuilder for accepting a prebuilt query from the user
type RawGQLQueryBuilder struct {
	connection rest
	query      string
}

// WithQuery a complete GraphQL query string
func (gql *RawGQLQueryBuilder) WithQuery(query string) *RawGQLQueryBuilder {
	gql.query = query
	return gql
}

func (gql *RawGQLQueryBuilder) Do(ctx context.Context) (*models.GraphQLResponse, error) {
	if gql.query == "" {
		return nil, errors.New("query cannot be empty")
	}

	response, err := runGraphQLQuery(ctx, gql.connection, gql.query)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// return the query string
func (gql *RawGQLQueryBuilder) build() string {
	return gql.query
}
