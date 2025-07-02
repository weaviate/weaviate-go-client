package alias

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
)

// Getter builder object to get a list of aliases
type Getter struct {
	connection *connection.Connection
	className  string
}

// WithClassName specifies the class to which the alias belongs to
func (s *Getter) WithClassName(className string) *Getter {
	s.className = strings.TrimSpace(className)
	return s
}

// Do get the list of alias
func (s *Getter) Do(ctx context.Context) ([]Alias, error) {
	return listAlias(ctx, s.connection, s.className)
}

func listAlias(ctx context.Context, conn *connection.Connection, className string) ([]Alias, error) {
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
			Aliases []Alias `json:"aliases"`
		}{}
		decodeErr := responseData.DecodeBodyIntoTarget(&resp)
		return resp.Aliases, decodeErr
	}
	return nil, except.NewWeaviateClientError(responseData.StatusCode, string(responseData.Body))
}
