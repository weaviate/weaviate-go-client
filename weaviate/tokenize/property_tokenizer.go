package tokenize

import (
	"context"
	"encoding/json"
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
	// Weaviate >=1.37.4 no longer echoes `tokenization`; backfill from the schema.
	// Only paid when the server omitted the field — older servers skip the lookup.
	if result.Tokenization == "" {
		if t, _ := p.fetchPropertyTokenization(ctx); t != "" {
			result.Tokenization = Tokenization(t)
		}
	}
	return &result, nil
}

// fetchPropertyTokenization returns "" on any failure: backfill is best-effort
// and must not mask the tokenize result with a schema error.
func (p *PropertyTokenizer) fetchPropertyTokenization(ctx context.Context) (string, error) {
	resp, err := p.connection.RunREST(ctx, "/schema/"+p.className, http.MethodGet, nil)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", err
	}
	var class struct {
		Properties []struct {
			Name         string `json:"name"`
			Tokenization string `json:"tokenization"`
		} `json:"properties"`
	}
	if err := json.Unmarshal(resp.Body, &class); err != nil {
		return "", err
	}
	for _, prop := range class.Properties {
		if prop.Name == p.propertyName {
			return prop.Tokenization, nil
		}
	}
	return "", nil
}
