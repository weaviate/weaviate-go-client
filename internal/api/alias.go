package api

import (
	"net/http"

	"github.com/weaviate/weaviate-go-client/v6/internal/transports"
)

var DeleteAliasRequest = transports.IdentityEndpoint[string](http.MethodDelete, "/aliases/%s")
