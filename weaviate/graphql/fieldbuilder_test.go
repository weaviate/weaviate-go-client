package graphql

import "testing"

func TestFieldsBuilder_build(t *testing.T) {
	tests := []struct {
		name   string
		fields []Field
		want   string
	}{
		{
			name: "a b c",
			fields: []Field{
				{Name: "a"},
				{Name: "b"},
				{Name: "c"},
			},
			want: "a b c",
		},
		{
			name: "a{b} c{d{e}}",
			fields: []Field{
				{
					Name: "a",
					Fields: []Field{
						{Name: "b"},
					},
				},
				{
					Name: "c",
					Fields: []Field{
						{
							Name: "d",
							Fields: []Field{
								{
									Name: "e",
								},
							},
						},
					},
				},
			},
			want: "a{b} c{d{e}}",
		},
		{
			name: "_additional{classification{basedOn classifiedFields completed id scope}}",
			fields: []Field{
				{
					Name: "_additional",
					Fields: []Field{
						{
							Name: "classification",
							Fields: []Field{
								{
									Name: "basedOn",
								},
								{
									Name: "classifiedFields",
								},
								{
									Name: "completed",
								},
								{
									Name: "id",
								},
								{
									Name: "scope",
								},
							},
						},
					},
				},
			},
			want: "_additional{classification{basedOn classifiedFields completed id scope}}",
		},
		{
			name: "inPublication{... on Publication{name}}",
			fields: []Field{
				{
					Name: "inPublication",
					Fields: []Field{
						{
							Name: "... on Publication",
							Fields: []Field{
								{
									Name: "name",
								},
							},
						},
					},
				},
			},
			want: "inPublication{... on Publication{name}}",
		},
		{
			name: "_additional{certainty}",
			fields: []Field{
				{
					Name: "_additional",
					Fields: []Field{
						{
							Name: "certainty",
						},
					},
				},
			},
			want: "_additional{certainty}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &FieldsBuilder{
				fields: tt.fields,
			}
			if got := b.build(); got != tt.want {
				t.Errorf("FieldsBuilder.build() = %v, want %v", got, tt.want)
			}
		})
	}
}
