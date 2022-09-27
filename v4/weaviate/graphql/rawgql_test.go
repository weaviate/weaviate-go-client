package graphql

import "testing"

func TestRawGQLBuilder_build(t *testing.T) {
	type fields struct {
		query string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "with query",
			fields: fields{
				query: "query { Get { Things { Thing { name } } } }",
			},
			want: "query { Get { Things { Thing { name } } } }",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Raw{query: tt.fields.query}
			if got := e.build(); got != tt.want {
				t.Errorf("Raw.build() = %v, want %v", got, tt.want)
			}
		})
	}
}
