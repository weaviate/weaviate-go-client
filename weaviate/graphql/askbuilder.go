package graphql

import (
	"encoding/json"
	"fmt"
	"strings"
)

type AskArgumentBuilder struct {
	question        string
	properties      []string
	withCertainty   bool
	certainty       float32
	withDistance    bool
	distance        float32
	withAutocorrect bool
	autocorrect     bool
	withRerank      bool
	rerank          bool
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

// WithCertainty specifies how accurate the answer should be
func (e *AskArgumentBuilder) WithCertainty(certainty float32) *AskArgumentBuilder {
	e.withCertainty = true
	e.certainty = certainty
	return e
}

// WithDistance specifies how accurate the answer should be
func (e *AskArgumentBuilder) WithDistance(distance float32) *AskArgumentBuilder {
	e.withDistance = true
	e.distance = distance
	return e
}

// WithAutocorrect this is a setting enabling autocorrect of question text
func (e *AskArgumentBuilder) WithAutocorrect(autocorrect bool) *AskArgumentBuilder {
	e.withAutocorrect = true
	e.autocorrect = autocorrect
	return e
}

// WithRerank this is a setting enabling re-ranking of results based on certainty
func (e *AskArgumentBuilder) WithRerank(rerank bool) *AskArgumentBuilder {
	e.withRerank = true
	e.rerank = rerank
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
	if e.withDistance {
		clause = append(clause, fmt.Sprintf("distance: %v", e.distance))
	}
	if e.withAutocorrect {
		clause = append(clause, fmt.Sprintf("autocorrect: %v", e.autocorrect))
	}
	if e.withRerank {
		clause = append(clause, fmt.Sprintf("rerank: %v", e.rerank))
	}
	return fmt.Sprintf("ask:{%s}", strings.Join(clause, " "))
}
