package graphql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHybridBuilder_build(t *testing.T) {

	hybrid := HybridArgumentBuilder{}
	str := hybrid.WithQuery("query").WithVector(1, 2, 3).WithAlpha(0.6).build()
	expected := `hybrid:{query: "query", vector: [1, 2, 3], alpha: 0.6}`
	require.Equal(t, expected, str)
}
