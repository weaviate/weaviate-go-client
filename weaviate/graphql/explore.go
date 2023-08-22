package graphql

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/weaviate/weaviate/entities/models"
)

// Explore query builder
type Explore struct {
	connection           rest
	includesFilterClause bool // true if brackets behind class is needed
	includesLimit        bool
	limit                int
	includesOffset       bool
	offset               int
	fields               []ExploreFields
	withNearVector       *NearVectorArgumentBuilder
	withNearObject       *NearObjectArgumentBuilder
	withNearText         *NearTextArgumentBuilder
	withAsk              *AskArgumentBuilder
	withNearImage        *NearImageArgumentBuilder
	withNearAudio        *NearAudioArgumentBuilder
	withNearVideo        *NearVideoArgumentBuilder
	withNearDepth        *NearDepthArgumentBuilder
	withNearThermal      *NearThermalArgumentBuilder
	withNearImu          *NearImuArgumentBuilder
}

// WithNearText adds nearText to clause
func (e *Explore) WithNearText(nearText *NearTextArgumentBuilder) *Explore {
	e.includesFilterClause = true
	e.withNearText = nearText
	return e
}

// WithNearObject adds nearObject to clause
func (e *Explore) WithNearObject(nearObject *NearObjectArgumentBuilder) *Explore {
	e.includesFilterClause = true
	e.withNearObject = nearObject
	return e
}

// WithAsk adds ask to clause
func (e *Explore) WithAsk(ask *AskArgumentBuilder) *Explore {
	e.includesFilterClause = true
	e.withAsk = ask
	return e
}

// WithNearImage adds nearImage to clause
func (e *Explore) WithNearImage(nearImage *NearImageArgumentBuilder) *Explore {
	e.includesFilterClause = true
	e.withNearImage = nearImage
	return e
}

// WithNearAudio adds nearAudio to clause
func (e *Explore) WithNearAudio(nearAudio *NearAudioArgumentBuilder) *Explore {
	e.includesFilterClause = true
	e.withNearAudio = nearAudio
	return e
}

// WithNearVideo adds nearVideo to clause
func (e *Explore) WithNearVideo(nearVideo *NearVideoArgumentBuilder) *Explore {
	e.includesFilterClause = true
	e.withNearVideo = nearVideo
	return e
}

// WithNearDepth adds nearDepth to clause
func (e *Explore) WithNearDepth(nearDepth *NearDepthArgumentBuilder) *Explore {
	e.includesFilterClause = true
	e.withNearDepth = nearDepth
	return e
}

// WithNearThermal adds nearThermal to clause
func (e *Explore) WithNearThermal(nearThermal *NearThermalArgumentBuilder) *Explore {
	e.includesFilterClause = true
	e.withNearThermal = nearThermal
	return e
}

// WithNearImu adds nearIMU to clause
func (e *Explore) WithNearImu(nearImu *NearImuArgumentBuilder) *Explore {
	e.includesFilterClause = true
	e.withNearImu = nearImu
	return e
}

// WithNearVector clause to find close objects
func (e *Explore) WithNearVector(nearVector *NearVectorArgumentBuilder) *Explore {
	e.includesFilterClause = true
	e.withNearVector = nearVector
	return e
}

// WithFields that should be included in the result set
func (e *Explore) WithFields(fields ...ExploreFields) *Explore {
	e.fields = fields
	return e
}

// WithLimit of objects in the result set
func (e *Explore) WithLimit(limit int) *Explore {
	e.includesFilterClause = true
	e.includesLimit = true
	e.limit = limit
	return e
}

// WithOffset of objects in the result set
func (e *Explore) WithOffset(offset int) *Explore {
	e.includesFilterClause = true
	e.includesOffset = true
	e.offset = offset
	return e
}

func (e *Explore) createFilterClause() string {
	if e.includesFilterClause {
		filters := []string{}
		for _, b := range []argumentBuilder{
			e.withAsk, e.withNearText, e.withNearObject, e.withNearVector, e.withNearImage,
			e.withNearAudio, e.withNearVideo, e.withNearDepth, e.withNearThermal, e.withNearImu,
		} {
			bVal := reflect.ValueOf(b)
			if bVal.Kind() == reflect.Ptr && !bVal.IsNil() {
				filters = append(filters, b.build())
			}
		}
		if e.includesLimit {
			filters = append(filters, fmt.Sprintf("limit: %v", e.limit))
		}
		if e.includesOffset {
			filters = append(filters, fmt.Sprintf("offset: %v", e.offset))
		}
		return fmt.Sprintf("(%s)", strings.Join(filters, ", "))
	}
	return ""
}

// build query
func (e *Explore) build() string {
	fields := ""
	for _, field := range e.fields {
		fields += fmt.Sprintf("%v ", field)
	}

	filterClause := e.createFilterClause()

	query := fmt.Sprintf("{Explore%v{%v}}", filterClause, fields)

	return query
}

// Do execute explore search
func (e *Explore) Do(ctx context.Context) (*models.GraphQLResponse, error) {
	return runGraphQLQuery(ctx, e.connection, e.build())
}
