package schema

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
)

// ClassExistenceChecker builder to check if a class is part of a weaviate schema
type ClassExistenceChecker struct {
	connection *connection.Connection
	className  string
}

// WithClassName defines the name of the class that should be checked
func (cd *ClassExistenceChecker) WithClassName(className string) *ClassExistenceChecker {
	cd.className = className
	return cd
}

// Do check if the class is part of the weaviate schema
func (cd *ClassExistenceChecker) Do(ctx context.Context) (bool, error) {
	responseData, err := cd.connection.RunREST(ctx, fmt.Sprintf("/schema/%s", cd.className), http.MethodGet, nil)
	if err != nil {
		return false, except.NewDerivedWeaviateClientError(err)
	}
	if responseData.StatusCode == 200 {
		return true, nil
	}
	return false, nil
}
