package graphql

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/filters"
	"github.com/weaviate/weaviate/entities/models"
)

// AggregateBuilder for the aggregate GraphQL query string
type AggregateBuilder struct {
	connection                rest
	fields                    []Field
	className                 string
	includesFilterClause      bool // true if brackets behind class is needed
	groupByClausePropertyName string
	tenant                    string
	withWhereFilter           *filters.WhereBuilder
	withNearVectorFilter      *NearVectorArgumentBuilder
	withNearObjectFilter      *NearObjectArgumentBuilder
	withNearTextFilter        *NearTextArgumentBuilder
	withAsk                   *AskArgumentBuilder
	withNearImage             *NearImageArgumentBuilder
	withNearAudio             *NearAudioArgumentBuilder
	withNearVideo             *NearVideoArgumentBuilder
	withNearDepth             *NearDepthArgumentBuilder
	withNearThermal           *NearThermalArgumentBuilder
	withNearImu               *NearImuArgumentBuilder
	withHybrid                *HybridArgumentBuilder
	includesObjectLimit       bool
	objectLimit               int
	includesLimit             bool
	limit                     int
}

// WithFields that should be included in the aggregation query e.g. `meta{count}`
func (ab *AggregateBuilder) WithFields(fields ...Field) *AggregateBuilder {
	ab.fields = fields
	return ab
}

// WithClassName that should be aggregated
func (ab *AggregateBuilder) WithClassName(name string) *AggregateBuilder {
	ab.className = name
	return ab
}

// WithWhere adds the where filter.
func (ab *AggregateBuilder) WithWhere(where *filters.WhereBuilder) *AggregateBuilder {
	ab.includesFilterClause = true
	ab.withWhereFilter = where
	return ab
}

// WithGroupBy adds the group by property clause as the filter.
//
//	The group by value/path clause still needs to be set in the WithFields field.
func (ab *AggregateBuilder) WithGroupBy(propertyName string) *AggregateBuilder {
	ab.includesFilterClause = true
	ab.groupByClausePropertyName = propertyName
	return ab
}

// WithNearText clause to find close objects
func (ab *AggregateBuilder) WithNearText(nearText *NearTextArgumentBuilder) *AggregateBuilder {
	ab.includesFilterClause = true
	ab.withNearTextFilter = nearText
	return ab
}

// WithNearObject clause to find close objects
func (ab *AggregateBuilder) WithNearObject(nearObject *NearObjectArgumentBuilder) *AggregateBuilder {
	ab.includesFilterClause = true
	ab.withNearObjectFilter = nearObject
	return ab
}

// WithNearVector clause to find close objects
func (ab *AggregateBuilder) WithNearVector(nearVector *NearVectorArgumentBuilder) *AggregateBuilder {
	ab.includesFilterClause = true
	ab.withNearVectorFilter = nearVector
	return ab
}

// WithObjectLimit specifies max number of vector search results to return
func (ab *AggregateBuilder) WithObjectLimit(objectLimit int) *AggregateBuilder {
	ab.objectLimit = objectLimit
	ab.includesObjectLimit = true
	return ab
}

// WithLimit specifies limit to group by argument
func (ab *AggregateBuilder) WithLimit(limit int) *AggregateBuilder {
	ab.limit = limit
	ab.includesLimit = true
	return ab
}

// WithAsk adds ask to clause
func (ab *AggregateBuilder) WithAsk(ask *AskArgumentBuilder) *AggregateBuilder {
	ab.includesFilterClause = true
	ab.withAsk = ask
	return ab
}

// WithNearImage adds nearImage to clause
func (ab *AggregateBuilder) WithNearImage(nearImage *NearImageArgumentBuilder) *AggregateBuilder {
	ab.includesFilterClause = true
	ab.withNearImage = nearImage
	return ab
}

