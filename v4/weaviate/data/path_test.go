package data

import (
	"testing"

	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/util"
	"github.com/stretchr/testify/assert"
)

func Test_buildObjectPath(t *testing.T) {
	version := "1.17.0"

	tests := []struct {
		name         string
		getter       *ObjectsGetter
		expectedPath string
	}{
		{
			name: "with consistency level only",
			getter: newTestGetter(version).WithID("123").
				WithClassName("SomeClass").
				WithConsistencyLevel("QUORUM"),
			expectedPath: "/objects/SomeClass/123?consistency_level=QUORUM",
		},
		{
			name: "with node name only",
			getter: newTestGetter(version).WithID("123").
				WithClassName("SomeClass").
				WithNodeName("node1"),
			expectedPath: "/objects/SomeClass/123?node_name=node1",
		},
		{
			name: "with consistency level and with vector and classification",
			getter: newTestGetter(version).WithID("123").
				WithClassName("SomeClass").
				WithConsistencyLevel("QUORUM").
				WithAdditional("classification").
				WithVector(),
			expectedPath: "/objects/SomeClass/123?include=classification,vector&consistency_level=QUORUM",
		},
		{
			name: "with node name and with vector and classification",
			getter: newTestGetter(version).WithID("123").
				WithClassName("SomeClass").
				WithNodeName("node1").
				WithAdditional("classification").
				WithVector(),
			expectedPath: "/objects/SomeClass/123?include=classification,vector&node_name=node1",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expectedPath, test.getter.buildPath())
	}
}

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
				version:           newDBVersionSupportForTests("1.13.2"),
			},
			want: "/objects/some-uuid/references/ref-prop",
		},
		{
			name: "build references path without class name with Weaviate version >=1.14",
			args: args{
				id:                "some-uuid",
				referenceProperty: "ref-prop",
				version:           newDBVersionSupportForTests("1.14.0"),
			},
			want: "/objects/some-uuid/references/ref-prop",
		},
		{
			name: "build references path with class name with Weaviate version <1.14",
			args: args{
				id:                "some-uuid",
				className:         "class-name",
				referenceProperty: "ref-prop",
				version:           newDBVersionSupportForTests("1.13.2"),
			},
			want: "/objects/some-uuid/references/ref-prop",
		},
		{
			name: "build references path with class name with Weaviate version >=1.14",
			args: args{
				id:                "some-uuid",
				className:         "class-name",
				referenceProperty: "ref-prop",
				version:           newDBVersionSupportForTests("1.14.0"),
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

type dbVersionProviderMock struct {
	version string
}

func (s *dbVersionProviderMock) Version() string {
	return s.version
}

func newDBVersionProviderMock(version string) *dbVersionProviderMock {
	return &dbVersionProviderMock{version}
}

func newDBVersionSupportForTests(version string) *util.DBVersionSupport {
	return util.NewDBVersionSupport(newDBVersionProviderMock(version))
}

func newTestGetter(version string) *ObjectsGetter {
	return &ObjectsGetter{
		dbVersionSupport: newDBVersionSupportForTests(version),
	}
}
