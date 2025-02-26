package graphql

import (
	"context"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
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

// Get queries with Multiple Class
func (api *API) MultiClassGet() *MultiClassBuilder {
	return &MultiClassBuilder{
		connection:    api.connection,
		classBuilders: make(map[string]*GetBuilder),
	}
}

// Explore queries
func (api *API) Explore() *Explore {
	return &Explore{connection: api.connection}
}

// Aggregate queries
func (api *API) Aggregate() *AggregateBuilder {
	return &AggregateBuilder{connection: api.connection}
}

// Raw creates a raw GraphQL query
func (api *API) Raw() *Raw {
	return &Raw{connection: api.connection}
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

// NearImageArgBuilder nearImage clause
func (api *API) NearAudioArgBuilder() *NearAudioArgumentBuilder {
	return &NearAudioArgumentBuilder{}
}

// NearImageArgBuilder nearImage clause
func (api *API) NearVideoArgBuilder() *NearVideoArgumentBuilder {
	return &NearVideoArgumentBuilder{}
}

// NearImageArgBuilder nearImage clause
func (api *API) NearDepthArgBuilder() *NearDepthArgumentBuilder {
	return &NearDepthArgumentBuilder{}
}

// NearImageArgBuilder nearImage clause
func (api *API) NearThermalArgBuilder() *NearThermalArgumentBuilder {
	return &NearThermalArgumentBuilder{}
}

// NearImageArgBuilder nearImage clause
func (api *API) NearImuArgBuilder() *NearImuArgumentBuilder {
	return &NearImuArgumentBuilder{}
}

// GroupArgBuilder nearImage clause
func (api *API) GroupArgBuilder() *GroupArgumentBuilder {
	return &GroupArgumentBuilder{}
}

// Bm25ArgBuilder bm25 clause
func (api *API) Bm25ArgBuilder() *BM25ArgumentBuilder {
	return &BM25ArgumentBuilder{}
}

// HybridArgumentBuilder hybrid clause
func (api *API) HybridArgumentBuilder() *HybridArgumentBuilder {
	return &HybridArgumentBuilder{}
}

// MultiTargetArgumentBuilder targets clause
func (api *API) MultiTargetArgumentBuilder() *MultiTargetArgumentBuilder {
	return &MultiTargetArgumentBuilder{}
}

// HybridSearchesArgumentBuilder hybrid.searches clause
func (api *API) HybridSearchesArgumentBuilder() *HybridSearchesArgumentBuilder {
	return &HybridSearchesArgumentBuilder{}
}

// GroupByArgBuilder groupBy clause
func (api *API) GroupByArgBuilder() *GroupByArgumentBuilder {
	return &GroupByArgumentBuilder{}
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
