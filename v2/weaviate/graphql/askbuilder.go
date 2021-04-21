package graphql

import (
	"encoding/json"
	"fmt"
	"strings"
)

type AskArgumentBuilder struct {
	question      string
	properties    []string
	withCertainty bool
	certainty     float32
}

// WithQuestion the question to be asked
func (e *AskArgumentBuilder) WithQuestion(question string) *AskArgumentBuilder {
	e.question = question
	return e
}

// WithProperties the list of properties that contain a text to look in for an answer
func (e *AskArgumentBuilder) WithProperties(properties []string) *AskArgumentBuilder {
	e.properties = properties
	return e
}

// WithCertainty that specifies how accurate the answer should be
func (e *AskArgumentBuilder) WithCertainty(certainty float32) *AskArgumentBuilder {
	e.withCertainty = true
	e.certainty = certainty
	return e
}

// Build build the given clause
func (e *AskArgumentBuilder) build() string {
	clause := []string{}
	if len(e.question) > 0 {
		clause = append(clause, fmt.Sprintf("question: \"%s\"", e.question))
	}
	if len(e.properties) > 0 {
		properties, _ := json.Marshal(e.properties)
		clause = append(clause, fmt.Sprintf("properties: %v", string(properties)))
	}
	if e.withCertainty {
		clause = append(clause, fmt.Sprintf("certainty: %v", e.certainty))
	}
	return fmt.Sprintf("ask:{%s}", strings.Join(clause, " "))
}