// WithNearAudio adds nearAudio to clause
func (ab *AggregateBuilder) WithNearAudio(nearAudio *NearAudioArgumentBuilder) *AggregateBuilder {
	ab.includesFilterClause = true
	ab.withNearAudio = nearAudio
	return ab
}

// WithNearVideo adds nearVideo to clause
func (ab *AggregateBuilder) WithNearVideo(nearVideo *NearVideoArgumentBuilder) *AggregateBuilder {
	ab.includesFilterClause = true
	ab.withNearVideo = nearVideo
	return ab
}

// WithNearDepth adds nearDepth to clause
func (ab *AggregateBuilder) WithNearDepth(nearDepth *NearDepthArgumentBuilder) *AggregateBuilder {
	ab.includesFilterClause = true
	ab.withNearDepth = nearDepth
	return ab
}

// WithNearThermal adds nearThermal to clause
func (ab *AggregateBuilder) WithNearThermal(nearThermal *NearThermalArgumentBuilder) *AggregateBuilder {
	ab.includesFilterClause = true
	ab.withNearThermal = nearThermal
	return ab
}

// WithNearImu adds nearIMU to clause
func (ab *AggregateBuilder) WithNearImu(nearImu *NearImuArgumentBuilder) *AggregateBuilder {
	ab.includesFilterClause = true
	ab.withNearImu = nearImu
	return ab
}

// WithHybrid to combine multiple searches
func (ab *AggregateBuilder) WithHybrid(hybrid *HybridArgumentBuilder) *AggregateBuilder {
	ab.includesFilterClause = true
	ab.withHybrid = hybrid
	return ab
}

// WithTenant to indicate which tenant aggregated objects belong to
func (ab *AggregateBuilder) WithTenant(tenant string) *AggregateBuilder {
	ab.includesFilterClause = true
	ab.tenant = tenant
	return ab
}

// Do execute the aggregation query
func (ab *AggregateBuilder) Do(ctx context.Context) (*models.GraphQLResponse, error) {
	return runGraphQLQuery(ctx, ab.connection, ab.build())
}

func (ab *AggregateBuilder) createFilterClause() string {
	if ab.includesFilterClause {
		filters := []string{}
		if ab.tenant != "" {
			filters = append(filters, fmt.Sprintf("tenant: %q", ab.tenant))
		}
		if len(ab.groupByClausePropertyName) > 0 {
			filters = append(filters, fmt.Sprintf(`groupBy: "%v"`, ab.groupByClausePropertyName))
		}
		if ab.withWhereFilter != nil {
			filters = append(filters, ab.withWhereFilter.String())
		}
		for _, b := range []argumentBuilder{
			ab.withAsk, ab.withNearTextFilter, ab.withNearObjectFilter, ab.withNearVectorFilter, ab.withNearImage,
			ab.withNearAudio, ab.withNearVideo, ab.withNearDepth, ab.withNearThermal, ab.withNearImu, ab.withHybrid,
		} {
			bVal := reflect.ValueOf(b)
			if bVal.Kind() == reflect.Ptr && !bVal.IsNil() {
				filters = append(filters, b.build())
			}
		}
		if ab.includesObjectLimit {
			filters = append(filters, fmt.Sprintf("objectLimit: %d", ab.objectLimit))
		}
		if ab.includesLimit {
			filters = append(filters, fmt.Sprintf("limit: %d", ab.limit))
		}
		return fmt.Sprintf("(%s)", strings.Join(filters, ", "))
	}
	return ""
}

func (ab *AggregateBuilder) createFieldsClause() string {
	if len(ab.fields) > 0 {
		fields := make([]string, len(ab.fields))
		for i := range ab.fields {
			fields[i] = ab.fields[i].build()
		}
		return strings.Join(fields, " ")
	}
	return ""
}

// build the query string
func (ab *AggregateBuilder) build() string {
	filterClause := ab.createFilterClause()
	fields := ab.createFieldsClause()
	return fmt.Sprintf(`{Aggregate{%v%v{%v}}}`, ab.className, filterClause, fields)
}
