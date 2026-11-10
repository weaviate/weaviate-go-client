package graphql

import (
	"context"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

// API group for GraphQL
//
// Deprecated: the GraphQL API is deprecated and is planned to be disabled and removed.
// Migrate to github.com/weaviate/weaviate-go-client/v6, which queries over gRPC.
// See https://docs.weaviate.io/weaviate/api/graphql/migration.
type API struct {
	connection *connection.Connection
}

// New GraphQL api group from connection
func New(con *connection.Connection) *API {
	return &API{connection: con}
}

// Get queries
//
// Deprecated: the GraphQL API is deprecated and is planned to be disabled and removed.
// Migrate to github.com/weaviate/weaviate-go-client/v6, which queries over gRPC.
// See https://docs.weaviate.io/weaviate/api/graphql/migration.
func (api *API) Get() *GetBuilder {
	return &GetBuilder{connection: api.connection}
}

// Get queries with Multiple Class
//
// Deprecated: the GraphQL API is deprecated and is planned to be disabled and removed.
// Migrate to github.com/weaviate/weaviate-go-client/v6, which queries over gRPC.
// See https://docs.weaviate.io/weaviate/api/graphql/migration.
func (api *API) MultiClassGet() *MultiClassBuilder {
	return &MultiClassBuilder{
		connection:    api.connection,
		classBuilders: make(map[string]*GetBuilder),
	}
}

// Explore queries
//
// Deprecated: the GraphQL API is deprecated and is planned to be disabled and removed.
// Migrate to github.com/weaviate/weaviate-go-client/v6, which queries over gRPC.
// See https://docs.weaviate.io/weaviate/api/graphql/migration.
func (api *API) Explore() *Explore {
	return &Explore{connection: api.connection}
}

// Aggregate queries
//
// Deprecated: the GraphQL API is deprecated and is planned to be disabled and removed.
// Migrate to github.com/weaviate/weaviate-go-client/v6, which queries over gRPC.
// See https://docs.weaviate.io/weaviate/api/graphql/migration.
func (api *API) Aggregate() *AggregateBuilder {
	return &AggregateBuilder{connection: api.connection}
}

// Raw creates a raw GraphQL query
//
// Deprecated: the GraphQL API is deprecated and is planned to be disabled and removed.
// Migrate to github.com/weaviate/weaviate-go-client/v6, which queries over gRPC.
// See https://docs.weaviate.io/weaviate/api/graphql/migration.
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
	warnGraphQLDeprecated()
	responseData, responseErr := rest.RunREST(ctx, "/graphql", http.MethodPost, &gqlQuery)
	err := except.CheckResponseDataErrorAndStatusCode(responseData, responseErr, 200)
	if err != nil {
		if graphQLDisabled(responseData) {
			return nil, except.NewWeaviateClientError(responseData.StatusCode, graphQLDisabledMessage)
		}
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	var gqlResponse models.GraphQLResponse
	parseErr := responseData.DecodeBodyIntoTarget(&gqlResponse)
	return &gqlResponse, parseErr
}

// graphQLDeprecationNotice points v5 users to the client that does not use GraphQL.
const graphQLDeprecationNotice = "the GraphQL API is deprecated and is planned to be disabled and removed. " +
	"weaviate-go-client/v5 sends Get, Aggregate, Explore, and Raw queries over GraphQL; " +
	"migrate to github.com/weaviate/weaviate-go-client/v6, which queries over gRPC. " +
	"See https://docs.weaviate.io/weaviate/api/graphql/migration"

const graphQLDisabledMessage = "graphql api is disabled on this Weaviate instance. " + graphQLDeprecationNotice

var graphQLWarning sync.Once

// warnGraphQLDeprecated logs the deprecation notice once per process.
func warnGraphQLDeprecated() {
	graphQLWarning.Do(func() {
		log.Printf("weaviate-go-client: %s", graphQLDeprecationNotice)
	})
}

// graphQLDisabled reports whether Weaviate refused the query because
// DISABLE_GRAPHQL is set.
func graphQLDisabled(responseData *connection.ResponseData) bool {
	return responseData != nil && strings.Contains(string(responseData.Body), "graphql api is disabled")
}
