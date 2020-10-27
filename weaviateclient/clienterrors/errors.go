package clienterrors

import (
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
)

// UnexpectedStatusCodeError is returned if the status code was not the one indicating a successful request
// Contains both the code and the error message
type UnexpectedStatusCodeError struct {
	StatusCode int
	msg        string
}

// NewUnexpectedStatusCodeError from status code and error message
func NewUnexpectedStatusCodeError(statusCode int, format string, args ...interface{}) *UnexpectedStatusCodeError {
	return &UnexpectedStatusCodeError{
		StatusCode: statusCode,
		msg:        fmt.Sprintf(format, args...),
	}
}

// NewUnexpectedStatusCodeErrorFromRESTResponse creates the error based on a repsonse data object
func NewUnexpectedStatusCodeErrorFromRESTResponse(responseData *connection.ResponseData) *UnexpectedStatusCodeError {
	return NewUnexpectedStatusCodeError(responseData.StatusCode, string(responseData.Body))
}

// Error message of the unexpected status code error
func (uce *UnexpectedStatusCodeError) Error() string {
	return fmt.Sprintf("status code: %v, error: %v", uce.StatusCode, uce.msg)
}

// CheckResponnseDataErrorAndStatusCode returns the response error if it is not nil,
//  and an UnexpectedStatusCodeError if the status code is not matching
func CheckResponnseDataErrorAndStatusCode(responseData *connection.ResponseData, responseErr error, expectedStatusCode int) error {
	if responseErr != nil {
		return responseErr
	}
	if responseData.StatusCode == expectedStatusCode {
		return nil
	}
	return NewUnexpectedStatusCodeErrorFromRESTResponse(responseData)
}
