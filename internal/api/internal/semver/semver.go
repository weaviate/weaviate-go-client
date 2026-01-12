package semver

import (
	"strings"

	"golang.org/x/mod/semver"
)

// Before returns whether version v is before version w.
// If v == w, Before returns false.
func Before(v, w string) bool {
	v, w = prefix(v), prefix(w)
	return semver.Compare(v, w) < 0
}

// After returns whether version v after version w.
// If v == w, After returns true.
func After(v, w string) bool {
	v, w = prefix(v), prefix(w)
	return semver.Compare(v, w) >= 0
}

// BeforeMajorMinor returns whether 'major.minor' version v is before w.
// If v == w, BeforeMajorMinor returns false.
func BeforeMajorMinor(v, w string) bool {
	v, w = prefix(v), prefix(w)
	v, w = semver.MajorMinor(v), semver.MajorMinor(w)
	return Before(v, w)
}

// AfterMajorMinor returns whether 'major.minor' version v is after w.
// If v == w, After returns true.
func AfterMajorMinor(v, w string) bool {
	v, w = prefix(v), prefix(w)
	v, w = semver.MajorMinor(v), semver.MajorMinor(w)
	return After(v, w)
}

// EqualMajorMinor returns whether 'major.minor' version v is equal to w.
func EqualMajorMinor(v, w string) bool {
	v, w = prefix(v), prefix(w)
	v, w = semver.MajorMinor(v), semver.MajorMinor(w)
	return semver.Compare(v, w) == 0
}

// prefix adds a "v" prefix to the version if it doesn't have one.
func prefix(v string) string {
	if !strings.HasPrefix(v, "v") {
		v = "v" + v
	}
	return v
}
