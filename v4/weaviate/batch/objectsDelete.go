package batch

import (
	"context"
	"fmt"
	"net/http"

	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/except"
	"github.com/semi-technologies/weaviate/entities/models"
)

// TODO: delete when batch deletion merged to weaviate master
// 		 then use *models.BatchDelete
type BatchDeleteFilter struct {
	DryRun *bool             `json:"dryRun,omitempty"`
	Match  *BatchDeleteMatch `json:"match"`
	Output *string           `json:"output,omitempty"`
}

// TODO: delete when batch deletion merged to weaviate master
type BatchDeleteMatch struct {
	Class string              `json:"class"`
	Where *models.WhereFilter `json:"where"`
}

// TODO: delete when batch deletion merged to weaviate master
type BatchDeleteResponse struct {
	Match   BatchDeleteMatch   `json:"match"`
	DryRun  bool               `json:"dryRun"`
	Output  string             `json:"output"`
	Results BatchDeleteResults `json:"results"`
}

// TODO: delete when batch deletion merged to weaviate master
type BatchDeleteResults struct {
	Matches    int64                     `json:"matches"`
	Limit      int64                     `json:"limit"`
	Failed     int64                     `json:"failed"`
	Successful int64                     `json:"successful"`
	Objects    []BatchDeleteResultObject `json:"objects"`
}

// TODO: delete when batch deletion merged to weaviate master
type BatchDeleteResultObject struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Errors *struct {
		Error []*BatchDeleteResultObjectError `json:"error,omitempty"`
	} `json:"errors,omitempty"`
}

// TODO: delete when batch deletion merged to weaviate master
type BatchDeleteResultObjectError struct {
	Message string `json:"message,omitempty"`
}

type ObjectsBatchDeleter struct {
	connection *connection.Connection
	filter     *BatchDeleteFilter
}

// WithFilter sets the filter for matching which objects to delete
func (ob *ObjectsBatchDeleter) WithFilter(filter *BatchDeleteFilter) *ObjectsBatchDeleter {
	ob.filter = filter
	return ob
}

// Do delete's all the objects which match the builder's filter
func (ob *ObjectsBatchDeleter) Do(ctx context.Context) (*BatchDeleteResponse, error) {
	if ob.filter == nil {
		return nil, fmt.Errorf("filter must be set prior to deletion, use WithFilter")
	}

	responseData, responseErr := ob.connection.RunREST(ctx, "/batch/objects", http.MethodDelete, ob.filter)
	err := except.CheckResponseDataErrorAndStatusCode(responseData, responseErr, 200)
	if err != nil {
		return nil, err
	}

	var parsedResponse BatchDeleteResponse
	parseErr := responseData.DecodeBodyIntoTarget(&parsedResponse)
	return &parsedResponse, parseErr
}
