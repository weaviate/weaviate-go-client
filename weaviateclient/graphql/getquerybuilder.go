package graphql

import (
	"context"
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/clienterrors"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/paragons"
	"github.com/semi-technologies/weaviate/entities/models"
	"net/http"
	"strings"
)

type GetBuilder struct {
	connection rest
	semanticKind paragons.SemanticKind
	className string
	withFields string

	includesFilterClause bool // true if brackets behind class is needed
	withWhereFilter      string
	includesLimit        bool
	limit                int
	withExploreFilter string
	withGroupFilter string
}

func (gb *GetBuilder) WithClassName(name string) *GetBuilder {
	gb.className = name
	return gb
}

func (gb *GetBuilder) WithFields(fields string) *GetBuilder {
	gb.withFields = fields
	return gb
}

func (gb *GetBuilder) WithWhere(filter string) *GetBuilder {
	gb.includesFilterClause = true
	gb.withWhereFilter = filter
	return gb
}

func (gb *GetBuilder) WithLimit(limit int) *GetBuilder {
	gb.includesFilterClause = true
	gb.includesLimit = true
	gb.limit = limit
	return gb
}

func (gb *GetBuilder) WithExplore(explore string) *GetBuilder {
	gb.includesFilterClause = true
	gb.withExploreFilter = explore
	return gb
}

func (gb *GetBuilder) WithGroup(group string) *GetBuilder {
	gb.includesFilterClause = true
	gb.withGroupFilter = group
	return gb
}

// Do execute the GraphQL query
func (gb *GetBuilder) Do(ctx context.Context) (*models.GraphQLResponse, error) {
	query := models.GraphQLQuery{
		Query:         gb.build(),
	}
	responseData, responseErr := gb.connection.RunREST(ctx, "/graphql", http.MethodPost, &query)
	err := clienterrors.CheckResponnseDataErrorAndStatusCode(responseData, responseErr, 200)
	if err != nil {
		return nil, err
	}
	var gqlResponse models.GraphQLResponse
	parseErr := responseData.DecodeBodyIntoTarget(&gqlResponse)
	return &gqlResponse, parseErr
}

// build the GraphQL query string (not needed when Do is executed)
func (gb *GetBuilder) build() string {
	semanticKind := strings.Title(string(gb.semanticKind))

	filterClause := ""
	if gb.includesFilterClause {
		filterClause = gb.createFilterClause()
	}

	query := fmt.Sprintf("{Get {%v {%v %v {%v}}}}", semanticKind, gb.className, filterClause, gb.withFields)

	return query
}

func (gb *GetBuilder) createFilterClause() string {
	clause := "("
	if len(gb.withWhereFilter) > 0 {
		clause += fmt.Sprintf("where: %v", gb.withWhereFilter)
	}
	if len(gb.withExploreFilter) > 0 {
		if string(clause[len(clause)-1]) == "(" {
			clause += fmt.Sprintf("explore: %v", gb.withExploreFilter)
		} else {
			clause += fmt.Sprintf(", explore: %v", gb.withExploreFilter)
		}
	}
	if len(gb.withGroupFilter) > 0 {
		if string(clause[len(clause)-1]) == "(" {
			clause += fmt.Sprintf("group: %v", gb.withGroupFilter)
		} else {
			clause += fmt.Sprintf(", group: %v", gb.withGroupFilter)
		}
	}
	if gb.includesLimit {
		if string(clause[len(clause)-1]) == "(" {
			clause += fmt.Sprintf("limit: %v", gb.limit)
		} else {
			clause += fmt.Sprintf(", limit: %v", gb.limit)
		}
	}
	clause += ")"
	return clause
}

