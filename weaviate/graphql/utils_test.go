package graphql

import (
	"testing"
)

func TestMarshalStringList(t *testing.T) {
	tests := []struct {
		in  []string
		out string
	}{
		{nil, "[]"},
		{[]string{}, "[]"},
		{[]string{"a"}, `["a"]`},
		{[]string{"ab", "ac"}, `["ab","ac"]`},
		{[]string{"ab", "ac", "abc"}, `["ab","ac","abc"]`},
	}
	for _, tc := range tests {
		got := string(marshalStrings(tc.in))
		if got != tc.out {
			t.Errorf("marshal(%v) got %v want %v", tc.in, got, tc.out)
		}

	}
}
