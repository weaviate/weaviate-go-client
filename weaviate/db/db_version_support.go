package db

import (
	"fmt"
	"strconv"
	"strings"
)

type versionProvider interface {
	Version() string
}

type VersionSupport struct {
	dbVersionProvider versionProvider
}

func NewDBVersionSupport(dbVersionProvider versionProvider) *VersionSupport {
	return &VersionSupport{dbVersionProvider}
}

func (v *VersionSupport) SupportsClassNameNamespacedEndpoints() bool {
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

	return (major == 1 && minor >= 14) || major >= 2
}

func (v *VersionSupport) WarnDeprecatedNonClassNameNamespacedEndpointsForObjects() {
	fmt.Printf("WARNING: Usage of objects paths without className is deprecated in Weaviate %v."+
		" Please provide className parameter\n", v.dbVersionProvider.Version())
}

func (v *VersionSupport) WarnDeprecatedNonClassNameNamespacedEndpointsForReferences() {
	fmt.Printf("WARNING: Usage of references paths without className is deprecated in Weaviate %v."+
		" Please provide className parameter\n", v.dbVersionProvider.Version())
}

func (v *VersionSupport) WarnDeprecatedNonClassNameNamespacedEndpointsForBeacons() {
	fmt.Printf("WARNING: Usage of beacon paths without className is deprecated in Weaviate %v."+
		" Please provide className parameter\n", v.dbVersionProvider.Version())
}

func (v *VersionSupport) WarnUsageOfNotSupportedClassNamespacedEndpointsForObjects() {
	fmt.Printf("WARNING: Usage of objects paths with className is not supported in Weaviate %v."+
		" className parameter is ignored\n", v.dbVersionProvider.Version())
}

func (v *VersionSupport) WarnUsageOfNotSupportedClassNamespacedEndpointsForReferences() {
	fmt.Printf("WARNING: Usage of references paths with className is not supported in Weaviate %v."+
		" className parameter is ignored\n", v.dbVersionProvider.Version())
}

func (v *VersionSupport) WarnUsageOfNotSupportedClassNamespacedEndpointsForBeacons() {
	fmt.Printf("WARNING: Usage of beacons paths with className is not supported in Weaviate %v."+
		" className parameter is ignored\n", v.dbVersionProvider.Version())
}

func (v *VersionSupport) WarnNotSupportedClassParameterInEndpointsForObjects() {
	fmt.Printf("WARNING: Usage of objects paths with class query parameter is not supported in Weaviate %v."+
		" class query parameter is ignored\n", v.dbVersionProvider.Version())
}
