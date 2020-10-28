package clienterror

import "fmt"

// WeaviateClientError is returned if the client experienced an error.
//  If the error is due to weaviate returning an unexpected status code the IsUnexpectedStatusCode field will be true
//   and the StatusCode field will be set
//  If the error occurred for another reason the DerivedFromError will be set and IsUnexpectedStatusCode will be false
type WeaviateClientError struct {
	IsUnexpectedStatusCode bool
	StatusCode int
	Msg        string
	DerivedFromError error
}

// Error message of the unexpected status code error
func (uce *WeaviateClientError) Error() string {
	return fmt.Sprintf("status code: %v, error: %v", uce.StatusCode, uce.Msg)
}


