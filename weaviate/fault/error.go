package fault

import (
	"fmt"
)

// WeaviateClientError is returned if the client experienced an error.
//
//	If the error is due to weaviate returning an unexpected status code the IsUnexpectedStatusCode field will be true
//	 and the StatusCode field will be set
//	If the error occurred for another reason the DerivedFromError will be set and IsUnexpectedStatusCode will be false
type WeaviateClientError struct {
	IsUnexpectedStatusCode bool
	StatusCode             int
	Msg                    string
	DerivedFromError       error
}

// Error message of the unexpected status code error
func (uce *WeaviateClientError) Error() string {
	msg := uce.Msg
	if uce.DerivedFromError != nil {
		msg = fmt.Sprintf("%s: %s", uce.Msg, uce.DerivedFromError.Error())
	}
	return fmt.Sprintf("status code: %v, error: %v", uce.StatusCode, msg)
}

// GoString makes WeaviateClientError satisfy the GoStringer interface. This allows the correct display when using the formatting %#v, used in assert.Nil().
func (uce *WeaviateClientError) GoString() string {
	return uce.Error()
}
