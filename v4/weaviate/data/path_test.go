package data

import (
	"testing"

	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/util"
)

func Test_buildReferencesPath(t *testing.T) {
	type args struct {
		id                string
		className         string
		referenceProperty string
		version           *util.DBVersionSupport
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "build references path without class name with Weaviate version <1.14",
			args: args{
				id:                "some-uuid",
				referenceProperty: "ref-prop",
				version:           util.NewDBVersionSupport("1.13.2"),
			},
			want: "/objects/some-uuid/references/ref-prop",
		},
		{
			name: "build references path without class name with Weaviate version >=1.14",
			args: args{
				id:                "some-uuid",
				referenceProperty: "ref-prop",
				version:           util.NewDBVersionSupport("1.14.0"),
			},
			want: "/objects/some-uuid/references/ref-prop",
		},
		{
			name: "build references path with class name with Weaviate version <1.14",
			args: args{
				id:                "some-uuid",
				className:         "class-name",
				referenceProperty: "ref-prop",
				version:           util.NewDBVersionSupport("1.13.2"),
			},
			want: "/objects/some-uuid/references/ref-prop",
		},
		{
			name: "build references path with class name with Weaviate version >=1.14",
			args: args{
				id:                "some-uuid",
				className:         "class-name",
				referenceProperty: "ref-prop",
				version:           util.NewDBVersionSupport("1.14.0"),
			},
			want: "/objects/class-name/some-uuid/references/ref-prop",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildReferencesPath(tt.args.id, tt.args.className, tt.args.referenceProperty, tt.args.version); got != tt.want {
				t.Errorf("buildReferencesPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
