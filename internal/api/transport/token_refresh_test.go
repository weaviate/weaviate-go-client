package transport

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"golang.org/x/oauth2"
)

func TestTokenKeepalive(t *testing.T) {
	t.Run("nil token source", func(t *testing.T) {
		require.NotPanics(t, func() {
			tokenKeepalive(t.Context(), nil, time.After)
		})
	})

	t.Run("refreshes after expiry", func(t *testing.T) {
		// Arrange
		ctx, cancel := context.WithCancel(t.Context())
		t.Cleanup(cancel)

		src := tokenSource{tok: oauth2.Token{
			ExpiresIn: 92,
			Expiry:    testkit.Now.Add(92 * time.Second),
		}}

		// Act
		tokenKeepalive(ctx, &src, func(d time.Duration) <-chan time.Time {
			assert.Equal(t, time.Duration(92)*time.Second, d, "must try to sleep for %ds", 92)
			return time.After(5 * time.Millisecond)
		})

		time.Sleep(7 * time.Millisecond)
		cancel()

		// Assert
		require.Equal(t, 2, src.used, "expect src.Token() to be used twice")

		src.used = 0
		time.Sleep(5 * time.Millisecond)
		require.Zero(t, src.used, "no src.Token() after context is canceled")
	})
}

// tokenSource is a fake [oauth2.Token] that always returns the same token.
// Similar to [oauth2.StaticTokenSource], but with configurable [Token.ExpiresIn].
type tokenSource struct {
	tok  oauth2.Token
	used int
}

func (src *tokenSource) Token() (*oauth2.Token, error) {
	src.used++
	return (*oauth2.Token)(&src.tok), nil
}
