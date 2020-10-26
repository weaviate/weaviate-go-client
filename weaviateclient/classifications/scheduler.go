package classifications

import (
	"context"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/clienterrors"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/paragons"
	"github.com/semi-technologies/weaviate/entities/models"
	"net/http"
)

type Scheduler struct {
	connection *connection.Connection
	classificationType paragons.Classification
	withClassName string
	withClassifyProperties []string
	withBasedOnProperties []string
	withK int32
	withSourceWhereFilter *models.WhereFilter
	withTrainingSetWhereFilter *models.WhereFilter
	withTargetWhereFilter *models.WhereFilter
}

func (s *Scheduler) WithType(classificationType paragons.Classification) *Scheduler {
	s.classificationType = classificationType
	return s
}

func (s *Scheduler) WithClassName(name string) *Scheduler {
	s.withClassName = name
	return s
}

func (s *Scheduler) WithClassifyProperties(classifyProperties []string) *Scheduler {
	s.withClassifyProperties = classifyProperties
	return s
}

func (s *Scheduler) WithBasedOnProperties(basedOnProperties []string) *Scheduler {
	s.withBasedOnProperties = basedOnProperties
	return s
}

func (s *Scheduler) WithSourceWhereFilter(whereFilter *models.WhereFilter) *Scheduler {
	s.withSourceWhereFilter = whereFilter
	return s
}

func (s *Scheduler) WithTrainingSetWhereFilter(whereFilter *models.WhereFilter) *Scheduler {
	s.withTrainingSetWhereFilter = whereFilter
	return s
}

func (s *Scheduler) WithTargetWhereFilter(whereFilter *models.WhereFilter) *Scheduler {
	s.withTargetWhereFilter = whereFilter
	return s
}

func (s *Scheduler) Do(ctx context.Context) (*models.Classification, error) {
	classType := string(s.classificationType)
	config := models.Classification{
		BasedOnProperties:               s.withBasedOnProperties,
		Class:                           s.withClassName,
		ClassifyProperties:              s.withClassifyProperties,
		SourceWhere:                     s.withSourceWhereFilter,
		TargetWhere:                     s.withTargetWhereFilter,
		TrainingSetWhere:                s.withTrainingSetWhereFilter,
		Type:                            &classType,
	}
	if s.classificationType == paragons.KNN {
		config.K = &s.withK
	}
	responseData, responseErr := s.connection.RunREST(ctx, "/classifications", http.MethodPost, config)
	err := clienterrors.CheckResponnseDataErrorAndStatusCode(responseData, responseErr, 201)
	if err != nil {
		return nil, err
	}

	var responseClassification models.Classification
	parseErr := responseData.DecodeBodyIntoTarget(&responseClassification)
	return &responseClassification, parseErr
}