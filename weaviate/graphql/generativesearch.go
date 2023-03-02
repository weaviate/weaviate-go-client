package graphql

import "fmt"

type GenerativeSearchBuilder struct {
	prompt string
	task   string
}

func NewGenerativeSearch() *GenerativeSearchBuilder {
	return &GenerativeSearchBuilder{}
}

func (gs *GenerativeSearchBuilder) SingleResult(prompt string) *GenerativeSearchBuilder {
	gs.prompt = prompt
	return gs
}

func (gs *GenerativeSearchBuilder) GroupedResult(task string) *GenerativeSearchBuilder {
	gs.task = task
	return gs
}

func (gs *GenerativeSearchBuilder) build() string {
	resultFields := []Field{}
	query := ""
	if gs.prompt != "" {
		gs.prompt = "\"\"\"" + gs.prompt + "\"\"\""
		query += fmt.Sprintf("singleResult:{ prompt: %v }", gs.prompt)
		resultFields = append(resultFields, Field{Name: "singleResult"})
	}
	if gs.task != "" {
		gs.task = "\"\"\"" + gs.task + "\"\"\""
		query += fmt.Sprintf("groupedResult:{ task: %v }", gs.task)
		resultFields = append(resultFields, Field{Name: "groupedResult"})
	}
	resultFields = append(resultFields, Field{Name: "error"})
	finalFields := FieldsBuilder{fields: []Field{
		{
			Name: "_additional",
			Fields: []Field{
				{
					Name:   fmt.Sprintf("generate(%v)", query),
					Fields: resultFields,
				},
			},
		},
	}}
	return finalFields.build()
}
