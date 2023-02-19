package graphql

type GenerativeSearch interface {
	build() string
}

type singleResult struct {
	prompt string
}

type groupedResult struct {
	task string
}

func NewGSWithSingleResult(prompt string) GenerativeSearch {
	return &singleResult{prompt: prompt}
}

func NewGSWithGroupedResult(task string) GenerativeSearch {
	return &groupedResult{task: task}
}

func (sr *singleResult) build() string {
	sr.prompt = "\"\"\"" + sr.prompt + "\"\"\""
	fields := FieldsBuilder{fields: []Field{
		{
			Name: "_additional",
			Fields: []Field{
				{
					Name: "generate(singleResult:{ prompt: " + sr.prompt + " })",
					Fields: []Field{
						{Name: "singleResult"},
						{Name: "error"},
					},
				},
			},
		},
	}}
	return fields.build()
}

func (gr *groupedResult) build() string {
	gr.task = "\"\"\"" + gr.task + "\"\"\""
	fields := FieldsBuilder{fields: []Field{
		{
			Name: "_additional",
			Fields: []Field{
				{
					Name: "generate(groupedResult:{ task: " + gr.task + " })",
					Fields: []Field{
						{Name: "groupedResult"},
						{Name: "error"},
					},
				},
			},
		},
	}}
	return fields.build()
}
