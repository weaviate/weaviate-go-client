package tokenize

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v6/collections"
	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
)

func NewClient(t internal.Transport) *Client {
	dev.AssertNotNil(t, "transport")
	return &Client{transport: t}
}

type Client struct {
	transport internal.Transport
}

type Result api.TokenizeResponse

// PropertyOptions specify a collection property
// to inherit the tokenization configuration from.
type PropertyOptions struct {
	CollectionName string
	PropetyName    string
}

// Tokenize a text string using existing property configuration.
func (c *Client) ForProperty(ctx context.Context, text string, options PropertyOptions) (*Result, error) {
	req := api.TokenizeRequest{
		RequestDefaults: api.RequestDefaults{CollectionName: options.CollectionName},
		PropetyName:     options.PropetyName,
		Text:            text,
	}

	var resp api.TokenizeResponse
	if err := c.transport.Do(ctx, &req, &resp); err != nil {
		return nil, fmt.Errorf("tokenize: %w", err)
	}

	return (*Result)(&resp), nil
}

// ExplicitOptions provide custom tokenization configuration.
type ExplicitOptions struct {
	Tokenization    collections.Tokenization
	TextAnalyzer    *collections.TextAnalyzerConfig
	Stopwords       *collections.StopwordConfig
	StopwordPresets map[string][]string
}

// Tokenize a text string using explicit tokenization configuration.
func (c *Client) Text(ctx context.Context, text string, options ExplicitOptions) (*Result, error) {
	req := api.TokenizeRequest{
		Text:            text,
		Tokenization:    api.Tokenization(options.Tokenization),
		StopwordPresets: options.StopwordPresets,
		TextAnalyzer:    (*api.TextAnalyzerConfig)(options.TextAnalyzer),
		Stopwords:       (*api.StopwordConfig)(options.Stopwords),
	}

	var resp api.TokenizeResponse
	if err := c.transport.Do(ctx, &req, &resp); err != nil {
		return nil, fmt.Errorf("tokenize: %w", err)
	}

	return (*Result)(&resp), nil
}
