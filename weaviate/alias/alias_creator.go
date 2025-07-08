package alias

import (
	"context"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
)

// AliasCreator builds object to create an alias
type AliasCreator struct {
	connection *connection.Connection
	alias      *Alias
}

// WithClass specifies the alias that will be added to the schema
func (cc *AliasCreator) WithAlias(alias *Alias) *AliasCreator {
	cc.alias = alias
	return cc
}

// Do create a alias in the schema as specified in the builder
func (cc *AliasCreator) Do(ctx context.Context) error {
	responseData, err := cc.connection.RunREST(ctx, "/aliases", http.MethodPost, cc.alias)
	return except.CheckResponseDataErrorAndStatusCode(responseData, err, 200)
}
