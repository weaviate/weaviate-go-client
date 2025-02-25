package classifications

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/filters"
	"github.com/weaviate/weaviate/entities/models"
)

// Scheduler builder to schedule a classification
type Scheduler struct {
	connection                 *connection.Connection
	classificationType         string
	withClassName              string
	withClassifyProperties     []string
	withBasedOnProperties      []string
	withSourceWhereFilter      *filters.WhereBuilder
	withTrainingSetWhereFilter *filters.WhereBuilder
	withTargetWhereFilter      *filters.WhereBuilder
	withWaitForCompletion      bool
	withSettings               interface{}
}

// WithType of classification e.g. knn or contextual
func (s *Scheduler) WithType(classificationType string) *Scheduler {
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
func (s *Scheduler) WithSourceWhereFilter(whereFilter *filters.WhereBuilder) *Scheduler {
	s.withSourceWhereFilter = whereFilter
	return s
}

// WithTrainingSetWhereFilter filter the objects that are used as training data. E.g. in a knn classification
func (s *Scheduler) WithTrainingSetWhereFilter(whereFilter *filters.WhereBuilder) *Scheduler {
	s.withTrainingSetWhereFilter = whereFilter
	return s
}

// WithTargetWhereFilter filter the label objects
func (s *Scheduler) WithTargetWhereFilter(whereFilter *filters.WhereBuilder) *Scheduler {
	s.withTargetWhereFilter = whereFilter
	return s
}

// WithSettings sets the classification settings
func (s *Scheduler) WithSettings(settings interface{}) *Scheduler {
	s.withSettings = settings
	return s
}

// WithWaitForCompletion block while classification is running (until classification succeeded or failed)
func (s *Scheduler) WithWaitForCompletion() *Scheduler {
	s.withWaitForCompletion = true
	return s
}

// Do schedule the classification in weaviate
func (s *Scheduler) Do(ctx context.Context) (*models.Classification, error) {
	responseData, responseErr := s.connection.RunREST(ctx, "/classifications", http.MethodPost, s.buildConfig())
	err := except.CheckResponseDataErrorAndStatusCode(responseData, responseErr, 201)
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

func (s *Scheduler) buildConfig() *models.Classification {
	config := &models.Classification{
		BasedOnProperties:  s.withBasedOnProperties,
		Class:              s.withClassName,
		ClassifyProperties: s.withClassifyProperties,
		Type:               s.classificationType,
		Settings:           s.withSettings,
	}

	// if any where filters are present, instantiate
	// the ClassificationFilters
	if s.withSourceWhereFilter != nil ||
		s.withTargetWhereFilter != nil ||
		s.withTrainingSetWhereFilter != nil {
		config.Filters = &models.ClassificationFilters{}
	}

	if s.withSourceWhereFilter != nil {
		config.Filters.SourceWhere = s.withSourceWhereFilter.Build()
	}
	if s.withTargetWhereFilter != nil {
		config.Filters.TargetWhere = s.withTargetWhereFilter.Build()
	}
	if s.withTrainingSetWhereFilter != nil {
		config.Filters.TrainingSetWhere = s.withTrainingSetWhereFilter.Build()
	}

	return config
}
