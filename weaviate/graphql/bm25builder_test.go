package graphql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBM25Builder_build(t *testing.T) {
	t.Run("all parameters", func(t *testing.T) {
		bm25 := BM25ArgumentBuilder{}
		str := bm25.WithQuery("query").WithProperties("title", "document", "date").build()
		expected := `bm25:{query: "query", properties: ["title","document","date"]}`
		require.Equal(t, expected, str)
	})

	t.Run("only query", func(t *testing.T) {
		bm25 := BM25ArgumentBuilder{}
		str := bm25.WithQuery("query").build()
		expected := `bm25:{query: "query"}`
		require.Equal(t, expected, str)
	})

	t.Run("query with escaping", func(t *testing.T) {
		bm25 := BM25ArgumentBuilder{}
		str := bm25.WithQuery("\"I'm a complex string\" says the string").build()
		expected := `bm25:{query: "\"I'm a complex string\" says the string"}`
		require.Equal(t, expected, str)
	})
}
