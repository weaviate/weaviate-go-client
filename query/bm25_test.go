package query_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weaviate/weaviate-go-client/v6/query"
)

func TestKeywordSimilarity(t *testing.T) {
	t.Run("all tokens match", func(t *testing.T) {
		assert.True(t, query.AllTokensMatch.AllTokensMatch(), "all tokens match")
		assert.Nil(t, query.AllTokensMatch.MinimumTokensMatch(), "minimum tokens match")
	})
	t.Run("minimum tokens match", func(t *testing.T) {
		similarity := query.MinimumTokensMatch(5)
		assert.False(t, similarity.AllTokensMatch(), "all tokens match")
		assert.NotNil(t, similarity.MinimumTokensMatch(), "minimum tokens match")
	})
}
