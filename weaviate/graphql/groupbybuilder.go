package graphql

import (
	"encoding/json"
	"fmt"
	"strings"
)

type GroupByArgumentBuilder struct {
	path                []string
	groups              int
	withGroups          bool
	objectsPerGroup     int
	withObjectsPerGroup bool
}

// WithPath the property by which is should be grouped by
func (b *GroupByArgumentBuilder) WithPath(path []string) *GroupByArgumentBuilder {
	b.path = path
	return b
}

// WithGroups maximum number of groups
func (b *GroupByArgumentBuilder) WithGroups(groups int) *GroupByArgumentBuilder {
	b.withGroups = true
	b.groups = groups
	return b
}

// WithObjectsPerGroup maximum number of elements in group
func (b *GroupByArgumentBuilder) WithObjectsPerGroup(objectsPerGroup int) *GroupByArgumentBuilder {
	b.withObjectsPerGroup = true
	b.objectsPerGroup = objectsPerGroup
	return b
}

// Build build the given clause
func (b *GroupByArgumentBuilder) build() string {
	clause := []string{}
	path, _ := json.Marshal(b.path)

	clause = append(clause, fmt.Sprintf("path:%s", path))
	if b.withGroups {
		clause = append(clause, fmt.Sprintf("groups:%v", b.groups))
	}
	if b.withObjectsPerGroup {
		clause = append(clause, fmt.Sprintf("objectsPerGroup:%v", b.objectsPerGroup))
	}
	return fmt.Sprintf("groupBy:{%s}", strings.Join(clause, " "))
}
