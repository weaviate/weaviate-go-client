package graphql

import (
	"context"
	"fmt"
	"strings"

	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/filters"
	"github.com/semi-technologies/weaviate/entities/models"
)

type MultiClassBuilder struct {
	connection    rest
	classBuilders map[string]*builderBase
}

// ClassName that should be queried
func NewQueryClassBuilder(className string) *builderBase {
	return &builderBase{className: className}
}

func (mb *MultiClassBuilder) AddQueryClass(class *builderBase) *MultiClassBuilder {
	mb.classBuilders[class.className] = class
	return mb
}

// Do execute the GraphQL query
func (mb *MultiClassBuilder) Do(ctx context.Context) (*models.GraphQLResponse, error) {
	return runGraphQLQuery(ctx, mb.connection, mb.build())
}

// build the GraphQL query string (not needed when Do is executed)
func (mb *MultiClassBuilder) build() string {
	var query string
	for _, class := range mb.classBuilders {
		filterClause := ""
		if class.includesFilterClause {
			filterClause = mb.createFilterClause(class)
		}
		fieldsClause := mb.createFieldsClause(class)
		query += fmt.Sprintf("%v %v {%v}", class.className, filterClause, fieldsClause) + " "
	}
	query = strings.TrimSpace(query)
	query = fmt.Sprintf("{Get {%v}}", query)
	return query
}

func (mb *MultiClassBuilder) createFilterClause(class *builderBase) string {
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

func (mb *MultiClassBuilder) createFieldsClause(class *builderBase) string {
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
func (bb *builderBase) WithFields(fields ...Field) *builderBase {
	bb.withFields = fields
	return bb
}

// WithWhere filter
func (bb *builderBase) WithWhere(where *filters.WhereBuilder) *builderBase {
	bb.includesFilterClause = true
	bb.withWhereFilter = where
	return bb
}

// WithLimit of objects in the result set
func (bb *builderBase) WithLimit(limit int) *builderBase {
	bb.includesFilterClause = true
	bb.includesLimit = true
	bb.limit = limit
	return bb
}

// WithOffset of objects in the result set
func (bb *builderBase) WithOffset(offset int) *builderBase {
	bb.includesFilterClause = true
	bb.includesOffset = true
	bb.offset = offset
	return bb
}

// WithBM25 to search the inverted index
func (bb *builderBase) WithBM25(bm25 *BM25ArgumentBuilder) *builderBase {
	bb.includesFilterClause = true
	bb.withBM25 = bm25
	return bb
}

// WithHybrid to combine multiple searches
func (bb *builderBase) WithHybrid(hybrid *HybridArgumentBuilder) *builderBase {
	bb.includesFilterClause = true
	bb.withHybrid = hybrid
	return bb
}

// WithNearText clause to find close objects
func (bb *builderBase) WithNearText(nearText *NearTextArgumentBuilder) *builderBase {
	bb.includesFilterClause = true
	bb.withNearTextFilter = nearText
	return bb
}

// WithNearObject clause to find close objects
func (bb *builderBase) WithNearImage(nearImage *NearImageArgumentBuilder) *builderBase {
	bb.includesFilterClause = true
	bb.withNearImageFilter = nearImage
	return bb
}

// WithNearVector clause to find close objects
func (bb *builderBase) WithNearVector(nearVector *NearVectorArgumentBuilder) *builderBase {
	bb.includesFilterClause = true
	bb.withNearVectorFilter = nearVector
	return bb
}

// WithGroup statement
func (bb *builderBase) WithGroup(group *GroupArgumentBuilder) *builderBase {
	bb.includesFilterClause = true
	bb.withGroupFilter = group
	return bb
}

// WithAsk clause to find an aswer to the question
func (bb *builderBase) WithAsk(ask *AskArgumentBuilder) *builderBase {
	bb.includesFilterClause = true
	bb.withAskFilter = ask
	return bb
}

// WithNearObject clause to find close objects
func (bb *builderBase) WithNearObject(nearObject *NearObjectArgumentBuilder) *builderBase {
	bb.includesFilterClause = true
	bb.withNearObjectFilter = nearObject
	return bb
}

// WithSort included in the result set
func (bb *builderBase) WithSort(sort ...Sort) *builderBase {
	bb.includesFilterClause = true
	bb.withSort = &SortBuilder{sort}
	return bb
}
