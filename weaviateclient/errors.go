package weaviateclient

import "fmt"

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

func (uce *UnexpectedStatusCodeError) Error() string {
	return fmt.Sprintf("status code: %v, error: %v", uce.StatusCode, uce.msg)
}
