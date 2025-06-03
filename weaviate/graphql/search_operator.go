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
	bm25.minimumMatch = times
	return bm25
}

func (bm25 *BM25SearchOperatorBuilder) build() string {
	query := fmt.Sprintf("{operator:%s", bm25.operator)
	if bm25.operator != "And" {
		query += fmt.Sprintf(" minimumOrTokensMatch:%d", bm25.minimumMatch)
	}
	query += "}"
	return query
}
