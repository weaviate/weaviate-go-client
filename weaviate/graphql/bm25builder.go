package graphql

import (
	"encoding/json"
	"fmt"
	"strings"

	pb "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

type BM25ArgumentBuilder struct {
	query          string
	properties     []string
	searchOperator *BM25SearchOperatorBuilder
}

// WithQuery the search string
func (b *BM25ArgumentBuilder) WithQuery(query string) *BM25ArgumentBuilder {
	b.query = query
	return b
}

func (b *BM25ArgumentBuilder) WithSearchOperator(searchOperator BM25SearchOperatorBuilder) *BM25ArgumentBuilder {
	b.searchOperator = &searchOperator
	return b
}

// WithProperties the properties to search. Leave blank for all
func (b *BM25ArgumentBuilder) WithProperties(properties ...string) *BM25ArgumentBuilder {
	b.properties = properties
	return b
}

// Build build the given clause
func (b *BM25ArgumentBuilder) build() string {
	clause := []string{}
	if b.query != "" {
		clause = append(clause, fmt.Sprintf("query: %q", b.query))
	}
	if len(b.properties) > 0 {
		propStr, err := json.Marshal(b.properties)
		if err != nil {
			panic(fmt.Errorf("failed to unmarshal bm25 properties: %s", err))
		}
		clause = append(clause, fmt.Sprintf("properties: %v", string(propStr)))
	}
	if b.searchOperator != nil {
		clause = append(clause, fmt.Sprintf("searchOperator:%s", b.searchOperator.build()))
	}
	return fmt.Sprintf("bm25:{%v}", strings.Join(clause, ", "))
}

func (b *BM25ArgumentBuilder) togrpc() *pb.BM25 {
	bm25 := &pb.BM25{
		Query:      b.query,
		Properties: b.properties,
	}
	return bm25
}
