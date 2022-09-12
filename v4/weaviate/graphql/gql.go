package graphql

import (
	"context"
	"net/http"

	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/except"
	"github.com/semi-technologies/weaviate/entities/models"
)

// API group for GraphQL
type API struct {
	connection *connection.Connection
}

// New GraphQL api group from connection
func New(con *connection.Connection) *API {
	return &API{connection: con}
}

// Get queries
func (api *API) Get() *GetBuilder {
	return &GetBuilder{connection: api.connection}
}

// Explore queries
func (api *API) Explore() *Explore {
	return &Explore{connection: api.connection}
}

// Aggregate queries
func (api *API) Aggregate() *AggregateBuilder {
	return &AggregateBuilder{connection: api.connection}
}

// RawGraphQLQueryBuilder creates a raw GraphQL query
func (api *API) Raw(query string) *RawGraphQLQueryBuilder {
	return &RawGraphQLQueryBuilder{connection: api.connection, query: query}
}

// NearTextArgBuilder nearText clause
func (api *API) NearTextArgBuilder() *NearTextArgumentBuilder {
	return &NearTextArgumentBuilder{}
}

// NearObjectArgBuilder nearObject clause
func (api *API) NearObjectArgBuilder() *NearObjectArgumentBuilder {
	return &NearObjectArgumentBuilder{}
}

// NearVectorArgBuilder nearVector clause
func (api *API) NearVectorArgBuilder() *NearVectorArgumentBuilder {
	return &NearVectorArgumentBuilder{}
}

// AskArgBuilder ask clause
func (api *API) AskArgBuilder() *AskArgumentBuilder {
	return &AskArgumentBuilder{}
}

// NearImageArgBuilder nearImage clause
func (api *API) NearImageArgBuilder() *NearImageArgumentBuilder {
	return &NearImageArgumentBuilder{}
}

// GroupArgBuilder nearImage clause
func (api *API) GroupArgBuilder() *GroupArgumentBuilder {
	return &GroupArgumentBuilder{}
}

// rest requests abstraction
type rest interface {
	// RunREST request to weaviate
	RunREST(ctx context.Context, path string, restMethod string, requestBody interface{}) (*connection.ResponseData, error)
}

func runGraphQLQuery(ctx context.Context, rest rest, query string) (*models.GraphQLResponse, error) {
	// Do execute the GraphQL query
	gqlQuery := models.GraphQLQuery{
		Query: query,
	}
	responseData, responseErr := rest.RunREST(ctx, "/graphql", http.MethodPost, &gqlQuery)
	err := except.CheckResponseDataErrorAndStatusCode(responseData, responseErr, 200)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	var gqlResponse models.GraphQLResponse
	parseErr := responseData.DecodeBodyIntoTarget(&gqlResponse)
	return &gqlResponse, parseErr
}
