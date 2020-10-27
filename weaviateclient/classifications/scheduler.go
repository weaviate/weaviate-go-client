package classifications

import (
	"context"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/clienterrors"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/paragons"
	"github.com/semi-technologies/weaviate/entities/models"
	"net/http"
)

// Scheduler builder to schedule a classification
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

// WithType of classification e.g. knn or contextual
func (s *Scheduler) WithType(classificationType paragons.Classification) *Scheduler {
	s.classificationType = classificationType
	return s
}

// WithClassName that should be classified
func (s *Scheduler) WithClassName(name string) *Scheduler {
	s.withClassName = name
	return s
}

// WithClassifyProperties defines the properties that will be labeled through the classification
func (s *Scheduler) WithClassifyProperties(classifyProperties []string) *Scheduler {
	s.withClassifyProperties = classifyProperties
	return s
}

// WithBasedOnProperties defines the properties that will be considered for the classification
func (s *Scheduler) WithBasedOnProperties(basedOnProperties []string) *Scheduler {
	s.withBasedOnProperties = basedOnProperties
	return s
}

// WithSourceWhereFilter filter the data objects to be labeled
func (s *Scheduler) WithSourceWhereFilter(whereFilter *models.WhereFilter) *Scheduler {
	s.withSourceWhereFilter = whereFilter
	return s
}

// WithTrainingSetWhereFilter filter the objects that are used as training data. E.g. in a knn classification
func (s *Scheduler) WithTrainingSetWhereFilter(whereFilter *models.WhereFilter) *Scheduler {
	s.withTrainingSetWhereFilter = whereFilter
	return s
}

// WithTargetWhereFilter filter the label objects
func (s *Scheduler) WithTargetWhereFilter(whereFilter *models.WhereFilter) *Scheduler {
	s.withTargetWhereFilter = whereFilter
	return s
}

// WithK set the number of neighbours considered by a knn classification
func (s *Scheduler) WithK(k int32) *Scheduler {
	s.withK = k
	return s
}

// Do schedule the classification in weaviate
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