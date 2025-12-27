package api

import (
	"net/http"

	"github.com/weaviate/weaviate-go-client/v6/internal/transport"
)

var DeleteAliasRequest = transport.IdentityEndpoint[string](http.MethodDelete, "/aliases/%s")
