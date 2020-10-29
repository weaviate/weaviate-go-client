package except

import (
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/fault"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
)


// NewWeaviateClientError from status code and error message
func NewWeaviateClientError(statusCode int, format string, args ...interface{}) *fault.WeaviateClientError {
	return &fault.WeaviateClientError{
		IsUnexpectedStatusCode: true,
		StatusCode: statusCode,
		Msg:        fmt.Sprintf(format, args...),
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

// NewUnexpectedStatusCodeErrorFromRESTResponse creates the error based on a repsonse data object
func NewUnexpectedStatusCodeErrorFromRESTResponse(responseData *connection.ResponseData) *fault.WeaviateClientError {
	return NewWeaviateClientError(responseData.StatusCode, string(responseData.Body))
}

// CheckResponnseDataErrorAndStatusCode returns the response error if it is not nil,
//  and an WeaviateClientError if the status code is not matching
func CheckResponnseDataErrorAndStatusCode(responseData *connection.ResponseData, responseErr error, expectedStatusCode int) error {
	if responseErr != nil {
		return NewDerivedWeaviateClientError(responseErr)
	}
	if responseData.StatusCode == expectedStatusCode {
		return nil
	}
	return NewUnexpectedStatusCodeErrorFromRESTResponse(responseData)
}
