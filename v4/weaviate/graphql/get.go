package graphql

import (
	"context"
	"fmt"
	"strings"

	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/filters"
	"github.com/semi-technologies/weaviate/entities/models"
)

// GetBuilder for GraphQL
type GetBuilder struct {
	connection    rest
	ClassBuilders map[string]*ClassBuilder
}

type ClassBuilder struct {
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
func NewQueryClassBuilder(className string) *ClassBuilder {
	return &ClassBuilder{className: className}
}

// WithFields included in the result set
func (cb *ClassBuilder) WithFields(fields ...Field) *ClassBuilder {
	cb.withFields = fields
	return cb
}

// WithWhere filter
func (cb *ClassBuilder) WithWhere(where *filters.WhereBuilder) *ClassBuilder {
	cb.includesFilterClause = true
	cb.withWhereFilter = where
	return cb
}

// WithLimit of objects in the result set
func (cb *ClassBuilder) WithLimit(limit int) *ClassBuilder {
	cb.includesFilterClause = true
	cb.includesLimit = true
	cb.limit = limit
	return cb
}

// WithOffset of objects in the result set
func (cb *ClassBuilder) WithOffset(offset int) *ClassBuilder {
	cb.includesFilterClause = true
	cb.includesOffset = true
	cb.offset = offset
	return cb
}

// WithBM25 to search the inverted index
func (cb *ClassBuilder) WithBM25(bm25 *BM25ArgumentBuilder) *ClassBuilder {
	cb.includesFilterClause = true
	cb.withBM25 = bm25
	return cb
}

// WithHybrid to combine multiple searches
func (cb *ClassBuilder) WithHybrid(hybrid *HybridArgumentBuilder) *ClassBuilder {
	cb.includesFilterClause = true
	cb.withHybrid = hybrid
	return cb
}

// WithNearText clause to find close objects
func (cb *ClassBuilder) WithNearText(nearText *NearTextArgumentBuilder) *ClassBuilder {
	cb.includesFilterClause = true
	cb.withNearTextFilter = nearText
	return cb
}

// WithNearObject clause to find close objects
func (cb *ClassBuilder) WithNearImage(nearImage *NearImageArgumentBuilder) *ClassBuilder {
	cb.includesFilterClause = true
	cb.withNearImageFilter = nearImage
	return cb
}

// WithNearVector clause to find close objects
func (cb *ClassBuilder) WithNearVector(nearVector *NearVectorArgumentBuilder) *ClassBuilder {
	cb.includesFilterClause = true
	cb.withNearVectorFilter = nearVector
	return cb
}

// WithGroup statement
func (cb *ClassBuilder) WithGroup(group *GroupArgumentBuilder) *ClassBuilder {
	cb.includesFilterClause = true
	cb.withGroupFilter = group
	return cb
}

// WithAsk clause to find an aswer to the question
func (cb *ClassBuilder) WithAsk(ask *AskArgumentBuilder) *ClassBuilder {
	cb.includesFilterClause = true
	cb.withAskFilter = ask
	return cb
}

// WithNearObject clause to find close objects
func (cb *ClassBuilder) WithNearObject(nearObject *NearObjectArgumentBuilder) *ClassBuilder {
	cb.includesFilterClause = true
	cb.withNearObjectFilter = nearObject
	return cb
}

// WithSort included in the result set
func (cb *ClassBuilder) WithSort(sort ...Sort) *ClassBuilder {
	cb.includesFilterClause = true
	cb.withSort = &SortBuilder{sort}
	return cb
}

// Do execute the GraphQL query
func (gb *GetBuilder) Do(ctx context.Context) (*models.GraphQLResponse, error) {
	return runGraphQLQuery(ctx, gb.connection, gb.build())
}

func (gb *GetBuilder) AddQueryClass(class *ClassBuilder) *GetBuilder {
	gb.ClassBuilders[class.className] = class
	return gb
}

// build the GraphQL query string (not needed when Do is executed)
func (gb *GetBuilder) build() string {
	var query string
	for _, class := range gb.ClassBuilders {
		filterClause := ""
		if class.includesFilterClause {
			filterClause = gb.createFilterClause(class)
		}
		fieldsClause := gb.createFieldsClause(class)
		query = query + fmt.Sprintf("%v %v {%v}", class.className, filterClause, fieldsClause) + " "
	}
	query = strings.TrimSpace(query)
	query = fmt.Sprintf("{Get {%v}}", query)
	return query
}

func (gb *GetBuilder) createFilterClause(class *ClassBuilder) string {
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

func (gb *GetBuilder) createFieldsClause(class *ClassBuilder) string {
	if len(class.withFields) > 0 {
		fields := make([]string, len(class.withFields))
		for i := range class.withFields {
			fields[i] = class.withFields[i].build()
		}
		return strings.Join(fields, " ")
	}
	return ""
}
