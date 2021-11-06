package graphql

import (
	"fmt"
	"strings"
)

// GroupType filter
type GroupType string

// Merge group type filter
const Merge GroupType = "merge"

// Closest group type filter
const Closest GroupType = "closest"

type GroupArgumentBuilder struct {
	withType  GroupType
	withForce bool
	force     float32
}

// WithType the type of group argument
func (b *GroupArgumentBuilder) WithType(groupType GroupType) *GroupArgumentBuilder {
	b.withType = groupType
	return b
}

// WithBeacon the beacon of the object
func (b *GroupArgumentBuilder) WithForce(force float32) *GroupArgumentBuilder {
	b.withForce = true
	b.force = force
	return b
}

// Build build the given clause
func (b *GroupArgumentBuilder) build() string {
	clause := []string{}
	if len(b.withType) > 0 {
		clause = append(clause, fmt.Sprintf("type: %s", b.withType))
	}
	if b.withForce {
		clause = append(clause, fmt.Sprintf("force: %v", b.force))
	}
	return fmt.Sprintf("group:{%s}", strings.Join(clause, " "))
}
