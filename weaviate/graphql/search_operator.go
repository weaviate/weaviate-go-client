package graphql

import "fmt"

const (
	BM25SearchOperatorAnd = "And"
	BM25SearchOperatorOr  = "Or"
)

type BM25SearchOperatorBuilder struct {
	operator     string
	minimumMatch int
}

func (bm25 *BM25SearchOperatorBuilder) WithOperator(operator string) *BM25SearchOperatorBuilder {
	bm25.operator = operator
	return bm25
}

// WithMinimumMatch is only relevant for BM25SearchOperatorOr operator.
func (bm25 *BM25SearchOperatorBuilder) WithMinimumMatch(times int) *BM25SearchOperatorBuilder {
	if bm25.operator != BM25SearchOperatorAnd {
		bm25.minimumMatch = times
	}
	return bm25
}

func (bm25 *BM25SearchOperatorBuilder) build() string {
	return fmt.Sprintf("{operator:%s minimumOrTokensMatch:%d}", bm25.operator, bm25.minimumMatch)
}
