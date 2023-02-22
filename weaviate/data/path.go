package data

import (
	"fmt"
	"net/url"
	"path"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/util"
)

type pathComponents struct {
	id                string
	class             string
	dbVersion         *util.DBVersionSupport
	consistencyLevel  string
	referenceProperty string
}

func buildObjectsGetPath(comp pathComponents) string {
	return buildObjectsPath(comp.id, comp.class, comp.dbVersion, "get")
}

func buildObjectsDeletePath(comp pathComponents) string {
	path := buildObjectsPath(comp.id, comp.class, comp.dbVersion, "delete")
	return appendURLParams(path, comp)
}

func buildObjectsCheckPath(comp pathComponents) string {
	return buildObjectsPath(comp.id, comp.class, comp.dbVersion, "check")
}

func buildObjectsUpdatePath(comp pathComponents) string {
	path := buildObjectsPath(comp.id, comp.class, comp.dbVersion, "update")
	return appendURLParams(path, comp)
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

// func buildReferencesPath(id, className, referenceProperty string, dbVersion *util.DBVersionSupport) string {
func buildReferencesPath(comp pathComponents) string {
	if comp.dbVersion.SupportsClassNameNamespacedEndpoints() {
		if len(comp.class) > 0 {
			path := fmt.Sprintf("/objects/%v/%v/references/%v", comp.class, comp.id, comp.referenceProperty)
			return appendURLParams(path, comp)
		} else {
			comp.dbVersion.WarnDeprecatedNonClassNameNamespacedEndpointsForReferences()
		}
	} else if len(comp.class) > 0 {
		comp.dbVersion.WarnUsageOfNotSupportedClassNamespacedEndpointsForReferences()
	}
	path := fmt.Sprintf("/objects/%v/references/%v", comp.id, comp.referenceProperty)
	return appendURLParams(path, comp)
}

func appendURLParams(path string, comp pathComponents) string {
	u := url.URL{Path: path}
	if comp.consistencyLevel != "" {
		params := url.Values{}
		params.Add("consistency_level", comp.consistencyLevel)
		u.RawQuery = params.Encode()
	}
	return u.String()
}
