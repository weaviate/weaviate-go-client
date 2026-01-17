package testkit

import "testing"

type ExclusiveTest interface {
	// Exclusive returns true if other test cases should be skipped.
	Exclusive() bool
}

// Only marks the test as exclusive if set to true.
// Only should be embedded in the test case struct.
// See the example for [WithOnly].
type Only bool

var _ ExclusiveTest = (*Only)(nil)

func (o Only) Exclusive() bool { return bool(o) }

// WithOnly filters a collection of table tests which contains [ExclusiveTest] cases.
//
// Example:
//
//	type test struct {
//		testkit.Only
//
//		name string
//		input, want int
//	}
//
//	// Only the second test case will run.
//	for _, tt := range testkit.WithOnly(t, []test{
//		{
//			name: "good test",
//			input: 10, want: 10,
//		},
//		{
//			name: "needs debugging",
//			input: 5, want: 43
//			Only: true,
//		},
//	}) {
//		t.Run(tt.name, func(t *testing.T) {
//			require.Equal(t, tt.want, tt.input)
//		})
//	}
//
// Any number of tests can be marked "exclusive" by setting Only: true.
func WithOnly[T ExclusiveTest](t *testing.T, tests []T) []T {
	var only []T
	for _, tt := range tests {
		if tt.Exclusive() {
			only = append(only, tt)
		}
	}

	if only == nil {
		only = tests
	}
	return only
}
