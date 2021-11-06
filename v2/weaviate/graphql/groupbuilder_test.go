package graphql

import "testing"

func TestGroupArgumentBuilder_build(t *testing.T) {
	type fields struct {
		withType  GroupType
		withForce *float32
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "with type closest and force",
			fields: fields{
				withType:  Closest,
				withForce: ptFloat32(0.4),
			},
			want: "group:{type: closest force: 0.4}",
		},
		{
			name: "with type merge and force",
			fields: fields{
				withType:  Merge,
				withForce: ptFloat32(0.49),
			},
			want: "group:{type: merge force: 0.49}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &GroupArgumentBuilder{}
			e = e.WithType(tt.fields.withType)
			if tt.fields.withForce != nil {
				e = e.WithForce(*tt.fields.withForce)
			}
			if got := e.build(); got != tt.want {
				t.Errorf("GroupArgumentBuilder.build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ptFloat32(val float32) *float32 {
	return &val
}
