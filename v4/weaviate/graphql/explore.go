package graphql

import (
	"context"
	"fmt"

	"github.com/semi-technologies/weaviate/entities/models"
)

// Explore query builder
type Explore struct {
	connection     rest
	fields         []ExploreFields
	withNearText   *NearTextArgumentBuilder
	withNearObject *NearObjectArgumentBuilder
	withAsk        *AskArgumentBuilder
	withNearImage  *NearImageArgumentBuilder
	withNearVector *NearVectorArgumentBuilder
}

// WithNearText adds nearText to clause
func (e *Explore) WithNearText(nearText *NearTextArgumentBuilder) *Explore {
	e.withNearText = nearText
	return e
}

// WithNearObject adds nearObject to clause
func (e *Explore) WithNearObject(nearObject *NearObjectArgumentBuilder) *Explore {
	e.withNearObject = nearObject
	return e
}

// WithAsk adds ask to clause
func (e *Explore) WithAsk(ask *AskArgumentBuilder) *Explore {
	e.withAsk = ask
	return e
}

// WithNearImage adds nearImage to clause
func (e *Explore) WithNearImage(nearImage *NearImageArgumentBuilder) *Explore {
	e.withNearImage = nearImage
	return e
}

// WithNearVector clause to find close objects
func (e *Explore) WithNearVector(nearVector *NearVectorArgumentBuilder) *Explore {
	e.withNearVector = nearVector
	return e
}

// WithFields that should be included in the result set
func (e *Explore) WithFields(fields ...ExploreFields) *Explore {
	e.fields = fields
	return e
}

func (e *Explore) createFilterClause() string {
	if e.withNearText != nil {
		return e.withNearText.build()
	}
	if e.withNearObject != nil {
		return e.withNearObject.build()
	}
	if e.withAsk != nil {
		return e.withAsk.build()
	}
	if e.withNearImage != nil {
		return e.withNearImage.build()
	}
	if e.withNearVector != nil {
		return e.withNearVector.build()
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

	query := fmt.Sprintf("{Explore(%v){%v}}", filterClause, fields)

	return query
}

// Do execute explore search
func (e *Explore) Do(ctx context.Context) (*models.GraphQLResponse, error) {
	return runGraphQLQuery(ctx, e.connection, e.build())
}
