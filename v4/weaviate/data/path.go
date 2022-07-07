package data

import (
	"fmt"
	"path"

	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/util"
)

func buildObjectsGetPath(id, className string, dbVersion *util.DBVersionSupport) string {
	return buildObjectsPath(id, className, dbVersion, "get")
}

func buildObjectsDeletePath(id, className string, dbVersion *util.DBVersionSupport) string {
	return buildObjectsPath(id, className, dbVersion, "delete")
}

func buildObjectsCheckPath(id, className string, dbVersion *util.DBVersionSupport) string {
	return buildObjectsPath(id, className, dbVersion, "check")
}

func buildObjectsUpdatePath(id, className string, dbVersion *util.DBVersionSupport) string {
	return buildObjectsPath(id, className, dbVersion, "update")
}

func buildObjectsPath(id, className string, dbVersion *util.DBVersionSupport, action string) string {
	p := "/objects"
	if len(id) > 0 {
		if dbVersion.SupportsClassNameNamespacedEndpoints() {
			if len(className) > 0 {
				p = path.Join(p, className)
			} else {
				dbVersion.WarnDeprecatedNonClassNameNamespacedEndpointsForObjects()
			}
		} else if len(className) > 0 && action != "update" {
			dbVersion.WarnUsageOfNotSupportedClassNamespacedEndpointsForObjects()
		}
		p = path.Join(p, id)
	}
	return p
}

func buildReferencesPath(id, className, referenceProperty string, dbVersion *util.DBVersionSupport) string {
	if dbVersion.SupportsClassNameNamespacedEndpoints() {
		if len(className) > 0 {
			return fmt.Sprintf("/objects/%v/%v/references/%v", className, id, referenceProperty)
		} else {
			dbVersion.WarnDeprecatedNonClassNameNamespacedEndpointsForReferences()
		}
	} else if len(className) > 0 {
		dbVersion.WarnUsageOfNotSupportedClassNamespacedEndpointsForReferences()
	}
	return fmt.Sprintf("/objects/%v/references/%v", id, referenceProperty)
}
