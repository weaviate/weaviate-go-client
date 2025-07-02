package alias

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

// ShardsGetter builder object to get a class' shards
type AliasList struct {
	connection *connection.Connection
	className  string
}

// WithClassName specifies the class to which the shards belong
func (s *AliasList) WithClassName(className string) *AliasList {
	s.className = strings.TrimSpace(className)
	return s
}

// Do get the status of the shards of the class specified in AliasList
func (s *AliasList) Do(ctx context.Context) ([]models.Alias, error) {
	return listAlias(ctx, s.connection, s.className)
}

func listAlias(ctx context.Context, conn *connection.Connection, className string) ([]models.Alias, error) {
	url := "/aliases"
	if className != "" {
		url = fmt.Sprintf("/aliases?class=%s", className)
	}

	responseData, err := conn.RunREST(ctx, url, http.MethodGet, nil)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}

	if responseData.StatusCode == 200 {
		resp := struct {
			Aliases []models.Alias `json:"aliases"`
		}{}
		decodeErr := responseData.DecodeBodyIntoTarget(&resp)
		return resp.Aliases, decodeErr
	}
	return nil, except.NewWeaviateClientError(responseData.StatusCode, string(responseData.Body))
}
