package schema

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
)

// VectorIndexDeleter is a builder to delete a vector index from a schema class
type VectorIndexDeleter struct {
	connection      *connection.Connection
	className       string
	vectorIndexName string
}

// WithClassName defines the name of the class for which the vector index will be deleted
func (v *VectorIndexDeleter) WithClassName(className string) *VectorIndexDeleter {
	v.className = className
	return v
}

// WithVectorIndexName defines the name of the vector index to be deleted
func (v *VectorIndexDeleter) WithVectorIndexName(vectorIndexName string) *VectorIndexDeleter {
	v.vectorIndexName = vectorIndexName
	return v
}

// Do deletes the vector index
func (v *VectorIndexDeleter) Do(ctx context.Context) error {
	path := fmt.Sprintf("/schema/%v/vectors/%s/index", v.className, v.vectorIndexName)
	responseData, err := v.connection.RunREST(ctx, path, http.MethodDelete, nil)
	return except.CheckResponseDataErrorAndStatusCode(responseData, err, 200)
}
