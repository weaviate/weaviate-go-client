package data

import (
	"path"
	"strconv"
	"strings"
)

func buildObjectsPath(id, className, version string) string {
	p := "/objects"
	if supportsClassNameNamespacedEndpoints(version) && className != "" {
		p = path.Join(p, className)
	}
	if id != "" {
		p = path.Join(p, id)
	}
	return p
}

func supportsClassNameNamespacedEndpoints(version string) bool {
	versionNumbers := strings.Split(version, ".")
	if len(versionNumbers) < 3 {
		return false
	}

	major, err := strconv.Atoi(versionNumbers[0])
	if err != nil {
		return false
	}

	minor, err := strconv.Atoi(versionNumbers[1])
	if err != nil {
		return false
	}

	return major >= 1 && minor >= 14
}
