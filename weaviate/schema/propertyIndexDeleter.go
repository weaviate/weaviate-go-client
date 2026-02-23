package schema

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
)

// PropertyIndexDeleter is a builder to delete a property's index from a schema class
type PropertyIndexDeleter struct {
	connection   *connection.Connection
	className    string
	propertyName string
	indexName    string
}

// WithClassName defines the name of the class for which a given property's index will be deleted
func (p *PropertyIndexDeleter) WithClassName(className string) *PropertyIndexDeleter {
	p.className = className
	return p
}

// WithPropertyName defines the name of the class's property for which the index will be deleted
func (p *PropertyIndexDeleter) WithPropertyName(propertyName string) *PropertyIndexDeleter {
	p.propertyName = propertyName
	return p
}

// WithFilterable defines filterable property index to be deleted
func (p *PropertyIndexDeleter) WithFilterable() *PropertyIndexDeleter {
	p.indexName = "filterable"
	return p
}

// WithSearchable defines searchable property index to be deleted
func (p *PropertyIndexDeleter) WithSearchable() *PropertyIndexDeleter {
	p.indexName = "searchable"
	return p
}

// WithRangeFilters defines rangeFilters property index to be deleted
func (p *PropertyIndexDeleter) WithRangeFilters() *PropertyIndexDeleter {
	p.indexName = "rangeFilters"
	return p
}

// Do deletes the property's index
func (p *PropertyIndexDeleter) Do(ctx context.Context) error {
	path := fmt.Sprintf("/schema/%v/properties/%s/index/%s", p.className, p.propertyName, p.indexName)
	responseData, err := p.connection.RunREST(ctx, path, http.MethodDelete, nil)
	return except.CheckResponseDataErrorAndStatusCode(responseData, err, 200)
}
