package graphql

import (
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/connection"
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
