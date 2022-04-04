package graphql

import (
	"fmt"
	"strings"
)

type Field struct {
	Name   string
	Fields []Field
}

func (f Field) build() string {
	clause := []string{}
	if len(f.Name) > 0 {
		clause = append(clause, f.Name)
	}
	if len(f.Fields) > 0 {
		fields := make([]string, len(f.Fields))
		for i := range f.Fields {
			fields[i] = f.Fields[i].build()
		}
		clause = append(clause, fmt.Sprintf("{%s}", strings.Join(fields, " ")))
	}
	return strings.Join(clause, "")
}

type FieldsBuilder struct {
	fields []Field
}

func (b *FieldsBuilder) WithFields(fields []Field) *FieldsBuilder {
	b.fields = fields
	return b
}

func (b *FieldsBuilder) build() string {
	clause := []string{}
	if len(b.fields) > 0 {
		for i := range b.fields {
			clause = append(clause, b.fields[i].build())
		}
	}
	return strings.Join(clause, " ")
}
