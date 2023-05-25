package graphql

import (
	"encoding/json"
	"fmt"
	"strings"
)

type GenerativeSearchBuilder struct {
	prompt     string
	task       string
	properties []string
}

func NewGenerativeSearch() *GenerativeSearchBuilder {
	return &GenerativeSearchBuilder{}
}

func (gsb *GenerativeSearchBuilder) SingleResult(prompt string) *GenerativeSearchBuilder {
	gsb.prompt = prompt
	return gsb
}

func (gsb *GenerativeSearchBuilder) GroupedResult(task string, properties ...string) *GenerativeSearchBuilder {
	gsb.task = task
	gsb.properties = properties
	return gsb
}

func (gsb *GenerativeSearchBuilder) build() Field {
	nameParts := []string{}
	fieldNames := []string{}

	if gsb.prompt != "" {
		nameParts = append(nameParts, fmt.Sprintf("singleResult:{prompt:\"\"\"%s\"\"\"}", gsb.prompt))
		fieldNames = append(fieldNames, "singleResult")
	}
	if gsb.task != "" || len(gsb.properties) > 0 {
		argParts := []string{}
		if gsb.task != "" {
			argParts = append(argParts, fmt.Sprintf("task:\"\"\"%s\"\"\"", gsb.task))
		}
		if len(gsb.properties) > 0 {
			properties, _ := json.Marshal(gsb.properties)
			argParts = append(argParts, fmt.Sprintf("properties:%v", string(properties)))
		}
		nameParts = append(nameParts, fmt.Sprintf("groupedResult:{%s}", strings.Join(argParts, ",")))
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
