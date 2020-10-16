package clienterrors

import (
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
)

type UnexpectedStatusCodeError struct {
	StatusCode int
	msg string
}

func NewUnexpectedStatusCodeError(statusCode int, format string, args ...interface{}) *UnexpectedStatusCodeError {
	return &UnexpectedStatusCodeError{
		StatusCode: statusCode,
		msg: fmt.Sprintf(format, args...),
	}
}

func NewUnexpectedStatusCodeErrorFromRESTResponse(responseData *connection.ResponseData) *UnexpectedStatusCodeError {
	return NewUnexpectedStatusCodeError(responseData.StatusCode, string(responseData.Body))
}

func (uce *UnexpectedStatusCodeError) Error() string {
	return fmt.Sprintf("status code: %v, error: %v", uce.StatusCode, uce.msg)
}