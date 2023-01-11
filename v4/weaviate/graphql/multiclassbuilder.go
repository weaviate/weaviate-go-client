package graphql

import (
	"context"
	"fmt"
	"strings"

	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/filters"
	"github.com/semi-technologies/weaviate/entities/models"
)

// MultiClassBuilder for GraphQL
type MultiClassBuilder struct {
	connection    rest
	classBuilders map[string]*classBuilder
}

type classBuilder struct {
	className  string
	withFields []Field

	includesFilterClause bool // true if brackets behind class is needed
	includesLimit        bool
	limit                int
	includesOffset       bool
	offset               int
	withWhereFilter      *filters.WhereBuilder
	withNearTextFilter   *NearTextArgumentBuilder
	withNearVectorFilter *NearVectorArgumentBuilder
	withNearObjectFilter *NearObjectArgumentBuilder
	withGroupFilter      *GroupArgumentBuilder
	withAskFilter        *AskArgumentBuilder
	withNearImageFilter  *NearImageArgumentBuilder
	withSort             *SortBuilder
	withBM25             *BM25ArgumentBuilder
	withHybrid           *HybridArgumentBuilder
}

// ClassName that should be queried
func NewQueryClassBuilder(className string) *classBuilder {
	return &classBuilder{className: className}
}

func (cb *MultiClassBuilder) AddQueryClass(class *classBuilder) *MultiClassBuilder {
	cb.classBuilders[class.className] = class
	return cb
}

// Do execute the GraphQL query
func (cb *MultiClassBuilder) Do(ctx context.Context) (*models.GraphQLResponse, error) {
	return runGraphQLQuery(ctx, cb.connection, cb.build())
}

// build the GraphQL query string (not needed when Do is executed)
func (cb *MultiClassBuilder) build() string {
	var query string
	for _, class := range cb.classBuilders {
		filterClause := ""
		if class.includesFilterClause {
			filterClause = cb.createFilterClause(class)
		}
		fieldsClause := cb.createFieldsClause(class)
		query += fmt.Sprintf("%v %v {%v}", class.className, filterClause, fieldsClause) + " "
	}
	query = strings.TrimSpace(query)
	query = fmt.Sprintf("{Get {%v}}", query)
	return query
}

func (cb *MultiClassBuilder) createFilterClause(class *classBuilder) string {
	filters := []string{}
	if class.withWhereFilter != nil {
		filters = append(filters, class.withWhereFilter.String())
	}
	if class.withNearTextFilter != nil {
		filters = append(filters, class.withNearTextFilter.build())
	}
	if class.withBM25 != nil {
		filters = append(filters, class.withBM25.build())
	}
	if class.withHybrid != nil {
		filters = append(filters, class.withHybrid.build())
	}
	if class.withNearVectorFilter != nil {
		filters = append(filters, class.withNearVectorFilter.build())
	}
	if class.withNearObjectFilter != nil {
		filters = append(filters, class.withNearObjectFilter.build())
	}
	if class.withAskFilter != nil {
		filters = append(filters, class.withAskFilter.build())
	}
	if class.withNearImageFilter != nil {
		filters = append(filters, class.withNearImageFilter.build())
	}
	if class.withGroupFilter != nil {
		filters = append(filters, class.withGroupFilter.build())
	}
	if class.includesLimit {
		filters = append(filters, fmt.Sprintf("limit: %v", class.limit))
	}
	if class.includesOffset {
		filters = append(filters, fmt.Sprintf("offset: %v", class.offset))
	}
	if class.withSort != nil {
		filters = append(filters, class.withSort.build())
	}
	return fmt.Sprintf("(%s)", strings.Join(filters, ", "))
}

func (cb *MultiClassBuilder) createFieldsClause(class *classBuilder) string {
	if len(class.withFields) > 0 {
		fields := make([]string, len(class.withFields))
		for i := range class.withFields {
			fields[i] = class.withFields[i].build()
		}
		return strings.Join(fields, " ")
	}
	return ""
}

// WithFields included in the result set
func (cb *classBuilder) WithFields(fields ...Field) *classBuilder {
	cb.withFields = fields
	return cb
}

// WithWhere filter
func (cb *classBuilder) WithWhere(where *filters.WhereBuilder) *classBuilder {
	cb.includesFilterClause = true
	cb.withWhereFilter = where
	return cb
}

// WithLimit of objects in the result set
func (cb *classBuilder) WithLimit(limit int) *classBuilder {
	cb.includesFilterClause = true
	cb.includesLimit = true
	cb.limit = limit
	return cb
}

// WithOffset of objects in the result set
func (cb *classBuilder) WithOffset(offset int) *classBuilder {
	cb.includesFilterClause = true
	cb.includesOffset = true
	cb.offset = offset
	return cb
}

// WithBM25 to search the inverted index
func (cb *classBuilder) WithBM25(bm25 *BM25ArgumentBuilder) *classBuilder {
	cb.includesFilterClause = true
	cb.withBM25 = bm25
	return cb
}

// WithHybrid to combine multiple searches
func (cb *classBuilder) WithHybrid(hybrid *HybridArgumentBuilder) *classBuilder {
	cb.includesFilterClause = true
	cb.withHybrid = hybrid
	return cb
}

// WithNearText clause to find close objects
func (cb *classBuilder) WithNearText(nearText *NearTextArgumentBuilder) *classBuilder {
	cb.includesFilterClause = true
	cb.withNearTextFilter = nearText
	return cb
}

// WithNearObject clause to find close objects
func (cb *classBuilder) WithNearImage(nearImage *NearImageArgumentBuilder) *classBuilder {
	cb.includesFilterClause = true
	cb.withNearImageFilter = nearImage
	return cb
}

// WithNearVector clause to find close objects
func (cb *classBuilder) WithNearVector(nearVector *NearVectorArgumentBuilder) *classBuilder {
	cb.includesFilterClause = true
	cb.withNearVectorFilter = nearVector
	return cb
}

// WithGroup statement
func (cb *classBuilder) WithGroup(group *GroupArgumentBuilder) *classBuilder {
	cb.includesFilterClause = true
	cb.withGroupFilter = group
	return cb
}

// WithAsk clause to find an aswer to the question
func (cb *classBuilder) WithAsk(ask *AskArgumentBuilder) *classBuilder {
	cb.includesFilterClause = true
	cb.withAskFilter = ask
	return cb
}

// WithNearObject clause to find close objects
func (cb *classBuilder) WithNearObject(nearObject *NearObjectArgumentBuilder) *classBuilder {
	cb.includesFilterClause = true
	cb.withNearObjectFilter = nearObject
	return cb
}

// WithSort included in the result set
func (cb *classBuilder) WithSort(sort ...Sort) *classBuilder {
	cb.includesFilterClause = true
	cb.withSort = &SortBuilder{sort}
	return cb
}
