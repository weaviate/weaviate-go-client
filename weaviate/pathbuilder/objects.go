package pathbuilder

import (
	"fmt"
	"net/url"
	"path"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/db"
)

func ObjectsGet(comp Components) string {
	return objectsPath(comp.ID, comp.Class, comp.DBVersion, "get")
}

func ObjectsDelete(comp Components) string {
	path := objectsPath(comp.ID, comp.Class, comp.DBVersion, "delete")
	return appendURLParams(path, comp)
}

func ObjectsCheck(comp Components) string {
	return objectsPath(comp.ID, comp.Class, comp.DBVersion, "check")
}

func ObjectsUpdate(comp Components) string {
	path := objectsPath(comp.ID, comp.Class, comp.DBVersion, "update")
	return appendURLParams(path, comp)
}

func objectsPath(id, className string, dbVersion *db.VersionSupport, action string) string {
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

func References(comp Components) string {
	if comp.DBVersion.SupportsClassNameNamespacedEndpoints() {
		if len(comp.Class) > 0 {
			path := fmt.Sprintf("/objects/%v/%v/references/%v", comp.Class, comp.ID, comp.ReferenceProperty)
			return appendURLParams(path, comp)
		} else {
			comp.DBVersion.WarnDeprecatedNonClassNameNamespacedEndpointsForReferences()
		}
	} else if len(comp.Class) > 0 {
		comp.DBVersion.WarnUsageOfNotSupportedClassNamespacedEndpointsForReferences()
	}
	path := fmt.Sprintf("/objects/%v/references/%v", comp.ID, comp.ReferenceProperty)
	return appendURLParams(path, comp)
}

func appendURLParams(path string, comp Components) string {
	u := url.URL{Path: path}
	params := url.Values{}

	if comp.ConsistencyLevel != "" {
		params.Add("consistency_level", comp.ConsistencyLevel)
	}
	if comp.Tenant != "" {
		params.Add("tenant", comp.Tenant)
	}
	if len(params) > 0 {
		u.RawQuery = params.Encode()
	}

	return u.String()
}
