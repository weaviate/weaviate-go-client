package graphql

import (
	"fmt"

	pb "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

const (
	BM25SearchOperatorAnd = "And"
	BM25SearchOperatorOr  = "Or"
)

type BM25SearchOperatorBuilder struct {
	operator     string
	minimumMatch int32
}

func (bm25 *BM25SearchOperatorBuilder) WithOperator(operator string) *BM25SearchOperatorBuilder {
	bm25.operator = operator
	return bm25
}

// WithMinimumMatch is only relevant for BM25SearchOperatorOr operator.
func (bm25 *BM25SearchOperatorBuilder) WithMinimumMatch(times int) *BM25SearchOperatorBuilder {
	if bm25.operator != BM25SearchOperatorAnd {
		bm25.minimumMatch = int32(times)
	}
	return bm25
}

func (bm25 *BM25SearchOperatorBuilder) build() string {
	return fmt.Sprintf("{operator:%s minimumOrTokensMatch:%d}", bm25.operator, bm25.minimumMatch)
}

func (bm25 *BM25SearchOperatorBuilder) togrpc() *pb.SearchOperatorOptions {
	switch bm25.operator {
	case BM25SearchOperatorAnd:
		return &pb.SearchOperatorOptions{Operator: pb.SearchOperatorOptions_OPERATOR_AND}
	case BM25SearchOperatorOr:
		return &pb.SearchOperatorOptions{Operator: pb.SearchOperatorOptions_OPERATOR_OR, MinimumOrTokensMatch: &bm25.minimumMatch}
	default:
		return nil
	}
}
