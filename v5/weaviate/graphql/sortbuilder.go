package graphql

import (
	"fmt"
	"strings"
)

type Sort struct {
	Path  []string
	Order SortOrder
}

func (s Sort) build() string {
	clause := []string{}
	if len(s.Path) > 0 {
		path := make([]string, len(s.Path))
		for i := range s.Path {
			path[i] = fmt.Sprintf("\"%s\"", s.Path[i])
		}
		clause = append(clause, fmt.Sprintf("path:[%s]", strings.Join(path, ", ")))
	}
	if len(s.Order) > 0 {
		clause = append(clause, fmt.Sprintf("order:%s", string(s.Order)))
	}
	return fmt.Sprintf("{%s}", strings.Join(clause, " "))
}

type SortBuilder struct {
	sort []Sort
}

func (b *SortBuilder) build() string {
	clause := []string{}
	if len(b.sort) > 0 {
		for i := range b.sort {
			clause = append(clause, b.sort[i].build())
		}
	}
	return fmt.Sprintf("sort:[%s]", strings.Join(clause, ", "))
}
