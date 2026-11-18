package tokenize_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/collections"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/tokenize"
)

func TestNewClient(t *testing.T) {
	require.Panics(t, func() {
		tokenize.NewClient(nil)
	}, "nil transport")
}

func TestClient(t *testing.T) {
	for _, tt := range []struct {
		name     string
		tokenize func(ctx context.Context, c *tokenize.Client) (*tokenize.Result, error)
		stubs    []testkit.Stub[api.TokenizeRequest, api.TokenizeResponse]
		want     *tokenize.Result
		err      testkit.Error // Expected error.
	}{
		{
			name: "for property",
			tokenize: func(ctx context.Context, c *tokenize.Client) (*tokenize.Result, error) {
				return c.ForProperty(ctx, "superduperfactorial", tokenize.PropertyOptions{
					CollectionName: "Songs",
					PropetyName:    "title",
				})
			},
			stubs: []testkit.Stub[api.TokenizeRequest, api.TokenizeResponse]{
				{
					Request: &api.TokenizeRequest{
						RequestDefaults: api.RequestDefaults{
							CollectionName: "Songs",
						},
						PropetyName: "title",
						Text:        "superduperfactorial",
					},
					Response: api.TokenizeResponse{
						Indexed: []string{"super", "duper", "factorial"},
					},
				},
			},
			want: &tokenize.Result{
				Indexed: []string{"super", "duper", "factorial"},
			},
		},
		{
			name: "text",
			tokenize: func(ctx context.Context, c *tokenize.Client) (*tokenize.Result, error) {
				return c.Text(ctx, "superduperfactorial", tokenize.ExplicitOptions{
					Tokenization: collections.TokenizationTrigram,
					TextAnalyzer: &collections.TextAnalyzerConfig{
						FoldASCII:      true,
						StopwordPreset: "en-US",
					},
				})
			},
			stubs: []testkit.Stub[api.TokenizeRequest, api.TokenizeResponse]{
				{
					Request: &api.TokenizeRequest{
						Text:         "superduperfactorial",
						Tokenization: api.TokenizationTrigram,
						TextAnalyzer: &api.TextAnalyzerConfig{
							FoldASCII:      true,
							StopwordPreset: "en-US",
						},
					},
					Response: api.TokenizeResponse{
						Indexed: []string{"super", "duper", "factorial"},
					},
				},
			},
			want: &tokenize.Result{
				Indexed: []string{"super", "duper", "factorial"},
			},
		},
		{
			name: "with error",
			tokenize: func(ctx context.Context, c *tokenize.Client) (*tokenize.Result, error) {
				return c.ForProperty(ctx, "superduperfactorial", tokenize.PropertyOptions{
					CollectionName: "Songs",
					PropetyName:    "title",
				})
			},
			stubs: []testkit.Stub[api.TokenizeRequest, api.TokenizeResponse]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			require.NotNil(t, tt.tokenize, "incomplete test case: missing tokenize function")

			transport := testkit.NewTransport(t, tt.stubs)
			c := tokenize.NewClient(transport)
			require.NotNil(t, c, "nil client")

			r, err := tt.tokenize(t.Context(), c)
			tt.err.Require(t, err, "tokenize error")

			require.Equal(t, tt.want, r, "bad result")
		})
	}
}
