package graphql

// RawGraphQLQueryBuilder for accepting a prebuilt query from the user
type RawGraphQLQueryBuilder struct {
	connection rest
	query      string
}

// return the query string
func (gql *RawGraphQLQueryBuilder) build() string {
	return gql.query
}
