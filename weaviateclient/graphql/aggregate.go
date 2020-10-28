package graphql

import (
	"context"
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/paragons"
	"github.com/semi-technologies/weaviate/entities/models"
	"strings"
)

type Aggregate struct {
	connection *connection.Connection
}

func (a *Aggregate) Things() *AggregateBuilder {
	return &AggregateBuilder{
		connection:   a.connection,
		semanticKind: paragons.SemanticKindThings,
	}
}

func (a *Aggregate) Actions() *AggregateBuilder {
	return &AggregateBuilder{
		connection:   a.connection,
		semanticKind: paragons.SemanticKindActions,
	}
}


type AggregateBuilder struct {
	connection rest
	semanticKind paragons.SemanticKind
	fields string
	className string
}

func (ab *AggregateBuilder) WithFields(fields string) *AggregateBuilder {
	ab.fields = fields
	return ab
}

func (ab *AggregateBuilder) WithClassName(name string) *AggregateBuilder {
	ab.className = name
	return ab
}

func (ab *AggregateBuilder) Do(ctx context.Context) (*models.GraphQLResponse, error) {
	return runGraphQLQuery(ctx, ab.connection, ab.build())
}

func (ab *AggregateBuilder) build() string {
	semanticKind := strings.Title(string(ab.semanticKind))
	return 	fmt.Sprintf("{Aggregate{%v{%v{%v}}}}", semanticKind, ab.className, ab.fields)
}