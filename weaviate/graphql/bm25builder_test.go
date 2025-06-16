package graphql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBM25Builder_build(t *testing.T) {
	t.Run("all parameters", func(t *testing.T) {
		var bm25 BM25ArgumentBuilder
		got := bm25.WithQuery("query").WithProperties("title", "document", "date").build()
		expected := `bm25:{query: "query", properties: ["title","document","date"]}`
		require.Equal(t, expected, got)
	})

	t.Run("only query", func(t *testing.T) {
		var bm25 BM25ArgumentBuilder
		got := bm25.WithQuery("query").build()
		expected := `bm25:{query: "query"}`
		require.Equal(t, expected, got)
	})

	t.Run("query with escaping", func(t *testing.T) {
		var bm25 BM25ArgumentBuilder
		got := bm25.WithQuery("\"I'm a complex string\" says the string").build()
		expected := `bm25:{query: "\"I'm a complex string\" says the string"}`
		require.Equal(t, expected, got)
	})

	t.Run("query with searchOperator (OR)", func(t *testing.T) {
		var (
			bm25     BM25ArgumentBuilder
			operator BM25SearchOperatorBuilder
		)
		operator.WithOperator(BM25SearchOperatorOr).WithMinimumMatch(4)
		got := bm25.WithQuery("hello").WithSearchOperator(operator).build()
		expected := `bm25:{query: "hello", searchOperator:{operator:Or minimumOrTokensMatch:4}}`
		require.Equal(t, expected, got)
	})

	t.Run("query with searchOperator (AND)", func(t *testing.T) {
		var (
			bm25     BM25ArgumentBuilder
			operator BM25SearchOperatorBuilder
		)
		operator.WithOperator(BM25SearchOperatorAnd).WithMinimumMatch(4)
		got := bm25.WithQuery("hello").WithSearchOperator(operator).build()
		expected := `bm25:{query: "hello", searchOperator:{operator:And minimumOrTokensMatch:0}}`
		require.Equal(t, expected, got)
	})
}
