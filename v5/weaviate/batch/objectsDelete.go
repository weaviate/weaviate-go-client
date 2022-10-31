package batch

import (
	"context"
	"fmt"
	"net/http"

	"github.com/semi-technologies/weaviate-go-client/v5/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/v5/weaviate/except"
	"github.com/semi-technologies/weaviate-go-client/v5/weaviate/filters"
	"github.com/semi-technologies/weaviate/entities/models"
)

type ObjectsBatchDeleter struct {
	connection  *connection.Connection
	className   string
	dryRun      bool
	output      string
	whereFilter *filters.WhereBuilder
}

func (b *ObjectsBatchDeleter) WithClassName(className string) *ObjectsBatchDeleter {
	b.className = className
	return b
}

func (b *ObjectsBatchDeleter) WithDryRun(dryRun bool) *ObjectsBatchDeleter {
	b.dryRun = dryRun
	return b
}

func (b *ObjectsBatchDeleter) WithOutput(output string) *ObjectsBatchDeleter {
	b.output = output
	return b
}

func (b *ObjectsBatchDeleter) WithWhere(whereFilter *filters.WhereBuilder) *ObjectsBatchDeleter {
	b.whereFilter = whereFilter
	return b
}

// Do delete's all the objects which match the builder's filter
func (ob *ObjectsBatchDeleter) Do(ctx context.Context) (*models.BatchDeleteResponse, error) {
	if ob.whereFilter == nil {
		return nil, fmt.Errorf("filter must be set prior to deletion, use WithFilter")
	}

	body := &models.BatchDelete{
		DryRun: &ob.dryRun,
		Output: &ob.output,
		Match: &models.BatchDeleteMatch{
			Class: ob.className,
			Where: ob.whereFilter.Build(),
		},
	}

	responseData, responseErr := ob.connection.RunREST(ctx, "/batch/objects", http.MethodDelete, body)
	err := except.CheckResponseDataErrorAndStatusCode(responseData, responseErr, 200)
	if err != nil {
		return nil, err
	}

	var parsedResponse models.BatchDeleteResponse
	parseErr := responseData.DecodeBodyIntoTarget(&parsedResponse)
	return &parsedResponse, parseErr
}
