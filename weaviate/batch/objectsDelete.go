package batch

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/pathbuilder"
	"github.com/weaviate/weaviate/entities/models"
)

type ObjectsBatchDeleter struct {
	connection       *connection.Connection
	className        string
	dryRun           *bool
	output           *string
	whereFilter      *filters.WhereBuilder
	consistencyLevel string
	tenant           string
}

func (b *ObjectsBatchDeleter) WithClassName(className string) *ObjectsBatchDeleter {
	b.className = className
	return b
}

func (b *ObjectsBatchDeleter) WithDryRun(dryRun bool) *ObjectsBatchDeleter {
	b.dryRun = &dryRun
	return b
}

func (b *ObjectsBatchDeleter) WithOutput(output string) *ObjectsBatchDeleter {
	b.output = &output
	return b
}

func (b *ObjectsBatchDeleter) WithWhere(whereFilter *filters.WhereBuilder) *ObjectsBatchDeleter {
	b.whereFilter = whereFilter
	return b
}

// WithConsistencyLevel determines how many replicas must acknowledge a request
// before it is considered successful. Mutually exclusive with node_name param.
// Can be one of 'ALL', 'ONE', or 'QUORUM'.
func (b *ObjectsBatchDeleter) WithConsistencyLevel(cl string) *ObjectsBatchDeleter {
	b.consistencyLevel = cl
	return b
}

// WithTenant sets tenant, objects should be deleted from
func (b *ObjectsBatchDeleter) WithTenant(tenant string) *ObjectsBatchDeleter {
	b.tenant = tenant
	return b
}

// Do delete's all the objects which match the builder's filter
func (ob *ObjectsBatchDeleter) Do(ctx context.Context) (*models.BatchDeleteResponse, error) {
	if ob.whereFilter == nil {
		return nil, fmt.Errorf("filter must be set prior to deletion, use WithWhere")
	}

	body := &models.BatchDelete{
		DryRun: ob.dryRun,
		Output: ob.output,
		Match: &models.BatchDeleteMatch{
			Class: ob.className,
			Where: ob.whereFilter.Build(),
		},
	}

	path := pathbuilder.BatchObjects(pathbuilder.Components{
		ConsistencyLevel: ob.consistencyLevel,
		Tenant:           ob.tenant,
	})
	responseData, responseErr := ob.connection.RunREST(ctx, path, http.MethodDelete, body)
	err := except.CheckResponseDataErrorAndStatusCode(responseData, responseErr, 200)
	if err != nil {
		return nil, err
	}

	var parsedResponse models.BatchDeleteResponse
	parseErr := responseData.DecodeBodyIntoTarget(&parsedResponse)
	return &parsedResponse, parseErr
}
