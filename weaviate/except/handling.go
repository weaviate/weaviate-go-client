package except

import (
	"fmt"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/fault"
)

// NewWeaviateClientError from status code and error message
func NewWeaviateClientError(statusCode int, format string, args ...interface{}) *fault.WeaviateClientError {
	return &fault.WeaviateClientError{
		IsUnexpectedStatusCode: true,
		StatusCode:             statusCode,
		Msg:                    fmt.Sprintf(format, args...),
	}
}

// NewDerivedWeaviateClientError wraps an error into a WeviateClientError as derived error
func NewDerivedWeaviateClientError(err error) *fault.WeaviateClientError {
	return &fault.WeaviateClientError{
		IsUnexpectedStatusCode: false,
		StatusCode:             -1,
		Msg:                    "check the DerivedFromError field for more information",
		DerivedFromError:       err,
	}
}

// NewUnexpectedStatusCodeErrorFromRESTResponse creates the error based on a response data object
func NewUnexpectedStatusCodeErrorFromRESTResponse(responseData *connection.ResponseData) *fault.WeaviateClientError {
	return NewWeaviateClientError(responseData.StatusCode, string(responseData.Body))
}

// CheckResponseDataErrorAndStatusCode returns the response error if it is not nil,
//
//	and an WeaviateClientError if the status code is not matching
func CheckResponseDataErrorAndStatusCode(responseData *connection.ResponseData, responseErr error, expectedStatusCodes ...int) error {
	if responseErr != nil {
		return NewDerivedWeaviateClientError(responseErr)
	}
	for i := range expectedStatusCodes {
		if responseData.StatusCode == expectedStatusCodes[i] {
			return nil
		}
	}
	return NewUnexpectedStatusCodeErrorFromRESTResponse(responseData)
}
