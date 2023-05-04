package graphql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroupByBuilder_build(t *testing.T) {
	newGroupByBuilder := func() *GroupByArgumentBuilder {
		return &GroupByArgumentBuilder{}
	}
	tests := []struct {
		name    string
		groupBy *GroupByArgumentBuilder
		want    string
	}{
		{
			name: "all params",
			groupBy: newGroupByBuilder().WithPath([]string{"property"}).
				WithGroups(1).
				WithObjectsPerGroup(2),
			want: "groupBy:{path:[\"property\"] groups:1 objectsPerGroup:2}",
		},
		{
			name: "all params with cross ref path",
			groupBy: newGroupByBuilder().WithPath([]string{"ofClass", "class", "property1"}).
				WithGroups(10).
				WithObjectsPerGroup(11),
			want: "groupBy:{path:[\"ofClass\",\"class\",\"property1\"] groups:10 objectsPerGroup:11}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.groupBy.build())
		})
	}
}
