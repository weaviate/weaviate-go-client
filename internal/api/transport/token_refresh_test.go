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
			// We expect this to exit immediately, no need for a goroutine.
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

		c := make(chan time.Time)
		defer close(c)

		// tick asserts that the caller tried to set the timer for the right duration,
		// then returns a channel that can be read from immediately.
		tick := func(d time.Duration) <-chan time.Time {
			assert.Equal(t, time.Duration(92)*time.Second, d, "must try to sleep for %ds", 92)

			// time.After(0) and the underlying time.NewTicker(0) panic on 0 input.
			return c
		}

		// Act
		go tokenKeepalive(ctx, &src, tick)

		// Wait until tokenKeepaline has received from this channel at least once.
		// Then cancel the context.
		c <- time.Now()
		cancel()

		// Assert
		assert.GreaterOrEqual(t, src.used, 1, "expect src.Token() to be used in the background")

		src.used = 0
		assert.Zero(t, src.used, "no src.Token() after context is canceled")
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
