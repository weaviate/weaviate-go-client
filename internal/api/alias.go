package api

import (
	"net/http"

	"github.com/weaviate/weaviate-go-client/v6/internal/transport"
)

type DeleteAliasRequest struct {
	transport.BaseEndpoint

	Alias string
}

var _ transport.Endpoint = (*DeleteAliasRequest)(nil)

func (d *DeleteAliasRequest) Method() string { return http.MethodDelete }
func (d *DeleteAliasRequest) Path() string   { return "/aliases/" + d.Alias }
