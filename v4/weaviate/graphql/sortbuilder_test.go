package graphql

import "testing"

func TestSortBuilder_build(t *testing.T) {
	tests := []struct {
		name string
		sort []Sort
		want string
	}{
		{
			name: "simple sort",
			sort: []Sort{
				{Path: []string{"property"}},
			},
			want: "sort:[{path:[\"property\"]}]",
		},
		{
			name: "simple sort with ASC",
			sort: []Sort{
				{Path: []string{"property"}, Order: Asc},
			},
			want: "sort:[{path:[\"property\"] order:asc}]",
		},
		{
			name: "simple sort with DESC and double paths",
			sort: []Sort{
				{Path: []string{"property1", "property2"}, Order: Desc},
			},
			want: "sort:[{path:[\"property1\", \"property2\"] order:desc}]",
		},
		{
			name: "complex sort",
			sort: []Sort{
				{Path: []string{"property1"}},
				{Path: []string{"property2"}},
			},
			want: "sort:[{path:[\"property1\"]}, {path:[\"property2\"]}]",
		},
		{
			name: "complex sort with ASC and DESC",
			sort: []Sort{
				{Path: []string{"property1"}, Order: Asc},
				{Path: []string{"property2"}, Order: Desc},
			},
			want: "sort:[{path:[\"property1\"] order:asc}, {path:[\"property2\"] order:desc}]",
		},
		{
			name: "complex sort with ASC and DESC and double paths",
			sort: []Sort{
				{Path: []string{"property1"}, Order: Asc},
				{Path: []string{"property2", "property3"}},
				{Path: []string{"property4"}, Order: Desc},
			},
			want: "sort:[{path:[\"property1\"] order:asc}, {path:[\"property2\", \"property3\"]}, {path:[\"property4\"] order:desc}]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &SortBuilder{
				sort: tt.sort,
			}
			if got := b.build(); got != tt.want {
				t.Errorf("SortBuilder.build() = %v, want %v", got, tt.want)
			}
		})
	}
}
