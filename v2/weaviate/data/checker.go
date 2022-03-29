package data

import (
	"context"
	"fmt"
	"net/http"

	"github.com/semi-technologies/weaviate-go-client/v2/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/v2/weaviate/except"
)

// Checker builder to check data object's existence
type Checker struct {
	connection *connection.Connection
	uuid       string
}

// WithID specifies the uuid of an data object to be checked
func (checker *Checker) WithID(uuid string) *Checker {
	checker.uuid = uuid
	return checker
}

// Do check the specified data object if it exists in weaviate
func (checker *Checker) Do(ctx context.Context) (bool, error) {
	path := fmt.Sprintf("/objects/%v", checker.uuid)
	responseData, err := checker.connection.RunREST(ctx, path, http.MethodHead, nil)
	exists := responseData.StatusCode == 204
	return exists, except.CheckResponseDataErrorAndStatusCode(responseData, err, 204, 404)
}
