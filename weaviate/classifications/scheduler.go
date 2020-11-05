package classifications

import (
	"context"
	"github.com/go-openapi/strfmt"
	"github.com/semi-technologies/weaviate-go-client/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviate/except"
	"github.com/semi-technologies/weaviate/entities/models"
	"net/http"
	"time"
)

// Scheduler builder to schedule a classification
type Scheduler struct {
	connection                 *connection.Connection
	classificationType         Type
	withClassName              string
	withClassifyProperties     []string
	withBasedOnProperties      []string
	withK                      int32
	withSourceWhereFilter      *models.WhereFilter
	withTrainingSetWhereFilter *models.WhereFilter
	withTargetWhereFilter      *models.WhereFilter
	withWaitForCompletion bool
}

// WithType of classification e.g. knn or contextual
func (s *Scheduler) WithType(classificationType Type) *Scheduler {
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

// WithWaitForCompletion block while classification is running (until classification succeeded or failed)
func (s *Scheduler) WithWaitForCompletion() *Scheduler {
	s.withWaitForCompletion = true
	return s
}

// Do schedule the classification in weaviate
func (s *Scheduler) Do(ctx context.Context) (*models.Classification, error) {
	classType := string(s.classificationType)
	config := models.Classification{
		BasedOnProperties:  s.withBasedOnProperties,
		Class:              s.withClassName,
		ClassifyProperties: s.withClassifyProperties,
		SourceWhere:        s.withSourceWhereFilter,
		TargetWhere:        s.withTargetWhereFilter,
		TrainingSetWhere:   s.withTrainingSetWhereFilter,
		Type:               &classType,
	}
	if s.classificationType == KNN {
		config.K = &s.withK
	}
	responseData, responseErr := s.connection.RunREST(ctx, "/classifications", http.MethodPost, config)
	err := except.CheckResponnseDataErrorAndStatusCode(responseData, responseErr, 201)
	if err != nil {
		return nil, err
	}

	var responseClassification models.Classification
	parseErr := responseData.DecodeBodyIntoTarget(&responseClassification)
	if parseErr != nil {
		return nil, parseErr
	}
	if !s.withWaitForCompletion {
		return &responseClassification, nil
	}
	return s.waitForCompletion(ctx, responseClassification.ID)
}

func (s *Scheduler) waitForCompletion(ctx context.Context, uuid strfmt.UUID) (*models.Classification, error) {
	getter := Getter{
		connection: s.connection,
		withID:     string(uuid),
	}
	classification, err := getter.Do(ctx)
	if err != nil {
		return nil, err
	}
	for classification.Status == "running" {
		time.Sleep(2.0 * time.Second)
		classification, err = getter.Do(ctx)
		if err != nil {
			return nil, err
		}
	}
	return classification, nil
}
