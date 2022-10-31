package util

import (
	"fmt"
	"strconv"
	"strings"
)

type VersionProvider interface {
	Version() string
}

type DBVersionSupport struct {
	dbVersionProvider VersionProvider
}

func NewDBVersionSupport(dbVersionProvider VersionProvider) *DBVersionSupport {
	return &DBVersionSupport{dbVersionProvider}
}

func (s *DBVersionSupport) SupportsClassNameNamespacedEndpoints() bool {
	versionNumbers := strings.Split(s.dbVersionProvider.Version(), ".")
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

	return (major == 1 && minor >= 14) || major >= 2
}

func (s *DBVersionSupport) WarnDeprecatedNonClassNameNamespacedEndpointsForObjects() {
	fmt.Printf("WARNING: Usage of objects paths without className is deprecated in Weaviate %s."+
		" Please provide className parameter\n", s.dbVersionProvider.Version())
}

func (s *DBVersionSupport) WarnDeprecatedNonClassNameNamespacedEndpointsForReferences() {
	fmt.Printf("WARNING: Usage of references paths without className is deprecated in Weaviate %s."+
		" Please provide className parameter\n", s.dbVersionProvider.Version())
}

func (s *DBVersionSupport) WarnDeprecatedNonClassNameNamespacedEndpointsForBeacons() {
	fmt.Printf("WARNING: Usage of beacon paths without className is deprecated in Weaviate %s."+
		" Please provide className parameter\n", s.dbVersionProvider.Version())
}

func (s *DBVersionSupport) WarnUsageOfNotSupportedClassNamespacedEndpointsForObjects() {
	fmt.Printf("WARNING: Usage of objects paths with className is not supported in Weaviate %s."+
		" className parameter is ignored\n", s.dbVersionProvider.Version())
}

func (s *DBVersionSupport) WarnUsageOfNotSupportedClassNamespacedEndpointsForReferences() {
	fmt.Printf("WARNING: Usage of references paths with className is not supported in Weaviate %s."+
		" className parameter is ignored\n", s.dbVersionProvider.Version())
}

func (s *DBVersionSupport) WarnUsageOfNotSupportedClassNamespacedEndpointsForBeacons() {
	fmt.Printf("WARNING: Usage of beacons paths with className is not supported in Weaviate %s."+
		" className parameter is ignored\n", s.dbVersionProvider.Version())
}

func (s *DBVersionSupport) WarnNotSupportedClassParameterInEndpointsForObjects() {
	fmt.Printf("WARNING: Usage of objects paths with class query parameter is not supported in Weaviate %s."+
		" class query parameter is ignored\n", s.dbVersionProvider.Version())
}
