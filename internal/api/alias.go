package api

import (
	"encoding/json"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/rest"
	"github.com/weaviate/weaviate-go-client/v6/internal/transports"
)

var DeleteAliasRequest = transports.IdentityEndpoint[string](http.MethodDelete, "/aliases/%s")

type CreateAliasRequest struct {
	transports.BaseEndpoint

	Alias      string // Required: alias name.
	Collection string // Required: collection.
}

var (
	_ transports.Endpoint = (*CreateAliasRequest)(nil)
	_ json.Marshaler      = (*CreateAliasRequest)(nil)
)

func (*CreateAliasRequest) Method() string { return http.MethodPost }
func (r *CreateAliasRequest) Path() string { return "/aliases" }
func (r *CreateAliasRequest) Body() any    { return r }

// MarshalJSON implements json.Marshaler via rest.Alias.
func (r *CreateAliasRequest) MarshalJSON() ([]byte, error) {
	req := &rest.Alias{Alias: r.Alias, Class: r.Collection}
	return json.Marshal(req)
}

type Alias struct {
	Alias, Collection string
}

var _ json.Unmarshaler = (*Alias)(nil)

// UnmarshalJSON implements json.Unmarshaler.
func (b *Alias) UnmarshalJSON(data []byte) error {
	var resp rest.Alias
	if err := json.Unmarshal(data, &resp); err != nil {
		return err
	}
	*b = Alias{Alias: resp.Alias, Collection: resp.Class}
	return nil
}
