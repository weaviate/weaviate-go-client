package graphql

import (
	"context"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/except"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate/entities/models"
	"net/http"
)

// API group for GrapQL
type API struct {
	Connection *connection.Connection
}

// Get queries
func (api *API) Get() *Get {
	return &Get{connection: api.Connection}
}

// Explore queries
func (api *API) Explore() *Explore {
	return &Explore{connection:    api.Connection}
}

// Aggregate queries
func (api *API) Aggregate() *Aggregate {
	return &Aggregate{connection: api.Connection}
}

// rest requests abstraction
type rest interface {
	//RunREST request to weaviate
	RunREST(ctx context.Context, path string, restMethod string, requestBody interface{}) (*connection.ResponseData, error)
}

func runGraphQLQuery(ctx context.Context, rest rest, query string) (*models.GraphQLResponse, error) {
	// Do execute the GraphQL query
	gqlQuery := models.GraphQLQuery{
		Query:         query,
	}
	responseData, responseErr := rest.RunREST(ctx, "/graphql", http.MethodPost, &gqlQuery)
	err := except.CheckResponnseDataErrorAndStatusCode(responseData, responseErr, 200)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	var gqlResponse models.GraphQLResponse
	parseErr := responseData.DecodeBodyIntoTarget(&gqlResponse)
	return &gqlResponse, except.NewDerivedWeaviateClientError(parseErr)
}