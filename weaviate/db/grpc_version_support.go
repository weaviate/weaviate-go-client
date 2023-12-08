package db

import (
	"strconv"
	"strings"
)

type GRPCVersionSupport struct {
	dbVersionProvider versionProvider
}

func NewGRPCVersionSupport(dbVersionProvider versionProvider) *GRPCVersionSupport {
	return &GRPCVersionSupport{dbVersionProvider}
}

func (v *GRPCVersionSupport) SupportsVectorBytesField() bool {
	versionNumbers := strings.Split(v.dbVersionProvider.Version(), ".")
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
	if major == 1 && minor == 22 && !strings.Contains(versionNumbers[2], "rc") {
		patch, err := strconv.Atoi(versionNumbers[2])
		if err != nil {
			return false
		}
		return patch >= 6
	}
	return (major == 1 && minor >= 23) || major >= 2
}
