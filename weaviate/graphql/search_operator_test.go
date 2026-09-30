package graphql

import (
	"testing"

	"github.com/stretchr/testify/require"
	pb "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

func TestSearchOperatorBuilder(t *testing.T) {
	t.Run("BM25SearchOperatorAnd", func(t *testing.T) {
		var builder BM25SearchOperatorBuilder
		builder.WithOperator(BM25SearchOperatorAnd).WithMinimumMatch(4)

		require.Equal(t, "{operator:And minimumOrTokensMatch:0}", builder.build())
		grpcOpts := builder.togrpc()
		require.NotNil(t, grpcOpts)
		require.Equal(t, pb.SearchOperatorOptions_OPERATOR_AND, grpcOpts.Operator)
		require.Nil(t, grpcOpts.MinimumOrTokensMatch)
	})

	t.Run("BM25SearchOperatorAndCross", func(t *testing.T) {
		var builder BM25SearchOperatorBuilder
		builder.WithOperator(BM25SearchOperatorAndCross).WithMinimumMatch(4)

		require.Equal(t, "{operator:AndCross minimumOrTokensMatch:0}", builder.build())
		grpcOpts := builder.togrpc()
		require.NotNil(t, grpcOpts)
		require.Equal(t, pb.SearchOperatorOptions_OPERATOR_AND_CROSS, grpcOpts.Operator)
		require.Nil(t, grpcOpts.MinimumOrTokensMatch)
	})

	t.Run("BM25SearchOperatorOr", func(t *testing.T) {
		var builder BM25SearchOperatorBuilder
		builder.WithOperator(BM25SearchOperatorOr).WithMinimumMatch(3)

		require.Equal(t, "{operator:Or minimumOrTokensMatch:3}", builder.build())
		grpcOpts := builder.togrpc()
		require.NotNil(t, grpcOpts)
		require.Equal(t, pb.SearchOperatorOptions_OPERATOR_OR, grpcOpts.Operator)
		require.NotNil(t, grpcOpts.MinimumOrTokensMatch)
		require.Equal(t, int32(3), *grpcOpts.MinimumOrTokensMatch)
	})
}
