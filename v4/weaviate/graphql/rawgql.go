package graphql

// RawGQLQueryBuilder for accepting a prebuilt query from the user
type RawGQLQueryBuilder struct {
	connection rest
	query      string
}

// return the query string
func (gql *RawGQLQueryBuilder) build() string {
	return gql.query
}
