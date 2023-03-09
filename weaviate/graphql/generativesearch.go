package graphql

import (
	"fmt"
	"strings"
)

type GenerativeSearchBuilder struct {
	prompt string
	task   string
}

func NewGenerativeSearch() *GenerativeSearchBuilder {
	return &GenerativeSearchBuilder{}
}

func (gsb *GenerativeSearchBuilder) SingleResult(prompt string) *GenerativeSearchBuilder {
	gsb.prompt = prompt
	return gsb
}

func (gsb *GenerativeSearchBuilder) GroupedResult(task string) *GenerativeSearchBuilder {
	gsb.task = task
	return gsb
}

func (gsb *GenerativeSearchBuilder) build() Field {
	nameParts := []string{}
	fieldNames := []string{}

	if gsb.prompt != "" {
		nameParts = append(nameParts, fmt.Sprintf("singleResult:{prompt:\"\"\"%s\"\"\"}", gsb.prompt))
		fieldNames = append(fieldNames, "singleResult")
	}
	if gsb.task != "" {
		nameParts = append(nameParts, fmt.Sprintf("groupedResult:{task:\"\"\"%s\"\"\"}", gsb.task))
		fieldNames = append(fieldNames, "groupedResult")
	}

	if len(nameParts) == 0 {
		return Field{}
	}

	fieldNames = append(fieldNames, "error")

	fields := make([]Field, len(fieldNames))
	for i, fieldName := range fieldNames {
		fields[i] = Field{Name: fieldName}
	}

	return Field{
		Name:   fmt.Sprintf("generate(%s)", strings.Join(nameParts, " ")),
		Fields: fields,
	}
}
