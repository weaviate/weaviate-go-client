package api

import (
	"encoding/json"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/rest"
	"github.com/weaviate/weaviate-go-client/v6/internal/transports"
)

type Alias struct {
	Collection string
	Alias      string
}

type CreateAliasRequest struct {
	transports.BaseEndpoint
	Alias
}

var _ transports.Endpoint = (*CreateAliasRequest)(nil)

func (CreateAliasRequest) Method() string { return http.MethodPost }
func (CreateAliasRequest) Path() string   { return "/aliases" }
func (r *CreateAliasRequest) Body() any {
	return &rest.AliasesCreateJSONRequestBody{
		Class: r.Collection,
		Alias: r.Alias.Alias,
	}
}

var (
	// GetAliasRequest retrieves an alias by name.
	// Use with [Alias] response dest.
	GetAliasRequest = transports.IdentityEndpoint[string](http.MethodGet, "/aliases/%s")
	// DeleteAliasesRequests deletes an alias by name.
	DeleteAliasRequest = transports.IdentityEndpoint[string](http.MethodDelete, "/aliases/%s")
	// ListAliasesRequests lists all defined aliases.
	// Use with [ListAliasesResponse].
	ListAliasesRequest = transports.StaticEndpoint(http.MethodGet, "/aliases")
)

// ListAliasesResponse unmarshals a collection of aliases.
type ListAliasesResponse []Alias

// UpdateAliasRequests assigns an existing alias to a different collection.
type UpdateAliasRequest struct {
	transports.BaseEndpoint
	Alias
}

var _ transports.Endpoint = (*UpdateAliasRequest)(nil)

func (UpdateAliasRequest) Method() string  { return http.MethodPut }
func (r *UpdateAliasRequest) Path() string { return "/aliases/" + r.Alias.Alias }
func (r *UpdateAliasRequest) Body() any {
	return &rest.AliasesUpdateJSONRequestBody{
		Class: r.Collection,
	}
}

func (a *Alias) UnmarshalJSON(data []byte) error {
	var alias rest.Alias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*a = Alias{
		Collection: alias.Class,
		Alias:      alias.Alias,
	}
	return nil
}

func (r *ListAliasesResponse) UnmarshalJSON(data []byte) error {
	var list rest.AliasResponse
	if err := json.Unmarshal(data, &list); err != nil {
		return err
	}
	for _, a := range list.Aliases {
		*r = append(*r, Alias{
			Collection: a.Class,
			Alias:      a.Alias,
		})
	}
	return nil
}
