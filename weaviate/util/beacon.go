package util

import "fmt"

func BuildBeacon(id, className string, dbVersion *DBVersionSupport) string {
	if dbVersion.SupportsClassNameNamespacedEndpoints() {
		if len(className) > 0 {
			return fmt.Sprintf("weaviate://localhost/%v/%v", className, id)
		} else {
			dbVersion.WarnDeprecatedNonClassNameNamespacedEndpointsForBeacons()
		}
	} else if len(className) > 0 {
		dbVersion.WarnUsageOfNotSupportedClassNamespacedEndpointsForBeacons()
	}
	return fmt.Sprintf("weaviate://localhost/%v", id)
}
