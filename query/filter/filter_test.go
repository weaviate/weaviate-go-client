package filter_test

import (
	"testing"
	"time"

	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/query/filter"
)

func TestFilter(t *testing.T) {
	filter.Property[int]("price").Equal(92)
	filter.Property[string]("title").Like("%box")
	filter.Property[string]("tags").ContainsAll("foo", "bar")

	var exprs []filter.Expr

	exprs = append(exprs, filter.And{
		filter.UUID.Equal(testkit.UUID),
	}, filter.And{
		filter.CreatedAt.LessThanEqual(testkit.Now.Add(-time.Hour)),
		filter.LastUpdatedAt.GreaterThan(testkit.Now),
	})

	_ = filter.Or(exprs)
}
