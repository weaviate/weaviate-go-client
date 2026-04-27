package tokenize

import (
	"context"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
)

// TextTokenizer is a builder for tokenizing arbitrary text via
// POST /v1/tokenize.
type TextTokenizer struct {
	connection      *connection.Connection
	text            string
	tokenization    Tokenization
	analyzerConfig  *AnalyzerConfig
	stopwordPresets map[string]*StopwordConfig
}

// WithText sets the text to tokenize.
func (b *TextTokenizer) WithText(text string) *TextTokenizer {
	b.text = text
	return b
}

// WithTokenization sets the tokenization method to apply.
func (b *TextTokenizer) WithTokenization(t Tokenization) *TextTokenizer {
	b.tokenization = t
	return b
}

// WithAnalyzerConfig sets the optional analyzer configuration (ASCII folding,
// stopword preset, ...).
func (b *TextTokenizer) WithAnalyzerConfig(cfg *AnalyzerConfig) *TextTokenizer {
	b.analyzerConfig = cfg
	return b
}

// WithStopwordPresets sets custom named stopword presets that can be
// referenced by AnalyzerConfig.StopwordPreset.
func (b *TextTokenizer) WithStopwordPresets(presets map[string]*StopwordConfig) *TextTokenizer {
	b.stopwordPresets = presets
	return b
}

// Do performs the tokenize request.
func (b *TextTokenizer) Do(ctx context.Context) (*TokenizeResult, error) {
	payload := tokenizeRequest{
		Text:            b.text,
		Tokenization:    b.tokenization,
		AnalyzerConfig:  b.analyzerConfig,
		StopwordPresets: b.stopwordPresets,
	}

	responseData, err := b.connection.RunREST(ctx, "/tokenize", http.MethodPost, payload)
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
