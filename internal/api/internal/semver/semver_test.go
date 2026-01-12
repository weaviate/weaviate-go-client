package semver_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/internal/semver"
)

func Test(t *testing.T) {
	for _, tt := range []struct {
		v, w   string
		before bool
	}{
		{v: "v1.29.0", w: "v1.32.0", before: true},
		{v: "v1.32.0", w: "v1.32.0", before: false},
		{v: "v1.33.0", w: "v1.32.5", before: false},

		{v: "v1.29", w: "v1.32", before: true},
		{v: "v1.32", w: "v1.32", before: false},
		{v: "v1.33", w: "v1.32", before: false},

		// golang.org/x/mod/semver expects a leading "v".
		{v: "1.29.0", w: "1.32.0", before: true},
		{v: "1.32.0", w: "1.32.0", before: false},
		{v: "1.33.0", w: "1.32.5", before: false},
	} {
		t.Run(fmt.Sprintf("v=%s w=%s", tt.v, tt.w), func(t *testing.T) {
			assert.Equal(t, tt.before, semver.Before(tt.v, tt.w), "check v is before w")
			assert.Equal(t, !tt.before, semver.After(tt.v, tt.w), "check v is after w")
		})
	}
}

func TestMajorMinor(t *testing.T) {
	for _, tt := range []struct {
		v, w   string
		before bool
	}{
		{v: "v1.29.1", w: "v1.30.5", before: true},
		{v: "v1.29.1", w: "v1.29.1", before: false},
		{v: "v1.29.1", w: "v1.29.5", before: false},

		// golang.org/x/mod/semver expects a leading "v".
		{v: "1.29.1", w: "1.30.5", before: true},
		{v: "1.29.1", w: "1.29.1", before: false},
		{v: "1.29.1", w: "1.29.5", before: false},
	} {
		t.Run(fmt.Sprintf("v=%s w=%s", tt.v, tt.w), func(t *testing.T) {
			assert.Equal(t, tt.before, semver.BeforeMajorMinor(tt.v, tt.w), "check v is before w")
			assert.Equal(t, !tt.before, semver.AfterMajorMinor(tt.v, tt.w), "check v is after w")
		})
	}

	for _, tt := range []struct {
		v, w  string
		equal bool
	}{
		{v: "v1.29.1", w: "v1.28.0", equal: false},
		{v: "v1.29.1", w: "v1.29.1", equal: true},
		{v: "v1.29.1", w: "v1.29.6", equal: true},
		{v: "v1.29.1", w: "v1.33.0", equal: false},

		// golang.org/x/mod/semver expects a leading "v".
		{v: "1.29.1", w: "1.28.0", equal: false},
		{v: "1.29.1", w: "1.29.1", equal: true},
		{v: "1.29.1", w: "1.29.6", equal: true},
		{v: "1.29.1", w: "1.33.0", equal: false},
	} {
		t.Run(fmt.Sprintf("equal: v=%s w=%s", tt.v, tt.w), func(t *testing.T) {
			assert.Equal(t, tt.equal, semver.EqualMajorMinor(tt.v, tt.w), "check v == w")
		})
	}
}
