package tokenize

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
)

// PropertyTokenizer is a builder for
// POST /v1/schema/{className}/properties/{propertyName}/tokenize.
// It tokenizes text using the configuration of an existing property on a
// collection (requires Weaviate >= 1.37.0).
type PropertyTokenizer struct {
	connection   *connection.Connection
	className    string
	propertyName string
	text         string
}

// WithClassName sets the name of the class that owns the property.
func (p *PropertyTokenizer) WithClassName(className string) *PropertyTokenizer {
	p.className = className
	return p
}

// WithPropertyName sets the name of the property whose tokenization
// configuration to apply.
func (p *PropertyTokenizer) WithPropertyName(propertyName string) *PropertyTokenizer {
	p.propertyName = propertyName
	return p
}

// WithText sets the text to tokenize.
func (p *PropertyTokenizer) WithText(text string) *PropertyTokenizer {
	p.text = text
	return p
}

// Do performs the tokenize request.
func (p *PropertyTokenizer) Do(ctx context.Context) (*TokenizeResult, error) {
	path := fmt.Sprintf("/schema/%s/properties/%s/tokenize", p.className, p.propertyName)
	payload := struct {
		Text string `json:"text"`
	}{Text: p.text}

	responseData, err := p.connection.RunREST(ctx, path, http.MethodPost, payload)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	if responseData.StatusCode != http.StatusOK {
		return nil, except.NewUnexpectedStatusCodeErrorFromRESTResponse(responseData)
	}

	var result TokenizeResult
	if decodeErr := responseData.DecodeBodyIntoTarget(&result); decodeErr != nil {
		return nil, decodeErr
	}
	return &result, nil
}
