package crossref

import (
	"fmt"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/db"
)

func BuildBeacon(id, className string, dbVersion *db.VersionSupport) string {
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
