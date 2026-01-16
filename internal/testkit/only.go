package testkit

import "testing"

type ExclusiveTest interface {
	// Exclusive returns true if other test cases should be skipped.
	Exclusive() bool
}

// Only marks the test as exclusive if set to true.
// Only should be embedded in the test case struct.
// See the example for [RunOnly].
type Only bool

var _ ExclusiveTest = (*Only)(nil)

func (o Only) Exclusive() bool { return bool(o) }

// RunOnly runs a collection of table tests which contains [ExclusiveTest] cases.
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
//	// The first test case will be reported as SKIP'ed.
//	// Only second test case will run.
//	testkit.RunOnly(t, []test{
//		{
//			name: "good test",
//			input: 10, want: 10,
//		},
//		{
//			name: "needs debugging",
//			input: 5, want: 43
//			Only: true,
//		},
//	}, func(t *testing.T, tt test) {
//		t.Run(tt.name, func(t *testing.T) {
//			require.Equal(t, tt.want, tt.input)
//		})
//	})
//
// Any number of tests can be marked "exclusive" by setting Only: true.
func RunOnly[T ExclusiveTest](t *testing.T, tests []T, f func(*testing.T, T)) {
	t.Helper()

	var only []T
	var skip []T
	for _, tt := range tests {
		if tt.Exclusive() {
			only = append(only, tt)
		} else {
			skip = append(skip, tt)
		}
	}

	if only == nil {
		only = tests
		skip = nil
	}

	for _, tt := range only {
		f(t, tt)
	}

	for range skip {
		t.Run("testkit.Only=false", func(t *testing.T) { t.Skip() })
	}
}
