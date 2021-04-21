package graphql

import (
	"fmt"
	"strings"
)

type NearObjectArgumentBuilder struct {
	id            string
	beacon        string
	withCertainty bool
	certainty     float32
}

// WithID the id of the object
func (e *NearObjectArgumentBuilder) WithID(id string) *NearObjectArgumentBuilder {
	e.id = id
	return e
}

// WithBeacon the beacon of the object
func (e *NearObjectArgumentBuilder) WithBeacon(beacon string) *NearObjectArgumentBuilder {
	e.beacon = beacon
	return e
}

// WithCertainty that is minimally required for an object to be included in the result set
func (e *NearObjectArgumentBuilder) WithCertainty(certainty float32) *NearObjectArgumentBuilder {
	e.withCertainty = true
	e.certainty = certainty
	return e
}

// Build build the given clause
func (e *NearObjectArgumentBuilder) build() string {
	clause := []string{}
	if len(e.id) > 0 {
		clause = append(clause, fmt.Sprintf("id: \"%s\"", e.id))
	}
	if len(e.beacon) > 0 {
		clause = append(clause, fmt.Sprintf("beacon: \"%s\"", e.beacon))
	}
	if e.withCertainty {
		clause = append(clause, fmt.Sprintf("certainty: %v", e.certainty))
	}
	return fmt.Sprintf("nearObject:{%s} ", strings.Join(clause, " "))
}
